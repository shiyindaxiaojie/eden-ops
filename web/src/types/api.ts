// 基础响应接口
export interface BaseResponse<T> {
  code: number
  message: string
  data: T
}

// 分页查询参数
export interface PageQuery {
  page: number
  pageSize: number
}

// 分页结果
export interface PageResult<T> {
  list: T[]
  total: number
}

// 用户接口
export interface User {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  status: string
  roleIds: number[]
  roles: Role[]
  created_at: string
  updated_at: string
}

// 角色接口
export interface Role {
  id: number
  name: string
  code: string
  status: string
  remark: string
  menuIds: number[]
  menus: Menu[]
  created_at: string
  updated_at: string
}

// 菜单接口
export interface Menu {
  id: number
  parent_id: number
  name: string
  path: string
  component: string
  icon: string
  sort: number
  type: string
  permission: string
  status: string
  children?: Menu[]
  created_at: string
  updated_at: string
}

// IP 定位接口
export interface IPLocation {
  ip: string
  country: string
  region: string
  city: string
  isp: string
  latitude: number
  longitude: number
}

// 云账号接口
export interface CloudAccount {
  id: number
  name: string
  providerId?: number
  accessKey: string
  secretKey: string
  description?: string
  status: number
  createdAt: string
  updatedAt: string
}

// 数据库配置接口
export interface DatabaseConfig {
  id: number
  name: string
  host: string
  port: number
  username: string
  password: string
  database: string
  type: string
  status: string
  remark: string
  created_at: string
  updated_at: string
}

// Kubernetes配置接口
export interface K8sConfig {
  id: number
  name: string
  kubeconfig: string
  status: string
  remark: string
  created_at: string
  updated_at: string
}

// 服务器配置接口
export interface ServerConfig {
  id: number
  name: string
  host: string
  port: number
  username: string
  password: string
  type: string
  status: string
  remark: string
  created_at: string
  updated_at: string
}

// Kubernetes工作负载
export interface K8sWorkload {
  id: number
  config_id: number
  name: string
  namespace: string
  kind: string
  replicas: number
  ready_replicas: number
  status: string
  created_at: string
  updated_at: string
  deleted_at: string | null
}

// Kubernetes Deployment
export interface K8sDeployment {
  name: string
  namespace: string
  replicas: number
  ready_replicas: number
  created_at: string
}

// Kubernetes StatefulSet
export interface K8sStatefulSet {
  name: string
  namespace: string
  replicas: number
  ready_replicas: number
  created_at: string
}

// Kubernetes DaemonSet
export interface K8sDaemonSet {
  name: string
  namespace: string
  desired_number_scheduled: number
  number_ready: number
  created_at: string
}

// Kubernetes Job
export interface K8sJob {
  name: string
  namespace: string
  completions: number
  succeeded: number
  created_at: string
}

// Kubernetes CronJob
export interface K8sCronJob {
  name: string
  namespace: string
  schedule: string
  last_schedule: string
  created_at: string
} 