import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login, logout, getUserInfo } from '@/api/auth'
import { getToken, setToken, removeToken } from '@/utils/auth'
import type { User } from '@/types/api'

interface LoginResponse {
  token: string
  user: User
}

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(getToken())
  const user = ref<User | null>(null)

  // 登录
  async function loginAction(username: string, password: string) {
    try {
      const res = await login({ username, password })
      const { token: newToken, user: userInfo } = res.data
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
      removeToken()
    } catch (error) {
      console.error('登出失败:', error)
      throw error
    }
  }

  // 获取用户信息
  async function getUserInfoAction() {
    try {
      const res = await getUserInfo()
      user.value = res.data
      return res
    } catch (error) {
      console.error('获取用户信息失败:', error)
      throw error
    }
  }

  return {
    token,
    user,
    loginAction,
    logoutAction,
    getUserInfoAction
  }
})