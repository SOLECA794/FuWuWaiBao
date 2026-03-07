import { API_BASE } from '../../config/api'
import { studentCoursewareApi } from './coursewareApi'
import { studentSessionApi } from './sessionApi'
import { studentQaApi } from './qaApi'
import { studentWeakPointApi } from './weakPointApi'
import { studentKnowledgeApi } from './knowledgeApi'

export const studentV1Api = {
  health: () => fetch(`${API_BASE}/health`, { cache: 'no-store' }),
  coursewares: studentCoursewareApi,
  sessions: studentSessionApi,
  qa: studentQaApi,
  weakPoints: studentWeakPointApi,
  knowledge: studentKnowledgeApi
}
