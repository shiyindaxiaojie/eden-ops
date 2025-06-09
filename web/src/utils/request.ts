import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { getToken } from '@/utils/auth'
import type { BaseResponse } from '@/types/api'

// 创建 axios 实例
const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 10000
})

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    // 添加 token
    const token = getToken()
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    console.error(error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse<BaseResponse<any>>) => {
    const res = response.data
    if (res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return response
  },
  (error) => {
    console.error(error)
    ElMessage.error(error.message || '请求失败')
    return Promise.reject(error)
  }
)

// 请求函数
export default function request<T>(config: AxiosRequestConfig): Promise<BaseResponse<T>> {
  return service.request<BaseResponse<T>>(config).then(response => response.data)
}