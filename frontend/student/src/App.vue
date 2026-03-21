<template>
  <div class="app-container">
    <StudentHeaderOverview
      :student-id="studentId"
      :current-course-name="currentCourseName"
      :backend-status="backendStatus"
      :backend-status-label="backendStatusLabel"
      :progress-percent="progressPercent"
      :current-page="currentPage"
      :is-play="isPlay"
    />

    <main class="main">
      <section class="left-section">
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

      <section class="right-section">
        <el-tabs v-model="activeTab" class="smart-tab">
          <el-tab-pane label="多模态提问" name="ask">
            <StudentAskPanel
              :question="question"
              :ask-loading="askLoading"
              :ai-reply="aiReply"
              :stream-typing-active="streamTypingActive"
              :qa-history="qaHistory"
              :latest-answer-meta="latestAnswerMeta"
              @update:question="question = $event"
              @open-upload="openUpload"
              @send-question="sendMultiModalQuestion"
            />
          </el-tab-pane>

          <el-tab-pane label="学习数据" name="data">
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
          </el-tab-pane>

          <el-tab-pane label="溯源定位" name="trace">
            <StudentTracePanel
              :trace-log="traceLog"
              @open-trace-mode="openTraceMode"
            />
          </el-tab-pane>

          <el-tab-pane label="知识点拆解" name="parse">
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
          </el-tab-pane>
        </el-tabs>
      </section>
    </main>

    <footer class="footer">© 2025 智能学习课堂系统 - 版权所有</footer>

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
import StudentHeaderOverview from './components/student/StudentHeaderOverview.vue'
import StudentCoursePanel from './components/student/StudentCoursePanel.vue'
import StudentAskPanel from './components/student/StudentAskPanel.vue'
import StudentStudyPanel from './components/student/StudentStudyPanel.vue'
import StudentTracePanel from './components/student/StudentTracePanel.vue'
import StudentKnowledgePanel from './components/student/StudentKnowledgePanel.vue'
import StudentBreakpointDialog from './components/student/StudentBreakpointDialog.vue'

const backendStatus = ref('checking')
const backendStatusLabel = computed(() => {
  if (backendStatus.value === 'online') return '在线'
  if (backendStatus.value === 'offline') return '离线'
  return '检测中'
})
let backendHealthTimer = null

const studentId = ref('student001')
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
const activeTab = ref('ask')
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

const normalizeTimeSec = (value, fallback = 0) => {
  const numeric = Number(value)
  if (!Number.isFinite(numeric)) return fallback
  return Math.max(0, Math.floor(numeric))
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
    currentPageMarkdown.value = String(
      payload.script || payload.content || payload.markdown || payload.raw_script || ''
    )
    playbackAudioMeta.value = payload.audio_meta || null
    playbackMode.value = payload.playback_mode || payload.audio_meta?.playback_mode || 'duration_timeline'
    applyPlaybackPosition({ nodeId: currentNodeId.value })
  } catch (error) {
    playbackNodes.value = []
    pageSummary.value = ''
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
    showBreakpointDialog.value = true
  } catch (error) {
    breakpointPage.value = 1
    showBreakpointDialog.value = true
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
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
.app-container {
  width: 100%;
  height: 100vh;
  background:
    radial-gradient(circle at top left, rgba(250, 204, 21, 0.18), transparent 24%),
    radial-gradient(circle at top right, rgba(14, 165, 233, 0.16), transparent 22%),
    linear-gradient(180deg, #fffdf7 0%, #f4f8fc 45%, #edf3f9 100%);
  display: flex;
  flex-direction: column;
  font-family: 'HarmonyOS Sans SC', 'Microsoft YaHei', sans-serif;
}
.main {
  flex: 1;
  display: flex;
  gap: 16px;
  padding: 16px 20px 20px;
  overflow: hidden;
}
.left-section {
  flex: 6;
  min-width: 0;
}
.right-section {
  flex: 4;
  overflow: hidden;
  min-width: 360px;
}
.smart-tab {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.7);
  border: 1px solid rgba(148, 163, 184, 0.18);
  border-radius: 24px;
  backdrop-filter: blur(10px);
  box-shadow: 0 18px 42px rgba(15, 23, 42, 0.08);
  padding: 0 14px 14px;
}
:deep(.el-tabs__header) {
  margin: 0;
  padding-top: 6px;
}
:deep(.el-tabs__nav-wrap::after) {
  display: none;
}
:deep(.el-tabs__item) {
  color: #475569;
  font-weight: 600;
}
:deep(.el-tabs__item.is-active) {
  color: #0f172a;
}
:deep(.el-tabs__active-bar) {
  background: linear-gradient(90deg, #0f766e 0%, #0284c7 100%);
  height: 3px;
  border-radius: 999px;
}
:deep(.el-tabs__content) {
  flex: 1;
  overflow: auto;
}
.footer {
  height: 40px;
  background: rgba(255, 255, 255, 0.62);
  text-align: center;
  line-height: 40px;
  font-size: 12px;
  color: #64748b;
  border-top: 1px solid rgba(226, 232, 240, 0.8);
  backdrop-filter: blur(8px);
}

@media (max-width: 1100px) {
	.main {
		flex-direction: column;
		overflow: auto;
	}
	.right-section {
		min-width: 0;
	}
}

@media (max-width: 720px) {
	.main {
		padding: 12px;
	}
}
</style>