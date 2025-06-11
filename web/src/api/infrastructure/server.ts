import request from '@/utils/request'

// 获取服务器配置列表
export function getServerList(params: any) {
  return request({
    url: '/api/v1/infrastructure/server',
    method: 'get',
    params
  })
}

// 创建服务器配置
export function createServer(data: any) {
  return request({
    url: '/api/v1/infrastructure/server',
    method: 'post',
    data
  })
}

// 获取服务器配置详情
export function getServerDetail(id: number) {
  return request({
    url: `/api/v1/infrastructure/server/${id}`,
    method: 'get'
  })
}

// 更新服务器配置
export function updateServer(id: number, data: any) {
  return request({
    url: `/api/infrastructure/server/${id}`,
    method: 'put',
    data
  })
}

// 删除服务器配置
export function deleteServer(id: number) {
  return request({
    url: `/api/infrastructure/server/${id}`,
    method: 'delete'
  })
} 