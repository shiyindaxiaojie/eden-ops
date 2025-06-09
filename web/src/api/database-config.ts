import request from '@/utils/request'
import type { DatabaseConfig, PageQuery, PageResult, BaseResponse } from '@/types/api'

export function getDatabaseConfigs(params: PageQuery) {
  return request<BaseResponse<PageResult<DatabaseConfig>>>({
    url: '/api/v1/database-configs',
    method: 'get',
    params
  })
}

export function getDatabaseConfig(id: number) {
  return request<BaseResponse<DatabaseConfig>>({
    url: `/api/v1/database-configs/${id}`,
    method: 'get'
  })
}

export function createDatabaseConfig(data: Partial<DatabaseConfig>) {
  return request<BaseResponse<DatabaseConfig>>({
    url: '/api/v1/database-configs',
    method: 'post',
    data
  })
}

export function updateDatabaseConfig(id: number, data: Partial<DatabaseConfig>) {
  return request<BaseResponse<DatabaseConfig>>({
    url: `/api/v1/database-configs/${id}`,
    method: 'put',
    data
  })
}

export function deleteDatabaseConfig(id: number) {
  return request<BaseResponse<null>>({
    url: `/api/v1/database-configs/${id}`,
    method: 'delete'
  })
} 