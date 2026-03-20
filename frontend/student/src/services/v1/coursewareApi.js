import { requestJson } from './request'

const buildQuery = (params) => {
  return Object.entries(params)
    .filter(([, value]) => value !== undefined && value !== null && value !== '')
    .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
    .join('&')
}

export const studentCoursewareApi = {
  list: () => requestJson('/api/v1/student/coursewares'),
  getPlaybackScript: (courseId, pageNum) => requestJson(`/api/v1/student/coursewares/${encodeURIComponent(courseId)}/scripts/${pageNum}`),
  getBreakpoint: (studentId, courseId) => requestJson(`/api/v1/student/coursewares/${encodeURIComponent(courseId)}/breakpoint?studentId=${encodeURIComponent(studentId)}`),
  updateBreakpoint: ({ studentId, courseId, pageNum }) => requestJson(`/api/v1/student/coursewares/${encodeURIComponent(courseId)}/breakpoint`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ studentId, pageNum })
  }),
  getStats: (studentId, courseId) => requestJson(`/api/v1/student/coursewares/${encodeURIComponent(courseId)}/stats?studentId=${encodeURIComponent(studentId)}`),
  getWeakPoints: (studentId, courseId) => requestJson(`/api/v1/student/coursewares/${encodeURIComponent(courseId)}/weak-points?studentId=${encodeURIComponent(studentId)}`),
  saveNote: ({ studentId, courseId, pageNum, content, x = 0, y = 0 }) => requestJson(`/api/v1/student/coursewares/${encodeURIComponent(courseId)}/notes`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ studentId, pageNum, content, x, y })
  }),

  // 新增接口：分页查笔记
  listNotes: ({ studentId, courseId, pageNum, pageSize }) => {
    const query = buildQuery({ studentId, courseId, pageNum, pageSize })
    return requestJson(`/api/v1/student/notes?${query}`)
  },

  // 新增接口：新增收藏
  addFavorite: ({ studentId, courseId, nodeId, pageNum, title }) => requestJson('/api/v1/student/favorites', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ studentId, courseId, nodeId, pageNum, title })
  }),

  // 新增接口：分页查收藏
  listFavorites: ({ studentId, courseId, page, pageSize }) => {
    const query = buildQuery({ studentId, courseId, page, pageSize })
    return requestJson(`/api/v1/student/favorites?${query}`)
  },

  // 新增接口：删除收藏
  deleteFavorite: (favoriteId) => requestJson(`/api/v1/student/favorites/${encodeURIComponent(favoriteId)}`, {
    method: 'DELETE'
  }),

  // 新增接口：生成练习题
  generatePractice: ({ studentId, courseId, nodeId, pageNum, difficulty, count }) => requestJson('/api/v1/student/practice/generate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ studentId, courseId, nodeId, pageNum, difficulty, count })
  }),

  // 新增接口：提交练习答案
  submitPractice: ({ taskId, studentId, answers }) => requestJson('/api/v1/student/practice/submit', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ taskId, studentId, answers })
  }),

  // 新增接口：节点专项讲解
  explainNode: ({ nodeId, courseId, pageNum, question }) => requestJson(`/api/v1/student/nodes/${encodeURIComponent(nodeId)}/explain`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ courseId, pageNum, question })
  })
}
