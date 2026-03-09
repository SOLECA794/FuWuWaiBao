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
          <button class="tab-btn" :class="{ active: activeTab === 'preview' }" @click="activeTab = 'preview'">课件预览</button>
          <button class="tab-btn" :class="{ active: activeTab === 'stats' }" @click="activeTab = 'stats'">学情分析</button>
          <button class="tab-btn" :class="{ active: activeTab === 'questions' }" @click="activeTab = 'questions'">提问统计</button>
          <button class="tab-btn" :class="{ active: activeTab === 'card' }" @click="activeTab = 'card'">学习卡点可视化</button>
        </div>

        <!-- 预览面板：完全移除404接口请求，改为本地状态提示 -->
        <div v-if="activeTab === 'preview'" class="preview-panel">
          <div class="preview-header">
            <h3>{{ currentCourseName }} - 第{{ currentEditPage }}页</h3>
            <div class="preview-controls">
              <button 
                class="preview-btn" 
                @click="prevPage" 
                :disabled="currentEditPage <= 1"
              >上一页</button>
              <button 
                class="preview-btn" 
                @click="nextPage" 
                :disabled="currentEditPage >= currentCourseTotalPages"
              >下一页</button>
            </div>
          </div>
          <div class="preview-content">
            <!-- 本地加载动画（不发请求） -->
            <div v-if="previewLoading" class="preview-loading">
              <div class="spinner"></div>
              <p>准备预览内容...</p>
            </div>

            <!-- 核心修改：动态绑定图片URL，确保插值生效 -->
            <div v-else class="mock-preview" v-show="currentCourseId">
              <img 
                :src="mockPreviewUrl" 
                class="preview-img"
                alt="课件预览占位图"
                @load="onPreviewLoad"
                @error="handlePreviewError"
              >
              <div class="preview-tip">
                <p>📌 预览功能说明</p>
                <p>当前后端预览接口暂未实现（404），此处为模拟预览图。</p>
                <p>课件ID: {{ currentCourseId }} | 页码: {{ currentEditPage }}</p>
              </div>
            </div>

            <!-- 无课件提示 -->
            <div v-if="!currentCourseId && !previewLoading" class="preview-empty">
              请先在左侧选择或上传一个课件
            </div>
          </div>
        </div>

        <!-- 修复：补全闭合标签 + 移除属性行注释 -->
        <TeacherScriptPanel
          v-else-if="activeTab === 'script'"
          :preview-url="realPreviewUrl"
          :current-course-id="currentCourseId"
          :current-edit-page="currentEditPage"
          :current-script="currentScript"
          :script-generating="scriptGenerating"
          :script-saving="scriptSaving"
          @generate-ai-script="generateAIScript"
          @save-script="saveScript"
          @update:current-script="currentScript = $event"
        ></TeacherScriptPanel>

        <TeacherStatsPanel
          v-else-if="activeTab === 'stats'"
          :current-course-id="currentCourseId"
          :current-course-name="currentCourseName"
          :student-stats="studentStats"
        ></TeacherStatsPanel>

        <TeacherQuestionsPanel
          v-else-if="activeTab === 'questions'"
          :current-course-id="currentCourseId"
          :current-course-name="currentCourseName"
          :current-course-total-pages="currentCourseTotalPages"
          :filter-page="filterPage"
          :filtered-questions="filteredQuestions"
          @update:filter-page="filterPage = $event"
        ></TeacherQuestionsPanel>

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
import { ref, onMounted, computed, onUnmounted } from 'vue'
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

// --- 状态管理 ---
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
const previewLoading = ref(false) // 仅用于本地动画，不发请求

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

// 新增：真实预览URL（假设后端返回图片）
const realPreviewUrl = computed(() => {
  if (!currentCourseId.value) return ''
  return `${API_BASE}/courseware/${currentCourseId.value}/page/${currentEditPage.value}?t=${Date.now()}`
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
  await checkBackendHealth()
  backendHealthTimer = window.setInterval(checkBackendHealth, 30000)
  await loadCoursewareList()
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
  } catch (err) {
    currentScript.value = ''
  }
}

const saveScript = async () => {
  if (!currentScript.value.trim()) return alert('请输入讲稿内容！')
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
    alert('课件上传成功！（预览接口待后端实现）')
    showUploadModal.value = false
    selectedFileName.value = ''
    await loadCoursewareList()
    activeTab.value = 'preview' // 上传后切到预览页
  } catch (err) {
    alert('上传失败：' + err.message)
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
      keyDifficulties: (payload.keywords || []).map(item => item.word).slice(0, 3).join('、') || '暂无'
    }
  } catch (err) {
    studentStats.value = { totalQuestions: 0, hotPages: [], keyDifficulties: '加载失败' }
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
      content: item.question || '',
      answer: item.answer || '',
      time: item.created_at ? new Date(item.created_at).toLocaleString() : ''
    }))
  } catch (err) {
    questionRecords.value = []
  }
}
</script>

<style scoped>
.teacher-app {
  width: 100%;
  height: 100vh;
  overflow: hidden;
  font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;
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
  transition: all 0.2s;
}
.tab-btn.active {
  background: #2563eb;
  color: #fff;
  box-shadow: 0 8px 18px rgba(37, 99, 235, 0.2);
}

/* --- 预览面板样式优化 --- */
.preview-panel {
  width: 100%;
  height: calc(100% - 40px);
  display: flex;
  flex-direction: column;
}
.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e2e8f0;
}
.preview-header h3 {
  margin: 0;
  color: #1e3a8a;
  font-size: 18px;
}
.preview-controls {
  display: flex;
  gap: 8px;
}
.preview-btn {
  border: none;
  border-radius: 6px;
  padding: 6px 12px;
  background: #2563eb;
  color: #fff;
  cursor: pointer;
  border: 1px solid transparent;
}
.preview-btn:disabled {
  background: #94a3b8;
  cursor: not-allowed;
}
.preview-content {
  flex: 1;
  width: 100%;
  height: 100%;
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
  background: #ffffff;
}

/* 本地加载动画 */
.preview-loading {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: #f8fafc;
  color: #64748b;
}
.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e2e8f0;
  border-top: 4px solid #2563eb;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 12px;
}
@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 空状态 */
.preview-empty {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #94a3b8;
  font-size: 16px;
  text-align: center;
}
</style>