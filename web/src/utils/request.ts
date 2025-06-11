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
    
    // 调试日志
    console.log('请求配置:', config)
    
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
    // 调试日志
    console.log('响应数据:', response)
    
    const res = response.data
    
    // 检查响应是否存在
    if (!res) {
      ElMessage.error('响应数据为空')
      return Promise.reject(new Error('响应数据为空'))
    }
    
    // 检查响应状态码
    if (res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    
    // 返回整个响应，而不仅仅是data部分
    return res
  },
  (error) => {
    // 调试日志
    console.error('响应错误:', error)
    console.error('错误详情:', error.response)
    
    const message = error.response?.data?.message || error.message || '请求失败'
    ElMessage.error(message)
    return Promise.reject(error)
  }
)

// 请求函数
export default function request<T = any>(config: AxiosRequestConfig): Promise<BaseResponse<T>> {
  return service.request(config)
}