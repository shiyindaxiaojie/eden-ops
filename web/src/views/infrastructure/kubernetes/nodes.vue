<template>
  <div class="app-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h3 class="title">节点管理</h3>
            <p class="subtitle">集群：{{ clusterName }}</p>
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
        <el-form-item v-if="showStatusFilter" label="状态" prop="status">
          <el-select v-model="queryParams.status" placeholder="请选择状态" clearable style="min-width: 120px;">
            <el-option label="全部" value="" />
            <el-option label="就绪" value="Ready" />
            <el-option label="未就绪" value="NotReady" />
            <el-option label="未知" value="Unknown" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table :data="tableData" v-loading="loading" style="width: 100%; min-width: 1400px;" @sort-change="handleSortChange">
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="name" label="节点名称" width="200" show-overflow-tooltip sortable="custom" />
        <el-table-column prop="internalIP" label="内部IP" width="130" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100" sortable="custom">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="schedulable" label="可调度" width="100">
          <template #default="{ row }">
            <el-tag :type="row.schedulable ? 'success' : 'warning'" size="small">{{ row.schedulable ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="资源容量" width="150">
          <template #default="{ row }">
            <div style="font-size: 12px; line-height: 1.4;">
              <div><span style="color: #409EFF;">CPU：</span>{{ formatCPU(row.cpuCapacity) }}</div>
              <div><span style="color: #67C23A;">内存：</span>{{ formatMemory(row.memoryCapacity) }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Pod 使用" width="120" sortable="custom">
          <template #default="{ row }">
            {{ row.podsUsage }}/{{ row.podsCapacity }}
          </template>
        </el-table-column>
        <el-table-column prop="osImage" label="操作系统" width="120" show-overflow-tooltip />
        <el-table-column prop="kubeletVersion" label="Kubelet 版本" width="130" show-overflow-tooltip />
        <el-table-column prop="createdAt" label="创建时间" width="160" sortable="custom">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
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
    <el-dialog v-model="detailVisible" title="节点详情" width="80%" :before-close="handleCloseDetail">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="节点名称">{{ currentNode?.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="内部 IP">{{ currentNode?.internalIP || '-' }}</el-descriptions-item>
        <el-descriptions-item label="外部 IP">{{ currentNode?.externalIP || '-' }}</el-descriptions-item>
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
        <el-descriptions-item label="操作系统">{{ currentNode?.osImage || '-' }}</el-descriptions-item>
        <el-descriptions-item label="内核版本">{{ currentNode?.kernelVersion || '-' }}</el-descriptions-item>
        <el-descriptions-item label="容器运行时">{{ currentNode?.containerRuntime || '-' }}</el-descriptions-item>
        <el-descriptions-item label="Kubelet 版本">{{ currentNode?.kubeletVersion || '-' }}</el-descriptions-item>
        <el-descriptions-item label="KubeProxy 版本">{{ currentNode?.kubeProxyVersion || '-' }}</el-descriptions-item>
        <el-descriptions-item label="CPU 容量">{{ formatCPU(currentNode?.cpuCapacity) }}</el-descriptions-item>
        <el-descriptions-item label="内存容量">{{ formatMemory(currentNode?.memoryCapacity) }}</el-descriptions-item>
        <el-descriptions-item label="Pod 容量">{{ currentNode?.podsCapacity || '-' }}</el-descriptions-item>
        <el-descriptions-item label="CPU 可分配">{{ formatCPU(currentNode?.cpuAllocatable) }}</el-descriptions-item>
        <el-descriptions-item label="内存可分配">{{ formatMemory(currentNode?.memoryAllocatable) }}</el-descriptions-item>
        <el-descriptions-item label="Pod 可分配">{{ currentNode?.podsAllocatable || '-' }}</el-descriptions-item>
        <el-descriptions-item label="Pod 使用量">{{ currentNode?.podsUsage || 0 }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatTime(currentNode?.createdAt) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getK8sNodes } from '@/api/infrastructure/kubernetes'

const route = useRoute()
const router = useRouter()

// 响应式数据
const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const detailVisible = ref(false)
const currentNode = ref(null)
const clusterName = ref('')
const configId = ref('')
const showStatusFilter = ref(true) // 控制状态过滤器显示

// 查询参数
const queryParams = ref({
  page: 1,
  pageSize: 10,
  configId: null as number | null,
  name: '',
  internalIP: '',
  status: '',
  sortBy: '',
  sortOrder: 'asc'
})

// 获取列表数据
const getList = async () => {
  loading.value = true
  try {
    const response = await getK8sNodes(queryParams.value)
    if (response.code === 0) {
      tableData.value = response.data.list || []
      total.value = response.data.total || 0
    } else {
      ElMessage.error(response.message || '获取节点列表失败')
    }
  } catch (error) {
    console.error('获取节点列表失败:', error)
    ElMessage.error('获取节点列表失败')
  } finally {
    loading.value = false
  }
}

// 查询
const handleQuery = () => {
  queryParams.value.page = 1
  getList()
}

// 重置查询
const resetQuery = () => {
  queryParams.value = {
    page: 1,
    pageSize: 10,
    configId: queryParams.value.configId, // 保持configId
    name: '',
    internalIP: '',
    status: '',
    sortBy: '',
    sortOrder: 'asc'
  }
  getList()
}

// 分页变化
const handlePageChange = (page: number) => {
  queryParams.value.page = page
  getList()
}

const handleSizeChange = (size: number) => {
  queryParams.value.pageSize = size
  queryParams.value.page = 1
  getList()
}

// 排序变化
const handleSortChange = ({ prop, order }: { prop: string; order: string | null }) => {
  if (order) {
    queryParams.value.sortBy = prop
    queryParams.value.sortOrder = order === 'ascending' ? 'asc' : 'desc'
  } else {
    queryParams.value.sortBy = ''
    queryParams.value.sortOrder = 'asc'
  }
  queryParams.value.page = 1
  getList()
}

// 查看详情
const handleViewDetail = (row: any) => {
  currentNode.value = row
  detailVisible.value = true
}

// 返回集群列表
const goBack = () => {
  router.push('/infrastructure/kubernetes')
}

// 关闭详情对话框
const handleCloseDetail = () => {
  detailVisible.value = false
  currentNode.value = null
}

// 格式化状态类型
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

// 格式化状态文本
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



// 格式化CPU
const formatCPU = (cpu: string) => {
  if (!cpu) return '-'
  if (cpu.endsWith('m')) {
    const value = parseInt(cpu.replace('m', ''))
    return `${(value / 1000).toFixed(2)}核`
  }
  return `${cpu}核`
}

// 格式化内存
const formatMemory = (memory: string) => {
  if (!memory) return '-'
  if (memory.endsWith('Ki')) {
    const value = parseInt(memory.replace('Ki', ''))
    return `${(value / 1024 / 1024).toFixed(2)}Gi`
  }
  return memory
}

// 格式化时间
const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

// 初始化
onMounted(() => {
  // 从路由参数获取集群信息
  configId.value = route.query.configId as string
  clusterName.value = route.query.clusterName as string || '未知集群'
  queryParams.value.configId = Number(configId.value)

  // 从URL参数获取状态过滤
  const statusParam = route.query.status as string
  if (statusParam) {
    // 如果有状态参数，说明是从集群管理页面跳转过来的，隐藏状态过滤器
    showStatusFilter.value = false
    queryParams.value.status = statusParam
  }

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
