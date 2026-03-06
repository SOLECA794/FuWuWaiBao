<template>
  <div class="teacher-app">
    <TeacherTopBar
      :backend-status-class="backendStatusClass"
      :backend-status-text="backendStatusText"
    />
    <TeacherOverviewStrip
      :current-course-name="currentCourseName"
      :current-edit-page="currentEditPage"
      :current-course-total-pages="currentCourseTotalPages"
      :current-course-published="currentCoursePublished"
    />

    <div class="main-content">
      <TeacherCoursewareSidebar
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
      />

      <div class="editor-section">
        <div class="tabs">
          <button class="tab-btn" :class="{ active: activeTab === 'script' }" @click="activeTab = 'script'">讲稿编辑</button>
          <button class="tab-btn" :class="{ active: activeTab === 'stats' }" @click="activeTab = 'stats'">学情分析</button>
          <button class="tab-btn" :class="{ active: activeTab === 'questions' }" @click="activeTab = 'questions'">提问统计</button>
          <button class="tab-btn" :class="{ active: activeTab === 'card' }" @click="activeTab = 'card'">学习卡点可视化</button>
        </div>

        <TeacherScriptPanel
          v-if="activeTab === 'script'"
          :preview-url="previewUrl"
          :current-course-id="currentCourseId"
          :current-edit-page="currentEditPage"
          :current-script="currentScript"
          :script-generating="scriptGenerating"
          :script-saving="scriptSaving"
          @generate-ai-script="generateAIScript"
          @save-script="saveScript"
          @update:current-script="currentScript = $event"
        />

        <TeacherStatsPanel
          v-else-if="activeTab === 'stats'"
          :current-course-id="currentCourseId"
          :current-course-name="currentCourseName"
          :student-stats="studentStats"
        />

        <TeacherQuestionsPanel
          v-else-if="activeTab === 'questions'"
          :current-course-id="currentCourseId"
          :current-course-name="currentCourseName"
          :current-course-total-pages="currentCourseTotalPages"
          :filter-page="filterPage"
          :filtered-questions="filteredQuestions"
          @update:filter-page="filterPage = $event"
        />

        <TeacherCardAnalysisPanel
          v-else
          :current-course-id="currentCourseId"
          :current-course-name="currentCourseName"
          :chart-type="chartType"
          :card-data="cardData"
          @update:chart-type="chartType = $event"
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
    />

    <TeacherPublishModal
      :visible="showPublishModal"
      :current-course-name="currentCourseName"
      :publish-scope="publishScope"
      @close="showPublishModal = false"
      @submit="publishCourseware"
      @update:publish-scope="publishScope = $event"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue'
import { API_BASE } from './config/api'
import { teacherV1Api } from './services/v1'
import TeacherTopBar from './components/teacher/TeacherTopBar.vue'
import TeacherOverviewStrip from './components/teacher/TeacherOverviewStrip.vue'
import TeacherCoursewareSidebar from './components/teacher/TeacherCoursewareSidebar.vue'
import TeacherScriptPanel from './components/teacher/TeacherScriptPanel.vue'
import TeacherStatsPanel from './components/teacher/TeacherStatsPanel.vue'
import TeacherQuestionsPanel from './components/teacher/TeacherQuestionsPanel.vue'
import TeacherCardAnalysisPanel from './components/teacher/TeacherCardAnalysisPanel.vue'
import TeacherUploadModal from './components/teacher/TeacherUploadModal.vue'
import TeacherPublishModal from './components/teacher/TeacherPublishModal.vue'

const coursewareList = ref([])
const currentCourseId = ref('')
const currentCourseName = ref('')
const currentCourseTotalPages = ref(0)
const currentEditPage = ref(1)
const currentScript = ref('')

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

const activeTab = ref('script')
const studentStats = ref({
  totalQuestions: 0,
  hotPages: [],
  keyDifficulties: '暂无'
})

const questionRecords = ref([])
const filterPage = ref('')
const filteredQuestions = computed(() => {
  if (!filterPage.value) return questionRecords.value
  return questionRecords.value.filter(q => q.page === Number(filterPage.value))
})
const currentCoursePublished = computed(() => {
  const currentCourse = coursewareList.value.find(item => item.id === currentCourseId.value)
  return !!currentCourse?.published
})
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
const previewUrl = computed(() => {
  if (!currentCourseId.value) return ''
  return `${API_BASE}/api/courseware/${currentCourseId.value}/page/${currentEditPage.value}`
})

const chartType = ref('bar')
let backendHealthTimer = null
const cardData = ref([])

onMounted(async () => {
  await checkBackendHealth()
  backendHealthTimer = window.setInterval(checkBackendHealth, 30000)
  await loadCoursewareList()
  if (currentCourseId.value) {
    await loadCourseContext(currentCourseId.value)
  }
})

onUnmounted(() => {
  if (backendHealthTimer) {
    window.clearInterval(backendHealthTimer)
    backendHealthTimer = null
  }
})

const checkBackendHealth = async () => {
  try {
    const res = await teacherV1Api.health()
    backendStatus.value = res.ok ? 'online' : 'offline'
  } catch (error) {
    backendStatus.value = 'offline'
  }
}

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
  await loadScript(courseId, 1)
  await loadStudentStats(courseId)
  await loadCardData(courseId)
  await loadQuestionRecords(courseId)
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
  currentEditPage.value = page
  await loadScript(currentCourseId.value, page)
}

const loadScript = async (courseId, page) => {
  try {
    const data = await teacherV1Api.coursewares.getScript(courseId, page)
    currentScript.value = data?.data?.content || ''
  } catch (err) {
    currentScript.value = ''
  }
}

const saveScript = async () => {
  if (!currentScript.value.trim()) {
    alert('请输入讲稿内容！')
    return
  }
  scriptSaving.value = true
  try {
    await teacherV1Api.coursewares.saveScript({
      courseId: currentCourseId.value,
      pageNum: currentEditPage.value,
      content: currentScript.value
    })
    alert('讲稿保存成功！')
  } catch (err) {
    alert('保存讲稿失败：' + err.message)
  } finally {
    scriptSaving.value = false
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
  } catch (err) {
    currentScript.value = '生成失败：' + err.message
  } finally {
    scriptGenerating.value = false
  }
}

const handleFileSelect = (event) => {
  const file = event.target.files?.[0]
  if (file) {
    selectedFileName.value = file.name
  }
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
    alert('课件上传并解析成功！')
    showUploadModal.value = false
    selectedFileName.value = ''
    await loadCoursewareList()
  } catch (err) {
    alert('上传失败：' + err.message)
  } finally {
    uploadLoading.value = false
    if (fileInput.value) {
      fileInput.value.value = ''
    }
  }
}

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

const loadStudentStats = async (courseId) => {
  try {
    const data = await teacherV1Api.analytics.getStats(courseId)
    const payload = data?.data || {}
    studentStats.value = {
      totalQuestions: payload.totalQuestions || 0,
      hotPages: (payload.pageStats || []).map(item => item.page).slice(0, 3),
      keyDifficulties: (payload.keywords || []).map(item => item.word).slice(0, 3).join('、') || '暂无'
    }
  } catch (err) {
    studentStats.value = {
      totalQuestions: 0,
      hotPages: [],
      keyDifficulties: '加载失败'
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
    console.error('加载卡点数据失败', err)
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
      content: item.question || '',
      answer: item.answer || '',
      time: item.created_at ? new Date(item.created_at).toLocaleString() : ''
    }))
  } catch (err) {
    questionRecords.value = []
    console.error('加载提问记录失败', err)
  }
}
</script>

<style scoped>
.teacher-app {
  width: 100%;
  height: 100vh;
  overflow: hidden;
  font-family: 'PingFang SC', 'Microsoft YaHei', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  background: linear-gradient(180deg, #f4f8ff 0%, #eef3fb 100%);
}
.main-content {
  display: flex;
  height: calc(100vh - 108px);
}
.editor-section {
  flex: 7;
  padding: 20px;
  overflow: auto;
}
.tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
}
.tab-btn {
  border: none;
  border-radius: 10px;
  padding: 10px 14px;
  background: #dbe7ff;
  color: #1e3a8a;
  cursor: pointer;
  font-weight: 500;
}
.tab-btn.active {
  background: #2563eb;
  color: #fff;
  box-shadow: 0 8px 18px rgba(37, 99, 235, 0.2);
}
</style>