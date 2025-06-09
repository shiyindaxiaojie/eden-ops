import request from '@/utils/request'
import type { IPLocation } from '@/types/api'

// IP 定位响应数据
interface IPLocation {
  ip: string
  country: string
  region: string
  city: string
  isp: string
  latitude: number
  longitude: number
}

// IP 定位
export function locateIP(ip: string) {
  return request<IPLocation>({
    url: `/api/v1/tools/ip/${ip}`,
    method: 'get'
  })
}