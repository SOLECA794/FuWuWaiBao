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

          <!-- 课件内容 -->
          <div class="course-content">
            <img :src="courseImg" alt="课件" class="course-img">

            <!-- 溯源定位高亮（创新点1）-->
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
          <!-- 1. 多模态提问面板（创新点2）-->
          <el-tab-pane label="多模态提问" name="ask">
            <div class="panel-box">
              <div class="question-header">
                <span>当前溯源定位：</span>
                <el-tag type="primary" size="small">第 {{ currentPage }} 页</el-tag>
                <el-tag type="warning" size="small" v-if="tracePoint">已圈选知识点</el-tag>
              </div>

              <!-- 多模态输入：文本 + 图片 + 语音 -->
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

              <!-- AI 对话记录 -->
              <div class="ai-chat" v-if="aiReply">
                <div class="chat-item teacher">
                  <div class="title">AI 智能答疑</div>
                  <div class="content">{{ aiReply }}</div>
                </div>
              </div>
            </div>
          </el-tab-pane>

          <!-- 2. 学习数据统计面板（创新点3）-->
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

          <!-- 3. 溯源定位过程（创新点4）-->
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
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

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

// 翻页
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

// 多模态提问
const openUpload = () => {
  ElMessage.info('已打开截图/圈图提问')
}
const sendMultiModalQuestion = () => {
  aiReply.value = `【AI 答疑】你在第 ${currentPage.value} 页提问：${question.value}\n
→ 溯源定位：缺失值填充方法\n
→ 回答：fillna() 适合连续数据，dropna() 适合少量缺失，interpolate() 用于时序数据。`
  ElMessage.success('AI 答疑完成')
}

// 溯源定位
const openTraceMode = () => {
  tracePoint.value = true
  traceTop.value = 150
  traceLeft.value = 200
  traceLog.value = `已定位：第 ${currentPage.value} 页 → 缺失值处理区域`
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

/* 头部 */
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

/* 主体 */
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

/* 课件卡片 */
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

/* 课件内容 */
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

/* 溯源高亮 */
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

/* 右侧面板 */
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

/* 多模态提问 */
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

/* AI 聊天 */
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

/* 数据统计 */
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

/* 溯源日志 */
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
</style>
