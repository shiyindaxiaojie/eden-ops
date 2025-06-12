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

    // 检查HTTP状态码，200表示成功
    if (response.status === 200) {
      return response.data
    } else {
      ElMessage.error(response.statusText || '请求失败')
      return Promise.reject(new Error(response.statusText || '请求失败'))
    }
  },
  (error) => {
    // 调试日志
    console.error('响应错误:', error)
    console.error('错误详情:', error.response)

    // 处理HTTP错误状态码
    let message = error.message
    if (error.response) {
      const status = error.response.status
      switch (status) {
        case 401:
          message = '未授权，请重新登录'
          break
        case 403:
          message = '拒绝访问'
          break
        case 404:
          message = '请求错误，未找到该资源'
          break
        case 500:
          message = '服务器内部错误'
          break
        default:
          message = error.response.data?.message || error.message
      }
    }

    ElMessage.error(message)
    return Promise.reject(error)
  }
)

// 请求函数
export default function request<T = any>(config: AxiosRequestConfig): Promise<BaseResponse<T>> {
  return service.request(config)
}