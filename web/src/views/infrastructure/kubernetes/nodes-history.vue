<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h2 class="title">节点历史</h2>
            <div class="subtitle">集群：{{ clusterName }}</div>
          </div>
          <el-button @click="goBack">返回集群列表</el-button>
        </div>
      </template>

      <el-form :model="queryParams" ref="queryForm" :inline="true" class="search-form">
        <el-form-item label="节点名称" prop="name">
          <el-input v-model="queryParams.name" placeholder="请输入节点名称" clearable />
        </el-form-item>
        <el-form-item label="内部IP" prop="internalIP">
          <el-input v-model="queryParams.internalIP" placeholder="请输入内部IP" clearable />
        </el-form-item>
        <el-form-item label="主机名" prop="hostname">
          <el-input v-model="queryParams.hostname" placeholder="请输入主机名" clearable />
        </el-form-item>
        <el-form-item label="删除原因" prop="archiveReason">
          <el-select v-model="queryParams.archiveReason" placeholder="请选择删除原因" clearable style="min-width: 150px;">
            <el-option label="全部" value="" />
            <el-option label="同步清理" value="sync_cleanup" />
            <el-option label="手动删除" value="manual" />
          </el-select>
        </el-form-item>
        <el-form-item label="删除时间" prop="dateRange">
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="min-width: 350px;"
            @change="handleDateChange"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="nodeList" v-loading="loading" style="width: 100%; min-width: 1600px;" @sort-change="handleSortChange">
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="name" label="节点名称" width="200" show-overflow-tooltip sortable="custom" />
        <el-table-column prop="internal_ip" label="内部IP" width="130" show-overflow-tooltip />
        <el-table-column prop="external_ip" label="外部IP" width="130" show-overflow-tooltip />
        <el-table-column prop="hostname" label="主机名" width="150" show-overflow-tooltip />
        <el-table-column prop="status" label="最后状态" width="100" sortable="custom">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="就绪状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.ready ? 'success' : 'danger'" size="small">{{ row.ready ? '就绪' : '未就绪' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="可调度" width="100">
          <template #default="{ row }">
            <el-tag :type="row.schedulable ? 'success' : 'warning'" size="small">{{ row.schedulable ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="资源容量" width="150">
          <template #default="{ row }">
            <div style="font-size: 12px; line-height: 1.4;">
              <div><span style="color: #409EFF;">CPU：</span>{{ formatCPU(row.cpu_capacity) }}</div>
              <div><span style="color: #67C23A;">内存：</span>{{ formatMemory(row.memory_capacity) }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Pod 使用" width="120" sortable="custom" sort-by="pods_usage">
          <template #default="{ row }">
            {{ row.pods_usage || 0 }}/{{ row.pods_capacity || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="os_image" label="操作系统" width="120" show-overflow-tooltip />
        <el-table-column prop="kubelet_version" label="Kubelet 版本" width="130" show-overflow-tooltip />
        <el-table-column prop="archive_reason" label="删除原因" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getArchiveReasonType(row.archive_reason)" size="small">
              {{ getArchiveReasonText(row.archive_reason) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160" align="center" sortable="custom">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="archived_at" label="删除时间" width="160" align="center" sortable="custom">
          <template #default="{ row }">
            {{ formatTime(row.archived_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-if="total > 0"
        :current-page="queryParams.page"
        :page-size="queryParams.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        class="pagination"
        background
        layout="total, sizes, prev, pager, next, jumper"
        :total-text="`共 ${total} 条`"
        :page-size-text="'条/页'"
        :goto-text="'前往'"
        :page-text="'页'"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </el-card>

    <!-- 节点详情对话框 -->
    <el-dialog v-model="detailVisible" title="节点历史详情" width="80%" :before-close="handleCloseDetail">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="节点名称">{{ currentNode?.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="内部 IP">{{ currentNode?.internal_ip || '-' }}</el-descriptions-item>
        <el-descriptions-item label="外部 IP">{{ currentNode?.external_ip || '-' }}</el-descriptions-item>
        <el-descriptions-item label="主机名">{{ currentNode?.hostname || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentNode?.status)">{{ currentNode?.status || '-' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="就绪状态">
          <el-tag :type="currentNode?.ready ? 'success' : 'danger'">{{ currentNode?.ready ? '就绪' : '未就绪' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="可调度">
          <el-tag :type="currentNode?.schedulable ? 'success' : 'warning'">{{ currentNode?.schedulable ? '是' : '否' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="操作系统">{{ currentNode?.os_image || '-' }}</el-descriptions-item>
        <el-descriptions-item label="内核版本">{{ currentNode?.kernel_version || '-' }}</el-descriptions-item>
        <el-descriptions-item label="容器运行时">{{ currentNode?.container_runtime || '-' }}</el-descriptions-item>
        <el-descriptions-item label="Kubelet 版本">{{ currentNode?.kubelet_version || '-' }}</el-descriptions-item>
        <el-descriptions-item label="KubeProxy 版本">{{ currentNode?.kube_proxy_version || '-' }}</el-descriptions-item>
        <el-descriptions-item label="CPU 容量">{{ formatCPU(currentNode?.cpu_capacity) }}</el-descriptions-item>
        <el-descriptions-item label="内存容量">{{ formatMemory(currentNode?.memory_capacity) }}</el-descriptions-item>
        <el-descriptions-item label="Pod 容量">{{ currentNode?.pods_capacity || '-' }}</el-descriptions-item>
        <el-descriptions-item label="CPU 可分配">{{ formatCPU(currentNode?.cpu_allocatable) }}</el-descriptions-item>
        <el-descriptions-item label="内存可分配">{{ formatMemory(currentNode?.memory_allocatable) }}</el-descriptions-item>
        <el-descriptions-item label="Pod 可分配">{{ currentNode?.pods_allocatable || '-' }}</el-descriptions-item>
        <el-descriptions-item label="Pod 使用量">{{ currentNode?.pods_usage || 0 }}</el-descriptions-item>
        <el-descriptions-item label="删除原因">
          <el-tag :type="getArchiveReasonType(currentNode?.archive_reason)">{{ getArchiveReasonText(currentNode?.archive_reason) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatTime(currentNode?.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="删除时间">{{ formatTime(currentNode?.archived_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getK8sHistory } from '@/api/k8s-history'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const nodeList = ref([])
const total = ref(0)
const detailVisible = ref(false)
const currentNode = ref<any>(null)
const clusterName = ref('')
const configId = ref('')
const dateRange = ref([])

const queryParams = ref({
  page: 1,
  pageSize: 10,
  name: '',
  internalIP: '',
  hostname: '',
  archiveReason: '',
  sortBy: '',
  sortOrder: 'desc',
  startTime: '',
  endTime: ''
})

const getList = async () => {
  loading.value = true
  try {
    const params = {
      page: queryParams.value.page,
      pageSize: queryParams.value.pageSize,
      startTime: queryParams.value.startTime,
      endTime: queryParams.value.endTime
    }
    const res = await getK8sHistory(Number(configId.value), 'nodes', params)
    nodeList.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('Failed to fetch node history:', error)
    ElMessage.error('获取节点历史失败')
  } finally {
    loading.value = false
  }
}

const handleQuery = () => {
  queryParams.value.page = 1
  getList()
}

const resetQuery = () => {
  queryParams.value = {
    page: 1,
    pageSize: 10,
    name: '',
    internalIP: '',
    hostname: '',
    archiveReason: '',
    sortBy: '',
    sortOrder: 'desc',
    startTime: '',
    endTime: ''
  }
  dateRange.value = []
  getList()
}

const handleSortChange = ({ prop, order }: { prop: string; order: string | null }) => {
  if (order) {
    queryParams.value.sortBy = prop
    queryParams.value.sortOrder = order === 'ascending' ? 'asc' : 'desc'
  } else {
    queryParams.value.sortBy = ''
    queryParams.value.sortOrder = 'desc'
  }
  queryParams.value.page = 1
  getList()
}

const handleDateChange = (dates: string[]) => {
  if (dates && dates.length === 2) {
    queryParams.value.startTime = dates[0]
    queryParams.value.endTime = dates[1]
  } else {
    queryParams.value.startTime = ''
    queryParams.value.endTime = ''
  }
}

const handlePageChange = (page: number) => {
  queryParams.value.page = page
  getList()
}

const handleSizeChange = (size: number) => {
  queryParams.value.pageSize = size
  queryParams.value.page = 1
  getList()
}

const handleViewDetail = (row: any) => {
  currentNode.value = row
  detailVisible.value = true
}

const handleCloseDetail = () => {
  detailVisible.value = false
  currentNode.value = null
}

const goBack = () => {
  router.push('/infrastructure/kubernetes')
}

const getStatusType = (status: string) => {
  switch (status) {
    case 'Ready':
      return 'success'
    case 'NotReady':
      return 'danger'
    default:
      return 'warning'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'Ready':
      return '就绪'
    case 'NotReady':
      return '未就绪'
    case 'Unknown':
      return '未知'
    default:
      return status || '-'
  }
}

const getArchiveReasonType = (reason: string) => {
  switch (reason) {
    case 'sync_cleanup':
      return 'warning'
    case 'manual':
      return 'danger'
    default:
      return 'info'
  }
}

const getArchiveReasonText = (reason: string) => {
  switch (reason) {
    case 'sync_cleanup':
      return '同步清理'
    case 'manual':
      return '手动删除'
    default:
      return reason || '-'
  }
}

const formatCPU = (cpu: string) => {
  if (!cpu) return '-'
  if (cpu.endsWith('m')) {
    const value = parseInt(cpu.replace('m', ''))
    return `${(value / 1000).toFixed(2)}核`
  }
  return `${cpu}核`
}

const formatMemory = (memory: string) => {
  if (!memory) return '-'
  if (memory.endsWith('Ki')) {
    const value = parseInt(memory.replace('Ki', ''))
    return `${(value / 1024 / 1024).toFixed(2)}Gi`
  }
  return memory
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  configId.value = route.query.configId as string
  clusterName.value = route.query.clusterName as string || '未知集群'

  if (!configId.value) {
    ElMessage.error('缺少集群配置ID')
    router.push('/infrastructure/kubernetes')
    return
  }

  getList()
})
</script>

<style lang="scss" scoped>
.app-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .header-left {
    .title {
      font-size: 18px;
      margin: 0 0 4px 0;
    }
    
    .subtitle {
      font-size: 14px;
      color: #666;
      margin: 0;
    }
  }
}

.search-form {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}
</style>
