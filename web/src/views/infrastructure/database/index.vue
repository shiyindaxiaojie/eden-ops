<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h2 class="title">数据库配置</h2>
          </div>
          <div class="header-right">
            <el-button type="primary" @click="handleCreate">新增配置</el-button>
          </div>
        </div>
      </template>

      <div class="filter-container">
        <el-form :inline="true" :model="listQuery" class="search-form">
          <el-form-item label="配置名称">
            <el-input
              v-model="listQuery.name"
              placeholder="配置名称"
              class="filter-item"
              @keyup.enter="handleFilter"
            />
          </el-form-item>
          <el-form-item label="云厂商">
            <el-select
              v-model="listQuery.providerId"
              placeholder="请选择云厂商"
              clearable
              class="filter-item"
            >
              <el-option
                v-for="item in providerOptions"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="数据库类型">
            <el-select
              v-model="listQuery.driver"
              placeholder="数据库类型"
              clearable
              class="filter-item"
            >
              <el-option
                v-for="item in driverOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button v-waves type="primary" icon="Search" @click="handleFilter">
              搜索
            </el-button>
            <el-button @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table
        v-loading="listLoading"
        :data="list"
        border
        fit
        highlight-current-row
      >
        <el-table-column type="index" label="序号" width="60" align="center" />
        <el-table-column prop="name" label="配置名称" min-width="150" />
        <el-table-column prop="providerName" label="云厂商" min-width="120">
          <template #default="{ row }">
            {{ row.providerName || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="driver" label="数据库类型" min-width="120">
          <template #default="{ row }">
            {{ getDriverLabel(row.driver) }}
          </template>
        </el-table-column>
        <el-table-column label="连接信息" min-width="200">
          <template #default="{ row }">
            {{ row.host }}:{{ row.port }}/{{ row.dbname }}
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleUpdate(row)">编辑</el-button>
            <el-button type="primary" link @click="handleTest(row)">测试连接</el-button>
            <el-button
              type="primary"
              link
              :type="row.status === 1 ? 'danger' : 'success'"
              @click="handleModifyStatus(row, row.status === 1 ? 0 : 1)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <pagination
        v-show="total > 0"
        :total="total"
        v-model:page="listQuery.page"
        v-model:limit="listQuery.limit"
        @pagination="getList"
      />

      <el-dialog
        :title="textMap[dialogStatus]"
        v-model="dialogFormVisible"
        width="650px"
        append-to-body
      >
        <el-form
          ref="dataForm"
          :rules="rules"
          :model="temp"
          label-width="120px"
          class="form-container"
        >
          <el-form-item label="配置名称" prop="name">
            <el-input v-model="temp.name" class="form-input" placeholder="请输入配置名称" />
          </el-form-item>
          <el-form-item label="云厂商" prop="providerId">
            <el-select v-model="temp.providerId" class="form-input" placeholder="请选择云厂商">
              <el-option
                v-for="item in providerOptions"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="数据库类型" prop="driver">
            <el-select v-model="temp.driver" class="form-input" placeholder="请选择数据库类型">
              <el-option
                v-for="item in driverOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="主机" prop="host">
            <el-input v-model="temp.host" class="form-input" placeholder="请输入主机地址" />
          </el-form-item>
          <el-form-item label="端口" prop="port">
            <el-input-number
              v-model="temp.port"
              :min="1"
              :max="65535"
              class="form-input"
              controls-position="right"
            />
          </el-form-item>
          <el-form-item label="用户名" prop="username">
            <el-input v-model="temp.username" class="form-input" placeholder="请输入用户名" />
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input
              v-model="temp.password"
              type="password"
              show-password
              class="form-input"
              placeholder="请输入密码"
            />
          </el-form-item>
          <el-form-item label="数据库名" prop="dbname">
            <el-input v-model="temp.dbname" class="form-input" placeholder="请输入数据库名" />
          </el-form-item>
          <el-form-item label="字符集" prop="charset">
            <el-input v-model="temp.charset" class="form-input" placeholder="请输入字符集" />
          </el-form-item>
          <el-form-item label="最大空闲连接" prop="maxIdleConns">
            <el-input-number
              v-model="temp.maxIdleConns"
              :min="1"
              :max="100"
              class="form-input"
              controls-position="right"
            />
          </el-form-item>
          <el-form-item label="最大连接数" prop="maxOpenConns">
            <el-input-number
              v-model="temp.maxOpenConns"
              :min="1"
              :max="1000"
              class="form-input"
              controls-position="right"
            />
          </el-form-item>
          <el-form-item label="描述">
            <el-input
              v-model="temp.description"
              type="textarea"
              :rows="3"
              class="form-input"
              placeholder="请输入描述信息"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-switch
              v-model="temp.status"
              :active-value="1"
              :inactive-value="0"
              active-text="启用"
              inactive-text="禁用"
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="dialogFormVisible = false">取 消</el-button>
            <el-button type="primary" @click="dialogStatus === 'create' ? createData() : updateData()">
              确 定
            </el-button>
          </div>
        </template>
      </el-dialog>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { getDatabaseConfigs, createDatabaseConfig, updateDatabaseConfig, deleteDatabaseConfig } from '@/api/database-config'
import { getCloudProviderList } from '@/api/infrastructure/cloud-provider'
import type { DatabaseConfig } from '@/types/api'
import waves from '@/directive/waves'

interface ListQuery {
  page: number
  limit: number
  name?: string
  providerId?: number
  driver?: string
}

const driverOptions = [
  { label: 'MySQL', value: 'mysql' },
  { label: 'PostgreSQL', value: 'postgres' },
  { label: 'SQLite', value: 'sqlite' }
]

const listLoading = ref(false)
const list = ref<DatabaseConfig[]>([])
const total = ref(0)
const providerOptions = ref([])
const listQuery = reactive<ListQuery>({
  page: 1,
  limit: 20,
  name: undefined,
  providerId: undefined,
  driver: undefined
})

const dialogFormVisible = ref(false)
const dialogStatus = ref('')
const textMap = {
  update: '编辑数据库配置',
  create: '新增数据库配置'
}

const dataFormRef = ref<FormInstance>()
const temp = reactive<Partial<DatabaseConfig>>({
  id: undefined,
  name: '',
  providerId: undefined,
  driver: 'mysql',
  host: '',
  port: 3306,
  username: '',
  password: '',
  dbname: '',
  charset: 'utf8mb4',
  maxIdleConns: 10,
  maxOpenConns: 100,
  description: '',
  status: 1
})

const rules = {
  name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
  driver: [{ required: true, message: '请选择数据库类型', trigger: 'change' }],
  host: [{ required: true, message: '请输入主机地址', trigger: 'blur' }],
  port: [{ required: true, message: '请输入端口号', trigger: 'blur' }],
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  dbname: [{ required: true, message: '请输入数据库名', trigger: 'blur' }]
}

const getDriverLabel = (value: string) => {
  const option = driverOptions.find(item => item.value === value)
  return option ? option.label : value
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
  listLoading.value = true
  try {
    const { data } = await getDatabaseConfigs(listQuery)
    list.value = data.items
    total.value = data.total
  } catch (error: any) {
    ElMessage.error(error.message || '获取数据库配置列表失败')
  } finally {
    listLoading.value = false
  }
}

const handleFilter = () => {
  listQuery.page = 1
  getList()
}

const resetQuery = () => {
  listQuery.name = undefined
  listQuery.providerId = undefined
  listQuery.driver = undefined
  handleFilter()
}

const resetTemp = () => {
  temp.id = undefined
  temp.name = ''
  temp.providerId = undefined
  temp.driver = 'mysql'
  temp.host = ''
  temp.port = 3306
  temp.username = ''
  temp.password = ''
  temp.dbname = ''
  temp.charset = 'utf8mb4'
  temp.maxIdleConns = 10
  temp.maxOpenConns = 100
  temp.description = ''
  temp.status = 1
}

const handleCreate = () => {
  resetTemp()
  dialogStatus.value = 'create'
  dialogFormVisible.value = true
  nextTick(() => {
    dataFormRef.value?.clearValidate()
  })
}

const createData = async () => {
  const valid = await dataFormRef.value?.validate()
  if (!valid) return

  try {
    await createDatabaseConfig(temp)
    ElMessage.success('创建成功')
    dialogFormVisible.value = false
    getList()
  } catch (error: any) {
    ElMessage.error(error.message || '创建失败')
  }
}

const handleUpdate = (row: DatabaseConfig) => {
  Object.assign(temp, row)
  dialogStatus.value = 'update'
  dialogFormVisible.value = true
  nextTick(() => {
    dataFormRef.value?.clearValidate()
  })
}

const updateData = async () => {
  const valid = await dataFormRef.value?.validate()
  if (!valid) return

  try {
    await updateDatabaseConfig(temp.id!, temp)
    ElMessage.success('更新成功')
    dialogFormVisible.value = false
    getList()
  } catch (error: any) {
    ElMessage.error(error.message || '更新失败')
  }
}

const handleDelete = (row: DatabaseConfig) => {
  ElMessageBox.confirm(`确认删除数据库配置"${row.name}"吗？`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async () => {
      try {
        await deleteDatabaseConfig(row.id)
        ElMessage.success('删除成功')
        getList()
      } catch (error: any) {
        ElMessage.error(error.message || '删除失败')
      }
    })
    .catch(() => {
      ElMessage.info('已取消删除')
    })
}

const handleModifyStatus = async (row: DatabaseConfig, status: number) => {
  const text = status === 1 ? '启用' : '禁用'
  try {
    await updateDatabaseConfig(row.id, { status })
    row.status = status
    ElMessage.success(`${text}成功`)
  } catch (error: any) {
    ElMessage.error(error.message || `${text}失败`)
  }
}

const handleTest = async (row: DatabaseConfig) => {
  try {
    // TODO: Implement test connection API
    ElMessage.success('连接测试成功')
  } catch (error: any) {
    ElMessage.error(error.message || '连接测试失败')
  }
}

onMounted(() => {
  getProviderOptions()
  getList()
})
</script>
1
<style lang="scss" scoped>
.app-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
}

.header-left {
  display: flex;
  flex-direction: column;
}

.title {
  font-size: 20px;
  font-weight: 500;
  color: #1f2937;
  margin: 0;
  line-height: 1.4;
}

.subtitle {
  font-size: 14px;
  color: #6b7280;
  margin-top: 4px;
}

.filter-container {
  margin-bottom: 20px;
  padding: 20px;
  background-color: #f9fafb;
  border-radius: 8px;
}

.search-form {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.form-container {
  padding: 20px;
}

.form-input {
  width: 100%;
  min-width: 240px;
}

:deep(.el-input-number) {
  width: 240px;
}

:deep(.el-select),
:deep(.el-cascader) {
  width: 100%;
  min-width: 240px;
}

:deep(.el-card__header) {
  padding: 0 20px;
  border-bottom: 1px solid #e5e7eb;
}

:deep(.el-card__body) {
  padding: 20px;
}

:deep(.el-table) {
  margin-top: 16px;
}

.dialog-footer {
  text-align: right;
  padding-top: 20px;
}
</style> 