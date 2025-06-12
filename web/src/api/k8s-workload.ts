import request from '@/utils/request'

export interface WorkloadQueryParams {
  page: number
  pageSize: number
  name?: string
  namespace?: string
  workloadType?: string
  configId?: string
}

export interface Workload {
  id: number
  configId: number
  name: string
  namespace: string
  kind: string
  replicas: number
  readyReplicas: number
  status: string
  labels: Record<string, string>
  selector: Record<string, string>
  images: string[]
  cpuRequest?: string
  cpuLimit?: string
  memoryRequest?: string
  memoryLimit?: string
  createdAt: string
  updatedAt: string
}

export interface WorkloadListResponse {
  list: Workload[]
  total: number
}

// 获取工作负载列表
export function getWorkloadList(params: WorkloadQueryParams) {
  return request<WorkloadListResponse>({
    url: '/k8s-workloads',
    method: 'get',
    params
  })
}

// 获取工作负载详情
export function getWorkloadDetail(id: number) {
  return request<Workload>({
    url: `/k8s-workloads/${id}`,
    method: 'get'
  })
}
