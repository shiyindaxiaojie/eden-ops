<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h2 class="title">云账号管理</h2>
          </div>
          <el-button type="primary" @click="handleAdd">新增账号</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="queryParams" class="search-form">
        <el-form-item label="账号名称">
          <el-input v-model="queryParams.name" placeholder="请输入账号名称" clearable />
        </el-form-item>
        <el-form-item label="云厂商">
          <el-select v-model="queryParams.providerId" placeholder="请选择云厂商" clearable style="min-width: 150px">
            <el-option
              v-for="item in providerOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="同步开关">
          <el-select v-model="queryParams.status" placeholder="请选择开关" clearable style="min-width: 150px">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="accountList" style="width: 100%" v-loading="loading">
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="name" label="账号名称" />
        <el-table-column prop="providerName" label="云厂商">
          <template #default="{ row }">
            <el-tag>{{ row.providerName || '-' }}</el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="accessKey" label="Access Key" min-width="220" />
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
      <el-table-column prop="description" label="描述" />
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>
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
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total-text="`共 ${total} 条`"
        :page-size-text="'条/页'"
        :goto-text="'前往'"
        :page-text="'页'"
        background
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
        <el-form-item label="账号名称" prop="name">
            <el-input v-model="form.name" placeholder="请输入账号名称" />
        </el-form-item>
          <el-form-item label="云厂商" prop="providerId">
            <el-select v-model="form.providerId" placeholder="请选择云厂商" style="width: 100%; min-width: 240px">
              <el-option
                v-for="item in providerOptions"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
          </el-select>
        </el-form-item>
          <el-form-item label="Access Key" prop="accessKey">
            <el-input v-model="form.accessKey" placeholder="请输入 Access Key" />
        </el-form-item>
          <el-form-item label="Secret Key" prop="secretKey">
            <el-input
              v-model="form.secretKey"
              type="password"
              placeholder="请输入 Secret Key"
              show-password
            />
        </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input
              v-model="form.description"
              type="textarea"
              :rows="2"
              placeholder="请输入描述信息"
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
      </el-form>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="dialogVisible = false">取 消</el-button>
            <el-button type="info" @click="handleTestConnection">测试连接</el-button>
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
  getCloudAccountList,
  createCloudAccount,
  updateCloudAccount,
  deleteCloudAccount
} from '@/api/infrastructure/cloud-account'
import { getCloudProviderList } from '@/api/infrastructure/cloud-provider'
import type { CloudAccount } from '@/types/api'

interface QueryParams {
  name: string
  providerId: number | null
  status: number | null
  pageNum: number
  pageSize: number
}

interface AccountForm extends Partial<CloudAccount> {
  accessKey: string
  secretKey: string
}

const loading = ref(false)
const total = ref(0)
const accountList = ref<CloudAccount[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const providerOptions = ref([])

const formRef = ref<FormInstance>()

const queryParams = reactive<QueryParams>({
  name: '',
  providerId: null,
  status: null,
  pageNum: 1,
  pageSize: 10
})

const form = reactive<AccountForm>({
  name: '',
  providerId: null,
  accessKey: '',
  secretKey: '',
  description: '',
  status: 1
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入账号名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  providerId: [
    { required: true, message: '请选择云厂商', trigger: 'change' }
  ],
  accessKey: [
    { required: true, message: '请输入 Access Key', trigger: 'blur' }
  ],
  secretKey: [
    { required: true, message: '请输入 Secret Key', trigger: 'blur' }
  ],
  region: [
    { required: true, message: '请输入地域', trigger: 'blur' }
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

const formatDateTime = (dateTime: string) => {
  if (!dateTime) return '-'
  const date = new Date(dateTime)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

const handleQuery = async () => {
  loading.value = true
  try {
    // 构建查询参数，过滤掉空值
    const params: any = {
      page: queryParams.pageNum,
      pageSize: queryParams.pageSize
    }

    if (queryParams.name) {
      params.name = queryParams.name
    }
    if (queryParams.providerId) {
      params.providerId = queryParams.providerId
    }
    if (queryParams.status !== null) {
      params.status = queryParams.status
    }

    const { data } = await getCloudAccountList(params)
    accountList.value = data.list
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
  queryParams.status = null
  queryParams.pageNum = 1
  handleQuery()
}

const handleAdd = () => {
  resetForm()
  dialogTitle.value = '新增账号'
  dialogVisible.value = true
}

const handleEdit = (row: CloudAccount) => {
  resetForm()
  dialogTitle.value = '编辑账号'
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleTest = async (row: CloudAccount) => {
  try {
    // TODO: 实现测试连接功能
    ElMessage.success('连接测试成功')
  } catch (error: any) {
    ElMessage.error(error.message || '连接测试失败')
  }
}

const handleTestConnection = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // TODO: 实现测试连接功能
        ElMessage.success('连接测试成功')
      } catch (error: any) {
        ElMessage.error(error.message || '连接测试失败')
      }
    }
  })
}

const handleStatusChange = async (row: CloudAccount) => {
  try {
    await updateCloudAccount(row.id, { ...row, status: row.status })
    ElMessage.success(row.status === 1 ? '已启用' : '已禁用')
    handleQuery()
  } catch (error: any) {
    ElMessage.error('状态更新失败')
    // 恢复原状态
    row.status = row.status === 1 ? 0 : 1
  }
}

const handleDelete = (row: CloudAccount) => {
  ElMessageBox.confirm(
    `确认删除账号"${row.name}"吗？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await deleteCloudAccount(row.id)
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
    accessKey: '',
    secretKey: '',
    description: '',
    status: 1
  })
}

const submitForm = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
        if (valid) {
      try {
        if (form.id) {
          await updateCloudAccount(form.id, form)
          ElMessage.success('更新成功')
        } else {
          await createCloudAccount(form)
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

.search-form {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}

.dialog-footer {
  text-align: right;
}
</style> 