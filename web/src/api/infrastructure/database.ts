import request from '@/utils/request'

// 获取数据库配置列表
export function getDatabaseList(params: any) {
  return request({
    url: '/api/infrastructure/database',
    method: 'get',
    params
  })
}

// 创建数据库配置
export function createDatabase(data: any) {
  return request({
    url: '/api/infrastructure/database',
    method: 'post',
    data
  })
}

// 获取数据库配置详情
export function getDatabaseDetail(id: number) {
  return request({
    url: `/api/infrastructure/database/${id}`,
    method: 'get'
  })
}

// 更新数据库配置
export function updateDatabase(id: number, data: any) {
  return request({
    url: `/api/infrastructure/database/${id}`,
    method: 'put',
    data
  })
}

// 删除数据库配置
export function deleteDatabase(id: number) {
  return request({
    url: `/api/infrastructure/database/${id}`,
    method: 'delete'
  })
} 