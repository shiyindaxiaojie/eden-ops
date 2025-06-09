import request from '@/utils/request'
import type { Role, PageQuery, PageResult } from '@/types/api'

// 获取角色列表
export function getRoles(params?: PageQuery) {
  return request<PageResult<Role>>({
    url: '/api/v1/roles',
    method: 'get',
    params
  })
}

// 获取角色详情
export function getRole(id: number) {
  return request<Role>({
    url: `/api/v1/roles/${id}`,
    method: 'get'
  })
}

// 创建角色
export function createRole(data: Partial<Role>) {
  return request<Role>({
    url: '/api/v1/roles',
    method: 'post',
    data
  })
}

// 更新角色
export function updateRole(id: number, data: Partial<Role>) {
  return request<Role>({
    url: `/api/v1/roles/${id}`,
    method: 'put',
    data
  })
}

// 删除角色
export function deleteRole(id: number) {
  return request<null>({
    url: `/api/v1/roles/${id}`,
    method: 'delete'
  })
}

// 获取角色的菜单权限
export function getRoleMenus(roleId: number) {
  return request<number[]>({
    url: `/api/v1/roles/${roleId}/menus`,
    method: 'get'
  })
}

// 分配角色菜单权限
export function assignRoleMenus(roleId: number, menuIds: number[]) {
  return request<null>({
    url: `/api/v1/roles/${roleId}/menus`,
    method: 'put',
    data: { menuIds }
  })
}