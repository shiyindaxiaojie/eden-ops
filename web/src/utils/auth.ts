// Token key
export const TOKEN_KEY = 'token'

export const USER_INFO_KEY = 'eden-user'

// 获取 token
export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

// 设置 token
export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

// 删除 token
export function removeToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

export function getUserInfo(): any | null {
  const userInfo = localStorage.getItem(USER_INFO_KEY)
  return userInfo ? JSON.parse(userInfo) : null
}

export function setUserInfo(userInfo: any): void {
  localStorage.setItem(USER_INFO_KEY, JSON.stringify(userInfo))
}

export function removeUserInfo(): void {
  localStorage.removeItem(USER_INFO_KEY)
}

export function clearAuth(): void {
  removeToken()
  removeUserInfo()
} 