import request from '@/utils/request'
import type { Menu, PageQuery, PageResult } from '@/types/api'

// 获取菜单列表
export function getMenus(params?: PageQuery) {
  return request<PageResult<Menu>>({
    url: '/api/v1/menus',
    method: 'get',
    params
  })
}

// 获取菜单树
export function getMenuTree() {
  return request<Menu[]>({
    url: '/api/v1/menus/tree',
    method: 'get'
  })
}

// 获取菜单详情
export function getMenu(id: number) {
  return request<Menu>({
    url: `/api/v1/menus/${id}`,
    method: 'get'
  })
}

// 创建菜单
export function createMenu(data: Partial<Menu>) {
  return request<Menu>({
    url: '/api/v1/menus',
    method: 'post',
    data
  })
}

// 更新菜单
export function updateMenu(id: number, data: Partial<Menu>) {
  return request<Menu>({
    url: `/api/v1/menus/${id}`,
    method: 'put',
    data
  })
}

// 删除菜单
export function deleteMenu(id: number) {
  return request<null>({
    url: `/api/v1/menus/${id}`,
    method: 'delete'
  })
}