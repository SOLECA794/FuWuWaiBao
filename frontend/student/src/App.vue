<template>
  <div class="app-container">
    <!-- 顶部导航栏 -->
    <header class="header">
      <div class="header-left">
        <img src="https://picsum.photos/40/40?random=1" alt="logo" class="logo" />
        <div class="system-name">智能学习课堂系统</div>
      </div>
      <div class="header-right">
        <span>用户名：学生001</span>
        <span>当前课程：Python数据处理</span>
      </div>
    </header>

    <!-- 主体内容 -->
    <main class="main">
      <!-- 左侧课件区域 -->
      <section class="left-section">
        <div class="course-card">
          <div class="course-header">
            <h3>Python数据处理 - 第{{ currentPage }}页</h3>
            <div class="course-tag">
              <el-tag size="small">数据清洗</el-tag>
              <span class="page-count">{{ currentPage }}/{{ totalPage }}</span>
            </div>
          </div>

          <div class="course-content">
            <img :src="courseImg" alt="课件内容" class="course-img" />
            <!-- 溯源定位高亮框 -->
            <div
              v-if="tracePoint"
              class="trace-highlight"
              :style="{ top: traceTop + 'px', left: traceLeft + 'px' }"
            ></div>
          </div>

          <div class="course-control">
            <el-button @click="prevPage" icon="el-icon-arrow-left" size="small">上一页</el-button>
            <el-button @click="togglePlay" :icon="isPlay ? 'el-icon-pause' : 'el-icon-play'" size="small">
              {{ isPlay ? '暂停' : '播放' }}
            </el-button>
            <el-button @click="nextPage" icon="el-icon-arrow-right" size="small">下一页</el-button>
          </div>
        </div>
      </section>

      <!-- 右侧智能交互区域 -->
      <section class="right-section">
        <el-tabs v-model="activeTab" class="smart-tab">
          <!-- 多模态提问 Tab -->
          <el-tab-pane label="多模态提问" name="ask">
            <div class="panel-box">
              <div class="question-header">基于课件内容精准提问</div>
              <div class="multi-modal-input">
                <el-input
                  v-model="question"
                  type="textarea"
                  placeholder="请输入你的问题..."
                  :rows="4"
                ></el-input>
                <div class="modal-tools">
                  <el-button size="small" @click="openUpload" icon="el-icon-upload">上传截图</el-button>
                  <el-button size="small" icon="el-icon-crop">圈图提问</el-button>
                  <el-button type="primary" size="small" @click="sendMultiModalQuestion" icon="el-icon-send">
                    发送提问
                  </el-button>
                </div>
              </div>

              <!-- AI 回答区域 -->
              <div class="ai-chat" v-if="aiReply">
                <div class="chat-item teacher">
                  <div class="title">AI 助教回答</div>
                  <div>{{ aiReply }}</div>
                </div>
              </div>
            </div>
          </el-tab-pane>

          <!-- 学习数据 Tab -->
          <el-tab-pane label="学习数据" name="data">
            <div class="panel-box">
              <div class="data-grid">
                <div class="data-item">
                  <div class="num">87%</div>
                  <div class="label">专注度</div>
                </div>
                <div class="data-item">
                  <div class="num">12 次</div>
                  <div class="label">本课提问</div>
                </div>
                <div class="data-item">
                  <div class="num">3 个</div>
                  <div class="label">薄弱点</div>
                </div>
                <div class="data-item">
                  <div class="num">92%</div>
                  <div class="label">掌握率</div>
                </div>
              </div>

              <div class="weak-point">
                <div class="title">AI 诊断薄弱点（点击可学习）</div>
                <el-tag
                  type="danger"
                  size="small"
                  style="cursor: pointer; margin: 3px"
                  @click="startWeakPointLearn('缺失值填充')"
                >
                  缺失值填充
                </el-tag>
                <el-tag
                  type="danger"
                  size="small"
                  style="cursor: pointer; margin: 3px"
                  @click="startWeakPointLearn('异常值识别')"
                >
                  异常值识别
                </el-tag>
                <el-tag
                  type="danger"
                  size="small"
                  style="cursor: pointer; margin: 3px"
                  @click="startWeakPointLearn('重复值处理')"
                >
                  重复值处理
                </el-tag>
              </div>

              <!-- 薄弱点讲解 -->
              <div v-if="currentExplain" class="explain-card">
                <h4>📘 {{ currentWeakPoint }} · 知识点讲解</h4>
                <p>{{ currentExplain }}</p>
                <el-button type="primary" @click="generateTest" style="margin-top: 10px">
                  已学会，开始习题检测
                </el-button>
              </div>

              <!-- 习题检测 -->
              <div v-if="currentTest" class="test-card">
                <h4>📝 练习题：{{ currentWeakPoint }}</h4>
                <p>{{ currentTest.question }}</p>
                <el-button
                  v-for="(opt, idx) in currentTest.options"
                  :key="idx"
                  style="margin: 5px"
                  @click="checkAnswer(opt)"
                >
                  {{ opt }}
                </el-button>
                <div v-if="testResult" :style="{ color: testResult.correct ? 'green' : 'red', marginTop: '10px' }">
                  {{ testResult.msg }}
                </div>
                <div v-if="testResult && testResult.analysis">
                  <small style="display: block; margin-top: 5px">解析：{{ testResult.analysis }}</small>
                </div>
              </div>
            </div>
          </el-tab-pane>

          <!-- 溯源定位 Tab -->
          <el-tab-pane label="溯源定位" name="trace">
            <div class="panel-box">
              <p>点击课件任意位置 → 圈选知识点 → 精准提问</p>
              <el-button type="primary" plain size="small" @click="openTraceMode">开启溯源定位</el-button>
              <div v-if="traceLog" class="trace-log">{{ traceLog }}</div>
            </div>
          </el-tab-pane>

          <!-- 知识点拆解 Tab -->
          <el-tab-pane label="知识点拆解" name="parse">
            <div class="panel-box">
              <!-- 文件上传区域 -->
              <div class="upload-area">
                <el-upload
                  class="upload-demo"
                  drag
                  action="#"
                  :auto-upload="false"
                  :on-change="handleFileChange"
                  accept=".pdf,.pptx"
                  :limit="1"
                >
                  <i class="el-icon-upload"></i>
                  <div class="el-upload__text">
                    拖拽文件到此处，或<em>点击上传</em><br />
                    <span style="font-size: 12px; color: #999">支持 PDF / PPTX 格式</span>
                  </div>
                </el-upload>
                <el-button type="primary" @click="parseKnowledge" :disabled="!uploadedFile" style="margin-top: 10px">
                  开始拆解知识点
                </el-button>
              </div>

              <!-- 知识点树形展示 -->
              <div class="knowledge-tree" v-if="knowledgeList.length > 0">
                <h4 style="margin: 10px 0">知识点结构（点击可定位）</h4>
                <el-tree
                  :data="knowledgeList"
                  :props="treeProps"
                  node-key="id"
                  @node-click="handleNodeClick"
                  default-expand-all
                ></el-tree>
              </div>

              <!-- 状态提示 -->
              <el-alert
                v-if="isParsing"
                title="正在拆解知识点，请稍候..."
                type="info"
                show-icon
                style="margin-top: 20px"
              ></el-alert>
              <el-alert
                v-if="parseResult"
                :title="parseResult"
                type="success"
                show-icon
                closable
                style="margin-top: 20px"
              ></el-alert>
            </div>
          </el-tab-pane>
        </el-tabs>
      </section>
    </main>

    <!-- 底部版权 -->
    <footer class="footer">© 2025 智能学习课堂系统 - 版权所有</footer>

    <!-- 断点续播弹窗 -->
    <el-dialog
      v-model="showBreakpointDialog"
      title="断点续播"
      width="30%"
      :close-on-click-modal="false"
      :show-close="false"
    >
      <p>检测到你上次学习到第 {{ breakpointPage }} 页，是否继续学习？</p>
      <template #footer>
        <el-button @click="restartStudy">重新开始</el-button>
        <el-button type="primary" @click="continueStudy">继续学习</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
// 1. 导入所有需要的依赖
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'

// 2. 基础课件变量
const currentDocId = ref(localStorage.getItem('last_doc_id') || '')
const currentPage = ref(1)
const totalPage = ref(1)
const isPlay = ref(false)
const courseImg = ref('https://picsum.photos/900/500?random=1')
const activeTab = ref('ask')
const currentScript = ref('')
const scriptLoading = ref(false)
const lessonsCache = ref({})

// 3. 多模态提问变量
const question = ref('')
const aiReply = ref('')

// 4. 溯源定位变量
const tracePoint = ref(false)
const traceTop = ref(0)
const traceLeft = ref(0)
const traceLog = ref('')

// 5. 断点续播变量
const showBreakpointDialog = ref(false)
const breakpointPage = ref(3)

// 6. 知识点拆解变量
const uploadedFile = ref(null)
const isParsing = ref(false)
const parseResult = ref('')
const knowledgeList = ref([])
const treeProps = ref({
  label: 'name',
  children: 'children'
})

// 7. 薄弱点学习 + 习题检测变量
const currentWeakPoint = ref('')
const currentExplain = ref('')
const currentTest = ref(null)
const testResult = ref(null)

// 8. 基础课件方法
const prevPage = () => {
  if (currentPage.value <= 1) return
  currentPage.value--
  updateCourseContent()
}

const nextPage = () => {
  if (currentPage.value >= totalPage.value) return
  currentPage.value++
  updateCourseContent()
}

const updateCourseContent = () => {
  courseImg.value = `https://picsum.photos/900/500?random=${currentPage.value}`
  fetchLessonContent()
}
  currentPage.value++
  courseImg.value = `https://picsum.photos/900/500?random=${currentPage.value}`
}

const togglePlay = () => {
  isPlay.value = !isPlay.value
}

// 9. 多模态提问方法
const openUpload = () => {
  ElMessage.info('已打开截图/圈图提问')
}

const sendMultiModalQuestion = () => {
  aiReply.value = `【AI 答疑】你在第 ${currentPage.value} 页提问：${question.value}\n
→ 溯源定位：缺失值填充方法\n
→ 回答：fillna() 适合连续数据，dropna() 适合少量缺失，interpolate() 用于时序数据。`
  ElMessage.success('AI 答疑完成')
}

// 10. 溯源定位方法
const openTraceMode = () => {
  tracePoint.value = true
  traceTop.value = 150
  traceLeft.value = 200
  traceLog.value = `已定位：第 ${currentPage.value} 页 → 缺失值处理区域`
}

// 11. 断点续播方法
onMounted(() => {
  showBreakpointDialog.value = true
})

const continueStudy = () => {
  currentPage.value = breakpointPage.value
  courseImg.value = `https://picsum.photos/900/500?random=${breakpointPage.value}`
  showBreakpointDialog.value = false
  ElMessage.success(`已为你跳转到第 ${breakpointPage.value} 页`)
}

const restartStudy = () => {
  currentPage.value = 1
  courseImg.value = `https://picsum.photos/900/500?random=1`
  showBreakpointDialog.value = false
  ElMessage.info('已回到第1页重新开始学习')
}

// 12. 知识点拆解方法
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
    // 模拟AI拆解的知识点数据
    const mockKnowledge = [
      {
        name: 'Python 缺失值处理',
        children: [
          { name: '缺失值检测', children: [{ name: 'isnull() 方法' }, { name: 'info() 方法' }] },
          { name: '缺失值填充', children: [{ name: 'fillna() 均值填充' }, { name: 'interpolate() 插值填充' }] }
        ]
      },
      {
        name: 'Python 异常值处理',
        children: [
          { name: '异常值识别', children: [{ name: '箱线图法' }, { name: 'Z分数法' }] },
          { name: '异常值处理', children: [{ name: '删除法' }, { name: '替换法' }] }
        ]
      }
    ]

    // 格式化数据适配树形组件
    knowledgeList.value = mockKnowledge.map((item, index) => ({
      id: index + 1,
      name: item.name,
      children: item.children || []
    }))

    // 统计知识点数量
    parseResult.value = `拆解成功！共识别出 ${countNodes(knowledgeList.value)} 个知识点`
    ElMessage.success('知识点结构拆解完成！')
  } catch (error) {
    parseResult.value = '拆解失败，请重试！'
    ElMessage.error('知识点拆解失败')
  } finally {
    isParsing.value = false
  }
}

// 辅助：统计知识点节点数量
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

// 点击知识点联动溯源定位
const handleNodeClick = (data) => {
  ElMessage.info(`已定位到知识点：${data.name}`)
  tracePoint.value = true
  traceTop.value = 200
  traceLeft.value = 300
  traceLog.value = `已定位知识点：${data.name}`
}

// 13. 薄弱点学习 + 习题检测核心逻辑
// 知识点讲解库
const weakPointExplain = {
  '缺失值填充': `缺失值是数据中为空的部分，常用方法：
1. fillna() 填充常数、均值、中位数
2. interpolate() 线性插值（适合时序）
3. dropna() 直接删除行/列`,

  '异常值识别': `异常值是明显偏离正常范围的数据，常用方法：
1. 箱线图法：超过 Q3+1.5IQR 或低于 Q1-1.5IQR 判为异常
2. Z分数法：绝对值大于3视为异常
3. 直方图/散点图观察偏离点`,

  '重复值处理': `重复值是完全相同的行，处理步骤：
1. duplicated() 查找重复
2. drop_duplicates() 删除重复
3. 按关键字段去重（如 id、时间）`
}

// 题库
const testBank = {
  '缺失值填充': {
    question: '以下哪种方法适合时序数据的缺失值填充？',
    options: ['fillna(均值)', 'dropna()', 'interpolate()', '直接填0'],
    answer: 'interpolate()',
    analysis: 'interpolate 是线性插值，最适合时间顺序数据'
  },
  '异常值识别': {
    question: '箱线图中，超过哪个范围被认为是异常值？',
    options: ['±2σ', 'Q1-1.5IQR ～ Q3+1.5IQR', '平均值±标准差', '95%置信区间'],
    answer: 'Q1-1.5IQR ～ Q3+1.5IQR',
    analysis: '箱线图异常判定标准就是 1.5 倍四分位距'
  },
  '重复值处理': {
    question: 'pandas中删除重复值用哪个方法？',
    options: ['unique()', 'drop_duplicates()', 'delete()', 'remove()'],
    answer: 'drop_duplicates()',
    analysis: 'drop_duplicates() 是官方去重方法'
  }
}

// 开始学习薄弱点
const startWeakPointLearn = (point) => {
  currentWeakPoint.value = point
  currentExplain.value = weakPointExplain[point]
  currentTest.value = null
  testResult.value = null
}

// 生成测试题
const generateTest = () => {
  currentTest.value = testBank[currentWeakPoint.value]
  testResult.value = null
}

// 检查答案
const checkAnswer = (option) => {
  const isCorrect = option === currentTest.value.answer
  testResult.value = {
    correct: isCorrect,
    msg: isCorrect ? '✅ 回答正确！' : '❌ 回答错误',
    analysis: currentTest.value.analysis
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
  background: #f5f7fa;
  display: flex;
  flex-direction: column;
  font-family: "Microsoft YaHei", sans-serif;
}
.header {
  height: 60px;
  background: #c4d5e6;
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}
.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}
.logo {
  border-radius: 50%;
}
.system-name {
  font-size: 16px;
  font-weight: 500;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
}
.main {
  flex: 1;
  display: flex;
  gap: 20px;
  padding: 20px;
  overflow: hidden;
}
.left-section {
  flex: 7;
}
.right-section {
  flex: 3;
}
.course-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  height: 100%;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  display: flex;
  flex-direction: column;
}
.course-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.course-header h3 {
  font-size: 16px;
  color: #333;
}
.course-tag {
  display: flex;
  gap: 8px;
  align-items: center;
}
.page-count {
  font-size: 13px;
  color: #666;
}
.course-content {
  flex: 1;
  background: #f9f9f9;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}
.course-img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}
.trace-highlight {
  position: absolute;
  width: 180px;
  height: 100px;
  border: 3px solid #ff6633;
  background: rgba(255, 102, 51, 0.1);
  pointer-events: none;
  border-radius: 6px;
  animation: flash 1.2s infinite;
}
@keyframes flash {
  0% {
    opacity: 0.4;
  }
  50% {
    opacity: 0.8;
  }
  100% {
    opacity: 0.4;
  }
}
.course-control {
  margin-top: 16px;
  display: flex;
  justify-content: center;
  gap: 12px;
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
.panel-box {
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}
.question-header {
  margin-bottom: 10px;
  font-size: 14px;
}
.multi-modal-input {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.modal-tools {
  display: flex;
  gap: 8px;
  align-items: center;
}
.ai-chat {
  margin-top: 16px;
}
.chat-item {
  padding: 12px;
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.5;
}
.teacher {
  background: #e6f7ff;
  border-left: 4px solid #97c2ed;
}
.chat-item .title {
  font-weight: bold;
  margin-bottom: 4px;
  color: #1989fa;
}
.data-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 16px;
}
.data-item {
  background: #f7f8fa;
  padding: 16px;
  border-radius: 8px;
  text-align: center;
}
.data-item .num {
  font-size: 20px;
  font-weight: bold;
  color: #1989fa;
}
.data-item .label {
  font-size: 13px;
  color: #666;
  margin-top: 4px;
}
.weak-point {
  margin-top: 10px;
}
.weak-point .title {
  font-size: 14px;
  margin-bottom: 8px;
  font-weight: 500;
}
.trace-log {
  margin-top: 12px;
  padding: 10px;
  background: #fff7e6;
  border-radius: 6px;
  font-size: 13px;
  color: #d48806;
}
.footer {
  height: 40px;
  background: white;
  text-align: center;
  line-height: 40px;
  font-size: 12px;
  color: #999;
  border-top: 1px solid #eee;
}

/* 知识点拆解样式 */
.upload-area {
  padding: 20px;
  border: 1px dashed #dcdfe6;
  border-radius: 8px;
  margin-bottom: 20px;
  text-align: center;
}
.knowledge-tree {
  margin-top: 20px;
  max-height: 400px;
  overflow-y: auto;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 10px;
}
:deep(.el-tree) {
  --el-tree-node-content-hover-bg-color: #e6f7ff;
}

/* 薄弱点学习 + 习题检测样式 */
.explain-card {
  margin-top: 16px;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
  border-left: 4px solid #1989fa;
}
.test-card {
  margin-top: 16px;
  padding: 12px;
  background: #fff7e6;
  border-radius: 8px;
  border-left: 4px solid #faad14;
}
</style>