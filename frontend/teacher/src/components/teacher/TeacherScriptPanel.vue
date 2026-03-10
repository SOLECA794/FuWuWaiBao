<template>
  <div class="tab-content">
    <div class="course-preview">
      <div class="preview-title">
        <h4>第{{ currentEditPage }}页预览</h4>
      </div>
      <div class="preview-img-wrap" v-if="currentCourseId">
        <img
          v-if="!imgError"
          :src="previewUrl"
          alt="课件预览"
          class="preview-img"
          @error="imgError = true"
        />
        <div class="preview-placeholder" v-else>
          <span class="placeholder-icon">📄</span>
          <p>第 {{ currentEditPage }} 页</p>
          <small>完整预览请切换至「课件预览」标签</small>
        </div>
      </div>
      <div class="no-preview" v-else>请先选择课件</div>
    </div>

    <div class="script-editor">
      <div class="editor-actions">
        <div class="action-info">
          <span class="page-badge">第 {{ currentEditPage }} 页讲稿</span>
          <span class="action-hint">AI 生成后可手动调整</span>
        </div>
        <div class="action-buttons">
          <button @click="$emit('generate-ai-script')" class="ai-btn" :disabled="!currentCourseId || scriptGenerating || scriptSaving">
            {{ scriptGenerating ? 'AI 生成中...' : 'AI 生成讲稿' }}
          </button>
          <button @click="$emit('save-script')" class="save-btn" :disabled="!currentCourseId || scriptSaving || scriptGenerating">
            {{ scriptSaving ? '保存中...' : '保存讲稿' }}
          </button>
        </div>
      </div>

      <div class="editor-layout">
        <textarea
          :value="currentScript"
          placeholder="请输入本页讲稿内容，支持 AI 生成..."
          class="script-textarea"
          @input="$emit('update:currentScript', $event.target.value)"
        ></textarea>

        <button class="toggle-nodes-btn" @click="isNodesVisible = !isNodesVisible" :title="isNodesVisible ? '收起节点' : '展开节点'">
          {{ isNodesVisible ? '❯' : '❮' }}
        </button>

        <aside class="node-preview" v-show="isNodesVisible">
          <div class="node-preview-header">
            <h5>讲授节点</h5>
            <span class="node-count">{{ currentScriptNodes.length }}</span>
          </div>
          <div v-if="currentScriptNodes.length" class="node-list">
            <div
              v-for="node in currentScriptNodes"
              :key="node.nodeId"
              class="node-item"
              :class="`node-${node.type}`"
            >
              <div class="node-meta">
                <span class="node-id">{{ node.nodeId }}</span>
                <span class="node-type-tag">{{ node.type === 'opening' ? '开场' : node.type === 'transition' ? '过渡' : '讲解' }}</span>
              </div>
              <p>{{ node.text }}</p>
            </div>
          </div>
          <div v-else class="node-empty">
            <span class="node-empty-icon">📝</span>
            <p>讲稿生成后，这里会显示按教学节奏拆分的节点。</p>
          </div>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

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

const isNodesVisible = ref(window.innerWidth > 1200)

const imgError = ref(false)
watch(() => props.currentEditPage, () => { imgError.value = false })
watch(() => props.currentCourseId, () => { imgError.value = false })

// 监听窗口大小自动收起节点
const handleResize = () => {
  if (window.innerWidth < 1100 && isNodesVisible.value) {
    isNodesVisible.value = false
  }
}
window.addEventListener('resize', handleResize)

defineEmits(['generate-ai-script', 'save-script', 'update:currentScript'])
</script>

<style scoped>
/* 根元素填满 flex 父容器并内部水平裁剪 */
.tab-content {
  flex: 1;
  min-height: 0;
  display: flex;
  overflow: hidden;
  background: #f4f8ff;
}

/* 左侧课件预览栅：固定宽度，独立滚动 */
.course-preview {
  flex: 0 0 240px;
  min-width: 180px;
  overflow-y: auto;
  padding: 20px 16px;
  border-right: 1px solid #e6ecf5;
  background: #fff;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.preview-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.preview-title h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #334155;
}

.preview-img-wrap {
  flex: 1;
}

.preview-img {
  width: 100%;
  border-radius: 10px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  object-fit: contain;
  display: block;
  max-height: 360px;
}

.preview-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 16px;
  text-align: center;
  border: 2px dashed #e2e8f0;
  border-radius: 12px;
  background: #f8fafc;
  gap: 8px;
  color: #94a3b8;
}

.placeholder-icon {
  font-size: 36px;
  opacity: 0.7;
}

.preview-placeholder p {
  margin: 0;
  font-size: 15px;
  font-weight: 500;
  color: #64748b;
}

.preview-placeholder small {
  font-size: 12px;
  color: #b0bccc;
  line-height: 1.5;
}

.no-preview {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
  font-size: 13px;
  text-align: center;
  padding: 40px 0;
}

/* 右侧讲稿编辑区： flex 列方向，完整占满剩余宽度 */
.script-editor {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 操作栏：固定高度，不压缩 */
.editor-actions {
  flex-shrink: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  border-bottom: 1px solid #e6ecf5;
  background: #fff;
  gap: 8px;
  min-width: 0;
}

.action-info {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  overflow: hidden;
  flex-shrink: 1;
}

.page-badge {
  font-size: 12px;
  font-weight: 600;
  color: #2563eb;
  background: #eff6ff;
  padding: 3px 10px;
  border-radius: 6px;
  border: 1px solid #bfdbfe;
  white-space: nowrap;
  flex-shrink: 0;
}

.action-hint {
  font-size: 12px;
  color: #b0bccc;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}

.action-buttons {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.ai-btn,
.save-btn {
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  color: #fff;
  font-size: 13px;
  font-weight: 500;
  transition: opacity 0.2s, transform 0.1s;
  white-space: nowrap;
}

.ai-btn:active:not(:disabled),
.save-btn:active:not(:disabled) {
  transform: scale(0.97);
}

.ai-btn {
  background: #0ea5e9;
}

.save-btn {
  background: #2563eb;
}

.ai-btn:disabled,
.save-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

/* 编辑布局：消耗剩余高度，两列 */
.editor-layout {
  flex: 1;
  min-height: 0;
  display: flex;
  overflow: hidden;
}

/* 讲稿文本域：占满剩余宽度，无 resize */
.script-textarea {
  flex: 1;
  min-width: 0;
  height: 100%;
  resize: none;
  border: none;
  border-right: 1px solid #e6ecf5;
  padding: 20px;
  font-size: 14px;
  line-height: 1.85;
  outline: none;
  background: #fafcff;
  font-family: inherit;
  color: #1e293b;
  transition: background 0.2s;
}

.script-textarea:focus {
  background: #fff;
}

.script-textarea::placeholder {
  color: #b0bccc;
}

.toggle-nodes-btn {
  width: 20px;
  background: #fff;
  border: none;
  border-left: 1px solid #e6ecf5;
  color: #94a3b8;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  transition: all 0.2s;
  flex-shrink: 0;
}

.toggle-nodes-btn:hover {
  background: #f8fafc;
  color: #2563eb;
}

/* 节点预览切：固定宽度，独立滚动 */
.node-preview {
  flex: 0 0 230px;
  min-width: 160px;
  overflow-y: auto;
  padding: 14px 16px;
  background: #f8fbff;
}

.node-preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
  padding-bottom: 10px;
  border-bottom: 1px solid #e6ecf5;
}

.node-preview-header h5 {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  color: #334155;
}

.node-count {
  font-size: 11px;
  font-weight: 600;
  color: #fff;
  background: #2563eb;
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  border-radius: 10px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.node-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.node-item {
  border-radius: 8px;
  padding: 10px 12px;
  background: #fff;
  border: 1px solid #e2e8f0;
  transition: box-shadow 0.15s;
}

.node-item:hover {
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.09);
}

.node-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.node-id {
  font-size: 11px;
  color: #94a3b8;
  font-family: monospace;
}

.node-type-tag {
  font-size: 11px;
  font-weight: 600;
  padding: 1px 7px;
  border-radius: 4px;
}

.node-opening .node-type-tag {
  color: #0369a1;
  background: #e0f2fe;
}

.node-explain .node-type-tag {
  color: #1d4ed8;
  background: #eff6ff;
}

.node-transition .node-type-tag {
  color: #065f46;
  background: #d1fae5;
}

.node-item p {
  margin: 0;
  color: #334155;
  line-height: 1.6;
  font-size: 13px;
}

.node-opening {
  border-left: 3px solid #0ea5e9;
}

.node-explain {
  border-left: 3px solid #2563eb;
}

.node-transition {
  border-left: 3px solid #10b981;
}

.node-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 12px;
  text-align: center;
  color: #94a3b8;
  gap: 10px;
}

.node-empty-icon {
  font-size: 28px;
  opacity: 0.6;
}

.node-empty p {
  font-size: 13px;
  line-height: 1.6;
  margin: 0;
}
</style>