import request from '@/utils/request'
import type { PageQuery, PageResult, BaseResponse, CloudAccount } from '@/types/api'

// 云账号接口
export interface CloudAccount {
  id: number
  name: string
  provider: string
  accessKey: string
  secretKey: string
  region: string
  description: string
  status: number
}

// 获取云账号列表
export function getCloudAccounts(params: any) {
  return request({
    url: '/api/v1/cloud-accounts',
    method: 'get',
    params
  })
}

// 获取云账号详情
export function getCloudAccount(id: number) {
  return request({
    url: `/api/v1/cloud-accounts/${id}`,
    method: 'get'
  })
}

// 创建云账号
export function createCloudAccount(data: Partial<CloudAccount>) {
  return request({
    url: '/api/v1/cloud-accounts',
    method: 'post',
    data
  })
}

// 更新云账号
export function updateCloudAccount(id: number, data: Partial<CloudAccount>) {
  return request({
    url: `/api/v1/cloud-accounts/${id}`,
    method: 'put',
    data
  })
}

// 删除云账号
export function deleteCloudAccount(id: number) {
  return request({
    url: `/api/v1/cloud-accounts/${id}`,
    method: 'delete'
  })
}

// 测试云账号连接
export function testCloudAccount(id: number) {
  return request({
    url: `/api/v1/cloud-accounts/${id}/test`,
    method: 'post'
  })
}

// 更新云账号状态
export function updateCloudAccountStatus(id: number, status: number) {
  return request<BaseResponse<null>>({
    url: `/api/v1/cloud-accounts/${id}`,
    method: 'put',
    data: { status }
  })
} 