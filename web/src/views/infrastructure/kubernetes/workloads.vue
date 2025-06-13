<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h2 class="title">工作负载</h2>
            <div class="subtitle">集群：{{ clusterName }}</div>
          </div>
          <el-button @click="goBack">返回集群列表</el-button>
        </div>
      </template>

      <el-form :model="queryParams" ref="queryForm" :inline="true" class="search-form">
        <el-form-item label="名称" prop="name">
          <el-input v-model="queryParams.name" placeholder="请输入工作负载名称" clearable />
        </el-form-item>
        <el-form-item label="命名空间" prop="namespace">
          <el-select v-model="queryParams.namespace" placeholder="请选择命名空间" clearable style="min-width: 200px;">
            <el-option
              v-for="item in namespaceOptions"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="类型" prop="workloadType">
          <el-select v-model="queryParams.workloadType" placeholder="请选择类型" clearable style="min-width: 150px;">
            <el-option label="全部" value="" />
            <el-option label="Deployment" value="Deployment" />
            <el-option label="StatefulSet" value="StatefulSet" />
            <el-option label="DaemonSet" value="DaemonSet" />
            <el-option label="Job" value="Job" />
            <el-option label="CronJob" value="CronJob" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="queryParams.status" placeholder="请选择状态" clearable style="min-width: 120px;">
            <el-option label="全部" value="" />
            <el-option label="运行中" value="Running" />
            <el-option label="等待中" value="Pending" />
            <el-option label="失败" value="Failed" />
            <el-option label="错误" value="Error" />
            <el-option label="进行中" value="Progressing" />
            <el-option label="可用" value="Available" />
            <el-option label="完成" value="Complete" />
          </el-select>
        </el-form-item>
        <el-form-item label="创建时间" prop="dateRange">
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

      <el-table :data="workloadList" v-loading="loading" style="width: 100%; min-width: 1400px;" @sort-change="handleSortChange">
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="name" label="工作负载名称" width="180" show-overflow-tooltip sortable="custom" />
        <el-table-column prop="namespace" label="命名空间" width="120" show-overflow-tooltip sortable="custom" />
        <el-table-column prop="kind" label="类型" width="120" />
        <el-table-column label="Pod" width="120" align="center" sortable="custom" sort-by="replicas">
          <template #default="{ row }">
            <span v-if="row.pod_status">{{ row.pod_status }}</span>
            <span v-else-if="row.replicas !== undefined">{{ row.ready_replicas || 0 }}/{{ row.replicas }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="Request/Limits" width="180" sortable="custom" sort-by="cpu_request">
          <template #default="{ row }">
            <div style="font-size: 12px; line-height: 1.4;">
              <div v-if="row.cpu_request_limits">
                <span style="color: #409EFF;">CPU：</span>{{ row.cpu_request_limits }}
              </div>
              <div v-if="row.memory_request_limits">
                <span style="color: #67C23A;">内存：</span>{{ row.memory_request_limits }}
              </div>
              <div v-if="!row.cpu_request_limits && !row.memory_request_limits">-</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center" sortable="custom">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
                <el-table-column prop="images" label="镜像" width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <div v-if="row.images && row.images.length > 0">
              <div v-for="(image, index) in row.images.slice(0, 2)" :key="index" style="font-size: 12px;">
                {{ image }}
              </div>
              <div v-if="row.images.length > 2" style="font-size: 11px; color: #999;">
                +{{ row.images.length - 2 }} 更多...
              </div>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160" align="center" sortable="custom">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
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
import { getWorkloadList } from '@/api/infrastructure/kubernetes'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const workloadList = ref([])
const total = ref(0)
const namespaceOptions = ref([])
const workloadTypeOptions = ref(['Deployment', 'StatefulSet', 'DaemonSet', 'Job', 'CronJob'])
const clusterName = ref('')
const configId = ref('')
const dateRange = ref([])

const queryParams = ref({
  page: 1,
  pageSize: 10,
  name: '',
  namespace: '',
  workloadType: '',
  status: '',
  replicas: '', // 新增replicas过滤参数
  sortBy: '',
  sortOrder: 'asc',
  startTime: '',
  endTime: '',
  configId: ''
})

const getList = async () => {
  loading.value = true
  try {
    const res = await getWorkloadList(queryParams.value)
    workloadList.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('Failed to fetch workloads:', error)
  } finally {
    loading.value = false
  }
}

const getNamespaces = async () => {
  if (!configId.value) return

  try {
    const response = await fetch(`/api/v1/k8s-namespaces?configId=${configId.value}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      }
    })
    const result = await response.json()
    if (result.code === 200) {
      namespaceOptions.value = result.data || []
    }
  } catch (error) {
    console.error('Failed to fetch namespaces:', error)
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
    workloadType: '',
    status: '',
    replicas: '',
    sortBy: '',
    sortOrder: 'asc',
    startTime: '',
    endTime: '',
    configId: configId.value
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
    queryParams.value.sortOrder = 'asc'
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
  // 这里可以打开详情弹窗或跳转到详情页面
  ElMessage.info('工作负载详情功能开发中...')
}

const goBack = () => {
  router.push('/infrastructure/kubernetes')
}

const getStatusType = (status: string) => {
  switch (status) {
    case 'Running':
    case 'Available':
    case 'Complete':
      return 'success'
    case 'Pending':
    case 'Progressing':
      return 'warning'
    case 'Failed':
    case 'Error':
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
    case 'Error':
      return '错误'
    case 'Progressing':
      return '进行中'
    case 'Available':
      return '可用'
    case 'Complete':
      return '完成'
    default:
      return status || '-'
  }
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  // 从路由参数获取集群信息
  configId.value = route.query.configId as string
  clusterName.value = route.query.clusterName as string || '未知集群'
  queryParams.value.configId = configId.value

  // 从URL参数获取状态过滤
  const statusParam = route.query.status as string
  if (statusParam) {
    queryParams.value.status = statusParam
  }

  // 从URL参数获取replicas过滤
  const replicasParam = route.query.replicas as string
  if (replicasParam) {
    queryParams.value.replicas = replicasParam
  }

  if (!configId.value) {
    ElMessage.error('缺少集群配置ID')
    router.push('/infrastructure/kubernetes')
    return
  }

  // 获取命名空间选项和工作负载列表
  getNamespaces()
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
