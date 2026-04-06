<template>
  <div class="panel-box">
    <div class="panel-hero">
      <div class="panel-head">
        <div>
          <div class="eyebrow">问答工作区</div>
          <div class="question-header">自然对话模式</div>
          <div class="header-subtitle">保留连续对话体验，只在需要时查看统计。</div>
        </div>
        <div class="head-actions">
          <button class="stats-toggle-btn" @click="showStats = !showStats">
            <span class="dot"></span>
            {{ showStats ? '收起统计' : '统计概览' }}
          </button>
          <div class="session-chip" v-if="latestAnswerMeta.sessionId">
            会话 {{ latestAnswerMeta.sessionId.slice(0, 8) }}
          </div>
        </div>
      </div>

      <div class="stats-flyout" v-if="showStats">
        <div class="dashboard-cards compact">
          <div class="dashboard-card" v-for="item in dashboardStats" :key="item.label" :class="item.tone">
            <div class="card-label">{{ item.label }}</div>
            <div class="card-value">{{ item.value }}</div>
            <div class="card-note">{{ item.note }}</div>
            <div class="card-track">
              <span :style="{ width: `${item.progress}%` }"></span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="conversation-board">
      <div class="conversation-header">
        <div>
          <div class="conversation-title">流式对话</div>
          <div class="conversation-subtitle">上方显示状态信息，下方保持简洁输入，按基础聊天界面组织。</div>
        </div>
        <div class="typing-status" :class="{ active: askLoading || streamTypingActive }">
          <span class="typing-dot"></span>
          {{ askLoading ? 'AI 正在思考' : streamTypingActive ? '逐字输出中' : '对话已就绪' }}
        </div>
      </div>

      <div class="conversation-thread">
        <div class="empty-state" v-if="!recentQuestionText && !aiReply && !askLoading && !streamTypingActive">
          在下方输入问题开始对话，例如：这一步为什么这样推导？
        </div>

        <div class="bubble-row user" v-if="recentQuestionText">
          <div class="bubble-avatar user-avatar">我</div>
          <div class="bubble user-bubble">
            <div class="bubble-meta">学生提问</div>
            <div class="bubble-text">{{ recentQuestionText }}</div>
          </div>
        </div>

        <div class="bubble-row assistant">
          <div class="bubble-avatar assistant-avatar">AI</div>
          <div class="bubble assistant-bubble" :class="{ streaming: askLoading || streamTypingActive }">
            <div class="bubble-meta-row">
              <span class="bubble-meta">AI 助教</span>
              <span class="reply-state" :class="{ reteach: latestAnswerMeta.needReteach }">
                {{ askLoading ? '生成中' : latestAnswerMeta.needReteach ? '重讲模式' : '标准答疑' }}
              </span>
            </div>
            <div class="bubble-text">{{ aiReply || '正在结合当前课件和上文追问生成回答...' }}</div>
            <div class="typing-hint" v-if="streamTypingActive || askLoading">
              <span></span><span></span><span></span>
            </div>
            <div class="reply-meta" v-if="latestAnswerMeta.sourcePage">
              <span>来源页：第 {{ latestAnswerMeta.sourcePage }} 页</span>
              <span v-if="latestAnswerMeta.sourceNodeId">来源节点：{{ latestAnswerMeta.sourceNodeId }}</span>
              <span>{{ latestAnswerMeta.needReteach ? '系统判断你需要更通俗的解释' : '系统将按当前节奏继续讲解' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="chat-composer">
      <div class="composer-toolbar">
        <el-select
          :model-value="summaryMode"
          size="small"
          class="compact-mode-select"
          @update:model-value="$emit('update:summaryMode', $event)"
        >
          <el-option label="速览" value="quick" />
          <el-option label="考试" value="exam" />
          <el-option label="讲解" value="teach" />
        </el-select>
        <el-button size="small" text @click="$emit('clear-draft')">清空</el-button>
      </div>
      <div class="composer-main">
        <el-input
          :model-value="question"
          type="textarea"
          class="chat-input"
          placeholder="输入你的问题..."
          :autosize="{ minRows: 1, maxRows: 3 }"
          @update:model-value="$emit('update:question', $event)"
        ></el-input>
        <el-button type="primary" class="send-btn" :loading="askLoading" @click="$emit('send-question')">发送</el-button>
      </div>
    </div>

  </div>
</template>

<script setup>
/* eslint-disable no-undef */
import { computed, ref } from 'vue'

const showStats = ref(false)

const clampPercent = (value) => Math.max(0, Math.min(100, Number(value) || 0))

const props = defineProps({
  question: {
    type: String,
    default: ''
  },
  askLoading: {
    type: Boolean,
    default: false
  },
  aiReply: {
    type: String,
    default: ''
  },
  streamTypingActive: {
    type: Boolean,
    default: false
  },
  qaHistory: {
    type: Array,
    default: () => []
  },
  latestAnswerMeta: {
    type: Object,
    default: () => ({})
  },
  summaryMode: {
    type: String,
    default: 'quick'
  },
  mergedSummary: {
    type: String,
    default: ''
  }
})

defineEmits([
  'update:question',
  'open-upload',
  'send-question',
  'update:summaryMode',
  'generate-summary',
  'use-summary',
  'clear-draft'
])

const questionHeat = computed(() => {
  const text = String(props.question || '').trim()
  return clampPercent(Math.round((text.length / 140) * 100))
})

const summaryCoverage = computed(() => {
  const text = String(props.mergedSummary || '').trim()
  return clampPercent(Math.round((text.length / 320) * 100))
})

const historyIntensity = computed(() => clampPercent(Math.round((props.qaHistory.length / 6) * 100)))

const answerIntensity = computed(() => {
  const base = props.askLoading ? 42 : 18
  return clampPercent(base + Math.round(String(props.aiReply || '').trim().length / 12))
})

const dashboardStats = computed(() => [
  {
    label: '提问热度',
    value: `${questionHeat.value}%`,
    note: props.question ? '当前草稿已进入提问状态' : '输入问题后热度会立即提升',
    tone: 'primary',
    progress: questionHeat.value
  },
  {
    label: '摘要覆盖',
    value: `${summaryCoverage.value}%`,
    note: props.mergedSummary ? '摘要已准备好用于追问' : '先生成摘要可提升回答质量',
    tone: 'success',
    progress: summaryCoverage.value
  },
  {
    label: '历史轨迹',
    value: `${props.qaHistory.length}`,
    note: props.qaHistory.length ? '最近提问已形成学习轨迹' : '提问记录会在这里沉淀',
    tone: 'accent',
    progress: historyIntensity.value
  },
  {
    label: '回答活跃',
    value: props.askLoading ? '生成中' : '就绪',
    note: latestAnswerStateText.value,
    tone: 'warm',
    progress: answerIntensity.value
  }
])

const latestAnswerStateText = computed(() => {
  if (props.askLoading) return 'AI 正在整理当前课件上下文'
  if (props.latestAnswerMeta.needReteach) return '当前回答会更偏向重讲与拆解'
  if (props.aiReply) return '回答已落入可继续学习的节奏'
  return '等待你发起下一轮提问'
})

const recentQuestionText = computed(() => {
  const draftQuestion = String(props.question || '').trim()
  if (draftQuestion) return draftQuestion
  const lastQuestion = props.qaHistory?.[0]?.question
  return String(lastQuestion || '').trim()
})

</script>

<style scoped>
.panel-box {
  display: flex;
  flex-direction: column;
  gap: 12px;
  position: relative;
  overflow: hidden;
  background:
    radial-gradient(circle at top right, rgba(143, 193, 181, 0.22), transparent 34%),
    radial-gradient(circle at 8% 18%, rgba(47, 96, 90, 0.08), transparent 22%),
    linear-gradient(180deg, #ffffff 0%, #f6fbf8 100%);
  border-radius: 22px;
  padding: 18px;
  min-height: 620px;
  border: 1px solid #d9e7df;
  box-shadow: 0 18px 34px rgba(33, 61, 54, 0.08);
}

.panel-box::before {
  content: '';
  position: absolute;
  inset: 0;
  background-image: linear-gradient(135deg, rgba(47, 96, 90, 0.06) 0%, rgba(47, 96, 90, 0.02) 40%, transparent 40%), linear-gradient(90deg, rgba(143, 193, 181, 0.12) 0 1px, transparent 1px 100%);
  background-size: 100% 100%, 48px 48px;
  opacity: 0.6;
  pointer-events: none;
}

.panel-box > * {
  position: relative;
  z-index: 1;
}

.panel-hero {
  padding: 2px 2px 6px;
}

.head-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stats-toggle-btn {
  border: 1px solid #d2e4dc;
  background: rgba(255, 255, 255, 0.92);
  color: #2f605a;
  border-radius: 999px;
  padding: 7px 12px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

.stats-toggle-btn:hover {
  transform: translateY(-1px);
  border-color: #8fc1b5;
  box-shadow: 0 8px 16px rgba(47, 96, 90, 0.12);
}

.stats-toggle-btn .dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #2f605a;
  box-shadow: 0 0 0 0 rgba(47, 96, 90, 0.35);
  animation: pulse-dot 1.4s ease-in-out infinite;
}
.panel-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}
.eyebrow {
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #6f867d;
}
.question-header {
  margin-top: 4px;
  font-size: 18px;
  font-weight: 700;
  color: #23463f;
}

.header-subtitle {
  margin-top: 8px;
  max-width: 38em;
  font-size: 13px;
  line-height: 1.65;
  color: #6c8278;
}

.stats-flyout {
  margin-top: 10px;
  border: 1px solid #d9e7df;
  border-radius: 16px;
  padding: 10px;
  background: linear-gradient(180deg, rgba(250, 252, 251, 0.95) 0%, rgba(244, 248, 246, 0.96) 100%);
  animation: bubble-in 0.24s ease both;
}

.dashboard-cards {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.dashboard-cards.compact .dashboard-card {
  padding: 10px 11px;
}

.dashboard-card {
  border-radius: 16px;
  padding: 12px 13px 11px;
  border: 1px solid #d9e7df;
  background: linear-gradient(180deg, rgba(255,255,255,0.96) 0%, rgba(244,248,246,0.98) 100%);
  box-shadow: 0 10px 20px rgba(33, 61, 54, 0.05);
}

.dashboard-card.primary {
  background: linear-gradient(180deg, rgba(235, 245, 243, 0.98) 0%, rgba(255,255,255,0.98) 100%);
}

.dashboard-card.success {
  background: linear-gradient(180deg, rgba(233, 246, 240, 0.98) 0%, rgba(255,255,255,0.98) 100%);
}

.dashboard-card.accent {
  background: linear-gradient(180deg, rgba(236, 244, 248, 0.98) 0%, rgba(255,255,255,0.98) 100%);
}

.dashboard-card.warm {
  background: linear-gradient(180deg, rgba(252, 247, 236, 0.98) 0%, rgba(255,255,255,0.98) 100%);
}

.card-label {
  font-size: 12px;
  color: #6f867d;
  font-weight: 600;
}

.card-value {
  margin-top: 6px;
  font-size: 22px;
  line-height: 1;
  font-weight: 800;
  color: #23463f;
}

.card-note {
  margin-top: 6px;
  min-height: 32px;
  font-size: 12px;
  line-height: 1.5;
  color: #557068;
}

.card-track {
  height: 6px;
  margin-top: 8px;
  border-radius: 999px;
  background: rgba(47, 96, 90, 0.08);
  overflow: hidden;
}

.card-track span {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #8fc1b5 0%, #2f605a 100%);
}

.summary-actions {
  margin-top: 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.summary-actions :deep(.el-button),
.chat-composer :deep(.el-button) {
  border-radius: 999px;
  transition: transform 0.2s ease, box-shadow 0.2s ease, background-color 0.2s ease, border-color 0.2s ease;
}

.summary-actions :deep(.el-button:hover),
.chat-composer :deep(.el-button:hover) {
  transform: translateY(-1px);
  box-shadow: 0 8px 16px rgba(47, 96, 90, 0.12);
}

.summary-actions :deep(.el-button:active),
.chat-composer :deep(.el-button:active) {
  transform: translateY(0);
  box-shadow: none;
}

.session-chip {
  flex-shrink: 0;
  padding: 6px 10px;
  border-radius: 999px;
  background: #eef5f1;
  color: #2f605a;
  font-size: 12px;
  font-weight: 600;
  border: 1px solid #d9e7df;
}
.chat-composer {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 10px;
  border-radius: 16px;
  border: 1px solid #d9e7df;
  background: rgba(255, 255, 255, 0.92);
}

.composer-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.composer-main {
  display: flex;
  align-items: flex-end;
  gap: 8px;
}

.chat-input {
  flex: 1;
}

.chat-input :deep(.el-textarea__inner) {
  border-radius: 16px;
  border-color: #d7e6de;
  background: rgba(255,255,255,0.95);
  box-shadow: inset 0 1px 2px rgba(15, 23, 42, 0.03);
  padding: 10px 12px;
  font-size: 14px;
  line-height: 1.5;
}

.chat-input :deep(.el-textarea__inner:focus) {
  border-color: #8fc1b5;
  box-shadow: 0 0 0 3px rgba(143, 193, 181, 0.14);
}

.send-btn {
  height: 36px;
  padding-inline: 16px;
}

.compact-mode-select {
  width: 92px;
}

.conversation-board {
  margin-top: 0;
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  border-radius: 20px;
  border: 1px solid #d9e7df;
  background:
    radial-gradient(circle at top right, rgba(143, 193, 181, 0.14), transparent 28%),
    linear-gradient(180deg, rgba(250, 252, 251, 0.98) 0%, rgba(244, 248, 246, 0.98) 100%);
  box-shadow: inset 0 1px 0 rgba(255,255,255,0.85), 0 12px 24px rgba(33, 61, 54, 0.05);
  padding: 14px;
}

.conversation-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 12px;
}

.conversation-title {
  font-size: 13px;
  font-weight: 700;
  color: #24453f;
}

.conversation-subtitle {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.6;
  color: #6f867d;
}

.typing-status {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  border-radius: 999px;
  background: #eef5f1;
  color: #5d7169;
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
}

.typing-status.active {
  background: rgba(143, 193, 181, 0.16);
  color: #2f605a;
}

.typing-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
  box-shadow: 0 0 0 0 rgba(47, 96, 90, 0.35);
}

.typing-status.active .typing-dot {
  animation: pulse-dot 1.2s ease-in-out infinite;
}

@keyframes pulse-dot {
  0%, 100% { transform: scale(1); opacity: 0.9; }
  50% { transform: scale(1.3); opacity: 0.5; }
}

.conversation-thread {
  display: flex;
  flex-direction: column;
  gap: 12px;
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding-right: 4px;
}

.empty-state {
  border: 1px dashed #cfe0d8;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.85);
  padding: 14px;
  color: #6f867d;
  font-size: 13px;
}

.bubble-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.bubble-row.user {
  justify-content: flex-end;
}

.bubble-row.user .bubble-avatar {
  order: 2;
}

.bubble-row.user .bubble {
  order: 1;
}

.bubble-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
  letter-spacing: 0.04em;
}

.user-avatar {
  background: linear-gradient(180deg, #2f605a 0%, #5c8f87 100%);
  color: #fff;
  box-shadow: 0 10px 18px rgba(47, 96, 90, 0.18);
}

.assistant-avatar {
  background: linear-gradient(180deg, #edf7f4 0%, #d9ece5 100%);
  color: #2f605a;
  border: 1px solid #cfe0d8;
}

.bubble {
  max-width: min(100%, 82%);
  border-radius: 18px;
  padding: 12px 14px;
  border: 1px solid #d9e7df;
  box-shadow: 0 10px 18px rgba(33, 61, 54, 0.05);
  position: relative;
  animation: bubble-in 0.34s ease both;
}

@keyframes bubble-in {
  from {
    opacity: 0;
    transform: translateY(8px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.user-bubble {
  background: linear-gradient(180deg, #2f605a 0%, #41766f 100%);
  color: #f8fffd;
  border-color: rgba(47, 96, 90, 0.25);
  box-shadow: 0 14px 24px rgba(47, 96, 90, 0.18);
}

.assistant-bubble {
  background: linear-gradient(180deg, #ffffff 0%, #f7faf8 100%);
}

.assistant-bubble.streaming {
  border-color: #8fc1b5;
  box-shadow: 0 14px 26px rgba(47, 96, 90, 0.12);
}

.bubble-meta-row {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  align-items: center;
  margin-bottom: 6px;
}

.bubble-meta {
  font-size: 11px;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: #739086;
  font-weight: 700;
}

.user-bubble .bubble-meta {
  color: rgba(255,255,255,0.85);
}

.bubble-text {
  white-space: pre-wrap;
  line-height: 1.8;
  font-size: 14px;
}

.user-bubble .bubble-text {
  color: #f8fffd;
}

.assistant-bubble .bubble-text {
  color: #425a51;
}

.chat-item .title {
  font-weight: 700;
  color: #24453f;
}
.reply-state {
  padding: 4px 10px;
  border-radius: 999px;
  background: #eaf4ef;
  color: #2f605a;
  font-size: 12px;
  font-weight: 600;
}
.reply-state.reteach {
  background: rgba(245, 158, 11, 0.14);
  color: #b45309;
}
.chat-content {
  color: #425a51;
  white-space: pre-wrap;
}

.typing-hint {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-top: 10px;
  font-size: 12px;
  color: #2f605a;
}

.typing-hint span {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #2f605a;
  animation: typing-wave 1s ease-in-out infinite;
}

.typing-hint span:nth-child(2) {
  animation-delay: 0.15s;
}

.typing-hint span:nth-child(3) {
  animation-delay: 0.3s;
}

@keyframes typing-wave {
  0%, 80%, 100% {
    transform: translateY(0);
    opacity: 0.35;
  }
  40% {
    transform: translateY(-4px);
    opacity: 1;
  }
}
.reply-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 10px;
  font-size: 12px;
  color: #6f867d;
}
.follow-up-box {
  margin-top: 14px;
  padding: 14px;
  border-radius: 16px;
  background: linear-gradient(180deg, #fff8eb 0%, #fffdf8 100%);
  border: 1px solid #f0dfb7;
}
.follow-up-title {
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: #92400e;
  text-transform: uppercase;
  margin-bottom: 6px;
}
.follow-up-text {
  font-size: 13px;
  color: #7c2d12;
  line-height: 1.7;
}
.question-history {
  margin-top: 14px;
  background: linear-gradient(180deg, rgba(246, 250, 248, 0.98) 0%, rgba(255,255,255,0.98) 100%);
  border: 1px solid #d9e7df;
  border-radius: 18px;
  padding: 12px;
}
.history-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}
.question-history .title {
  font-size: 13px;
  color: #23463f;
  font-weight: 700;
}
.history-head span {
  font-size: 12px;
  color: #6f867d;
}
.history-item {
  display: flex;
  gap: 10px;
  padding: 10px 0;
  border-bottom: 1px dashed #d9e2ef;
  font-size: 12px;
}
.history-item:last-child {
  border-bottom: none;
}
.history-index {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: #deece4;
  color: #2f605a;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  flex-shrink: 0;
}
.history-body {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.history-item .q {
  color: #24453f;
  font-weight: 600;
}
.history-item .a {
  color: #50675e;
  line-height: 1.6;
}

.history-item .meta {
  color: #6f867d;
}

.history-pagination {
  margin-top: 8px;
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 720px) {
  .dashboard-cards {
    grid-template-columns: minmax(0, 1fr);
  }
  .head-actions {
    width: 100%;
    justify-content: flex-start;
  }
  .bubble {
    max-width: 100%;
  }
  .bubble-meta-row,
  .conversation-header {
    flex-direction: column;
    align-items: flex-start;
  }
	.panel-head,
	.history-head {
		flex-direction: column;
		align-items: flex-start;
	}
}
</style>