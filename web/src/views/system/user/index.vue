<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" @click="handleAdd">新增用户</el-button>
    </div>
        </template>

      <el-form :inline="true" :model="queryParams" class="search-form">
        <el-form-item label="用户名">
          <el-input v-model="queryParams.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="queryParams.phone" placeholder="请输入手机号" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="请选择状态" clearable style="min-width: 240px">
            <el-option label="启用" value="1" />
            <el-option label="禁用" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="userList" style="width: 100%" v-loading="loading">
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="nickname" label="昵称" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="phone" label="手机号" />
        <el-table-column prop="status" label="状态" align="center" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 || row.status === '1' ? 'success' : 'danger'">
              {{ row.status === 1 || row.status === '1' ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handleRole(row)">角色</el-button>
            <el-button type="primary" link @click="handleResetPwd(row)">重置密码</el-button>
            <el-button 
              type="danger" 
              link 
              @click="handleDelete(row)"
              v-if="row.username !== 'admin'"
            >删除</el-button>
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

      <!-- 用户表单对话框 -->
      <el-dialog
        :title="dialogTitle"
        v-model="dialogVisible"
        width="500px"
        append-to-body
        @close="resetForm"
      >
      <el-form
          ref="formRef"
          :model="form"
        :rules="rules"
          label-width="80px"
      >
        <el-form-item label="用户名" prop="username">
            <el-input v-model="form.username" placeholder="请输入用户名" :disabled="form.id !== undefined" />
        </el-form-item>
          <el-form-item label="密码" prop="password" v-if="!form.id">
            <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
            <el-input v-model="form.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
            <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
            <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-radio-group v-model="form.status">
              <el-radio label="1">正常</el-radio>
              <el-radio label="0">禁用</el-radio>
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

      <!-- 分配角色对话框 -->
      <el-dialog
        title="分配角色"
        v-model="roleDialogVisible"
        width="500px"
        append-to-body
      >
      <el-form
          ref="roleFormRef"
        :model="roleForm"
          label-width="80px"
      >
        <el-form-item label="用户名">
          <el-input v-model="roleForm.username" disabled />
        </el-form-item>
        <el-form-item label="角色">
            <el-select
              v-model="roleForm.roleIds"
              multiple
              placeholder="请选择角色"
              style="width: 100%; min-width: 240px"
            >
            <el-option
                v-for="role in roleOptions"
                :key="role.id"
                :label="role.name"
                :value="role.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="roleDialogVisible = false">取 消</el-button>
            <el-button type="primary" @click="submitRoleForm">确 定</el-button>
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
  getUsers,
  createUser,
  updateUser,
  deleteUser,
  getUserRoles,
  assignUserRoles,
  updateUserStatus
} from '@/api/user'
import { getRoles } from '@/api/role'
import type { User, Role } from '@/types/api'

interface QueryParams {
  username: string
  phone: string
  status: string
  pageNum: number
  pageSize: number
}

interface UserForm extends Partial<User> {
  password?: string
}

interface RoleForm {
  userId: number
  username: string
  roleIds: number[]
}

const loading = ref(false)
const total = ref(0)
const userList = ref<User[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const roleDialogVisible = ref(false)
const roleOptions = ref<Role[]>([])

const formRef = ref<FormInstance>()
const roleFormRef = ref<FormInstance>()

const queryParams = reactive<QueryParams>({
  username: '',
  phone: '',
  status: '',
  pageNum: 1,
  pageSize: 10
})

const form = reactive<UserForm>({
        username: '',
        password: '',
        nickname: '',
        email: '',
        phone: '',
        status: '1'
})

const roleForm = reactive<RoleForm>({
  userId: 0,
  username: '',
  roleIds: []
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
        email: [
          { required: true, message: '请输入邮箱', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ],
        phone: [
          { required: true, message: '请输入手机号', trigger: 'blur' },
          { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
        ]
}

const handleQuery = async () => {
  loading.value = true
  try {
    const response = await getUsers({
      ...queryParams,
      page: queryParams.pageNum,
      size: queryParams.pageSize
    })
    
    // 检查响应格式，适配不同的数据结构
    if (response.data && response.data.list) {
      // 标准分页格式
      userList.value = response.data.list
      total.value = response.data.total
    } else if (Array.isArray(response)) {
      // 数组格式
      userList.value = response
      total.value = response.length
    } else if (Array.isArray(response.data)) {
      // 数据直接是数组
      userList.value = response.data
      total.value = response.data.length
    } else {
      // 未知格式，清空数据
      userList.value = []
      total.value = 0
      console.error('未知的响应格式:', response)
    }
  } catch (error: any) {
    console.error('查询用户列表失败:', error)
    ElMessage.error(error.message || '查询失败')
    userList.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const resetQuery = () => {
  queryParams.username = ''
  queryParams.phone = ''
  queryParams.status = ''
  queryParams.pageNum = 1
  handleQuery()
}

const handleAdd = () => {
  resetForm()
  dialogTitle.value = '新增用户'
  dialogVisible.value = true
}

const handleEdit = (row: User) => {
  resetForm()
  dialogTitle.value = '编辑用户'
  const userData = {...row, status: row.status.toString()}
  Object.assign(form, userData)
  dialogVisible.value = true
}

const handleRole = async (row: User) => {
  roleForm.userId = row.id
  roleForm.username = row.username
  roleForm.roleIds = []
  roleDialogVisible.value = true

  try {
    // 获取所有角色
    const rolesResponse = await getRoles()
    // 获取用户角色
    const userRolesResponse = await getUserRoles(row.id)
    
    // 处理角色列表数据
    if (rolesResponse.data && rolesResponse.data.list) {
      roleOptions.value = rolesResponse.data.list
    } else if (Array.isArray(rolesResponse.data)) {
      roleOptions.value = rolesResponse.data
    } else if (Array.isArray(rolesResponse)) {
      roleOptions.value = rolesResponse
    } else {
      roleOptions.value = []
      console.error('获取角色列表格式错误:', rolesResponse)
    }
    
    // 处理用户角色数据
    if (Array.isArray(userRolesResponse)) {
      roleForm.roleIds = userRolesResponse
    } else if (Array.isArray(userRolesResponse.data)) {
      roleForm.roleIds = userRolesResponse.data
    } else {
      roleForm.roleIds = []
      console.error('获取用户角色格式错误:', userRolesResponse)
    }
  } catch (error: any) {
    console.error('获取角色信息失败:', error)
    ElMessage.error(error.message || '获取角色信息失败')
  }
}

const handleResetPwd = (row: User) => {
  ElMessageBox.prompt('请输入新密码', '重置密码', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputType: 'password',
    inputPattern: /^.{6,20}$/,
    inputErrorMessage: '密码长度在 6 到 20 个字符'
  }).then(({ value }) => {
    updateUser(row.id, { password: value }).then(() => {
      ElMessage.success('密码重置成功')
    })
  }).catch(() => {
    ElMessage.info('已取消重置')
  })
}

const handleDelete = (row: User) => {
  ElMessageBox.confirm(
    `确认删除用户"${row.username}"吗？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await deleteUser(row.id)
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
        username: '',
        password: '',
        nickname: '',
        email: '',
        phone: '',
        status: '1'
  })
}

const submitForm = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (form.id) {
          await updateUser(form.id, form)
          ElMessage.success('更新成功')
        } else {
          await createUser(form)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        handleQuery()
      } catch (error: any) {
        console.error('保存用户失败:', error)
        ElMessage.error(error.message || (form.id ? '更新失败' : '创建失败'))
      }
    }
  })
}

const submitRoleForm = async () => {
  try {
    await assignUserRoles(roleForm.userId, roleForm.roleIds)
    ElMessage.success('分配角色成功')
    roleDialogVisible.value = false
  } catch (error: any) {
    console.error('分配角色失败:', error)
    ElMessage.error(error.message || '分配角色失败')
  }
}

onMounted(() => {
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