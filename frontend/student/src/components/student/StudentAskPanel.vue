<template>
  <div class="panel-box">
    <div class="panel-head">
      <div>
        <div class="eyebrow">实时问答</div>
        <div class="question-header">围绕当前讲授节点发起追问</div>
      </div>
      <div class="session-chip" v-if="latestAnswerMeta.sessionId">
        会话 {{ latestAnswerMeta.sessionId.slice(0, 8) }}
      </div>
    </div>

    <div class="multi-modal-input">
      <el-input
        :model-value="question"
        type="textarea"
        placeholder="例如：这个公式为什么这样推导？能再举一个生活化例子吗？"
        :rows="5"
        @update:model-value="$emit('update:question', $event)"
      ></el-input>
      <div class="modal-tools">
        <el-button size="small" plain @click="$emit('open-upload')" icon="el-icon-upload">上传截图</el-button>
        <el-button size="small" plain icon="el-icon-crop">圈图提问</el-button>
        <el-button type="primary" size="small" :loading="askLoading" @click="$emit('send-question')" icon="el-icon-send">
          发送提问
        </el-button>
      </div>
    </div>

    <div class="ai-chat" v-if="aiReply || askLoading">
      <div class="chat-item teacher">
        <div class="chat-title-row">
          <div class="title">AI 助教回答</div>
          <span class="reply-state" :class="{ reteach: latestAnswerMeta.needReteach }">
            {{ askLoading ? '生成中' : latestAnswerMeta.needReteach ? '重讲模式' : '标准答疑' }}
          </span>
        </div>
        <div class="chat-content">{{ aiReply || '正在结合当前课件和上文追问生成回答...' }}</div>
        <div class="typing-hint" v-if="streamTypingActive">正在逐字生成中...</div>
        <div class="reply-meta" v-if="latestAnswerMeta.sourcePage">
          <span>来源页：第 {{ latestAnswerMeta.sourcePage }} 页</span>
          <span v-if="latestAnswerMeta.sourceNodeId">来源节点：{{ latestAnswerMeta.sourceNodeId }}</span>
          <span>{{ latestAnswerMeta.needReteach ? '系统判断你需要更通俗的解释' : '系统将按当前节奏继续讲解' }}</span>
        </div>
      </div>
    </div>

    <div class="follow-up-box" v-if="latestAnswerMeta.followUpSuggestion">
      <div class="follow-up-title">后续建议</div>
      <div class="follow-up-text">{{ latestAnswerMeta.followUpSuggestion }}</div>
    </div>

    <div class="question-history" v-if="qaHistory.length">
      <div class="history-head">
        <div class="title">最近提问</div>
        <span>保留最近 5 轮</span>
      </div>
      <div class="history-item" v-for="(item, idx) in qaHistory" :key="idx">
        <span class="history-index">0{{ idx + 1 }}</span>
        <div class="history-body">
          <span class="q">Q：{{ item.question }}</span>
          <span class="a">A：{{ item.answer }}</span>
          <span class="meta" v-if="item.sourcePage || item.sourceNodeId">
            来源：第 {{ item.sourcePage || '-' }} 页 {{ item.sourceNodeId ? `· ${item.sourceNodeId}` : '' }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
/* eslint-disable no-undef */
defineProps({
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
  }
})

defineEmits(['update:question', 'open-upload', 'send-question'])
</script>

<style scoped>
.panel-box {
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98) 0%, rgba(248, 250, 252, 0.98) 100%);
  border-radius: 20px;
  padding: 18px;
  border: 1px solid rgba(148, 163, 184, 0.18);
  box-shadow: 0 16px 30px rgba(15, 23, 42, 0.08);
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
  color: #0f766e;
}
.question-header {
  margin-top: 4px;
  font-size: 18px;
  font-weight: 700;
  color: #0f172a;
}
.session-chip {
  flex-shrink: 0;
  padding: 6px 10px;
  border-radius: 999px;
  background: rgba(15, 118, 110, 0.08);
  color: #0f766e;
  font-size: 12px;
  font-weight: 600;
}
.multi-modal-input {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.modal-tools {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}
.ai-chat {
  margin-top: 16px;
}
.chat-item {
  padding: 14px;
  border-radius: 16px;
  font-size: 14px;
  line-height: 1.75;
}
.teacher {
  background: linear-gradient(180deg, #f0fdfa 0%, #f8fafc 100%);
  border: 1px solid rgba(45, 212, 191, 0.2);
}
.chat-title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}
.chat-item .title {
  font-weight: 700;
  color: #0f172a;
}
.reply-state {
  padding: 4px 10px;
  border-radius: 999px;
  background: rgba(14, 165, 233, 0.12);
  color: #0369a1;
  font-size: 12px;
  font-weight: 600;
}
.reply-state.reteach {
  background: rgba(245, 158, 11, 0.14);
  color: #b45309;
}
.chat-content {
  color: #334155;
  white-space: pre-wrap;
}

.typing-hint {
  margin-top: 8px;
  font-size: 12px;
  color: #0f766e;
}
.reply-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 10px;
  font-size: 12px;
  color: #64748b;
}
.follow-up-box {
  margin-top: 14px;
  padding: 14px;
  border-radius: 16px;
  background: linear-gradient(180deg, rgba(250, 204, 21, 0.12) 0%, rgba(255, 255, 255, 0.96) 100%);
  border: 1px solid rgba(250, 204, 21, 0.2);
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
  background: rgba(248, 250, 252, 0.9);
  border: 1px solid rgba(148, 163, 184, 0.18);
  border-radius: 16px;
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
  color: #0f172a;
  font-weight: 700;
}
.history-head span {
  font-size: 12px;
  color: #64748b;
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
  background: #e2e8f0;
  color: #0f172a;
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
  color: #0f172a;
  font-weight: 600;
}
.history-item .a {
  color: #475569;
  line-height: 1.6;
}

.history-item .meta {
  color: #64748b;
}

@media (max-width: 720px) {
	.panel-head,
	.chat-title-row,
	.history-head {
		flex-direction: column;
		align-items: flex-start;
	}
}
</style>