import request from '@/utils/request'

export function getCloudProviders(params: any) {
  return request({
    url: '/api/v1/cloud-providers',
    method: 'get',
    params
  })
}

export function updateCloudProvider(id: number, data: any) {
  return request({
    url: `/api/v1/cloud-providers/${id}`,
    method: 'put',
    data
  })
} 