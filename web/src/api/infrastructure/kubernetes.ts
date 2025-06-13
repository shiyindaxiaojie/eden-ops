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

// 工作负载相关API
export function getWorkloadList(params: any) {
  return request({
    url: '/api/v1/k8s-workloads',
    method: 'get',
    params
  })
}

export function getWorkloadDetail(id: number) {
  return request({
    url: `/api/v1/k8s-workloads/${id}`,
    method: 'get'
  })
}

// Pod相关API
export function getPodList(params: any) {
  return request({
    url: '/api/v1/k8s-pods',
    method: 'get',
    params
  })
}

export function getPodDetail(id: number) {
  return request({
    url: `/api/v1/k8s-pods/${id}`,
    method: 'get'
  })
}

// 节点相关API
export function getK8sNodes(params: any) {
  return request({
    url: '/api/v1/k8s-nodes',
    method: 'get',
    params
  })
}

export function getK8sNodeDetail(id: number) {
  return request({
    url: `/api/v1/k8s-nodes/${id}`,
    method: 'get'
  })
}