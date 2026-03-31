<template>
  <div class="student-app">
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
            <button class="menu-item" :class="{ active: activeSection === 'trace' }" @click="activeSection = 'trace'" title="溯源定位">
              <span class="menu-icon">溯</span>
              <span v-show="!isMenuCollapsed">溯源定位</span>
            </button>
            <button class="menu-item" :class="{ active: activeSection === 'knowledge' }" @click="activeSection = 'knowledge'" title="知识拆解">
              <span class="menu-icon">知</span>
              <span v-show="!isMenuCollapsed">知识拆解</span>
            </button>
          </div>
        </aside>

        <section class="workspace-content">
          <section class="section-bar">
            <div class="section-title-wrap">
              <span class="section-tag">{{ activeSectionMeta.tag }}</span>
              <h2>{{ activeSectionMeta.title }}</h2>
            </div>
            <p>{{ activeSectionMeta.desc }}</p>
          </section>

          <transition name="page-fade" mode="out-in">
          <div v-if="activeSection === 'classroom'" key="classroom" class="page-layout classroom-grid">
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
                @prev-page="prevPage"
                @select-node="selectPlaybackNode"
                @toggle-play="togglePlay"
                @toggle-tts="toggleTts"
                @speak-current-node="speakCurrentNode"
                @next-page="nextPage"
              />
            </section>
            <aside class="ai-stage">
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
            </aside>
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

          <div v-else-if="activeSection === 'trace'" key="trace" class="page-layout two-col">
            <section class="left-stage">
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
                :playback-nodes="playbackNodes"
                :current-node-id="currentNodeId"
                :tts-enabled="ttsEnabled"
                :page-summary="pageSummary"
                :script-content="currentPageMarkdown"
                :is-script-loading="scriptLoading"
                :trace-point="tracePoint"
                :trace-top="traceTop"
                :trace-left="traceLeft"
                :is-play="isPlay"
                @prev-page="prevPage"
                @select-node="selectPlaybackNode"
                @toggle-play="togglePlay"
                @toggle-tts="toggleTts"
                @speak-current-node="speakCurrentNode"
                @next-page="nextPage"
              />
            </section>
            <section class="right-stage">
              <StudentTracePanel
                :trace-log="traceLog"
                @open-trace-mode="openTraceMode"
              />
            </section>
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
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { studentV1Api } from './services/v1'
import { API_BASE } from './config/api'
import StudentCoursePanel from './components/student/StudentCoursePanel.vue'
import StudentAskPanel from './components/student/StudentAskPanel.vue'
import StudentStudyPanel from './components/student/StudentStudyPanel.vue'
import StudentTracePanel from './components/student/StudentTracePanel.vue'
import StudentKnowledgePanel from './components/student/StudentKnowledgePanel.vue'
import StudentBreakpointDialog from './components/student/StudentBreakpointDialog.vue'

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

const backendStatus = ref('checking')
let backendHealthTimer = null

const studentId = ref(resolveStudentId())
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
const sectionMetas = {
  classroom: { tag: 'Learning', title: '课堂学习', desc: '左看结构，中间看课件，右侧直接提问，学习流程一屏完成。' },
  analytics: { tag: 'Analytics', title: '学习分析', desc: '聚焦薄弱点与练习反馈，快速找到下一步学习重点。' },
  trace: { tag: 'Trace', title: '溯源定位', desc: '把问题定位到具体页面与节点，避免泛泛追问。' },
  knowledge: { tag: 'Knowledge', title: '知识拆解', desc: '把资料拆成知识树，便于回看和定位关键概念。' }
}
const activeSectionMeta = computed(() => sectionMetas[activeSection.value] || sectionMetas.classroom)
const isMenuCollapsed = ref(false)
const progressPercent = computed(() => Math.round((currentPage.value / totalPage.value) * 100))
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
const traceLog = ref('')
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
    traceLog.value = `问答定位：第 ${latestAnswerMeta.value.sourcePage || currentPage.value} 页 / 节点 ${latestAnswerMeta.value.sourceNodeId || currentNodeId.value}，续接节点 ${resumeNodeId || currentNodeId.value}`
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
      traceLog.value = `问答定位：第 ${latestAnswerMeta.value.sourcePage || currentPage.value} 页 / 节点 ${latestAnswerMeta.value.sourceNodeId || currentNodeId.value}`
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

const openTraceMode = () => {
  tracePoint.value = true
  traceTop.value = 150
  traceLeft.value = 200
  const sourceNode = latestAnswerMeta.value?.sourceNodeId || currentNodeId.value
  traceLog.value = `已定位：第 ${currentPage.value} 页 → 当前节点 ${currentNodeId.value}${sourceNode ? `（最近问答来源节点 ${sourceNode}）` : ''}`
}

onMounted(() => {
  checkBackendHealth()
  backendHealthTimer = window.setInterval(checkBackendHealth, 30000)
  startPlaybackTimer()
  initializeCourseContext()
})

onUnmounted(() => {
  if (backendHealthTimer) {
    window.clearInterval(backendHealthTimer)
    backendHealthTimer = null
  }
  stopPlaybackTimer()
  stopSpeechNarration()
  stopStreamTypewriter()
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
    const data = await studentV1Api.coursewares.list()
    const list = Array.isArray(data?.data) ? data.data : []
    const published = list.filter(item => item.is_published)
    const target = published[0] || list[0]

    if (!target) {
      courseId.value = ''
      currentCourseName.value = ''
      totalPage.value = 1
      updateCourseContent()
      ElMessage.warning('暂无可学习课件，请联系教师发布课件')
      return
    }

    courseId.value = target.id
    currentCourseName.value = target.title || '未命名课件'
    totalPage.value = target.total_page || 1
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
  traceLog.value = `已定位知识点：${data.name}`
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
  min-height: calc(100vh - 56px);
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

.section-bar {
  border: 1px solid #d7e5dd;
  background: linear-gradient(135deg, #ffffff 0%, #f4f8f6 100%);
  border-radius: 14px;
  padding: 10px 12px;
}

.section-title-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
}

.section-tag {
  font-size: 11px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  border-radius: 999px;
  padding: 3px 9px;
  color: #2f605a;
  background: #e9f3ee;
  border: 1px solid #d2e4db;
  font-weight: 700;
}

.section-bar h2 {
  font-size: 17px;
  color: #24453f;
}

.section-bar p {
  margin-top: 4px;
  color: #6b8178;
  font-size: 12px;
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
  display: grid;
  grid-template-columns: 270px minmax(0, 1fr) 420px;
  gap: 12px;
  align-items: stretch;
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

.center-stage {
  min-width: 0;
}

.ai-stage {
  min-width: 0;
}

.ai-stage :deep(.panel-box) {
  height: 100%;
  min-height: 520px;
  display: flex;
  flex-direction: column;
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

  .ai-stage :deep(.panel-box) {
    height: auto;
    min-height: 0;
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
}
</style>