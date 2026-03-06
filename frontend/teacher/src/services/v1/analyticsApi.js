import { requestJson } from './request'

export const teacherAnalyticsApi = {
  getStats: (courseId) => requestJson(`/api/v1/teacher/coursewares/${encodeURIComponent(courseId)}/stats`),
  getCardData: (courseId) => requestJson(`/api/v1/teacher/coursewares/${encodeURIComponent(courseId)}/card-data`),
  getQuestionRecords: (courseId, page = 1, pageSize = 100) => requestJson(`/api/v1/teacher/coursewares/${encodeURIComponent(courseId)}/questions?page=${page}&pageSize=${pageSize}`)
}
