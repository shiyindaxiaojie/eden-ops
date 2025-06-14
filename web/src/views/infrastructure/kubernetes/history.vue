<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h2 class="title">{{ config?.name }} - 历史数据</h2>
            <span class="subtitle">{{ getResourceTypeText(resourceType) }}历史记录</span>
          </div>
          <div class="header-right">
            <el-button @click="goBack">返回</el-button>
          </div>
        </div>
      </template>

      <!-- 查询条件 -->
      <el-form :model="queryForm" :inline="true" class="query-form">
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            @change="handleDateRangeChange"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadHistoryData">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 统计信息 -->
      <el-row :gutter="20" class="stats-row">
        <el-col :span="6">
          <el-statistic title="总记录数" :value="total" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="最近归档" :value="lastArchivedTime" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="同步清理" :value="syncCleanupCount" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="手动归档" :value="manualArchiveCount" />
        </el-col>
      </el-row>

      <!-- 历史数据表格 -->
      <el-table :data="historyData" v-loading="loading" class="history-table">
        <el-table-column label="原始ID" prop="original_id" width="80" />
        <el-table-column label="名称" prop="name" min-width="200" show-overflow-tooltip />
        <el-table-column label="命名空间" prop="namespace" width="120" v-if="resourceType !== 'nodes'" />
        
        <!-- Pod特有字段 -->
        <template v-if="resourceType === 'pods'">
          <el-table-column label="工作负载" prop="workload_name" width="150" show-overflow-tooltip />
          <el-table-column label="状态" prop="status" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="节点" prop="node_name" width="120" show-overflow-tooltip />
        </template>

        <!-- Node特有字段 -->
        <template v-if="resourceType === 'nodes'">
          <el-table-column label="内部IP" prop="internal_ip" width="120" />
          <el-table-column label="状态" prop="status" width="100">
            <template #default="{ row }">
              <el-tag :type="getNodeStatusType(row.status)">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="就绪" prop="ready" width="80">
            <template #default="{ row }">
              <el-tag :type="row.ready ? 'success' : 'danger'">
                {{ row.ready ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
        </template>

        <!-- Workload特有字段 -->
        <template v-if="resourceType === 'workloads'">
          <el-table-column label="类型" prop="kind" width="100" />
          <el-table-column label="副本数" width="100">
            <template #default="{ row }">
              {{ row.ready_replicas }}/{{ row.replicas }}
            </template>
          </el-table-column>
          <el-table-column label="状态" prop="status" width="100">
            <template #default="{ row }">
              <el-tag :type="getWorkloadStatusType(row.status)">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
        </template>

        <el-table-column label="归档原因" prop="archive_reason" width="120">
          <template #default="{ row }">
            <el-tag :type="getArchiveReasonType(row.archive_reason)">
              {{ getArchiveReasonText(row.archive_reason) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="归档时间" prop="archived_at" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.archived_at) }}
          </template>
        </el-table-column>
        <el-table-column label="原始创建时间" prop="created_at" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="queryForm.page"
        v-model:page-size="queryForm.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        class="pagination"
      />
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getK8sConfig } from '@/api/k8s-config'
import { getK8sHistory } from '@/api/k8s-history'
import type { K8sConfig } from '@/types/api'
import { formatDateTime } from '@/utils/format'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const config = ref<K8sConfig>()
const historyData = ref<any[]>([])
const total = ref(0)
const dateRange = ref<[string, string] | null>(null)

// 查询参数
const queryForm = ref({
  page: 1,
  pageSize: 20,
  startTime: '',
  endTime: ''
})

// 统计数据
const syncCleanupCount = ref(0)
const manualArchiveCount = ref(0)
const lastArchivedTime = ref('')

// 从路由获取参数
const configId = computed(() => parseInt(route.params.configId as string))
const resourceType = computed(() => route.params.type as string)

// 获取资源类型文本
const getResourceTypeText = (type: string) => {
  switch (type) {
    case 'pods': return 'Pod'
    case 'nodes': return 'Node'
    case 'workloads': return '工作负载'
    default: return '未知'
  }
}

// 获取集群信息
const getClusterInfo = async () => {
  try {
    const res = await getK8sConfig(configId.value)
    config.value = res.data
  } catch (err: any) {
    ElMessage.error(err.message || '获取集群信息失败')
    router.push('/infrastructure/kubernetes')
  }
}

// 加载历史数据
const loadHistoryData = async () => {
  loading.value = true
  try {
    const params = {
      page: queryForm.value.page,
      pageSize: queryForm.value.pageSize,
      startTime: queryForm.value.startTime,
      endTime: queryForm.value.endTime
    }
    
    const res = await getK8sHistory(configId.value, resourceType.value, params)
    historyData.value = res.data.data
    total.value = res.data.total
    
    // 统计归档原因
    syncCleanupCount.value = historyData.value.filter(item => item.archive_reason === 'sync_cleanup').length
    manualArchiveCount.value = historyData.value.filter(item => item.archive_reason === 'manual').length
    
    // 获取最新归档时间
    if (historyData.value.length > 0) {
      lastArchivedTime.value = formatDateTime(historyData.value[0].archived_at)
    }
  } catch (err: any) {
    ElMessage.error(err.message || '获取历史数据失败')
  } finally {
    loading.value = false
  }
}

// 处理时间范围变化
const handleDateRangeChange = (value: [string, string] | null) => {
  if (value) {
    queryForm.value.startTime = value[0]
    queryForm.value.endTime = value[1]
  } else {
    queryForm.value.startTime = ''
    queryForm.value.endTime = ''
  }
}

// 重置查询
const resetQuery = () => {
  dateRange.value = null
  queryForm.value.startTime = ''
  queryForm.value.endTime = ''
  queryForm.value.page = 1
  loadHistoryData()
}

// 分页处理
const handleSizeChange = (size: number) => {
  queryForm.value.pageSize = size
  queryForm.value.page = 1
  loadHistoryData()
}

const handleCurrentChange = (page: number) => {
  queryForm.value.page = page
  loadHistoryData()
}

// 状态类型
const getStatusType = (status: string) => {
  switch (status) {
    case 'Running': return 'success'
    case 'Failed': return 'danger'
    case 'Pending': return 'warning'
    default: return 'info'
  }
}

const getNodeStatusType = (status: string) => {
  switch (status) {
    case 'Ready': return 'success'
    case 'NotReady': return 'danger'
    default: return 'warning'
  }
}

const getWorkloadStatusType = (status: string) => {
  switch (status) {
    case 'Running': return 'success'
    case 'Failed': return 'danger'
    case 'Completed': return 'success'
    default: return 'warning'
  }
}

const getArchiveReasonType = (reason: string) => {
  switch (reason) {
    case 'sync_cleanup': return 'primary'
    case 'manual': return 'warning'
    case 'expired': return 'info'
    default: return 'info'
  }
}

const getArchiveReasonText = (reason: string) => {
  switch (reason) {
    case 'sync_cleanup': return '同步清理'
    case 'manual': return '手动归档'
    case 'expired': return '过期清理'
    default: return reason
  }
}

// 返回
const goBack = () => {
  router.push('/infrastructure/kubernetes')
}

onMounted(() => {
  getClusterInfo()
  loadHistoryData()
})
</script>

<style lang="scss" scoped>
.app-container {
  padding: 20px;

  .box-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .header-left {
        .title {
          margin: 0;
          font-size: 18px;
          font-weight: 500;
        }

        .subtitle {
          margin-left: 10px;
          font-size: 14px;
          color: #666;
        }
      }
    }

    .query-form {
      margin-bottom: 20px;
    }

    .stats-row {
      margin-bottom: 20px;
    }

    .history-table {
      margin-bottom: 20px;
    }

    .pagination {
      text-align: right;
    }
  }
}
</style>
