<template>
  <div class="tab-content">
    <div class="vertical-layout">
      <!-- 上半部分：课件预览标题与图 -->
      <div class="course-preview-section">
        <div class="preview-header">
          <h3 class="section-title">课件预览</h3>
          <div class="preview-controls" v-if="currentCourseId">
             <button class="nav-btn" @click="$emit('prev-page')" :disabled="currentEditPage <= 1">上一页</button>
             <span class="page-info">{{ currentEditPage }} / {{ totalPages }}</span>
             <button class="nav-btn" @click="$emit('next-page')" :disabled="currentEditPage >= totalPages">下一页</button>
          </div>
        </div>
        <div class="preview-img-wrap" v-if="currentCourseId">
          <iframe
            v-if="previewUrl"
            :src="previewUrl"
            title="课件预览"
            class="preview-iframe"
          ></iframe>
          <div class="preview-placeholder" v-else>
            <span class="placeholder-icon">📄</span>
            <p>第 {{ currentEditPage }} 页</p>
          </div>
        </div>
        <div class="no-preview" v-else>请先在右侧选择课件</div>
      </div>

      <!-- 中部：操作按钮 -->
      <div class="editor-actions">
        <button @click="$emit('generate-ai-script')" class="action-btn ai-btn" :disabled="!currentCourseId || scriptGenerating || scriptSaving">
          {{ scriptGenerating ? 'AI 生成中...' : 'AI 生成讲稿' }}
        </button>
        <button @click="$emit('save-script')" class="action-btn save-btn" :disabled="!currentCourseId || scriptSaving || scriptGenerating">
          {{ scriptSaving ? '保存中...' : '保存讲稿' }}
        </button>
      </div>

      <!-- 下半部分：文本框 -->
      <div class="script-editor-section">
        <h3 class="section-title">讲稿内容</h3>
        <textarea
          :value="currentScript"
          placeholder="请输入本页讲稿内容，支持 AI 生成..."
          class="script-textarea"
          @input="$emit('update:currentScript', $event.target.value)"
        ></textarea>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  previewUrl: {
    type: String,
    default: ''
  },
  currentCourseId: {
    type: String,
    default: ''
  },
  currentEditPage: {
    type: Number,
    default: 1
  },
  totalPages: {
    type: Number,
    default: 1
  },
  currentScript: {
    type: String,
    default: ''
  },
  currentScriptNodes: {
    type: Array,
    default: () => []
  },
  scriptGenerating: {
    type: Boolean,
    default: false
  },
  scriptSaving: {
    type: Boolean,
    default: false
  }
})

defineEmits(['generate-ai-script', 'save-script', 'update:currentScript', 'prev-page', 'next-page'])
</script>

<style scoped>
.tab-content {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  background: #ffffff;
  padding: 30px;
  overflow-y: auto;
}

.vertical-layout {
  display: flex;
  flex-direction: column;
  gap: 24px;
  max-width: 900px;
  margin: 0 auto;
  width: 100%;
}

.course-preview-section {
  display: flex;
  flex-direction: column;
  align-items: stretch;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.section-title {
  font-size: 20px;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
  letter-spacing: 1px;
}

.preview-controls {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-info {
  font-size: 14px;
  color: #64748b;
  font-weight: 500;
}

.nav-btn {
  padding: 6px 14px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: #fff;
  color: #334155;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
}

.nav-btn:hover:not(:disabled) {
  background: #f1f5f9;
  border-color: #cbd5e1;
}

.nav-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.preview-img-wrap {
  width: 100%;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 12px;
  background: #f8fafc;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}
 
.preview-iframe {
  width: 100%;
  height: 420px;
  border: none;
  border-radius: 8px;
  background: #ffffff;
}

.preview-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #94a3b8;
}

.placeholder-icon {
  font-size: 40px;
  opacity: 0.5;
  margin-bottom: 10px;
}

.no-preview {
  color: #94a3b8;
  padding: 40px;
  text-align: center;
  background: #f8fafc;
  border-radius: 12px;
  border: 1px dashed #e2e8f0;
}

.editor-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding: 0;
}

.action-btn {
  padding: 12px 32px;
  border-radius: 8px;
  border: 1.5px solid #2F605A;
  font-size: 15px;
  cursor: pointer;
  background: transparent;
  color: #2F605A;
  font-weight: 600;
  transition: all 0.2s;
  letter-spacing: 0.5px;
}

.action-btn:hover:not(:disabled) {
  background: #2F605A;
  color: #fff;
  transform: translateY(-1px);
}

.action-btn:active:not(:disabled) {
  transform: translateY(0);
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  border-color: #94a3b8;
  color: #94a3b8;
}

.script-editor-section {
  display: flex;
  flex-direction: column;
}

.script-textarea {
  width: 100%;
  min-height: 200px;
  resize: vertical;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 20px;
  font-size: 16px;
  line-height: 1.8;
  background: #ffffff;
  color: #1e293b;
  outline: none;
  font-family: inherit;
  transition: all 0.2s;
  box-shadow: inset 0 2px 4px 0 rgba(0, 0, 0, 0.02);
  margin-top: 12px;
}

.script-textarea:focus {
  border-color: #2F605A;
  box-shadow: 0 0 0 3px rgba(47, 96, 90, 0.1);
}

.script-textarea::placeholder {
  color: #cbd5e1;
}
</style>