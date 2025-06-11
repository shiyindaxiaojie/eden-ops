import request from '@/utils/request'

// 获取云账号列表
export function getCloudAccountList(params: any) {
  return request({
    url: '/api/v1/infrastructure/cloud-accounts',
    method: 'get',
    params
  })
}

// 创建云账号
export function createCloudAccount(data: any) {
  return request({
    url: '/api/v1/infrastructure/cloud-accounts',
    method: 'post',
    data
  })
}

// 获取云账号详情
export function getCloudAccountDetail(id: number) {
  return request({
    url: `/api/infrastructure/cloud-accounts/${id}`,
    method: 'get'
  })
}

// 更新云账号
export function updateCloudAccount(id: number, data: any) {
  return request({
    url: `/api/infrastructure/cloud-accounts/${id}`,
    method: 'put',
    data
  })
}

// 删除云账号
export function deleteCloudAccount(id: number) {
  return request({
    url: `/api/infrastructure/cloud-accounts/${id}`,
    method: 'delete'
  })
} 