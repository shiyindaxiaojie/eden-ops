import request from '@/utils/request'
import type { BaseResponse } from '@/types/api'

// 获取K8s命名空间列表
export function getK8sNamespaces(configId: string) {
  return request<BaseResponse<string[]>>({
    url: '/api/v1/k8s-namespaces',
    method: 'get',
    params: {
      configId
    }
  })
}
