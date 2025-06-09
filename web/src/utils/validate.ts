/**
 * 验证用户名
 * @param rule 验证规则
 * @param value 用户名
 * @param callback 回调函数
 */
export function validateUsername(rule: any, value: string, callback: Function): void {
  if (value.trim().length === 0) {
    callback(new Error('请输入用户名'))
  } else {
    callback()
  }
}

/**
 * 验证URL
 * @param url 待验证的URL
 * @returns 是否为有效URL
 */
export function isValidURL(url: string): boolean {
  const pattern = new RegExp(
    '^(https?:\\/\\/)?' + // 协议
    '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|' + // 域名
    '((\\d{1,3}\\.){3}\\d{1,3}))' + // IP地址
    '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*' + // 端口和路径
    '(\\?[;&a-z\\d%_.~+=-]*)?' + // 查询字符串
    '(\\#[-a-z\\d_]*)?$', // 锚点
    'i'
  )
  return pattern.test(url)
}

/**
 * 验证是否为外部链接
 * @param path 路径
 * @returns 是否为外部链接
 */
export function isExternal(path: string): boolean {
  return /^(https?:|mailto:|tel:)/.test(path)
}

/**
 * 验证是否为有效的邮箱
 * @param email 邮箱
 * @returns 是否为有效的邮箱
 */
export function isValidEmail(email: string): boolean {
  const pattern = /^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$/
  return pattern.test(email)
}

/**
 * 验证是否为有效的手机号
 * @param phone 手机号
 * @returns 是否为有效的手机号
 */
export function isValidPhone(phone: string): boolean {
  const pattern = /^1[3456789]\d{9}$/
  return pattern.test(phone)
} 