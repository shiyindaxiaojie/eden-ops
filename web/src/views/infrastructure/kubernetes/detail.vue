<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h2 class="title">{{ config?.name }}</h2>
            <span class="subtitle">集群工作负载</span>
          </div>
          <div class="header-right">
            <el-button @click="goBack">返回</el-button>
          </div>
        </div>
      </template>

      <!-- 集群信息 -->
      <el-descriptions class="cluster-info" :column="4" border>
        <el-descriptions-item label="Context">{{ config?.context || '-' }}</el-descriptions-item>
        <el-descriptions-item label="云厂商">{{ config?.provider?.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ config?.version || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(config?.status)">{{ getStatusText(config?.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="节点数">{{ config?.node_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="Pod数">{{ config?.pod_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="CPU">
          {{ config ? `${config.cpu_used.toFixed(1)}/${config.cpu_total.toFixed(1)} Core` : '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="内存">
          {{ config ? `${formatMemory(config.memory_used)}/${formatMemory(config.memory_total)}` : '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="同步状态">
          <el-switch
            v-model="config.sync_enabled"
            :active-value="'1'"
            :inactive-value="'0'"
            @change="handleSyncChange"
          />
        </el-descriptions-item>
        <el-descriptions-item label="同步间隔">{{ config?.sync_interval || '-' }} 秒</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ config ? formatDateTime(config.created_at) : '-' }}</el-descriptions-item>
        <el-descriptions-item label="最后同步">{{ config?.last_sync_time ? formatDateTime(config.last_sync_time) : '-' }}</el-descriptions-item>
      </el-descriptions>

      <!-- 工作负载 -->
      <el-tabs v-model="activeTab" class="workload-tabs">
        <el-tab-pane label="Deployment" name="Deployment">
          <el-table :data="getWorkloadsByKind('Deployment')" v-loading="loading">
            <el-table-column label="名称" prop="name" min-width="200" show-overflow-tooltip />
            <el-table-column label="命名空间" prop="namespace" min-width="120" show-overflow-tooltip />
            <el-table-column label="副本数" min-width="100" align="center">
              <template #default="{ row }">
                {{ row.ready_replicas }}/{{ row.replicas }}
              </template>
            </el-table-column>
            <el-table-column label="状态" prop="status" min-width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getWorkloadStatusType(row.status)">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="StatefulSet" name="StatefulSet">
          <el-table :data="getWorkloadsByKind('StatefulSet')" v-loading="loading">
            <el-table-column label="名称" prop="name" min-width="200" show-overflow-tooltip />
            <el-table-column label="命名空间" prop="namespace" min-width="120" show-overflow-tooltip />
            <el-table-column label="副本数" min-width="100" align="center">
              <template #default="{ row }">
                {{ row.ready_replicas }}/{{ row.replicas }}
              </template>
            </el-table-column>
            <el-table-column label="状态" prop="status" min-width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getWorkloadStatusType(row.status)">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="DaemonSet" name="DaemonSet">
          <el-table :data="getWorkloadsByKind('DaemonSet')" v-loading="loading">
            <el-table-column label="名称" prop="name" min-width="200" show-overflow-tooltip />
            <el-table-column label="命名空间" prop="namespace" min-width="120" show-overflow-tooltip />
            <el-table-column label="就绪数" min-width="100" align="center">
              <template #default="{ row }">
                {{ row.ready_replicas }}/{{ row.replicas }}
              </template>
            </el-table-column>
            <el-table-column label="状态" prop="status" min-width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getWorkloadStatusType(row.status)">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="Job" name="Job">
          <el-table :data="getWorkloadsByKind('Job')" v-loading="loading">
            <el-table-column label="名称" prop="name" min-width="200" show-overflow-tooltip />
            <el-table-column label="命名空间" prop="namespace" min-width="120" show-overflow-tooltip />
            <el-table-column label="完成数" min-width="100" align="center">
              <template #default="{ row }">
                {{ row.ready_replicas }}/{{ row.replicas }}
              </template>
            </el-table-column>
            <el-table-column label="状态" prop="status" min-width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getWorkloadStatusType(row.status)">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="CronJob" name="CronJob">
          <el-table :data="getWorkloadsByKind('CronJob')" v-loading="loading">
            <el-table-column label="名称" prop="name" min-width="200" show-overflow-tooltip />
            <el-table-column label="命名空间" prop="namespace" min-width="120" show-overflow-tooltip />
            <el-table-column label="状态" prop="status" min-width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getWorkloadStatusType(row.status)">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getK8sConfig, getK8sWorkloads, updateK8sConfig } from '@/api/k8s-config'
import type { K8sConfig, K8sWorkload } from '@/types/api'
import { formatDateTime } from '@/utils/format'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const config = ref<K8sConfig>()
const workloads = ref<K8sWorkload[]>([])
const activeTab = ref('Deployment')

// 获取集群信息
const getClusterInfo = async () => {
  const id = parseInt(route.params.id as string)
  if (!id) {
    ElMessage.error('无效的集群ID')
    router.push('/infrastructure/kubernetes')
    return
  }

  loading.value = true
  try {
    const [configRes, workloadsRes] = await Promise.all([
      getK8sConfig(id),
      getK8sWorkloads(id)
    ])
    config.value = configRes.data
    workloads.value = workloadsRes.data
  } catch (err: any) {
    ElMessage.error(err.message || '获取集群信息失败')
    router.push('/infrastructure/kubernetes')
  } finally {
    loading.value = false
  }
}

// 按类型获取工作负载
const getWorkloadsByKind = (kind: string) => {
  return workloads.value.filter(w => w.kind === kind)
}

// 获取集群状态类型
const getStatusType = (status?: string) => {
  switch (status) {
    case 'running':
      return 'success'
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

// 获取集群状态文本
const getStatusText = (status?: string) => {
  switch (status) {
    case 'running':
      return '运行中'
    case 'error':
      return '错误'
    default:
      return '未连接'
  }
}

// 获取工作负载状态类型
const getWorkloadStatusType = (status: string) => {
  switch (status) {
    case 'Running':
      return 'success'
    case 'Failed':
      return 'danger'
    case 'Completed':
      return 'success'
    case 'Scheduled':
      return 'success'
    default:
      return 'warning'
  }
}

// 返回
const goBack = () => {
  router.push('/infrastructure/kubernetes')
}

// 格式化内存大小
const formatMemory = (bytes: number) => {
  if (bytes === 0) return '0 GB'
  const gb = bytes / (1024 * 1024 * 1024)
  return gb.toFixed(2) + ' GB'
}

// 处理同步状态变更
const handleSyncChange = async () => {
  if (!config.value) return
  try {
    await updateK8sConfig(config.value.id, { sync_enabled: config.value.sync_enabled })
    ElMessage.success('更新成功')
  } catch (error: any) {
    config.value.sync_enabled = config.value.sync_enabled === '1' ? '0' : '1'
    ElMessage.error(error.message || '更新失败')
  }
}

onMounted(() => {
  getClusterInfo()
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

    .cluster-info {
      margin-bottom: 20px;
    }

    .workload-tabs {
      margin-top: 20px;
    }
  }
}
</style> 