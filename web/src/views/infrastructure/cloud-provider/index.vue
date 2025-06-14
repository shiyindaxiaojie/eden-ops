<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h2 class="title">云厂商管理</h2>
          </div>
          <el-button type="primary" @click="handleAdd">新增云厂商</el-button>
        </div>
      </template>

      <div class="description-box">
        <el-alert
          title="云厂商管理用于 Kubernetes 集群关联，支持阿里云、腾讯云、华为云等主流云厂商"
          type="info"
          :closable="false"
          show-icon
        />
      </div>

      <el-form :inline="true" :model="queryParams" class="search-form">
        <el-form-item label="厂商名称">
          <el-input
            v-model="queryParams.name"
            placeholder="请输入厂商名称"
            clearable
            style="min-width: 150px"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="queryParams.status"
            placeholder="请选择状态"
            clearable
            style="min-width: 150px"
          >
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table
        v-loading="loading"
        :data="tableData"
        style="width: 100%"
      >
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column
          prop="name"
          label="云厂商名称"
          min-width="150"
        />
        <el-table-column
          prop="code"
          label="云厂商代码"
          min-width="150"
        />
        <el-table-column
          prop="description"
          label="描述"
          min-width="200"
        />
        <el-table-column
          label="状态"
          width="80"
          align="center"
        >
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column
          label="操作"
          min-width="150"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              link
              @click="handleDelete(row)"
            >
              删除
            </el-button>
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

    <!-- 新增/编辑对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item
          label="云厂商名称"
          prop="name"
        >
          <el-input v-model="form.name" placeholder="请输入云厂商名称" />
        </el-form-item>
        <el-form-item
          label="云厂商代码"
          prop="code"
        >
          <el-input v-model="form.code" placeholder="请输入云厂商代码" />
        </el-form-item>
        <el-form-item
          label="描述"
          prop="description"
        >
          <el-input
            v-model="form.description"
            type="textarea"
            placeholder="请输入描述"
          />
        </el-form-item>
        <el-form-item
          label="状态"
          prop="status"
        >
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
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
import { getCloudProviderList, createCloudProvider, updateCloudProvider, deleteCloudProvider } from '@/api/infrastructure/cloud-provider'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()

const queryParams = ref({
  page: 1,
  pageSize: 10,
  name: '',
  status: undefined
})

const form = ref({
  id: undefined,
  name: '',
  code: '',
  description: '',
  status: 1
})

const rules = {
  name: [
    { required: true, message: '请输入云厂商名称', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入云厂商代码', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

const getList = async () => {
  loading.value = true
  try {
    const res = await getCloudProviderList(queryParams.value)
    tableData.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('Failed to fetch cloud providers:', error)
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
    status: undefined
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
  dialogTitle.value = '新增云厂商'
  form.value = {
    id: undefined,
    name: '',
    code: '',
    description: '',
    status: 1
  }
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑云厂商'
  form.value = {
    ...row
  }
  dialogVisible.value = true
}

const handleStatusChange = async (row: any) => {
  try {
    const statusText = row.status === 1 ? '启用' : '禁用'
    await updateCloudProvider(row.id, {
      ...row,
      status: row.status
    })
    ElMessage.success(`${statusText}成功`)
    // 不需要重新加载列表，因为状态已经在前端更新了
  } catch (error) {
    console.error('Failed to update status:', error)
    // 如果更新失败，恢复原来的状态
    row.status = row.status === 1 ? 0 : 1
    ElMessage.error('状态更新失败，请重试')
  }
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm('确认删除该云厂商吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteCloudProvider(row.id)
      ElMessage.success('删除成功')
      getList()
    } catch (error) {
      console.error('Failed to delete cloud provider:', error)
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
          await updateCloudProvider(form.value.id, data)
          ElMessage.success('更新成功')
        } else {
          await createCloudProvider(data)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        getList()
      } catch (error) {
        console.error('Failed to submit form:', error)
      }
    }
  })
}

onMounted(() => {
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
}

.description-box {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}
</style> 