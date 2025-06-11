import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login, logout, getUserInfo } from '@/api/auth'
import { getToken, setToken, removeToken, clearAuth } from '@/utils/auth'
import type { User } from '@/types/api'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(getToken())
  const user = ref<User | null>(null)

  // 登录
  async function loginAction(username: string, password: string) {
    try {
      const res = await login({ username, password })
      // login 函数已经返回了 res.data，所以这里直接解构
      const { token: newToken, user: userInfo } = res
      token.value = newToken
      user.value = userInfo
      setToken(newToken)
      return res
    } catch (error) {
      token.value = null
      user.value = null
      removeToken()
      throw error
    }
  }

  // 登出
  async function logoutAction() {
    try {
      await logout()
      token.value = null
      user.value = null
      clearAuth() // 清除所有认证信息
      // 跳转到登录页
      router.push('/login')
    } catch (error) {
      console.error('登出失败:', error)
      // 即使后端登出失败，也要清除本地状态
      token.value = null
      user.value = null
      clearAuth() // 清除所有认证信息
      router.push('/login')
      throw error
    }
  }

  // 获取用户信息
  async function getUserInfoAction() {
    try {
      const res = await getUserInfo()
      user.value = res
      return res
    } catch (error) {
      console.error('获取用户信息失败:', error)
      throw error
    }
  }

  return {
    token,
    user,
    userInfo: user, // 添加别名以兼容现有代码
    loginAction,
    logoutAction,
    logout: logoutAction, // 添加别名以兼容现有代码
    getUserInfoAction
  }
})