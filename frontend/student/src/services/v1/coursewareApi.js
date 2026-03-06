import { requestJson } from './request'

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
  })
}
