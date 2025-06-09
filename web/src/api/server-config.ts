import request from '@/utils/request'
import type { ServerConfig, PageQuery, PageResult, BaseResponse } from '@/types/api'

export function getServerConfigs(params: PageQuery) {
  return request<BaseResponse<PageResult<ServerConfig>>>({
    url: '/api/v1/server-configs',
    method: 'get',
    params
  })
}

export function getServerConfig(id: number) {
  return request<BaseResponse<ServerConfig>>({
    url: `/api/v1/server-configs/${id}`,
    method: 'get'
  })
}

export function createServerConfig(data: Partial<ServerConfig>) {
  return request<BaseResponse<ServerConfig>>({
    url: '/api/v1/server-configs',
    method: 'post',
    data
  })
}

export function updateServerConfig(id: number, data: Partial<ServerConfig>) {
  return request<BaseResponse<ServerConfig>>({
    url: `/api/v1/server-configs/${id}`,
    method: 'put',
    data
  })
}

export function deleteServerConfig(id: number) {
  return request<BaseResponse<null>>({
    url: `/api/v1/server-configs/${id}`,
    method: 'delete'
  })
} 