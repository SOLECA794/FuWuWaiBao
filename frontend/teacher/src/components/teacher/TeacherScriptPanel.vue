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

      <footer class="stage-actions">
        <button class="ghost-btn" @click="splitCurrentBlock" :disabled="!editor || !currentCourseId">切分当前块</button>
        <button class="ghost-btn" @click="rebuildNodesFromScript" :disabled="!editor || !currentCourseId">按段落重建节点</button>
        <button class="ghost-btn" @click="$emit('generate-ai-script')" :disabled="!currentCourseId || scriptGenerating || scriptSaving">
          {{ scriptGenerating ? '智能生成中...' : '智能生成节点建议' }}
        </button>
        <button class="save-btn" @click="$emit('save-script')" :disabled="!currentCourseId || scriptSaving || scriptGenerating">
          {{ scriptSaving ? '保存中...' : '保存讲稿与节点' }}
        </button>
      </footer>
    </section>

    <section class="script-studio">
      <header class="studio-head">
        <h4>块级讲稿编辑器</h4>
        <span class="autosave-tip">映射自动保存：{{ autosaveState }}</span>
      </header>

      <div class="studio-split">
        <div class="editor-pane">
          <div class="mapping-board">
            <aside class="knowledge-list">
              <h5>知识节点</h5>
              <div
                v-for="(item, index) in blockSnapshots"
                :key="item.nodeId"
                class="knowledge-item"
                :class="{ active: activeNodeId === item.nodeId }"
                draggable="true"
                @dragstart="onNodeDragStart($event, item.nodeId)"
                @dragover.prevent
                @drop.prevent="onNodeDrop(index)"
                @click="focusNode(item.nodeId)"
              >
                <span class="node-index">{{ index + 1 }}</span>
                <div class="node-meta">
                  <strong>{{ item.nodeId }}</strong>
                  <p>{{ item.textPreview }}</p>
                </div>
              </div>
            </aside>

            <div class="editor-canvas" @dragover.prevent @drop.prevent="onEditorDrop">
              <EditorContent v-if="editor" :editor="editor" class="tiptap-root" />
            </div>
          </div>
        </div>

        <div class="preview-pane">
          <h5>Markdown 实时预览</h5>
          <div class="markdown-preview">
            <UnifiedMdRenderer :md-text="currentScript" />
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { Extension } from '@tiptap/core'
import StarterKit from '@tiptap/starter-kit'
import { Editor, EditorContent } from '@tiptap/vue-3'
import UnifiedMdRenderer from '../common/UnifiedMdRenderer.vue'

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

const emit = defineEmits([
  'generate-ai-script',
  'save-script',
  'update:current-script',
  'update:current-script-nodes',
  'autosave-mapping',
  'prev-page',
  'next-page'
])

const editor = ref(null)
const blockSnapshots = ref([])
const activeNodeId = ref('')
const draggedNodeId = ref('')
const draggedNodeIndex = ref(-1)
const autosaveState = ref('待触发')
let syncTimer = null
let autosaveTimer = null

const ParagraphNodeId = Extension.create({
  name: 'paragraphNodeId',
  addGlobalAttributes() {
    return [
      {
        types: ['paragraph'],
        attributes: {
          nodeId: {
            default: null,
            parseHTML: (element) => element.getAttribute('data-node-id') || null,
            renderHTML: (attributes) => {
              if (!attributes.nodeId) {
                return {}
              }
              return {
                'data-node-id': attributes.nodeId
              }
            }
          }
        }
      }
    ]
  }
})

watch(
  () => [props.currentCourseId, props.currentEditPage],
  () => {
    buildEditorFromNodes(props.currentScriptNodes, true)
  },
  { immediate: true }
)

watch(
  () => props.currentScriptNodes,
  (nodes) => {
    if (!editor.value) {
      buildEditorFromNodes(nodes, true)
      return
    }
    const currentNodeIds = blockSnapshots.value.map((item) => item.nodeId).join('|')
    const incomingNodeIds = (nodes || []).map((item, idx) => ensureNodeId(item?.nodeId, idx + 1)).join('|')
    if (currentNodeIds !== incomingNodeIds) {
      buildEditorFromNodes(nodes, false)
    }
  },
  { deep: true }
)

onBeforeUnmount(() => {
  if (syncTimer) window.clearTimeout(syncTimer)
  if (autosaveTimer) window.clearTimeout(autosaveTimer)
  if (editor.value) editor.value.destroy()
})

function buildEditorFromNodes(nodes, shouldCreate) {
  const safeNodes = normalizeNodes(nodes)
  const content = buildEditorContentFromNodes(safeNodes)

  if (shouldCreate || !editor.value) {
    if (editor.value) {
      editor.value.destroy()
    }
    editor.value = new Editor({
      extensions: [StarterKit, ParagraphNodeId],
      content,
      editorProps: {
        attributes: {
          class: 'block-editor'
        },
        handleClick(view, pos) {
          const node = view.state.doc.nodeAt(pos)
          if (node?.type?.name === 'paragraph') {
            activeNodeId.value = node.attrs.nodeId || ''
          }
          return false
        }
      },
      onSelectionUpdate: ({ editor: editorInstance }) => {
        const paragraph = findCurrentParagraph(editorInstance)
        if (paragraph?.nodeId) {
          activeNodeId.value = paragraph.nodeId
        }
      },
      onUpdate: ({ editor: editorInstance }) => {
        synchronizeState(editorInstance)
      }
    })
  } else {
    editor.value.commands.setContent(content, false)
    synchronizeState(editor.value)
  }
}

function synchronizeState(editorInstance) {
  if (syncTimer) window.clearTimeout(syncTimer)
  syncTimer = window.setTimeout(() => {
    normalizeParagraphNodeIds(editorInstance)
    const blocks = collectBlocks(editorInstance)
    blockSnapshots.value = blocks.map((block, index) => ({
      ...block,
      textPreview: String(block.text || '空段落').replace(/\s+/g, ' ').slice(0, 36) || '空段落',
      sortOrder: index + 1,
      type: inferNodeType(index, blocks.length)
    }))

    if (!activeNodeId.value && blockSnapshots.value.length > 0) {
      activeNodeId.value = blockSnapshots.value[0].nodeId
    }

    highlightActiveParagraph()

    const script = blockSnapshots.value
      .map((block) => String(block.text || '').trim())
      .filter(Boolean)
      .join('\n\n')

    emit('update:current-script', script)
    emit('update:current-script-nodes', blockSnapshots.value.map((item) => ({
      id: item.id || '',
      nodeId: item.nodeId,
      title: `节点${item.sortOrder}`,
      summary: item.textPreview,
      scriptText: item.text,
      reteachScript: '',
      transitionText: '',
      estimatedDuration: estimateDuration(item.text),
      sortOrder: item.sortOrder,
      type: item.type
    })))

    if (autosaveTimer) window.clearTimeout(autosaveTimer)
    autosaveState.value = '等待自动保存...'
    autosaveTimer = window.setTimeout(() => {
      autosaveState.value = '自动保存中...'
      emit('autosave-mapping', {
        content: script,
        nodes: blockSnapshots.value.map((item) => ({
          id: item.id || '',
          nodeId: item.nodeId,
          title: `节点${item.sortOrder}`,
          summary: item.textPreview,
          scriptText: item.text,
          reteachScript: '',
          transitionText: '',
          estimatedDuration: estimateDuration(item.text),
          sortOrder: item.sortOrder
        }))
      })
      window.setTimeout(() => {
        autosaveState.value = '已自动保存'
      }, 250)
    }, 900)
  }, 60)
}

function splitCurrentBlock() {
  if (!editor.value) return
  editor.value.chain().focus().splitBlock().run()
  normalizeParagraphNodeIds(editor.value)
  synchronizeState(editor.value)
}

function rebuildNodesFromScript() {
  if (!editor.value) return
  const lines = String(props.currentScript || '')
    .split(/\n{2,}/)
    .map((item) => item.trim())
    .filter(Boolean)

  const nodes = lines.length
    ? lines.map((text, idx) => ({
        nodeId: ensureNodeId('', idx + 1),
        scriptText: text
      }))
    : [{ nodeId: ensureNodeId('', 1), scriptText: '' }]

  buildEditorFromNodes(nodes, false)
}

function onNodeDragStart(event, nodeId) {
  draggedNodeId.value = nodeId
  if (event?.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('text/plain', nodeId)
  }
}

function onNodeDrop(targetIndex) {
  const sourceIndex = blockSnapshots.value.findIndex((item) => item.nodeId === draggedNodeId.value)
  if (sourceIndex < 0 || sourceIndex === targetIndex) return

  const next = [...blockSnapshots.value]
  const [moved] = next.splice(sourceIndex, 1)
  next.splice(targetIndex, 0, moved)
  applyBlocks(next)
}

function onEditorDrop(event) {
  const droppedNodeId = draggedNodeId.value || event?.dataTransfer?.getData('text/plain') || ''
  if (!droppedNodeId || !editor.value) return

  const targetParagraph = event.target?.closest?.('[data-node-id]')
  if (!targetParagraph) return

  const targetNodeId = targetParagraph.getAttribute('data-node-id') || ''
  if (!targetNodeId || targetNodeId === droppedNodeId) return

  const next = blockSnapshots.value.map((item) => ({ ...item }))
  const target = next.find((item) => item.nodeId === targetNodeId)
  const source = next.find((item) => item.nodeId === droppedNodeId)
  if (!target) return

  if (source) {
    source.nodeId = ensureNodeId('', Date.now())
  }
  target.nodeId = droppedNodeId
  activeNodeId.value = droppedNodeId
  applyBlocks(next)
}

function focusNode(nodeId) {
  activeNodeId.value = nodeId
  highlightActiveParagraph()

  if (!editor.value) return
  const block = blockSnapshots.value.find((item) => item.nodeId === nodeId)
  if (!block) return
  editor.value.chain().focus().setTextSelection(block.pos + 1).run()

  nextTick(() => {
    const paragraph = document.querySelector(`[data-node-id="${nodeId}"]`)
    paragraph?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  })
}

function applyBlocks(blocks) {
  if (!editor.value) return
  const sorted = blocks.map((item, index) => ({
    ...item,
    sortOrder: index + 1,
    nodeId: ensureNodeId(item.nodeId, index + 1)
  }))
  const content = buildEditorContentFromNodes(sorted)
  editor.value.commands.setContent(content, false)
  synchronizeState(editor.value)
}

function collectBlocks(editorInstance) {
  const result = []
  editorInstance.state.doc.descendants((node, pos) => {
    if (node.type.name !== 'paragraph') return true
    result.push({
      id: '',
      nodeId: ensureNodeId(node.attrs.nodeId, result.length + 1),
      text: node.textContent || '',
      pos
    })
    return true
  })
  if (result.length === 0) {
    result.push({ id: '', nodeId: ensureNodeId('', 1), text: '', pos: 0 })
  }
  return result
}

function normalizeParagraphNodeIds(editorInstance) {
  const tr = editorInstance.state.tr
  const seen = new Set()
  let index = 1

  editorInstance.state.doc.descendants((node, pos) => {
    if (node.type.name !== 'paragraph') return true
    let nodeId = node.attrs.nodeId
    if (!nodeId || seen.has(nodeId)) {
      nodeId = ensureNodeId('', index)
      tr.setNodeMarkup(pos, undefined, {
        ...node.attrs,
        nodeId
      })
    }
    seen.add(nodeId)
    index += 1
    return true
  })

  if (tr.docChanged) {
    editorInstance.view.dispatch(tr)
  }
}

function highlightActiveParagraph() {
  const root = document.querySelector('.block-editor')
  if (!root) return
  root.querySelectorAll('[data-node-id]').forEach((element) => {
    if (element.getAttribute('data-node-id') === activeNodeId.value) {
      element.classList.add('active-paragraph')
    } else {
      element.classList.remove('active-paragraph')
    }
  })
}

function findCurrentParagraph(editorInstance) {
  const { $from } = editorInstance.state.selection
  const node = $from.parent
  if (node.type.name !== 'paragraph') return null
  return {
    nodeId: node.attrs.nodeId,
    pos: $from.before($from.depth)
  }
}

function normalizeNodes(nodes) {
  const list = Array.isArray(nodes) ? nodes : []
  if (list.length > 0) {
    return list.map((node, index) => ({
      nodeId: ensureNodeId(node?.nodeId, index + 1),
      scriptText: String(node?.scriptText || node?.text || '').trim()
    }))
  }

  const fallback = String(props.currentScript || '').trim()
  if (!fallback) {
    return [{ nodeId: ensureNodeId('', 1), scriptText: '' }]
  }

  return fallback
    .split(/\n{2,}|(?<=[。！？])\s*/)
    .map((item) => item.trim())
    .filter(Boolean)
    .map((text, index) => ({
      nodeId: ensureNodeId('', index + 1),
      scriptText: text
    }))
}

function buildEditorContentFromNodes(nodes) {
  return {
    type: 'doc',
    content: nodes.map((item) => ({
      type: 'paragraph',
      attrs: {
        nodeId: ensureNodeId(item.nodeId, 1)
      },
      content: item.scriptText
        ? [
            {
              type: 'text',
              text: item.scriptText
            }
          ]
        : []
    }))
  }
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

function ensureNodeId(value, index) {
  const raw = String(value || '').trim()
  if (raw) return raw
  return `p${props.currentEditPage}_n${index}`
}
</script>

<style scoped>
.editor-workbench {
  flex: 1;
  min-height: 0;
  display: flex;
  overflow-y: auto;
  overflow-x: hidden;
  background: #f2f5f1;
}

.viewer-stage {
  width: 36%;
  min-width: 360px;
  border-right: 1px solid #d7e1d9;
  display: flex;
  flex-direction: column;
  padding: 16px;
  gap: 12px;
  background: linear-gradient(160deg, #f8fbf8 0%, #eef5f0 100%);
  overflow-y: auto;
}

.stage-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stage-head h4 {
  margin: 0;
  color: #27433b;
}

.slide-index {
  font-size: 14px;
  color: #4f6d64;
}

.slide-wrap {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.slide-nav-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 32px;
  height: 32px;
  border-radius: 999px;
  border: 1px solid #b8c9bf;
  background: #f6fbf8;
  color: #35554d;
  cursor: pointer;
  z-index: 2;
}

.slide-nav-btn.prev {
  left: 6px;
}

.slide-nav-btn.next {
  right: 6px;
}

.slide-canvas {
  width: 100%;
  min-height: 280px;
  border: 1px solid #d0ddd5;
  border-radius: 12px;
  background: #fff;
  overflow: hidden;
}

.preview-iframe {
  width: 100%;
  height: 100%;
  min-height: 280px;
  border: 0;
}

.preview-placeholder,
.no-preview {
  display: flex;
  min-height: 280px;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  color: #7d9289;
}

.stage-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.ghost-btn,
.save-btn {
  border: 1px solid #cad9cf;
  background: #f8fcf9;
  color: #2e4d43;
  border-radius: 8px;
  padding: 8px 12px;
  cursor: pointer;
}

.save-btn {
  background: #2f6158;
  border-color: #2f6158;
  color: #fff;
}

.script-studio {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  padding: 16px;
  gap: 12px;
  overflow-y: auto;
}

.studio-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.studio-head h4 {
  margin: 0;
  color: #244138;
}

.autosave-tip {
  color: #60796f;
  font-size: 13px;
}

.studio-split {
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.editor-pane,
.preview-pane {
  min-height: 0;
  border: 1px solid #d5e0d7;
  border-radius: 12px;
  background: #fbfdfc;
  overflow: hidden;
}

.mapping-board {
  height: 100%;
  display: grid;
  grid-template-columns: 240px 1fr;
  min-height: 0;
}

.knowledge-list {
  border-right: 1px solid #d9e5dc;
  padding: 10px;
  overflow: auto;
}

.knowledge-list h5,
.preview-pane h5 {
  margin: 0 0 10px;
  color: #466359;
}

.knowledge-item {
  display: flex;
  gap: 8px;
  padding: 8px;
  border: 1px solid #e2ece5;
  border-radius: 10px;
  background: #fff;
  margin-bottom: 8px;
  cursor: pointer;
}

.knowledge-item.active {
  border-color: #2f6158;
  box-shadow: inset 0 0 0 1px #2f6158;
  background: #f2f9f5;
}

.node-index {
  width: 22px;
  height: 22px;
  border-radius: 999px;
  background: #e4f0e8;
  color: #2f6158;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  flex-shrink: 0;
}

.node-meta {
  min-width: 0;
}

.node-meta strong {
  display: block;
  font-size: 12px;
  color: #38594f;
}

.node-meta p {
  margin: 3px 0 0;
  font-size: 12px;
  color: #72897f;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.editor-canvas {
  min-height: 0;
  overflow: auto;
  padding: 12px;
}

.tiptap-root :deep(.block-editor) {
  min-height: 100%;
  outline: none;
}

.tiptap-root :deep(.block-editor p) {
  border: 1px solid #d7e2da;
  border-radius: 10px;
  padding: 10px 12px;
  margin: 0 0 10px;
  background: #fff;
}

.tiptap-root :deep(.block-editor p::before) {
  content: attr(data-node-id);
  display: inline-block;
  font-size: 11px;
  color: #5d796f;
  background: #eef6f1;
  border: 1px solid #d7e8de;
  padding: 2px 6px;
  border-radius: 999px;
  margin-right: 8px;
}

.tiptap-root :deep(.block-editor p.active-paragraph) {
  border-color: #2f6158;
  box-shadow: 0 0 0 2px rgba(47, 97, 88, 0.14);
  background: #f4fbf7;
}

.preview-pane {
  display: flex;
  flex-direction: column;
  padding: 12px;
}

.markdown-preview {
  flex: 1;
  min-height: 0;
  overflow: auto;
  color: #2f433d;
  line-height: 1.62;
}

.markdown-preview :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 10px 0;
}

.markdown-preview :deep(th),
.markdown-preview :deep(td) {
  border: 1px solid #d4dfd7;
  padding: 6px 8px;
}

.markdown-preview :deep(img) {
  max-width: 100%;
  border-radius: 8px;
}

.markdown-preview :deep(pre) {
  background: #f4f7f5;
  border: 1px solid #d8e2da;
  border-radius: 8px;
  padding: 10px;
  overflow: auto;
}

.markdown-preview :deep(.mermaid-chart) {
  border: 1px solid #dae5dd;
  border-radius: 8px;
  padding: 10px;
  background: #fff;
}

@media (max-width: 1280px) {
  .editor-workbench {
    flex-direction: column;
  }

  .viewer-stage {
    width: 100%;
    min-width: 0;
    border-right: 0;
    border-bottom: 1px solid #d7e1d9;
  }

  .studio-split {
    grid-template-columns: 1fr;
  }

  .mapping-board {
    grid-template-columns: 1fr;
  }

  .knowledge-list {
    border-right: 0;
    border-bottom: 1px solid #d9e5dc;
    max-height: 220px;
  }
}
</style>
