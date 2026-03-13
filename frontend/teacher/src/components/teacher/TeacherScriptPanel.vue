<template>
  <div class="tab-content">
    <div class="course-preview">
      <div class="preview-title">
        <h4>第{{ currentEditPage }}页预览</h4>
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
          <small>完整预览请切换至「课件预览」标签</small>
        </div>
      </div>
      <div class="no-preview" v-else>请先选择课件</div>
    </div>

    <div class="script-editor">
      <div class="editor-actions">
        <div class="action-info">
          <span class="page-badge">第 {{ currentEditPage }} 页讲稿</span>
          <span class="action-hint">节点会跟随保存，支持标题、正文和讲解时长微调</span>
        </div>
        <div class="action-buttons">
          <button @click="$emit('generate-ai-script')" class="ai-btn" :disabled="!currentCourseId || scriptGenerating || scriptSaving">
            {{ scriptGenerating ? 'AI 生成中...' : 'AI 生成讲稿' }}
          </button>
          <button class="ghost-btn" :disabled="!currentCourseId" @click="rebuildNodesFromScript">重建节点</button>
          <button @click="$emit('save-script')" class="save-btn" :disabled="!currentCourseId || scriptSaving || scriptGenerating">
            {{ scriptSaving ? '保存中...' : '保存讲稿与节点' }}
          </button>
        </div>
      </div>

      <div class="editor-layout">
        <section class="script-pane">
          <div class="section-header">
            <h5>整页讲稿</h5>
            <span>{{ scriptLengthLabel }}</span>
          </div>
          <textarea
            :value="currentScript"
            placeholder="请输入本页讲稿内容，支持 AI 生成..."
            class="script-textarea"
            @input="$emit('update:current-script', $event.target.value)"
          ></textarea>
        </section>

        <button class="toggle-nodes-btn" @click="isNodesVisible = !isNodesVisible" :title="isNodesVisible ? '收起节点' : '展开节点'">
          {{ isNodesVisible ? '❯' : '❮' }}
        </button>

        <aside class="node-editor" v-show="isNodesVisible">
          <div class="node-editor-header">
            <div>
              <h5>讲授节点</h5>
              <p>建议按节奏拆成 3-6 个节点，学生端时间轴会直接使用这里的内容。</p>
            </div>
            <div class="header-actions">
              <span class="node-count">{{ localNodes.length }}</span>
              <button class="mini-btn" @click="addNode">新增节点</button>
            </div>
          </div>

          <div v-if="localNodes.length" class="node-list">
            <article
              v-for="(node, index) in localNodes"
              :key="node.id || node.nodeId || index"
              class="node-item"
              :class="`node-${node.type}`"
            >
              <div class="node-meta">
                <span class="node-id">{{ node.nodeId || `p${currentEditPage}_n${index + 1}` }}</span>
                <span class="node-type-tag">{{ typeLabel(node.type) }}</span>
              </div>

              <input
                :value="node.title"
                class="node-title-input"
                placeholder="节点标题"
                @input="updateNode(index, 'title', $event.target.value)"
              />

              <div class="node-inline-fields">
                <label>
                  <span>摘要</span>
                  <input
                    :value="node.summary"
                    placeholder="可选，给教师侧快速浏览用"
                    @input="updateNode(index, 'summary', $event.target.value)"
                  />
                </label>
                <label class="duration-field">
                  <span>时长</span>
                  <input
                    :value="node.estimatedDuration"
                    type="number"
                    min="10"
                    max="180"
                    @input="updateNode(index, 'estimatedDuration', $event.target.value)"
                  />
                </label>
              </div>

              <textarea
                :value="node.scriptText"
                class="node-textarea"
                placeholder="节点讲稿正文"
                @input="updateNode(index, 'scriptText', $event.target.value)"
              ></textarea>

              <div class="node-actions">
                <button class="mini-btn" :disabled="index === 0" @click="moveNode(index, -1)">上移</button>
                <button class="mini-btn" :disabled="index === localNodes.length - 1" @click="moveNode(index, 1)">下移</button>
                <button class="mini-btn danger" @click="removeNode(index)">删除</button>
              </div>
            </article>
          </div>

          <div v-else class="node-empty">
            <span class="node-empty-icon">📝</span>
            <p>还没有节点。可以先输入整页讲稿，然后点击“重建节点”。</p>
          </div>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
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

const emit = defineEmits(['generate-ai-script', 'save-script', 'update:current-script', 'update:current-script-nodes'])

const isNodesVisible = ref(window.innerWidth > 1200)
const imgError = ref(false)
const localNodes = ref([])

const syncNodesFromProps = (nodes) => {
  localNodes.value = (nodes || []).map((node, index, list) => ({
    id: node.id || '',
    nodeId: node.nodeId || `p${props.currentEditPage}_n${index + 1}`,
    type: node.type || inferNodeType(index, list.length),
    title: node.title || `第${props.currentEditPage}页节点${index + 1}`,
    summary: node.summary || '',
    scriptText: node.scriptText || '',
    reteachScript: node.reteachScript || '',
    transitionText: node.transitionText || '',
    estimatedDuration: Number(node.estimatedDuration) || estimateDuration(node.scriptText || node.summary || ''),
    sortOrder: Number(node.sortOrder) || index + 1
  }))
}

watch(() => props.currentScriptNodes, (value) => {
  syncNodesFromProps(value)
}, { immediate: true, deep: true })

watch(() => props.currentEditPage, () => { imgError.value = false })
watch(() => props.currentCourseId, () => { imgError.value = false })

const handleResize = () => {
  if (window.innerWidth < 1100 && isNodesVisible.value) {
    isNodesVisible.value = false
  }
}
window.addEventListener('resize', handleResize)
onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
})

const scriptLengthLabel = computed(() => {
  const count = String(props.currentScript || '').trim().length
  return count ? `${count} 字` : '未填写'
})

const emitNodes = () => {
  emit('update:current-script-nodes', localNodes.value.map((node, index, list) => ({
    ...node,
    type: inferNodeType(index, list.length),
    sortOrder: index + 1,
    estimatedDuration: Number(node.estimatedDuration) || estimateDuration(node.scriptText || node.summary || '')
  })))
}

const updateNode = (index, field, value) => {
  localNodes.value = localNodes.value.map((node, currentIndex) => {
    if (currentIndex !== index) return node
    return {
      ...node,
      [field]: field === 'estimatedDuration' ? Number(value) || 0 : value
    }
  })
  emitNodes()
}

const moveNode = (index, direction) => {
  const next = [...localNodes.value]
  const target = index + direction
  if (target < 0 || target >= next.length) return
  const current = next[index]
  next[index] = next[target]
  next[target] = current
  localNodes.value = next
  emitNodes()
}

const removeNode = (index) => {
  localNodes.value = localNodes.value.filter((_, currentIndex) => currentIndex !== index)
  emitNodes()
}

const addNode = () => {
  localNodes.value = [...localNodes.value, {
    id: '',
    nodeId: `p${props.currentEditPage}_n${localNodes.value.length + 1}`,
    type: inferNodeType(localNodes.value.length, localNodes.value.length + 1),
    title: `第${props.currentEditPage}页节点${localNodes.value.length + 1}`,
    summary: '',
    scriptText: '',
    reteachScript: '',
    transitionText: '',
    estimatedDuration: 30,
    sortOrder: localNodes.value.length + 1
  }]
  emitNodes()
}

const rebuildNodesFromScript = () => {
  const raw = String(props.currentScript || '').trim()
  if (!raw) {
    localNodes.value = []
    emitNodes()
    return
  }
  localNodes.value = raw
    .split(/\n{2,}|(?<=[。！？])\s*/)
    .map(item => item.trim())
    .filter(Boolean)
    .map((text, index, list) => ({
      id: '',
      nodeId: `p${props.currentEditPage}_n${index + 1}`,
      type: inferNodeType(index, list.length),
      title: `第${props.currentEditPage}页节点${index + 1}`,
      summary: text.length > 48 ? `${text.slice(0, 48)}...` : text,
      scriptText: text,
      reteachScript: '',
      transitionText: '',
      estimatedDuration: estimateDuration(text),
      sortOrder: index + 1
    }))
  emitNodes()
}

const typeLabel = (type) => {
  if (type === 'opening') return '开场'
  if (type === 'transition') return '过渡'
  return '讲解'
}

function inferNodeType(index, total) {
  if (total <= 1 || index === 0) return 'opening'
  if (index === total - 1) return 'transition'
  return 'explain'
}

function estimateDuration(text) {
  const size = Math.ceil(String(text || '').trim().length / 14)
  return Math.max(20, Math.min(90, size || 20))
}
</script>

<style scoped>
.tab-content {
  flex: 1;
  min-height: 0;
  display: flex;
  overflow: hidden;
  background: #f4f8ff;
}

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
 
.preview-iframe {
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

.script-editor {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

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
  color: #2F605A;
  background: #E8F0EF;
  padding: 3px 10px;
  border-radius: 6px;
  border: 1px solid #bfdbfe;
  white-space: nowrap;
}

.action-hint {
  font-size: 12px;
  color: #94a3b8;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.action-buttons {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.ai-btn,
.save-btn,
.ghost-btn,
.mini-btn {
  border: none;
  padding: 8px 14px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: opacity 0.2s, transform 0.1s;
}

.ai-btn,
.save-btn {
  color: #fff;
}

.ai-btn {
  background: #356F68;
}

.save-btn {
  background: #2F605A;
}

.ghost-btn,
.mini-btn {
  background: #E8F0EF;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
}

.mini-btn.danger {
  background: #fff1f2;
  color: #be123c;
  border-color: #fecdd3;
}

.ai-btn:disabled,
.save-btn:disabled,
.ghost-btn:disabled,
.mini-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.editor-layout {
  flex: 1;
  min-height: 0;
  display: flex;
  overflow: hidden;
}

.script-pane {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  background: #fff;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 18px 10px;
  color: #64748b;
  font-size: 12px;
  border-bottom: 1px solid #eef2f7;
}

.section-header h5 {
  margin: 0;
  color: #334155;
  font-size: 13px;
}

.script-textarea {
  flex: 1;
  min-width: 0;
  height: 100%;
  resize: none;
  border: none;
  padding: 20px;
  font-size: 14px;
  line-height: 1.85;
  outline: none;
  background: #fafcff;
  font-family: inherit;
  color: #1e293b;
}

.toggle-nodes-btn {
  width: 20px;
  background: #fff;
  border: none;
  border-left: 1px solid #e6ecf5;
  border-right: 1px solid #e6ecf5;
  color: #94a3b8;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  transition: all 0.2s;
  flex-shrink: 0;
}

.node-editor {
  flex: 0 0 380px;
  min-width: 280px;
  overflow-y: auto;
  padding: 14px 16px;
  background: #f8fbff;
}

.node-editor-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 14px;
  padding-bottom: 10px;
  border-bottom: 1px solid #e6ecf5;
}

.node-editor-header h5 {
  margin: 0 0 6px;
  font-size: 13px;
  font-weight: 600;
  color: #334155;
}

.node-editor-header p {
  margin: 0;
  font-size: 12px;
  color: #64748b;
  line-height: 1.5;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.node-count {
  font-size: 11px;
  font-weight: 600;
  color: #fff;
  background: #2F605A;
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
  gap: 10px;
}

.node-item {
  border-radius: 12px;
  padding: 12px;
  background: #fff;
  border: 1px solid #e2e8f0;
  box-shadow: 0 6px 18px rgba(15, 23, 42, 0.05);
}

.node-opening {
  border-left: 3px solid #356F68;
}

.node-explain {
  border-left: 3px solid #2F605A;
}

.node-transition {
  border-left: 3px solid #10b981;
}

.node-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.node-id {
  font-size: 11px;
  color: #94a3b8;
  font-family: Consolas, monospace;
}

.node-type-tag {
  font-size: 11px;
  font-weight: 600;
  padding: 1px 7px;
  border-radius: 8px;
  color: #1d4ed8;
  background: #E8F0EF;
}

.node-opening .node-type-tag {
  color: #0369a1;
  background: #e0f2fe;
}

.node-transition .node-type-tag {
  color: #065f46;
  background: #d1fae5;
}

.node-title-input,
.node-inline-fields input,
.node-textarea {
  width: 100%;
  border: 1px solid #dbe3ef;
  border-radius: 8px;
  padding: 9px 10px;
  font-family: inherit;
  box-sizing: border-box;
}

.node-title-input {
  margin-bottom: 8px;
  font-size: 13px;
  font-weight: 600;
}

.node-inline-fields {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 88px;
  gap: 8px;
  margin-bottom: 8px;
}

.node-inline-fields label {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 12px;
  color: #64748b;
}

.node-textarea {
  min-height: 96px;
  resize: vertical;
  line-height: 1.7;
  margin-bottom: 8px;
}

.node-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
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

@media (max-width: 1180px) {
  .node-editor {
    flex-basis: 320px;
    min-width: 240px;
  }
}

@media (max-width: 960px) {
  .course-preview {
    display: none;
  }
}
</style>