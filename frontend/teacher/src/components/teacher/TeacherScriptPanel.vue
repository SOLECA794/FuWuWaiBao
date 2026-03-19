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
          :class="{ active: selectedNodeIndex === index, focused: focusedNodeId === (node.nodeId || ''), uncovered: effectiveUncoveredNodeIdSet.has(String(node.nodeId || '').trim()) }"
          :ref="(el) => setTimelineNodeRef(el, node.nodeId || `idx_${index}`)"
          @click="selectedNodeIndex = index"
          @dragover.prevent
          @drop="handleDropBindToNode(node.nodeId)"
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
          <button class="ghost-btn" @click="jumpToPriorityUncoveredNode" :disabled="!hasUncoveredNodes">优先补齐高频未覆盖节点</button>
          <button class="ghost-btn" @click="jumpToNextUncoveredNode" :disabled="!hasUncoveredNodes">下一个未覆盖节点</button>
          <span class="coverage-badge" v-if="hasUncoveredNodes">剩余未覆盖 {{ effectiveUncoveredNodeIdSet.size }}</span>
          <span class="coverage-badge neutral" v-if="localCoverageSummary.totalNodes > 0">
            本地预估覆盖 {{ localCoverageSummary.coveredNodes }}/{{ localCoverageSummary.totalNodes }}（{{ localCoverageSummary.coveragePercent }}%）
          </span>
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
        <div class="coverage-hint" v-if="selectedNodeIsUncovered">
          <span>当前节点尚未完成映射，建议优先补齐。</span>
          <span v-if="selectedNodePriorityScore > 0">提问热度分：{{ selectedNodePriorityScore }}</span>
          <button type="button" class="tiny-btn" @click="fillSelectedNodeMappingDraft">一键补最小映射骨架</button>
        </div>
        <div class="coverage-hint" v-if="selectedNodeNeedPrerequisiteHint">
          <span>当前节点缺少前置知识设置，建议补充 prerequisites 以提升续接质量。</span>
        </div>
        <div class="coverage-hint success" v-else-if="selectedNodeHasLocalDraft">
          <span>当前节点已添加映射草稿，保存后会参与正式覆盖率计算。</span>
        </div>
        <div class="node-tabs" v-if="timelineNodes.length" ref="nodeTabsRef">
          <button
            v-for="(node, idx) in timelineNodes"
            :key="node.nodeId || idx"
            :class="{ active: selectedNodeIndex === idx, focused: focusedNodeId === (node.nodeId || ''), uncovered: effectiveUncoveredNodeIdSet.has(String(node.nodeId || '').trim()) }"
            :ref="(el) => setNodeTabRef(el, node.nodeId || `idx_${idx}`)"
            @click="selectedNodeIndex = idx"
            @dragover.prevent
            @drop="handleDropBindToNode(node.nodeId)"
          >
            节点{{ idx + 1 }}：{{ node.title || '未命名' }}
          </button>
        </div>
        <div class="outline-board" v-if="timelineNodes.length">
          <div class="mapping-head">
            <span>大纲重组（拖拽排序自动同步讲稿）</span>
          </div>
          <div class="outline-list">
            <div
              v-for="(node, idx) in timelineNodes"
              :key="`outline_${node.nodeId || idx}`"
              class="outline-item"
              :class="{ active: selectedNodeIndex === idx }"
              draggable="true"
              @dragstart="handleOutlineDragStart(idx)"
              @dragover.prevent
              @drop="handleOutlineDrop(idx)"
              @click="selectedNodeIndex = idx"
            >
              <strong>{{ idx + 1 }}. {{ node.title || node.nodeId }}</strong>
              <small>{{ (node.scriptText || '').slice(0, 28) || '未填写讲稿' }}</small>
            </div>
          </div>
        </div>
        <div class="linked-editor" ref="linkedEditorRef">
          <div class="editor-toolbar">
            <button type="button" class="tiny-btn" @click="toggleBold">B</button>
            <button type="button" class="tiny-btn" @click="toggleItalic">I</button>
            <button type="button" class="tiny-btn" @click="toggleBulletList">列表</button>
            <button type="button" class="tiny-btn" @click="splitCurrentNodeByCaret">断点切分</button>
          </div>
          <EditorContent :editor="tiptapEditor" class="tiptap-editor" />
        </div>
        <div class="block-editor" v-if="scriptParagraphs.length">
          <div class="mapping-head">
            <span>块级讲稿编辑器（双向高亮）</span>
            <span>{{ scriptParagraphs.length }} 段</span>
          </div>
          <div class="block-list">
            <div
              v-for="segment in scriptParagraphs"
              :key="segment.segmentId"
              class="block-item"
              :class="{
                active: activeBlockSegmentId === segment.segmentId,
                related: selectedNodeRelatedSegmentIdSet.has(segment.segmentId)
              }"
              :ref="(el) => setScriptBlockRef(el, segment.segmentId)"
              @click="handleScriptBlockClick(segment.segmentId)"
              draggable="true"
              @dragstart="handleSegmentDragStart(segment.segmentId)"
              @dragend="handleSegmentDragEnd"
            >
              <div class="block-head">
                <strong>{{ segment.segmentId }}</strong>
                <span>{{ segmentNodeSummaryMap.get(segment.segmentId) || '未绑定节点' }}</span>
              </div>
              <div class="block-actions">
                <button type="button" class="tiny-btn" @click.stop="splitScriptBlock(segment.segmentId)">断点切分</button>
              </div>
              <textarea
                class="block-input"
                :ref="(el) => setBlockTextareaRef(el, segment.segmentId)"
                :value="segment.text"
                @click.stop
                @input="updateScriptBlockText(segment.segmentId, $event.target.value)"
              ></textarea>
            </div>
          </div>
        </div>
        <div class="mapping-editor" v-if="scriptParagraphs.length">
          <div class="mapping-head">
            <span>段落-节点映射（当前节点）</span>
            <span>{{ selectedNodeMappedSegmentCount }}/{{ scriptParagraphs.length }}（全局覆盖 {{ segmentCoverageSummary.coveredSegments }}/{{ segmentCoverageSummary.totalSegments }}）</span>
          </div>
          <div class="mapping-actions">
            <button type="button" class="tiny-btn" @click="bindAllSegmentsToSelectedNode" :disabled="!selectedNodeId || scriptParagraphs.length === 0">全选段落</button>
            <button type="button" class="tiny-btn danger" @click="clearSelectedNodeSegments" :disabled="!selectedNodeId">清空当前节点映射</button>
            <button type="button" class="tiny-btn" @click="bindUnmappedSegmentsToSelectedNode" :disabled="!selectedNodeId || segmentCoverageSummary.uncoveredSegmentIds.length === 0">绑定未映射段落</button>
            <button type="button" class="tiny-btn" @click="rebuildSegmentIdsFromScript" :disabled="scriptParagraphs.length === 0">按讲稿重建段落ID</button>
          </div>
          <div class="mapping-note" v-if="segmentCoverageSummary.uncoveredSegmentIds.length > 0">
            未映射段落：{{ segmentCoverageSummary.uncoveredSegmentIds.join('、') }}
          </div>
          <div class="mapping-note warn" v-if="selectedNodeMappingIssues.hasIssue">
            映射异常：
            <span v-if="selectedNodeMappingIssues.invalidSegmentIds.length">无效段落ID {{ selectedNodeMappingIssues.invalidSegmentIds.join('、') }}</span>
            <span v-if="selectedNodeMappingIssues.duplicateSegmentIds.length">；重复段落ID {{ selectedNodeMappingIssues.duplicateSegmentIds.join('、') }}</span>
            <button type="button" class="tiny-btn" @click="repairSelectedNodeSegments">一键修复映射</button>
          </div>
          <div class="segment-list">
            <button
              v-for="segment in scriptParagraphs"
              :key="segment.segmentId"
              type="button"
              class="segment-item"
              :class="{ active: selectedNodeMappedSegmentIdSet.has(segment.segmentId) }"
              @click="toggleSegmentBinding(segment.segmentId, segment.text)"
            >
              <strong>{{ segment.segmentId }}</strong>
              <span>{{ segment.text }}</span>
            </button>
          </div>
        </div>
        <div class="structured-toggle-row" v-if="timelineNodes.length">
          <button type="button" class="tiny-btn" @click="showStructuredEditors = !showStructuredEditors">
            {{ showStructuredEditors ? '收起结构化编辑区' : '展开结构化编辑区' }}
          </button>
        </div>
        <div class="knowledge-meta-editor" v-if="timelineNodes.length">
          <div class="mapping-head">
            <span>知识属性编辑（当前节点）</span>
            <div class="mapping-actions">
              <button type="button" class="tiny-btn" @click="autoFillPrerequisite">自动补前置节点</button>
              <button type="button" class="tiny-btn" @click="autoFillPrerequisitesForPage">全页自动串联</button>
            </div>
          </div>
          <div class="meta-grid">
            <label>
              <span>难度</span>
              <select :value="selectedNodeKnowledgeMeta.difficulty" @change="updateKnowledgeMetaField('difficulty', $event.target.value)">
                <option value="">未设置</option>
                <option value="easy">easy</option>
                <option value="medium">medium</option>
                <option value="hard">hard</option>
              </select>
            </label>
            <label>
              <span>标签（逗号分隔）</span>
              <input
                type="text"
                :value="selectedNodeKnowledgeMeta.tagsText"
                placeholder="例如：核心概念,易错点"
                @change="updateKnowledgeMetaField('tagsText', $event.target.value)"
              />
            </label>
            <label>
              <span>前置节点（逗号分隔）</span>
              <input
                type="text"
                :value="selectedNodeKnowledgeMeta.prerequisitesText"
                placeholder="例如：p1_n1,p1_n2"
                @change="updateKnowledgeMetaField('prerequisitesText', $event.target.value)"
              />
            </label>
          </div>
        </div>
        <div class="structured-editors" v-if="timelineNodes.length && showStructuredEditors" ref="structuredEditorsRef">
          <label>
            <span>结构化 MD</span>
            <textarea
              class="mini-editor"
              ref="structuredMarkdownEditorRef"
              :value="selectedNodeStructuredMarkdown"
              placeholder="可编辑结构化 Markdown（支持公式与图表占位）"
              @input="updateSelectedNodeField('structuredMarkdown', $event.target.value)"
            ></textarea>
            <div class="md-preview" v-if="selectedNodeStructuredMarkdown.trim()">
              <div class="preview-title">MD 预览</div>
              <div class="preview-body" v-html="structuredMarkdownPreview"></div>
            </div>
          </label>
          <label>
            <span>知识节点 JSON</span>
            <textarea
              class="mini-editor"
              :value="selectedNodeKnowledgeNodesJson"
              placeholder='如: [{"node_id":"..."}]'
              @input="updateSelectedNodeField('knowledgeNodesJson', $event.target.value)"
              @blur="formatSelectedNodeJSON('knowledgeNodesJson')"
            ></textarea>
            <div class="json-actions">
              <button type="button" class="tiny-btn" @click="formatSelectedNodeJSON('knowledgeNodesJson')">格式化</button>
              <button type="button" class="tiny-btn danger" @click="resetSelectedNodeJSON('knowledgeNodesJson')">重置为 []</button>
            </div>
            <div class="json-tip" :class="knowledgeNodesState.valid ? 'ok' : 'error'">
              {{ knowledgeNodesState.message }}
            </div>
          </label>
          <label>
            <span>讲稿段落映射 JSON</span>
            <textarea
              class="mini-editor"
              :value="selectedNodeScriptSegmentsJson"
              placeholder='如: [{"segment_id":"seg_1","node_ids":["..."]}]'
              @input="updateSelectedNodeField('scriptSegmentsJson', $event.target.value)"
              @blur="formatSelectedNodeJSON('scriptSegmentsJson')"
            ></textarea>
            <div class="json-actions">
              <button type="button" class="tiny-btn" @click="formatSelectedNodeJSON('scriptSegmentsJson')">格式化</button>
              <button type="button" class="tiny-btn danger" @click="resetSelectedNodeJSON('scriptSegmentsJson')">重置为 []</button>
            </div>
            <div class="json-tip" :class="scriptSegmentsState.valid ? 'ok' : 'error'">
              {{ scriptSegmentsState.message }}
            </div>
          </label>
        </div>
      </section>
    </aside>
  </div>
</template>

<script setup>
import { computed, nextTick, onUnmounted, ref, watch } from 'vue'
import { EditorContent, useEditor } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
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
  mappingCoverage: {
    type: Object,
    default: null
  },
  nodeStats: {
    type: Array,
    default: () => []
  },
  focusNodeRequest: {
    type: Object,
    default: null
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
const focusedNodeId = ref('')
const timelineNodeRefs = ref({})
const nodeTabsRef = ref(null)
const nodeTabRefs = ref({})
const linkedEditorRef = ref(null)
const structuredEditorsRef = ref(null)
const structuredMarkdownEditorRef = ref(null)
const scriptBlockRefs = ref({})
const blockTextareaRefs = ref({})
const showStructuredEditors = ref(true)
const pendingFocusNodeId = ref('')
const localDraftCoveredNodeIds = ref([])
const activeBlockSegmentId = ref('')
const draggingSegmentId = ref('')
const draggingOutlineIndex = ref(-1)
const syncingFromTipTap = ref(false)
let focusFlashTimer = null

const tiptapEditor = useEditor({
  content: '',
  extensions: [StarterKit],
  autofocus: false,
  editorProps: {
    attributes: {
      class: 'tiptap-editor-body'
    }
  },
  onUpdate: ({ editor }) => {
    if (syncingFromTipTap.value) return
    const text = editor.getText({ blockSeparator: '\n\n' })
    updateSelectedNodeText(text)
  }
})

const syncNodesFromProps = (nodes) => {
  const mapped = (nodes || []).map((node, index, list) => ({
    id: node.id || '',
    nodeId: node.nodeId || `p${props.currentEditPage}_n${index + 1}`,
    schemaVersion: Number(node.schemaVersion) || 2,
    type: node.type || inferNodeType(index, list.length),
    title: node.title || `${index === 0 ? '节点1：开场' : `节点${index + 1}：讲解`}`,
    summary: node.summary || node.text || '',
    scriptText: node.scriptText || node.text || '',
    reteachScript: node.reteachScript || '',
    transitionText: node.transitionText || '',
    structuredMarkdown: node.structuredMarkdown || '',
    knowledgeNodesJson: node.knowledgeNodesJson || '',
    scriptSegmentsJson: node.scriptSegmentsJson || '',
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
  localDraftCoveredNodeIds.value = []
})

watch(() => props.mappingCoverage, () => {
  // 覆盖率刷新后清空本地草稿标记，避免和服务端状态叠加造成误导。
  localDraftCoveredNodeIds.value = []
}, { deep: true })

watch(() => props.focusNodeRequest, (request) => {
  if (!request || !request.nodeId) return
  focusNodeById(request.nodeId)
}, { deep: true })

watch(localNodes, () => {
  if (pendingFocusNodeId.value) {
    const target = pendingFocusNodeId.value
    pendingFocusNodeId.value = ''
    focusNodeById(target)
  }
}, { deep: true })

onUnmounted(() => {
  if (focusFlashTimer) {
    window.clearTimeout(focusFlashTimer)
    focusFlashTimer = null
  }
  if (tiptapEditor.value) {
    tiptapEditor.value.destroy()
  }
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

const uncoveredNodeIdSet = computed(() => {
  const ids = props.mappingCoverage?.uncoveredNodeIds
  if (!Array.isArray(ids) || ids.length === 0) return new Set()
  return new Set(ids.map(item => String(item || '').trim()).filter(Boolean))
})

const localDraftCoveredNodeIdSet = computed(() => new Set((localDraftCoveredNodeIds.value || []).map(item => String(item || '').trim()).filter(Boolean)))

const effectiveUncoveredNodeIdSet = computed(() => {
  if (uncoveredNodeIdSet.value.size === 0) return new Set()
  const clone = new Set(uncoveredNodeIdSet.value)
  localDraftCoveredNodeIdSet.value.forEach((id) => clone.delete(id))
  return clone
})

const hasUncoveredNodes = computed(() => effectiveUncoveredNodeIdSet.value.size > 0)

const nodePriorityScoreMap = computed(() => {
  const list = Array.isArray(props.nodeStats) ? props.nodeStats : []
  const result = new Map()
  list.forEach((item = {}) => {
    const nodeId = String(item.nodeId || '').trim()
    if (!nodeId) return
    const dialogueCount = Number(item.dialogueCount) || 0
    const needReteachCount = Number(item.needReteachCount) || 0
    // 重讲请求优先级更高，按 2x 权重参与排序
    const score = dialogueCount + needReteachCount * 2
    result.set(nodeId, score)
  })
  return result
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

watch(selectedNodeScript, (value) => {
  if (!tiptapEditor.value) return
  const nextText = String(value || '')
  const currentText = tiptapEditor.value.getText({ blockSeparator: '\n\n' })
  if (currentText === nextText) return
  syncingFromTipTap.value = true
  const html = nextText
    ? `<p>${escapeHtml(nextText).replace(/\n\n+/g, '</p><p>').replace(/\n/g, '<br />')}</p>`
    : '<p></p>'
  tiptapEditor.value.commands.setContent(html, false)
  window.setTimeout(() => {
    syncingFromTipTap.value = false
  }, 0)
}, { immediate: true })

const selectedNodeStructuredMarkdown = computed(() => {
  const node = timelineNodes.value[selectedNodeIndex.value]
  return node?.structuredMarkdown || ''
})

const selectedNodeKnowledgeNodesJson = computed(() => {
  const node = timelineNodes.value[selectedNodeIndex.value]
  return node?.knowledgeNodesJson || ''
})

const selectedNodeScriptSegmentsJson = computed(() => {
  const node = timelineNodes.value[selectedNodeIndex.value]
  return node?.scriptSegmentsJson || ''
})

const localCoverageSummary = computed(() => {
  const nodeIds = timelineNodes.value
    .map(node => String(node?.nodeId || '').trim())
    .filter(Boolean)
  if (nodeIds.length === 0) {
    return { totalNodes: 0, coveredNodes: 0, coveragePercent: 0 }
  }

  const mappedSet = new Set()
  timelineNodes.value.forEach((node = {}) => {
    const raw = String(node.scriptSegmentsJson || '').trim()
    if (!raw) return
    try {
      const parsed = JSON.parse(raw)
      const list = Array.isArray(parsed) ? parsed : [parsed]
      list.forEach((segment = {}) => {
        const ids = Array.isArray(segment.node_ids) ? segment.node_ids : []
        ids.forEach((id) => {
          const normalized = String(id || '').trim()
          if (normalized) mappedSet.add(normalized)
        })
      })
    } catch {
      // 无效 JSON 不参与本地覆盖率统计
    }
  })

  const coveredNodes = nodeIds.filter(id => mappedSet.has(id)).length
  const coveragePercent = Math.round((coveredNodes / nodeIds.length) * 100)
  return {
    totalNodes: nodeIds.length,
    coveredNodes,
    coveragePercent
  }
})

const scriptParagraphs = computed(() => {
  const raw = String(props.currentScript || '').trim()
  if (!raw) return []
  return raw
    .split(/\n{2,}|(?<=[。！？])\s*/)
    .map(item => item.trim())
    .filter(Boolean)
    .map((text, idx) => ({ segmentId: `seg_${idx + 1}`, text }))
})

const segmentNodeIdMap = computed(() => {
  const map = new Map()
  timelineNodes.value.forEach((node = {}) => {
    const nodeId = String(node.nodeId || '').trim()
    if (!nodeId) return
    const segments = parseSegmentList(node.scriptSegmentsJson)
    segments.forEach((segment = {}) => {
      const segmentId = String(segment.segment_id || '').trim()
      if (!segmentId) return
      const ids = map.get(segmentId) || []
      if (!ids.includes(nodeId)) {
        ids.push(nodeId)
      }
      map.set(segmentId, ids)
    })
  })
  return map
})

const selectedNodeId = computed(() => String(timelineNodes.value[selectedNodeIndex.value]?.nodeId || '').trim())
const selectedNodeIsUncovered = computed(() => !!selectedNodeId.value && effectiveUncoveredNodeIdSet.value.has(selectedNodeId.value))
const selectedNodeHasLocalDraft = computed(() => !!selectedNodeId.value && localDraftCoveredNodeIdSet.value.has(selectedNodeId.value))
const selectedNodePriorityScore = computed(() => nodePriorityScoreMap.value.get(selectedNodeId.value) || 0)
const selectedNodeNeedPrerequisiteHint = computed(() => {
  const nodeId = selectedNodeId.value
  if (!nodeId) return false
  const knowledge = parseKnowledgeNodes(selectedNodeKnowledgeNodesJson.value)
  const current = knowledge.find(item => String(item?.node_id || '').trim() === nodeId)
  if (!current) return false
  const difficulty = String(current?.difficulty || '').trim().toLowerCase()
  const prerequisites = Array.isArray(current?.prerequisites) ? current.prerequisites : []
  return (difficulty === 'medium' || difficulty === 'hard') && prerequisites.length === 0
})

watch(selectedNodeId, (value) => {
  if (!value) return
  focusNodeRelatedBlock(value)
})

const selectedNodeMappedSegmentIdSet = computed(() => {
  const nodeId = selectedNodeId.value
  if (!nodeId) return new Set()
  const segments = parseSegmentList(selectedNodeScriptSegmentsJson.value)
  const result = new Set()
  segments.forEach((segment = {}) => {
    const ids = Array.isArray(segment.node_ids) ? segment.node_ids : []
    const contains = ids.map(item => String(item || '').trim()).includes(nodeId)
    if (contains) {
      const segmentId = String(segment.segment_id || '').trim()
      if (segmentId) result.add(segmentId)
    }
  })
  return result
})

const selectedNodeMappedSegmentCount = computed(() => selectedNodeMappedSegmentIdSet.value.size)

const selectedNodeRelatedSegmentIdSet = computed(() => {
  const nodeId = selectedNodeId.value
  if (!nodeId) return new Set()
  const result = new Set()
  segmentNodeIdMap.value.forEach((nodeIDs = [], segmentId) => {
    if (Array.isArray(nodeIDs) && nodeIDs.includes(nodeId)) {
      result.add(segmentId)
    }
  })
  return result
})

const segmentNodeSummaryMap = computed(() => {
  const titleByNodeId = new Map(
    timelineNodes.value.map(node => [String(node.nodeId || '').trim(), String(node.title || '').trim()])
  )
  const result = new Map()
  scriptParagraphs.value.forEach((segment) => {
    const nodeIDs = segmentNodeIdMap.value.get(segment.segmentId) || []
    if (!nodeIDs.length) {
      result.set(segment.segmentId, '')
      return
    }
    const labels = nodeIDs.map(id => titleByNodeId.get(id) || id)
    result.set(segment.segmentId, labels.join(' / '))
  })
  return result
})

const segmentCoverageSummary = computed(() => {
  const allSegmentIds = scriptParagraphs.value.map(item => item.segmentId)
  if (allSegmentIds.length === 0) {
    return { totalSegments: 0, coveredSegments: 0, uncoveredSegmentIds: [] }
  }
  const mappedSet = new Set()
  timelineNodes.value.forEach((node = {}) => {
    const list = parseSegmentList(node.scriptSegmentsJson)
    list.forEach((segment = {}) => {
      const segmentId = String(segment.segment_id || '').trim()
      if (segmentId) mappedSet.add(segmentId)
    })
  })
  const uncoveredSegmentIds = allSegmentIds.filter(id => !mappedSet.has(id))
  return {
    totalSegments: allSegmentIds.length,
    coveredSegments: allSegmentIds.length - uncoveredSegmentIds.length,
    uncoveredSegmentIds
  }
})

const selectedNodeMappingIssues = computed(() => {
  const validSegmentIdSet = new Set(scriptParagraphs.value.map(item => item.segmentId))
  const segments = parseSegmentList(selectedNodeScriptSegmentsJson.value)
  const seen = new Set()
  const duplicate = new Set()
  const invalid = new Set()

  segments.forEach((item = {}) => {
    const segmentId = String(item.segment_id || '').trim()
    if (!segmentId) return
    if (seen.has(segmentId)) {
      duplicate.add(segmentId)
    }
    seen.add(segmentId)
    if (!validSegmentIdSet.has(segmentId)) {
      invalid.add(segmentId)
    }
  })

  const invalidSegmentIds = Array.from(invalid)
  const duplicateSegmentIds = Array.from(duplicate)
  return {
    invalidSegmentIds,
    duplicateSegmentIds,
    hasIssue: invalidSegmentIds.length > 0 || duplicateSegmentIds.length > 0
  }
})

const selectedNodeKnowledgeMeta = computed(() => {
  const nodeId = selectedNodeId.value
  if (!nodeId) {
    return { difficulty: '', tagsText: '', prerequisitesText: '' }
  }
  const knowledge = parseKnowledgeNodes(selectedNodeKnowledgeNodesJson.value)
  const hit = knowledge.find(item => String(item?.node_id || '').trim() === nodeId) || {}
  const tags = Array.isArray(hit.tags) ? hit.tags : []
  const prerequisites = Array.isArray(hit.prerequisites) ? hit.prerequisites : []
  return {
    difficulty: String(hit.difficulty || '').trim(),
    tagsText: tags.map(item => String(item || '').trim()).filter(Boolean).join(','),
    prerequisitesText: prerequisites.map(item => String(item || '').trim()).filter(Boolean).join(',')
  }
})

const knowledgeNodesState = computed(() => validateJSONArrayText(selectedNodeKnowledgeNodesJson.value, '知识节点 JSON'))
const scriptSegmentsState = computed(() => validateJSONArrayText(selectedNodeScriptSegmentsJson.value, '讲稿段落映射 JSON'))

const structuredMarkdownPreview = computed(() => renderSimpleMarkdown(selectedNodeStructuredMarkdown.value))

const updateNode = (index, field, value, options = { syncScript: true }) => {
  localNodes.value = localNodes.value.map((node, currentIndex) => {
    if (currentIndex !== index) return node
    return {
      ...node,
      [field]: field === 'estimatedDuration' ? Number(value) || 0 : value
    }
  })
  emitNodes()
  if (options.syncScript) {
    syncScriptFromNodes()
  }
}

const addNode = () => {
  localNodes.value = [...localNodes.value, {
    id: '',
    nodeId: `p${props.currentEditPage}_n${localNodes.value.length + 1}`,
    schemaVersion: 2,
    type: inferNodeType(localNodes.value.length, localNodes.value.length + 1),
    title: `节点${localNodes.value.length + 1}：讲解`,
    summary: '',
    scriptText: '',
    reteachScript: '',
    transitionText: '',
    structuredMarkdown: '',
    knowledgeNodesJson: '[]',
    scriptSegmentsJson: '[]',
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
      schemaVersion: 2,
      type: inferNodeType(index, list.length),
      title: `${index === 0 ? '节点1：开场' : `节点${index + 1}：核心代码`}`,
      summary: text.length > 48 ? `${text.slice(0, 48)}...` : text,
      scriptText: text,
      reteachScript: '',
      transitionText: '',
      structuredMarkdown: '',
      knowledgeNodesJson: '[]',
      scriptSegmentsJson: '[]',
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
      schemaVersion: 2,
      type: 'opening',
      title: '节点1：开场',
      summary: '',
      scriptText: value,
      reteachScript: '',
      transitionText: '',
      structuredMarkdown: '',
      knowledgeNodesJson: '[]',
      scriptSegmentsJson: '[]',
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

const updateSelectedNodeField = (field, value) => {
  updateNode(selectedNodeIndex.value, field, value, { syncScript: false })
}

const formatSelectedNodeJSON = (field) => {
  const currentNode = timelineNodes.value[selectedNodeIndex.value]
  if (!currentNode) return
  const raw = String(currentNode[field] || '').trim()
  if (!raw) {
    updateNode(selectedNodeIndex.value, field, '[]', { syncScript: false })
    return
  }
  try {
    const parsed = JSON.parse(raw)
    const normalized = Array.isArray(parsed) ? parsed : [parsed]
    updateNode(selectedNodeIndex.value, field, JSON.stringify(normalized, null, 2), { syncScript: false })
  } catch {
    // 保留用户输入，避免误覆盖
  }
}

const resetSelectedNodeJSON = (field) => {
  updateNode(selectedNodeIndex.value, field, '[]', { syncScript: false })
}

function validateJSONArrayText(rawText, label) {
  const raw = String(rawText || '').trim()
  if (!raw) {
    return { valid: true, message: `${label} 为空，将按 [] 处理` }
  }
  try {
    const parsed = JSON.parse(raw)
    if (!Array.isArray(parsed)) {
      return { valid: false, message: `${label} 必须是 JSON 数组` }
    }
    return { valid: true, message: `${label} 格式正确（${parsed.length} 条）` }
  } catch {
    return { valid: false, message: `${label} 格式错误，请检查逗号/引号/括号` }
  }
}

function renderSimpleMarkdown(md) {
  const raw = String(md || '').trim()
  if (!raw) return ''
  const escaped = escapeHtml(raw)
  const lines = escaped.split('\n')
  const html = lines.map((line) => {
    const text = line.trim()
    if (!text) return '<div class="md-spacer"></div>'
    if (text.startsWith('### ')) return `<h6>${text.slice(4)}</h6>`
    if (text.startsWith('## ')) return `<h5>${text.slice(3)}</h5>`
    if (text.startsWith('# ')) return `<h4>${text.slice(2)}</h4>`
    if (text.startsWith('- ')) return `<p class="md-li">• ${text.slice(2)}</p>`
    return `<p>${text}</p>`
  })
  return html.join('')
}

function escapeHtml(text) {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
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

function setTimelineNodeRef(el, nodeId) {
  if (!nodeId) return
  if (el) {
    timelineNodeRefs.value[nodeId] = el
    return
  }
  delete timelineNodeRefs.value[nodeId]
}

function setNodeTabRef(el, nodeId) {
  if (!nodeId) return
  if (el) {
    nodeTabRefs.value[nodeId] = el
    return
  }
  delete nodeTabRefs.value[nodeId]
}

function setScriptBlockRef(el, segmentId) {
  if (!segmentId) return
  if (el) {
    scriptBlockRefs.value[segmentId] = el
    return
  }
  delete scriptBlockRefs.value[segmentId]
}

function setBlockTextareaRef(el, segmentId) {
  if (!segmentId) return
  if (el) {
    blockTextareaRefs.value[segmentId] = el
    return
  }
  delete blockTextareaRefs.value[segmentId]
}

function toggleBold() {
  if (!tiptapEditor.value) return
  tiptapEditor.value.chain().focus().toggleBold().run()
}

function toggleItalic() {
  if (!tiptapEditor.value) return
  tiptapEditor.value.chain().focus().toggleItalic().run()
}

function toggleBulletList() {
  if (!tiptapEditor.value) return
  tiptapEditor.value.chain().focus().toggleBulletList().run()
}

function splitCurrentNodeByCaret() {
  if (!tiptapEditor.value) return
  const plainText = tiptapEditor.value.getText({ blockSeparator: '\n\n' })
  if (!plainText.trim()) return
  const state = tiptapEditor.value.state
  const pos = state.selection?.from || 0
  const left = plainText.slice(0, pos).trim()
  const right = plainText.slice(pos).trim()
  if (!left || !right) return
  splitNodeAtSelectedIndex(left, right)
}

function handleSegmentDragStart(segmentId) {
  draggingSegmentId.value = String(segmentId || '').trim()
}

function handleSegmentDragEnd() {
  draggingSegmentId.value = ''
}

function handleDropBindToNode(nodeId) {
  const targetNodeId = String(nodeId || '').trim()
  const segmentId = String(draggingSegmentId.value || '').trim()
  if (!targetNodeId || !segmentId) return
  bindSegmentToNode(segmentId, targetNodeId)
  draggingSegmentId.value = ''
}

function bindSegmentToNode(segmentId, nodeId) {
  const targetSegmentId = String(segmentId || '').trim()
  const targetNodeId = String(nodeId || '').trim()
  if (!targetSegmentId || !targetNodeId) return

  const index = localNodes.value.findIndex(item => String(item?.nodeId || '').trim() === targetNodeId)
  if (index < 0) return

  const node = localNodes.value[index] || {}
  const segments = parseSegmentList(node.scriptSegmentsJson)
  const paragraph = scriptParagraphs.value.find(item => item.segmentId === targetSegmentId)
  const existing = segments.find(item => String(item?.segment_id || '').trim() === targetSegmentId)
  if (!existing) {
    segments.push({
      segment_id: targetSegmentId,
      node_ids: [targetNodeId],
      text: String(paragraph?.text || '').trim(),
      manual: true
    })
  } else {
    const ids = Array.isArray(existing.node_ids) ? existing.node_ids : []
    const normalized = Array.from(new Set(ids.map(item => String(item || '').trim()).filter(Boolean)))
    if (!normalized.includes(targetNodeId)) {
      normalized.push(targetNodeId)
    }
    existing.node_ids = normalized
    existing.text = existing.text || String(paragraph?.text || '').trim()
    existing.manual = true
  }

  updateNode(index, 'scriptSegmentsJson', JSON.stringify(segments, null, 2), { syncScript: false })
  if (!localDraftCoveredNodeIdSet.value.has(targetNodeId)) {
    localDraftCoveredNodeIds.value = [...localDraftCoveredNodeIds.value, targetNodeId]
  }
}

function splitScriptBlock(segmentId) {
  const targetSegmentId = String(segmentId || '').trim()
  const textarea = blockTextareaRefs.value[targetSegmentId]
  const original = scriptParagraphs.value.find(item => item.segmentId === targetSegmentId)
  if (!textarea || !original) return

  const raw = String(textarea.value || original.text || '')
  const caret = Number(textarea.selectionStart || 0)
  const left = raw.slice(0, caret).trim()
  const right = raw.slice(caret).trim()
  if (!left || !right) return

  const updated = []
  scriptParagraphs.value.forEach((segment) => {
    if (segment.segmentId !== targetSegmentId) {
      updated.push(segment)
      return
    }
    updated.push({ segmentId: segment.segmentId, text: left })
    updated.push({ segmentId: `${segment.segmentId}_split`, text: right })
  })
  rebuildScriptFromParagraphDraft(updated)
}

function rebuildScriptFromParagraphDraft(paragraphDraft) {
  const cleaned = (paragraphDraft || [])
    .map(item => ({ text: String(item?.text || '').trim() }))
    .filter(item => item.text)
    .map((item, idx) => ({ segmentId: `seg_${idx + 1}`, text: item.text }))

  const merged = cleaned.map(item => item.text).join('\n\n')
  emit('update:current-script', merged)

  localNodes.value = localNodes.value.map((node = {}) => {
    const source = parseSegmentList(node.scriptSegmentsJson)
    const normalized = source
      .map((segment = {}) => {
        const previousText = String(segment.text || '').trim()
        const matched = cleaned.find(item => item.text === previousText)
        if (!matched) return null
        return {
          ...segment,
          segment_id: matched.segmentId,
          text: matched.text,
          manual: true
        }
      })
      .filter(Boolean)
    return {
      ...node,
      scriptSegmentsJson: JSON.stringify(normalized, null, 2)
    }
  })
  emitNodes()
}

function handleOutlineDragStart(index) {
  draggingOutlineIndex.value = Number(index)
}

function handleOutlineDrop(targetIndex) {
  const from = Number(draggingOutlineIndex.value)
  const to = Number(targetIndex)
  draggingOutlineIndex.value = -1
  if (from < 0 || to < 0 || from === to || from >= localNodes.value.length || to >= localNodes.value.length) {
    return
  }
  const next = [...localNodes.value]
  const [moved] = next.splice(from, 1)
  next.splice(to, 0, moved)
  reorderNodesAndSync(next)
  selectedNodeIndex.value = to
}

function reorderNodesAndSync(nextNodes) {
  const ordered = Array.isArray(nextNodes) ? nextNodes : []
  const oldToNew = new Map()
  ordered.forEach((node = {}, idx) => {
    const oldId = String(node.nodeId || '').trim()
    const newId = `p${props.currentEditPage}_n${idx + 1}`
    if (oldId) {
      oldToNew.set(oldId, newId)
    }
  })

  localNodes.value = ordered.map((node = {}, idx, list) => {
    const nextNodeId = `p${props.currentEditPage}_n${idx + 1}`
    const scriptSegments = parseSegmentList(node.scriptSegmentsJson).map((segment = {}) => {
      const ids = Array.isArray(segment.node_ids) ? segment.node_ids : []
      const remappedIDs = Array.from(new Set(ids
        .map(item => String(item || '').trim())
        .filter(Boolean)
        .map(item => oldToNew.get(item) || item)))
      return {
        ...segment,
        node_ids: remappedIDs
      }
    })

    const knowledge = parseKnowledgeNodes(node.knowledgeNodesJson).map((entry = {}) => {
      const currentNodeID = String(entry.node_id || '').trim()
      const prereq = Array.isArray(entry.prerequisites) ? entry.prerequisites : []
      return {
        ...entry,
        node_id: oldToNew.get(currentNodeID) || currentNodeID,
        prerequisites: Array.from(new Set(prereq
          .map(item => String(item || '').trim())
          .filter(Boolean)
          .map(item => oldToNew.get(item) || item)))
      }
    })

    return {
      ...node,
      nodeId: nextNodeId,
      type: inferNodeType(idx, list.length),
      sortOrder: idx + 1,
      scriptSegmentsJson: JSON.stringify(scriptSegments, null, 2),
      knowledgeNodesJson: JSON.stringify(knowledge, null, 2)
    }
  })

  emitNodes()
  syncScriptFromNodes()
}

function splitNodeAtSelectedIndex(leftText, rightText) {
  const index = selectedNodeIndex.value
  if (index < 0) return
  const source = [...localNodes.value]
  const current = source[index]
  if (!current) return

  const newNode = {
    ...current,
    id: '',
    scriptText: rightText,
    summary: rightText.length > 48 ? `${rightText.slice(0, 48)}...` : rightText,
    scriptSegmentsJson: '[]',
    knowledgeNodesJson: '[]'
  }

  source[index] = {
    ...current,
    scriptText: leftText,
    summary: leftText.length > 48 ? `${leftText.slice(0, 48)}...` : leftText
  }
  source.splice(index + 1, 0, newNode)
  reorderNodesAndSync(source)
  selectedNodeIndex.value = index + 1
}

function focusNodeById(nodeId) {
  const targetId = String(nodeId || '').trim()
  if (!targetId) return
  const index = timelineNodes.value.findIndex(node => String(node.nodeId || '').trim() === targetId)
  if (index < 0) {
    pendingFocusNodeId.value = targetId
    return
  }
  selectedNodeIndex.value = index
  focusedNodeId.value = targetId
  if (focusFlashTimer) {
    window.clearTimeout(focusFlashTimer)
    focusFlashTimer = null
  }
  showStructuredEditors.value = true
  focusFlashTimer = window.setTimeout(() => {
    if (focusedNodeId.value === targetId) {
      focusedNodeId.value = ''
    }
  }, 2000)
  window.requestAnimationFrame(() => {
    const el = timelineNodeRefs.value[targetId]
    if (el && typeof el.scrollIntoView === 'function') {
      el.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'center' })
    }
    const tabEl = nodeTabRefs.value[targetId]
    if (tabEl && typeof tabEl.scrollIntoView === 'function') {
      tabEl.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'center' })
    }
  })
  nextTick(() => {
    if (structuredEditorsRef.value && typeof structuredEditorsRef.value.scrollIntoView === 'function') {
      structuredEditorsRef.value.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
    }
    if (structuredMarkdownEditorRef.value && typeof structuredMarkdownEditorRef.value.focus === 'function') {
      structuredMarkdownEditorRef.value.focus()
      return
    }
    if (linkedEditorRef.value && typeof linkedEditorRef.value.focus === 'function') {
      linkedEditorRef.value.focus()
    }
    focusNodeRelatedBlock(targetId)
  })
}

function focusNodeRelatedBlock(nodeId) {
  const targetNodeId = String(nodeId || '').trim()
  if (!targetNodeId) return
  const firstMatched = scriptParagraphs.value.find((segment) => {
    const nodeIDs = segmentNodeIdMap.value.get(segment.segmentId) || []
    return nodeIDs.includes(targetNodeId)
  })
  if (!firstMatched) return
  activeBlockSegmentId.value = firstMatched.segmentId
  nextTick(() => {
    const targetEl = scriptBlockRefs.value[firstMatched.segmentId]
    if (targetEl && typeof targetEl.scrollIntoView === 'function') {
      targetEl.scrollIntoView({ behavior: 'smooth', block: 'center' })
    }
  })
}

function handleScriptBlockClick(segmentId) {
  const targetSegmentId = String(segmentId || '').trim()
  if (!targetSegmentId) return
  activeBlockSegmentId.value = targetSegmentId
  const nodeIDs = segmentNodeIdMap.value.get(targetSegmentId) || []
  if (!nodeIDs.length) return
  focusNodeById(nodeIDs[0])
}

function updateScriptBlockText(segmentId, value) {
  const targetSegmentId = String(segmentId || '').trim()
  if (!targetSegmentId) return

  const updatedParagraphs = scriptParagraphs.value.map((segment) => {
    if (segment.segmentId !== targetSegmentId) return segment
    return {
      ...segment,
      text: String(value || '')
    }
  })

  const merged = updatedParagraphs
    .map(item => String(item.text || '').trim())
    .filter(Boolean)
    .join('\n\n')
  emit('update:current-script', merged)

  localNodes.value = localNodes.value.map((node = {}) => {
    const segments = parseSegmentList(node.scriptSegmentsJson)
    if (!segments.length) return node
    const patched = segments.map((segment = {}) => {
      const currentSegmentId = String(segment.segment_id || '').trim()
      if (currentSegmentId !== targetSegmentId) return segment
      return {
        ...segment,
        text: String(value || '').trim()
      }
    })
    return {
      ...node,
      scriptSegmentsJson: JSON.stringify(patched, null, 2)
    }
  })
  emitNodes()
}

function jumpToNextUncoveredNode() {
  if (!hasUncoveredNodes.value || timelineNodes.value.length === 0) return
  const total = timelineNodes.value.length
  const start = Math.max(0, Number(selectedNodeIndex.value) || 0)
  for (let offset = 1; offset <= total; offset += 1) {
    const idx = (start + offset) % total
    const candidateId = String(timelineNodes.value[idx]?.nodeId || '').trim()
    if (candidateId && effectiveUncoveredNodeIdSet.value.has(candidateId)) {
      focusNodeById(candidateId)
      return
    }
  }
}

function jumpToPriorityUncoveredNode() {
  if (!hasUncoveredNodes.value || timelineNodes.value.length === 0) return
  const currentIndex = Math.max(0, Number(selectedNodeIndex.value) || 0)
  const candidates = timelineNodes.value
    .map((node, idx) => {
      const nodeId = String(node?.nodeId || '').trim()
      return {
        idx,
        nodeId,
        score: nodePriorityScoreMap.value.get(nodeId) || 0,
        distance: idx >= currentIndex ? idx - currentIndex : timelineNodes.value.length - currentIndex + idx
      }
    })
    .filter(item => item.nodeId && effectiveUncoveredNodeIdSet.value.has(item.nodeId))

  if (candidates.length === 0) return

  candidates.sort((a, b) => {
    if (b.score !== a.score) return b.score - a.score
    return a.distance - b.distance
  })

  focusNodeById(candidates[0].nodeId)
}

function fillSelectedNodeMappingDraft() {
  const nodeId = selectedNodeId.value
  if (!nodeId) return
  const segments = parseSegmentList(selectedNodeScriptSegmentsJson.value)
  const exists = segments.some(item => Array.isArray(item?.node_ids) && item.node_ids.map(v => String(v || '').trim()).includes(nodeId))
  if (exists) return
  const fallbackSegmentId = `seg_${selectedNodeIndex.value + 1}`
  segments.push({
    segment_id: fallbackSegmentId,
    node_ids: [nodeId],
    manual: true
  })
  updateSelectedNodeField('scriptSegmentsJson', JSON.stringify(segments, null, 2))
  if (!localDraftCoveredNodeIdSet.value.has(nodeId)) {
    localDraftCoveredNodeIds.value = [...localDraftCoveredNodeIds.value, nodeId]
  }
}

function toggleSegmentBinding(segmentId, segmentText) {
  const nodeId = selectedNodeId.value
  const targetSegmentId = String(segmentId || '').trim()
  if (!nodeId || !targetSegmentId) return

  const segments = parseSegmentList(selectedNodeScriptSegmentsJson.value)
  const idx = segments.findIndex(item => String(item?.segment_id || '').trim() === targetSegmentId)

  if (idx < 0) {
    segments.push({
      segment_id: targetSegmentId,
      node_ids: [nodeId],
      text: String(segmentText || '').trim(),
      manual: true
    })
  } else {
    const current = segments[idx] || {}
    const normalizedIds = Array.from(new Set((Array.isArray(current.node_ids) ? current.node_ids : []).map(item => String(item || '').trim()).filter(Boolean)))
    const hitIndex = normalizedIds.indexOf(nodeId)
    if (hitIndex >= 0) {
      normalizedIds.splice(hitIndex, 1)
    } else {
      normalizedIds.push(nodeId)
    }
    if (normalizedIds.length === 0) {
      segments.splice(idx, 1)
    } else {
      segments[idx] = {
        ...current,
        segment_id: targetSegmentId,
        node_ids: normalizedIds,
        manual: true
      }
    }
  }

  updateSelectedNodeField('scriptSegmentsJson', JSON.stringify(segments, null, 2))
  const stillMapped = segments.some(item => Array.isArray(item?.node_ids) && item.node_ids.map(v => String(v || '').trim()).includes(nodeId))
  if (stillMapped) {
    if (!localDraftCoveredNodeIdSet.value.has(nodeId)) {
      localDraftCoveredNodeIds.value = [...localDraftCoveredNodeIds.value, nodeId]
    }
    return
  }
  localDraftCoveredNodeIds.value = localDraftCoveredNodeIds.value.filter(item => String(item || '').trim() !== nodeId)
}

function bindAllSegmentsToSelectedNode() {
  const nodeId = selectedNodeId.value
  if (!nodeId || scriptParagraphs.value.length === 0) return
  const segments = parseSegmentList(selectedNodeScriptSegmentsJson.value)

  scriptParagraphs.value.forEach((paragraph) => {
    const targetSegmentId = String(paragraph.segmentId || '').trim()
    if (!targetSegmentId) return
    const idx = segments.findIndex(item => String(item?.segment_id || '').trim() === targetSegmentId)
    if (idx < 0) {
      segments.push({
        segment_id: targetSegmentId,
        node_ids: [nodeId],
        text: String(paragraph.text || '').trim(),
        manual: true
      })
      return
    }
    const current = segments[idx] || {}
    const ids = Array.isArray(current.node_ids) ? current.node_ids : []
    const normalized = Array.from(new Set(ids.map(item => String(item || '').trim()).filter(Boolean)))
    if (!normalized.includes(nodeId)) {
      normalized.push(nodeId)
    }
    segments[idx] = {
      ...current,
      segment_id: targetSegmentId,
      node_ids: normalized,
      text: current.text || String(paragraph.text || '').trim(),
      manual: true
    }
  })

  updateSelectedNodeField('scriptSegmentsJson', JSON.stringify(segments, null, 2))
  if (!localDraftCoveredNodeIdSet.value.has(nodeId)) {
    localDraftCoveredNodeIds.value = [...localDraftCoveredNodeIds.value, nodeId]
  }
}

function clearSelectedNodeSegments() {
  const nodeId = selectedNodeId.value
  if (!nodeId) return
  const segments = parseSegmentList(selectedNodeScriptSegmentsJson.value)

  const cleaned = segments
    .map((item = {}) => {
      const ids = Array.isArray(item.node_ids) ? item.node_ids : []
      const normalized = ids
        .map(v => String(v || '').trim())
        .filter(v => v && v !== nodeId)
      if (normalized.length === 0) {
        return null
      }
      return {
        ...item,
        node_ids: Array.from(new Set(normalized)),
        manual: true
      }
    })
    .filter(Boolean)

  updateSelectedNodeField('scriptSegmentsJson', JSON.stringify(cleaned, null, 2))
  localDraftCoveredNodeIds.value = localDraftCoveredNodeIds.value.filter(item => String(item || '').trim() !== nodeId)
}

function bindUnmappedSegmentsToSelectedNode() {
  const nodeId = selectedNodeId.value
  if (!nodeId) return
  const targets = segmentCoverageSummary.value.uncoveredSegmentIds || []
  if (targets.length === 0) return
  const segments = parseSegmentList(selectedNodeScriptSegmentsJson.value)

  targets.forEach((segmentId) => {
    const targetId = String(segmentId || '').trim()
    if (!targetId) return
    const paragraph = scriptParagraphs.value.find(item => item.segmentId === targetId)
    const idx = segments.findIndex(item => String(item?.segment_id || '').trim() === targetId)
    if (idx < 0) {
      segments.push({
        segment_id: targetId,
        node_ids: [nodeId],
        text: String(paragraph?.text || '').trim(),
        manual: true
      })
      return
    }
    const current = segments[idx] || {}
    const ids = Array.isArray(current.node_ids) ? current.node_ids : []
    const normalized = Array.from(new Set(ids.map(item => String(item || '').trim()).filter(Boolean)))
    if (!normalized.includes(nodeId)) {
      normalized.push(nodeId)
    }
    segments[idx] = {
      ...current,
      segment_id: targetId,
      node_ids: normalized,
      text: current.text || String(paragraph?.text || '').trim(),
      manual: true
    }
  })

  updateSelectedNodeField('scriptSegmentsJson', JSON.stringify(segments, null, 2))
  if (!localDraftCoveredNodeIdSet.value.has(nodeId)) {
    localDraftCoveredNodeIds.value = [...localDraftCoveredNodeIds.value, nodeId]
  }
}

function rebuildSegmentIdsFromScript() {
  if (scriptParagraphs.value.length === 0) return
  const segments = parseSegmentList(selectedNodeScriptSegmentsJson.value)
  const normalized = scriptParagraphs.value.map((paragraph) => {
    const existed = segments.find(item => String(item?.segment_id || '').trim() === String(paragraph.segmentId || '').trim())
    if (!existed) {
      return {
        segment_id: paragraph.segmentId,
        node_ids: [],
        text: String(paragraph.text || '').trim(),
        manual: true
      }
    }
    return {
      ...existed,
      segment_id: paragraph.segmentId,
      text: String(paragraph.text || '').trim(),
      node_ids: Array.isArray(existed.node_ids)
        ? Array.from(new Set(existed.node_ids.map(v => String(v || '').trim()).filter(Boolean)))
        : [],
      manual: true
    }
  })
  updateSelectedNodeField('scriptSegmentsJson', JSON.stringify(normalized, null, 2))
}

function repairSelectedNodeSegments() {
  const validSegmentIdSet = new Set(scriptParagraphs.value.map(item => item.segmentId))
  const paragraphTextById = new Map(scriptParagraphs.value.map(item => [item.segmentId, String(item.text || '').trim()]))
  const source = parseSegmentList(selectedNodeScriptSegmentsJson.value)
  const mergedById = new Map()

  source.forEach((item = {}) => {
    const segmentId = String(item.segment_id || '').trim()
    if (!segmentId || !validSegmentIdSet.has(segmentId)) {
      return
    }
    const existing = mergedById.get(segmentId) || {
      segment_id: segmentId,
      node_ids: [],
      text: paragraphTextById.get(segmentId) || '',
      manual: true
    }
    const ids = Array.isArray(item.node_ids) ? item.node_ids : []
    const mergedNodeIDs = new Set([...(existing.node_ids || []), ...ids])
    existing.node_ids = Array.from(mergedNodeIDs)
      .map(v => String(v || '').trim())
      .filter(Boolean)
    existing.text = paragraphTextById.get(segmentId) || String(existing.text || '').trim()
    existing.manual = true
    mergedById.set(segmentId, existing)
  })

  const repaired = Array.from(mergedById.values())
  updateSelectedNodeField('scriptSegmentsJson', JSON.stringify(repaired, null, 2))
}

function parseSegmentList(rawValue) {
  const raw = String(rawValue || '').trim()
  if (!raw) return []
  try {
    const parsed = JSON.parse(raw)
    return Array.isArray(parsed) ? parsed : [parsed]
  } catch {
    return []
  }
}

function updateKnowledgeMetaField(field, rawValue) {
  const nodeId = selectedNodeId.value
  if (!nodeId) return
  const knowledge = parseKnowledgeNodes(selectedNodeKnowledgeNodesJson.value)
  const idx = knowledge.findIndex(item => String(item?.node_id || '').trim() === nodeId)
  const fallback = {
    node_id: nodeId,
    title: timelineNodes.value[selectedNodeIndex.value]?.title || nodeId,
    level: 'point',
    difficulty: '',
    tags: [],
    prerequisites: []
  }
  const current = idx >= 0 ? { ...knowledge[idx] } : fallback

  if (field === 'difficulty') {
    current.difficulty = String(rawValue || '').trim()
  }
  if (field === 'tagsText') {
    current.tags = String(rawValue || '')
      .split(',')
      .map(item => item.trim())
      .filter(Boolean)
  }
  if (field === 'prerequisitesText') {
    current.prerequisites = String(rawValue || '')
      .split(',')
      .map(item => item.trim())
      .filter(Boolean)
  }

  if (idx >= 0) {
    knowledge[idx] = current
  } else {
    knowledge.push(current)
  }

  updateSelectedNodeField('knowledgeNodesJson', JSON.stringify(knowledge, null, 2))
}

function autoFillPrerequisite() {
  const nodeId = selectedNodeId.value
  if (!nodeId) return
  const currentIndex = timelineNodes.value.findIndex(node => String(node?.nodeId || '').trim() === nodeId)
  if (currentIndex < 0) return

  const previousNodeId = currentIndex > 0
    ? String(timelineNodes.value[currentIndex - 1]?.nodeId || '').trim()
    : ''

  const knowledge = parseKnowledgeNodes(selectedNodeKnowledgeNodesJson.value)
  const idx = knowledge.findIndex(item => String(item?.node_id || '').trim() === nodeId)
  const fallback = {
    node_id: nodeId,
    title: timelineNodes.value[selectedNodeIndex.value]?.title || nodeId,
    level: 'point',
    difficulty: '',
    tags: [],
    prerequisites: []
  }
  const current = idx >= 0 ? { ...knowledge[idx] } : fallback

  if (previousNodeId) {
    current.prerequisites = [previousNodeId]
  } else {
    current.prerequisites = []
  }

  if (idx >= 0) {
    knowledge[idx] = current
  } else {
    knowledge.push(current)
  }

  updateSelectedNodeField('knowledgeNodesJson', JSON.stringify(knowledge, null, 2))
}

function autoFillPrerequisitesForPage() {
  if (!localNodes.value.length) return

  localNodes.value = localNodes.value.map((node = {}, index) => {
    const nodeId = String(node.nodeId || '').trim()
    if (!nodeId) return node

    const previousNodeId = index > 0
      ? String(localNodes.value[index - 1]?.nodeId || '').trim()
      : ''

    const knowledge = parseKnowledgeNodes(node.knowledgeNodesJson)
    const entryIndex = knowledge.findIndex(item => String(item?.node_id || '').trim() === nodeId)
    const fallback = {
      node_id: nodeId,
      title: node.title || nodeId,
      level: 'point',
      difficulty: '',
      tags: [],
      prerequisites: []
    }
    const current = entryIndex >= 0 ? { ...knowledge[entryIndex] } : fallback
    current.prerequisites = previousNodeId ? [previousNodeId] : []

    if (entryIndex >= 0) {
      knowledge[entryIndex] = current
    } else {
      knowledge.push(current)
    }

    return {
      ...node,
      knowledgeNodesJson: JSON.stringify(knowledge, null, 2)
    }
  })

  emitNodes()
}

function parseKnowledgeNodes(rawValue) {
  const raw = String(rawValue || '').trim()
  if (!raw) return []
  try {
    const parsed = JSON.parse(raw)
    return Array.isArray(parsed) ? parsed : [parsed]
  } catch {
    return []
  }
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

.timeline-node.uncovered .node-bubble {
  border-style: dashed;
  border-color: #c0841a;
}

.timeline-node.focused .node-bubble {
  box-shadow: 0 0 0 2px #2f605a inset, 0 8px 18px rgba(47, 96, 90, 0.18);
}

.timeline-node.focused .node-dot {
  background: #2f605a;
  border-color: #e3f0eb;
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

.coverage-badge {
  display: inline-flex;
  align-items: center;
  border: 1px solid #fed7aa;
  background: #fff7ed;
  color: #9a3412;
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 12px;
}

.coverage-badge.neutral {
  border-color: #cbd5e1;
  background: #f8fafc;
  color: #334155;
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

.ai-inline input {
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

.node-tabs button.focused {
  box-shadow: 0 0 0 2px #2f605a inset;
}

.node-tabs button.uncovered {
  border-style: dashed;
  border-color: #c0841a;
}

.linked-editor {
  flex: 0 0 auto;
  min-height: 180px;
  margin-top: 8px;
  border: 1px solid #d9e4de;
  border-radius: 10px;
  background: #ffffff;
  overflow: hidden;
}

.editor-toolbar {
  display: flex;
  gap: 6px;
  padding: 6px;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.tiptap-editor {
  min-height: 130px;
  max-height: 220px;
  overflow-y: auto;
}

:deep(.tiptap-editor-body) {
  min-height: 130px;
  padding: 8px 10px;
  outline: none;
  color: #334155;
  font-size: 13px;
  line-height: 1.65;
}

:deep(.tiptap-editor-body p) {
  margin: 0 0 8px;
}

:deep(.tiptap-editor-body p:last-child) {
  margin-bottom: 0;
}

.outline-board {
  margin-top: 10px;
  border: 1px solid #dbe7e1;
  border-radius: 10px;
  background: #f9fcfa;
  padding: 8px;
}

.outline-list {
  display: grid;
  gap: 6px;
  max-height: 170px;
  overflow-y: auto;
}

.outline-item {
  border: 1px solid #dbe7e1;
  border-radius: 8px;
  padding: 6px 8px;
  background: #fff;
  display: flex;
  flex-direction: column;
  gap: 3px;
  cursor: move;
}

.outline-item.active {
  border-color: #2f605a;
  background: #eef6f2;
}

.outline-item strong {
  color: #355048;
  font-size: 12px;
}

.outline-item small {
  color: #64748b;
  font-size: 11px;
}

.structured-toggle-row {
  margin-top: 8px;
}

.block-editor {
  margin-top: 10px;
  border: 1px solid #dbe7e1;
  border-radius: 10px;
  background: #f9fcfa;
  padding: 8px;
}

.block-list {
  display: grid;
  gap: 8px;
  max-height: 260px;
  overflow-y: auto;
}

.block-item {
  border: 1px solid #dbe7e1;
  background: #ffffff;
  border-radius: 8px;
  padding: 6px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.block-item.active {
  border-color: #2f605a;
  box-shadow: 0 0 0 1px #2f605a inset;
}

.block-item.related {
  background: #f1f8f5;
}

.block-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.block-actions {
  display: flex;
  justify-content: flex-end;
}

.block-head strong {
  color: #355048;
  font-size: 12px;
}

.block-head span {
  color: #64748b;
  font-size: 12px;
}

.block-input {
  min-height: 54px;
  border: 1px solid #d9e4de;
  border-radius: 8px;
  padding: 6px 8px;
  font-size: 12px;
  line-height: 1.5;
  resize: vertical;
  font-family: inherit;
}

.mapping-editor {
  margin-top: 10px;
  border: 1px solid #dbe7e1;
  border-radius: 10px;
  background: #f9fcfa;
  padding: 8px;
}

.mapping-head {
  display: flex;
  justify-content: space-between;
  color: #4b635a;
  font-size: 12px;
  margin-bottom: 8px;
}

.mapping-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 8px;
}

.mapping-note {
  font-size: 12px;
  color: #b45309;
  margin-bottom: 8px;
}

.mapping-note.warn {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
  color: #9a3412;
}

.knowledge-meta-editor {
  margin-top: 10px;
  border: 1px solid #dbe7e1;
  border-radius: 10px;
  background: #f9fcfa;
  padding: 8px;
}

.meta-grid {
  display: grid;
  gap: 8px;
}

.meta-grid label {
  display: flex;
  flex-direction: column;
  gap: 4px;
  color: #4b635a;
  font-size: 12px;
}

.meta-grid input,
.meta-grid select {
  border: 1px solid #d9e4de;
  border-radius: 8px;
  padding: 6px 8px;
  font-size: 12px;
  background: #fff;
}

.segment-list {
  display: grid;
  gap: 6px;
  max-height: 180px;
  overflow-y: auto;
}

.segment-item {
  border: 1px solid #dbe7e1;
  background: #ffffff;
  border-radius: 8px;
  text-align: left;
  padding: 6px 8px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.segment-item strong {
  color: #355048;
  font-size: 12px;
}

.segment-item span {
  color: #64748b;
  font-size: 12px;
  line-height: 1.4;
}

.segment-item.active {
  border-color: #2f605a;
  background: #ecf6f2;
}

.coverage-hint {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  border: 1px solid #fed7aa;
  background: #fff7ed;
  color: #9a3412;
  border-radius: 9px;
  padding: 8px 10px;
  font-size: 12px;
}

.coverage-hint.success {
  border-color: #bbf7d0;
  background: #f0fdf4;
  color: #166534;
}

.structured-editors {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
  margin-top: 10px;
}

.structured-editors label {
  display: flex;
  flex-direction: column;
  gap: 6px;
  color: #587069;
  font-size: 12px;
}

.mini-editor {
  min-height: 72px;
  border: 1px solid #d9e4de;
  border-radius: 9px;
  padding: 8px 10px;
  resize: vertical;
  font-size: 12px;
  line-height: 1.5;
  background: #fbfdfc;
  font-family: inherit;
}

.json-tip {
  font-size: 12px;
  line-height: 1.4;
}

.json-actions {
  display: flex;
  gap: 6px;
  margin-top: 4px;
}

.tiny-btn {
  border: 1px solid #d0dfd7;
  background: #edf5f1;
  color: #3c524b;
  border-radius: 8px;
  padding: 3px 8px;
  font-size: 12px;
  cursor: pointer;
}

.tiny-btn.danger {
  border-color: #e0caca;
  background: #f8efef;
  color: #8d4f4f;
}

.json-tip.ok {
  color: #3f6f59;
}

.json-tip.error {
  color: #a44a4a;
}

.md-preview {
  border: 1px dashed #cadbd2;
  border-radius: 9px;
  background: #f7fbf9;
  padding: 8px;
}

.preview-title {
  font-size: 12px;
  color: #5c756c;
  margin-bottom: 6px;
}

.preview-body {
  font-size: 12px;
  color: #324842;
  line-height: 1.6;
}

.preview-body h4,
.preview-body h5,
.preview-body h6 {
  margin: 4px 0;
}

.preview-body p {
  margin: 2px 0;
}

.preview-body .md-li {
  padding-left: 4px;
}

.preview-body .md-spacer {
  height: 6px;
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