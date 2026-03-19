<template>
  <HomeLogin v-if="!isLoggedIn" @login-success="handleLoginSuccess" />
  <div v-else class="teacher-app">
    <div class="workspace-shell">
      <TeacherTopBar
        :backend-status-class="backendStatusClass"
        :backend-status-text="backendStatusText"
        :username="loggedInUsername"
        @logout="isLoggedIn = false"
      />

      <div class="main-content">
      <!-- 方案修改：左侧 MENU 导航栏 (带 ins 风图标) -->
      <div class="left-sidebar-menu" :class="{ 'collapsed': isLeftMenuCollapsed }">
        <div class="menu-header">
          <span v-show="!isLeftMenuCollapsed">菜单</span>
          <button class="menu-toggle-btn" @click="isLeftMenuCollapsed = !isLeftMenuCollapsed" :title="isLeftMenuCollapsed ? '展开菜单' : '收起菜单'">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 12h18"></path><path d="M3 6h18"></path><path d="M3 18h18"></path></svg>
          </button>
        </div>
        <div class="menu-list">
          <div class="menu-item" :class="{ active: activeTab === 'script' }" @click="activeTab = 'script'" title="编辑讲稿">
            <svg class="ins-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
            <span v-show="!isLeftMenuCollapsed">编辑讲稿</span>
          </div>
          <div class="menu-item" :class="{ active: activeTab === 'stats' }" @click="activeTab = 'stats'" title="学情分析">
            <svg class="ins-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="20" x2="18" y2="10"></line><line x1="12" y1="20" x2="12" y2="4"></line><line x1="6" y1="20" x2="6" y2="14"></line></svg>
            <span v-show="!isLeftMenuCollapsed">学情分析</span>
          </div>
          <div class="menu-item" :class="{ active: activeTab === 'questions' }" @click="activeTab = 'questions'" title="提问统计">
            <svg class="ins-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
            <span v-show="!isLeftMenuCollapsed">提问统计</span>
          </div>
          <div class="menu-item" :class="{ active: activeTab === 'card' }" @click="activeTab = 'card'" title="卡点可视化">
            <svg class="ins-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg>
            <span v-show="!isLeftMenuCollapsed">卡点可视化</span>
          </div>
          <div class="menu-item" :class="{ active: activeTab === 'platform' }" @click="activeTab = 'platform'" title="平台管理">
              <svg class="ins-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
               <svg class="ins-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="7" height="7"></rect>
                <rect x="14" y="3" width="7" height="7"></rect>
                <rect x="14" y="14" width="7" height="7"></rect>
                <rect x="3" y="14" width="7" height="7"></rect>
              </svg>

              </svg>
        <span v-show="!isLeftMenuCollapsed">平台管理 <small>(悬停后展开)</small></span>
          </div>

        </div>
      </div>

      <!-- 中间内容编辑区 -->
      <div class="editor-section">
        <div v-if="activeTab === 'script'" class="tab-container script-mode">
          <div v-if="!currentCourseId" class="empty-tip-container">
            <div class="empty-tip">请先在右侧选择或上传一个课件以编辑讲稿</div>
          </div>
          <TeacherScriptPanel
            v-else
            :preview-url="realPreviewUrl"
            :current-course-id="currentCourseId"
            :current-edit-page="currentEditPage"
            :total-pages="currentCourseTotalPages"
            :current-script="currentScript"
            :current-script-nodes="currentScriptNodes"
            :mapping-coverage="currentMappingCoverage"
            :node-stats="studentStats.nodeStats || []"
            :focus-node-request="focusNodeRequest"
            :script-generating="scriptGenerating"
            :script-saving="scriptSaving"
            @generate-ai-script="generateAIScript"
            @save-script="saveScript"
            @update:current-script="currentScript = $event"
            @update:current-script-nodes="currentScriptNodes = $event"
            @prev-page="prevPage"
            @next-page="nextPage"
          ></TeacherScriptPanel>
        </div>

        <div v-else-if="activeTab === 'platform'" class="tab-container">
            <PlatformManagementPanel />
        </div>
  
      

        <div v-else-if="activeTab === 'stats'" class="tab-container">
          <div v-if="!currentCourseId" class="empty-tip-container">
            <div class="empty-tip">请先在右侧选择或上传一个课件查看学情数据</div>
          </div>
          <TeacherStatsPanel
            v-else
            :current-course-id="currentCourseId"
            :current-course-name="currentCourseName"
            :student-stats="studentStats"
          ></TeacherStatsPanel>
        </div>

        <div v-else-if="activeTab === 'questions'" class="tab-container">
          <div v-if="!currentCourseId" class="empty-tip-container">
            <div class="empty-tip">请先在右侧选择或上传一个课件查看提问统计</div>
          </div>
          <TeacherQuestionsPanel
            v-else
            :current-course-id="currentCourseId"
            :current-course-name="currentCourseName"
            :current-course-total-pages="currentCourseTotalPages"
            :filter-page="filterPage"
            :filtered-questions="filteredQuestions"
            :uncovered-node-ids="studentStats.mappingCoverage?.uncoveredNodeIds || []"
            @update:filter-page="filterPage = $event"
            @focus-node="focusQuestionNode"
          ></TeacherQuestionsPanel>
        </div>

        <div v-else class="tab-container">
          <div v-if="!currentCourseId" class="empty-tip-container">
            <div class="empty-tip">请先在右侧选择或上传一个课件查看卡点分析</div>
          </div>
          <TeacherCardAnalysisPanel
            v-else
            :current-course-id="currentCourseId"
            :current-course-name="currentCourseName"
            :chart-type="chartType"
            :card-data="cardData"
            @update:chart-type="chartType = $event"
          ></TeacherCardAnalysisPanel>
        </div>
      </div>

      <!-- 右侧课件管理（原左侧侧边栏移到右侧） -->
      <TeacherCoursewareSidebar
        v-if="activeTab !== 'script'"
        v-show="isSidebarVisible"
        :courseware-list="coursewareList"
        :current-course-id="currentCourseId"
        :current-course-name="currentCourseName"
        :current-course-total-pages="currentCourseTotalPages"
        :current-edit-page="currentEditPage"
        :course-list-loading="courseListLoading"
        @open-publish="showPublishModal = true"
        @open-upload="showUploadModal = true"
        @select-course="selectCourse"
        @delete-course="deleteCourse"
        @select-page="selectEditPage"
        style="border-right: none; border-left: 1px solid #e2e8f0;"
      />
    </div>
    </div>

    <TeacherUploadModal
      :visible="showUploadModal"
      :selected-file-name="selectedFileName"
      :upload-loading="uploadLoading"
      @close="showUploadModal = false"
      @select-file="handleFileSelect"
      @submit="uploadCourseware"
      @file-input-ready="fileInput = $event"
    ></TeacherUploadModal>

    <TeacherPublishModal
      :visible="showPublishModal"
      :current-course-name="currentCourseName"
      :publish-scope="publishScope"
      @close="showPublishModal = false"
      @submit="publishCourseware"
      @update:publish-scope="publishScope = $event"
    ></TeacherPublishModal>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onUnmounted, watch } from 'vue'
import { API_BASE } from './config/api'
import { teacherV1Api } from './services/v1'
// 导入所有组件
import TeacherTopBar from './components/teacher/TeacherTopBar.vue'
import TeacherOverviewStrip from './components/teacher/TeacherOverviewStrip.vue'
import TeacherCoursewareSidebar from './components/teacher/TeacherCoursewareSidebar.vue'
import TeacherScriptPanel from './components/teacher/TeacherScriptPanel.vue'
import TeacherStatsPanel from './components/teacher/TeacherStatsPanel.vue'
import TeacherQuestionsPanel from './components/teacher/TeacherQuestionsPanel.vue'
import TeacherCardAnalysisPanel from './components/teacher/TeacherCardAnalysisPanel.vue'
import TeacherUploadModal from './components/teacher/TeacherUploadModal.vue'
import TeacherPublishModal from './components/teacher/TeacherPublishModal.vue'
import HomeLogin from './components/HomeLogin.vue'
import PlatformManagementPanel from './components/teacher/PlatformManagementPanel.vue'

// --- 状态管理 ---
const isLoggedIn = ref(false)
const loggedInUsername = ref('')
const activeTab = ref('script') // 或 'platform' 如果希望默认显示平台管理

const handleLoginSuccess = (user) => {
  if (user.role === 'student') {
    window.location.href = 'http://localhost:8080'
  } else {
    loggedInUsername.value = user.username
    isLoggedIn.value = true
  }
}

const coursewareList = ref([])
const currentCourseId = ref('')
const currentCourseName = ref('')
const currentCourseTotalPages = ref(0)
const currentEditPage = ref(1)
const currentScript = ref('')
const currentScriptNodes = ref([])
const currentMappingCoverage = ref(null)
const focusNodeRequest = ref(null)

const showUploadModal = ref(false)
const fileInput = ref(null)
const selectedFileName = ref('')
const uploadLoading = ref(false)

const showPublishModal = ref(false)
const publishScope = ref('all')
const backendStatus = ref('checking')
const courseListLoading = ref(false)
const scriptGenerating = ref(false)
const scriptSaving = ref(false)

const showPlatformManagement = ref(false)

// --- 生命周期钩子 ---
onMounted(async () => {
  // 检查登录状态
  const user = await teacherV1Api.checkLoginStatus()
  if (user) {
    handleLoginSuccess(user)
  }

  // 检查后端状态
  checkBackendStatus()

  // 检查课件列表
  await fetchCoursewareList()
})

onUnmounted(() => {
  clearInterval(backendHealthTimer)
})
const checkBackendStatus = async () => {
  try {
    const res = await teacherV1Api.health()
    backendStatus.value = res.ok ? 'online' : 'offline'
  } catch (error) {
    backendStatus.value = 'offline'
    console.error('后端状态检查失败:', error)
  }
}


const isSidebarVisible = ref(window.innerWidth > 1024)
const isLeftMenuCollapsed = ref(window.innerWidth <= 1024)
const studentStats = ref({
  totalQuestions: 0,
  hotPages: [],
  keyDifficulties: '暂无',
  nodeStats: [],
  mappingCoverage: null,
  activeSessions: 0,
  avgTurnsPerSession: 0,
  nodeHeatmap: [],
  masteryRadar: { indicators: [], values: [], avgMastery: 0 },
  classTrend: [],
  learningInsights: { reteachNodes: [], prerequisiteGaps: [], summary: '' }
})

watch(currentMappingCoverage, (coverage) => {
  studentStats.value = {
    ...studentStats.value,
    mappingCoverage: coverage || null
  }
})

const cardData = ref([])
const chartType = ref('bar')
const questionRecords = ref([])
const filterPage = ref('')
const previewLoading = ref(false) // 仅用于本地动画，不发请求
let backendHealthTimer = null

// --- 计算属性 ---
const filteredQuestions = computed(() => {
  if (!filterPage.value) return questionRecords.value
  return questionRecords.value.filter(q => q.page === Number(filterPage.value))
})

const currentCoursePublished = computed(() => {
  const currentCourse = coursewareList.value.find(item => item.id === currentCourseId.value)
  return !!currentCourse?.published
})

const backendStatusText = computed(() => {
  return backendStatus.value === 'online' ? '在线' : backendStatus.value === 'offline' ? '离线' : '检测中'
})

const backendStatusClass = computed(() => {
  return backendStatus.value === 'online' ? 'online' : backendStatus.value === 'offline' ? 'offline' : 'checking'
})

const buildNodesFromScriptText = (scriptText, page) => {
  const raw = String(scriptText || '').trim()
  if (!raw) return []
  return raw
    .split(/\n{2,}|(?<=[。！？])/)
    .map(item => item.trim())
    .filter(Boolean)
    .map((text, index, list) => ({
      id: '',
      nodeId: `p${page}_n${index + 1}`,
      type: index === 0 ? 'opening' : index === list.length - 1 ? 'transition' : 'explain',
      title: index === 0 ? '节点1：开场' : `节点${index + 1}：讲解`,
      summary: text.length > 48 ? `${text.slice(0, 48)}...` : text,
      scriptText: text,
      reteachScript: '',
      transitionText: '',
      structuredMarkdown: '',
      knowledgeNodesJson: '[]',
      scriptSegmentsJson: '[]',
      estimatedDuration: 30,
      sortOrder: index + 1
    }))
}

// 新增：真实预览URL（假设后端返回图片）
const realPreviewUrl = computed(() => {
  if (!currentCourseId.value) return ''
  // 后端统一前缀为 /api，保持与其他接口一致
  return `${API_BASE}/api/courseware/${currentCourseId.value}/page/${currentEditPage.value}?t=${Date.now()}`
})

// 新增：模拟预览URL（根据页码动态变化）
const mockPreviewUrl = computed(() => {
  if (!currentCourseId.value) return ''
  // 使用picsum的random参数，确保每页图片不同
  return `https://picsum.photos/800/600?random=${currentCourseId.value}_${currentEditPage.value}`
})

// 新增：预览加载成功处理
const onPreviewLoad = () => {
  previewLoading.value = false
  console.log('模拟预览图片加载成功:', mockPreviewUrl.value)
}

// 新增：预览错误处理
const handlePreviewError = () => {
  previewLoading.value = false
  console.error('模拟预览图片加载失败:', mockPreviewUrl.value)
  alert('模拟预览图片加载失败，请检查网络')
}

// --- 生命周期 ---
onMounted(async () => {
  // 先检查后端健康状态，但不阻塞列表加载
  checkBackendHealth()
  backendHealthTimer = window.setInterval(checkBackendHealth, 30 * 1000)
  
  // 确保列表在初始化时必定加载
  await loadCoursewareList()
  
  // 如果列表加载成功且有数据，加载首个课件的上下文
  if (currentCourseId.value) {
    await loadCourseContext(currentCourseId.value)
  }
})

onUnmounted(() => {
  if (backendHealthTimer) window.clearInterval(backendHealthTimer)
})

// --- 核心方法 ---
const checkBackendHealth = async () => {
  try {
    const res = await teacherV1Api.health()
    backendStatus.value = res.ok ? 'online' : 'offline'
  } catch (error) {
    backendStatus.value = 'offline'
  }
}

// 仅加载课件列表，无任何额外操作（保证稳定）
const loadCoursewareList = async () => {
  courseListLoading.value = true
  try {
    const data = await teacherV1Api.coursewares.list()
    const list = data.data || []
    coursewareList.value = list.map(item => ({
      id: item.id,
      name: item.title,
      totalPages: item.total_page || 1,
      published: !!item.is_published
    }))

    if (coursewareList.value.length > 0) {
      const first = coursewareList.value[0]
      currentCourseId.value = first.id
      currentCourseName.value = first.name
      currentCourseTotalPages.value = first.totalPages
    }
  } catch (err) {
    console.error('加载课件列表失败', err)
  } finally {
    courseListLoading.value = false
  }
}

const loadCourseContext = async (courseId) => {
  previewLoading.value = true
  // 仅加载必要数据，不加载预览接口
  await Promise.all([
    loadScript(courseId, 1),
    loadStudentStats(courseId),
    loadCardData(courseId),
    loadQuestionRecords(courseId)
  ])
  // 本地模拟加载时长
  setTimeout(() => previewLoading.value = false, 500)
}

const selectCourse = async (course) => {
  currentCourseId.value = course.id
  currentCourseName.value = course.name
  currentCourseTotalPages.value = course.totalPages
  currentEditPage.value = 1
  await loadCourseContext(course.id)
}

const deleteCourse = async (courseId) => {
  if (!confirm('确定删除该课件吗？')) return
  try {
    await teacherV1Api.coursewares.remove(courseId)
    await loadCoursewareList()
    if (currentCourseId.value === courseId) {
      currentCourseId.value = ''
      currentCourseName.value = ''
      currentCourseTotalPages.value = 0
      currentScript.value = ''
      questionRecords.value = []
      cardData.value = []
    }
  } catch (err) {
    alert('删除课件失败：' + err.message)
  }
}

const selectEditPage = async (page) => {
  previewLoading.value = true
  currentEditPage.value = page
  await loadScript(currentCourseId.value, page)
  setTimeout(() => previewLoading.value = false, 300)
}

// --- 预览翻页 ---
const prevPage = () => currentEditPage.value > 1 && selectEditPage(currentEditPage.value - 1)
const nextPage = () => currentEditPage.value < currentCourseTotalPages.value && selectEditPage(currentEditPage.value + 1)

// --- 讲稿相关 ---
const loadScript = async (courseId, page) => {
  try {
    const data = await teacherV1Api.coursewares.getScript(courseId, page)
    currentScript.value = data?.data?.content || ''
    const nodes = data?.data?.nodes
    currentMappingCoverage.value = data?.data?.mappingCoverage || null
    currentScriptNodes.value = Array.isArray(nodes) && nodes.length > 0
      ? normalizeNodesFromServer(nodes)
      : buildNodesFromScriptText(currentScript.value, page)
  } catch (err) {
    currentScript.value = ''
    currentScriptNodes.value = []
    currentMappingCoverage.value = null
  }
}

const normalizeNodesFromServer = (nodes) => {
  const list = Array.isArray(nodes) ? nodes : []
  return list.map((node = {}, index) => {
    const knowledgeNodes = parseJSONArrayOrDefault(node.knowledgeNodesJson, node.knowledgeNodes)
    const scriptSegments = parseJSONArrayOrDefault(node.scriptSegmentsJson, node.scriptSegments)
    return {
      ...node,
      schemaVersion: Number(node.schemaVersion) || 2,
      sortOrder: Number(node.sortOrder) || index + 1,
      knowledgeNodes,
      scriptSegments,
      knowledgeNodesJson: JSON.stringify(knowledgeNodes),
      scriptSegmentsJson: JSON.stringify(scriptSegments)
    }
  })
}

const saveScript = async () => {
  if (!currentScript.value.trim()) return alert('请输入讲稿内容！')

  const validationError = validateNodeJSONFields(currentScriptNodes.value)
  if (validationError) {
    alert(validationError)
    return
  }

  const normalizedNodes = normalizeNodesForSave(currentScriptNodes.value)
  currentScriptNodes.value = normalizedNodes

  scriptSaving.value = true
  try {
    const courseId = currentCourseId.value
    const pageNum = currentEditPage.value
    const [, nodeSaveResult] = await Promise.all([
      teacherV1Api.coursewares.saveScript({
        courseId,
        pageNum,
        content: currentScript.value
      }),
      teacherV1Api.coursewares.saveNodes({
        courseId,
        pageNum,
        nodes: normalizedNodes
      })
    ])
    const coverage = nodeSaveResult?.data?.mappingCoverage || null
    currentMappingCoverage.value = coverage
    await loadStudentStats(courseId)
    if (coverage && Number(coverage.coverageRate || 0) < 1) {
      alert(`保存成功，但当前节点映射覆盖率为 ${Math.round(Number(coverage.coverageRate || 0) * 100)}%，存在未覆盖节点：${(coverage.uncoveredNodeIds || []).join('、')}`)
      return
    }
    alert('讲稿保存成功！')
  } catch (err) {
    const detail = String(err?.message || '').trim()
    if (detail) {
      alert('保存讲稿失败：' + detail)
      return
    }
    alert('保存讲稿失败：请检查节点映射与依赖关系配置')
  } finally {
    scriptSaving.value = false
  }
}

const normalizeNodesForSave = (nodes) => {
  const list = Array.isArray(nodes) ? nodes : []
  return list.map((node = {}, index) => {
    const knowledgeNodes = parseJSONArrayOrDefault(node.knowledgeNodesJson, node.knowledgeNodes)
    const scriptSegments = parseJSONArrayOrDefault(node.scriptSegmentsJson, node.scriptSegments)
    return {
      ...node,
      sortOrder: Number(node.sortOrder) || index + 1,
      knowledgeNodes,
      scriptSegments,
      knowledgeNodesJson: JSON.stringify(knowledgeNodes),
      scriptSegmentsJson: JSON.stringify(scriptSegments)
    }
  })
}

const parseJSONArrayOrDefault = (raw, fallback) => {
  if (Array.isArray(fallback) && fallback.length > 0) {
    return fallback
  }
  const text = String(raw ?? '').trim()
  if (!text) return []
  try {
    const parsed = JSON.parse(text)
    if (Array.isArray(parsed)) {
      return parsed
    }
    return [parsed]
  } catch {
    return []
  }
}

const validateNodeJSONFields = (nodes) => {
  const list = Array.isArray(nodes) ? nodes : []
  const validNodeIdSet = new Set(list.map(item => String(item?.nodeId || '').trim()).filter(Boolean))
  const validSegmentIdSet = buildSegmentIdSetFromScript(currentScript.value)

  for (let i = 0; i < list.length; i += 1) {
    const node = list[i] || {}
    const nodeLabel = node.title || node.nodeId || `节点${i + 1}`

    const knowledgeError = validateJSONArray(node.knowledgeNodesJson, '知识节点 JSON')
    if (knowledgeError) {
      return `${nodeLabel} 的知识节点 JSON 无效：${knowledgeError}`
    }

    const knowledgeRefError = validateKnowledgeReferences(node.knowledgeNodesJson, validNodeIdSet, validSegmentIdSet)
    if (knowledgeRefError) {
      return `${nodeLabel} 的知识节点 JSON 引用无效：${knowledgeRefError}`
    }

    const segmentsError = validateJSONArray(node.scriptSegmentsJson, '讲稿段落映射 JSON')
    if (segmentsError) {
      return `${nodeLabel} 的讲稿段落映射 JSON 无效：${segmentsError}`
    }

    const segmentRefError = validateSegmentReferences(node.scriptSegmentsJson, validNodeIdSet, validSegmentIdSet)
    if (segmentRefError) {
      return `${nodeLabel} 的讲稿段落映射 JSON 引用无效：${segmentRefError}`
    }
  }
  return ''
}

const buildSegmentIdSetFromScript = (scriptText) => {
  const segments = String(scriptText || '')
    .split(/\n{2,}|(?<=[。！？])\s*/)
    .map(item => item.trim())
    .filter(Boolean)
  const ids = segments.map((_, index) => `seg_${index + 1}`)
  return new Set(ids)
}

const validateSegmentReferences = (rawValue, validNodeIdSet, validSegmentIdSet) => {
  const text = String(rawValue ?? '').trim()
  if (!text) return ''
  let list = []
  try {
    const parsed = JSON.parse(text)
    list = Array.isArray(parsed) ? parsed : [parsed]
  } catch {
    return ''
  }

  const seenSegmentIds = new Set()
  for (let i = 0; i < list.length; i += 1) {
    const item = list[i] || {}
    const segmentId = String(item.segment_id || '').trim()
    if (!segmentId) {
      return `第 ${i + 1} 项缺少 segment_id`
    }
    if (seenSegmentIds.has(segmentId)) {
      return `segment_id ${segmentId} 在同一节点内重复`
    }
    seenSegmentIds.add(segmentId)
    if (validSegmentIdSet.size > 0 && !validSegmentIdSet.has(segmentId)) {
      return `segment_id ${segmentId} 不在当前讲稿段落范围内`
    }

    const nodeIDs = Array.isArray(item.node_ids) ? item.node_ids : []
    if (nodeIDs.length === 0) {
      return `segment_id ${segmentId} 缺少 node_ids`
    }
    for (let j = 0; j < nodeIDs.length; j += 1) {
      const nodeId = String(nodeIDs[j] || '').trim()
      if (!nodeId) {
        return `segment_id ${segmentId} 包含空 node_id`
      }
      if (validNodeIdSet.size > 0 && !validNodeIdSet.has(nodeId)) {
        return `segment_id ${segmentId} 引用了不存在的 node_id ${nodeId}`
      }
    }
  }
  return ''
}

const validateKnowledgeReferences = (rawValue, validNodeIdSet, validSegmentIdSet) => {
  const text = String(rawValue ?? '').trim()
  if (!text) return ''
  let list = []
  try {
    const parsed = JSON.parse(text)
    list = Array.isArray(parsed) ? parsed : [parsed]
  } catch {
    return ''
  }

  for (let i = 0; i < list.length; i += 1) {
    const item = list[i] || {}
    const nodeId = String(item.node_id || '').trim()
    if (!nodeId) {
      return `第 ${i + 1} 项缺少 node_id`
    }
    if (validNodeIdSet.size > 0 && !validNodeIdSet.has(nodeId)) {
      return `node_id ${nodeId} 不在当前页节点范围内`
    }

    const prerequisites = Array.isArray(item.prerequisites) ? item.prerequisites : []
    for (let j = 0; j < prerequisites.length; j += 1) {
      const prereq = String(prerequisites[j] || '').trim()
      if (!prereq) {
        return `node_id ${nodeId} 的 prerequisites 包含空值`
      }
      if (validNodeIdSet.size > 0 && !validNodeIdSet.has(prereq)) {
        return `node_id ${nodeId} 的 prerequisites 引用了不存在节点 ${prereq}`
      }
    }

    const coverageSpan = Array.isArray(item.coverage_span) ? item.coverage_span : []
    for (let j = 0; j < coverageSpan.length; j += 1) {
      const segId = String(coverageSpan[j] || '').trim()
      if (!segId) {
        return `node_id ${nodeId} 的 coverage_span 包含空值`
      }
      if (validSegmentIdSet.size > 0 && !validSegmentIdSet.has(segId)) {
        return `node_id ${nodeId} 的 coverage_span 引用了不存在段落 ${segId}`
      }
    }
  }

  return ''
}

const validateJSONArray = (value, fieldLabel) => {
  const raw = String(value ?? '').trim()
  if (!raw) return ''
  try {
    const parsed = JSON.parse(raw)
    if (!Array.isArray(parsed)) {
      return `${fieldLabel} 必须是数组格式（[]）`
    }
    return ''
  } catch {
    return `${fieldLabel} 不是合法 JSON`
  }
}

const generateAIScript = async () => {
  scriptGenerating.value = true
  try {
    currentScript.value = 'AI正在生成讲稿...'
    const data = await teacherV1Api.coursewares.generateScript({
      courseId: currentCourseId.value,
      pageNum: currentEditPage.value
    })
    currentScript.value = data?.data?.content || 'AI生成失败，请重试'
    const payload = data?.data || {}
    if (Array.isArray(payload.nodes) && payload.nodes.length > 0) {
      currentScriptNodes.value = normalizeNodesFromServer(payload.nodes)
    }
    if (payload.mappingCoverage) {
      currentMappingCoverage.value = payload.mappingCoverage
    }
    await loadScript(currentCourseId.value, currentEditPage.value)
  } catch (err) {
    currentScript.value = '生成失败：' + err.message
  } finally {
    scriptGenerating.value = false
  }
}

// --- 上传相关 ---
const handleFileSelect = (event) => {
  const file = event.target.files?.[0]
  if (file) selectedFileName.value = file.name
}

const uploadCourseware = async () => {
  const file = fileInput.value?.files?.[0]
  if (!file) return

  uploadLoading.value = true
  const formData = new FormData()
  formData.append('file', file)
  formData.append('title', file.name.replace(/\.[^.]+$/, ''))

  try {
    await teacherV1Api.coursewares.upload(formData)
    alert('课件上传成功！AI 解析已在后台执行，稍后可查看讲稿内容。')
    showUploadModal.value = false
    selectedFileName.value = ''
    await loadCoursewareList()
    activeTab.value = 'preview' // 上传后切到预览页
  } catch (err) {
    alert('上传失败：' + (err.message || '未知错误，请检查后端服务是否正常'))
  } finally {
    uploadLoading.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}

// --- 发布相关 ---
const publishCourseware = async () => {
  try {
    await teacherV1Api.coursewares.publish({
      courseId: currentCourseId.value,
      scope: publishScope.value
    })
    const course = coursewareList.value.find(c => c.id === currentCourseId.value)
    if (course) course.published = true
    alert('课件发布成功！学生端已可查看。')
    showPublishModal.value = false
  } catch (err) {
    alert('发布失败：' + err.message)
  }
}

// --- 数据加载（统计/卡点/提问） ---
const loadStudentStats = async (courseId) => {
  try {
    const data = await teacherV1Api.analytics.getStats(courseId)
    const payload = data?.data || {}
    studentStats.value = {
      totalQuestions: payload.totalQuestions || 0,
      hotPages: (payload.pageStats || []).map(item => item.page).slice(0, 3),
      keyDifficulties: (payload.keywords || []).map(item => item.word).slice(0, 3).join('、') || '暂无',
      nodeStats: Array.isArray(payload.nodeStats) ? payload.nodeStats : [],
      mappingCoverage: payload.mappingCoverage || null,
      activeSessions: payload.activeSessions || 0,
      avgTurnsPerSession: Number(payload.avgTurnsPerSession || 0),
      nodeHeatmap: Array.isArray(payload.nodeHeatmap) ? payload.nodeHeatmap : [],
      masteryRadar: payload.masteryRadar || { indicators: [], values: [], avgMastery: 0 },
      classTrend: Array.isArray(payload.classTrend) ? payload.classTrend : [],
      learningInsights: payload.learningInsights || { reteachNodes: [], prerequisiteGaps: [], summary: '' }
    }
  } catch (err) {
    studentStats.value = {
      totalQuestions: 0,
      hotPages: [],
      keyDifficulties: '加载失败',
      nodeStats: [],
      mappingCoverage: null,
      activeSessions: 0,
      avgTurnsPerSession: 0,
      nodeHeatmap: [],
      masteryRadar: { indicators: [], values: [], avgMastery: 0 },
      classTrend: [],
      learningInsights: { reteachNodes: [], prerequisiteGaps: [], summary: '' }
    }
  }
}

const loadCardData = async (courseId) => {
  try {
    const data = await teacherV1Api.analytics.getCardData(courseId)
    const pageStats = data?.data?.pageStats || []
    cardData.value = pageStats.map(item => ({
      page: item.page,
      提问量: item.questionCount,
      停留时长: item.stayTime,
      卡点指数: item.cardIndex
    }))
  } catch (err) {
    cardData.value = []
  }
}

const loadQuestionRecords = async (courseId) => {
  try {
    const data = await teacherV1Api.analytics.getQuestionRecords(courseId, 1, 100)
    const list = data?.data?.list || []
    questionRecords.value = list.map((item) => ({
      id: item.id,
      studentId: item.user_id || item.userId || '未知',
      page: item.page_index || item.pageIndex || 1,
      nodeId: item.node_id || item.nodeId || '',
      nodeTitle: item.node_title || item.nodeTitle || '',
      content: item.question || '',
      answer: item.answer || '',
      time: item.created_at ? new Date(item.created_at).toLocaleString() : ''
    }))
  } catch (err) {
    questionRecords.value = []
  }
}

const focusQuestionNode = async (record) => {
  if (!record) return
  const page = Number(record.page || 1) || 1
  const targetNodeId = String(record.nodeId || '').trim()
  activeTab.value = 'script'
  if (page !== currentEditPage.value) {
    await selectEditPage(page)
  }
  if (targetNodeId) {
    focusNodeRequest.value = { nodeId: targetNodeId, at: Date.now() }
  }
  ElMessage.success(`已定位到第 ${page} 页，可继续按节点编辑讲稿`)
}
</script>

<style scoped>
.teacher-app {
  width: 100%;
  height: 100vh;
  overflow: hidden;
  font-family: 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  background: radial-gradient(circle at 12% 8%, #f5fbf8 0%, #edf3ef 45%, #e8efeb 100%);
  padding: 14px;
  box-sizing: border-box;
}

.workspace-shell {
  width: 100%;
  height: 100%;
  border-radius: 28px;
  overflow: hidden;
  background: #f7faf8;
  border: 1px solid #d8e4dc;
  box-shadow: 0 24px 48px rgba(45, 72, 66, 0.08);
}

.main-content {
  display: flex;
  height: calc(100% - 56px);
}
.editor-section {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
/* 左侧菜单导航栏 */
.left-sidebar-menu {
  flex: 0 0 180px;
  background: #ffffff;
  border-right: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  transition: flex-basis 0.3s ease;
}

.left-sidebar-menu.collapsed {
  flex-basis: 64px;
}

.menu-header {
  padding: 20px 20px 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  font-weight: 700;
  color: #64748b;
  letter-spacing: 1.5px;
  font-family: monospace;
}

.left-sidebar-menu.collapsed .menu-header {
  padding: 20px 0 10px;
  justify-content: center;
}

.menu-toggle-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  color: #94a3b8;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
  border-radius: 4px;
  transition: background 0.2s, color 0.2s;
}

.menu-toggle-btn:hover {
  background: #f1f5f9;
  color: #1e293b;
}

.menu-toggle-btn svg {
  width: 18px;
  height: 18px;
}

.menu-list {
  display: flex;
  flex-direction: column;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  cursor: pointer;
  color: #334155;
  transition: all 0.2s;
  border-right: 3px solid transparent;
  font-size: 15px;
  white-space: nowrap;
  overflow: hidden;
}

.left-sidebar-menu.collapsed .menu-item {
  padding: 16px 0;
  justify-content: center;
}

.menu-item.active {
  color: #2F605A;
  background: #F4F7F7;
  border-right-color: #2F605A;
  font-weight: 600;
}

.menu-item small {
  color: #94a3b8;
  font-size: 11px;
}

.ins-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  stroke: currentColor;
}

.tab-container {
  flex: 1;
  min-height: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}


.empty-tip-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 14px;
  border: 1px dashed #cbd5e1;
  margin: 20px;
}

.empty-tip {
  color: #94a3b8;
  font-size: 15px;
  text-align: center;
}

</style>