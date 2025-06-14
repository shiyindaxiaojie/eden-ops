<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h2 class="title">Kubernetes 集群</h2>
          </div>
          <el-button type="primary" @click="handleAdd">接入集群</el-button>
        </div>
      </template>

      <el-form :model="queryParams" ref="queryForm" :inline="true" class="search-form">
        <el-form-item label="集群名称" prop="name">
          <el-input v-model="queryParams.name" placeholder="请输入集群名称" clearable style="min-width: 120px;" />
        </el-form-item>
        <el-form-item label="集群ID" prop="clusterID">
          <el-input v-model="queryParams.clusterID" placeholder="请输入集群ID" clearable style="min-width: 120px;" />
        </el-form-item>
        <el-form-item label="云厂商" prop="providerId">
          <el-select v-model="queryParams.providerId" placeholder="请选择云厂商" clearable style="min-width: 150px;">
            <el-option
              v-for="item in providerOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="同步开关" prop="status">
          <el-select v-model="queryParams.status" placeholder="请选择开关" clearable style="min-width: 120px;">
            <el-option label="全部" value="" />
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="clusterList" v-loading="loading" style="width: 100%; min-width: 1200px;">
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="name" label="集群名称" width="120" show-overflow-tooltip />
        <el-table-column prop="clusterID" label="集群ID" width="120" show-overflow-tooltip />
        <el-table-column prop="providerName" label="云厂商" width="90" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.providerName || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="同步开关" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="工作负载" width="100">
          <template #default="{ row }">
            <div style="font-size: 12px; line-height: 1.4;">
              <div>
                <el-link type="primary" @click="handleViewWorkloads(row, '')">总数 {{ row.workloadCount || 0 }}</el-link>
              </div>
              <div>
                <el-link :type="(row.workloadRunning || 0) > 0 ? 'success' : 'info'" @click="handleViewWorkloads(row, 'gt0')">在线 {{ row.workloadRunning || 0 }}</el-link>
              </div>
              <div>
                <el-link :type="(row.workloadIdle || 0) > 0 ? 'danger' : 'info'" @click="handleViewWorkloads(row, 'eq0')">闲置 {{ row.workloadIdle || 0 }}</el-link>
              </div>
              <div>
                <el-link :type="(row.workloadDestroyedCount || 0) > 0 ? 'warning' : 'info'" @click="handleViewHistory(row, 'workloads')">删除 {{ row.workloadDestroyedCount || 0 }}</el-link>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Pod" width="90">
          <template #default="{ row }">
            <div style="font-size: 12px; line-height: 1.4;">
              <div>
                <el-link type="primary" @click="handleViewPods(row, '')">总数 {{ row.podTotal || 0 }}</el-link>
              </div>
              <div>
                <el-link :type="(row.podRunning || 0) > 0 ? 'success' : 'info'" @click="handleViewPods(row, 'Running')">运行 {{ row.podRunning || 0 }}</el-link>
              </div>
              <div>
                <el-link :type="(row.podError || 0) > 0 ? 'danger' : 'info'" @click="handleViewPods(row, 'Error')">异常 {{ row.podError || 0 }}</el-link>
              </div>
              <div>
                <el-link :type="(row.podDestroyedCount || 0) > 0 ? 'warning' : 'info'" @click="handleViewHistory(row, 'pods')">销毁 {{ row.podDestroyedCount || 0 }}</el-link>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="节点" width="90">
          <template #default="{ row }">
            <div style="font-size: 12px; line-height: 1.4;">
              <div>
                <el-link type="primary" @click="handleViewNodes(row, '')">总数 {{ row.nodeTotal || 0 }}</el-link>
              </div>
              <div>
                <el-link :type="(row.nodeRunning || 0) > 0 ? 'success' : 'info'" @click="handleViewNodes(row, 'Ready')">就绪 {{ row.nodeRunning || 0 }}</el-link>
              </div>
              <div>
                <el-link :type="(row.nodeError || 0) > 0 ? 'danger' : 'info'" @click="handleViewNodes(row, 'NotReady')">异常 {{ row.nodeError || 0 }}</el-link>
              </div>
              <div>
                <el-link :type="(row.nodeDestroyedCount || 0) > 0 ? 'warning' : 'info'" @click="handleViewHistory(row, 'nodes')">驱逐 {{ row.nodeDestroyedCount || 0 }}</el-link>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="CPU" width="120">
          <template #default="{ row }">
            <div v-if="row.cpuUsed && row.cpuTotal">
              <el-progress
                :percentage="getCpuPercentage(row)"
                :stroke-width="8"
                :show-text="false"
                style="margin-bottom: 2px;"
              />
              <div style="font-size: 11px; color: #666;">{{ row.cpuUsed }}/{{ row.cpuTotal }}</div>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="内存" width="120">
          <template #default="{ row }">
            <div v-if="row.memoryUsed && row.memoryTotal">
              <el-progress
                :percentage="getMemoryPercentage(row)"
                :stroke-width="8"
                :show-text="false"
                style="margin-bottom: 2px;"
              />
              <div style="font-size: 11px; color: #666;">{{ row.memoryUsed }}/{{ row.memoryTotal }}</div>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="120" show-overflow-tooltip />
        <el-table-column prop="description" label="描述" width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
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

    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="600px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="集群名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入集群名称" />
        </el-form-item>
        <el-form-item label="集群ID" v-if="form.clusterID">
          <el-input v-model="form.clusterID" readonly placeholder="从Kubeconfig自动解析" />
        </el-form-item>
        <el-form-item label="云厂商" prop="providerId">
          <el-select v-model="form.providerId" placeholder="请选择云厂商" clearable>
            <el-option
              v-for="item in providerOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="Kubeconfig" prop="kubeconfig">
          <el-input
            v-model="form.kubeconfig"
            type="textarea"
            :rows="8"
            placeholder="请输入Kubeconfig配置"
            @input="handleKubeconfigChange"
          />
          <div style="margin-top: 8px;">
            <el-button type="primary" size="small" @click="handleTestConnection" :loading="testLoading">
              测试连接
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            placeholder="请输入描述"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>
        <el-form-item label="同步间隔" prop="syncInterval">
          <el-input-number
            v-model="form.syncInterval"
            :min="30"
            :step="10"
            placeholder="同步间隔（秒）"
            style="width: 120px;"
          />
          <div style="font-size: 12px; color: #999; margin: 4px 4px;">最低30秒，用于定时同步集群状态</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { useRouter } from 'vue-router'
import { getKubernetesList, createKubernetes, updateKubernetes, deleteKubernetes } from '@/api/infrastructure/kubernetes'
import { getCloudProviderList } from '@/api/infrastructure/cloud-provider'

const router = useRouter()

const loading = ref(false)
const clusterList = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()
const providerOptions = ref([])
const testLoading = ref(false)

const queryParams = ref({
  page: 1,
  pageSize: 10,
  name: '',
  providerId: null,
  clusterID: '',
  status: ''
})

const form = ref({
  id: undefined,
  name: '',
  providerId: null,
  kubeconfig: '',
  clusterID: '',
  description: '',
  status: 1,
  syncInterval: 30
})

const rules = {
  name: [
    { required: true, message: '请输入集群名称', trigger: 'blur' }
  ],
  kubeconfig: [
    { required: true, message: '请输入Kubeconfig配置', trigger: 'blur' }
  ]
}

const getProviderOptions = async () => {
  try {
    const res = await getCloudProviderList({ pageSize: 100 })
    providerOptions.value = res.data.list
  } catch (error) {
    console.error('Failed to fetch cloud providers:', error)
  }
}

const getList = async () => {
  loading.value = true
  try {
    const res = await getKubernetesList(queryParams.value)
    clusterList.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('Failed to fetch clusters:', error)
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
    providerId: null,
    clusterID: '',
    status: ''
  }
  getList()
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

const handleAdd = () => {
  dialogTitle.value = '接入集群'
  form.value = {
    id: undefined,
    name: '',
    providerId: null,
    kubeconfig: '',
    clusterID: '',
    description: '',
    status: 1,
    syncInterval: 30
  }
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑集群'
  form.value = {
    ...row,
    status: row.status,
    syncInterval: row.syncInterval || 30
  }
  dialogVisible.value = true
}

const handleViewWorkloads = (row: any, replicasFilter: string = '') => {
  // 路由到工作负载页面，传递集群ID和replicas过滤
  const query: any = {
    configId: row.id,
    clusterName: row.name
  }
  if (replicasFilter) {
    query.replicas = replicasFilter
  }
  router.push({
    path: '/infrastructure/kubernetes/workloads',
    query
  })
}

const handleViewPods = (row: any, status: string = '') => {
  // 路由到Pod页面，传递集群ID和状态过滤
  const query: any = {
    configId: row.id,
    clusterName: row.name
  }
  if (status) {
    query.status = status
  }
  router.push({
    path: '/infrastructure/kubernetes/pods',
    query
  })
}

const handleViewNodes = (row: any, status: string = '') => {
  // 路由到节点页面，传递集群ID和状态过滤
  const query: any = {
    configId: row.id,
    clusterName: row.name
  }
  if (status) {
    query.status = status
  }
  router.push({
    path: '/infrastructure/kubernetes/nodes',
    query
  })
}

const handleViewHistory = (row: any, type: string) => {
  // 路由到历史数据页面
  router.push({
    path: `/infrastructure/kubernetes/history/${row.id}/${type}`,
    query: {
      clusterName: row.name
    }
  })
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm('确认删除该集群配置吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteKubernetes(row.id)
      ElMessage.success('删除成功')
      getList()
    } catch (error) {
      console.error('Failed to delete cluster:', error)
    }
  }).catch(() => {})
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const data = {
          ...form.value
        }
        if (form.value.id) {
          await updateKubernetes(form.value.id, data)
          ElMessage.success('更新成功')
        } else {
          await createKubernetes(data)
          ElMessage.success('接入成功')
        }
        dialogVisible.value = false
        getList()
      } catch (error) {
        console.error('Failed to submit form:', error)
        ElMessage.error('操作失败，请检查输入信息')
      }
    }
  })
}

const handleStatusChange = async (row: any) => {
  try {
    await updateKubernetes(row.id, { ...row, status: row.status })
    ElMessage.success(row.status === 1 ? '已启用' : '已禁用')
    getList()
  } catch (error) {
    console.error('Failed to update status:', error)
    ElMessage.error('状态更新失败')
    // 恢复原状态
    row.status = row.status === 1 ? 0 : 1
  }
}

const handleKubeconfigChange = () => {
  // 尝试解析kubeconfig获取集群ID
  if (form.value.kubeconfig) {
    try {
      const config = parseKubeconfig(form.value.kubeconfig)
      if (config && config.currentContext) {
        // 查找当前上下文对应的集群
        const context = config.contexts?.find(ctx => ctx.name === config.currentContext)
        if (context && context.context && context.context.cluster) {
          form.value.clusterID = context.context.cluster
        }
      }
    } catch (error) {
      // 解析失败时清空集群ID
      form.value.clusterID = ''
    }
  } else {
    form.value.clusterID = ''
  }
}

const parseKubeconfig = (kubeconfigText: string) => {
  try {
    // 简单的YAML解析，实际项目中建议使用专门的YAML解析库
    const lines = kubeconfigText.split('\n')
    const config: any = {}
    let currentSection = ''
    let currentItem: any = {}

    for (const line of lines) {
      const trimmed = line.trim()
      if (trimmed.startsWith('current-context:')) {
        config.currentContext = trimmed.split(':')[1].trim()
      } else if (trimmed === 'contexts:') {
        currentSection = 'contexts'
        config.contexts = []
      } else if (currentSection === 'contexts' && trimmed.startsWith('- context:')) {
        if (Object.keys(currentItem).length > 0) {
          config.contexts.push(currentItem)
        }
        currentItem = { context: {} }
      } else if (currentSection === 'contexts' && trimmed.startsWith('cluster:')) {
        currentItem.context.cluster = trimmed.split(':')[1].trim()
      } else if (currentSection === 'contexts' && trimmed.startsWith('name:')) {
        currentItem.name = trimmed.split(':')[1].trim()
      }
    }

    if (currentSection === 'contexts' && Object.keys(currentItem).length > 0) {
      config.contexts.push(currentItem)
    }

    return config
  } catch (error) {
    return null
  }
}

const handleTestConnection = async () => {
  if (!form.value.kubeconfig) {
    ElMessage.warning('请先输入Kubeconfig配置')
    return
  }

  testLoading.value = true
  try {
    // 这里调用测试连接API
    const response = await fetch('/api/v1/infrastructure/kubernetes/test', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        kubeconfig: form.value.kubeconfig
      })
    })

    if (response.ok) {
      ElMessage.success('连接测试成功')
    } else {
      ElMessage.error('连接测试失败')
    }
  } catch (error) {
    console.error('Failed to test connection:', error)
    ElMessage.error('连接测试失败')
  } finally {
    testLoading.value = false
  }
}

// 计算CPU使用百分比
const getCpuPercentage = (row: any) => {
  if (!row.cpuUsed || !row.cpuTotal) return 0

  // 解析CPU值，支持m（毫核）和k（千核）单位
  const parseValue = (value: string) => {
    if (value.endsWith('m')) {
      return parseInt(value.slice(0, -1))
    } else if (value.endsWith('k')) {
      return parseInt(value.slice(0, -1)) * 1000
    } else {
      return parseInt(value) * 1000 // 默认为核心数，转换为毫核
    }
  }

  const used = parseValue(row.cpuUsed)
  const total = parseValue(row.cpuTotal)

  return Math.round((used / total) * 100)
}

// 计算内存使用百分比
const getMemoryPercentage = (row: any) => {
  if (!row.memoryUsed || !row.memoryTotal) return 0

  // 解析内存值，支持GB、MB、KB、PB等单位
  const parseValue = (value: string) => {
    const units = {
      'KB': 1024,
      'MB': 1024 * 1024,
      'GB': 1024 * 1024 * 1024,
      'TB': 1024 * 1024 * 1024 * 1024,
      'PB': 1024 * 1024 * 1024 * 1024 * 1024
    }

    for (const [unit, multiplier] of Object.entries(units)) {
      if (value.endsWith(unit)) {
        return parseFloat(value.slice(0, -unit.length)) * multiplier
      }
    }

    return parseFloat(value) // 默认为字节
  }

  const used = parseValue(row.memoryUsed)
  const total = parseValue(row.memoryTotal)

  return Math.round((used / total) * 100)
}

onMounted(() => {
  getProviderOptions()
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

  .title {
    font-size: 18px;
    margin: 0;
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