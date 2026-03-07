import { API_BASE } from '../../config/api'

export async function requestJson(path, options = {}) {
  const response = await fetch(`${API_BASE}${path}`, options)
  const payload = await response.json().catch(() => ({}))
  if (!response.ok || (typeof payload.code !== 'undefined' && payload.code !== 200)) {
    throw new Error(payload.message || `请求失败: ${response.status}`)
  }
  return payload
}
