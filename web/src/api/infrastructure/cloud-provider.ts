import request from '@/utils/request'

// 获取云厂商列表
export function getCloudProviderList(params: any) {
  return request({
    url: '/api/infrastructure/cloud-providers',
    method: 'get',
    params
  })
}

// 创建云厂商
export function createCloudProvider(data: any) {
  return request({
    url: '/api/infrastructure/cloud-providers',
    method: 'post',
    data
  })
}

// 获取云厂商详情
export function getCloudProviderDetail(id: number) {
  return request({
    url: `/api/infrastructure/cloud-providers/${id}`,
    method: 'get'
  })
}

// 更新云厂商
export function updateCloudProvider(id: number, data: any) {
  return request({
    url: `/api/infrastructure/cloud-providers/${id}`,
    method: 'put',
    data
  })
}

// 删除云厂商
export function deleteCloudProvider(id: number) {
  return request({
    url: `/api/infrastructure/cloud-providers/${id}`,
    method: 'delete'
  })
} 