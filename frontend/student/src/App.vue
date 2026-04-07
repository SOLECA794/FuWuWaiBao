<template>
  <HomeLogin v-if="!isLoggedIn" @login-success="handleLoginSuccess" />
  <div v-else-if="!hasCourseSelected" class="course-selection-page">
    <StudentTopBar
      :backend-status-class="backendStatusClass"
      :backend-status-text="backendStatusText"
      :username="studentId"
      @logout="handleLogout"
    />
    <div class="selection-layout">
      <aside class="selection-user-sidebar">
        <div class="user-avatar">{{ (studentId || '学').slice(0, 1).toUpperCase() }}</div>
        <div class="user-name">{{ studentId || '学生' }}</div>
        <div class="user-subtitle">我的学习空间</div>

        <el-button class="refresh-btn" @click="loadCourseSelectionData" :loading="selectionLoading">刷新资源</el-button>
      </aside>

      <section class="selection-main-panel" v-loading="selectionLoading">
        <div class="selection-head">
          <div>
            <h2>选择你要学习的课件</h2>
            <p>左侧是用户栏，右侧平铺展示你的选课。点击任一卡片会直接进入学习页面。</p>
          </div>
        </div>

        <div class="selection-filters">
          <el-select v-model="selectedTeachingCourseId" placeholder="筛选课程" filterable>
            <el-option v-for="item in selectionCourseOptions" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
          <el-select v-model="selectedCourseClassId" placeholder="筛选教学班级" filterable>
            <el-option v-for="item in filteredSelectionClassOptions" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </div>

        <div class="course-tile-grid">
          <button
            v-for="card in selectionDisplayCards"
            :key="card.id"
            class="course-tile"
            :class="{ active: selectedCoursewareId === card.id, mock: card.mock }"
            @click="pickCoursewareCard(card)"
          >
            <div class="tile-badge">{{ card.mock ? '占位选课' : '我的选课' }}</div>
            <h3>{{ card.name }}</h3>
            <p>{{ card.desc }}</p>
            <div class="tile-meta">
              <span>{{ card.courseName || '未绑定课程' }}</span>
              <span>{{ card.className || '未绑定班级' }}</span>
            </div>
          </button>
        </div>

        <div class="selection-tip" v-if="selectionDisplayCards.length === 0">
          暂无可展示课件，请点击“刷新资源”重试。
        </div>
      </section>
    </div>
  </div>
  <div v-else class="student-app">
    <StudentTopBar
      :backend-status-class="backendStatusClass"
      :backend-status-text="backendStatusText"
      :username="studentId"
      @logout="handleLogout"
    />
    <div class="ambient-layer">
      <span class="orb orb-a"></span>
      <span class="orb orb-b"></span>
      <span class="orb orb-c"></span>
    </div>
    <div class="workspace-shell">
      <main class="main-layout">
        <aside class="left-sidebar-menu" :class="{ collapsed: isMenuCollapsed }">
          <div class="menu-header">
            <span v-show="!isMenuCollapsed">导航</span>
            <button class="menu-toggle-btn" @click="isMenuCollapsed = !isMenuCollapsed" :title="isMenuCollapsed ? '展开菜单' : '收起菜单'">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 12h18"></path><path d="M3 6h18"></path><path d="M3 18h18"></path></svg>
            </button>
          </div>
          <div class="menu-list">
            <button class="menu-item" :class="{ active: activeSection === 'classroom' }" @click="activeSection = 'classroom'" title="课堂学习">
              <span class="menu-icon">课</span>
              <span v-show="!isMenuCollapsed">课堂学习</span>
            </button>
            <button class="menu-item" :class="{ active: activeSection === 'analytics' }" @click="activeSection = 'analytics'" title="学习分析">
              <span class="menu-icon">析</span>
              <span v-show="!isMenuCollapsed">学习分析</span>
            </button>
            <button class="menu-item" :class="{ active: activeSection === 'recommend' }" @click="activeSection = 'recommend'" title="学习推荐">
              <span class="menu-icon">荐</span>
              <span v-show="!isMenuCollapsed">学习推荐</span>
            </button>
            <button class="menu-item" :class="{ active: activeSection === 'knowledge' }" @click="activeSection = 'knowledge'" title="知识拆解">
              <span class="menu-icon">知</span>
              <span v-show="!isMenuCollapsed">知识拆解</span>
            </button>
            <button class="menu-item" :class="{ active: showAskWorkspace }" @click="toggleAskWorkspace" title="问答浮窗">
              <span class="menu-icon">问</span>
              <span v-show="!isMenuCollapsed">问答</span>
            </button>
            <button class="menu-item" v-if="hasCourseSelected" @click="backToSelectionPage" title="返回选课页">
              <span class="menu-icon">返</span>
              <span v-show="!isMenuCollapsed">返回选课页</span>
            </button>
            <button class="menu-item" :class="{ active: activeSection === 'personal' }" @click="activeSection = 'personal'" title="个人中心">
              <span class="menu-icon">我</span>
              <span v-show="!isMenuCollapsed">个人中心</span>
            </button>
          </div>
        </aside>

        <section class="workspace-content">
          <transition name="page-fade" mode="out-in">
          <div v-if="activeSection === 'classroom'" key="classroom" class="page-layout classroom-grid">
            <section class="center-stage">
              <StudentCoursePanel
                :current-course-name="currentCourseName"
                :current-page="currentPage"
                :total-page="totalPage"
                :page-timeline-duration="pageTimelineDuration"
                :current-timeline-sec="currentTimelineSec"
                :active-node-elapsed-sec="activeNodeElapsedSec"
                :active-node-duration="activeNodeDuration"
                :current-node-title="currentNodeMeta?.title || ''"
                :active-node-type-label="activeNodeTypeLabel"
                :playback-mode="playbackMode"
                :playback-audio-meta="playbackAudioMeta"
                :progress-percent="progressPercent"
                :course-img="courseImg"
                :playback-nodes="[]"
                :current-node-id="currentNodeId"
                :tts-enabled="ttsEnabled"
                :page-summary="''"
                :script-content="currentPageMarkdown"
                :is-script-loading="scriptLoading"
                :trace-point="tracePoint"
                :trace-top="traceTop"
                :trace-left="traceLeft"
                :is-play="isPlay"
                :show-status-strip="false"
                @prev-page="prevPage"
                @select-node="selectPlaybackNode"
                @toggle-play="togglePlay"
                @toggle-tts="toggleTts"
                @speak-current-node="speakCurrentNode"
                @next-page="nextPage"
              />
            </section>

            <aside class="outline-stage">
              <div class="outline-header">
                <div>
                  <div class="outline-label">二级导航</div>
                  <h3>节点树大纲</h3>
                </div>
                <span>{{ filteredOutlineNodes.length }}/{{ playbackNodes.length }}</span>
              </div>
              <div class="outline-tools">
                <el-select v-model="outlineFilter" size="small" placeholder="筛选节点">
                  <el-option label="全部节点" value="all" />
                  <el-option label="关键讲解" value="core" />
                  <el-option label="开场节点" value="opening" />
                  <el-option label="过渡节点" value="transition" />
                </el-select>
                <el-button size="small" plain @click="focusCurrentNode">定位当前</el-button>
              </div>
              <div class="outline-list" v-if="filteredOutlineNodes.length">
                <button
                  v-for="(node, idx) in filteredOutlineNodes"
                  :key="node.node_id"
                  class="outline-item"
                  :class="{ active: node.node_id === currentNodeId }"
                  @click="selectPlaybackNode(node.node_id)"
                >
                  <span class="outline-index">{{ String(idx + 1).padStart(2, '0') }}</span>
                  <div class="outline-content">
                    <div class="outline-row">
                      <strong>{{ node.title || node.node_id }}</strong>
                      <span class="outline-time">{{ formatNodeTime(node.start_sec) }}</span>
                    </div>
                    <p>{{ node.text || '暂无节点说明' }}</p>
                  </div>
                </button>
              </div>
              <div class="outline-empty" v-else>当前页面暂无可用节点。</div>
            </aside>

            <section class="classroom-status-strip">
              <div class="status-row">
                <span class="status-pill">进度 {{ progressPercent }}%</span>
                <span class="status-pill">{{ isPlay ? '正在讲解' : '已暂停' }}</span>
                <span class="status-pill" v-if="currentNodeMeta?.title">节点 {{ currentNodeMeta.title }}</span>
                <span class="status-pill" v-if="pageTimelineDuration > 0">{{ formatNodeTime(currentTimelineSec) }} / {{ formatNodeTime(pageTimelineDuration) }}</span>
              </div>
              <div class="status-track" v-if="pageTimelineDuration > 0">
                <div class="status-fill" :style="{ width: timelinePercent + '%' }"></div>
              </div>
              <div class="status-track" v-else>
                <div class="progress-fill" :style="{ width: progressPercent + '%' }"></div>
              </div>
              <div class="status-note" v-if="courseAudioStatusText || activeNodeDuration > 0">
                <span v-if="activeNodeDuration > 0">节点 {{ formatNodeTime(activeNodeElapsedSec) }} / {{ formatNodeTime(activeNodeDuration) }}</span>
                <span>{{ activeNodeTypeLabel }}</span>
                <span v-if="courseAudioStatusText">{{ courseAudioStatusText }}</span>
              </div>
            </section>
          </div>

          <div v-else-if="activeSection === 'analytics'" key="analytics" class="page-layout single-col">
            <StudentStudyPanel
              :learning-stats="learningStats"
              :weak-point-tags="weakPointTags"
              :current-explain="currentExplain"
              :current-weak-point="currentWeakPoint"
              :current-test="currentTest"
              :test-result="testResult"
              @start-weak-point="startWeakPointLearn"
              @generate-test="generateTest"
              @check-answer="checkAnswer"
            />
          </div>

          <div v-else-if="activeSection === 'recommend'" key="recommend" class="page-layout single-col">
            <StudentRecommendPanel
              :course-name="currentCourseName"
              :current-node-title="currentNodeMeta?.title || ''"
              :current-page="currentPage"
            />
          </div>

          <div v-else-if="activeSection === 'personal'" key="personal" class="page-layout single-col">
            <StudentPersonalCenter
              :student-id="studentId"
              :course-id="courseId"
              :current-course-name="currentCourseName"
              :learning-stats="learningStats"
              :weak-point-tags="weakPointTags"
              @jump-classroom="activeSection = 'classroom'"
              @jump-analytics="activeSection = 'analytics'"
            />
          </div>

          <div v-else key="knowledge" class="page-layout single-col">
            <StudentKnowledgePanel
              :uploaded-file="uploadedFile"
              :is-parsing="isParsing"
              :parse-result="parseResult"
              :knowledge-list="knowledgeList"
              :tree-props="treeProps"
              @file-change="handleFileChange"
              @parse-knowledge="parseKnowledge"
              @node-click="handleNodeClick"
            />
          </div>
          </transition>
        </section>
      </main>
    </div>

    <transition name="qa-flyout-fade">
      <div v-if="showAskWorkspace" class="qa-flyout-backdrop" @click.self="closeAskWorkspace">
        <div class="qa-flyout-panel" :style="qaFlyoutStyle" role="dialog" aria-modal="true" aria-label="问答工作区悬浮窗">
          <div class="qa-flyout-header">
            <div class="qa-flyout-drag-handle" @pointerdown.prevent="startAskWorkspaceDrag">
              <div class="qa-flyout-kicker">问答工作区</div>
              <h3>悬浮答疑窗口</h3>
              <p>可随时收起，不影响当前课程浏览。</p>
            </div>
            <button class="qa-flyout-close" @click="closeAskWorkspace" aria-label="关闭问答悬浮窗">×</button>
          </div>

          <StudentAskPanel
            :question="question"
            :ask-loading="askLoading"
            :ai-reply="aiReply"
            :stream-typing-active="streamTypingActive"
            :qa-history="qaHistory"
            :latest-answer-meta="latestAnswerMeta"
            :summary-mode="summaryMode"
            :merged-summary="mergedSummary"
            @update:question="question = $event"
            @update:summaryMode="summaryMode = $event"
            @open-upload="openUpload"
            @generate-summary="generateMergedSummary"
            @use-summary="injectSummaryToQuestion"
            @clear-draft="clearQaDraft"
            @send-question="sendMultiModalQuestion"
          />
          <span class="qa-flyout-resize-handle" title="拖动调整大小" @pointerdown.prevent="startAskWorkspaceResize"></span>
        </div>
      </div>
    </transition>

    <footer class="footer">© 2025 智能学习课堂系统 · 学生端</footer>

    <StudentBreakpointDialog
      :model-value="showBreakpointDialog"
      :breakpoint-page="breakpointPage"
      @restart-study="restartStudy"
      @continue-study="continueStudy"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount, onUnmounted, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { studentV1Api } from './services/v1'
import { API_BASE } from './config/api'
import StudentCoursePanel from './components/student/StudentCoursePanel.vue'
import StudentAskPanel from './components/student/StudentAskPanel.vue'
import StudentStudyPanel from './components/student/StudentStudyPanel.vue'
import StudentRecommendPanel from './components/student/StudentRecommendPanel.vue'
import StudentKnowledgePanel from './components/student/StudentKnowledgePanel.vue'
import StudentBreakpointDialog from './components/student/StudentBreakpointDialog.vue'
import StudentPersonalCenter from './components/student/StudentPersonalCenter.vue'
import StudentTopBar from './components/StudentTopBar.vue'
import HomeLogin from './components/HomeLogin.vue'

const resolveStudentId = () => {
  const queryId = typeof window !== 'undefined'
    ? new URLSearchParams(window.location.search).get('studentId')
    : ''
  const normalizedQueryId = String(queryId || '').trim().toLowerCase()
  const cachedId = typeof window !== 'undefined'
    ? String(window.localStorage.getItem('fuww_student_id') || '').trim().toLowerCase()
    : ''
  const finalId = normalizedQueryId || cachedId || 'xuesheng'
  if (typeof window !== 'undefined') {
    window.localStorage.setItem('fuww_student_id', finalId)
  }
  return finalId
}

const resolveTeacherOrigin = () => {
  if (typeof window === 'undefined') return 'http://localhost:5173'
  const cached = String(window.localStorage.getItem('fuww_teacher_origin') || '').trim()
  if (cached) return cached
  const protocol = window.location.protocol || 'http:'
  const hostname = window.location.hostname || 'localhost'
  return `${protocol}//${hostname}:5173`
}

const isLoggedIn = ref(false)
const hasCourseSelected = ref(false)
const selectionLoading = ref(false)
const selectionCourseOptions = ref([])
const selectionClassOptions = ref([])
const selectionCoursewares = ref([])
const selectedTeachingCourseId = ref('')
const selectedCourseClassId = ref('')
const selectedCoursewareId = ref('')

const backendStatus = ref('checking')
let backendHealthTimer = null

const backendStatusText = computed(() => {
  if (backendStatus.value === 'online') return '在线'
  if (backendStatus.value === 'offline') return '离线'
  return '检测中'
})

const backendStatusClass = computed(() => {
  if (backendStatus.value === 'online') return 'online'
  if (backendStatus.value === 'offline') return 'offline'
  return 'checking'
})

const studentId = ref('')
const courseId = ref('')
const sessionId = ref('')
const currentNodeId = ref('p1_n1')
const playbackNodes = ref([])
const pageSummary = ref('')
const currentPageMarkdown = ref('')
const scriptLoading = ref(false)
const playbackMode = ref('duration_timeline')
const playbackAudioMeta = ref(null)
const playbackState = ref('paused')
const ttsEnabled = ref(true)
const currentCourseName = ref('')
const currentPage = ref(1)
const totalPage = ref(10)
const isPlay = ref(false)
const courseImg = ref('')
const activeSection = ref('classroom')
const isMenuCollapsed = ref(false)
const showAskWorkspace = ref(false)
const askWorkspaceLayout = reactive({
  left: 0,
  top: 0,
  width: 360,
  height: 620
})
const askWorkspaceInteraction = reactive({
  mode: '',
  pointerId: null,
  startX: 0,
  startY: 0,
  startLeft: 0,
  startTop: 0,
  startWidth: 0,
  startHeight: 0
})
const progressPercent = computed(() => Math.round((currentPage.value / totalPage.value) * 100))
const timelinePercent = computed(() => {
  if (!pageTimelineDuration.value) return 0
  return Math.min(100, Math.max(0, Math.round((currentTimelineSec.value / pageTimelineDuration.value) * 100)))
})
const filteredSelectionClassOptions = computed(() => {
  if (!selectedTeachingCourseId.value) return selectionClassOptions.value
  return selectionClassOptions.value.filter((item) => item.teachingCourseId === selectedTeachingCourseId.value)
})

const filteredSelectionCoursewares = computed(() => {
  return selectionCoursewares.value.filter((item) => {
    const courseMatch = !selectedTeachingCourseId.value || item.teachingCourseId === selectedTeachingCourseId.value
    const classMatch = !selectedCourseClassId.value || item.courseClassId === selectedCourseClassId.value
    return courseMatch && classMatch
  })
})

const selectionDisplayCards = computed(() => {
  const realCards = filteredSelectionCoursewares.value.map((item, index) => ({
    ...item,
    mock: false,
    desc: item.desc || `共 ${item.totalPage || 1} 页内容，点击卡片立即开始学习。`,
    order: index
  }))
  if (realCards.length > 0) return realCards

  const fallbackCourse = selectionCourseOptions.value[0]?.name || '示例课程'
  const fallbackClass = filteredSelectionClassOptions.value[0]?.name || '示例班级'
  return Array.from({ length: 6 }).map((_, index) => ({
    id: `mock-courseware-${index + 1}`,
    name: `占位课件 ${String(index + 1).padStart(2, '0')}`,
    courseName: fallbackCourse,
    className: fallbackClass,
    desc: '当前暂无真实选课数据，先用占位卡片展示平铺效果。',
    totalPage: 1,
    mock: true,
    order: index
  }))
})

const currentTimelineSec = ref(0)
let playbackTimer = null
let currentSpeechUtterance = null
let streamTypingTimer = null
const streamTypingQueue = ref([])
const streamTypingActive = ref(false)

const question = ref('')
const aiReply = ref('')
const askLoading = ref(false)
const qaHistory = ref([])
const latestAnswerMeta = ref({
  sourcePage: 0,
  sourceNodeId: '',
  needReteach: false,
  followUpSuggestion: '',
  sessionId: ''
})

const tracePoint = ref(false)
const traceTop = ref(0)
const traceLeft = ref(0)
const outlineFilter = ref('all')
const summaryMode = ref('quick')
const mergedSummary = ref('')

const showBreakpointDialog = ref(false)
const breakpointPage = ref(3)

const uploadedFile = ref(null)
const isParsing = ref(false)
const parseResult = ref('')
const knowledgeList = ref([])
const treeProps = ref({
  label: 'name',
  children: 'children'
})

const currentWeakPoint = ref('')
const currentExplain = ref('')
const currentTest = ref(null)
const testResult = ref(null)
const weakPointTags = ref([])
const currentQuestionId = ref('')
const learningStats = ref({
  focusScore: 0,
  totalQuestions: 0,
  weakPointCount: 0,
  masteryRate: 100
})

const pageTimelineDuration = computed(() => {
  const lastNode = playbackNodes.value[playbackNodes.value.length - 1]
  return Number(lastNode?.end_sec || 0)
})

const currentNodeMeta = computed(() => {
  return playbackNodes.value.find(node => node.node_id === currentNodeId.value) || null
})

const activeNodeDuration = computed(() => Number(currentNodeMeta.value?.duration_sec || 0))

const activeNodeElapsedSec = computed(() => {
  const node = currentNodeMeta.value
  if (!node) return 0
  return Math.max(0, Math.min(activeNodeDuration.value, currentTimelineSec.value - Number(node.start_sec || 0)))
})

const activeNodeTypeLabel = computed(() => {
  const type = currentNodeMeta.value?.type
  if (type === 'opening') return '开场讲解'
  if (type === 'transition') return '过渡收束'
  return '核心讲解'
})

const courseAudioStatusText = computed(() => {
  const status = playbackAudioMeta.value?.audio_status
  const duration = Number(playbackAudioMeta.value?.audio_duration_sec || 0)
  if (!status) return ''
  if (status === 'ready' && duration > 0) {
    return `音频已生成 ${formatNodeTime(duration)}`
  }
  if (status === 'processing') return '音频生成中'
  return '使用时长驱动讲解'
})

const filteredOutlineNodes = computed(() => {
  if (outlineFilter.value === 'all') return playbackNodes.value
  if (outlineFilter.value === 'core') {
    return playbackNodes.value.filter(node => (node.type || 'core') === 'core')
  }
  return playbackNodes.value.filter(node => node.type === outlineFilter.value)
})

const normalizeTimeSec = (value, fallback = 0) => {
  const numeric = Number(value)
  if (!Number.isFinite(numeric)) return fallback
  return Math.max(0, Math.floor(numeric))
}

const formatNodeTime = (seconds) => {
  const normalized = normalizeTimeSec(seconds)
  const mins = String(Math.floor(normalized / 60)).padStart(2, '0')
  const secs = String(normalized % 60).padStart(2, '0')
  return `${mins}:${secs}`
}

const trimText = (text, length = 56) => {
  const value = String(text || '').replace(/\s+/g, ' ').trim()
  if (!value) return ''
  if (value.length <= length) return value
  return `${value.slice(0, length)}...`
}

const focusCurrentNode = () => {
  if (!currentNodeId.value) return
  const exists = filteredOutlineNodes.value.some(node => node.node_id === currentNodeId.value)
  if (!exists) {
    outlineFilter.value = 'all'
  }
}

const generateMergedSummary = () => {
  const points = playbackNodes.value
    .slice(0, 3)
    .map((node) => `${node.title || node.node_id}：${trimText(node.text, 44)}`)
  const baseSummary = pageSummary.value || points.join('；') || '本页暂未解析出可用摘要。'

  if (summaryMode.value === 'exam') {
    mergedSummary.value = `【考试速记】${baseSummary}\n重点节点：${points.join(' | ') || '无'}\n建议：优先理解定义、过程和结论。`
    return
  }
  if (summaryMode.value === 'teach') {
    mergedSummary.value = `【讲解版】本页核心：${baseSummary}\n可先讲“${currentNodeMeta.value?.title || '当前节点'}”，再用例子解释。`
    return
  }
  mergedSummary.value = `【速览】${baseSummary}`
}

const injectSummaryToQuestion = () => {
  if (!mergedSummary.value) {
    ElMessage.info('请先生成摘要，再用于提问')
    return
  }
  question.value = `请基于以下摘要，帮我用更通俗的方式讲解：\n${mergedSummary.value}`
  ElMessage.success('摘要已写入提问框')
}

const closeAskWorkspace = () => {
  showAskWorkspace.value = false
}

const toggleAskWorkspace = () => {
  ensureAskWorkspaceLayout()
  showAskWorkspace.value = !showAskWorkspace.value
}

const ASK_WORKSPACE_LAYOUT_KEY = 'fuww_student_ask_workspace_layout'
const ASK_WORKSPACE_MARGIN = 12
const ASK_WORKSPACE_TOP = 68
const ASK_WORKSPACE_MIN_WIDTH = 320
const ASK_WORKSPACE_MIN_HEIGHT = 420

const clamp = (value, min, max) => Math.min(Math.max(value, min), max)

const getViewportBounds = () => {
  if (typeof window === 'undefined') {
    return { width: 1280, height: 720 }
  }
  return {
    width: window.innerWidth || 1280,
    height: window.innerHeight || 720
  }
}

const getDefaultAskWorkspaceLayout = () => {
  const viewport = getViewportBounds()
  const width = clamp(360, ASK_WORKSPACE_MIN_WIDTH, Math.max(ASK_WORKSPACE_MIN_WIDTH, viewport.width - ASK_WORKSPACE_MARGIN * 2))
  const height = clamp(620, ASK_WORKSPACE_MIN_HEIGHT, Math.max(ASK_WORKSPACE_MIN_HEIGHT, viewport.height - ASK_WORKSPACE_TOP - ASK_WORKSPACE_MARGIN))
  return {
    left: Math.max(ASK_WORKSPACE_MARGIN, viewport.width - width - ASK_WORKSPACE_MARGIN),
    top: ASK_WORKSPACE_TOP,
    width,
    height
  }
}

const clampAskWorkspaceLayout = (layout) => {
  const viewport = getViewportBounds()
  const widthLimit = Math.max(ASK_WORKSPACE_MIN_WIDTH, viewport.width - ASK_WORKSPACE_MARGIN * 2)
  const heightLimit = Math.max(ASK_WORKSPACE_MIN_HEIGHT, viewport.height - ASK_WORKSPACE_TOP - ASK_WORKSPACE_MARGIN)
  const width = clamp(Math.round(layout.width || 0), ASK_WORKSPACE_MIN_WIDTH, widthLimit)
  const height = clamp(Math.round(layout.height || 0), ASK_WORKSPACE_MIN_HEIGHT, heightLimit)
  const leftLimit = Math.max(ASK_WORKSPACE_MARGIN, viewport.width - width - ASK_WORKSPACE_MARGIN)
  const topLimit = Math.max(ASK_WORKSPACE_TOP, viewport.height - height - ASK_WORKSPACE_MARGIN)
  const left = clamp(Math.round(layout.left || 0), ASK_WORKSPACE_MARGIN, leftLimit)
  const top = clamp(Math.round(layout.top || 0), ASK_WORKSPACE_TOP, topLimit)
  return { left, top, width, height }
}

const ensureAskWorkspaceLayout = () => {
  const clamped = clampAskWorkspaceLayout(askWorkspaceLayout)
  askWorkspaceLayout.left = clamped.left
  askWorkspaceLayout.top = clamped.top
  askWorkspaceLayout.width = clamped.width
  askWorkspaceLayout.height = clamped.height
}

const persistAskWorkspaceLayout = () => {
  if (typeof window === 'undefined') return
  window.localStorage.setItem(ASK_WORKSPACE_LAYOUT_KEY, JSON.stringify({
    left: askWorkspaceLayout.left,
    top: askWorkspaceLayout.top,
    width: askWorkspaceLayout.width,
    height: askWorkspaceLayout.height
  }))
}

const loadAskWorkspaceLayout = () => {
  if (typeof window === 'undefined') return
  let parsed = null
  try {
    parsed = JSON.parse(window.localStorage.getItem(ASK_WORKSPACE_LAYOUT_KEY) || 'null')
  } catch (error) {
    parsed = null
  }
  const merged = parsed && typeof parsed === 'object'
    ? {
        left: Number(parsed.left),
        top: Number(parsed.top),
        width: Number(parsed.width),
        height: Number(parsed.height)
      }
    : getDefaultAskWorkspaceLayout()
  const clamped = clampAskWorkspaceLayout(merged)
  askWorkspaceLayout.left = clamped.left
  askWorkspaceLayout.top = clamped.top
  askWorkspaceLayout.width = clamped.width
  askWorkspaceLayout.height = clamped.height
}

const qaFlyoutStyle = computed(() => ({
  left: `${askWorkspaceLayout.left}px`,
  top: `${askWorkspaceLayout.top}px`,
  width: `${askWorkspaceLayout.width}px`,
  height: `${askWorkspaceLayout.height}px`
}))

const stopAskWorkspaceInteraction = () => {
  if (typeof window === 'undefined') return
  window.removeEventListener('pointermove', handleAskWorkspacePointerMove)
  window.removeEventListener('pointerup', handleAskWorkspacePointerUp)
  window.removeEventListener('pointercancel', handleAskWorkspacePointerUp)
  window.removeEventListener('blur', handleAskWorkspacePointerUp)
  document.body.style.userSelect = ''
  document.body.style.cursor = ''
  askWorkspaceInteraction.mode = ''
  askWorkspaceInteraction.pointerId = null
}

const handleAskWorkspacePointerMove = (event) => {
  if (!askWorkspaceInteraction.mode || typeof window === 'undefined') return
  const viewport = getViewportBounds()
  if (askWorkspaceInteraction.mode === 'drag') {
    const nextLayout = clampAskWorkspaceLayout({
      left: askWorkspaceInteraction.startLeft + (event.clientX - askWorkspaceInteraction.startX),
      top: askWorkspaceInteraction.startTop + (event.clientY - askWorkspaceInteraction.startY),
      width: askWorkspaceLayout.width,
      height: askWorkspaceLayout.height
    })
    askWorkspaceLayout.left = nextLayout.left
    askWorkspaceLayout.top = nextLayout.top
    return
  }

  const nextWidth = clamp(
    askWorkspaceInteraction.startWidth + (event.clientX - askWorkspaceInteraction.startX),
    ASK_WORKSPACE_MIN_WIDTH,
    Math.max(ASK_WORKSPACE_MIN_WIDTH, viewport.width - askWorkspaceLayout.left - ASK_WORKSPACE_MARGIN)
  )
  const nextHeight = clamp(
    askWorkspaceInteraction.startHeight + (event.clientY - askWorkspaceInteraction.startY),
    ASK_WORKSPACE_MIN_HEIGHT,
    Math.max(ASK_WORKSPACE_MIN_HEIGHT, viewport.height - askWorkspaceLayout.top - ASK_WORKSPACE_MARGIN)
  )
  askWorkspaceLayout.width = nextWidth
  askWorkspaceLayout.height = nextHeight
}

const handleAskWorkspacePointerUp = () => {
  stopAskWorkspaceInteraction()
  ensureAskWorkspaceLayout()
  persistAskWorkspaceLayout()
}

const startAskWorkspaceDrag = (event) => {
  if (!showAskWorkspace.value) return
  if (event.button !== 0) return
  ensureAskWorkspaceLayout()
  askWorkspaceInteraction.mode = 'drag'
  askWorkspaceInteraction.pointerId = event.pointerId
  askWorkspaceInteraction.startX = event.clientX
  askWorkspaceInteraction.startY = event.clientY
  askWorkspaceInteraction.startLeft = askWorkspaceLayout.left
  askWorkspaceInteraction.startTop = askWorkspaceLayout.top
  askWorkspaceInteraction.startWidth = askWorkspaceLayout.width
  askWorkspaceInteraction.startHeight = askWorkspaceLayout.height
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'move'
  window.addEventListener('pointermove', handleAskWorkspacePointerMove)
  window.addEventListener('pointerup', handleAskWorkspacePointerUp)
  window.addEventListener('pointercancel', handleAskWorkspacePointerUp)
  window.addEventListener('blur', handleAskWorkspacePointerUp)
}

const startAskWorkspaceResize = (event) => {
  if (!showAskWorkspace.value) return
  if (event.button !== 0) return
  ensureAskWorkspaceLayout()
  askWorkspaceInteraction.mode = 'resize'
  askWorkspaceInteraction.pointerId = event.pointerId
  askWorkspaceInteraction.startX = event.clientX
  askWorkspaceInteraction.startY = event.clientY
  askWorkspaceInteraction.startLeft = askWorkspaceLayout.left
  askWorkspaceInteraction.startTop = askWorkspaceLayout.top
  askWorkspaceInteraction.startWidth = askWorkspaceLayout.width
  askWorkspaceInteraction.startHeight = askWorkspaceLayout.height
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'nwse-resize'
  window.addEventListener('pointermove', handleAskWorkspacePointerMove)
  window.addEventListener('pointerup', handleAskWorkspacePointerUp)
  window.addEventListener('pointercancel', handleAskWorkspacePointerUp)
  window.addEventListener('blur', handleAskWorkspacePointerUp)
}

const clearQaDraft = () => {
  question.value = ''
  aiReply.value = ''
  stopStreamTypewriter()
}

const clampTimelineSec = (value) => {
  const normalized = normalizeTimeSec(value)
  if (pageTimelineDuration.value <= 0) return normalized
  return Math.min(pageTimelineDuration.value, normalized)
}

const normalizeTimelineForNode = (nodeId) => {
  const node = playbackNodes.value.find(item => item.node_id === nodeId)
  currentTimelineSec.value = node ? Number(node.start_sec || 0) : 0
}

const applyPlaybackPosition = ({ nodeId = '', timeSec = null } = {}) => {
  if (!playbackNodes.value.length) {
    currentNodeId.value = `p${currentPage.value}_n1`
    currentTimelineSec.value = 0
    return
  }

  const matchedNode = playbackNodes.value.find(item => item.node_id === nodeId)
  const fallbackNode = matchedNode || playbackNodes.value[0]
  currentNodeId.value = fallbackNode?.node_id || `p${currentPage.value}_n1`

  if (timeSec !== null && timeSec !== undefined) {
    currentTimelineSec.value = clampTimelineSec(timeSec)
    syncCurrentNodeWithTimeline()
    return
  }

  normalizeTimelineForNode(currentNodeId.value)
}

const syncCurrentNodeWithTimeline = () => {
  if (!playbackNodes.value.length) return
  const matched = playbackNodes.value.find((node) => {
    const start = Number(node.start_sec || 0)
    const end = Number(node.end_sec || 0)
    return currentTimelineSec.value >= start && currentTimelineSec.value < end
  }) || playbackNodes.value[playbackNodes.value.length - 1]
  if (matched?.node_id && matched.node_id !== currentNodeId.value) {
    currentNodeId.value = matched.node_id
  }
}

const stopPlaybackTimer = () => {
  if (playbackTimer) {
    window.clearInterval(playbackTimer)
    playbackTimer = null
  }
}

const stopStreamTypewriter = () => {
  if (streamTypingTimer) {
    window.clearInterval(streamTypingTimer)
    streamTypingTimer = null
  }
  streamTypingQueue.value = []
  streamTypingActive.value = false
}

const startStreamTypewriter = () => {
  if (streamTypingTimer || streamTypingQueue.value.length === 0) return
  streamTypingActive.value = true
  streamTypingTimer = window.setInterval(() => {
    if (!streamTypingQueue.value.length) {
      window.clearInterval(streamTypingTimer)
      streamTypingTimer = null
      streamTypingActive.value = false
      return
    }
    const nextChar = streamTypingQueue.value.shift()
    aiReply.value += nextChar
  }, 16)
}

const pushTypewriterText = (text) => {
  const value = String(text || '')
  if (!value) return
  streamTypingQueue.value.push(...value.split(''))
  startStreamTypewriter()
}

const waitTypewriterDrain = async () => {
  const startedAt = Date.now()
  while (streamTypingQueue.value.length > 0 || streamTypingActive.value) {
    if (Date.now() - startedAt > 3000) {
      aiReply.value += streamTypingQueue.value.join('')
      stopStreamTypewriter()
      break
    }
    await new Promise(resolve => window.setTimeout(resolve, 30))
  }
}

const stopSpeechNarration = () => {
  if (window.speechSynthesis) {
    window.speechSynthesis.cancel()
  }
  currentSpeechUtterance = null
}

const speakCurrentNode = () => {
  if (!ttsEnabled.value || !isPlay.value) return
  if (!window.speechSynthesis || typeof window.SpeechSynthesisUtterance === 'undefined') return
  const node = currentNodeMeta.value
  if (!node) return
  const text = String(node.text || node.title || '').trim()
  if (!text) return

  const speakingMark = `${currentPage.value}:${node.node_id}:${normalizeTimeSec(node.start_sec)}`
  if (currentSpeechUtterance?.__mark === speakingMark) {
    return
  }

  stopSpeechNarration()
  const utter = new SpeechSynthesisUtterance(text)
  utter.lang = 'zh-CN'
  utter.rate = 1
  utter.pitch = 1
  utter.volume = 1
  utter.__mark = speakingMark
  utter.onend = () => {
    if (currentSpeechUtterance === utter) {
      currentSpeechUtterance = null
    }
  }
  utter.onerror = () => {
    if (currentSpeechUtterance === utter) {
      currentSpeechUtterance = null
    }
  }
  currentSpeechUtterance = utter
  window.speechSynthesis.speak(utter)
}

const startPlaybackTimer = () => {
  stopPlaybackTimer()
  if (!playbackNodes.value.length) return
  playbackTimer = window.setInterval(async () => {
    if (!isPlay.value) return
    currentTimelineSec.value += 1
    syncCurrentNodeWithTimeline()

    if (pageTimelineDuration.value > 0 && currentTimelineSec.value >= pageTimelineDuration.value) {
      if (currentPage.value < totalPage.value) {
        currentPage.value += 1
        await refreshCurrentPageData({ preserveCurrentNode: false })
        await saveBreakpoint()
      } else {
        isPlay.value = false
        stopPlaybackTimer()
      }
    }
  }, 1000)
}

const prevPage = async () => {
  if (!courseId.value || currentPage.value <= 1) return
  isPlay.value = false
  currentPage.value--
  await refreshCurrentPageData({ preserveCurrentNode: false })
  await saveBreakpoint()
}

const nextPage = async () => {
  if (!courseId.value || currentPage.value >= totalPage.value) return
  isPlay.value = false
  currentPage.value++
  await refreshCurrentPageData({ preserveCurrentNode: false })
  await saveBreakpoint()
}

const updateCourseContent = () => {
  if (!courseId.value) {
    courseImg.value = ''
    return
  }
  courseImg.value = `${API_BASE}/api/courseware/${courseId.value}/page/${currentPage.value}`
}

const loadStudentScript = async () => {
  scriptLoading.value = true
  if (!courseId.value) {
    currentNodeId.value = 'p1_n1'
    playbackNodes.value = []
    pageSummary.value = ''
    currentPageMarkdown.value = ''
    playbackMode.value = 'duration_timeline'
    playbackAudioMeta.value = null
    currentTimelineSec.value = 0
    scriptLoading.value = false
    return
  }
  try {
    const data = await studentV1Api.coursewares.getPlaybackScript(courseId.value, currentPage.value)
    const payload = data?.data || {}
    const nodes = data?.data?.nodes || []
    playbackNodes.value = nodes
    pageSummary.value = payload.page_summary || ''
    mergedSummary.value = payload.page_summary || ''
    currentPageMarkdown.value = String(
      payload.script || payload.content || payload.markdown || payload.raw_script || ''
    )
    playbackAudioMeta.value = payload.audio_meta || null
    playbackMode.value = payload.playback_mode || payload.audio_meta?.playback_mode || 'duration_timeline'
    applyPlaybackPosition({ nodeId: currentNodeId.value })
  } catch (error) {
    playbackNodes.value = []
    pageSummary.value = ''
    mergedSummary.value = ''
    currentPageMarkdown.value = ''
    playbackMode.value = 'duration_timeline'
    playbackAudioMeta.value = null
    currentNodeId.value = `p${currentPage.value}_n1`
    currentTimelineSec.value = 0
  } finally {
    scriptLoading.value = false
  }
}

const selectPlaybackNode = async (nodeId) => {
  currentNodeId.value = nodeId
  normalizeTimelineForNode(nodeId)
  await saveBreakpoint()
}

const refreshCurrentPageData = async ({ preserveCurrentNode = true, targetNodeId = '', targetTimeSec = null } = {}) => {
  const nextNodeId = preserveCurrentNode ? currentNodeId.value : (targetNodeId || `p${currentPage.value}_n1`)
  currentNodeId.value = nextNodeId
  updateCourseContent()
  await loadStudentScript()
  applyPlaybackPosition({ nodeId: targetNodeId || nextNodeId, timeSec: targetTimeSec })
}

const togglePlay = () => {
  if (!playbackNodes.value.length) {
    ElMessage.warning('当前页暂无可播放的讲授节点')
    return
  }
  isPlay.value = !isPlay.value
  playbackState.value = isPlay.value ? 'lecturing' : 'paused'
  if (!isPlay.value) {
    stopSpeechNarration()
    return
  }
  speakCurrentNode()
}

const toggleTts = () => {
  ttsEnabled.value = !ttsEnabled.value
  if (!ttsEnabled.value) {
    stopSpeechNarration()
    ElMessage.info('已关闭语音讲稿')
    return
  }
  if (isPlay.value) {
    speakCurrentNode()
  }
  ElMessage.success('已开启语音讲稿')
}

const openUpload = () => {
  ElMessage.info('已打开截图/圈图提问')
}

const sendMultiModalQuestion = async () => {
  if (askLoading.value) {
    ElMessage.info('当前正在处理上一条提问，请稍后')
    return
  }
  if (!question.value.trim()) {
    ElMessage.warning('请输入问题后再发送')
    return
  }
  if (!courseId.value) {
    ElMessage.warning('暂无可用课件，无法提问')
    return
  }

  askLoading.value = true
  isPlay.value = false
  playbackState.value = 'tutoring'
  stopSpeechNarration()
  stopStreamTypewriter()
  try {
    const currentQuestion = question.value
    aiReply.value = ''
    latestAnswerMeta.value = {
      sourcePage: 0,
      sourceNodeId: '',
      needReteach: false,
      followUpSuggestion: '',
      sessionId: sessionId.value
    }
    let finalPayload = null
    await studentV1Api.qa.stream({
      sessionId: sessionId.value,
      courseId: courseId.value,
      page: currentPage.value,
      nodeId: currentNodeId.value,
      question: currentQuestion
    }, {
      token: (payload) => {
        pushTypewriterText(payload.text || '')
      },
      sentence: (payload) => {
        if (!aiReply.value && streamTypingQueue.value.length === 0) {
          pushTypewriterText(payload.text || '')
        }
      },
      final: (payload) => {
        finalPayload = payload
      }
    })

    await waitTypewriterDrain()

    if (finalPayload?.session_id) {
      sessionId.value = finalPayload.session_id
    }

    const resumePage = Number(finalPayload?.resume_page || currentPage.value) || currentPage.value
    const resumeNodeId = finalPayload?.resume_node_id || currentNodeId.value
    const resumeSec = finalPayload?.resume_sec
    if (resumePage !== currentPage.value) {
      currentPage.value = resumePage
      await refreshCurrentPageData({
        preserveCurrentNode: false,
        targetNodeId: resumeNodeId,
        targetTimeSec: resumeSec
      })
    } else if (resumeNodeId || resumeSec !== undefined) {
      applyPlaybackPosition({ nodeId: resumeNodeId, timeSec: resumeSec })
    }

    latestAnswerMeta.value = {
      sourcePage: finalPayload?.source_page || resumePage,
      sourceNodeId: finalPayload?.source_node_id || finalPayload?.sourceNodeId || (resumeNodeId || currentNodeId.value),
      needReteach: !!finalPayload?.need_reteach,
      followUpSuggestion: finalPayload?.follow_up_suggestion || '',
      sessionId: finalPayload?.session_id || sessionId.value
    }

    qaHistory.value.unshift({
      question: currentQuestion,
      answer: aiReply.value,
      sourcePage: latestAnswerMeta.value.sourcePage,
      sourceNodeId: latestAnswerMeta.value.sourceNodeId
    })
    if (qaHistory.value.length > 5) {
      qaHistory.value = qaHistory.value.slice(0, 5)
    }
    question.value = ''
    if (finalPayload?.need_reteach) {
      ElMessage.success('已按追问语境切换为重讲模式')
    } else {
      playbackState.value = 'resuming'
      isPlay.value = true
      playbackState.value = 'lecturing'
      speakCurrentNode()
      ElMessage.success('AI 答疑完成，并已准备继续讲解')
    }
  } catch (error) {
    try {
      const fallbackResp = await studentV1Api.qa.ask({
        courseId: courseId.value,
        studentId: studentId.value,
        pageNum: currentPage.value,
        nodeId: currentNodeId.value,
        question: question.value
      })
      const payload = fallbackResp?.data || {}
      aiReply.value = payload.answer || ''
      latestAnswerMeta.value = {
        sourcePage: payload.sourcePage || payload.source_page || currentPage.value,
        sourceNodeId: payload.sourceNodeId || payload.source_node_id || currentNodeId.value,
        needReteach: !!payload.needReteach,
        followUpSuggestion: payload.followUpSuggestion || '',
        sessionId: sessionId.value
      }
      qaHistory.value.unshift({
        question: question.value,
        answer: aiReply.value,
        sourcePage: latestAnswerMeta.value.sourcePage,
        sourceNodeId: latestAnswerMeta.value.sourceNodeId
      })
      if (qaHistory.value.length > 5) {
        qaHistory.value = qaHistory.value.slice(0, 5)
      }
      question.value = ''
      playbackState.value = latestAnswerMeta.value.needReteach ? 'tutoring' : 'resuming'
      if (!latestAnswerMeta.value.needReteach) {
        isPlay.value = true
        playbackState.value = 'lecturing'
        speakCurrentNode()
      }
      ElMessage.warning('流式问答失败，已切换到普通问答返回结果')
    } catch (fallbackError) {
      aiReply.value = ''
      ElMessage.error(`提问失败：${fallbackError.message || error.message}`)
    }
  } finally {
    if (!latestAnswerMeta.value.needReteach) {
      playbackState.value = isPlay.value ? 'lecturing' : 'paused'
      if (isPlay.value) {
        speakCurrentNode()
      }
    }
    askLoading.value = false
  }
}

const pickCoursewareCard = async (card) => {
  if (!card) return
  if (card.mock) {
    const fallback = filteredSelectionCoursewares.value[0] || selectionCoursewares.value[0]
    if (fallback) {
      selectedCoursewareId.value = fallback.id
      if (fallback.teachingCourseId) selectedTeachingCourseId.value = fallback.teachingCourseId
      if (fallback.courseClassId) selectedCourseClassId.value = fallback.courseClassId
    } else {
      // 无真实课件时仍允许进入后续页面，后续由课程初始化阶段给出提示。
      selectedCoursewareId.value = ''
    }
    await enterWorkspaceFromSelection({ allowPlaceholder: true })
    return
  }
  if (card.teachingCourseId) {
    selectedTeachingCourseId.value = card.teachingCourseId
  }
  if (card.courseClassId) {
    selectedCourseClassId.value = card.courseClassId
  }
  selectedCoursewareId.value = card.id
  await enterWorkspaceFromSelection()
}

const loadCourseSelectionData = async () => {
  selectionLoading.value = true
  try {
    const [courseRes, classRes, coursewareRes] = await Promise.all([
      studentV1Api.platform.listCourses({ page: 1, pageSize: 100 }),
      studentV1Api.platform.listClasses({ page: 1, pageSize: 200 }),
      studentV1Api.coursewares.list()
    ])

    const platformCourses = Array.isArray(courseRes?.data?.items) ? courseRes.data.items : []
    const platformClasses = Array.isArray(classRes?.data?.items) ? classRes.data.items : []
    const coursewareList = Array.isArray(coursewareRes?.data) ? coursewareRes.data : []

    selectionCourseOptions.value = platformCourses.map((item) => ({
      id: String(item.courseId || ''),
      name: item.title || '未命名课程'
    }))

    selectionClassOptions.value = platformClasses.map((item) => ({
      id: String(item.classId || ''),
      name: item.className || '未命名班级',
      teachingCourseId: String(item.teachingCourseId || '')
    }))

    selectionCoursewares.value = coursewareList.map((item) => ({
      id: String(item.id || item.courseId || ''),
      name: item.title || '未命名课件',
      totalPage: Number(item.total_page || 1),
      teachingCourseId: String(item.teaching_course_id || ''),
      courseName: String(item.teaching_course_title || ''),
      courseClassId: String(item.course_class_id || ''),
      className: String(item.course_class_name || ''),
      desc: item.is_published ? '已发布，可进入学习' : '未发布，暂为教师侧预览资源',
      published: Boolean(item.is_published)
    }))

    if (!selectedTeachingCourseId.value && selectionCourseOptions.value.length > 0) {
      selectedTeachingCourseId.value = selectionCourseOptions.value[0].id
    }
    if (!selectedCourseClassId.value && filteredSelectionClassOptions.value.length > 0) {
      selectedCourseClassId.value = filteredSelectionClassOptions.value[0].id
    }
    if (!selectedCoursewareId.value && filteredSelectionCoursewares.value.length > 0) {
      selectedCoursewareId.value = filteredSelectionCoursewares.value[0].id
    }
  } catch (error) {
    selectionCourseOptions.value = []
    selectionClassOptions.value = []
    selectionCoursewares.value = []
    ElMessage.error(`课件选择数据加载失败：${error.message}`)
  } finally {
    selectionLoading.value = false
  }
}

const startStudentWorkspace = async () => {
  checkBackendHealth()
  if (backendHealthTimer) {
    window.clearInterval(backendHealthTimer)
  }
  backendHealthTimer = window.setInterval(checkBackendHealth, 30000)
  startPlaybackTimer()
  await initializeCourseContext()
}

const enterWorkspaceFromSelection = async ({ allowPlaceholder = false } = {}) => {
  const selected = selectionCoursewares.value.find((item) => item.id === selectedCoursewareId.value)
  if (!selected && !allowPlaceholder) {
    ElMessage.warning('请先选择要学习的课件')
    return
  }

  if (selected) {
    courseId.value = selected.id
    currentCourseName.value = selected.name
    totalPage.value = selected.totalPage || 1
  } else {
    courseId.value = ''
    currentCourseName.value = '临时占位学习空间'
    totalPage.value = 1
  }
  currentPage.value = 1
  hasCourseSelected.value = true
  await startStudentWorkspace()
}

const backToSelectionPage = () => {
  hasCourseSelected.value = false
  isPlay.value = false
  playbackState.value = 'paused'
  showAskWorkspace.value = false
  stopSpeechNarration()
  void loadCourseSelectionData()
}

const handleLoginSuccess = (user) => {
  const role = String(user?.role || '').trim().toLowerCase()
  const username = String(user?.username || '').trim().toLowerCase() || 'xuesheng'

  if (role === 'teacher') {
    const teacherOrigin = resolveTeacherOrigin()
    const encodedUsername = encodeURIComponent(username)
    window.location.href = `${teacherOrigin}/?role=teacher&username=${encodedUsername}`
    return
  }

  isLoggedIn.value = true
  studentId.value = username
  hasCourseSelected.value = false
  if (typeof window !== 'undefined') {
    window.localStorage.setItem('fuww_student_id', username)
    window.localStorage.setItem('fuww_student_origin', window.location.origin)
  }
  void loadCourseSelectionData()
}

const handleLogout = () => {
  isLoggedIn.value = false
  hasCourseSelected.value = false
  showAskWorkspace.value = false
  studentId.value = ''
  selectedTeachingCourseId.value = ''
  selectedCourseClassId.value = ''
  selectedCoursewareId.value = ''
  selectionCourseOptions.value = []
  selectionClassOptions.value = []
  selectionCoursewares.value = []
  question.value = ''
  aiReply.value = ''
  qaHistory.value = []
  isPlay.value = false
  stopPlaybackTimer()
  stopSpeechNarration()
  stopStreamTypewriter()
  if (backendHealthTimer) {
    window.clearInterval(backendHealthTimer)
    backendHealthTimer = null
  }
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    loadAskWorkspaceLayout()
    window.localStorage.setItem('fuww_student_origin', window.location.origin)
    const params = new URLSearchParams(window.location.search)
    const role = String(params.get('role') || '').trim().toLowerCase()
    const username = String(params.get('username') || '').trim().toLowerCase()

    if (role && username) {
      if (role === 'teacher') {
        const teacherOrigin = resolveTeacherOrigin()
        window.location.replace(`${teacherOrigin}/?role=teacher&username=${encodeURIComponent(username)}`)
        return
      }
      isLoggedIn.value = true
      hasCourseSelected.value = false
      studentId.value = username
      window.localStorage.setItem('fuww_student_id', username)
      window.history.replaceState({}, document.title, window.location.pathname)
      void loadCourseSelectionData()
      return
    }

    studentId.value = resolveStudentId()
  }
})

onUnmounted(() => {
  if (backendHealthTimer) {
    window.clearInterval(backendHealthTimer)
    backendHealthTimer = null
  }
  stopPlaybackTimer()
  stopSpeechNarration()
  stopStreamTypewriter()
  stopAskWorkspaceInteraction()
})

onBeforeUnmount(() => {
  stopAskWorkspaceInteraction()
})

watch(selectedTeachingCourseId, () => {
  const classValid = filteredSelectionClassOptions.value.some((item) => item.id === selectedCourseClassId.value)
  if (!classValid) {
    selectedCourseClassId.value = filteredSelectionClassOptions.value[0]?.id || ''
  }
  const coursewareValid = filteredSelectionCoursewares.value.some((item) => item.id === selectedCoursewareId.value)
  if (!coursewareValid) {
    selectedCoursewareId.value = filteredSelectionCoursewares.value[0]?.id || ''
  }
})

watch(selectedCourseClassId, () => {
  const coursewareValid = filteredSelectionCoursewares.value.some((item) => item.id === selectedCoursewareId.value)
  if (!coursewareValid) {
    selectedCoursewareId.value = filteredSelectionCoursewares.value[0]?.id || ''
  }
})

watch(isPlay, (value) => {
  if (value) {
    playbackState.value = askLoading.value ? 'tutoring' : 'lecturing'
    startPlaybackTimer()
    speakCurrentNode()
    return
  }
  playbackState.value = askLoading.value ? 'tutoring' : 'paused'
  stopPlaybackTimer()
  stopSpeechNarration()
})

watch(playbackNodes, () => {
  if (isPlay.value) {
    startPlaybackTimer()
    speakCurrentNode()
  }
})

watch(currentNodeId, () => {
  if (isPlay.value) {
    speakCurrentNode()
  }
})

const initializeCourseContext = async () => {
  try {
    if (!courseId.value) {
      const data = await studentV1Api.coursewares.list({
        teachingCourseId: selectedTeachingCourseId.value,
        courseClassId: selectedCourseClassId.value
      })
      const list = Array.isArray(data?.data) ? data.data : []
      const published = list.filter(item => item.is_published)
      const target = published[0] || list[0]

      if (!target) {
        courseId.value = ''
        currentCourseName.value = ''
        totalPage.value = 1
        updateCourseContent()
        ElMessage.warning('当前课程/班级暂无可学习课件，请联系教师发布课件')
        return
      }

      courseId.value = String(target.id || target.courseId || '')
      currentCourseName.value = target.title || '未命名课件'
      totalPage.value = target.total_page || 1
    }

    currentPage.value = 1
    await refreshCurrentPageData({ preserveCurrentNode: false })

    const session = await studentV1Api.sessions.start({
      userId: studentId.value,
      courseId: courseId.value
    })
    sessionId.value = session?.data?.sessionId || ''

    await loadBreakpoint()
    await loadWeakPoints()
    await loadStudyData()
  } catch (error) {
    courseId.value = ''
    currentCourseName.value = ''
    updateCourseContent()
    ElMessage.error(`加载课程失败：${error.message}`)
  }
}

const checkBackendHealth = async () => {
  try {
    const res = await studentV1Api.health()
    backendStatus.value = res.ok ? 'online' : 'offline'
  } catch (error) {
    backendStatus.value = 'offline'
  }
}

const loadBreakpoint = async () => {
  if (!courseId.value) return
  try {
    const data = await studentV1Api.coursewares.getBreakpoint(studentId.value, courseId.value)
    breakpointPage.value = data?.data?.pageNum || data?.data?.lastPageNum || 1
    showBreakpointDialog.value = breakpointPage.value > 1
  } catch (error) {
    breakpointPage.value = 1
    showBreakpointDialog.value = false
  }
}

const saveBreakpoint = async () => {
  if (!courseId.value) return
  try {
    await studentV1Api.coursewares.updateBreakpoint({
      studentId: studentId.value,
      courseId: courseId.value,
      pageNum: currentPage.value
    })
    await studentV1Api.sessions.updateProgress({
      sessionId: sessionId.value,
      userId: studentId.value,
      courseId: courseId.value,
      currentPage: currentPage.value,
      currentNodeId: currentNodeId.value,
      currentTimeSec: currentTimelineSec.value
    })
  } catch (error) {
    console.warn('断点保存失败', error)
  }
}

const loadStudyData = async () => {
  if (!courseId.value) return
  try {
    const data = await studentV1Api.coursewares.getStats(studentId.value, courseId.value)
    const payload = data.data || {}
    const weakPoints = payload.weakPoints || []
    learningStats.value = {
      focusScore: payload.focusScore || 0,
      totalQuestions: payload.totalQuestions || 0,
      weakPointCount: weakPoints.length,
      masteryRate: Math.max(35, 100 - (weakPoints.length * 10))
    }
  } catch (error) {
    console.warn('学习数据加载失败', error)
  }
}

const loadWeakPoints = async () => {
  if (!courseId.value) return
  try {
    const data = await studentV1Api.coursewares.getWeakPoints(studentId.value, courseId.value)
    if (Array.isArray(data.data) && data.data.length > 0) {
      weakPointTags.value = data.data.map(item => ({ id: item.weakPointId, name: item.name }))
    } else {
      weakPointTags.value = []
    }
  } catch (error) {
    weakPointTags.value = []
    console.warn('加载薄弱点失败', error)
  }
}

const continueStudy = async () => {
  currentPage.value = breakpointPage.value
  await refreshCurrentPageData({ preserveCurrentNode: false })
  showBreakpointDialog.value = false
  await saveBreakpoint()
  ElMessage.success(`已为你跳转到第 ${breakpointPage.value} 页`)
}

const restartStudy = async () => {
  currentPage.value = 1
  await refreshCurrentPageData({ preserveCurrentNode: false })
  showBreakpointDialog.value = false
  await saveBreakpoint()
  ElMessage.info('已回到第1页重新开始学习')
}

const handleFileChange = (file) => {
  uploadedFile.value = file.raw
  parseResult.value = ''
  knowledgeList.value = []
}

const parseKnowledge = async () => {
  if (!uploadedFile.value) {
    ElMessage.warning('请先上传 PDF/PPTX 文件！')
    return
  }

  isParsing.value = true
  try {
    const fileText = await uploadedFile.value.text().catch(() => '')
    const textPayload = fileText.trim() || `文件名: ${uploadedFile.value.name}`

    const data = await studentV1Api.knowledge.parse({
      fileContent: textPayload,
      fileType: uploadedFile.value.name.split('.').pop() || 'unknown',
      studentId: studentId.value
    })

    const structure = data?.data?.structure || []
    knowledgeList.value = structure.map((item, index) => ({
      id: index + 1,
      name: item.name,
      children: item.children || []
    }))

    parseResult.value = `拆解成功！共识别出 ${countNodes(knowledgeList.value)} 个知识点`
    ElMessage.success('知识点结构拆解完成！')
  } catch (error) {
    parseResult.value = '拆解失败，请重试！'
    ElMessage.error('知识点拆解失败')
  } finally {
    isParsing.value = false
  }
}

const countNodes = (tree) => {
  let count = 0
  tree.forEach((node) => {
    count++
    if (node.children && node.children.length) {
      count += countNodes(node.children)
    }
  })
  return count
}

const handleNodeClick = (data) => {
  ElMessage.info(`已定位到知识点：${data.name}`)
  tracePoint.value = true
  traceTop.value = 200
  traceLeft.value = 300
}

const startWeakPointLearn = async (point) => {
  currentWeakPoint.value = point.name
  try {
    const data = await studentV1Api.weakPoints.explain(point.id, point.name)
    currentExplain.value = data?.data?.content || '暂无讲解内容'
  } catch (error) {
    currentExplain.value = '暂时无法获取讲解，请稍后重试。'
    ElMessage.error(`讲解加载失败：${error.message}`)
  }
  currentTest.value = null
  testResult.value = null
}

const generateTest = async () => {
  try {
    const currentPoint = weakPointTags.value.find(item => item.name === currentWeakPoint.value)
    const data = await studentV1Api.weakPoints.generateTest({
      weakPointId: currentPoint?.id || '',
      weakPointName: currentWeakPoint.value,
      studentId: studentId.value,
      questionType: 'single'
    })
    currentQuestionId.value = data?.data?.questionId || ''
    currentTest.value = {
      question: data?.data?.content || '暂无题目',
      options: data?.data?.options || []
    }
    testResult.value = null
  } catch (error) {
    ElMessage.error(`生成习题失败：${error.message}`)
  }
}

const checkAnswer = async (option) => {
  try {
    const data = await studentV1Api.weakPoints.checkAnswer({
      studentId: studentId.value,
      questionId: currentQuestionId.value,
      userAnswer: option
    })
    testResult.value = {
      correct: data?.data?.isCorrect,
      msg: data?.data?.isCorrect ? '✅ 回答正确！' : '❌ 回答错误',
      analysis: data?.data?.explanation || ''
    }
  } catch (error) {
    ElMessage.error(`答案校验失败：${error.message}`)
  }
}
</script>

<style scoped>
.course-selection-page {
  min-height: 100vh;
  padding: 14px;
  box-sizing: border-box;
  background: radial-gradient(circle at 12% 8%, #f5fbf8 0%, #edf3ef 45%, #e8efeb 100%);
}

.course-selection-page :deep(.top-nav) {
  margin-bottom: 14px;
}

.selection-layout {
  margin-top: 14px;
  min-height: calc(100vh - 102px);
  display: grid;
  grid-template-columns: 200px minmax(0, 1fr);
  gap: 14px;
}

.selection-user-sidebar {
  border-radius: 18px;
  border: 1px solid #cfe4da;
  background: linear-gradient(180deg, #f3faf6 0%, #e7f4ed 100%);
  color: #2f605a;
  padding: 18px 14px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.user-avatar {
  width: 58px;
  height: 58px;
  border-radius: 50%;
  background: linear-gradient(180deg, #ffffff 0%, #dceee6 100%);
  box-shadow: 0 10px 18px rgba(33, 61, 54, 0.16);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 700;
}

.user-name {
  margin-top: 12px;
  font-size: 20px;
  font-weight: 700;
  letter-spacing: 0.02em;
}

.user-subtitle {
  margin-top: 2px;
  font-size: 12px;
  color: #6d877d;
}

.refresh-btn {
  margin-top: auto;
  width: 100%;
  border: 1px solid #bdd8cb;
  background: #ffffff;
}

.refresh-btn :deep(span) {
  color: #2f605a;
  font-weight: 700;
}

.selection-main-panel {
  border-radius: 18px;
  border: 1px solid #d9e7df;
  background: linear-gradient(180deg, #ffffff 0%, #f7faf8 100%);
  box-shadow: 0 20px 42px rgba(33, 61, 54, 0.08);
  padding: 18px;
  overflow: hidden;
}

.selection-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 14px;
}

.selection-head h2 {
  margin: 0;
  font-size: 24px;
  color: #1f3f38;
}

.selection-head p {
  margin: 8px 0 0;
  font-size: 14px;
  color: #5f7b71;
}

.selection-filters {
  margin-top: 14px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 240px));
  gap: 10px;
}

.course-tile-grid {
  margin-top: 14px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(230px, 1fr));
  gap: 12px;
  max-height: calc(100vh - 290px);
  overflow: auto;
  padding-right: 4px;
}

.course-tile {
  border: 1px solid #d9e7df;
  border-radius: 16px;
  background: linear-gradient(180deg, #ffffff 0%, #f6fbf8 100%);
  padding: 14px;
  text-align: left;
  cursor: pointer;
  transition: transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease;
}

.course-tile:hover {
  transform: translateY(-2px);
  border-color: #90c0b5;
  box-shadow: 0 12px 20px rgba(33, 61, 54, 0.1);
}

.course-tile.active {
  border-color: #5d8f83;
  box-shadow: 0 0 0 2px rgba(93, 143, 131, 0.18);
}

.course-tile.mock {
  background: linear-gradient(180deg, #fcfcff 0%, #f4f6fd 100%);
}

.tile-badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 9px;
  border-radius: 999px;
  background: #edf5f2;
  color: #2f605a;
  font-size: 11px;
  font-weight: 700;
}

.course-tile h3 {
  margin: 10px 0 6px;
  font-size: 16px;
  color: #24453f;
}

.course-tile p {
  margin: 0;
  font-size: 12px;
  line-height: 1.65;
  color: #648177;
}

.tile-meta {
  margin-top: 10px;
  display: grid;
  gap: 4px;
  font-size: 12px;
  color: #4d665d;
}

.selection-tip {
  margin-top: 12px;
  font-size: 13px;
  color: #6f867d;
}

@media (max-width: 980px) {
  .selection-layout {
    grid-template-columns: 1fr;
  }

  .selection-user-sidebar {
    align-items: flex-start;
  }

  .user-avatar {
    width: 46px;
    height: 46px;
    font-size: 20px;
  }

  .selection-filters {
    grid-template-columns: minmax(0, 1fr);
  }

  .course-tile-grid {
    max-height: none;
  }
}

.student-app {
  position: relative;
  width: 100%;
  min-height: 100vh;
  padding: 14px;
  box-sizing: border-box;
  background: radial-gradient(circle at 12% 8%, #f5fbf8 0%, #edf3ef 45%, #e8efeb 100%);
  font-family: 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  overflow: hidden;
}

.student-app :deep(.top-nav) {
  position: relative;
  z-index: 2;
  margin-bottom: 14px;
}

.ambient-layer {
  position: absolute;
  inset: 0;
  pointer-events: none;
  z-index: 0;
}

.orb {
  position: absolute;
  border-radius: 999px;
  filter: blur(38px);
  opacity: 0.45;
  animation: floatOrb 14s ease-in-out infinite;
}

.orb-a {
  width: 220px;
  height: 220px;
  background: #9ccfc3;
  left: 4%;
  top: 8%;
}

.orb-b {
  width: 280px;
  height: 280px;
  background: #bddfd6;
  right: 6%;
  top: 18%;
  animation-delay: -4s;
}

.orb-c {
  width: 180px;
  height: 180px;
  background: #d4e8e1;
  right: 20%;
  bottom: 10%;
  animation-delay: -8s;
}

.workspace-shell {
  position: relative;
  z-index: 1;
  width: 100%;
  min-height: calc(100vh - 84px);
  border-radius: 28px;
  background: #f7faf8;
  border: 1px solid #d8e4dc;
  box-shadow: 0 24px 48px rgba(45, 72, 66, 0.08);
  overflow: hidden;
}

.main-layout {
  min-height: calc(100vh - 250px);
  padding: 12px 18px 18px;
  display: flex;
  gap: 14px;
}

.left-sidebar-menu {
  flex: 0 0 180px;
  background: #ffffff;
  border: 1px solid #d9e7df;
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.left-sidebar-menu.collapsed {
  flex-basis: 70px;
}

.menu-header {
  padding: 14px 14px 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  letter-spacing: 0.08em;
  color: #6b7f75;
  text-transform: uppercase;
  font-weight: 700;
}

.left-sidebar-menu.collapsed .menu-header {
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
}

.menu-toggle-btn svg {
  width: 18px;
  height: 18px;
}

.menu-list {
  display: flex;
  flex-direction: column;
  padding: 0 8px 8px;
  gap: 4px;
}

.menu-item {
  border: 1px solid #d4e4db;
  background: #fff;
  color: #536a61;
  font-size: 13px;
  font-weight: 600;
  border-radius: 10px;
  padding: 8px 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 10px;
  text-align: left;
  transform: translateY(0);
}

.left-sidebar-menu.collapsed .menu-item {
  justify-content: center;
}

.menu-icon {
  width: 22px;
  height: 22px;
  border-radius: 8px;
  background: #e8f2ed;
  color: #2f605a;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
}

.menu-item:hover {
  border-color: #7ea497;
  color: #2f605a;
  transform: translateY(-1px);
}

.menu-item.active {
  color: #fff;
  background: linear-gradient(135deg, #2f605a 0%, #4d8a80 100%);
  border-color: #2f605a;
  box-shadow: 0 8px 16px rgba(47, 96, 90, 0.25);
}

.menu-item.active .menu-icon {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
}

.workspace-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.logout-btn {
  border: 1px solid #cfe0d7;
  background: #ffffff;
  color: #2f605a;
  border-radius: 999px;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
}

.logout-btn:hover {
  border-color: #2f605a;
}

.page-layout {
  height: 100%;
}

.page-fade-enter-active,
.page-fade-leave-active {
  transition: all 0.26s ease;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
  transform: translateY(8px) scale(0.995);
}

.page-layout.two-col {
  display: flex;
  gap: 14px;
}

.page-layout.single-col {
  display: block;
}

.left-stage {
  flex: 1 1 62%;
  min-width: 0;
}

.right-stage {
  flex: 1 1 38%;
  min-width: 360px;
  max-width: 560px;
}

.page-layout.classroom-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.outline-stage {
  border: 1px solid #d8e5de;
  border-radius: 18px;
  background: linear-gradient(180deg, #ffffff 0%, #f6faf8 100%);
  padding: 12px;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.outline-header {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  align-items: flex-start;
}

.outline-label {
  font-size: 11px;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: #6a8278;
  font-weight: 700;
}

.outline-header h3 {
  margin-top: 3px;
  font-size: 17px;
  color: #23463f;
}

.outline-header span {
  font-size: 12px;
  border-radius: 999px;
  border: 1px solid #d1e2da;
  padding: 3px 8px;
  color: #42665d;
  background: #eef5f1;
}

.outline-tools {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.outline-list {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  overflow: auto;
  padding-right: 2px;
}

.outline-item {
  border: 1px solid #d7e5dd;
  background: #fff;
  border-radius: 12px;
  padding: 8px;
  text-align: left;
  display: flex;
  gap: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.outline-item:hover {
  border-color: #8caea2;
  transform: translateY(-1px);
}

.outline-item.active {
  border-color: #2f605a;
  box-shadow: 0 10px 16px rgba(47, 96, 90, 0.16);
  background: linear-gradient(180deg, #eff6f3 0%, #ffffff 100%);
}

.outline-index {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: #e4efe9;
  color: #2f605a;
  font-size: 11px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.outline-content {
  min-width: 0;
}

.outline-row {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  align-items: center;
}

.outline-row strong {
  font-size: 13px;
  color: #274d46;
}

.outline-time {
  color: #6d847b;
  font-size: 11px;
}

.outline-content p {
  margin-top: 5px;
  font-size: 12px;
  line-height: 1.45;
  color: #577068;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.outline-empty {
  margin-top: 12px;
  border: 1px dashed #cfddd5;
  border-radius: 12px;
  padding: 12px;
  color: #70857c;
  font-size: 13px;
}

.classroom-status-strip {
  border: 1px solid rgba(148, 163, 184, 0.16);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.95) 0%, rgba(244, 250, 247, 0.96) 100%);
  border-radius: 12px;
  padding: 10px 12px;
}

.classroom-status-strip .status-row {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 7px;
}

.classroom-status-strip .status-pill {
  display: inline-flex;
  align-items: center;
  padding: 3px 9px;
  border-radius: 999px;
  font-size: 12px;
  color: #3b5d54;
  background: #edf4f0;
  border: 1px solid #d5e4dc;
}

.classroom-status-strip .status-track,
.classroom-status-strip .progress-track {
  height: 8px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.7);
  overflow: hidden;
}

.classroom-status-strip .status-fill,
.classroom-status-strip .progress-fill {
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #0f766e 0%, #0284c7 100%);
}

.classroom-status-strip .status-note {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 7px;
  font-size: 12px;
  color: #628075;
}

.center-stage {
  min-width: 0;
}

.qa-flyout-backdrop {
  position: fixed;
  inset: 56px 0 0 0;
  z-index: 40;
  background: rgba(255, 255, 255, 0.02);
  backdrop-filter: none;
}

.qa-flyout-panel {
  position: fixed;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  border-radius: 22px;
  background: rgba(247, 250, 248, 0.9);
  border: 1px solid rgba(216, 229, 222, 0.88);
  box-shadow: 0 18px 34px rgba(45, 72, 66, 0.12);
  overflow: hidden;
}

.qa-flyout-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  cursor: move;
  user-select: none;
}

.qa-flyout-drag-handle {
  min-width: 0;
}

.qa-flyout-kicker {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #6a8278;
}

.qa-flyout-header h3 {
  margin-top: 4px;
  font-size: 18px;
  color: #23463f;
}

.qa-flyout-header p {
  margin-top: 6px;
  font-size: 13px;
  color: #6f867d;
}

.qa-flyout-close {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  border: 1px solid #d0dfd7;
  background: #fff;
  color: #50695f;
  font-size: 20px;
  line-height: 1;
  cursor: pointer;
}

.qa-flyout-resize-handle {
  position: absolute;
  right: 8px;
  bottom: 8px;
  width: 18px;
  height: 18px;
  border-right: 2px solid rgba(47, 96, 90, 0.42);
  border-bottom: 2px solid rgba(47, 96, 90, 0.42);
  border-radius: 0 0 14px 0;
  cursor: nwse-resize;
  opacity: 0.85;
}

.qa-flyout-panel :deep(.panel-box) {
  height: 100%;
  min-height: 0;
  background:
    radial-gradient(circle at top right, rgba(143, 193, 181, 0.14), transparent 34%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.88) 0%, rgba(246, 251, 248, 0.84) 100%);
}

.qa-flyout-panel :deep(.conversation-board) {
  flex: 1;
  min-height: 0;
}

.qa-flyout-panel :deep(.conversation-thread) {
  max-height: none;
}

.qa-flyout-panel:active .qa-flyout-resize-handle,
.qa-flyout-panel:hover .qa-flyout-resize-handle {
  opacity: 1;
}

@keyframes qa-flyout-pop {
  from {
    opacity: 0;
    transform: translateY(16px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.qa-flyout-fade-enter-active,
.qa-flyout-fade-leave-active {
  transition: opacity 0.22s ease;
}

.qa-flyout-fade-enter-from,
.qa-flyout-fade-leave-to {
  opacity: 0;
}

.qa-flyout-fade-enter-active .qa-flyout-panel {
  animation: qa-flyout-pop 0.24s ease both;
}

.qa-flyout-fade-leave-active .qa-flyout-panel {
  animation: qa-flyout-pop 0.18s ease reverse both;
}

.footer {
  height: 42px;
  line-height: 42px;
  text-align: center;
  color: #70847a;
  font-size: 12px;
  position: relative;
  z-index: 1;
}

@keyframes floatOrb {
  0%, 100% { transform: translateY(0) translateX(0); }
  50% { transform: translateY(-10px) translateX(8px); }
}

@media (max-width: 1280px) {
  .main-layout {
    gap: 10px;
  }

  .page-layout.two-col {
    flex-direction: column;
  }

  .right-stage {
    min-width: 0;
    max-width: 100%;
  }

  .page-layout.classroom-grid {
    grid-template-columns: minmax(0, 1fr);
  }

  .outline-stage {
    max-height: 280px;
  }

  .qa-flyout-panel {
    width: min(360px, calc(100vw - 24px));
    height: min(620px, calc(100vh - 80px));
  }
}

@media (max-width: 768px) {
  .student-app {
    padding: 8px;
  }

  .workspace-shell {
    border-radius: 18px;
  }

  .main-layout {
    padding: 10px;
    flex-direction: column;
  }

  .left-sidebar-menu {
    flex: 0 0 auto;
    width: 100%;
  }

  .left-sidebar-menu.collapsed {
    flex-basis: auto;
  }

  .menu-list {
    flex-direction: row;
    flex-wrap: wrap;
  }

  .menu-item {
    flex: 1 1 calc(50% - 4px);
  }

  .qa-flyout-backdrop {
    inset: 56px 0 0 0;
  }

  .qa-flyout-panel {
    width: min(calc(100vw - 16px), 420px);
    height: min(calc(100vh - 72px), 620px);
    border-radius: 18px;
  }

  .qa-flyout-header {
    cursor: default;
  }
}
</style>