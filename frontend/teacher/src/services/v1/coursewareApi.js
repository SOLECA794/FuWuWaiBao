import { requestJson } from './request'

export const teacherCoursewareApi = {
  list: () => requestJson('/api/v1/teacher/coursewares'),
  upload: (formData) => requestJson('/api/v1/teacher/coursewares/upload', {
    method: 'POST',
    body: formData
  }),
  remove: (courseId) => requestJson(`/api/v1/teacher/coursewares/${encodeURIComponent(courseId)}`, {
    method: 'DELETE'
  }),
  publish: ({ courseId, scope }) => requestJson(`/api/v1/teacher/coursewares/${encodeURIComponent(courseId)}/publish`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ scope })
  }),
  getScript: (courseId, pageNum) => requestJson(`/api/v1/teacher/coursewares/${encodeURIComponent(courseId)}/scripts/${pageNum}`),
  saveScript: ({ courseId, pageNum, content }) => requestJson(`/api/v1/teacher/coursewares/${encodeURIComponent(courseId)}/scripts/${pageNum}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ content })
  }),
  generateScript: ({ courseId, pageNum, mode = 'llm' }) => requestJson(`/api/v1/teacher/coursewares/${encodeURIComponent(courseId)}/scripts/ai-generate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ pageNum, mode })
  })
}
