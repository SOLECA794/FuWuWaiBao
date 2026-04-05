import { requestJson } from './v1/request'

function generateTempUuid(prefix = 'tmp') {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID()
  }
  const seed = `${Date.now().toString(16)}_${Math.random().toString(16).slice(2)}`
  return `${prefix}_${seed}`
}

function normalizeItemsWithId(items, prefix) {
  if (!Array.isArray(items)) return []
  return items.map((item, index) => {
    const safeItem = item && typeof item === 'object' ? item : { value: item }
    return {
      ...safeItem,
      id: safeItem.id || generateTempUuid(`${prefix}_${index}`)
    }
  })
}

/**
 * 获取教学迭代概览。
 * @param {string} courseId
 * @returns {Promise<{basicNodeTree: Array, pendingNodes: Array, pendingCases: Array}>}
 */
export async function getIterationOverview(courseId) {
  const normalizedCourseId = encodeURIComponent(String(courseId || '').trim())
  const payload = await requestJson(`/api/v1/teacher/iteration-overview/${normalizedCourseId}`, {
    method: 'GET'
  })

  const source = payload && typeof payload === 'object' ? (payload.data || payload) : {}

  return {
    ...source,
    basicNodeTree: normalizeItemsWithId(source.basicNodeTree, 'basic_node'),
    pendingNodes: normalizeItemsWithId(source.pendingNodes, 'pending_node'),
    pendingCases: normalizeItemsWithId(source.pendingCases, 'pending_case')
  }
}

/**
 * 基于教学迭代节点顺序生成讲稿。
 * @param {{ courseId?: string, nodeOrder: string[] }} payload
 * @returns {Promise<{ data: string }>}
 */
export async function generateIterationScript(payload) {
  const body = payload && typeof payload === 'object' ? payload : { nodeOrder: [] }
  const response = await requestJson('/api/v1/teacher/iteration/script-generate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  })

  const data = response && typeof response === 'object' ? (response.data || response) : ''
  if (typeof data === 'string') return { data }
  if (data && typeof data === 'object') {
    return { data: String(data.markdown || data.content || '') }
  }
  return { data: '' }
}

export const teacherV1Service = {
  getIterationOverview,
  generateIterationScript
}
