import request from '@/utils/request'
import type { User, PageQuery, PageResult, BaseResponse } from '@/types/api'

// 用户登录接口
interface LoginData {
  username: string
  password: string
}

export function login(data: LoginData) {
  return request({
    url: '/api/v1/login',
    method: 'post',
    data
  })
}

// 获取用户列表
export function getUsers(params?: PageQuery) {
  return request<PageResult<User>>({
    url: '/api/v1/users',
    method: 'get',
    params
  })
}

// 获取用户详情
export function getUser(id: number) {
  return request<User>({
    url: `/api/v1/users/${id}`,
    method: 'get'
  })
}

// 创建用户
export function createUser(data: Partial<User>) {
  return request<User>({
    url: '/api/v1/users',
    method: 'post',
    data
  })
}

// 更新用户
export function updateUser(id: number, data: Partial<User>) {
  return request<User>({
    url: `/api/v1/users/${id}`,
    method: 'put',
    data
  })
}

// 删除用户
export function deleteUser(id: number) {
  return request<null>({
    url: `/api/v1/users/${id}`,
    method: 'delete'
  })
}

// 更新用户状态
export function updateUserStatus(id: number, status: string) {
  return request({
    url: `/api/v1/users/${id}`,
    method: 'put',
    data: { status }
  })
}

// 获取用户信息
export function getUserInfo() {
  return request<User>({
    url: '/api/v1/users/info',
    method: 'get'
  })
}

// 获取用户角色
export function getUserRoles(userId: number) {
  return request<number[]>({
    url: `/api/v1/users/${userId}/roles`,
    method: 'get'
  })
}

// 分配用户角色
export function assignUserRoles(userId: number, roleIds: number[]) {
  return request<null>({
    url: `/api/v1/users/${userId}/roles`,
    method: 'put',
    data: roleIds
  })
}

export function getInfo() {
  return request({
    url: '/api/v1/users/info',
    method: 'get'
  })
}

export function logout() {
  return request({
    url: '/api/v1/logout',
    method: 'post'
  })
}