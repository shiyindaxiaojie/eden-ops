import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import NProgress from 'nprogress'
import Layout from '@/layout/index.vue'
import { getToken } from '@/utils/auth'

export const constantRoutes: Array<RouteRecordRaw> = [
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '仪表盘', icon: 'Odometer' }
      }
    ]
  },
  {
    path: '/tools',
    component: Layout,
    redirect: '/tools/ip-locator',
    meta: { title: '常用工具', icon: 'Tools' },
    children: [
      {
        path: 'ip-locator',
        name: 'IpLocator',
        component: () => import('@/views/tools/ip-locator/index.vue'),
        meta: { title: 'IP定位器', icon: 'Location' }
      }
    ]
  },
  {
    path: '/infrastructure',
    component: Layout,
    redirect: '/infrastructure/kubernetes',
    meta: { title: '基础设施', icon: 'Monitor' },
    children: [
      {
        path: 'kubernetes',
        component: () => import('@/views/infrastructure/kubernetes/index.vue'),
        name: 'Kubernetes',
        meta: { title: 'Kubernetes 集群', icon: 'kubernetes' }
      },
      {
        path: 'kubernetes/detail/:id',
        component: () => import('@/views/infrastructure/kubernetes/detail.vue'),
        name: 'KubernetesDetail',
        meta: { title: 'Kubernetes 集群明细', activeMenu: '/infrastructure/kubernetes' },
        hidden: true
      },
      {
        path: 'kubernetes/workloads',
        component: () => import('@/views/infrastructure/kubernetes/workloads.vue'),
        name: 'KubernetesWorkloads',
        meta: { title: 'Kubernetes 工作负载', activeMenu: '/infrastructure/kubernetes' },
        hidden: true
      },
      {
        path: 'kubernetes/pods',
        component: () => import('@/views/infrastructure/kubernetes/pods.vue'),
        name: 'KubernetesPods',
        meta: { title: 'Kubernetes Pod', activeMenu: '/infrastructure/kubernetes' },
        hidden: true
      },
      {
        path: 'kubernetes/nodes',
        component: () => import('@/views/infrastructure/kubernetes/nodes.vue'),
        name: 'KubernetesNodes',
        meta: { title: 'Kubernetes 节点', activeMenu: '/infrastructure/kubernetes' },
        hidden: true
      },
      {
        path: 'kubernetes/history/:configId/:type',
        component: () => import('@/views/infrastructure/kubernetes/history.vue'),
        name: 'KubernetesHistory',
        meta: { title: 'Kubernetes 历史数据', activeMenu: '/infrastructure/kubernetes' },
        hidden: true
      },
      {
        path: 'cloud-provider',
        name: 'CloudProvider',
        component: () => import('@/views/infrastructure/cloud-provider/index.vue'),
        meta: { title: '云厂商管理', icon: 'cloud' }
      },
      {
        path: 'cloud-account',
        name: 'CloudAccount',
        component: () => import('@/views/infrastructure/cloud-account/index.vue'),
        meta: { title: '云账号管理', icon: 'Cloudy' }
      },
      {
        path: 'database',
        name: 'Database',
        component: () => import('@/views/infrastructure/database/index.vue'),
        meta: { title: '数据库配置', icon: 'DataLine' }
      },
      {
        path: 'server',
        name: 'Server',
        component: () => import('@/views/infrastructure/server/index.vue'),
        meta: { title: '服务器配置', icon: 'Monitor' }
      }
    ]
  },
  {
    path: '/system',
    component: Layout,
    redirect: '/system/user',
    meta: { title: '系统管理', icon: 'Setting' },
    children: [
      {
        path: 'user',
        name: 'User',
        component: () => import('@/views/system/user/index.vue'),
        meta: { title: '用户管理', icon: 'User' }
      },
      {
        path: 'role',
        name: 'Role',
        component: () => import('@/views/system/role/index.vue'),
        meta: { title: '角色管理', icon: 'UserFilled' }
      },
      {
        path: 'menu',
        name: 'Menu',
        component: () => import('@/views/system/menu/index.vue'),
        meta: { title: '菜单管理', icon: 'Menu' }
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes: constantRoutes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  NProgress.start()
  
  // 设置页面标题
  const title = to.meta.title ? `${to.meta.title} - eden*` : 'eden*'
  document.title = title

  const token = getToken()
  console.log('路由守卫检查Token:', token, '路由:', to.path)

  if (to.path === '/login') {
    if (token) {
      // 已登录时访问登录页，重定向到首页
      next({ path: '/' })
    } else {
      next()
    }
  } else {
    if (token) {
      next()
    } else {
      // 未登录访问其他页面，重定向到登录页
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    }
  }
})

router.afterEach(() => {
  NProgress.done()
})

export default router