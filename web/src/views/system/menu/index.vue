<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h2 class="title">菜单管理</h2>
            <span class="subtitle">管理系统的菜单结构和权限</span>
          </div>
          <div class="header-right">
            <el-button type="primary" @click="handleAdd">新增菜单</el-button>
          </div>
        </div>
      </template>

      <el-table
        :data="menuList"
        style="width: 100%"
        v-loading="loading"
        row-key="id"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
      >
        <el-table-column prop="name" label="菜单名称" min-width="180" />
        <el-table-column prop="icon" label="图标" width="100">
          <template #default="{ row }">
            <el-icon v-if="row.icon">
              <component :is="row.icon" />
            </el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路由地址" min-width="180" />
        <el-table-column prop="component" label="组件路径" min-width="180" />
        <el-table-column prop="permission" label="权限标识" min-width="150" />
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === '1' ? 'success' : 'danger'">
              {{ row.status === '1' ? '显示' : '隐藏' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleAdd(row)">新增</el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 菜单表单对话框 -->
      <el-dialog
        :title="dialogTitle"
        v-model="dialogVisible"
        width="650px"
        append-to-body
        @close="resetForm"
      >
        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-width="100px"
          class="menu-form"
        >
          <el-form-item label="上级菜单">
            <el-tree-select
              v-model="form.parent_id"
              :data="menuOptions"
              :props="{ label: 'name', children: 'children' }"
              value-key="id"
              placeholder="选择上级菜单"
              check-strictly
              clearable
              class="form-select"
            />
          </el-form-item>
          <el-form-item label="菜单类型" prop="type">
            <el-radio-group v-model="form.type">
              <el-radio label="M">目录</el-radio>
              <el-radio label="C">菜单</el-radio>
              <el-radio label="F">按钮</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="菜单名称" prop="name">
            <el-input v-model="form.name" placeholder="请输入菜单名称" class="form-input" />
          </el-form-item>
          <el-form-item label="图标" prop="icon" v-if="form.type !== 'F'">
            <el-input v-model="form.icon" placeholder="请输入图标名称" class="form-input" />
          </el-form-item>
          <el-form-item label="路由地址" prop="path" v-if="form.type !== 'F'">
            <el-input v-model="form.path" placeholder="请输入路由地址" class="form-input" />
          </el-form-item>
          <el-form-item label="组件路径" prop="component" v-if="form.type === 'C'">
            <el-input v-model="form.component" placeholder="请输入组件路径" class="form-input" />
          </el-form-item>
          <el-form-item label="权限标识" prop="permission" v-if="form.type === 'F'">
            <el-input v-model="form.permission" placeholder="请输入权限标识" class="form-input" />
          </el-form-item>
          <el-form-item label="显示排序" prop="sort">
            <el-input-number v-model="form.sort" :min="0" :max="999" controls-position="right" class="form-input" />
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-radio-group v-model="form.status">
              <el-radio label="1">显示</el-radio>
              <el-radio label="0">隐藏</el-radio>
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
  getMenuTree,
  getMenu,
  createMenu,
  updateMenu,
  deleteMenu
} from '@/api/menu'
import type { Menu } from '@/types/api'

interface MenuForm extends Partial<Menu> {
  parent_id?: number
}

const loading = ref(false)
const menuList = ref<Menu[]>([])
const menuOptions = ref<Menu[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref('')

const formRef = ref<FormInstance>()

const form = reactive<MenuForm>({
  parent_id: undefined,
  type: 'M',
  name: '',
  icon: '',
  path: '',
  component: '',
  permission: '',
  sort: 0,
  status: '1'
})

const rules: FormRules = {
  type: [
    { required: true, message: '请选择菜单类型', trigger: 'change' }
  ],
  name: [
    { required: true, message: '请输入菜单名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  path: [
    { required: true, message: '请输入路由地址', trigger: 'blur' }
  ],
  component: [
    { required: true, message: '请输入组件路径', trigger: 'blur' }
  ],
  permission: [
    { required: true, message: '请输入权限标识', trigger: 'blur' }
  ],
  sort: [
    { required: true, message: '请输入显示排序', trigger: 'blur' }
  ]
}

const loadMenuData = async () => {
  loading.value = true
  try {
    const { data } = await getMenuTree()
    menuList.value = data || []
    menuOptions.value = [{ id: 0, name: '主目录', children: data || [] }]
  } catch (error: any) {
    menuList.value = []
    menuOptions.value = [{ id: 0, name: '主目录', children: [] }]
    ElMessage.error(error.message || '获取菜单列表失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = (row?: Menu) => {
  resetForm()
  if (row) {
    form.parent_id = row.id
  }
  dialogTitle.value = '新增菜单'
  dialogVisible.value = true
}

const handleEdit = async (row: Menu) => {
  resetForm()
  dialogTitle.value = '编辑菜单'
  try {
    const { data } = await getMenu(row.id)
    Object.assign(form, data)
    dialogVisible.value = true
  } catch (error: any) {
    ElMessage.error(error.message || '获取菜单详情失败')
  }
}

const handleDelete = (row: Menu) => {
  if (!row || !row.id || !row.name) {
    ElMessage.warning('无效的菜单数据')
    return
  }

  ElMessageBox.confirm(
    `确认删除菜单"${row.name}"吗？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await deleteMenu(row.id)
      ElMessage.success('删除成功')
      loadMenuData()
    } catch (error: any) {
      ElMessage.error(error.message || '删除失败')
    }
  }).catch(() => {
    ElMessage.info('已取消删除')
  })
}

const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  Object.assign(form, {
    id: undefined,
    parent_id: undefined,
    type: 'M',
    name: '',
    icon: '',
    path: '',
    component: '',
    permission: '',
    sort: 0,
    status: '1'
  })
}

const submitForm = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const submitData = { ...form }
        if (submitData.type === 'M') {
          submitData.component = ''
          submitData.permission = ''
        } else if (submitData.type === 'C') {
          submitData.permission = ''
        } else if (submitData.type === 'F') {
          submitData.path = ''
          submitData.component = ''
          submitData.icon = ''
        }

        if (submitData.id) {
          await updateMenu(submitData.id, submitData)
          ElMessage.success('更新成功')
        } else {
          await createMenu(submitData)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        loadMenuData()
      } catch (error: any) {
        ElMessage.error(error.message || (form.id ? '更新失败' : '创建失败'))
      }
    }
  })
}

onMounted(() => {
  loadMenuData()
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

.menu-form {
  padding: 20px;
}

.form-input,
.form-select {
  width: 100%;
  min-width: 240px;
}

.dialog-footer {
  text-align: right;
  padding-top: 20px;
}

:deep(.el-tree-select) {
  width: 100%;
  min-width: 240px;
}

:deep(.el-select),
:deep(.el-cascader) {
  width: 100%;
  min-width: 240px;
}

:deep(.el-input-number) {
  width: 240px;
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
</style> 