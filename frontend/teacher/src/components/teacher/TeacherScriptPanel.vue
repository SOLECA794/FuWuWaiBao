<template>
  <div class="editor-workbench">
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
          <iframe v-if="previewUrl" :src="previewUrl" title="课件预览" class="preview-iframe"></iframe>
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

    <aside class="copilot-panel">
      <h4>上下文智能助手与编辑器</h4>

      <section class="panel-block">
        <h5>智能助手</h5>
        <div class="ai-inline">
          <input v-model="aiPrompt" placeholder="请输入智能指令..." />
        </div>
        <button class="copilot-action" @click="$emit('generate-ai-script')" :disabled="!currentCourseId || scriptGenerating || scriptSaving">
          智能推荐资源
        </button>
      </section>

      <section class="panel-block">
        <h5>助手推荐</h5>
        <ul class="recommend-list">
          <li v-for="(item, idx) in recommendationItems" :key="`${item.title}_${idx}`" @click="selectedNodeIndex = idx">
            <strong>{{ item.title }}</strong>
            <span>{{ item.desc }}</span>
          </li>
        </ul>
      </section>

      <section class="panel-block grow">
        <h5>关联内容编辑器</h5>
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
const aiPrompt = ref('')

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

const recommendationItems = computed(() => {
  return timelineNodes.value.slice(0, 4).map((node, index) => ({
    title: node.title || `代码片段 ${index + 1}`,
    desc: (node.summary || node.scriptText || '可根据当前页内容做精讲和提问设计').slice(0, 44)
  }))
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
  flex: 1;
  min-height: 0;
  display: flex;
  overflow: hidden;
  background: linear-gradient(180deg, #f7fbf9 0%, #edf3ef 100%);
}

.viewer-stage {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  padding: 18px 20px 10px;
}

.stage-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
}

.stage-head h4 {
  margin: 0;
  font-size: 34px;
  font-weight: 500;
  color: #314641;
}

.slide-index {
  font-size: 32px;
  color: #6b847b;
}

.slide-wrap {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 44px;
}

.slide-nav-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 34px;
  height: 34px;
  border-radius: 999px;
  border: 1px solid #c6d8ce;
  background: rgba(244, 251, 247, 0.95);
  color: #406056;
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
  width: min(100%, 620px);
  min-height: 290px;
  background: #eef4f1;
  border: 1px solid #ceddd4;
  border-radius: 12px;
  padding: 12px;
  box-sizing: border-box;
}

.preview-iframe {
  width: 100%;
  min-height: 260px;
  border-radius: 8px;
  background: #fbfdfc;
  border: 1px solid #d6e3dc;
  object-fit: contain;
  display: block;
  height: 100%;
}

.preview-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 56px 16px;
  text-align: center;
  border: 2px dashed #d6e3dc;
  border-radius: 12px;
  background: #f6fbf8;
  gap: 8px;
  color: #6b7f78;
}

.preview-placeholder small {
  font-size: 13px;
  color: #93a69f;
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
  margin-top: 20px;
  padding: 26px 10px 14px;
  min-height: 176px;
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
  top: 104px;
  border-top: 2px solid #bccdc4;
  pointer-events: none;
}

.timeline-curve {
  position: absolute;
  left: 22px;
  right: 22px;
  top: 84px;
  height: 48px;
  border-bottom: 2px solid #d4e0da;
  border-radius: 0 0 56px 56px;
  opacity: 0.75;
  pointer-events: none;
}

.timeline-node {
  position: relative;
  z-index: 1;
  text-align: center;
  cursor: pointer;
  flex: 0 0 176px;
}

.node-bubble {
  min-height: 64px;
  border-radius: 12px;
  border: 1px solid #d4e1da;
  background: #e8f0ec;
  padding: 8px 10px;
  color: #3d4f49;
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
  color: #6f867d;
}

.node-dot {
  display: inline-block;
  margin-top: 16px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #dbe7e1;
  border: 3px solid #f6faf8;
  box-shadow: 0 2px 8px rgba(68, 96, 86, 0.18);
}

.timeline-node.active .node-bubble {
  border-color: #9bb8ad;
  background: #d9e7df;
}

.timeline-node.active .node-dot {
  background: #b7d1c4;
}

.stage-actions {
  margin-top: auto;
  padding: 8px 0;
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
.save-btn,
.copilot-action {
  border: 1px solid #d0dfd7;
  background: #edf5f1;
  color: #3c524b;
  padding: 8px 12px;
  border-radius: 10px;
  cursor: pointer;
  font-size: 13px;
}

.save-btn,
.copilot-action {
  background: #dbe8e1;
  border-color: #c4d7cd;
  font-weight: 600;
}

button:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.copilot-panel {
  flex: 0 0 340px;
  width: 340px;
  border-left: 1px solid #d8e4dc;
  background: #f8fbf9;
  padding: 14px 12px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 0;
  overflow-y: auto;
  padding-bottom: 16px;
  scrollbar-width: thin;
  scrollbar-color: #b4c8bc #edf4ef;
}

.copilot-panel::-webkit-scrollbar {
  width: 8px;
}

.copilot-panel::-webkit-scrollbar-track {
  background: #edf4ef;
  border-radius: 999px;
}

.copilot-panel::-webkit-scrollbar-thumb {
  background: #b4c8bc;
  border-radius: 999px;
}

.copilot-panel h4 {
  margin: 0;
  font-size: 18px;
  color: #314641;
}

.panel-block {
  border: 1px solid #dce8e1;
  border-radius: 12px;
  padding: 10px;
  background: #ffffff;
  flex-shrink: 0;
}

.panel-block h5 {
  margin: 0 0 8px;
  font-size: 14px;
  color: #425a51;
}

.ai-inline input,
.linked-editor {
  width: 100%;
  border: 1px solid #d9e4de;
  border-radius: 9px;
  padding: 8px 10px;
  box-sizing: border-box;
  font-size: 13px;
  font-family: inherit;
}

.recommend-list {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.recommend-list li {
  border: 1px solid #e2ebe6;
  background: #f7faf8;
  border-radius: 10px;
  padding: 8px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.recommend-list strong {
  font-size: 13px;
  color: #385049;
}

.recommend-list span {
  font-size: 12px;
  color: #768d84;
}

.grow {
  flex: 0 0 auto;
  min-height: 230px;
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
  border: 1px solid #dae6df;
  border-radius: 999px;
  background: #f5faf7;
  color: #546b62;
  padding: 4px 10px;
  font-size: 12px;
  cursor: pointer;
  white-space: nowrap;
}

.node-tabs button.active {
  background: #dfebe4;
  border-color: #c5d8ce;
  color: #355048;
}

.linked-editor {
  flex: 0 0 auto;
  height: 150px;
  min-height: 120px;
  margin-top: 8px;
  resize: none;
  line-height: 1.7;
  overflow-y: auto;
}

@media (max-width: 1200px) {
  .copilot-panel {
    flex-basis: 300px;
    width: 300px;
  }
}

@media (max-width: 980px) {
  .editor-workbench {
    flex-direction: column;
  }

  .copilot-panel {
    width: 100%;
    flex-basis: auto;
    border-left: none;
    border-top: 1px solid #d8e4dc;
  }
}

</style>