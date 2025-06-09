import router from './router'
import { useUserStore } from '@/stores/user'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

NProgress.configure({ showSpinner: false })

const whiteList = ['/login', '/auth-redirect']

router.beforeEach(async (to, from, next) => {
  NProgress.start()

  // 直接从localStorage获取token
  const token = localStorage.getItem('token')
  const hasToken = !!token
  
  const userStore = useUserStore()

  if (hasToken) {
    if (to.path === '/login') {
      next({ path: '/' })
      NProgress.done()
    } else {
      const hasUserInfo = userStore.userInfo !== null
      if (hasUserInfo) {
        next()
      } else {
        try {
          // 获取用户信息
          await userStore.getUserInfo()
          next({ ...to, replace: true })
        } catch (error) {
          // 获取用户信息失败，清除token并跳转到登录页
          userStore.clearAuth()
          next(`/login?redirect=${to.path}`)
          NProgress.done()
        }
      }
    }
  } else {
    if (whiteList.indexOf(to.path) !== -1) {
      next()
    } else {
      next(`/login?redirect=${to.path}`)
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  NProgress.done()
}) 