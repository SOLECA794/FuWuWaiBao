<template>
  <div class="editor-workbench">
    <header class="window-toolbar">
      <button class="window-toggle-btn" :class="{ active: isPreviewWindowOpen }" @click="toggleWindow('preview')">
        <span class="toggle-dot"></span>
        课件预览窗口
      </button>
      <button class="window-toggle-btn" :class="{ active: isEditorWindowOpen }" @click="toggleWindow('editor')">
        <span class="toggle-dot"></span>
        内容编辑器窗口
      </button>
    </header>

    <section v-if="!isPreviewWindowOpen && !isEditorWindowOpen" class="workbench-empty-state">
      <div class="empty-state-copy">
        <h4>暂无打开窗口</h4>
        <p>点击上方按钮打开课件预览或内容编辑器，窗口支持独立关闭与再次打开。</p>
      </div>
    </section>

    <div v-else class="window-grid" :class="windowGridClass">
      <section v-if="isPreviewWindowOpen" class="dock-window">
        <header class="dock-head">
          <div class="dock-title">
            <span class="dock-signal"></span>
          </div>
          <button class="dock-close-btn" @click="isPreviewWindowOpen = false" title="关闭预览窗口">×</button>
        </header>

        <section class="viewer-stage">
          <header class="stage-head">
            <h4>课件预览</h4>
            <div class="slide-index">第 {{ currentEditPage }}/{{ totalPages || 1 }} 页</div>
          </header>

          <div class="slide-wrap" v-if="currentCourseId">
            <button
              class="slide-nav-btn prev"
              @click="$emit('prev-page')"
              :disabled="currentEditPage <= 1"
              title="上一页"
            >
              ‹
            </button>
            <div class="slide-canvas">
              <img v-if="previewUrl" :src="previewUrl" alt="课件预览" class="preview-image" />
              <div class="preview-placeholder" v-else>
                <span>第 {{ currentEditPage }} 页</span>
                <small>当前课件暂未生成预览图</small>
              </div>
            </div>
            <button
              class="slide-nav-btn next"
              @click="$emit('next-page')"
              :disabled="currentEditPage >= totalPages"
              title="下一页"
            >
              ›
            </button>
          </div>
          <div class="no-preview" v-else>请先选择课件</div>

          <div class="timeline-area">
            <div class="timeline-curve"></div>
            <div class="timeline-line"></div>
            <div
              v-for="(node, index) in timelineNodes"
              :key="node.nodeId || index"
              class="timeline-node"
              :class="{ active: selectedNodeIndex === index }"
              @click="selectedNodeIndex = index"
            >
              <div class="node-bubble">
                <strong>{{ node.title || `节点 ${index + 1}` }}</strong>
                <span>{{ Number(node.estimatedDuration) || 20 }}秒</span>
              </div>
              <span class="node-dot"></span>
            </div>
          </div>

          <footer class="stage-actions">
            <div class="left-actions">
              <button class="ghost-btn" @click="addNode" :disabled="!currentCourseId">+ 新增节点</button>
              <button class="ghost-btn" @click="$emit('generate-ai-script')" :disabled="!currentCourseId || scriptGenerating || scriptSaving">
                {{ scriptGenerating ? '智能生成中...' : '智能生成节点建议' }}
              </button>
              <button class="ghost-btn" @click="rebuildNodesFromScript" :disabled="!currentCourseId">重建节点</button>
            </div>
            <div class="right-actions">
              <button class="pager-btn" @click="$emit('prev-page')" :disabled="currentEditPage <= 1">上一页</button>
              <button class="pager-btn" @click="$emit('next-page')" :disabled="currentEditPage >= totalPages">下一页</button>
              <button class="save-btn" @click="$emit('save-script')" :disabled="!currentCourseId || scriptSaving || scriptGenerating">
                {{ scriptSaving ? '保存中...' : '保存讲稿与节点' }}
              </button>
            </div>
          </footer>
        </section>
      </section>

      <section v-if="isEditorWindowOpen" class="dock-window">
        <header class="dock-head">
          <div class="dock-title">
            <span class="dock-signal"></span>
          </div>
          <button class="dock-close-btn" @click="isEditorWindowOpen = false" title="关闭内容编辑器">×</button>
        </header>

        <aside class="copilot-panel">
          <h4>关联内容编辑器</h4>

          <section class="panel-block grow">
            <h5>节点编辑</h5>
            <div class="node-tabs" v-if="timelineNodes.length">
              <button
                v-for="(node, idx) in timelineNodes"
                :key="node.nodeId || idx"
                :class="{ active: selectedNodeIndex === idx }"
                @click="selectedNodeIndex = idx"
              >
                节点{{ idx + 1 }}：{{ node.title || '未命名' }}
              </button>
            </div>
            <textarea
              class="linked-editor"
              :value="selectedNodeScript"
              placeholder="在这里编辑当前节点讲稿..."
              @input="updateSelectedNodeText($event.target.value)"
            ></textarea>
          </section>
        </aside>
      </section>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
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

const emit = defineEmits(['generate-ai-script', 'save-script', 'update:current-script', 'update:current-script-nodes', 'prev-page', 'next-page'])

const localNodes = ref([])
const selectedNodeIndex = ref(0)
const isPreviewWindowOpen = ref(false)
const isEditorWindowOpen = ref(false)

const windowGridClass = computed(() => {
  if (isPreviewWindowOpen.value && isEditorWindowOpen.value) return 'two-open'
  if (isPreviewWindowOpen.value) return 'single-open preview-only'
  if (isEditorWindowOpen.value) return 'single-open editor-only'
  return 'single-open'
})

const toggleWindow = (windowType) => {
  if (windowType === 'preview') {
    isPreviewWindowOpen.value = !isPreviewWindowOpen.value
    return
  }
  if (windowType === 'editor') {
    isEditorWindowOpen.value = !isEditorWindowOpen.value
  }
}

const syncNodesFromProps = (nodes) => {
  const mapped = (nodes || []).map((node, index, list) => ({
    id: node.id || '',
    nodeId: node.nodeId || `p${props.currentEditPage}_n${index + 1}`,
    type: node.type || inferNodeType(index, list.length),
    title: node.title || `${index === 0 ? '节点1：开场' : `节点${index + 1}：讲解`}`,
    summary: node.summary || node.text || '',
    scriptText: node.scriptText || node.text || '',
    reteachScript: node.reteachScript || '',
    transitionText: node.transitionText || '',
    estimatedDuration: Number(node.estimatedDuration) || estimateDuration(node.scriptText || node.summary || node.text || ''),
    sortOrder: Number(node.sortOrder) || index + 1
  }))

  localNodes.value = mapped
  if (selectedNodeIndex.value > mapped.length - 1) {
    selectedNodeIndex.value = Math.max(0, mapped.length - 1)
  }
}

watch(() => props.currentScriptNodes, (value) => {
  syncNodesFromProps(value)
}, { immediate: true, deep: true })

watch(() => props.currentEditPage, () => {
  selectedNodeIndex.value = 0
})

const emitNodes = () => {
  emit('update:current-script-nodes', localNodes.value.map((node, index, list) => ({
    ...node,
    type: inferNodeType(index, list.length),
    sortOrder: index + 1,
    estimatedDuration: Number(node.estimatedDuration) || estimateDuration(node.scriptText || node.summary || '')
  })))
}

const syncScriptFromNodes = () => {
  const merged = localNodes.value
    .map(node => String(node.scriptText || node.summary || '').trim())
    .filter(Boolean)
    .join('\n')
  emit('update:current-script', merged)
}

const timelineNodes = computed(() => {
  if (localNodes.value.length > 0) return localNodes.value
  return [{
    nodeId: `p${props.currentEditPage}_n1`,
    type: 'opening',
    title: '节点1：开场',
    scriptText: props.currentScript || '',
    estimatedDuration: 20
  }]
})

const selectedNodeScript = computed(() => {
  const node = timelineNodes.value[selectedNodeIndex.value]
  return node?.scriptText || ''
})

const updateNode = (index, field, value) => {
  localNodes.value = localNodes.value.map((node, currentIndex) => {
    if (currentIndex !== index) return node
    return {
      ...node,
      [field]: field === 'estimatedDuration' ? Number(value) || 0 : value
    }
  })
  emitNodes()
  syncScriptFromNodes()
}

const addNode = () => {
  localNodes.value = [...localNodes.value, {
    id: '',
    nodeId: `p${props.currentEditPage}_n${localNodes.value.length + 1}`,
    type: inferNodeType(localNodes.value.length, localNodes.value.length + 1),
    title: `节点${localNodes.value.length + 1}：讲解`,
    summary: '',
    scriptText: '',
    reteachScript: '',
    transitionText: '',
    estimatedDuration: 30,
    sortOrder: localNodes.value.length + 1
  }]
  selectedNodeIndex.value = localNodes.value.length - 1
  emitNodes()
  syncScriptFromNodes()
}

const rebuildNodesFromScript = () => {
  const raw = String(props.currentScript || '').trim()
  if (!raw) {
    localNodes.value = []
    emitNodes()
    emit('update:current-script', '')
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
      title: `${index === 0 ? '节点1：开场' : `节点${index + 1}：核心代码`}`,
      summary: text.length > 48 ? `${text.slice(0, 48)}...` : text,
      scriptText: text,
      reteachScript: '',
      transitionText: '',
      estimatedDuration: estimateDuration(text),
      sortOrder: index + 1
    }))
  selectedNodeIndex.value = 0
  emitNodes()
  syncScriptFromNodes()
}

const updateSelectedNodeText = (value) => {
  if (!localNodes.value.length) {
    localNodes.value = [{
      id: '',
      nodeId: `p${props.currentEditPage}_n1`,
      type: 'opening',
      title: '节点1：开场',
      summary: '',
      scriptText: value,
      reteachScript: '',
      transitionText: '',
      estimatedDuration: estimateDuration(value),
      sortOrder: 1
    }]
  } else {
    updateNode(selectedNodeIndex.value, 'scriptText', value)
    return
  }
  emitNodes()
  syncScriptFromNodes()
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
.editor-workbench {
  --wb-bg-top: #ffffff;
  --wb-bg-bottom: #f8fcfa;
  --wb-surface: #ffffff;
  --wb-surface-soft: #f4faf7;
  --wb-border: rgba(125, 162, 146, 0.22);
  --wb-border-strong: rgba(86, 130, 112, 0.34);
  --wb-text: #0f172a;
  --wb-muted: #5f7467;
  --wb-accent: #5ca68f;
  --wb-shadow: 0 12px 28px rgba(45, 79, 66, 0.12);

  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
  overflow: hidden;
  padding: 12px;
  box-sizing: border-box;
  background: linear-gradient(180deg, var(--wb-bg-top) 0%, var(--wb-bg-bottom) 100%);
}

.window-toolbar {
  /* 出现动画：用于卡片进入时的淡入 + 轻微上移 */
  animation: dock-appear 240ms cubic-bezier(.16,.84,.44,1) both;
  will-change: transform, opacity;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  padding: 0;
  border: 0;
  border-radius: 18px;
  background: transparent;
}

.window-toggle-btn {
  border: 0;
  border-radius: 18px;
  background: transparent;
  color: #5f7467;
  padding: 10px 16px;
  font-size: 14px;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.window-toggle-btn:hover {
  background: rgba(236, 247, 242, 0.9);
  color: #2f6052;
}

.window-toggle-btn.active {
  border: 0;
  background: rgba(227, 245, 238, 0.96);
  color: #2f6052;
  box-shadow: 0 6px 14px rgba(92, 166, 143, 0.2);
}

.toggle-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #94a3b8;
  box-shadow: 0 0 0 2px rgba(125, 162, 146, 0.2);
}

.window-toggle-btn.active .toggle-dot {
  background: var(--wb-accent);
  box-shadow: 0 0 0 2px rgba(92, 166, 143, 0.28);
}

.workbench-empty-state {
  flex: 1;
  min-height: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 28px 12px 14px;
  color: var(--wb-muted);
}

.empty-state-copy {
  text-align: center;
  max-width: 460px;
}

.empty-state-copy h4 {
  margin: 0 0 8px;
  font-size: 21px;
  font-weight: 600;
  color: #1e293b;
}

.empty-state-copy p {
  margin: 0;
  font-size: 13px;
  line-height: 1.65;
}

.window-grid {
  flex: 1;
  min-height: 0;
  display: grid;
  gap: 12px;
  grid-template-columns: minmax(0, 1.28fr) minmax(0, 1fr);
}

.window-grid.two-open {
  grid-template-columns: minmax(0, 1.28fr) minmax(0, 1fr);
}

.window-grid.single-open {
  align-items: stretch;
}

.window-grid.single-open .dock-window {
  height: 100%;
}

.window-grid.single-open.preview-only {
  grid-template-columns: minmax(0, 1.28fr) minmax(0, 1fr);
}
.window-grid.single-open.preview-only .dock-window {
  grid-column: 1 / 2;
}

.window-grid.single-open.editor-only {
  grid-template-columns: minmax(0, 1.28fr) minmax(0, 1fr);
}
.window-grid.single-open.editor-only .dock-window {
  grid-column: 2 / 3;
}

.dock-window {
  min-width: 0;
  min-height: 0;
  border: 0;
  border-radius: 22px;
  background: var(--wb-surface);
  box-shadow: var(--wb-shadow);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}

.dock-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px 8px;
  border-bottom: 0;
  background: transparent;
}

.dock-title {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: #334155;
  font-size: 15px;
  font-weight: 650;
}

.dock-signal {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: var(--wb-accent);
  box-shadow: 0 0 0 3px rgba(92, 166, 143, 0.22);
}

.dock-close-btn {
  width: 30px;
  height: 30px;
  border-radius: 10px;
  border: 0;
  background: rgba(255, 255, 255, 0.6);
  color: #64748b;
  font-size: 20px;
  line-height: 1;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.dock-close-btn:hover {
  background: rgba(255, 255, 255, 0.86);
  color: #334155;
}

.viewer-stage {
  flex: 1;
  min-width: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  padding: 4px 12px 12px;
  overflow-y: auto;
  overflow-x: hidden;
  scrollbar-width: thin;
  scrollbar-color: #c7cde8 #edf2f7;
}

.viewer-stage::-webkit-scrollbar {
  width: 8px;
}

.viewer-stage::-webkit-scrollbar-track {
  background: #edf2f7;
  border-radius: 999px;
}

.viewer-stage::-webkit-scrollbar-thumb {
  background: #c7cde8;
  border-radius: 999px;
}

.stage-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.stage-head h4 {
  margin: 0;
  font-size: clamp(28px, 2.5vw, 40px);
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--wb-text);
}

.slide-index {
  font-size: clamp(22px, 2.1vw, 32px);
  font-weight: 500;
  color: var(--wb-muted);
}

.slide-wrap {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 clamp(12px, 4vw, 40px);
}

.slide-nav-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 36px;
  height: 36px;
  border-radius: 999px;
  border: 0;
  background: rgba(255, 255, 255, 0.82);
  color: #475569;
  font-size: 24px;
  line-height: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 3;
}

.slide-nav-btn.prev {
  left: 0;
}

.slide-nav-btn.next {
  right: 0;
}

.slide-nav-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.slide-canvas {
  /* 宽度按列宽与视口比例自适应，最大不超过 820px；高度由 aspect-ratio 决定，保证按比例缩放 */ 
  width: min(100%, clamp(360px, 44vw, 820px));
  aspect-ratio: 16 / 9;
  max-height: 80vh;
  background: #ffffff;
  border: 1px solid var(--wb-border);
  border-radius: 16px;
  padding: 10px;
  box-sizing: border-box;
  transition: width 180ms ease, height 180ms ease;
}

.preview-image {
  width: 100%;
  height: 100%;
  border-radius: 14px;
  background: #ffffff;
  border: 0;
  object-fit: contain;
  display: block;
}

.preview-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 36px 14px;
  text-align: center;
  border-radius: 14px;
  background: #f7fcf9;
  gap: 8px;
  color: #6b7f78;
}

.preview-placeholder small {
  font-size: 13px;
  color: #8aa396;
  line-height: 1.5;
}

.no-preview {
  color: #93a69f;
  font-size: 14px;
  text-align: center;
  padding: 80px 0;
}

.timeline-area {
  position: relative;
  margin-top: 12px;
  padding: 18px 8px 10px;
  min-height: 136px;
  display: flex;
  align-items: flex-start;
  gap: 12px;
  overflow-x: auto;
  overflow-y: visible;
}

.timeline-line {
  position: absolute;
  left: 22px;
  right: 22px;
  top: 80px;
  border-top: 2px solid rgba(125, 162, 146, 0.36);
  pointer-events: none;
}

.timeline-curve {
  position: absolute;
  left: 22px;
  right: 22px;
  top: 64px;
  height: 34px;
  border-bottom: 2px solid rgba(125, 162, 146, 0.28);
  border-radius: 0 0 56px 56px;
  opacity: 0.75;
  pointer-events: none;
}

.timeline-node {
  position: relative;
  z-index: 1;
  text-align: center;
  cursor: pointer;
  flex: 0 0 150px;
}

.node-bubble {
  min-height: 54px;
  border-radius: 16px;
  border: 0;
  background: rgba(255, 255, 255, 0.95);
  padding: 8px 10px;
  color: #334155;
  display: flex;
  flex-direction: column;
  gap: 6px;
  justify-content: center;
}

.node-bubble strong {
  font-size: 14px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.node-bubble span {
  font-size: 12px;
  color: var(--wb-muted);
}

.node-dot {
  display: inline-block;
  margin-top: 12px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #d8e8e2;
  border: 0;
  box-shadow: 0 2px 8px rgba(87, 124, 109, 0.2);
}

.timeline-node.active .node-bubble {
  border-color: rgba(92, 166, 143, 0.38);
  background: rgba(228, 246, 239, 0.88);
}

.timeline-node.active .node-dot {
  background: rgba(92, 166, 143, 0.7);
}

.stage-actions {
  margin-top: auto;
  padding: 10px 0 2px;
  display: flex;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.left-actions,
.right-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.ghost-btn,
.pager-btn,
.save-btn {
  border: 0;
  background: rgba(255, 255, 255, 0.96);
  color: #334155;
  padding: 9px 13px;
  border-radius: 12px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.08);
  transition: all 0.2s ease;
}

.save-btn {
  background: linear-gradient(180deg, #e8f6ef 0%, #d9efe5 100%);
  font-weight: 600;
  color: #2f6052;
}

.ghost-btn:hover:not(:disabled),
.pager-btn:hover:not(:disabled),
.save-btn:hover:not(:disabled) {
  background: #ffffff;
  transform: translateY(-1px);
  box-shadow: 0 8px 18px rgba(15, 23, 42, 0.1);
}

button:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.copilot-panel {
  flex: 1;
  width: auto;
  border-left: none;
  background: transparent;
  padding: 10px 10px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 0;
  overflow-y: auto;
  padding-bottom: 16px;
  scrollbar-width: thin;
  scrollbar-color: #b8d1c6 #edf7f2;
}

.copilot-panel::-webkit-scrollbar {
  width: 8px;
}

.copilot-panel::-webkit-scrollbar-track {
  background: #edf7f2;
  border-radius: 999px;
}

.copilot-panel::-webkit-scrollbar-thumb {
  background: #b8d1c6;
  border-radius: 999px;
}

.copilot-panel h4 {
  margin: 0;
  font-size: 16px;
  color: #1e293b;
}

.panel-block {
  border: 1px solid var(--wb-border);
  border-radius: 16px;
  padding: 12px;
  background: #ffffff;
  flex-shrink: 0;
}

.panel-block h5 {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 700;
  color: #1f2937;
}

.linked-editor {
  width: 100%;
  border: 1px solid var(--wb-border);
  border-radius: 14px;
  background: #ffffff;
  padding: 12px;
  box-sizing: border-box;
  font-size: 15px;
  font-family: inherit;
}

.grow {
  flex: 0 0 auto;
  min-height: 180px;
  display: flex;
  flex-direction: column;
  overflow: visible;
}

.node-tabs {
  display: flex;
  gap: 6px;
  overflow-x: auto;
  padding-bottom: 8px;
}

.node-tabs button {
  border: 0;
  border-radius: 999px;
  background: rgba(242, 248, 245, 0.92);
  color: #5f7467;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
  white-space: nowrap;
}

.node-tabs button.active {
  background: rgba(226, 244, 237, 0.95);
  color: #2f6052;
}

.linked-editor {
  flex: 0 0 auto;
  height: 150px;
  min-height: 130px;
  margin-top: 10px;
  resize: none;
  line-height: 1.8;
  overflow-y: auto;
}

@media (max-width: 1200px) {
  .window-grid.two-open {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 980px) {
  .editor-workbench {
    padding: 8px;
  }

  .window-toolbar {
    gap: 6px;
  }

  .window-toggle-btn {
    flex: 1;
    justify-content: center;
  }
  
  /* 窄屏时单卡回退为全宽 */
  .window-grid.single-open {
    grid-template-columns: 1fr;
    justify-content: stretch;
    align-items: stretch;
  }

  .window-grid.single-open.preview-only .dock-window,
  .window-grid.single-open.editor-only .dock-window {
    width: 100%;
    max-width: none;
    grid-column: auto;
  }
}

@keyframes orbit-spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

@keyframes dock-appear {
  from {
    opacity: 0;
    transform: translateY(8px) scale(0.995);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@keyframes orbit-spin-rev {
  from {
    transform: rotate(360deg);
  }
  to {
    transform: rotate(0deg);
  }
}

@keyframes pulse {
  0%,
  100% {
    transform: scale(1);
    opacity: 0.86;
  }
  50% {
    transform: scale(1.24);
    opacity: 1;
  }
}

</style>