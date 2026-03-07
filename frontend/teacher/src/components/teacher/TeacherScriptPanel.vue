<template>
  <div class="tab-content">
    <div class="course-preview">
      <h4>第{{ currentEditPage }}页课件预览</h4>
      <img
        :src="previewUrl"
        alt="课件预览"
        class="preview-img"
        v-if="currentCourseId"
      />
      <div class="no-preview" v-else>请先上传并选择课件</div>
    </div>

    <div class="script-editor">
      <div class="editor-actions">
        <button @click="$emit('generate-ai-script')" class="ai-btn" :disabled="!currentCourseId || scriptGenerating || scriptSaving">
          {{ scriptGenerating ? 'AI 生成中...' : 'AI 生成讲稿' }}
        </button>
        <button @click="$emit('save-script')" class="save-btn" :disabled="!currentCourseId || scriptSaving || scriptGenerating">
          {{ scriptSaving ? '保存中...' : '保存讲稿' }}
        </button>
      </div>
      <textarea
        :value="currentScript"
        placeholder="请输入本页讲稿内容，支持AI生成..."
        class="script-textarea"
        @input="$emit('update:currentScript', $event.target.value)"
      ></textarea>
    </div>
  </div>
</template>

<script setup>
defineProps({
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
  currentScript: {
    type: String,
    default: ''
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

defineEmits(['generate-ai-script', 'save-script', 'update:currentScript'])
</script>

<style scoped>
.tab-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}
.course-preview,
.script-editor {
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid #e6ecf5;
  border-radius: 14px;
  padding: 16px;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.06);
}
.course-preview h4 {
  margin-bottom: 12px;
}
.preview-img {
  width: 100%;
  max-height: 520px;
  object-fit: contain;
  border-radius: 10px;
  background: #f8fafc;
}
.no-preview {
  min-height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
}
.editor-actions {
  display: flex;
  gap: 10px;
  margin-bottom: 12px;
}
.ai-btn,
.save-btn {
  border: none;
  padding: 10px 14px;
  border-radius: 8px;
  cursor: pointer;
  color: #fff;
}
.ai-btn {
  background: #0ea5e9;
}
.save-btn {
  background: #2563eb;
}
.ai-btn:disabled,
.save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.script-textarea {
  width: 100%;
  min-height: 420px;
  border: 1px solid #dbe3ef;
  border-radius: 10px;
  padding: 14px;
  resize: vertical;
  font-size: 14px;
  line-height: 1.7;
  outline: none;
}
</style>