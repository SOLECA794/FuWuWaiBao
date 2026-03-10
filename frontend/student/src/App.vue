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
          :course-img="courseImg"
          :playback-nodes="playbackNodes"
          :current-node-id="currentNodeId"
          :page-summary="pageSummary"
          :trace-point="tracePoint"
          :trace-top="traceTop"
          :trace-left="traceLeft"
          :is-play="isPlay"
          @prev-page="prevPage"
          @select-node="selectPlaybackNode"
          @toggle-play="togglePlay"
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
              :qa-history="qaHistory"
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
import { ref, onMounted, onUnmounted, computed } from 'vue'
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
const currentCourseName = ref('')
const currentPage = ref(1)
const totalPage = ref(10)
const isPlay = ref(false)
const courseImg = ref('')
const activeTab = ref('ask')
const progressPercent = computed(() => Math.round((currentPage.value / totalPage.value) * 100))

const question = ref('')
const aiReply = ref('')
const askLoading = ref(false)
const qaHistory = ref([])

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

const prevPage = async () => {
  if (!courseId.value || currentPage.value <= 1) return
  currentPage.value--
  await refreshCurrentPageData()
  await saveBreakpoint()
}

const nextPage = async () => {
  if (!courseId.value || currentPage.value >= totalPage.value) return
  currentPage.value++
  await refreshCurrentPageData()
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
  if (!courseId.value) {
    currentNodeId.value = 'p1_n1'
    playbackNodes.value = []
    pageSummary.value = ''
    return
  }
  try {
    const data = await studentV1Api.coursewares.getPlaybackScript(courseId.value, currentPage.value)
    const nodes = data?.data?.nodes || []
    playbackNodes.value = nodes
    pageSummary.value = data?.data?.page_summary || ''
    currentNodeId.value = nodes[0]?.node_id || `p${currentPage.value}_n1`
  } catch (error) {
    playbackNodes.value = []
    pageSummary.value = ''
    currentNodeId.value = `p${currentPage.value}_n1`
  }
}

const selectPlaybackNode = async (nodeId) => {
  currentNodeId.value = nodeId
  await saveBreakpoint()
}

const refreshCurrentPageData = async () => {
  updateCourseContent()
  await loadStudentScript()
}

const togglePlay = () => {
  isPlay.value = !isPlay.value
}

const openUpload = () => {
  ElMessage.info('已打开截图/圈图提问')
}

const sendMultiModalQuestion = async () => {
  if (!question.value.trim()) {
    ElMessage.warning('请输入问题后再发送')
    return
  }
  if (!courseId.value) {
    ElMessage.warning('暂无可用课件，无法提问')
    return
  }

  askLoading.value = true
  try {
    const currentQuestion = question.value
    aiReply.value = ''
    let finalPayload = null
    await studentV1Api.qa.stream({
      sessionId: sessionId.value,
      courseId: courseId.value,
      page: currentPage.value,
      nodeId: currentNodeId.value,
      question: currentQuestion
    }, {
      token: (payload) => {
        aiReply.value += payload.text || ''
      },
      sentence: (payload) => {
        if (!aiReply.value) aiReply.value = payload.text || ''
      },
      final: (payload) => {
        finalPayload = payload
      }
    })

    qaHistory.value.unshift({
      question: currentQuestion,
      answer: aiReply.value
    })
    if (qaHistory.value.length > 3) {
      qaHistory.value = qaHistory.value.slice(0, 3)
    }
    if (finalPayload?.resume_node_id) {
      currentNodeId.value = finalPayload.resume_node_id
    }
    question.value = ''
    ElMessage.success('AI 答疑完成')
  } catch (error) {
    aiReply.value = ''
    ElMessage.error(`提问失败：${error.message}`)
  } finally {
    askLoading.value = false
  }
}

const openTraceMode = () => {
  tracePoint.value = true
  traceTop.value = 150
  traceLeft.value = 200
  traceLog.value = `已定位：第 ${currentPage.value} 页 → 节点 ${currentNodeId.value}`
}

onMounted(() => {
  checkBackendHealth()
  backendHealthTimer = window.setInterval(checkBackendHealth, 30000)
  initializeCourseContext()
})

onUnmounted(() => {
  if (backendHealthTimer) {
    window.clearInterval(backendHealthTimer)
    backendHealthTimer = null
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
    await refreshCurrentPageData()

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
      currentNodeId: currentNodeId.value
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
  await refreshCurrentPageData()
  showBreakpointDialog.value = false
  await saveBreakpoint()
  ElMessage.success(`已为你跳转到第 ${breakpointPage.value} 页`)
}

const restartStudy = async () => {
  currentPage.value = 1
  await refreshCurrentPageData()
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
  background: linear-gradient(180deg, #f6f9ff 0%, #eef3fb 100%);
  display: flex;
  flex-direction: column;
  font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;
}
.main {
  flex: 1;
  display: flex;
  gap: 16px;
  padding: 12px 20px 20px;
  overflow: hidden;
}
.left-section {
  flex: 6;
}
.right-section {
  flex: 4;
  overflow: hidden;
}
.smart-tab {
  height: 100%;
  display: flex;
  flex-direction: column;
}
:deep(.el-tabs__content) {
  flex: 1;
  overflow: auto;
}
.footer {
  height: 40px;
  background: rgba(255, 255, 255, 0.85);
  text-align: center;
  line-height: 40px;
  font-size: 12px;
  color: #999;
  border-top: 1px solid #eee;
}
</style>