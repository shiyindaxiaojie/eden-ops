import request from '@/utils/request'
import type { User } from '@/types/api'

// 登录请求参数
interface LoginRequest {
  username: string
  password: string
}

// 登录响应数据
interface LoginResponse {
  token: string
  user: User
}

// 登录
export function login(data: LoginRequest) {
  return request<LoginResponse>({
    url: '/api/v1/login',
    method: 'post',
    data
  })
}

// 登出
export function logout() {
  return request<null>({
    url: '/api/v1/logout',
    method: 'post'
  })
}

// 获取用户信息
export function getUserInfo() {
  return request<User>({
    url: '/api/v1/users/info',
    method: 'get'
  })
} 