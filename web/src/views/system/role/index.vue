<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>角色管理</span>
          <el-button type="primary" @click="handleAdd">新增角色</el-button>
        </div>
      </template>
      
      <el-form :inline="true" :model="queryParams" class="search-form">
        <el-form-item label="角色名称">
          <el-input v-model="queryParams.name" placeholder="请输入角色名称" clearable />
        </el-form-item>
        <el-form-item label="角色编码">
          <el-input v-model="queryParams.code" placeholder="请输入角色编码" clearable />
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

      <el-table :data="roleList" style="width: 100%" v-loading="loading">
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="name" label="角色名称" />
        <el-table-column prop="code" label="角色编码" />
        <el-table-column prop="remark" label="备注" />
        <el-table-column prop="status" label="状态" align="center" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === '1' ? 'success' : 'danger'">
              {{ row.status === '1' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handlePermission(row)">权限</el-button>
            <el-button 
              type="danger" 
              link 
              @click="handleDelete(row)"
              v-if="row.code !== 'ROLE_ADMIN'"
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

      <!-- 角色表单对话框 -->
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
          <el-form-item label="角色名称" prop="name">
            <el-input v-model="form.name" placeholder="请输入角色名称" />
          </el-form-item>
          <el-form-item label="角色编码" prop="code">
            <el-input v-model="form.code" placeholder="请输入角色编码" :disabled="form.id !== undefined" />
          </el-form-item>
          <el-form-item label="备注" prop="remark">
            <el-input v-model="form.remark" type="textarea" placeholder="请输入备注" />
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-radio-group v-model="form.status">
              <el-radio label="1">启用</el-radio>
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

      <!-- 分配权限对话框 -->
      <el-dialog
        title="分配权限"
        v-model="menuDialogVisible"
        width="500px"
        append-to-body
      >
        <el-form
          ref="menuFormRef"
          :model="menuForm"
          label-width="80px"
        >
          <el-form-item label="角色名称">
            <el-input v-model="menuForm.name" disabled />
          </el-form-item>
          <el-form-item label="菜单权限">
            <el-tree
              ref="menuTreeRef"
              :data="menuTreeData"
              :props="{ label: 'name', children: 'children' }"
              show-checkbox
              node-key="id"
              empty-text="加载中，请稍候"
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="menuDialogVisible = false">取 消</el-button>
            <el-button type="primary" @click="submitMenuForm">确 定</el-button>
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
import type { ElTree } from 'element-plus'
import {
  getRoles,
  createRole,
  updateRole,
  deleteRole,
  getRoleMenus,
  assignRoleMenus
} from '@/api/role'
import { getMenuTree } from '@/api/menu'
import type { Role, Menu } from '@/types/api'

interface QueryParams {
  name: string
  code: string
  status: string
  pageNum: number
  pageSize: number
}

interface RoleForm extends Partial<Role> {}

interface MenuForm {
  roleId: number
  name: string
}

const loading = ref(false)
const total = ref(0)
const roleList = ref<Role[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const menuDialogVisible = ref(false)
const menuTreeData = ref<Menu[]>([])

const formRef = ref<FormInstance>()
const menuFormRef = ref<FormInstance>()
const menuTreeRef = ref<InstanceType<typeof ElTree>>()

const queryParams = reactive<QueryParams>({
  name: '',
  code: '',
  status: '',
  pageNum: 1,
  pageSize: 10
})

const form = reactive<RoleForm>({
  name: '',
  code: '',
  remark: '',
  status: '1'
})

const menuForm = reactive<MenuForm>({
  roleId: 0,
  name: ''
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入角色编码', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ]
}

const handleQuery = async () => {
  loading.value = true
  try {
    const { data } = await getRoles({
      ...queryParams,
      page: queryParams.pageNum,
      size: queryParams.pageSize
    })
    roleList.value = data.list
    total.value = data.total
  } catch (error: any) {
    ElMessage.error(error.message || '查询失败')
  } finally {
    loading.value = false
  }
}

const resetQuery = () => {
  queryParams.name = ''
  queryParams.code = ''
  queryParams.status = ''
  queryParams.pageNum = 1
  handleQuery()
}

const handleAdd = () => {
  resetForm()
  dialogTitle.value = '新增角色'
  dialogVisible.value = true
}

const handleEdit = (row: Role) => {
  resetForm()
  dialogTitle.value = '编辑角色'
  Object.assign(form, row)
  dialogVisible.value = true
}

const handlePermission = async (row: Role) => {
  menuForm.roleId = row.id
  menuForm.name = row.name
  menuDialogVisible.value = true

  try {
    const [menuTreeRes, roleMenusRes] = await Promise.all([
      getMenuTree(),
      getRoleMenus(row.id)
    ])
    menuTreeData.value = menuTreeRes.data
    // 等待下一个 tick，确保树组件已经渲染
    await nextTick()
    if (menuTreeRef.value) {
      menuTreeRef.value.setCheckedKeys(roleMenusRes.data)
    }
  } catch (error: any) {
    ElMessage.error(error.message || '获取菜单权限失败')
  }
}

const handleDelete = (row: Role) => {
  ElMessageBox.confirm(
    `确认删除角色"${row.name}"吗？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await deleteRole(row.id)
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
    code: '',
    remark: '',
    status: '1'
  })
}

const submitForm = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (form.id) {
          await updateRole(form.id, form)
          ElMessage.success('更新成功')
        } else {
          await createRole(form)
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

const submitMenuForm = async () => {
  if (!menuTreeRef.value) return

  try {
    const checkedKeys = menuTreeRef.value.getCheckedKeys(false) as number[]
    const halfCheckedKeys = menuTreeRef.value.getHalfCheckedKeys() as number[]
    await assignRoleMenus(menuForm.roleId, [...checkedKeys, ...halfCheckedKeys])
    ElMessage.success('分配权限成功')
    menuDialogVisible.value = false
  } catch (error: any) {
    ElMessage.error(error.message || '分配权限失败')
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