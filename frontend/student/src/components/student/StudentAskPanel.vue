<template>
  <div class="panel-box">
    <div class="question-header">基于课件内容精准提问</div>
    <div class="multi-modal-input">
      <el-input
        :model-value="question"
        type="textarea"
        placeholder="请输入你的问题..."
        :rows="4"
        @update:model-value="$emit('update:question', $event)"
      ></el-input>
      <div class="modal-tools">
        <el-button size="small" @click="$emit('open-upload')" icon="el-icon-upload">上传截图</el-button>
        <el-button size="small" icon="el-icon-crop">圈图提问</el-button>
        <el-button type="primary" size="small" :loading="askLoading" @click="$emit('send-question')" icon="el-icon-send">
          发送提问
        </el-button>
      </div>
    </div>

    <div class="ai-chat" v-if="aiReply">
      <div class="chat-item teacher">
        <div class="title">AI 助教回答</div>
        <div>{{ aiReply }}</div>
      </div>
    </div>

    <div class="question-history" v-if="qaHistory.length">
      <div class="title">最近提问</div>
      <div class="history-item" v-for="(item, idx) in qaHistory" :key="idx">
        <span class="q">Q：{{ item.question }}</span>
        <span class="a">A：{{ item.answer }}</span>
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
  qaHistory: {
    type: Array,
    default: () => []
  }
})

defineEmits(['update:question', 'open-upload', 'send-question'])
</script>

<style scoped>
.panel-box {
  background: rgba(255, 255, 255, 0.96);
  border-radius: 14px;
  padding: 16px;
  border: 1px solid #e6ecf5;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.06);
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
.question-history {
  margin-top: 14px;
  background: #f8fbff;
  border: 1px solid #e6ecf5;
  border-radius: 8px;
  padding: 10px;
}
.question-history .title {
  font-size: 13px;
  color: #475569;
  margin-bottom: 8px;
}
.history-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 6px 0;
  border-bottom: 1px dashed #d9e2ef;
  font-size: 12px;
}
.history-item:last-child {
  border-bottom: none;
}
.history-item .q {
  color: #0f172a;
}
.history-item .a {
  color: #475569;
}
</style>