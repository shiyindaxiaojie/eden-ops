import request from '@/utils/request'

// 获取历史数据
export function getK8sHistory(configId: number, type: string, params: any) {
  return request({
    url: `/api/k8s-history/${configId}/${type}`,
    method: 'get',
    params
  })
}

// 获取历史数据统计
export function getK8sHistoryStatistics(configId: number) {
  return request({
    url: `/api/k8s-history/${configId}/statistics`,
    method: 'get'
  })
}

// 手动清理历史数据
export function cleanupK8sHistory(beforeDate: string) {
  return request({
    url: '/api/k8s-history/cleanup',
    method: 'post',
    data: {
      beforeDate
    }
  })
}

// 历史数据类型定义
export interface K8sHistoryQuery {
  page: number
  pageSize: number
  startTime?: string
  endTime?: string
}

export interface K8sPodHistory {
  id: number
  original_id: number
  config_id: number
  name: string
  namespace: string
  workload_name?: string
  workload_kind?: string
  status: string
  phase?: string
  node_name?: string
  pod_ip?: string
  host_ip?: string
  restart_count: number
  start_time?: string
  created_at: string
  updated_at: string
  deleted_at?: string
  archived_at: string
  archive_reason: string
}

export interface K8sNodeHistory {
  id: number
  original_id: number
  config_id: number
  name: string
  internal_ip?: string
  external_ip?: string
  hostname?: string
  os_image?: string
  kernel_version?: string
  container_runtime?: string
  kubelet_version?: string
  kube_proxy_version?: string
  cpu_capacity?: string
  memory_capacity?: string
  pods_capacity?: string
  cpu_allocatable?: string
  memory_allocatable?: string
  pods_allocatable?: string
  cpu_usage?: string
  memory_usage?: string
  pods_usage: number
  labels?: string
  annotations?: string
  taints?: string
  conditions?: string
  status: string
  ready: boolean
  schedulable: boolean
  created_at: string
  updated_at: string
  deleted_at?: string
  archived_at: string
  archive_reason: string
}

export interface K8sWorkloadHistory {
  id: number
  original_id: number
  config_id: number
  name: string
  namespace: string
  kind: string
  replicas: number
  ready_replicas: number
  status?: string
  labels?: string
  selector?: string
  images?: string
  cpu_request?: string
  cpu_limit?: string
  memory_request?: string
  memory_limit?: string
  created_at: string
  updated_at: string
  deleted_at?: string
  archived_at: string
  archive_reason: string
}

export interface K8sHistoryStatistics {
  podHistoryCount: number
  nodeHistoryCount: number
  workloadHistoryCount: number
  totalHistoryCount: number
  lastArchivedAt?: string
}
