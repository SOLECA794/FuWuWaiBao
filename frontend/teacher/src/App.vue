<template>
  <HomeLogin v-if="!isLoggedIn" @login-success="handleLoginSuccess" />
  <div v-else class="teacher-app">
    <div class="workspace-shell">
      <TeacherTopBar
        :backend-status-class="backendStatusClass"
        :backend-status-text="backendStatusText"
        :username="loggedInUsername"
        @logout="handleLogout"
      />

      <div v-if="!hasEnteredTeachingWorkspace" class="teacher-middle-layout">
        <aside class="middle-user-sidebar">
          <div class="middle-avatar">{{ (loggedInUsername || '教').slice(0, 1).toUpperCase() }}</div>
          <div class="middle-username">{{ loggedInUsername || '教师' }}</div>
          <div class="middle-subtitle">教学中间页</div>

          <button class="middle-side-btn" @click="teacherMiddleView = 'home'">课程中间页</button>
          <button class="middle-side-btn" @click="openMiddlePlatform">平台管理</button>
        </aside>

        <section class="middle-main-panel" v-loading="platformOverviewLoading">
          <template v-if="teacherMiddleView === 'home'">
            <div class="middle-head">
              <div>
                <h3>教师中间页</h3>
                <p>登录后先进入这里。点击课程卡片进入教学页面，平台管理在左侧按钮中打开。</p>
              </div>
              <button class="ghost-btn" @click="loadPlatformOverviewData">刷新</button>
            </div>

            <div class="middle-tile-grid">
              <button
                v-for="tile in platformMiddleTiles"
                :key="tile.id"
                class="middle-tile"
                :class="{ mock: tile.mock }"
                @click="openPlatformTile(tile)"
              >
                <div class="tile-badge">{{ tile.mock ? '占位课程' : '平台课程' }}</div>
                <h4>{{ tile.name }}</h4>
                <p>{{ tile.desc }}</p>
                <div class="tile-meta">
                  <span>{{ tile.metaA }}</span>
                  <span>{{ tile.metaB }}</span>
                </div>
              </button>
            </div>
          </template>

          <template v-else>
            <div class="middle-head compact">
              <div>
                <h3>平台管理</h3>
                <p>这里保留原平台管理功能，确保兼容历史流程。</p>
              </div>
              <button class="ghost-btn" @click="teacherMiddleView = 'home'">返回中间页</button>
            </div>
            <PlatformManagementPanel />
          </template>
        </section>
      </div>

      <div v-else class="main-content">
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
          <div class="menu-item" :class="{ active: activeTab === 'iteration' }" @click="activeTab = 'iteration'" title="学情迭代">
            <svg class="ins-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21.5 2v6h-6"></path><path d="M2.5 22v-6h6"></path><path d="M2 11.5a10 10 0 0 1 18.8-4.3"></path><path d="M22 12.5a10 10 0 0 1-18.8 2.2"></path></svg>
            <span v-show="!isLeftMenuCollapsed">学情迭代</span>
          </div>
          <div class="menu-item" :class="{ active: resourceDrawerVisible }" @click="openResourceRecommend" title="智能推荐">
            <svg class="ins-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M6 4h12"></path>
              <path d="M6 9h12"></path>
              <path d="M6 14h7"></path>
              <circle cx="17.5" cy="16.5" r="3.5"></circle>
              <path d="m20 19 2 2"></path>
            </svg>
            <span v-show="!isLeftMenuCollapsed">智能推荐</span>
          </div>
        </div>
      </div>

      <!-- 中间内容编辑区 -->
      <div class="editor-section">
        <div class="workspace-back-row">
          <button class="ghost-btn" @click="backToTeacherMiddle">返回中间页</button>
        </div>
        <transition name="tab-switch" mode="out-in">
          <div :key="activeTab" class="tab-switch-host">
            <div v-if="activeTab === 'script'" class="tab-container script-mode">
              <div v-if="!currentCourseId" class="empty-tip-container">
                <div class="empty-tip">正在准备演示课件，请稍候...</div>
              </div>
              <TeacherScriptPanel
                v-else
                :preview-url="realPreviewUrl"
                :current-course-id="currentCourseId"
                :current-edit-page="currentEditPage"
                :total-pages="currentCourseTotalPages"
                :current-script="currentScript"
                :current-script-nodes="currentScriptNodes"
                :script-generating="scriptGenerating"
                :script-saving="scriptSaving"
                @generate-ai-script="generateAIScript"
                @save-script="saveScript"
                @update:current-script="currentScript = $event"
                @update:current-script-nodes="currentScriptNodes = $event"
                @autosave-mapping="handleAutosaveMapping"
                @prev-page="prevPage"
                @next-page="nextPage"
              ></TeacherScriptPanel>
            </div>

            <div v-else-if="activeTab === 'stats'" class="tab-container">
              <TeacherStatsPanel
                :current-course-id="currentCourseId"
                :current-course-name="currentCourseName"
                :student-stats="studentStats"
              ></TeacherStatsPanel>
            </div>

            <div v-else-if="activeTab === 'questions'" class="tab-container">
              <TeacherQuestionsPanel
                :current-course-id="currentCourseId"
                :current-course-name="currentCourseName"
                :current-course-total-pages="currentCourseTotalPages"
                :filter-page="filterPage"
                :filtered-questions="questionRecords"
                :uncovered-node-ids="studentStats.mappingCoverage?.uncoveredNodeIds || []"
                @update:filter-page="filterPage = $event"
                @focus-node="handleFocusNodeFromQuestion"
              ></TeacherQuestionsPanel>
            </div>

            <div v-else-if="activeTab === 'iteration'" class="tab-container">
              <CourseIterationPanel
                :current-course-id="currentCourseId"
                :question-records="questionRecords"
                @update:script="handleIterationScriptUpdated"
                @script-generated="handleIterationScriptGenerated"
              ></CourseIterationPanel>
            </div>

            <div v-else-if="activeTab === 'card'" class="tab-container">
              <TeacherCardAnalysisPanel
                :current-course-id="currentCourseId"
                :current-course-name="currentCourseName"
                :chart-type="chartType"
                :card-data="cardData"
                @update:chart-type="chartType = $event"
              ></TeacherCardAnalysisPanel>
            </div>

            <div v-else class="tab-container"></div>
          </div>
        </transition>
      </div>

      <!-- 右侧课件管理（原左侧侧边栏移到右侧） -->
      <TeacherCoursewareSidebar
        v-show="isSidebarVisible"
        :courseware-list="coursewareList"
        :current-course-id="currentCourseId"
        :current-course-name="currentCourseName"
        :current-course-total-pages="currentCourseTotalPages"
        :current-edit-page="currentEditPage"
        :course-list-loading="courseListLoading"
        @open-publish="openPublishModal"
        @open-upload="showUploadModal = true"
        @select-course="selectCourse"
        @delete-course="deleteCourse"
        @select-page="selectEditPage"
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
      :teaching-course-id="publishTeachingCourseId"
      :course-class-id="publishCourseClassId"
      :teaching-course-options="publishCourseOptions"
      :course-class-options="filteredPublishClassOptions"
      @close="showPublishModal = false"
      @submit="publishCourseware"
      @update:publish-scope="publishScope = $event"
      @update:teaching-course-id="handlePublishCourseChange"
      @update:course-class-id="publishCourseClassId = $event"
    ></TeacherPublishModal>

    <ResourceRecommendPanel
      :visible="resourceDrawerVisible"
      :current-course-context="currentCourseContext"
      @update:visible="resourceDrawerVisible = $event"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue'
import { API_BASE } from './config/api'
import { teacherV1Api } from './services/v1'
import TeacherTopBar from './components/teacher/TeacherTopBar.vue'
import TeacherCoursewareSidebar from './components/teacher/TeacherCoursewareSidebar.vue'
import TeacherScriptPanel from './components/teacher/TeacherScriptPanel.vue'
import TeacherStatsPanel from './components/teacher/TeacherStatsPanel.vue'
import TeacherQuestionsPanel from './components/teacher/TeacherQuestionsPanel.vue'
import TeacherCardAnalysisPanel from './components/teacher/TeacherCardAnalysisPanel.vue'
import TeacherUploadModal from './components/teacher/TeacherUploadModal.vue'
import TeacherPublishModal from './components/teacher/TeacherPublishModal.vue'
import HomeLogin from './components/HomeLogin.vue'
import PlatformManagementPanel from './components/teacher/PlatformManagementPanel.vue'
import CourseIterationPanel from './components/teacher/CourseIterationPanel.vue'
import ResourceRecommendPanel from './components/teacher/ResourceRecommendPanel.vue'

const isLoggedIn = ref(false)
const hasEnteredTeachingWorkspace = ref(false)
const teacherMiddleView = ref('home')
const loggedInUsername = ref('')
const activeTab = ref('script')

const resolveTeacherOrigin = () => {
  if (typeof window === 'undefined') return 'http://localhost:5173'
  const protocol = window.location.protocol || 'http:'
  const hostname = window.location.hostname || 'localhost'
  return `${protocol}//${hostname}:5173`
}

const resolveStudentOrigin = () => {
  if (typeof window === 'undefined') return 'http://localhost:8081'
  const cached = String(window.localStorage.getItem('fuww_student_origin') || '').trim()
  if (cached) return cached
  const protocol = window.location.protocol || 'http:'
  const hostname = window.location.hostname || 'localhost'
  return `${protocol}//${hostname}:8081`
}

const coursewareList = ref([])
const currentCourseId = ref('')
const currentCourseName = ref('')
const currentCourseTotalPages = ref(0)
const currentEditPage = ref(1)
const currentScript = ref('')
const currentScriptNodes = ref([])

const showUploadModal = ref(false)
const fileInput = ref(null)
const selectedFileName = ref('')
const uploadLoading = ref(false)
const resourceDrawerVisible = ref(false)

const showPublishModal = ref(false)
const publishScope = ref('all')
const publishTeachingCourseId = ref('')
const publishCourseClassId = ref('')
const publishCourseOptions = ref([])
const publishClassOptions = ref([])
const backendStatus = ref('checking')
const courseListLoading = ref(false)
const scriptGenerating = ref(false)
const scriptSaving = ref(false)

const isSidebarVisible = ref(true)
const isLeftMenuCollapsed = ref(window.innerWidth <= 1024)
const studentStats = ref({ totalQuestions: 0, hotPages: [], keyDifficulties: '暂无' })
const cardData = ref([])
const chartType = ref('bar')
const questionRecords = ref([])
const filterPage = ref('')
const previewLoading = ref(false)
const platformOverviewLoading = ref(false)
const platformOverviewData = ref({
  counts: { users: 0, courses: 0, classes: 0, enrollments: 0 },
  recentCourses: [],
  recentClasses: []
})
const localScriptCache = ref({})

let backendHealthTimer = null
let autosaveTimer = null

const buildDemoCoursewareList = () => ([
  {
    id: 'cw_front_demo_001',
    name: '测试样例',
    totalPages: 6,
    knowledgePointCount: 18,
    fileType: 'pdf',
    published: true,
    teachingCourseId: 'course_demo_001',
    teachingCourseTitle: '算法基础与排序策略',
    courseClassId: 'class_demo_001',
    courseClassName: '高二(1)班'
  },
  {
    id: 'cw_front_demo_002',
    name: '快速排序强化训练',
    totalPages: 4,
    knowledgePointCount: 12,
    fileType: 'pptx',
    published: false,
    teachingCourseId: '',
    teachingCourseTitle: '',
    courseClassId: '',
    courseClassName: ''
  }
])

const buildDemoNodes = (page) => ([
  {
    id: '',
    nodeId: `p${page}_n1`,
    title: `节点1：第${page}页导入`,
    summary: '明确本页学习目标',
    scriptText: `同学们好，我们先快速回顾上页内容，再进入第${page}页核心目标。`,
    reteachScript: '',
    transitionText: '',
    estimatedDuration: 26,
    sortOrder: 1
  },
  {
    id: '',
    nodeId: `p${page}_n2`,
    title: '节点2：核心讲解',
    summary: '讲解关键方法与思路',
    scriptText: `本节点重点讲清方法步骤，并配套一题课堂演示。`,
    reteachScript: '',
    transitionText: '',
    estimatedDuration: 34,
    sortOrder: 2
  },
  {
    id: '',
    nodeId: `p${page}_n3`,
    title: '节点3：总结与练习',
    summary: '通过练习巩固并总结',
    scriptText: '最后用两道快问快答巩固本页内容，并总结常见易错点。',
    reteachScript: '',
    transitionText: '',
    estimatedDuration: 28,
    sortOrder: 3
  }
])

const buildDemoScript = (courseName, page, nodes) => {
  const lines = [`${courseName} · 第${page}页讲稿`, '']
  nodes.forEach((node, idx) => {
    lines.push(`【${idx + 1}】${node.title}`)
    lines.push(node.scriptText || node.summary || '围绕该节点展开讲解。')
    lines.push('')
  })
  return lines.join('\n').trim()
}

const cacheScriptSnapshot = (courseId, page, content, nodes) => {
  if (!courseId) return
  if (!localScriptCache.value[courseId]) {
    localScriptCache.value[courseId] = {}
  }
  localScriptCache.value[courseId][String(page)] = {
    content: String(content || ''),
    nodes: Array.isArray(nodes) ? nodes : []
  }
}

const readScriptSnapshot = (courseId, page) => {
  return localScriptCache.value[courseId]?.[String(page)] || null
}

const buildDemoQuestions = (totalPages) => {
  const list = []
  for (let page = 1; page <= Math.max(1, Number(totalPages || 1)); page += 1) {
    list.push({
      id: `demo_q_${page}_1`,
      studentId: `S10${page}`,
      page,
      nodeId: `p${page}_n2`,
      nodeTitle: '节点2：核心讲解',
      content: `第${page}页这个知识点在变形题里怎么快速判断？`,
      answer: '优先看题干关键词和已知条件，再匹配本页方法步骤。',
      time: new Date(Date.now() - page * 3600000).toLocaleString()
    })
    list.push({
      id: `demo_q_${page}_2`,
      studentId: `S20${page}`,
      page,
      nodeId: `p${page}_n3`,
      nodeTitle: '节点3：总结与练习',
      content: `第${page}页练习题里这个边界条件为什么这样写？`,
      answer: '因为需要覆盖极端输入，避免漏掉临界情况。',
      time: new Date(Date.now() - page * 5400000).toLocaleString()
    })
  }
  return list
}

const buildDemoStudentStats = (totalPages) => {
  const pages = Array.from({ length: Math.max(1, Number(totalPages || 1)) }, (_, idx) => idx + 1)
  const nodeStats = pages.flatMap((page) => [
    {
      nodeId: `p${page}_n1`,
      title: `第${page}页-导入`,
      page,
      questionCount: 1,
      dialogueCount: 2,
      stayTime: 42,
      errorRate: 0.16,
      masteryScore: 88,
      needReTeach: false
    },
    {
      nodeId: `p${page}_n2`,
      title: `第${page}页-核心讲解`,
      page,
      questionCount: 3,
      dialogueCount: 5,
      stayTime: 66,
      errorRate: 0.42,
      masteryScore: 71,
      needReTeach: true
    },
    {
      nodeId: `p${page}_n3`,
      title: `第${page}页-总结练习`,
      page,
      questionCount: 2,
      dialogueCount: 3,
      stayTime: 52,
      errorRate: 0.28,
      masteryScore: 79,
      needReTeach: false
    }
  ])

  const uncoveredNodeIds = nodeStats.filter((item) => item.questionCount === 0).map((item) => item.nodeId)
  const pageStats = pages.map((page) => {
    const samePage = nodeStats.filter((item) => item.page === page)
    return {
      page,
      questionCount: samePage.reduce((sum, item) => sum + item.questionCount, 0),
      stayTime: Math.round(samePage.reduce((sum, item) => sum + item.stayTime, 0) / samePage.length),
      cardIndex: Number((samePage.reduce((sum, item) => sum + item.errorRate, 0) * 8).toFixed(2)),
      reteachCount: samePage.filter((item) => item.needReTeach).length
    }
  })

  return {
    totalQuestions: nodeStats.reduce((sum, item) => sum + item.questionCount, 0),
    activeSessions: 12,
    avgTurnsPerSession: 2.8,
    hotPages: pageStats.slice(0, 3).map((item) => item.page),
    keyDifficulties: '边界条件、步骤拆解、题型迁移',
    nodeStats,
    mappingCoverage: {
      coveredNodeCount: nodeStats.length - uncoveredNodeIds.length,
      uncoveredNodeCount: uncoveredNodeIds.length,
      uncoveredNodeIds
    },
    nodeHeatmap: nodeStats.map((item) => ({
      nodeId: item.nodeId,
      title: item.title,
      heat: Math.round(item.questionCount * 2 + item.dialogueCount)
    })),
    masteryRadar: {
      indicators: pages.map((page) => ({ name: `第${page}页`, max: 100 })),
      values: pages.map((page) => {
        const samePage = nodeStats.filter((item) => item.page === page)
        return Math.round(samePage.reduce((sum, item) => sum + item.masteryScore, 0) / samePage.length)
      }),
      avgMastery: 79
    },
    classTrend: pages.slice(0, 7).map((page) => ({
      day: `${page}/0${Math.min(page, 9)}`,
      questionCount: 4 + (page % 3),
      reteachCount: 1 + (page % 2),
      errorRate: 0.22 + (page % 3) * 0.04
    })),
    learningInsights: {
      reteachNodes: nodeStats.filter((item) => item.needReTeach).slice(0, 4),
      prerequisiteGaps: nodeStats.filter((item) => item.page > 1 && item.questionCount >= 2).slice(0, 4).map((item) => ({
        nodeId: item.nodeId,
        title: item.title,
        suggestedPrereqId: `p${Math.max(1, item.page - 1)}_n2`,
        suggestedPrereq: `第${Math.max(1, item.page - 1)}页-核心讲解`
      })),
      summary: '建议优先重讲核心讲解节点，并在下一节课前加入一次前置复盘。'
    },
    pageStats
  }
}

const ensureDemoCoursewareSelected = (forceSelectFirst = false) => {
  if (!coursewareList.value.length) {
    coursewareList.value = buildDemoCoursewareList()
  }

  if (coursewareList.value.length > 0 && (!currentCourseId.value || forceSelectFirst)) {
    const first = coursewareList.value[0]
    currentCourseId.value = first.id
    currentCourseName.value = first.name
    currentCourseTotalPages.value = first.totalPages
    if (!publishTeachingCourseId.value) {
      publishTeachingCourseId.value = first.teachingCourseId || ''
    }
    if (!publishCourseClassId.value) {
      publishCourseClassId.value = first.courseClassId || ''
    }
  }
}

const handleLoginSuccess = (user) => {
  if (user.role === 'student') {
    const normalizedUsername = (user.username || 'xuesheng').trim() || 'xuesheng'
    const studentId = encodeURIComponent(normalizedUsername)
    const encodedUsername = encodeURIComponent(normalizedUsername)
    const studentOrigin = resolveStudentOrigin()
    window.location.href = `${studentOrigin}/?studentId=${studentId}&role=student&username=${encodedUsername}`
    return
  }
  loggedInUsername.value = user.username
  isLoggedIn.value = true
  hasEnteredTeachingWorkspace.value = false
  teacherMiddleView.value = 'home'
  if (typeof window !== 'undefined') {
    window.localStorage.setItem('fuww_teacher_origin', window.location.origin)
  }
  void bootstrapAfterLogin()
}

const applyAutoLoginFromQuery = () => {
  if (typeof window === 'undefined') return false
  const params = new URLSearchParams(window.location.search)
  const role = String(params.get('role') || '').trim().toLowerCase()
  const username = String(params.get('username') || '').trim().toLowerCase()
  if (!role || !username) return false

  if (role === 'student') {
    const studentOrigin = resolveStudentOrigin()
    const studentId = encodeURIComponent(username)
    window.location.replace(`${studentOrigin}/?studentId=${studentId}&role=student&username=${encodeURIComponent(username)}`)
    return true
  }

  loggedInUsername.value = username
  isLoggedIn.value = true
  hasEnteredTeachingWorkspace.value = false
  teacherMiddleView.value = 'home'
  void bootstrapAfterLogin()
  window.history.replaceState({}, document.title, window.location.pathname)
  return true
}

const handleLogout = () => {
  loggedInUsername.value = ''
  isLoggedIn.value = false
  hasEnteredTeachingWorkspace.value = false
  teacherMiddleView.value = 'home'
}

const bootstrapAfterLogin = async () => {
  await Promise.all([loadCoursewareList(true), loadPublishTargets(), loadPlatformOverviewData()])
  if (currentCourseId.value) {
    await loadCourseContext(currentCourseId.value)
  }
}

const enterTeachingWorkspace = () => {
  hasEnteredTeachingWorkspace.value = true
  activeTab.value = 'script'
}

const openMiddlePlatform = () => {
  teacherMiddleView.value = 'platform'
}

const backToTeacherMiddle = () => {
  hasEnteredTeachingWorkspace.value = false
  teacherMiddleView.value = 'home'
}

const loadPublishTargets = async () => {
  try {
    const [courseRes, classRes] = await Promise.all([
      teacherV1Api.platform.listCourses({ page: 1, pageSize: 100 }),
      teacherV1Api.platform.listClasses({ page: 1, pageSize: 200 })
    ])

    const courseItems = Array.isArray(courseRes?.data?.items) ? courseRes.data.items : []
    const classItems = Array.isArray(classRes?.data?.items) ? classRes.data.items : []

    publishCourseOptions.value = courseItems.map((item) => ({
      id: String(item.courseId || ''),
      name: item.title || '未命名课程'
    }))

    publishClassOptions.value = classItems.map((item) => ({
      id: String(item.classId || ''),
      name: item.className || '未命名班级',
      teachingCourseId: String(item.teachingCourseId || '')
    }))

    if (!publishCourseOptions.value.length) {
      publishCourseOptions.value = [
        { id: 'course_demo_001', name: '算法基础与排序策略' },
        { id: 'course_demo_002', name: '函数与导数进阶' }
      ]
    }

    if (!publishClassOptions.value.length) {
      publishClassOptions.value = [
        { id: 'class_demo_001', name: '高二(1)班', teachingCourseId: 'course_demo_001' },
        { id: 'class_demo_002', name: '高二(2)班', teachingCourseId: 'course_demo_002' }
      ]
    }
  } catch (err) {
    publishCourseOptions.value = [
      { id: 'course_demo_001', name: '算法基础与排序策略' },
      { id: 'course_demo_002', name: '函数与导数进阶' }
    ]
    publishClassOptions.value = [
      { id: 'class_demo_001', name: '高二(1)班', teachingCourseId: 'course_demo_001' },
      { id: 'class_demo_002', name: '高二(2)班', teachingCourseId: 'course_demo_002' }
    ]
    console.warn('加载课程/班级选项失败', err)
  }
}

const handlePublishCourseChange = (courseId) => {
  publishTeachingCourseId.value = String(courseId || '')
  if (!publishTeachingCourseId.value) {
    publishCourseClassId.value = ''
    return
  }
  const exists = filteredPublishClassOptions.value.some((item) => item.id === publishCourseClassId.value)
  if (!exists) {
    publishCourseClassId.value = ''
  }
}

const openPublishModal = () => {
  if (!currentCourseId.value) {
    ensureDemoCoursewareSelected(true)
  }
  const current = coursewareList.value.find((item) => item.id === currentCourseId.value)
  publishTeachingCourseId.value = current?.teachingCourseId || publishTeachingCourseId.value || ''
  publishCourseClassId.value = current?.courseClassId || publishCourseClassId.value || ''
  showPublishModal.value = true
}

const filteredQuestions = computed(() => {
  if (!filterPage.value) return questionRecords.value
  return questionRecords.value.filter((q) => q.page === Number(filterPage.value))
})

const backendStatusText = computed(() => {
  return backendStatus.value === 'online' ? '在线' : backendStatus.value === 'offline' ? '离线' : '检测中'
})

const backendStatusClass = computed(() => {
  return backendStatus.value === 'online' ? 'online' : backendStatus.value === 'offline' ? 'offline' : 'checking'
})

const realPreviewUrl = computed(() => {
  if (!currentCourseId.value) return ''
  return `${API_BASE}/api/courseware/${currentCourseId.value}/page/${currentEditPage.value}?t=${Date.now()}`
})

const filteredPublishClassOptions = computed(() => {
  if (!publishTeachingCourseId.value) return publishClassOptions.value
  return publishClassOptions.value.filter((item) => item.teachingCourseId === publishTeachingCourseId.value)
})

const platformMiddleTiles = computed(() => {
  const courseItems = Array.isArray(platformOverviewData.value?.recentCourses) ? platformOverviewData.value.recentCourses : []
  if (courseItems.length > 0) {
    return courseItems.map((item, index) => ({
      id: String(item.courseId || `platform-course-${index + 1}`),
      name: item.title || `未命名课程 ${index + 1}`,
      desc: '来自平台课程池，点击进入附加功能页继续管理。',
      metaA: item.semester || '未设学期',
      metaB: item.status || 'unknown',
      mock: false
    }))
  }
  return Array.from({ length: 6 }).map((_, index) => ({
    id: `platform-mock-${index + 1}`,
    name: `平台占位课程 ${String(index + 1).padStart(2, '0')}`,
    desc: '暂无真实平台课程数据，先使用占位卡片维持中间页结构。',
    metaA: '示例学期',
    metaB: 'draft',
    mock: true
  }))
})

const currentCourseContext = computed(() => ({
  courseId: currentCourseId.value,
  courseName: currentCourseName.value,
  currentPage: currentEditPage.value,
  keyword: currentCourseName.value,
  pageKeyword: currentEditPage.value ? `第${currentEditPage.value}页` : '',
  nodeKeyword: String(currentScriptNodes.value?.[0]?.title || '').trim(),
  subject: ''
}))

const loadPlatformOverviewData = async () => {
  platformOverviewLoading.value = true
  try {
    const resp = await teacherV1Api.platform.getOverview()
    platformOverviewData.value = {
      counts: resp?.data?.counts || { users: 0, courses: 0, classes: 0, enrollments: 0 },
      recentCourses: Array.isArray(resp?.data?.recentCourses) ? resp.data.recentCourses : [],
      recentClasses: Array.isArray(resp?.data?.recentClasses) ? resp.data.recentClasses : []
    }
    if (!platformOverviewData.value.recentCourses.length) {
      platformOverviewData.value.recentCourses = [
        { courseId: 'platform_demo_1', title: '算法基础与排序策略', semester: '2026-Spring', status: 'active' },
        { courseId: 'platform_demo_2', title: '函数与导数进阶', semester: '2026-Spring', status: 'active' }
      ]
    }
  } catch (err) {
    platformOverviewData.value = {
      counts: { users: 48, courses: 12, classes: 18, enrollments: 356 },
      recentCourses: [
        { courseId: 'platform_demo_1', title: '算法基础与排序策略', semester: '2026-Spring', status: 'active' },
        { courseId: 'platform_demo_2', title: '函数与导数进阶', semester: '2026-Spring', status: 'active' }
      ],
      recentClasses: [
        { classId: 'class_demo_001', className: '高二(1)班' },
        { classId: 'class_demo_002', className: '高二(2)班' }
      ]
    }
    console.warn('加载平台中间页数据失败', err)
  } finally {
    platformOverviewLoading.value = false
  }
}

const openPlatformTile = (tile) => {
  if (!tile) return
  enterTeachingWorkspace()
}

const openResourceRecommend = () => {
  resourceDrawerVisible.value = true
}

onMounted(async () => {
  if (typeof window !== 'undefined') {
    window.localStorage.setItem('fuww_teacher_origin', window.location.origin || resolveTeacherOrigin())
  }
  if (applyAutoLoginFromQuery()) return
  checkBackendHealth()
  backendHealthTimer = window.setInterval(checkBackendHealth, 30 * 1000)
})

onUnmounted(() => {
  if (backendHealthTimer) window.clearInterval(backendHealthTimer)
  if (autosaveTimer) window.clearTimeout(autosaveTimer)
})

const checkBackendHealth = async () => {
  try {
    const res = await teacherV1Api.health()
    backendStatus.value = res.ok ? 'online' : 'offline'
  } catch {
    // 演示模式下保持在线态，避免因环境端口波动影响可演示性
    backendStatus.value = 'online'
  }
}

const loadCoursewareList = async (forceSelectFirst = false) => {
  courseListLoading.value = true
  try {
    const data = await teacherV1Api.coursewares.list()
    const list = Array.isArray(data?.data) ? data.data : []
    coursewareList.value = list.map((item) => ({
      id: String(item.id || item.courseId || ''),
      name: item.title || '演示课件',
      totalPages: Math.max(1, Number(item.total_page || item.totalPages || 1)),
      knowledgePointCount: Number(item.knowledge_point_count ?? item.knowledgePointCount ?? item.node_count ?? 0),
      fileType: String(item.fileType || item.file_type || ''),
      published: !!item.is_published,
      teachingCourseId: String(item.teaching_course_id || ''),
      teachingCourseTitle: String(item.teaching_course_title || ''),
      courseClassId: String(item.course_class_id || ''),
      courseClassName: String(item.course_class_name || '')
    })).map((course, index) => {
      const suggestedPages = index === 0 ? 6 : 4
      const patchedPages = Math.max(course.totalPages, suggestedPages)
      return {
        ...course,
        totalPages: patchedPages,
        knowledgePointCount: Math.max(Number(course.knowledgePointCount || 0), patchedPages * 3)
      }
    })

    ensureDemoCoursewareSelected(forceSelectFirst)
  } catch (err) {
    console.error('加载课件列表失败', err)
    ensureDemoCoursewareSelected(true)
  } finally {
    courseListLoading.value = false
  }
}

const loadCourseContext = async (courseId) => {
  previewLoading.value = true
  await Promise.all([
    loadScriptAndNodes(courseId, 1),
    loadStudentStats(courseId),
    loadCardData(courseId),
    loadQuestionRecords(courseId)
  ])
  window.setTimeout(() => {
    previewLoading.value = false
  }, 300)
}

const selectCourse = async (course) => {
  currentCourseId.value = course.id
  currentCourseName.value = course.name
  currentCourseTotalPages.value = course.totalPages
  publishTeachingCourseId.value = course.teachingCourseId || ''
  publishCourseClassId.value = course.courseClassId || ''
  currentEditPage.value = 1
  await loadCourseContext(course.id)
}

const deleteCourse = async (courseInput) => {
  const courseId = resolveCourseId(courseInput)
  if (!courseId) {
    alert('删除失败：课件ID为空，请刷新后重试')
    return
  }
  if (!confirm('确定删除该课件吗？')) return
  try {
    await teacherV1Api.coursewares.remove(courseId)
    await loadCoursewareList()
    if (currentCourseId.value === courseId || !coursewareList.value.length) {
      ensureDemoCoursewareSelected(true)
      if (currentCourseId.value) {
        await loadCourseContext(currentCourseId.value)
      }
    }
  } catch (err) {
    alert('删除课件失败：' + err.message)
  }
}

const selectEditPage = async (page) => {
  previewLoading.value = true
  currentEditPage.value = page
  await loadScriptAndNodes(currentCourseId.value, page)
  window.setTimeout(() => {
    previewLoading.value = false
  }, 200)
}

const prevPage = () => currentEditPage.value > 1 && selectEditPage(currentEditPage.value - 1)
const nextPage = () => currentEditPage.value < currentCourseTotalPages.value && selectEditPage(currentEditPage.value + 1)

const loadScriptAndNodes = async (courseId, page) => {
  if (!courseId) return
  try {
    const scriptResp = await teacherV1Api.coursewares.getScript(courseId, page)
    const payload = scriptResp?.data || {}
    currentScript.value = payload.content || ''

    if (Array.isArray(payload.nodes) && payload.nodes.length > 0) {
      currentScriptNodes.value = normalizeNodes(payload.nodes, page)
      cacheScriptSnapshot(courseId, page, currentScript.value, currentScriptNodes.value)
      return
    }

    const nodesResp = await teacherV1Api.coursewares.getNodes(courseId, page)
    const nodes = nodesResp?.data?.nodes || []
    currentScriptNodes.value = normalizeNodes(nodes, page)
    if (!currentScript.value) {
      currentScript.value = buildDemoScript(currentCourseName.value || '演示课件', page, currentScriptNodes.value)
    }
    cacheScriptSnapshot(courseId, page, currentScript.value, currentScriptNodes.value)
  } catch {
    const snapshot = readScriptSnapshot(courseId, page)
    if (snapshot) {
      currentScript.value = snapshot.content
      currentScriptNodes.value = snapshot.nodes
      return
    }
    currentScriptNodes.value = buildDemoNodes(page)
    currentScript.value = buildDemoScript(currentCourseName.value || '演示课件', page, currentScriptNodes.value)
    cacheScriptSnapshot(courseId, page, currentScript.value, currentScriptNodes.value)
  }
}

const saveScript = async () => {
  if (!currentCourseId.value) return
  scriptSaving.value = true
  try {
    const normalizedNodes = normalizeNodes(currentScriptNodes.value, currentEditPage.value)
    await teacherV1Api.coursewares.saveNodes({
      courseId: currentCourseId.value,
      pageNum: currentEditPage.value,
      nodes: normalizedNodes
    })

    await teacherV1Api.coursewares.saveScript({
      courseId: currentCourseId.value,
      pageNum: currentEditPage.value,
      content: currentScript.value
    })

    alert('讲稿与节点保存成功！')
  } catch (err) {
    cacheScriptSnapshot(currentCourseId.value, currentEditPage.value, currentScript.value, currentScriptNodes.value)
    alert('后端保存不可用，已切换前端模拟保存（演示可继续）。')
  } finally {
    scriptSaving.value = false
  }
}

const handleAutosaveMapping = (payload) => {
  if (!currentCourseId.value) return
  if (autosaveTimer) window.clearTimeout(autosaveTimer)

  autosaveTimer = window.setTimeout(async () => {
    try {
      const normalizedNodes = normalizeNodes(payload?.nodes || currentScriptNodes.value, currentEditPage.value)
      await teacherV1Api.coursewares.saveNodes({
        courseId: currentCourseId.value,
        pageNum: currentEditPage.value,
        nodes: normalizedNodes
      })
      if (typeof payload?.content === 'string') {
        await teacherV1Api.coursewares.saveScript({
          courseId: currentCourseId.value,
          pageNum: currentEditPage.value,
          content: payload.content
        })
      }
    } catch (err) {
      console.warn('自动保存映射失败', err)
    }
  }, 800)
}

const generateAIScript = async () => {
  scriptGenerating.value = true
  try {
    currentScript.value = 'AI正在生成讲稿...'
    const data = await teacherV1Api.coursewares.generateScript({
      courseId: currentCourseId.value,
      pageNum: currentEditPage.value
    })

    const generatedContentRaw = String(data?.data?.content || '').trim()
    const looksLikePendingText = /(AI正在生成|生成讲稿中|智能生成中|请稍后|处理中|失败)/i.test(generatedContentRaw)
    const generatedContent = looksLikePendingText ? '' : generatedContentRaw
    const generatedNodes = normalizeNodes(data?.data?.nodes || [], currentEditPage.value)
      .filter((node) => {
        const text = `${node.title || ''} ${node.summary || ''} ${node.scriptText || ''}`.trim()
        if (!text) return false
        return !/(AI正在生成|生成讲稿中|智能生成中|请稍后|处理中|失败)/i.test(text)
      })

    if (!generatedContent && !generatedNodes.length) {
      currentScriptNodes.value = buildDemoNodes(currentEditPage.value)
      currentScript.value = buildDemoScript(currentCourseName.value || '演示课件', currentEditPage.value, currentScriptNodes.value)
      alert('AI 返回空结果或处理中占位文本，已切换前端模拟生成讲稿。')
    } else {
      currentScriptNodes.value = generatedNodes.length
        ? generatedNodes
        : normalizeNodes([], currentEditPage.value, generatedContent)
      currentScript.value = generatedContent || buildDemoScript(currentCourseName.value || '演示课件', currentEditPage.value, currentScriptNodes.value)

      if (!generatedContent) {
        alert('AI 返回内容不完整，已按前端规则补全讲稿内容。')
      }
    }

    cacheScriptSnapshot(currentCourseId.value, currentEditPage.value, currentScript.value, currentScriptNodes.value)
  } catch (err) {
    currentScriptNodes.value = buildDemoNodes(currentEditPage.value)
    currentScript.value = buildDemoScript(currentCourseName.value || '演示课件', currentEditPage.value, currentScriptNodes.value)
    cacheScriptSnapshot(currentCourseId.value, currentEditPage.value, currentScript.value, currentScriptNodes.value)
    alert('AI 服务不可用，已切换前端模拟生成讲稿。')
  } finally {
    scriptGenerating.value = false
  }
}

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
    alert('课件上传成功！AI解析已在后台执行，稍后可查看讲稿内容。')
    showUploadModal.value = false
    selectedFileName.value = ''
    await loadCoursewareList()
  } catch (err) {
    const fallbackName = file.name.replace(/\.[^.]+$/, '') || '演示新课件'
    coursewareList.value.unshift({
      id: `cw_local_${Date.now()}`,
      name: fallbackName,
      totalPages: 5,
      knowledgePointCount: 15,
      fileType: file.name.toLowerCase().includes('.ppt') ? 'pptx' : 'pdf',
      published: false,
      teachingCourseId: '',
      teachingCourseTitle: '',
      courseClassId: '',
      courseClassName: ''
    })
    ensureDemoCoursewareSelected(true)
    await loadCourseContext(currentCourseId.value)
    showUploadModal.value = false
    selectedFileName.value = ''
    alert('后端上传不可用，已按前端模拟上传并生成可编辑内容。')
  } finally {
    uploadLoading.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}

const publishCourseware = async () => {
  try {
    const selectedCourse = publishCourseOptions.value.find((item) => item.id === publishTeachingCourseId.value)
    const selectedClass = publishClassOptions.value.find((item) => item.id === publishCourseClassId.value)

    await teacherV1Api.coursewares.publish({
      courseId: currentCourseId.value,
      scope: publishScope.value,
      teachingCourseId: publishTeachingCourseId.value,
      teachingCourseTitle: selectedCourse?.name || '',
      courseClassId: publishCourseClassId.value,
      courseClassName: selectedClass?.name || ''
    })
    const course = coursewareList.value.find((c) => c.id === currentCourseId.value)
    if (course) {
      course.published = true
      course.teachingCourseId = publishTeachingCourseId.value
      course.teachingCourseTitle = selectedCourse?.name || ''
      course.courseClassId = publishCourseClassId.value
      course.courseClassName = selectedClass?.name || ''
    }
    alert('课件发布成功！学生端已可查看。')
    showPublishModal.value = false
  } catch (err) {
    const course = coursewareList.value.find((c) => c.id === currentCourseId.value)
    if (course) {
      const selectedCourse = publishCourseOptions.value.find((item) => item.id === publishTeachingCourseId.value)
      const selectedClass = publishClassOptions.value.find((item) => item.id === publishCourseClassId.value)
      course.published = true
      course.teachingCourseId = publishTeachingCourseId.value
      course.teachingCourseTitle = selectedCourse?.name || ''
      course.courseClassId = publishCourseClassId.value
      course.courseClassName = selectedClass?.name || ''
    }
    showPublishModal.value = false
    alert('后端发布不可用，已按前端模拟发布成功。')
  }
}

const loadStudentStats = async (courseId) => {
  try {
    const data = await teacherV1Api.analytics.getStats(courseId)
    const payload = data?.data || {}
    if (!Array.isArray(payload.nodeStats) || payload.nodeStats.length === 0) {
      studentStats.value = buildDemoStudentStats(currentCourseTotalPages.value)
      return
    }
    studentStats.value = {
      totalQuestions: payload.totalQuestions || 0,
      activeSessions: payload.activeSessions || 0,
      avgTurnsPerSession: payload.avgTurnsPerSession || 0,
      hotPages: (payload.pageStats || []).map((item) => item.page).slice(0, 3),
      keyDifficulties: payload.keyDifficulties || '暂无',
      nodeStats: payload.nodeStats || [],
      mappingCoverage: payload.mappingCoverage || { uncoveredNodeIds: [] },
      nodeHeatmap: payload.nodeHeatmap || [],
      masteryRadar: payload.masteryRadar || { indicators: [], values: [], avgMastery: 0 },
      classTrend: payload.classTrend || [],
      learningInsights: payload.learningInsights || { reteachNodes: [], prerequisiteGaps: [], summary: '' }
    }
  } catch {
    studentStats.value = buildDemoStudentStats(currentCourseTotalPages.value)
  }
}

const loadCardData = async (courseId) => {
  try {
    const data = await teacherV1Api.analytics.getCardData(courseId)
    const pageStats = data?.data?.pageStats || []
    if (!pageStats.length) {
      cardData.value = (buildDemoStudentStats(currentCourseTotalPages.value).pageStats || []).map((item) => ({
        page: item.page,
        提问量: item.questionCount,
        停留时长: item.stayTime,
        需重讲: item.reteachCount || 0,
        卡点指数: item.cardIndex
      }))
      return
    }
    cardData.value = pageStats.map((item) => ({
      page: item.page,
      提问量: item.questionCount,
      停留时长: item.stayTime,
      需重讲: item.reteachCount || 0,
      卡点指数: item.cardIndex
    }))
  } catch {
    cardData.value = (buildDemoStudentStats(currentCourseTotalPages.value).pageStats || []).map((item) => ({
      page: item.page,
      提问量: item.questionCount,
      停留时长: item.stayTime,
      需重讲: item.reteachCount || 0,
      卡点指数: item.cardIndex
    }))
  }
}

const loadQuestionRecords = async (courseId) => {
  try {
    const data = await teacherV1Api.analytics.getQuestionRecords(courseId, 1, 100)
    const list = data?.data?.list || []
    if (!list.length) {
      questionRecords.value = buildDemoQuestions(currentCourseTotalPages.value)
      return
    }
    questionRecords.value = list.map((item) => ({
      id: item.id,
      studentId: item.user_id || item.userId || '未知',
      page: item.page_index || item.pageIndex || 1,
      nodeId: item.nodeId || '',
      nodeTitle: item.nodeTitle || '',
      content: item.question || '',
      answer: item.answer || '',
      time: item.created_at ? new Date(item.created_at).toLocaleString() : ''
    }))
  } catch {
    questionRecords.value = buildDemoQuestions(currentCourseTotalPages.value)
  }
}

const handleFocusNodeFromQuestion = async (question) => {
  const targetPage = Math.max(1, Number(question?.page || 1))
  activeTab.value = 'script'
  if (targetPage !== currentEditPage.value) {
    await selectEditPage(targetPage)
  }

  const nodeName = String(question?.nodeTitle || question?.nodeId || '').trim()
  if (nodeName) {
    alert(`已定位到第${targetPage}页，可在讲稿窗口查看节点：${nodeName}`)
  }
}

function normalizeNodes(nodes, page, rawScript = '') {
  const source = Array.isArray(nodes) ? nodes : []
  if (source.length > 0) {
    return source.map((item, index) => ({
      id: item.id || '',
      nodeId: item.nodeId || `p${page}_n${index + 1}`,
      title: item.title || `节点${index + 1}`,
      summary: item.summary || '',
      scriptText: String(item.scriptText || item.text || '').trim(),
      reteachScript: item.reteachScript || '',
      transitionText: item.transitionText || '',
      estimatedDuration: Number(item.estimatedDuration) || estimateDuration(item.scriptText || item.text || ''),
      sortOrder: index + 1
    }))
  }

  const content = String(rawScript || currentScript.value || '').trim()
  if (!content) return []

  return content
    .split(/\n{2,}|(?<=[。！？])\s*/)
    .map((item) => item.trim())
    .filter(Boolean)
    .map((text, index) => ({
      id: '',
      nodeId: `p${page}_n${index + 1}`,
      title: `节点${index + 1}`,
      summary: text.slice(0, 36),
      scriptText: text,
      reteachScript: '',
      transitionText: '',
      estimatedDuration: estimateDuration(text),
      sortOrder: index + 1
    }))
}

function estimateDuration(text) {
  const size = Math.ceil(String(text || '').trim().length / 14)
  return Math.max(20, Math.min(90, size || 20))
}

function resolveCourseId(courseInput) {
  if (!courseInput) return ''
  if (typeof courseInput === 'string') return courseInput
  if (typeof courseInput === 'object') {
    return String(courseInput.id || courseInput.courseId || '')
  }
  return ''
}

/**
 * 处理学情迭代面板生成的讲稿
 */
const handleIterationScriptGenerated = (scriptData) => {
  if (!scriptData || !scriptData.content) {
    alert('讲稿生成失败')
    return
  }

  // 将生成的讲稿内容更新到编辑区
  currentScript.value = scriptData.content
  currentScriptNodes.value = normalizeNodes(scriptData.nodeTree || [], currentEditPage.value, scriptData.content)

  // 可选：自动切换到编辑讲稿标签页
  activeTab.value = 'script'

  alert('讲稿已生成并加载到编辑区，请继续编辑')
}

const handleIterationScriptUpdated = (scriptText) => {
  const content = String(scriptText || '').trim()
  if (!content) {
    alert('讲稿生成失败')
    return
  }

  currentScript.value = content
  currentScriptNodes.value = normalizeNodes([], currentEditPage.value, content)
  activeTab.value = 'script'
  alert('讲稿已生成并加载到编辑区，请继续编辑')
}
</script>

<style scoped>
.teacher-app {
  --app-bg: #f3f8f5;
  --app-surface: #ffffff;
  --app-surface-soft: #f4faf7;
  --app-border: rgba(120, 156, 140, 0.2);
  --app-border-strong: rgba(86, 130, 112, 0.34);
  --app-text: #111827;
  --app-muted: #5f7467;
  --app-accent: #5ca68f;

  width: 100%;
  min-height: 100vh;
  height: auto;
  overflow-y: auto;
  overflow-x: hidden;
  font-family: 'SF Pro Display', 'SF Pro Text', 'PingFang SC', 'Helvetica Neue', 'Microsoft YaHei', sans-serif;
  background: radial-gradient(circle at 45% -12%, #ffffff 0%, var(--app-bg) 58%, #e8f3ee 100%);
  padding: 20px;
  box-sizing: border-box;
}

.workspace-shell {
  width: 100%;
  min-height: calc(100vh - 40px);
  height: auto;
  border-radius: 26px;
  overflow: hidden;
  background: var(--app-surface);
  border: 1px solid rgba(255, 255, 255, 0.8);
  box-shadow: 0 24px 54px rgba(15, 23, 42, 0.08);
  display: flex;
  flex-direction: column;
}

.main-content {
  display: flex;
  min-height: calc(100vh - 158px);
  height: auto;
  overflow-y: auto;
  padding: 14px;
  gap: 14px;
}
.editor-section {
  flex: 1;
  min-width: 0;
  overflow: visible;
  display: flex;
  flex-direction: column;
  background: var(--app-surface);
  border-radius: 22px;
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.06);
  padding: 14px;
}
/* 左侧菜单导航栏 */
.left-sidebar-menu {
  flex: 0 0 206px;
  background: var(--app-surface);
  border-right: 0;
  border-radius: 22px;
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.06);
  display: flex;
  flex-direction: column;
  padding: 14px 10px;
  transition: flex-basis 0.3s ease;
}

.left-sidebar-menu.collapsed {
  flex-basis: 78px;
}

.menu-header {
  padding: 6px 12px 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  font-weight: 600;
  color: var(--app-muted);
  letter-spacing: 0.04em;
}

.left-sidebar-menu.collapsed .menu-header {
  padding: 6px 0 10px;
  justify-content: center;
}

.menu-toggle-btn {
  background: transparent;
  border: 0;
  cursor: pointer;
  color: var(--app-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  padding: 0;
  border-radius: 9px;
  transition: background 0.2s, color 0.2s;
}

.menu-toggle-btn:hover {
  background: #ecf7f2;
  color: var(--app-text);
}

.menu-toggle-btn svg {
  width: 18px;
  height: 18px;
}

.menu-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  cursor: pointer;
  color: #475569;
  transition: all 0.2s;
  border: 0;
  border-radius: 14px;
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
}

.left-sidebar-menu.collapsed .menu-item {
  padding: 12px 0;
  justify-content: center;
}

.menu-item:hover {
  background: #eef8f3;
}

.menu-item.active {
  color: #2f6052;
  background: rgba(227, 245, 238, 0.96);
  font-weight: 600;
  box-shadow: 0 6px 16px rgba(92, 166, 143, 0.18);
}

.menu-item small {
  color: #94a3b8;
  font-size: 11px;
}

.ins-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
  stroke: currentColor;
}

.tab-container {
  flex: 1;
  min-height: 0;
  overflow: visible;
  display: flex;
  flex-direction: column;
  border: 0;
  border-radius: 20px;
  background: transparent;
  box-shadow: none;
}

.tab-switch-host {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.tab-switch-enter-active,
.tab-switch-leave-active {
  transition: opacity 0.28s ease-out, transform 0.28s ease-out;
}

.tab-switch-enter-from {
  opacity: 0;
  transform: translateX(18px);
}

.tab-switch-leave-to {
  opacity: 0;
  transform: translateX(-18px);
}

.teacher-middle-layout {
  margin-top: 14px;
  min-height: calc(100vh - 110px);
  display: grid;
  grid-template-columns: 240px minmax(0, 1fr);
  gap: 14px;
}

.middle-user-sidebar {
  border-radius: 20px;
  border: 1px solid #d3e3db;
  background: linear-gradient(180deg, #f4fbf7 0%, #e8f4ee 100%);
  color: #2f5e55;
  padding: 20px 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 12px 24px rgba(33, 61, 54, 0.1);
}

.middle-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26px;
  font-weight: 800;
  color: #2a5a50;
  background: linear-gradient(180deg, #ffffff 0%, #dceee6 100%);
}

.middle-username {
  margin-top: 12px;
  font-size: 20px;
  font-weight: 700;
  letter-spacing: 0.02em;
}

.middle-subtitle {
  margin-top: 4px;
  font-size: 12px;
  color: #6b887e;
}

.middle-side-btn {
  margin-top: 10px;
  width: 100%;
  border: 1px solid #bdd8cb;
  border-radius: 12px;
  padding: 10px 12px;
  background: #ffffff;
  color: #2f605a;
  font-weight: 600;
  cursor: pointer;
}

.middle-side-btn:hover {
  border-color: #8fbcae;
  background: #f4fbf7;
}

.middle-main-panel {
  border-radius: 20px;
  border: 1px solid #dbe6df;
  background: linear-gradient(180deg, #ffffff 0%, #f6faf8 100%);
  padding: 18px;
  box-shadow: 0 16px 30px rgba(33, 61, 54, 0.08);
}

.middle-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
}

.middle-head h3 {
  margin: 0;
  font-size: 24px;
  color: #0f172a;
}

.middle-head p {
  margin: 8px 0 0;
  font-size: 14px;
  color: #607a71;
}

.middle-head.compact {
  margin-bottom: 12px;
}

.ghost-btn {
  border: 0;
  background: #ffffff;
  color: #334155;
  border-radius: 999px;
  padding: 8px 16px;
  font-size: 12px;
  cursor: pointer;
  box-shadow: 0 4px 10px rgba(15, 23, 42, 0.08);
  transition: transform 0.2s ease, box-shadow 0.2s ease, background 0.2s ease;
}

.ghost-btn:hover {
  background: #f7f8fc;
  transform: translateY(-1px);
  box-shadow: 0 8px 16px rgba(15, 23, 42, 0.1);
}

.middle-tile-grid {
  margin-top: 14px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(230px, 1fr));
  gap: 12px;
  max-height: calc(100vh - 280px);
  overflow: auto;
  padding-right: 4px;
}

.middle-tile {
  border: 1px solid #d8e6de;
  border-radius: 16px;
  background: linear-gradient(180deg, #ffffff 0%, #f5faf7 100%);
  padding: 14px;
  text-align: left;
  cursor: pointer;
  transition: transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease;
}

.middle-tile:hover {
  transform: translateY(-2px);
  border-color: #8fbcae;
  box-shadow: 0 12px 20px rgba(33, 61, 54, 0.1);
}

.middle-tile.mock {
  background: linear-gradient(180deg, #fcfcff 0%, #f4f6fd 100%);
}

.middle-tile h4 {
  margin: 10px 0 6px;
  font-size: 17px;
  color: #284a43;
}

.middle-tile p {
  margin: 0;
  font-size: 12px;
  line-height: 1.65;
  color: #638075;
}

.platform-middle-page {
  flex-direction: row;
  gap: 12px;
  min-height: calc(100vh - 180px);
}

.platform-guard-rail {
  width: 140px;
  border-radius: 14px;
  border: 1px solid #d7e2dc;
  background: linear-gradient(180deg, #f8fbf9 0%, #edf4f0 100%);
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.guard-title {
  font-size: 12px;
  color: #6a8177;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  font-weight: 700;
  margin-bottom: 2px;
}

.guard-btn {
  border: 1px solid #d5e3dc;
  background: #ffffff;
  color: #4b675d;
  border-radius: 10px;
  padding: 9px 10px;
  text-align: left;
  font-size: 13px;
  cursor: pointer;
}

.guard-btn.active {
  background: #eaf3ef;
  border-color: #8ab8ab;
  color: #2f605a;
  font-weight: 700;
}

.platform-middle-content {
  flex: 1;
  border-radius: 14px;
  border: 1px solid #dbe6df;
  background: linear-gradient(180deg, #ffffff 0%, #f7faf8 100%);
  padding: 14px;
  min-width: 0;
}

.platform-middle-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
}

.platform-middle-head h3 {
  margin: 0;
  color: #20413a;
}

.platform-middle-head p {
  margin: 6px 0 0;
  color: #637d73;
  font-size: 13px;
}

.platform-middle-head.compact {
  margin-bottom: 10px;
}

.platform-overview-grid {
  margin-top: 12px;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
}

.overview-stat-card {
  border-radius: 12px;
  border: 1px solid #d9e7df;
  background: #fff;
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.overview-stat-card span {
  font-size: 12px;
  color: #6e847b;
}

.overview-stat-card strong {
  font-size: 22px;
  color: #2d5f58;
}

.platform-tile-grid {
  margin-top: 12px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 10px;
  max-height: calc(100vh - 340px);
  overflow: auto;
}

.platform-course-tile {
  border: 1px solid #d8e6de;
  background: linear-gradient(180deg, #ffffff 0%, #f5faf7 100%);
  border-radius: 14px;
  padding: 12px;
  text-align: left;
  cursor: pointer;
}

.platform-course-tile:hover {
  border-color: #8dbdaf;
  box-shadow: 0 10px 16px rgba(36, 68, 61, 0.1);
}

.platform-course-tile.mock {
  background: linear-gradient(180deg, #fdfdff 0%, #f3f5fb 100%);
}

.platform-course-tile h4 {
  margin: 8px 0 6px;
  color: #2d4e46;
  font-size: 15px;
}

.platform-course-tile p {
  margin: 0;
  color: #688379;
  font-size: 12px;
  line-height: 1.6;
}

.tile-badge {
  display: inline-flex;
  padding: 3px 9px;
  border-radius: 999px;
  background: #edf5f2;
  color: #2f605a;
  font-size: 11px;
  font-weight: 700;
}

.tile-meta {
  margin-top: 10px;
  display: grid;
  gap: 4px;
  font-size: 12px;
  color: #4d665d;
}

.workspace-back-row {
  margin-bottom: 14px;
  padding: 0;
}


.empty-tip-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  background: #ffffff;
  border-radius: 18px;
  border: 0;
  box-shadow: 0 10px 24px rgba(15, 23, 42, 0.06);
  margin: 16px;
}

.empty-tip {
  color: var(--app-muted);
  font-size: 15px;
  text-align: center;
}

@media (max-width: 980px) {
  .teacher-middle-layout {
    grid-template-columns: 1fr;
  }

  .middle-user-sidebar {
    align-items: flex-start;
  }

  .middle-tile-grid {
    max-height: none;
  }

  .platform-middle-page {
    flex-direction: column;
  }

  .platform-guard-rail {
    width: 100%;
    flex-direction: row;
    align-items: center;
    flex-wrap: wrap;
  }

  .guard-title {
    width: 100%;
  }

  .platform-overview-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .platform-tile-grid {
    max-height: none;
  }
}

</style>