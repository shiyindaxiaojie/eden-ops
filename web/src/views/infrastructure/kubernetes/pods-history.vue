<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h2 class="title">Pod 历史</h2>
            <div class="subtitle">集群：{{ clusterName }}</div>
          </div>
          <el-button @click="goBack">返回集群列表</el-button>
        </div>
      </template>

      <el-form :model="queryParams" ref="queryForm" :inline="true" class="search-form">
        <el-form-item label="名称" prop="name">
          <el-input v-model="queryParams.name" placeholder="请输入Pod名称" clearable />
        </el-form-item>
        <el-form-item label="命名空间" prop="namespace">
          <el-input v-model="queryParams.namespace" placeholder="请输入命名空间" clearable />
        </el-form-item>
        <el-form-item label="工作负载" prop="workloadName">
          <el-input v-model="queryParams.workloadName" placeholder="请输入工作负载名称" clearable />
        </el-form-item>
        <el-form-item label="节点" prop="nodeName">
          <el-input v-model="queryParams.nodeName" placeholder="请输入节点名称" clearable />
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

      <el-table :data="podList" v-loading="loading" style="width: 100%; min-width: 1600px;" @sort-change="handleSortChange">
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="name" label="Pod 名称" width="250" show-overflow-tooltip sortable="custom" />
        <el-table-column prop="namespace" label="命名空间" width="120" show-overflow-tooltip sortable="custom" />
        <el-table-column prop="workload_name" label="工作负载" width="150" show-overflow-tooltip sortable="custom" />
        <el-table-column prop="workload_kind" label="类型" width="100" />
        <el-table-column prop="status" label="最后状态" width="120" align="center" sortable="custom">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="phase" label="阶段" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getPhaseType(row.phase)" size="small">
              {{ getPhaseText(row.phase) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="pod_ip" label="Pod IP" width="120" />
        <el-table-column prop="host_ip" label="主机 IP" width="120" />
        <el-table-column prop="restart_count" label="重启次数" width="110" sortable="custom" />
        <el-table-column prop="node_name" label="节点" width="150" show-overflow-tooltip />
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
const podList = ref([])
const total = ref(0)
const clusterName = ref('')
const configId = ref('')
const dateRange = ref([])

const queryParams = ref({
  page: 1,
  pageSize: 10,
  name: '',
  namespace: '',
  workloadName: '',
  nodeName: '',
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
    const res = await getK8sHistory(Number(configId.value), 'pods', params)
    podList.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('Failed to fetch pod history:', error)
    ElMessage.error('获取Pod历史失败')
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
    namespace: '',
    workloadName: '',
    nodeName: '',
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

const handleViewDetail = (_row: any) => {
  ElMessage.info('Pod历史详情功能开发中...')
}

const goBack = () => {
  router.push('/infrastructure/kubernetes')
}

const getStatusType = (status: string) => {
  switch (status) {
    case 'Running':
    case 'Succeeded':
      return 'success'
    case 'Pending':
      return 'warning'
    case 'Failed':
    case 'ImagePullBackOff':
    case 'CrashLoopBackOff':
      return 'danger'
    default:
      return 'info'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'Running':
      return '运行中'
    case 'Pending':
      return '等待中'
    case 'Failed':
      return '失败'
    case 'Succeeded':
      return '成功'
    case 'ImagePullBackOff':
      return '镜像拉取失败'
    case 'CrashLoopBackOff':
      return '崩溃循环'
    default:
      return status || '-'
  }
}

const getPhaseType = (phase: string) => {
  switch (phase) {
    case 'Running':
      return 'success'
    case 'Pending':
      return 'warning'
    case 'Failed':
      return 'danger'
    case 'Succeeded':
      return 'success'
    default:
      return 'info'
  }
}

const getPhaseText = (phase: string) => {
  switch (phase) {
    case 'Running':
      return '运行中'
    case 'Pending':
      return '等待中'
    case 'Failed':
      return '失败'
    case 'Succeeded':
      return '成功'
    default:
      return phase || '-'
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
