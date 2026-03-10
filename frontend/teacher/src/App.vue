<template>
  <div class="teacher-app">
    <!-- 顶部导航栏：保留 #ddf5f5 浅绿 -->
    <div style="background-color: #ddf5f5 !important; color: #2d3748 !important;">
      <TeacherTopBar
        :backend-status-class="backendStatusClass"
        :backend-status-text="backendStatusText"
      />
    </div>

    <!-- 标题栏：强制纯白 -->
    <div style="background-color: #ffffff !important; border-bottom: 1px solid #e2e8f0 !important;">
      <TeacherOverviewStrip
        :current-course-name="currentCourseName"
        :current-edit-page="currentEditPage"
        :current-course-total-pages="currentCourseTotalPages"
        :current-course-published="currentCoursePublished"
      />
    </div>

    <!-- 主布局 -->
    <div class="main-layout">
      <!-- 左侧菜单：纯白 -->
      <div class="left-sidebar" style="background-color: #ffffff !important;">
        <div class="menu-list">
          <button class="menu-btn" :class="{ active: activeTab === 'script' }" @click="activeTab = 'script'" style="background-color: #ffffff !important;">
            <span class="menu-icon icon-edit"></span>
            <span class="menu-text">讲稿编辑</span>
          </button>
          <button class="menu-btn" :class="{ active: activeTab === 'stats' }" @click="activeTab = 'stats'" style="background-color: #ffffff !important;">
            <span class="menu-icon icon-chart"></span>
            <span class="menu-text">学情分析</span>
          </button>
          <button class="menu-btn" :class="{ active: activeTab === 'questions' }" @click="activeTab = 'questions'" style="background-color: #ffffff !important;">
            <span class="menu-icon icon-question"></span>
            <span class="menu-text">提问统计</span>
          </button>
          <button class="menu-btn" :class="{ active: activeTab === 'card' }" @click="activeTab = 'card'" style="background-color: #ffffff !important;">
            <span class="menu-icon icon-card"></span>
            <span class="menu-text">学习卡点可视化</span>
          </button>
        </div>
      </div>

      <!-- 中间预览区（选择排序复习）：强制纯白 -->
      <div class="middle-preview" style="background-color: #ffffff !important; border: none !important; box-shadow: none !important;">
        <div class="preview-panel">
          <div class="preview-header">
            <h3>{{ currentCourseName }} - 第{{ currentEditPage }}页</h3>
            <div class="preview-controls">
              <button class="preview-btn" @click="prevPage" :disabled="currentEditPage <= 1">上一页</button>
              <button class="preview-btn" @click="nextPage" :disabled="currentEditPage >= currentCourseTotalPages">下一页</button>
            </div>
          </div>
          <div class="preview-content" style="background-color: #ffffff !important; background-image: none !important;">
            <div v-if="previewLoading" class="preview-loading" style="background-color: #ffffff !important;">
              <div class="spinner"></div>
              <p>准备预览内容...</p>
            </div>
            <div v-else class="mock-preview" v-show="currentCourseId" style="background-color: #ffffff !important;">
              <img :src="mockPreviewUrl" class="preview-img" alt="课件预览占位图" @load="onPreviewLoad" @error="handlePreviewError">
              <div class="preview-tip">
                <p>📌 预览功能说明</p>
                <p>当前后端预览接口暂未实现（404），此处为模拟预览图。</p>
                <p>课件ID: {{ currentCourseId }} | 页码: {{ currentEditPage }}</p>
              </div>
            </div>
            <div v-if="!currentCourseId && !previewLoading" class="preview-empty" style="background-color: #ffffff !important;">请先在右侧选择或上传一个课件</div>
          </div>
        </div>
        <div class="function-panel" v-if="activeTab !== 'preview'" style="background-color: #ffffff !important;">
          <TeacherScriptPanel v-if="activeTab === 'script'" :preview-url="realPreviewUrl" :current-course-id="currentCourseId" :current-edit-page="currentEditPage" :current-script="currentScript" :script-generating="scriptGenerating" :script-saving="scriptSaving" @generate-ai-script="generateAIScript" @save-script="saveScript" @update:current-script="currentScript = $event"/>
          <TeacherStatsPanel v-if="activeTab === 'stats'" :current-course-id="currentCourseId" :current-course-name="currentCourseName" :student-stats="studentStats"/>
          <TeacherQuestionsPanel v-if="activeTab === 'questions'" :current-course-id="currentCourseId" :current-course-name="currentCourseName" :current-course-total-pages="currentCourseTotalPages" :filter-page="filterPage" :filtered-questions="filteredQuestions" @update:filter-page="filterPage = $event"/>
          <TeacherCardAnalysisPanel v-if="activeTab === 'card'" :current-course-id="currentCourseId" :current-course-name="currentCourseName" :chart-type="chartType" :card-data="cardData" @update:chart-type="chartType = $event"/>
        </div>
      </div>

      <!-- 右侧课件管理区：强制纯白 -->
      <div class="right-sidebar" style="background-color: #ffffff !important;">
        <TeacherCoursewareSidebar :courseware-list="coursewareList" :current-course-id="currentCourseId" :current-course-name="currentCourseName" :current-course-total-pages="currentCourseTotalPages" :current-edit-page="currentEditPage" :course-list-loading="courseListLoading" @open-publish="showPublishModal = true" @open-upload="showUploadModal = true" @select-course="selectCourse" @delete-course="deleteCourse" @select-page="selectEditPage" style="background-color: #ffffff !important;"/>
        <div class="ai-script-box" style="background-color: #ffffff !important; border: 1px solid #e2e8f0 !important; box-shadow: none !important;">
          <h4>AI讲稿生成</h4>
          <textarea class="script-input" v-model="currentScript" placeholder="请输入讲稿草稿，或点击下方按钮AI生成..." :disabled="scriptGenerating" style="background-color: #ffffff !important;"></textarea>
          <div class="script-btns">
            <button class="ai-generate-btn" @click="generateAIScript" :disabled="!currentCourseId || scriptGenerating" style="background-color: #2563eb !important; color: #fff !important;">{{ scriptGenerating ? '生成中...' : 'AI生成讲稿' }}</button>
            <button class="save-script-btn" @click="saveScript" :disabled="!currentScript.trim() || scriptSaving" style="background-color: #10b981 !important; color: #fff !important;">{{ scriptSaving ? '保存中...' : '保存讲稿' }}</button>
          </div>
        </div>
      </div>
    </div>

    <TeacherUploadModal :visible="showUploadModal" :selected-file-name="selectedFileName" :upload-loading="uploadLoading" @close="showUploadModal = false" @select-file="handleFileSelect" @submit="uploadCourseware" @file-input-ready="fileInput = $event"/>
    <TeacherPublishModal :visible="showPublishModal" :current-course-name="currentCourseName" :publish-scope="publishScope" @close="showPublishModal = false" @submit="publishCourseware" @update:publish-scope="publishScope = $event"/>
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

// 状态管理
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
const studentStats = ref({ totalQuestions: 0, hotPages: [], keyDifficulties: '暂无' })
const questionRecords = ref([])
const filterPage = ref('')
const previewLoading = ref(false)

// 计算属性
const filteredQuestions = computed(() => filterPage.value ? questionRecords.value.filter(q => q.page === Number(filterPage.value)) : questionRecords.value)
const currentCoursePublished = computed(() => coursewareList.value.find(item => item.id === currentCourseId.value)?.published || false)
const backendStatusText = computed(() => backendStatus.value === 'online' ? '在线' : '离线')
const backendStatusClass = computed(() => backendStatus.value)

// 预览URL
const realPreviewUrl = computed(() => currentCourseId.value ? `${API_BASE}/courseware/${currentCourseId.value}/page/${currentEditPage.value}?t=${Date.now()}` : '')
const mockPreviewUrl = computed(() => currentCourseId.value ? `https://picsum.photos/800/600?random=${currentCourseId.value}_${currentEditPage.value}` : '')

// 预览回调
const onPreviewLoad = () => previewLoading.value = false
const handlePreviewError = () => {
  previewLoading.value = false
  console.error('预览图片加载失败')
  alert('模拟预览图片加载失败，请检查网络')
}

// 生命周期
// 生命周期
let backendHealthTimer = null
onMounted(async () => {
  await checkBackendHealth()
  backendHealthTimer = setInterval(checkBackendHealth, 30000)
  
  // 强制先加载课件列表
  await loadCoursewareList()

  // 等待列表加载完成后，再自动选第一个课件（修复刷新消失）
  if (coursewareList.value.length > 0) {
    const first = coursewareList.value[0]
    await selectCourse(first)
  }
})
onUnmounted(() => {
  if (backendHealthTimer) clearInterval(backendHealthTimer)
})

// 核心方法
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
    coursewareList.value = data.data.map(item => ({ id: item.id, name: item.title, totalPages: item.total_page, published: !!item.is_published }))
    if (coursewareList.value.length) {
      currentCourseId.value = coursewareList.value[0].id
      currentCourseName.value = coursewareList.value[0].name
      currentCourseTotalPages.value = coursewareList.value[0].totalPages
    }
  } catch (err) {
    console.error('加载课件列表失败', err)
  } finally {
    courseListLoading.value = false
  }
}

const loadCourseContext = async (courseId) => {
  previewLoading.value = true
  try {
    await Promise.all([
      loadScript(courseId, 1),
      loadStudentStats(courseId),
      loadCardData(courseId),
      loadQuestionRecords(courseId)
    ])
  } catch (err) {
    console.error('加载课程上下文失败', err)
  } finally {
    setTimeout(() => previewLoading.value = false, 500)
  }
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
    if (courseId === currentCourseId.value) {
      currentCourseId.value = ''
      currentCourseName.value = ''
      currentCourseTotalPages.value = 0
      currentScript.value = ''
      questionRecords.value = []
      studentStats.value = { totalQuestions: 0, hotPages: [], keyDifficulties: '暂无' }
    }
  } catch (err) {
    alert('课件删除失败：' + err.message)
  }
}

const selectEditPage = async (page) => {
  previewLoading.value = true
  currentEditPage.value = page
  await loadScript(currentCourseId.value, page)
  setTimeout(() => previewLoading.value = false, 300)
}

const prevPage = () => currentEditPage.value > 1 && selectEditPage(currentEditPage.value - 1)
const nextPage = () => currentEditPage.value < currentCourseTotalPages.value && selectEditPage(currentEditPage.value + 1)

// 讲稿相关
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
    await teacherV1Api.coursewares.saveScript({ courseId: currentCourseId.value, pageNum: currentEditPage.value, content: currentScript.value })
    alert('讲稿保存成功！')
  } catch (err) {
    alert('讲稿保存失败：' + err.message)
  } finally {
    scriptSaving.value = false
  }
}

const generateAIScript = async () => {
  scriptGenerating.value = true
  try {
    currentScript.value = 'AI正在生成讲稿...'
    const data = await teacherV1Api.coursewares.generateScript({ courseId: currentCourseId.value, pageNum: currentEditPage.value })
    currentScript.value = data?.data?.content || 'AI生成失败，请重试'
  } catch (err) {
    currentScript.value = '生成失败：' + err.message
  } finally {
    scriptGenerating.value = false
  }
}

// 上传相关
const handleFileSelect = (event) => {
  const file = event.target.files?.[0]
  if (file) selectedFileName.value = file.name
}

const uploadCourseware = async () => {
  const file = fileInput.value?.files?.[0]
  if (!file) return
  uploadLoading.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('title', file.name.replace(/\.[^.]+$/, ''))
    await teacherV1Api.coursewares.upload(formData)
    alert('课件上传成功！（预览接口待后端实现）')
    showUploadModal.value = false
    selectedFileName.value = ''
    await loadCoursewareList()
    activeTab.value = 'preview'
  } catch (err) {
    alert('上传失败：' + err.message)
  } finally {
    uploadLoading.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}

// 发布相关
const publishCourseware = async () => {
  try {
    await teacherV1Api.coursewares.publish({ courseId: currentCourseId.value, scope: publishScope.value })
    const course = coursewareList.value.find(c => c.id === currentCourseId.value)
    if (course) course.published = true
    alert('课件发布成功！学生端已可查看。')
    showPublishModal.value = false
  } catch (err) {
    alert('发布失败：' + err.message)
  }
}

// 数据加载
const loadStudentStats = async (courseId) => {
  try {
    const data = await teacherV1Api.analytics.getStats(courseId)
    const payload = data?.data || {}
    studentStats.value = {
      totalQuestions: payload.totalQuestions || 0,
      hotPages: payload.pageStats?.map(item => item.page).slice(0, 3) || [],
      keyDifficulties: payload.keywords?.map(item => item.word).slice(0, 3).join('、') || '暂无'
    }
  } catch (err) {
    studentStats.value = { totalQuestions: 0, hotPages: [], keyDifficulties: '加载失败' }
  }
}

const loadCardData = async (courseId) => {
  try {
    const data = await teacherV1Api.analytics.getCardData(courseId)
    const pageStats = data?.data?.pageStats || []
    studentStats.value = pageStats.map(item => ({ page: item.page, 提问量: item.questionCount, 停留时长: item.stayTime, 卡点指数: item.cardIndex }))
  } catch (err) {
    studentStats.value = []
  }
}

const loadQuestionRecords = async (courseId) => {
  try {
    const data = await teacherV1Api.analytics.getQuestionRecords(courseId, 1, 100)
    questionRecords.value = data?.data?.list.map(item => ({
      id: item.id,
      studentId: item.user_id || item.userId || '未知',
      page: item.page_index || item.pageIndex || 1,
      content: item.question || '',
      answer: item.answer || '',
      time: item.created_at ? new Date(item.created_at).toLocaleString() : ''
    })) || []
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

.main-layout {
  display: flex;
  height: calc(100vh - 108px);
  gap: 16px;
  padding: 0 16px;
}

.left-sidebar {
  width: 180px;
  background: #ffffff;
  border-radius: 8px;
  padding: 20px 0;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.menu-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 0 12px;
}

.menu-btn {
  width: 100%;
  border: none;
  border-radius: 6px;
  padding: 12px 8px;
  background: #ffffff;
  color: #666666;
  cursor: pointer;
  font-size: 14px;
  text-align: left;
  padding-left: 16px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 12px;
}

.menu-btn.active {
  background: #ddf5f5;
  color: #34beae;
  font-weight: 500;
}

.menu-btn:hover:not(.active) {
  background: #f5f5f5;
}

.menu-icon {
  width: 16px;
  height: 16px;
  display: inline-block;
  position: relative;
  flex-shrink: 0;
}

.menu-btn .menu-icon {
  color: #666666;
}

.menu-btn.active .menu-icon {
  color: #34beae;
}

/* 图标样式 */
.icon-edit::after {
  content: '';
  position: absolute;
  width: 14px;
  height: 14px;
  border: 1px solid currentColor;
  border-radius: 2px;
  top: 0;
  left: 0;
}
.icon-edit::before {
  content: '';
  position: absolute;
  width: 8px;
  height: 2px;
  background: currentColor;
  transform: rotate(45deg);
  top: 12px;
  left: 10px;
  border-radius: 1px;
}

.icon-chart::before {
  content: '';
  position: absolute;
  width: 16px;
  height: 10px;
  border-bottom: 1px solid currentColor;
  border-left: 1px solid currentColor;
  bottom: 0;
  left: 0;
}
.icon-chart::after {
  content: '';
  position: absolute;
  width: 3px;
  height: 8px;
  background: currentColor;
  bottom: 0;
  left: 2px;
  border-radius: 1px 1px 0 0;
}
.icon-chart span {
  position: absolute;
  width: 3px;
  height: 5px;
  background: currentColor;
  bottom: 0;
  left: 7px;
  border-radius: 1px 1px 0 0;
}
.icon-chart span::after {
  content: '';
  position: absolute;
  width: 3px;
  height: 10px;
  background: currentColor;
  bottom: 0;
  left: 5px;
  border-radius: 1px 1px 0 0;
}

.icon-question::before {
  content: '?';
  position: absolute;
  font-size: 16px;
  font-weight: bold;
  line-height: 16px;
  text-align: center;
  width: 100%;
  height: 100%;
}

.icon-card::before {
  content: '';
  position: absolute;
  width: 14px;
  height: 10px;
  border: 1px solid currentColor;
  border-radius: 2px;
  top: 2px;
  left: 1px;
}
.icon-card::after {
  content: '';
  position: absolute;
  width: 10px;
  height: 1px;
  background: currentColor;
  top: 5px;
  left: 3px;
}

.menu-text {
  flex: 1;
}

/* 中间预览区 */
.middle-preview {
  flex: 1;
  background: #ffffff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.preview-panel {
  width: 100%;
  height: 100%;
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

.mock-preview {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  box-sizing: border-box;
}

.preview-img {
  max-width: 100%;
  max-height: 70%;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  margin-bottom: 15px;
}

.preview-tip {
  text-align: center;
  color: #64748b;
  font-size: 14px;
  line-height: 1.6;
}

.preview-tip p:first-child {
  font-weight: bold;
  color: #2563eb;
  margin-bottom: 5px;
}

.preview-empty {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #94a3b8;
  font-size: 16px;
  text-align: center;
}

.function-panel {
  width: 100%;
  height: 100%;
  overflow: auto;
}

/* 右侧功能区 */
.right-sidebar {
  width: 300px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.ai-script-box {
  background: #ffffff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.ai-script-box h4 {
  margin: 0 0 12px 0;
  color: #1e3a8a;
  font-size: 16px;
}

.script-input {
  width: 100%;
  height: 120px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 8px 12px;
  font-size: 14px;
  resize: none;
  margin-bottom: 12px;
  box-sizing: border-box;
}

.script-input:disabled {
  background: #f8fafc;
  color: #94a3b8;
}

.script-btns {
  display: flex;
  gap: 8px;
}

.ai-generate-btn {
  flex: 1;
  border: none;
  border-radius: 6px;
  padding: 8px 0;
  background: #2563eb;
  color: #fff;
  cursor: pointer;
}

.ai-generate-btn:disabled {
  background: #94a3b8;
  cursor: not-allowed;
}

.save-script-btn {
  flex: 1;
  border: none;
  border-radius: 6px;
  padding: 8px 0;
  background: #10b981;
  color: #fff;
  cursor: pointer;
}

.save-script-btn:disabled {
  background: #94a3b8;
  cursor: not-allowed;
}
</style>