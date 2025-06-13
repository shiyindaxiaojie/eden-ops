<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h2 class="title">服务器配置</h2>
          </div>
          <div class="header-right">
            <el-button type="primary" @click="handleAdd">新增配置</el-button>
          </div>
        </div>
      </template>

      <el-form :inline="true" :model="queryParams" class="search-form">
        <el-form-item label="配置名称">
      <el-input
            v-model="queryParams.name"
            placeholder="请输入配置名称"
            clearable
            style="min-width: 240px"
          />
        </el-form-item>
        <el-form-item label="云厂商">
          <el-select v-model="queryParams.providerId" placeholder="请选择云厂商" clearable style="min-width: 240px;">
            <el-option
              v-for="item in providerOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="服务器类型">
          <el-select
            v-model="queryParams.type"
            placeholder="请选择服务器类型"
            clearable
            style="min-width: 240px"
          >
            <el-option label="Linux" value="linux" />
            <el-option label="Windows" value="windows" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="queryParams.status"
            placeholder="请选择状态"
            clearable
            style="min-width: 240px"
          >
            <el-option label="正常" value="1" />
            <el-option label="异常" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="serverList" style="width: 100%" v-loading="loading">
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="name" label="配置名称" min-width="120" />
        <el-table-column prop="providerName" label="云厂商" min-width="120">
          <template #default="{ row }">
            {{ row.providerName || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="type" label="服务器类型" min-width="120">
          <template #default="{ row }">
            <el-tag>{{ row.type === 'linux' ? 'Linux' : 'Windows' }}</el-tag>
        </template>
      </el-table-column>
        <el-table-column prop="host" label="主机" min-width="150" />
        <el-table-column prop="port" label="端口" width="100" align="center" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === '1' ? 'success' : 'danger'">
              {{ row.status === '1' ? '正常' : '异常' }}
          </el-tag>
        </template>
      </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handleTest(row)">测试连接</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

      <el-pagination
        v-if="total > 0"
        class="pagination"
      :total="total"
        v-model:current-page="queryParams.pageNum"
        v-model:page-size="queryParams.pageSize"
        :page-sizes="[10, 20, 30, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />

      <el-dialog
        :title="dialogTitle"
        v-model="dialogVisible"
        width="600px"
        append-to-body
        @close="resetForm"
      >
      <el-form
          ref="formRef"
          :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="配置名称" prop="name">
            <el-input v-model="form.name" placeholder="请输入配置名称" />
          </el-form-item>
          <el-form-item label="云厂商" prop="providerId">
            <el-select v-model="form.providerId" placeholder="请选择云厂商" clearable style="width: 100%;">
              <el-option
                v-for="item in providerOptions"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="服务器类型" prop="type">
            <el-select
              v-model="form.type"
              placeholder="请选择服务器类型"
              style="width: 100%; min-width: 240px"
            >
              <el-option label="Linux" value="linux" />
              <el-option label="Windows" value="windows" />
            </el-select>
        </el-form-item>
        <el-form-item label="主机" prop="host">
            <el-input v-model="form.host" placeholder="请输入主机地址" />
        </el-form-item>
        <el-form-item label="端口" prop="port">
            <el-input-number
              v-model="form.port"
              :min="1"
              :max="65535"
              controls-position="right"
              style="width: 100%"
            />
        </el-form-item>
        <el-form-item label="用户名" prop="username">
            <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              show-password
            />
        </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-radio-group v-model="form.status">
              <el-radio label="1">正常</el-radio>
              <el-radio label="0">异常</el-radio>
            </el-radio-group>
        </el-form-item>
      </el-form>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="dialogVisible = false">取 消</el-button>
            <el-button type="primary" @click="submitForm">确 定</el-button>
      </div>
        </template>
    </el-dialog>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  getServerConfigs,
  createServerConfig,
  updateServerConfig,
  deleteServerConfig,
  testServerConfig
} from '@/api/server-config'
import { getCloudProviderList } from '@/api/infrastructure/cloud-provider'
import type { ServerConfig } from '@/types/api'

interface QueryParams {
  name: string
  providerId: number | null
  type: string
  status: string
  pageNum: number
  pageSize: number
}

interface ServerForm extends Partial<ServerConfig> {
  password: string
}

const loading = ref(false)
const total = ref(0)
const serverList = ref<ServerConfig[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const providerOptions = ref([])

const formRef = ref<FormInstance>()

const queryParams = reactive<QueryParams>({
  name: '',
  providerId: null,
  type: '',
  status: '',
  pageNum: 1,
  pageSize: 10
})

const form = reactive<ServerForm>({
  name: '',
  providerId: null,
  type: 'linux',
  host: '',
  port: 22,
  username: '',
  password: '',
  status: '1'
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入配置名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择服务器类型', trigger: 'change' }
  ],
  host: [
    { required: true, message: '请输入主机地址', trigger: 'blur' }
  ],
  port: [
    { required: true, message: '请输入端口号', trigger: 'blur' }
  ],
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
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

const handleQuery = async () => {
  loading.value = true
  try {
    const { data } = await getServerConfigs({
      ...queryParams,
      page: queryParams.pageNum,
      size: queryParams.pageSize
    })
    serverList.value = data.list
    total.value = data.total
  } catch (error: any) {
    ElMessage.error(error.message || '查询失败')
  } finally {
    loading.value = false
  }
}

const resetQuery = () => {
  queryParams.name = ''
  queryParams.providerId = null
  queryParams.type = ''
  queryParams.status = ''
  queryParams.pageNum = 1
  handleQuery()
}

const handleAdd = () => {
  resetForm()
  dialogTitle.value = '新增配置'
  dialogVisible.value = true
}

const handleEdit = (row: ServerConfig) => {
  resetForm()
  dialogTitle.value = '编辑配置'
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleTest = async (row: ServerConfig) => {
  try {
    await testServerConfig(row.id)
    ElMessage.success('连接测试成功')
  } catch (error: any) {
    ElMessage.error(error.message || '连接测试失败')
  }
}

const handleDelete = (row: ServerConfig) => {
  ElMessageBox.confirm(
    `确认删除配置"${row.name}"吗？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await deleteServerConfig(row.id)
      ElMessage.success('删除成功')
      handleQuery()
    } catch (error: any) {
      ElMessage.error(error.message || '删除失败')
    }
  }).catch(() => {
    ElMessage.info('已取消删除')
  })
}

const handleSizeChange = (size: number) => {
  queryParams.pageSize = size
  handleQuery()
}

const handleCurrentChange = (page: number) => {
  queryParams.pageNum = page
  handleQuery()
}

const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  Object.assign(form, {
    id: undefined,
    name: '',
    providerId: null,
    type: 'linux',
    host: '',
    port: 22,
    username: '',
    password: '',
    status: '1'
  })
}

const submitForm = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
        if (valid) {
      try {
        if (form.id) {
          await updateServerConfig(form.id, form)
          ElMessage.success('更新成功')
        } else {
          await createServerConfig(form)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        handleQuery()
      } catch (error: any) {
        ElMessage.error(error.message || (form.id ? '更新失败' : '创建失败'))
      }
    }
  })
}

onMounted(() => {
  getProviderOptions()
  handleQuery()
})
</script>

<style scoped>
.app-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title {
  margin: 0;
  font-size: 20px;
  font-weight: 500;
}

.search-form {
  background-color: #f5f7fa;
  padding: 24px;
  border-radius: 4px;
  margin-bottom: 24px;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}

.dialog-footer {
  text-align: right;
}

:deep(.el-input),
:deep(.el-select),
:deep(.el-cascader) {
  width: 100%;
  min-width: 240px;
}
</style> 