import request from '@/utils/request'
import type {
  K8sConfig,
  K8sDeployment,
  K8sStatefulSet,
  K8sDaemonSet,
  K8sJob,
  K8sCronJob,
  PageQuery,
  PageResult,
  BaseResponse,
  K8sWorkload
} from '@/types/api'

// 获取 K8s 配置列表
export function getK8sConfigs(params?: any) {
  return request<K8sConfig[]>({
    url: '/api/v1/k8s-configs',
    method: 'get',
    params
  })
}

// 获取 K8s 配置详情
export function getK8sConfig(id: number) {
  return request<K8sConfig>({
    url: `/api/v1/k8s-configs/${id}`,
    method: 'get'
  })
}

// 创建 K8s 配置
export function createK8sConfig(data: K8sConfig) {
  return request<K8sConfig>({
    url: '/api/v1/k8s-configs',
    method: 'post',
    data
  })
}

// 更新 K8s 配置
export function updateK8sConfig(id: number, data: K8sConfig) {
  return request({
    url: `/api/v1/k8s-configs/${id}`,
    method: 'put',
    data
  })
}

// 删除 K8s 配置
export function deleteK8sConfig(id: number) {
  return request({
    url: `/api/v1/k8s-configs/${id}`,
    method: 'delete'
  })
}

// 测试 K8s 配置连接
export function testK8sConfig(configContent: string) {
  return request({
    url: '/api/v1/k8s-configs/test',
    method: 'post',
    data: { config_content: configContent }
  })
}

// 获取 Deployment 列表
export function getK8sDeployments(clusterId: string) {
  return request<BaseResponse<K8sDeployment[]>>({
    url: `/api/v1/k8s-configs/${clusterId}/deployments`,
    method: 'get'
  })
}

// 获取 StatefulSet 列表
export function getK8sStatefulSets(clusterId: string) {
  return request<BaseResponse<K8sStatefulSet[]>>({
    url: `/api/v1/k8s-configs/${clusterId}/statefulsets`,
    method: 'get'
  })
}

// 获取 DaemonSet 列表
export function getK8sDaemonSets(clusterId: string) {
  return request<BaseResponse<K8sDaemonSet[]>>({
    url: `/api/v1/k8s-configs/${clusterId}/daemonsets`,
    method: 'get'
  })
}

// 获取 Job 列表
export function getK8sJobs(clusterId: string) {
  return request<BaseResponse<K8sJob[]>>({
    url: `/api/v1/k8s-configs/${clusterId}/jobs`,
    method: 'get'
  })
}

// 获取 CronJob 列表
export function getK8sCronJobs(clusterId: string) {
  return request<BaseResponse<K8sCronJob[]>>({
    url: `/api/v1/k8s-configs/${clusterId}/cronjobs`,
    method: 'get'
  })
}

// 获取工作负载列表
export function getK8sWorkloads(id: number) {
  return request<K8sWorkload[]>({
    url: `/api/v1/k8s-configs/${id}/workloads`,
    method: 'get'
  })
} 