import request from '@/utils/request'

// 获取Kubernetes集群配置列表
export function getKubernetesList(params: any) {
  return request({
    url: '/api/v1/infrastructure/kubernetes',
    method: 'get',
    params
  })
}

// 创建Kubernetes集群配置
export function createKubernetes(data: any) {
  return request({
    url: '/api/v1/infrastructure/kubernetes',
    method: 'post',
    data
  })
}

// 获取Kubernetes集群配置详情
export function getKubernetesDetail(id: number) {
  return request({
    url: `/api/v1/infrastructure/kubernetes/${id}`,
    method: 'get'
  })
}

// 更新Kubernetes集群配置
export function updateKubernetes(id: number, data: any) {
  return request({
    url: `/api/v1/infrastructure/kubernetes/${id}`,
    method: 'put',
    data
  })
}

// 删除Kubernetes集群配置
export function deleteKubernetes(id: number) {
  return request({
    url: `/api/v1/infrastructure/kubernetes/${id}`,
    method: 'delete'
  })
}

// 测试Kubernetes集群连接
export function testKubernetesConnection(data: any) {
  return request({
    url: '/api/v1/infrastructure/kubernetes/test',
    method: 'post',
    data
  })
}

// 获取Kubernetes集群工作负载
export function getKubernetesWorkloads(id: number) {
  return request({
    url: `/api/infrastructure/kubernetes/${id}/workloads`,
    method: 'get'
  })
} 