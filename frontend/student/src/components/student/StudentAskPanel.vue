<template>
  <div class="panel-box">
    <div class="panel-head">
      <div>
        <div class="eyebrow">问答工作区</div>
        <div class="question-header">摘要与追问一体化</div>
      </div>
      <div class="session-chip" v-if="latestAnswerMeta.sessionId">
        会话 {{ latestAnswerMeta.sessionId.slice(0, 8) }}
      </div>
    </div>

    <div class="summary-workbench">
      <div class="summary-head">
        <span>摘要模式</span>
        <el-select
          :model-value="summaryMode"
          size="small"
          class="summary-mode-select"
          @update:model-value="$emit('update:summaryMode', $event)"
        >
          <el-option label="速览模式" value="quick" />
          <el-option label="考试模式" value="exam" />
          <el-option label="讲解模式" value="teach" />
        </el-select>
      </div>
      <div class="summary-actions">
        <el-button size="small" type="primary" plain @click="$emit('generate-summary')">生成摘要</el-button>
        <el-button size="small" plain @click="$emit('use-summary')">用于提问</el-button>
        <el-button size="small" plain @click="$emit('clear-draft')">清空输入</el-button>
      </div>
      <div class="summary-content" v-if="mergedSummary">{{ mergedSummary }}</div>
      <div class="summary-empty" v-else>先生成摘要，再一键转成提问草稿。</div>
    </div>

    <div class="multi-modal-input">
      <el-input
        :model-value="question"
        type="textarea"
        placeholder="把你没听懂的地方直接问出来，例如：这一步为什么这么推？"
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
        <span>{{ historyStartIndex }}-{{ historyEndIndex }} / {{ qaHistory.length }}</span>
      </div>
      <div class="history-item" v-for="(item, idx) in pagedHistory" :key="historyStartIndex + idx">
        <span class="history-index">{{ String(historyStartIndex + idx).padStart(2, '0') }}</span>
        <div class="history-body">
          <span class="q">Q：{{ item.question }}</span>
          <span class="a">A：{{ item.answer }}</span>
          <span class="meta" v-if="item.sourcePage || item.sourceNodeId">
            来源：第 {{ item.sourcePage || '-' }} 页 {{ item.sourceNodeId ? `· ${item.sourceNodeId}` : '' }}
          </span>
        </div>
      </div>
      <div class="history-pagination" v-if="qaHistory.length > historyPageSize">
        <el-pagination
          small
          background
          layout="prev, pager, next"
          :current-page="historyPage"
          :page-size="historyPageSize"
          :total="qaHistory.length"
          @current-change="historyPage = $event"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
/* eslint-disable no-undef */
import { computed, ref, watch } from 'vue'

const historyPage = ref(1)
const historyPageSize = 3

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

const pagedHistory = computed(() => {
  const start = (historyPage.value - 1) * historyPageSize
  return props.qaHistory.slice(start, start + historyPageSize)
})

const historyStartIndex = computed(() => {
  if (!props.qaHistory.length) return 0
  return (historyPage.value - 1) * historyPageSize + 1
})

const historyEndIndex = computed(() => {
  if (!props.qaHistory.length) return 0
  return Math.min(historyStartIndex.value + pagedHistory.value.length - 1, props.qaHistory.length)
})

watch(() => props.qaHistory.length, (len) => {
  const pageCount = Math.max(1, Math.ceil(len / historyPageSize))
  if (historyPage.value > pageCount) {
    historyPage.value = pageCount
  }
})
</script>

<style scoped>
.panel-box {
  background: linear-gradient(180deg, #ffffff 0%, #f7faf8 100%);
  border-radius: 20px;
  padding: 18px;
  border: 1px solid #d9e7df;
  box-shadow: 0 16px 30px rgba(33, 61, 54, 0.08);
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
.summary-workbench {
  border: 1px solid #d9e7df;
  border-radius: 14px;
  background: linear-gradient(180deg, #f8fcfa 0%, #f2f8f5 100%);
  padding: 10px;
  margin-bottom: 12px;
}

.summary-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.summary-head span {
  font-size: 12px;
  font-weight: 700;
  color: #385c53;
}

.summary-mode-select {
  width: 120px;
}

.summary-actions {
  margin-top: 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.summary-content,
.summary-empty {
  margin-top: 8px;
  border-radius: 10px;
  padding: 9px 10px;
  font-size: 12px;
  line-height: 1.6;
}

.summary-content {
  border: 1px solid #d7e6de;
  background: #fff;
  color: #365950;
  white-space: pre-wrap;
}

.summary-empty {
  border: 1px dashed #cfddd5;
  color: #6f847c;
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
  background: linear-gradient(180deg, #f2f7f4 0%, #fbfdfc 100%);
  border: 1px solid #d8e7df;
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
  margin-top: 8px;
  font-size: 12px;
  color: #2f605a;
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
  background: #f6faf8;
  border: 1px solid #d9e7df;
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
	.panel-head,
	.chat-title-row,
	.history-head {
		flex-direction: column;
		align-items: flex-start;
	}
}
</style>