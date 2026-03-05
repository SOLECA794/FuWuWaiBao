<template>
  <div class="teacher-app">
    <!-- 顶部导航栏 -->
    <el-header height="56px" class="top-nav">
      <div class="nav-left">
        <span class="app-title">智能互动教学平台 · 教师端</span>
      </div>
      <div class="nav-right">
        <el-avatar :size="36" src="https://picsum.photos/id/1005/40/40">
          <el-icon><UserFilled /></el-icon>
        </el-avatar>
        <span class="teacher-name ml-2">教师 2025T001</span>
      </div>
    </el-header>

    <!-- 主体内容 -->
    <el-container class="main-container">
      <!-- 左侧课件管理区 -->
      <el-aside width="320px" class="courseware-manage-section">
        <div class="section-header">
          <h3>课件管理</h3>
          <div class="header-actions">
            <el-button 
              type="success" 
              size="small" 
              @click="showPublishModal = true"
              :disabled="!currentCourseId"
            >
              <el-icon><UploadFilled /></el-icon> 发布课件
            </el-button>
            <el-button 
              type="primary" 
              size="small" 
              @click="showUploadModal = true"
            >
              <el-icon><Plus /></el-icon> 上传课件
            </el-button>
          </div>
        </div>

        <!-- 课件列表 -->
        <el-card shadow="hover" class="courseware-list-card mt-2">
          <el-list border :data="coursewareList" empty-text="暂无课件，请点击上方按钮上传">
            <el-list-item 
              v-for="course in coursewareList" 
              :key="course.courseId"
              :class="{ active: course.courseId === currentCourseId }"
              @click="selectCourse(course)"
            >
              <el-list-item-meta>
                <el-list-item-meta-title>
                  {{ course.title }}
                  <el-tag 
                    :type="getStatusTagType(course.status)" 
                    size="small" 
                    class="ml-2"
                  >
                    {{ getStatusText(course.status) }}
                  </el-tag>
                </el-list-item-meta-title>
                <el-list-item-meta-description>
                  {{ course.fileType }} | {{ formatTime(course.createdAt) }}
                </el-list-item-meta-description>
              </el-list-item-meta>
              <el-space>
                <el-tag v-if="course.published" size="small" type="info">已发布</el-tag>
                <el-button 
                  v-if="course.status === 'failed'"
                  size="mini" 
                  type="text" 
                  @click.stop="retryParseCourse(course.courseId)"
                >
                  重试解析
                </el-button>
                <el-button 
                  size="mini" 
                  type="text" 
                  text 
                  @click.stop="deleteCourse(course.courseId)"
                >
                  <el-icon color="#ff4d4f"><Delete /></el-icon>
                </el-button>
              </el-space>
            </el-list-item>
          </el-list>
        </el-card>

        <!-- 页码选择器 -->
        <el-card shadow="hover" class="page-selector-card mt-2" v-if="currentCourseId">
          <div class="page-selector-header">
            <span>当前课件：{{ currentCourseTitle }}</span>
          </div>
          <el-space wrap class="page-buttons mt-2">
            <el-button 
              v-for="page in currentCourseTotalPages"
              :key="page"
              size="small"
              :type="page === currentEditPageNum ? 'primary' : 'default'"
              @click="selectEditPage(page)"
            >
              第{{ page }}页
            </el-button>
          </el-space>
        </el-card>
      </el-aside>

      <!-- 右侧编辑区 -->
      <el-main class="editor-section">
        <el-tabs v-model="activeTab" type="card" class="editor-tabs">
          <!-- 讲稿编辑 -->
          <el-tab-pane label="讲稿编辑" name="script">
            <!-- 课件预览 -->
            <el-card class="course-preview-card mb-4">
              <template #header>
                <span>第{{ currentEditPageNum }}页课件预览</span>
              </template>
              <div v-if="currentCourseId" class="preview-container">
                <div v-if="coursewarePreviewLoading" class="preview-loading">
                  <el-skeleton :rows="8" animated />
                </div>
                <div v-else class="preview-content">
                  <div ref="previewRef" class="preview-viewport"></div>
                </div>
              </div>
              <el-empty v-else description="请先上传并选择课件"></el-empty>
            </el-card>

            <!-- 富文本编辑 -->
            <el-card>
              <template #header>
                <div class="editor-actions">
                  <el-button 
                    type="success" 
                    size="small"
                    @click="generateAIScript"
                    :loading="aiGenerating"
                  >
                    <el-icon><RobotFilled /></el-icon> AI 生成讲稿
                  </el-button>
                  <el-button 
                    type="primary" 
                    size="small"
                    @click="saveScript"
                    :loading="savingScript"
                  >
                    <el-icon><Save /></el-icon> 保存讲稿
                  </el-button>
                </div>
              </template>
              <QuillEditor
                v-model:content="currentScriptContent"
                placeholder="请输入本页讲稿内容，支持富文本编辑..."
                :options="editorOptions"
                class="script-editor"
              />
            </el-card>
          </el-tab-pane>

          <!-- 学情分析 -->
          <el-tab-pane label="学情分析" name="stats">
            <el-card v-if="currentCourseId">
              <template #header>
                <span>学情分析 - {{ currentCourseTitle }}</span>
              </template>
              <el-descriptions :column="2" border>
                <el-descriptions-item label="总提问数">
                  {{ studentStats.totalQuestions }}
                </el-descriptions-item>
                <el-descriptions-item label="高频提问页码">
                  {{ studentStats.hotPages.join('、') || '暂无' }}
                </el-descriptions-item>
                <el-descriptions-item label="重点难点" :span="2">
                  {{ studentStats.keyDifficulties }}
                </el-descriptions-item>
              </el-descriptions>
              <!-- 学情分析图表 - 加宽处理 -->
              <div class="stats-chart mt-4" ref="statsChartRef" style="width:100%;min-width:900px;height:400px;"></div>
            </el-card>
            <el-empty v-else description="请先选择一个课件查看学情数据"></el-empty>
          </el-tab-pane>

          <!-- 提问记录 -->
          <el-tab-pane label="提问记录" name="questions">
            <el-card v-if="currentCourseId">
              <template #header>
                <div class="questions-header">
                  <span>提问记录 - {{ currentCourseTitle }}</span>
                  <el-select 
                    v-model="filterPageNum" 
                    placeholder="按页码筛选" 
                    size="small"
                    class="ml-2"
                    style="width: 120px"
                  >
                    <el-option label="全部" value="" />
                    <el-option 
                      v-for="page in currentCourseTotalPages"
                      :key="page"
                      :label="`第${page}页`"
                      :value="page"
                    />
                  </el-select>
                </div>
              </template>
              <el-table 
                :data="filteredQuestions" 
                border 
                stripe
                :loading="questionsLoading"
              >
                <el-table-column prop="studentId" label="学生ID" width="100" />
                <el-table-column prop="pageNum" label="页码" width="80" />
                <el-table-column prop="content" label="提问内容" min-width="300" />
                <el-table-column prop="answer" label="AI回复" min-width="300" />
                <el-table-column prop="time" label="提问时间" width="180" />
              </el-table>
              <el-pagination
                @size-change="handlePageSizeChange"
                @current-change="handleCurrentPageChange"
                :current-page="questionsPage.current"
                :page-sizes="[10, 20, 50]"
                :page-size="questionsPage.size"
                layout="total, sizes, prev, pager, next, jumper"
                :total="questionsPage.total"
                class="mt-3 pagination"
              />
            </el-card>
            <el-empty v-else description="请先选择一个课件查看提问记录"></el-empty>
          </el-tab-pane>

          <!-- 卡点可视化 - 加宽处理 -->
          <el-tab-pane label="学习卡点可视化" name="card-analysis">
            <el-card v-if="currentCourseId">
              <template #header>
                <div class="chart-header">
                  <span>学习卡点分析 - {{ currentCourseTitle }}</span>
                  <el-space class="ml-2">
                    <span>图表类型：</span>
                    <el-button 
                      size="small"
                      :type="chartType === 'bar' ? 'primary' : 'default'"
                      @click="chartType = 'bar'; renderCardChart()"
                    >
                      柱状图
                    </el-button>
                    <el-button 
                      size="small"
                      :type="chartType === 'line' ? 'primary' : 'default'"
                      @click="chartType = 'line'; renderCardChart()"
                    >
                      折线图
                    </el-button>
                    <el-button 
                      size="small"
                      :type="chartType === 'pie' ? 'primary' : 'default'"
                      @click="chartType = 'pie'; renderCardChart()"
                    >
                      饼图
                    </el-button>
                  </el-space>
                </div>
              </template>
              <!-- 卡点可视化图表 - 加宽处理 -->
              <div 
                ref="cardChartRef" 
                class="card-chart-container"
                style="width:100%;min-width:900px;height:500px;margin:0 auto;"
              ></div>
              <el-alert 
                title="数据说明" 
                type="info" 
                :closable="false"
                class="mt-3"
              >
                <ul class="chart-tip">
                  <li>提问量：该页面学生发起的提问总数</li>
                  <li>停留时长：学生平均停留时长（秒）</li>
                  <li>卡点指数：综合提问量+停留时长计算的卡点程度（0-10）</li>
                </ul>
              </el-alert>
            </el-card>
            <el-empty v-else description="请先选择一个课件查看卡点分析"></el-empty>
          </el-tab-pane>
        </el-tabs>
      </el-main>
    </el-container>

    <!-- 上传课件弹窗 -->
    <el-dialog 
      v-model="showUploadModal" 
      title="上传课件（PDF/PPTX）" 
      width="500px"
      destroy-on-close
    >
      <el-upload
        ref="uploadRef"
        class="upload-demo"
        drag
        action="#"
        :auto-upload="false"
        :file-list="uploadFileList"
        accept=".pdf,.pptx"
        :on-change="handleFileChange"
        :on-remove="handleFileRemove"
      >
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            只能上传 PDF/PPTX 文件，且不超过 10MB
          </div>
        </template>
      </el-upload>
      <el-form :model="uploadForm" class="mt-3" label-width="80px">
        <el-form-item label="课件标题" prop="title">
          <el-input 
            v-model="uploadForm.title" 
            placeholder="可选，默认使用文件名"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showUploadModal = false">取消</el-button>
        <el-button 
          type="primary" 
          @click="uploadCourseware"
          :loading="uploading"
          :disabled="uploadFileList.length === 0"
        >
          上传并解析
        </el-button>
      </template>
    </el-dialog>

    <!-- 发布课件弹窗 -->
    <el-dialog 
      v-model="showPublishModal" 
      title="发布课件给学生端" 
      width="400px"
    >
      <el-form :model="publishForm" label-width="80px">
        <el-form-item label="当前课件">
          <el-input v-model="currentCourseTitle" disabled />
        </el-form-item>
        <el-form-item label="发布范围" prop="scope">
          <el-select v-model="publishForm.scope" placeholder="请选择发布范围">
            <el-option label="全部学生" value="all" />
            <el-option label="班级1" value="class1" />
            <el-option label="班级2" value="class2" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPublishModal = false">取消</el-button>
        <el-button 
          type="primary" 
          @click="publishCourseware"
          :loading="publishing"
        >
          确认发布
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onUnmounted, watch } from 'vue'
import { QuillEditor } from '@vueup/vue-quill'
import 'quill/dist/quill.snow.css' // 富文本编辑器样式
import * as echarts from 'echarts'
import * as pdfjsLib from 'pdfjs-dist'
import { getDocument } from 'pdfjs-dist/build/pdf'
import { ElMessageBox, ElMessage } from 'element-plus'
import { get, post, put, del } from './utils/request'

// PDF.js 配置
pdfjsLib.GlobalWorkerOptions.workerSrc = '//unpkg.com/pdfjs-dist@3.4.120/build/pdf.worker.min.js'

// ========== 核心变量 ==========
const coursewareList = ref([])
const currentCourseId = ref('')
const currentCourseTitle = ref('')
const currentCourseTotalPages = ref(0)
const currentEditPageNum = ref(1)
const currentScriptContent = ref('')

// 上传相关
const showUploadModal = ref(false)
const uploadRef = ref(null)
const uploadFileList = ref([])
const uploadForm = ref({ title: '' })
const uploading = ref(false)

// 发布相关
const showPublishModal = ref(false)
const publishForm = ref({ scope: 'all' })
const publishing = ref(false)

// 标签页
const activeTab = ref('script')

// 学情分析
const studentStats = ref({
  totalQuestions: 0,
  hotPages: [],
  keyDifficulties: ''
})

// 提问记录
const questionRecords = ref([])
const filterPageNum = ref('')
const questionsPage = ref({
  current: 1,
  size: 20,
  total: 0
})
const questionsLoading = ref(false)

// 卡点可视化
const chartType = ref('bar')
let cardChartInstance = null
let statsChartInstance = null
const cardData = ref([])

// 加载状态
const coursewarePreviewLoading = ref(false)
const aiGenerating = ref(false)
const savingScript = ref(false)

// 图表Ref
const previewRef = ref(null)
const cardChartRef = ref(null)
const statsChartRef = ref(null)

// ========== 计算属性 ==========
const filteredQuestions = computed(() => {
  if (!filterPageNum.value) return questionRecords.value
  return questionRecords.value.filter(q => q.pageNum === Number(filterPageNum.value))
})

// ========== 工具方法 ==========
const getStatusText = (status) => {
  const statusMap = {
    'parsing': '解析中',
    'ready': '解析完成',
    'failed': '解析失败'
  }
  return statusMap[status] || '未知状态'
}

const getStatusTagType = (status) => {
  const typeMap = {
    'parsing': 'warning',
    'ready': 'success',
    'failed': 'danger'
  }
  return typeMap[status] || 'info'
}

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  return new Date(timeStr).toLocaleString('zh-CN')
}

// ========== 生命周期 ==========
onMounted(async () => {
  await loadCoursewareList()
  initCharts()
})

onUnmounted(() => {
  if (cardChartInstance) cardChartInstance.dispose()
  if (statsChartInstance) statsChartInstance.dispose()
})

// 监听课件切换
watch(currentCourseId, async (newVal) => {
  if (newVal) {
    await loadCoursewareDetail(newVal)
    await loadScript(newVal, currentEditPageNum.value)
    await loadStudentStats(newVal)
    await loadQuestionRecords(newVal)
    await loadCardData(newVal)
    await loadCoursewarePreview(newVal, currentEditPageNum.value)
  } else {
    currentCourseTitle.value = ''
    currentCourseTotalPages.value = 0
    currentScriptContent.value = ''
  }
})

// 监听页码切换
watch(currentEditPageNum, async (newVal) => {
  if (currentCourseId.value) {
    await loadScript(currentCourseId.value, newVal)
    await loadCoursewarePreview(currentCourseId.value, newVal)
  }
})

// ========== 课件管理 ==========
const loadCoursewareList = async () => {
  try {
    // 模拟数据
    coursewareList.value = [
      {
        courseId: 'uuid-1001',
        title: 'Python 基础教程',
        fileType: 'PDF',
        status: 'ready',
        published: true,
        totalPages: 10,
        createdAt: '2025-03-01 10:00:00'
      },
      {
        courseId: 'uuid-1002',
        title: 'Vue3 实战开发',
        fileType: 'PPTX',
        status: 'parsing',
        published: false,
        totalPages: 8,
        createdAt: '2025-03-02 14:30:00'
      },
      {
        courseId: 'uuid-1003',
        title: 'JavaScript 高级语法',
        fileType: 'PDF',
        status: 'failed',
        published: false,
        totalPages: 12,
        createdAt: '2025-03-03 09:15:00'
      }
    ]
    if (coursewareList.value.length > 0) {
      selectCourse(coursewareList.value[0])
    }
  } catch (err) {
    console.error('加载课件列表失败', err)
  }
}

const loadCoursewareDetail = async (courseId) => {
  try {
    const course = coursewareList.value.find(item => item.courseId === courseId)
    if (course) {
      currentCourseTotalPages.value = course.totalPages || 0
    }
  } catch (err) {
    console.error('加载课件详情失败', err)
  }
}

const selectCourse = (course) => {
  currentCourseId.value = course.courseId
  currentCourseTitle.value = course.title
  currentCourseTotalPages.value = course.totalPages || 0
  currentEditPageNum.value = 1
}

const deleteCourse = async (courseId) => {
  try {
    await ElMessageBox.confirm(
      '确定删除该课件吗？删除后将无法恢复',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    coursewareList.value = coursewareList.value.filter(item => item.courseId !== courseId)
    ElMessage.success('删除成功')
    
    if (currentCourseId.value === courseId) {
      currentCourseId.value = ''
      currentCourseTitle.value = ''
      currentCourseTotalPages.value = 0
    }
  } catch (err) {
    console.error('删除课件失败', err)
  }
}

const retryParseCourse = async (courseId) => {
  try {
    const course = coursewareList.value.find(item => item.courseId === courseId)
    if (course) {
      course.status = 'parsing'
      setTimeout(() => {
        course.status = 'ready'
        ElMessage.success('解析成功')
      }, 2000)
    }
  } catch (err) {
    console.error('重试解析失败', err)
  }
}

const selectEditPage = (pageNum) => {
  currentEditPageNum.value = pageNum
}

// ========== 课件上传 ==========
const handleFileChange = (file) => {
  uploadFileList.value = [file]
  if (!uploadForm.title) {
    uploadForm.title = file.name.split('.').slice(0, -1).join('.')
  }
}

const handleFileRemove = () => {
  uploadFileList.value = []
}

const uploadCourseware = async () => {
  if (uploadFileList.length === 0) return
  
  uploading.value = true
  try {
    const newCourse = {
      courseId: `uuid-${Date.now()}`,
      title: uploadForm.title || uploadFileList[0].name.split('.').slice(0, -1).join('.'),
      fileType: uploadFileList[0].name.split('.').pop().toUpperCase(),
      status: 'parsing',
      published: false,
      totalPages: Math.floor(Math.random() * 10) + 5,
      createdAt: new Date().toLocaleString('zh-CN')
    }
    
    coursewareList.value.unshift(newCourse)
    ElMessage.success('上传成功，正在解析...')
    
    showUploadModal.value = false
    uploadFileList.value = []
    uploadForm.title = ''
    
    setTimeout(() => {
      newCourse.status = 'ready'
      ElMessage.success('课件解析完成')
    }, 3000)
  } catch (err) {
    console.error('上传课件失败', err)
  } finally {
    uploading.value = false
  }
}

// ========== 讲稿管理 ==========
const loadScript = async (courseId, pageNum) => {
  try {
    currentScriptContent.value = `第${pageNum}页讲稿内容：
欢迎学习 ${currentCourseTitle.value} 的第${pageNum}页内容，
本节主要讲解：
1. 核心知识点
2. 实战案例
3. 常见问题解答`
  } catch (err) {
    currentScriptContent.value = ''
  }
}

const saveScript = async () => {
  if (!currentScriptContent.value.trim()) {
    ElMessage.warning('请输入讲稿内容！')
    return
  }
  
  savingScript.value = true
  try {
    ElMessage.success('讲稿保存成功')
  } catch (err) {
    console.error('保存讲稿失败', err)
  } finally {
    savingScript.value = false
  }
}

const generateAIScript = async () => {
  aiGenerating.value = true
  try {
    const contentList = [
      `第${currentEditPageNum.value}页 ${currentCourseTitle.value} 详细讲解：\n`,
      '一、知识点梳理\n',
      '1. 核心概念：\n',
      '    - 定义：\n',
      '    - 应用场景：\n',
      '2. 实战案例：\n',
      '    - 案例1：\n',
      '    - 案例2：\n',
      '二、常见问题解答\n',
      '1. 问题1：\n',
      '2. 问题2：\n',
      '三、总结与拓展\n'
    ]
    
    currentScriptContent.value = ''
    let index = 0
    
    const timer = setInterval(() => {
      if (index < contentList.length) {
        currentScriptContent.value += contentList[index]
        index++
      } else {
        clearInterval(timer)
        ElMessage.success('AI 讲稿生成完成')
        aiGenerating.value = false
      }
    }, 300)
  } catch (err) {
    ElMessage.error('AI 讲稿生成失败')
    console.error('生成失败', err)
    aiGenerating.value = false
  }
}

// ========== 课件发布 ==========
const publishCourseware = async () => {
  publishing.value = true
  try {
    const course = coursewareList.value.find(item => item.courseId === currentCourseId.value)
    if (course) {
      course.published = true
      ElMessage.success('课件发布成功')
      showPublishModal.value = false
    }
  } catch (err) {
    console.error('发布课件失败', err)
  } finally {
    publishing.value = false
  }
}

// ========== 课件预览 ==========
const loadCoursewarePreview = async (courseId, pageNum) => {
  if (!previewRef.value) return
  
  coursewarePreviewLoading.value = true
  try {
    setTimeout(() => {
      previewRef.value.innerHTML = `
        <div style="width:100%;height:400px;background:#f5f5f5;display:flex;align-items:center;justify-content:center;">
          <div style="text-align:center;">
            <img src="https://picsum.photos/id/${pageNum * 10}/600/400" style="max-width:100%;max-height:400px;" />
            <p style="margin-top:10px;color:#666;">${currentCourseTitle.value} - 第${pageNum}页预览</p>
          </div>
        </div>
      `
      coursewarePreviewLoading.value = false
    }, 1000)
  } catch (err) {
    console.error('加载预览失败', err)
    previewRef.value.innerHTML = '<div class="preview-error">预览加载失败</div>'
    coursewarePreviewLoading.value = false
  }
}

// ========== 学情分析 ==========
const loadStudentStats = async (courseId) => {
  try {
    studentStats.value = {
      totalQuestions: Math.floor(Math.random() * 50) + 10,
      hotPages: [2, 4, 5, 7, 8],
      keyDifficulties: '1. 函数闭包的理解；2. 异步编程的执行顺序；3. 组件通信的多种方式'
    }
    renderStatsChart()
  } catch (err) {
    studentStats.value = {
      totalQuestions: 0,
      hotPages: [],
      keyDifficulties: '加载失败'
    }
  }
}

// ========== 提问记录 ==========
const loadQuestionRecords = async (courseId) => {
  questionsLoading.value = true
  try {
    const mockQuestions = []
    for (let i = 1; i <= 25; i++) {
      mockQuestions.push({
        studentId: `stu-${1000 + i}`,
        pageNum: Math.floor(Math.random() * currentCourseTotalPages.value) + 1,
        content: `第${Math.floor(Math.random() * currentCourseTotalPages.value) + 1}页的这个知识点不太理解，能再解释一下吗？`,
        answer: '好的，这个知识点的核心是...（AI 详细解答内容）',
        time: new Date(Date.now() - Math.random() * 86400000 * 7).toLocaleString('zh-CN')
      })
    }
    
    questionRecords.value = mockQuestions
    questionsPage.total = mockQuestions.length
  } catch (err) {
    questionRecords.value = []
    questionsPage.total = 0
  } finally {
    questionsLoading.value = false
  }
}

const handlePageSizeChange = (size) => {
  questionsPage.size = size
}

const handleCurrentPageChange = (page) => {
  questionsPage.current = page
}

// ========== 图表相关 ==========
const initCharts = () => {
  // 初始化时确保容器已渲染，避免图表缩成一团
  setTimeout(() => {
    if (cardChartRef.value) {
      cardChartInstance = echarts.init(cardChartRef.value)
    }
    if (statsChartRef.value) {
      statsChartInstance = echarts.init(statsChartRef.value)
    }
  }, 500)
}

const loadCardData = async (courseId) => {
  try {
    cardData.value = [
      { pageNum: 1, 提问量: 1, 停留时长: 20, 卡点指数: 1.2 },
      { pageNum: 2, 提问量: 5, 停留时长: 60, 卡点指数: 5.5 },
      { pageNum: 3, 提问量: 2, 停留时长: 25, 卡点指数: 2.1 },
      { pageNum: 4, 提问量: 8, 停留时长: 90, 卡点指数: 8.5 },
      { pageNum: 5, 提问量: 3, 停留时长: 45, 卡点指数: 3.8 },
      { pageNum: 6, 提问量: 7, 停留时长: 80, 卡点指数: 7.2 },
      { pageNum: 7, 提问量: 4, 停留时长: 50, 卡点指数: 4.3 },
      { pageNum: 8, 提问量: 6, 停留时长: 75, 卡点指数: 6.8 },
      { pageNum: 9, 提问量: 1, 停留时长: 18, 卡点指数: 1.0 },
      { pageNum: 10, 提问量: 0, 停留时长: 12, 卡点指数: 0.5 }
    ]
    renderCardChart()
  } catch (err) {
    console.error('加载卡点数据失败', err)
  }
}

const renderCardChart = () => {
  if (!cardChartInstance || cardData.value.length === 0) return
  cardChartInstance.clear()
  const pages = cardData.value.map(item => `第${item.pageNum}页`)
  const questionCounts = cardData.value.map(item => item.提问量)
  const stayTimes = cardData.value.map(item => item.停留时长)
  const cardScores = cardData.value.map(item => item.卡点指数)

  let option = {}

  switch (chartType.value) {
    case 'bar':
      option = {
        title: { text: '各页面学习卡点数据', left: 'center', fontSize: 16 },
        tooltip: { trigger: 'axis', textStyle: { fontSize: 14 } },
        legend: { data: ['提问量', '停留时长(秒)', '卡点指数'], top: 30, textStyle: { fontSize: 14 } },
        grid: { left: '5%', right: '5%', top: '15%', bottom: '10%' },
        xAxis: { type: 'category', data: pages, axisLabel: { fontSize: 14 } },
        yAxis: { type: 'value', axisLabel: { fontSize: 14 } },
        series: [
          { name: '提问量', type: 'bar', data: questionCounts, barWidth: '20%' },
          { name: '停留时长(秒)', type: 'bar', data: stayTimes, barWidth: '20%' },
          { name: '卡点指数', type: 'bar', data: cardScores, barWidth: '20%', itemStyle: { color: '#ff4d4f' } }
        ]
      }
      break
    case 'line':
      option = {
        title: { text: '各页面学习卡点趋势', left: 'center', fontSize: 16 },
        tooltip: { trigger: 'axis', textStyle: { fontSize: 14 } },
        legend: { data: ['提问量', '停留时长(秒)', '卡点指数'], top: 30, textStyle: { fontSize: 14 } },
        grid: { left: '5%', right: '5%', top: '15%', bottom: '10%' },
        xAxis: { type: 'category', data: pages, axisLabel: { fontSize: 14 } },
        yAxis: { type: 'value', axisLabel: { fontSize: 14 } },
        series: [
          { name: '提问量', type: 'line', data: questionCounts, lineWidth: 3, symbolSize: 8 },
          { name: '停留时长(秒)', type: 'line', data: stayTimes, lineWidth: 3, symbolSize: 8 },
          { name: '卡点指数', type: 'line', data: cardScores, lineWidth: 3, symbolSize: 8, lineStyle: { color: '#ff4d4f' }, itemStyle: { color: '#ff4d4f' } }
        ]
      }
      break
    case 'pie':
      const top5CardData = [...cardData.value].sort((a, b) => b.卡点指数 - a.卡点指数).slice(0, 5)
      option = {
        title: { text: 'TOP5 卡点页面占比', left: 'center', fontSize: 16 },
        tooltip: { trigger: 'item', textStyle: { fontSize: 14 } },
        legend: { orient: 'vertical', left: 'left', top: 'center', textStyle: { fontSize: 14 } },
        series: [
          {
            name: '卡点指数',
            type: 'pie',
            radius: ['30%', '70%'],
            center: ['60%', '50%'],
            data: top5CardData.map(item => ({
              name: `第${item.pageNum}页`,
              value: item.卡点指数
            })),
            label: {
              formatter: '{b}: {c} ({d}%)',
              fontSize: 14
            },
            itemStyle: {
              borderRadius: 5,
              borderColor: '#fff',
              borderWidth: 2
            }
          }
        ]
      }
      break
  }

  cardChartInstance.setOption(option)
  
  // 强制自适应，解决缩成一团问题
  const resizeChart = () => {
    cardChartInstance?.resize()
    statsChartInstance?.resize()
  }
  window.addEventListener('resize', resizeChart)
  // 初始化时强制resize一次
  resizeChart()
}

const renderStatsChart = () => {
  if (!statsChartInstance || studentStats.value.hotPages.length === 0) return
  
  const option = {
    title: { text: '高频提问页码分布', left: 'center', fontSize: 16 },
    tooltip: { trigger: 'axis', textStyle: { fontSize: 14 } },
    grid: { left: '5%', right: '5%', top: '15%', bottom: '10%' },
    xAxis: { type: 'category', data: studentStats.value.hotPages.map(p => `第${p}页`), axisLabel: { fontSize: 14 } },
    yAxis: { type: 'value', name: '提问次数', axisLabel: { fontSize: 14 } },
    series: [
      {
        name: '提问次数',
        type: 'bar',
        data: studentStats.value.hotPages.map(() => Math.floor(Math.random() * 10) + 1),
        barWidth: '30%',
        itemStyle: { color: '#1677ff' }
      }
    ]
  }
  
  statsChartInstance.setOption(option)
  // 强制自适应
  statsChartInstance.resize()
}

// ========== 富文本编辑器配置 ==========
const editorOptions = {
  theme: 'snow',
  modules: {
    toolbar: [
      ['bold', 'italic', 'underline', 'strike'],
      ['blockquote', 'code-block'],
      [{ 'header': [1, 2, 3, false] }],
      [{ 'list': 'ordered'}, { 'list': 'bullet' }],
      [{ 'script': 'sub'}, { 'script': 'super' }],
      [{ 'indent': '-1'}, { 'indent': '+1' }],
      [{ 'direction': 'rtl' }],
      [{ 'size': ['small', false, 'large', 'huge'] }],
      [{ 'color': [] }, { 'background': [] }],
      [{ 'font': [] }],
      [{ 'align': [] }],
      ['clean'],
      ['link', 'image']
    ]
  }
}
</script>

<style scoped>
.teacher-app {
  width: 100%;
  height: 100vh;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

/* 顶部导航 */
.top-nav {
  background: #1677ff;
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}
.app-title {
  font-size: 16px;
  font-weight: 500;
}
.nav-right {
  display: flex;
  align-items: center;
}
.teacher-name {
  margin-left: 8px;
  font-size: 14px;
}

/* 主体容器 */
.main-container {
  height: calc(100vh - 56px);
  overflow: hidden;
}

/* 左侧课件管理区 */
.courseware-manage-section {
  background: #f5f7fa;
  padding: 16px;
  overflow-y: auto;
}
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}
.section-header h3 {
  font-size: 16px;
  font-weight: 500;
  margin: 0;
}
.courseware-list-card {
  --el-card-padding: 12px;
  margin-bottom: 12px;
}
.el-list-item.active {
  background-color: #e6f7ff;
}
.page-selector-card {
  --el-card-padding: 12px;
}
.page-selector-header {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}
.page-buttons {
  width: 100%;
  flex-wrap: wrap;
  gap: 8px;
}

/* 右侧编辑区 */
.editor-section {
  padding: 16px;
  background: #f5f7fa;
  overflow-y: auto;
  height: calc(100vh - 56px);
}
.editor-tabs {
  --el-tabs-card-border-color: #e8e8e8;
  --el-tabs-active-color: #1677ff;
}
.course-preview-card {
  --el-card-padding: 16px;
  margin-bottom: 16px;
}
.preview-container {
  width: 100%;
  min-height: 300px;
}
.preview-loading {
  width: 100%;
  height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.preview-content {
  width: 100%;
  display: flex;
  justify-content: center;
}
.preview-viewport {
  max-width: 100%;
  max-height: 500px;
}
.script-editor {
  min-height: 400px;
  --el-quill-editor-border-color: #e8e8e8;
}
.questions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}
.pagination {
  text-align: right;
  margin-top: 16px;
}
/* 图表容器样式 - 加宽处理 */
.card-chart-container {
  background: white;
  padding: 16px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}
.stats-chart {
  background: white;
  padding: 16px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}
.chart-tip {
  margin: 8px 0 0 20px;
  padding: 0;
  font-size: 14px;
  color: #666;
}

/* 响应式适配 - 保证图表不挤压 */
@media (min-width: 1200px) {
  .editor-section {
    padding: 24px;
  }
  .card-chart-container {
    min-width: 1000px;
    height: 550px !important;
  }
  .stats-chart {
    min-width: 1000px;
    height: 450px !important;
  }
}
@media (max-width: 1199px) {
  .card-chart-container {
    min-width: 700px;
    height: 450px !important;
  }
  .stats-chart {
    min-width: 700px;
    height: 350px !important;
  }
}
</style>