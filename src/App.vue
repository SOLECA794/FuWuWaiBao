<template>
  <div class="app-container">
    <!-- 顶部导航 -->
    <header class="header">
      <div class="header-left">
        <img src="https://picsum.photos/40/40" alt="logo" class="logo">
        <span class="system-name">智能互动教学平台 · 学生端</span>
      </div>
      <div class="header-right">
        <el-badge :value="3" type="danger">
          <el-avatar :size="36" src="https://picsum.photos/100/100"></el-avatar>
        </el-badge>
        <span class="user-name">学生 2025001</span>
      </div>
    </header>

    <!-- 主体 -->
    <main class="main">
      <!-- 左侧：课件区 -->
      <section class="left-section">
        <div class="course-card">
          <div class="course-header">
            <h3>Python 数据分析 · 第 3 课：缺失值处理</h3>
            <div class="course-tag">
              <el-tag type="success" size="small">学习中</el-tag>
              <span class="page-count">第 {{ currentPage }} / {{ totalPage }} 页</span>
            </div>
          </div>

          <!-- 新增：断点续播弹窗 -->
          <el-dialog
            title="学习进度提醒"
            v-model="showBreakpointDialog"
            width="30%"
            center
            :close-on-click-modal="false"
          >
            <p>你上次学习停留到 <span style="color: #1989fa; font-weight: bold;">第 {{ breakpointPage }} 页</span>，是否继续从该页学习？</p>
            <template v-slot:footer>
              <el-button @click="restartStudy">重新开始</el-button>
              <el-button type="primary" @click="continueStudy">继续学习</el-button>
            </template>
          </el-dialog>

          <!-- 课件内容 -->
          <div class="course-content">
            <img :src="courseImg" alt="课件" class="course-img">
            <div v-if="tracePoint" class="trace-highlight"
                 :style="{top: traceTop+'px', left: traceLeft+'px'}">
            </div>
          </div>

          <!-- 翻页控制 -->
          <div class="course-control">
            <el-button icon="el-icon-arrow-left" @click="prevPage">上一页</el-button>
            <el-button :icon="isPlay?'el-icon-video-pause':'el-icon-video-play'" @click="togglePlay">
              {{ isPlay ? '暂停' : '自动播放' }}
            </el-button>
            <el-button icon="el-icon-arrow-right" @click="nextPage">下一页</el-button>
          </div>
        </div>
      </section>

      <!-- 右侧：智能互动区 -->
      <section class="right-section">
        <el-tabs v-model="activeTab" class="smart-tab">
          <el-tab-pane label="多模态提问" name="ask">
            <div class="panel-box">
              <div class="question-header">
                <span>当前溯源定位：</span>
                <el-tag type="primary" size="small">第 {{ currentPage }} 页</el-tag>
                <el-tag type="warning" size="small" v-if="tracePoint">已圈选知识点</el-tag>
              </div>
              <div class="multi-modal-input">
                <el-input v-model="question" type="textarea" :rows="2" placeholder="输入问题..."></el-input>
                <div class="modal-tools">
                  <el-button size="mini" icon="el-icon-picture" @click="openUpload">截图/圈图提问</el-button>
                  <el-button size="mini" icon="el-icon-microphone" type="info">语音提问</el-button>
                  <el-button size="mini" type="success" @click="sendMultiModalQuestion" :disabled="!question">
                    发送 AI 问答
                  </el-button>
                </div>
              </div>
              <div class="ai-chat" v-if="aiReply">
                <div class="chat-item teacher">
                  <div class="title">AI 智能答疑</div>
                  <div class="content">{{ aiReply }}</div>
                </div>
              </div>
            </div>
          </el-tab-pane>
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
                <div class="title">AI 诊断薄弱点</div>
                <el-tag type="danger" size="small">缺失值填充</el-tag>
                <el-tag type="danger" size="small">异常值识别</el-tag>
                <el-tag type="danger" size="small">重复值处理</el-tag>
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane label="溯源定位" name="trace">
            <div class="panel-box">
              <p>点击课件任意位置 → 圈选知识点 → 精准提问</p>
              <el-button type="primary" plain size="small" @click="openTraceMode">
                开启溯源定位
              </el-button>
              <div v-if="traceLog" class="trace-log">
                {{ traceLog }}
              </div>
            </div>
          </el-tab-pane>
          <!-- 新增：知识点拆解 Tab -->
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
                    拖拽文件到此处，或<em>点击上传</em><br>
                    <span style="font-size: 12px; color: #999;">支持 PDF / PPTX 格式</span>
                  </div>
                </el-upload>
                <el-button type="primary" @click="parseKnowledge" :disabled="!uploadedFile" style="margin-top: 10px;">
                  开始拆解知识点
                </el-button>
              </div>

              <!-- 知识点树形展示 -->
              <div class="knowledge-tree" v-if="knowledgeList.length > 0">
                <h4 style="margin: 10px 0;">知识点结构（点击可定位）</h4>
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
                style="margin-top: 20px;"
              ></el-alert>
              <el-alert
                v-if="parseResult"
                :title="parseResult"
                type="success"
                show-icon
                closable
                style="margin-top: 20px;"
              ></el-alert>
            </div>
          </el-tab-pane>
        </el-tabs>
      </section>
    </main>

    <!-- 底部 -->
    <footer class="footer">
      2026 服务外包创新创业大赛 · 智能教学系统 · 前端企业级设计
    </footer>
  </div>
</template>

<script setup>
// 1. 导入所有需要的依赖
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

// 2. 原有变量（完全保留）
const currentPage = ref(3)
const totalPage = ref(10)
const isPlay = ref(false)
const courseImg = ref('https://picsum.photos/900/500?random=3')
const activeTab = ref('ask')
const question = ref('')
const aiReply = ref('')
const tracePoint = ref(false)
const traceTop = ref(0)
const traceLeft = ref(0)
const traceLog = ref('')

// 3. 断点续播变量（新增）
const showBreakpointDialog = ref(false)
const breakpointPage = ref(3) // 模拟断点页码

// 新增：知识点拆解核心变量
const uploadedFile = ref(null) // 上传的文件对象
const isParsing = ref(false) // 拆解中状态
const parseResult = ref('') // 拆解结果提示
const knowledgeList = ref([]) // 知识点树形数据
const treeProps = ref({ // 树形组件配置
  label: 'name',
  children: 'children'
})

// 4. 原有方法（完全保留）
const prevPage = () => {
  if (currentPage.value <= 1) return
  currentPage.value--
  courseImg.value = `https://picsum.photos/900/500?random=${currentPage.value}`
}
const nextPage = () => {
  if (currentPage.value >= totalPage.value) return
  currentPage.value++
  courseImg.value = `https://picsum.photos/900/500?random=${currentPage.value}`
}
const togglePlay = () => {
  isPlay.value = !isPlay.value
}
const openUpload = () => {
  ElMessage.info('已打开截图/圈图提问')
}
const sendMultiModalQuestion = () => {
  aiReply.value = `【AI 答疑】你在第 ${currentPage.value} 页提问：${question.value}\n
→ 溯源定位：缺失值填充方法\n
→ 回答：fillna() 适合连续数据，dropna() 适合少量缺失，interpolate() 用于时序数据。`
  ElMessage.success('AI 答疑完成')
}
const openTraceMode = () => {
  tracePoint.value = true
  traceTop.value = 150
  traceLeft.value = 200
  traceLog.value = `已定位：第 ${currentPage.value} 页 → 缺失值处理区域`
}

// 5. 断点续播方法（新增）
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

// 新增：知识点拆解核心方法
// 1. 监听文件上传
const handleFileChange = (file) => {
  uploadedFile.value = file.raw
  parseResult.value = ''
  knowledgeList.value = []
}

// 2. 知识点拆解（模拟数据版）
const parseKnowledge = async () => {
  if (!uploadedFile.value) {
    ElMessage.warning('请先上传 PDF/PPTX 文件！')
    return
  }

  isParsing.value = true
  try {
    // 模拟AI拆解的知识点数据
    const mockKnowledge = [
      { name: 'Python 缺失值处理', children: [
        { name: '缺失值检测', children: [{ name: 'isnull() 方法' }, { name: 'info() 方法' }] },
        { name: '缺失值填充', children: [{ name: 'fillna() 均值填充' }, { name: 'interpolate() 插值填充' }] }
      ]},
      { name: 'Python 异常值处理', children: [
        { name: '异常值识别', children: [{ name: '箱线图法' }, { name: 'Z分数法' }] },
        { name: '异常值处理', children: [{ name: '删除法' }, { name: '替换法' }] }
      ]}
    ]

    // 格式化数据
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

// 3. 辅助：统计知识点节点数量
const countNodes = (tree) => {
  let count = 0
  tree.forEach(node => {
    count++
    if (node.children && node.children.length) {
      count += countNodes(node.children)
    }
  })
  return count
}

// 4. 点击知识点：联动溯源定位
const handleNodeClick = (data) => {
  ElMessage.info(`已定位到知识点：${data.name}`)
  tracePoint.value = true
  traceTop.value = 200
  traceLeft.value = 300
  traceLog.value = `已定位知识点：${data.name}`
}
</script>

<!-- ========== 这里就是你要找的 <style scoped> 部分 ========== -->
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
  background: #1989fa;
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
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
  box-shadow: 0 2px 12px rgba(0,0,0,0.06);
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
  background: rgba(255,102,51,0.1);
  pointer-events: none;
  border-radius: 6px;
  animation: flash 1.2s infinite;
}
@keyframes flash {
  0% { opacity: 0.4; }
  50% { opacity: 0.8; }
  100% { opacity: 0.4; }
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
  box-shadow: 0 2px 10px rgba(0,0,0,0.05);
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
  border-left: 4px solid #1989fa;
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

/* ========== 新增：知识点拆解样式 ========== */
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
</style>
