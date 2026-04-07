<template>
  <div class="knowledge-workbench">
    <header class="workbench-head">
      <div>
        <p class="eyebrow">个人知识拆解</p>
        <h3>构建你的专属知识体系</h3>
        <p class="head-tip">上传复习资料后自动拆解为知识树，支持问答、笔记、习题和复习计划演示。</p>
      </div>
      <div class="head-actions">
        <el-button plain @click="showHistory = true">我的拆解历史</el-button>
        <el-button plain @click="emit('reset-current')">新建拆解</el-button>
        <el-button type="danger" plain @click="clearCurrent">清空当前</el-button>
      </div>
    </header>

    <section v-if="!hasKnowledge" class="init-state">
      <div class="upload-card">
        <el-upload
          drag
          action="#"
          :auto-upload="false"
          :on-change="handleChange"
          accept=".pdf,.ppt,.pptx,.doc,.docx,.png,.jpg,.jpeg"
          :limit="1"
        >
          <div class="el-upload__text">拖拽文件到这里，或点击上传</div>
          <p class="upload-sub">支持 PDF / PPT / Word / 图片，演示模式仅做前端模拟</p>
        </el-upload>

        <div v-if="uploadedFile" class="file-brief">
          <strong>{{ uploadedFile.name }}</strong>
          <span>{{ fileSizeLabel }}</span>
        </div>

        <div class="parse-options">
          <div class="option-row">
            <span>拆解粒度</span>
            <el-radio-group v-model="parseOptions.granularity" size="small">
              <el-radio-button label="fine">精细拆解</el-radio-button>
              <el-radio-button label="coarse">粗略拆解</el-radio-button>
            </el-radio-group>
          </div>
          <div class="option-row">
            <span>自动生成讲稿</span>
            <el-switch v-model="parseOptions.autoScript" />
          </div>
          <div class="option-row">
            <span>自动生成习题</span>
            <el-switch v-model="parseOptions.autoQuiz" />
          </div>
        </div>

        <el-button class="parse-btn" type="primary" :disabled="!uploadedFile || isParsing" @click="emit('parse-knowledge')">
          开始拆解知识点
        </el-button>

        <el-steps v-if="isParsing || parseResult" :active="parseStepActive" finish-status="success" simple class="parse-steps">
          <el-step title="文件解析" />
          <el-step title="知识拆分" />
          <el-step title="知识树生成" />
          <el-step title="内容优化" />
        </el-steps>

        <el-alert v-if="isParsing" title="正在拆解中，请稍候..." type="info" show-icon :closable="false" />
        <el-alert v-if="parseResult" :title="parseResult" type="success" show-icon />
      </div>

      <div class="feature-grid">
        <div class="feature-item">AI 自动拆解知识点并生成结构化知识树</div>
        <div class="feature-item">支持知识点编辑、问答、习题与笔记联动</div>
        <div class="feature-item">演示数据可直接驱动完整学习闭环</div>
      </div>

      <div class="history-mini">
        <div class="block-head">
          <h4>历史拆解记录</h4>
          <el-button text @click="showHistory = true">查看全部</el-button>
        </div>
        <div v-if="historyRecords.length" class="history-list">
          <article v-for="item in historyRecords.slice(0, 4)" :key="item.id" class="history-card">
            <strong>{{ item.fileName }}</strong>
            <p>{{ item.timeLabel }}</p>
            <span>{{ item.chapterCount }} 章 · {{ item.pointCount }} 点</span>
          </article>
        </div>
        <el-empty v-else description="你还没有拆解过讲义，上传文件开始第一次拆解吧" />
      </div>
    </section>

    <section v-else class="parsed-state">
      <aside class="left-pane">
        <div class="left-tools">
          <el-switch v-model="editMode" inline-prompt active-text="编辑" inactive-text="查看" />
          <el-button size="small" @click="expandAll">展开</el-button>
          <el-button size="small" @click="collapseAll">折叠</el-button>
          <el-input v-model="searchKeyword" size="small" placeholder="搜索知识点" clearable />
          <el-button size="small" type="primary" @click="saveTreeEdits">保存修改</el-button>
        </div>

        <el-tree
          ref="treeRef"
          class="knowledge-tree"
          :data="knowledgeList"
          :props="treeProps"
          node-key="id"
          :default-expand-all="treeExpanded"
          :highlight-current="true"
          :current-node-key="currentNodeId"
          :filter-node-method="filterNode"
          :draggable="editMode"
          @node-click="handleNodeClick"
        >
          <template #default="{ data }">
            <div class="tree-node-row">
              <span class="node-title">{{ data.name }}</span>
              <div class="node-badges">
                <el-tag size="small" :type="masteryTagType(data.id)">{{ masteryLabel(data.id) }}</el-tag>
                <span class="node-counter">📝{{ noteCount(data.id) }}</span>
                <span class="node-counter">✏️{{ questionCount(data.id) }}</span>
              </div>
            </div>
          </template>
        </el-tree>

        <div class="left-foot">
          <span>章节：{{ knowledgeList.length }}</span>
          <span>知识点：{{ totalNodes }}</span>
          <span>更新时间：{{ new Date().toLocaleTimeString() }}</span>
        </div>
      </aside>

      <section class="right-pane">
        <div class="detail-head">
          <div>
            <p class="eyebrow">当前知识点</p>
            <h4>{{ currentNodeName || '请选择知识点' }}</h4>
          </div>
          <div class="detail-actions">
            <el-tag>{{ masteryLabel(currentNodeId) }}</el-tag>
            <el-button size="small" @click="markMastered">标记已掌握</el-button>
          </div>
        </div>

        <el-tabs v-model="activeTab">
          <el-tab-pane label="知识点详情" name="detail">
            <el-input v-model="detailDraft" type="textarea" :rows="7" />
            <div class="inline-actions">
              <el-button @click="toggleFavorite">{{ currentMeta.favorite ? '取消收藏' : '收藏知识点' }}</el-button>
              <el-button @click="addToReview">添加到复习计划</el-button>
              <el-button type="primary" @click="regenPractice">生成配套习题</el-button>
              <el-button type="success" @click="saveDetail">保存内容</el-button>
            </div>
          </el-tab-pane>

          <el-tab-pane label="AI问答" name="qa">
            <div class="qa-list" v-if="currentMeta.qa.length">
              <article v-for="(item, idx) in currentMeta.qa" :key="`${item.q}-${idx}`" class="qa-item">
                <p><strong>Q:</strong> {{ item.q }}</p>
                <p><strong>A:</strong> {{ item.a }}</p>
              </article>
            </div>
            <el-empty v-else description="还没有问答记录" />
            <div class="qa-input-row">
              <el-input v-model="qaInput" placeholder="输入问题，默认绑定当前知识点" @keyup.enter="sendMockQa" />
              <el-button type="primary" @click="sendMockQa">发送</el-button>
            </div>
          </el-tab-pane>

          <el-tab-pane label="我的笔记" name="notes">
            <el-input v-model="noteInput" type="textarea" :rows="4" placeholder="记录你的理解、公式或例题总结" />
            <div class="inline-actions">
              <el-button type="primary" @click="saveNote">保存笔记</el-button>
            </div>
            <div class="note-list" v-if="currentMeta.notes.length">
              <article v-for="item in currentMeta.notes" :key="item.id" class="note-item">
                <strong>{{ item.time }}</strong>
                <p>{{ item.content }}</p>
              </article>
            </div>
            <el-empty v-else description="当前知识点暂无笔记" />
          </el-tab-pane>

          <el-tab-pane label="配套习题" name="practice">
            <div class="practice-list" v-if="currentMeta.questions.length">
              <article v-for="q in currentMeta.questions" :key="q.id" class="practice-item">
                <h5>{{ q.stem }}</h5>
                <el-radio-group v-model="q.userAnswer" @change="submitQuestion(q)">
                  <el-radio v-for="opt in q.options" :key="opt" :label="opt">{{ opt }}</el-radio>
                </el-radio-group>
                <p v-if="q.result" :class="['q-result', q.result.ok ? 'ok' : 'bad']">{{ q.result.text }}</p>
              </article>
            </div>
            <el-empty v-else description="当前知识点暂无习题，可点击“生成配套习题”" />
          </el-tab-pane>
        </el-tabs>

        <div class="data-strip">
          <el-progress :percentage="masteryProgress" :stroke-width="14" />
          <p>已掌握 {{ masteredCount }} / {{ totalNodes }} 个知识点，建议优先复习“未掌握”节点。</p>
        </div>
      </section>
    </section>

    <div v-if="hasKnowledge" class="floating-actions">
      <el-button type="primary" circle @click="globalAsk">问</el-button>
      <el-button circle @click="generateReviewPack">包</el-button>
      <el-button circle @click="exportMockDoc">导</el-button>
      <el-button circle @click="syncToCourse">联</el-button>
    </div>

    <el-dialog v-model="showHistory" title="我的拆解历史" width="720px">
      <div v-if="historyRecords.length" class="history-dialog-list">
        <article v-for="item in historyRecords" :key="item.id" class="history-dialog-card">
          <div>
            <strong>{{ item.fileName }}</strong>
            <p>{{ item.timeLabel }}</p>
          </div>
          <div>
            <span>{{ item.chapterCount }} 章 / {{ item.pointCount }} 点</span>
          </div>
        </article>
      </div>
      <el-empty v-else description="暂无历史记录" />
    </el-dialog>
  </div>
</template>

<script setup>
/* eslint-disable no-undef */
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'

const emit = defineEmits(['file-change', 'parse-knowledge', 'node-click', 'reset-current'])

const props = defineProps({
  uploadedFile: {
    type: [Object, File],
    default: null
  },
  isParsing: {
    type: Boolean,
    default: false
  },
  parseResult: {
    type: String,
    default: ''
  },
  knowledgeList: {
    type: Array,
    default: () => []
  },
  treeProps: {
    type: Object,
    default: () => ({})
  }
})

const hasKnowledge = computed(() => Array.isArray(props.knowledgeList) && props.knowledgeList.length > 0)
const fileSizeLabel = computed(() => {
  const size = Number(props.uploadedFile?.size || 0)
  if (!size) return '未知大小'
  if (size > 1024 * 1024) return `${(size / 1024 / 1024).toFixed(2)} MB`
  return `${Math.round(size / 1024)} KB`
})

const parseOptions = reactive({
  granularity: 'fine',
  autoScript: true,
  autoQuiz: true
})

const parseStepActive = ref(0)
let parseStepTimer = null

const treeRef = ref(null)
const treeExpanded = ref(true)
const searchKeyword = ref('')
const editMode = ref(false)
const activeTab = ref('detail')
const currentNodeId = ref('')
const currentNodeName = ref('')
const detailDraft = ref('')
const qaInput = ref('')
const noteInput = ref('')
const showHistory = ref(false)

const nodeStateMap = reactive({})
const historyRecords = ref([])
const HISTORY_KEY = 'fuww_student_knowledge_history'

const totalNodes = computed(() => countNodes(props.knowledgeList))
const masteredCount = computed(() => {
  const ids = Object.keys(nodeStateMap)
  return ids.filter(id => nodeStateMap[id]?.mastery === 'mastered').length
})
const masteryProgress = computed(() => {
  if (!totalNodes.value) return 0
  return Math.round((masteredCount.value / totalNodes.value) * 100)
})

const currentMeta = computed(() => {
  if (!currentNodeId.value) {
    return {
      detail: '',
      favorite: false,
      mastery: 'todo',
      qa: [],
      notes: [],
      questions: []
    }
  }
  return ensureNodeState({ id: currentNodeId.value, name: currentNodeName.value })
})

const handleChange = (file) => {
  parseStepActive.value = 0
  emit('file-change', file)
}

const clearCurrent = () => {
  emit('reset-current')
  ElMessage.success('已清空当前拆解内容')
}

const hashSeed = (text) => {
  const source = String(text || '')
  let hash = 0
  for (let i = 0; i < source.length; i++) {
    hash = ((hash << 5) - hash) + source.charCodeAt(i)
    hash |= 0
  }
  return Math.abs(hash)
}

const ensureNodeState = (node) => {
  if (!node?.id) {
    return {
      detail: '',
      favorite: false,
      mastery: 'todo',
      qa: [],
      notes: [],
      questions: []
    }
  }
  if (!nodeStateMap[node.id]) {
    const seed = hashSeed(node.name || node.id)
    const mastery = seed % 3 === 0 ? 'mastered' : (seed % 2 === 0 ? 'weak' : 'todo')
    nodeStateMap[node.id] = {
      detail: `【${node.name || '知识点'}】\n1. 该知识点的定义与边界\n2. 常见题型与易错点\n3. 推荐复习顺序：概念 -> 例题 -> 迁移训练`,
      favorite: false,
      mastery,
      qa: [],
      notes: [],
      questions: buildMockQuestions(node)
    }
  }
  return nodeStateMap[node.id]
}

const buildMockQuestions = (node) => {
  const title = node?.name || '当前知识点'
  return [
    {
      id: `${node?.id || 'q'}-1`,
      stem: `关于“${title}”的理解，哪项描述最准确？`,
      options: ['只需背定义即可', '理解定义并结合场景应用', '只做题不总结', '记住关键词即可'],
      answer: '理解定义并结合场景应用',
      userAnswer: '',
      result: null
    },
    {
      id: `${node?.id || 'q'}-2`,
      stem: `遇到“${title}”相关综合题时，首要步骤是？`,
      options: ['直接代公式', '先识别题目条件与约束', '先看答案', '跳过难题'],
      answer: '先识别题目条件与约束',
      userAnswer: '',
      result: null
    }
  ]
}

const flattenNodes = (nodes, acc = []) => {
  (nodes || []).forEach((node) => {
    acc.push(node)
    if (Array.isArray(node.children) && node.children.length) {
      flattenNodes(node.children, acc)
    }
  })
  return acc
}

const pickFirstLeaf = (nodes) => {
  const list = flattenNodes(nodes, [])
  return list.find(item => !Array.isArray(item.children) || item.children.length === 0) || list[0] || null
}

const handleNodeClick = (node) => {
  if (!node) return
  currentNodeId.value = String(node.id || '')
  currentNodeName.value = String(node.name || '未命名知识点')
  const state = ensureNodeState(node)
  detailDraft.value = state.detail
  emit('node-click', node)
}

const countNodes = (tree) => {
  let count = 0
  ;(tree || []).forEach((node) => {
    count += 1
    if (Array.isArray(node.children) && node.children.length) {
      count += countNodes(node.children)
    }
  })
  return count
}

const filterNode = (keyword, data) => {
  if (!keyword) return true
  return String(data?.name || '').toLowerCase().includes(String(keyword).toLowerCase())
}

const noteCount = (id) => {
  if (!id || !nodeStateMap[id]) return 0
  return nodeStateMap[id].notes.length
}

const questionCount = (id) => {
  if (!id || !nodeStateMap[id]) return 0
  return nodeStateMap[id].questions.length
}

const masteryLabel = (id) => {
  const level = nodeStateMap[id]?.mastery || 'todo'
  if (level === 'mastered') return '已掌握'
  if (level === 'weak') return '未掌握'
  return '待学习'
}

const masteryTagType = (id) => {
  const level = nodeStateMap[id]?.mastery || 'todo'
  if (level === 'mastered') return 'success'
  if (level === 'weak') return 'danger'
  return 'warning'
}

const expandAll = () => {
  treeExpanded.value = true
}

const collapseAll = () => {
  treeExpanded.value = false
  window.setTimeout(() => {
    treeExpanded.value = true
  }, 0)
}

const saveTreeEdits = () => {
  ElMessage.success('知识树结构已保存（演示模式）')
}

const markMastered = () => {
  if (!currentNodeId.value) return
  ensureNodeState({ id: currentNodeId.value, name: currentNodeName.value }).mastery = 'mastered'
  ElMessage.success('已标记为已掌握')
}

const saveDetail = () => {
  if (!currentNodeId.value) return
  ensureNodeState({ id: currentNodeId.value, name: currentNodeName.value }).detail = detailDraft.value
  ElMessage.success('知识点内容已保存')
}

const toggleFavorite = () => {
  if (!currentNodeId.value) return
  const state = ensureNodeState({ id: currentNodeId.value, name: currentNodeName.value })
  state.favorite = !state.favorite
  ElMessage.success(state.favorite ? '已收藏知识点' : '已取消收藏')
}

const addToReview = () => {
  ElMessage.success('已添加到复习计划（演示模式）')
}

const regenPractice = () => {
  if (!currentNodeId.value) return
  const state = ensureNodeState({ id: currentNodeId.value, name: currentNodeName.value })
  state.questions = buildMockQuestions({ id: currentNodeId.value, name: currentNodeName.value })
  activeTab.value = 'practice'
  ElMessage.success('已重新生成配套习题')
}

const sendMockQa = () => {
  const q = String(qaInput.value || '').trim()
  if (!q || !currentNodeId.value) return
  const state = ensureNodeState({ id: currentNodeId.value, name: currentNodeName.value })
  const answer = `基于“${currentNodeName.value}”建议你按“概念-例题-错因复盘”三步学习。当前问题可先从关键词拆解入手。`
  state.qa.unshift({ q, a: answer })
  qaInput.value = ''
}

const saveNote = () => {
  const content = String(noteInput.value || '').trim()
  if (!content || !currentNodeId.value) return
  const state = ensureNodeState({ id: currentNodeId.value, name: currentNodeName.value })
  state.notes.unshift({
    id: `${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
    time: new Date().toLocaleString(),
    content
  })
  noteInput.value = ''
  ElMessage.success('笔记已保存')
}

const submitQuestion = (question) => {
  const ok = question.userAnswer === question.answer
  question.result = {
    ok,
    text: ok ? '回答正确，已同步掌握进度。' : `回答错误，正确答案：${question.answer}`
  }
  if (ok && currentNodeId.value) {
    ensureNodeState({ id: currentNodeId.value, name: currentNodeName.value }).mastery = 'mastered'
  }
}

const globalAsk = () => {
  ElMessage.success('已进入全局问答（演示模式）')
}

const generateReviewPack = () => {
  ElMessage.success('完整复习包已生成并同步到个人中心（演示模式）')
}

const exportMockDoc = () => {
  ElMessage.success('知识树与笔记已导出（演示模式）')
}

const syncToCourse = () => {
  ElMessage.success('已同步到课程学习进度（演示模式）')
}

const saveHistoryRecord = () => {
  if (!props.parseResult || !hasKnowledge.value) return
  const item = {
    id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
    fileName: props.uploadedFile?.name || '未命名资料',
    timeLabel: new Date().toLocaleString(),
    chapterCount: props.knowledgeList.length,
    pointCount: totalNodes.value
  }
  const next = [item, ...historyRecords.value].slice(0, 20)
  historyRecords.value = next
  try {
    window.localStorage.setItem(HISTORY_KEY, JSON.stringify(next))
  } catch (error) {
    console.warn('保存知识拆解历史失败', error)
  }
}

watch(() => props.knowledgeList, (next) => {
  if (!Array.isArray(next) || next.length === 0) {
    currentNodeId.value = ''
    currentNodeName.value = ''
    detailDraft.value = ''
    return
  }
  const first = pickFirstLeaf(next)
  if (!first) return
  handleNodeClick(first)
}, { immediate: true, deep: true })

watch(searchKeyword, (keyword) => {
  if (treeRef.value) {
    treeRef.value.filter(keyword)
  }
})

watch(() => props.isParsing, (parsing) => {
  if (parseStepTimer) {
    window.clearInterval(parseStepTimer)
    parseStepTimer = null
  }
  if (parsing) {
    parseStepActive.value = 0
    parseStepTimer = window.setInterval(() => {
      parseStepActive.value = Math.min(2, parseStepActive.value + 1)
    }, 900)
    return
  }
  if (props.parseResult) {
    parseStepActive.value = 3
    saveHistoryRecord()
  }
})

watch(() => props.parseResult, (value) => {
  if (value && !props.isParsing) {
    parseStepActive.value = 3
    saveHistoryRecord()
  }
})

onMounted(() => {
  try {
    const parsed = JSON.parse(window.localStorage.getItem(HISTORY_KEY) || '[]')
    historyRecords.value = Array.isArray(parsed) ? parsed : []
  } catch (error) {
    historyRecords.value = []
  }
})

onBeforeUnmount(() => {
  if (parseStepTimer) {
    window.clearInterval(parseStepTimer)
  }
})
</script>

<style scoped>
.knowledge-workbench {
  position: relative;
  background: linear-gradient(180deg, #ffffff 0%, #f6fbf8 100%);
  border-radius: 22px;
  border: 1px solid #d7e6de;
  padding: 16px;
  box-shadow: 0 18px 36px rgba(24, 55, 46, 0.08);
}

.workbench-head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: flex-start;
  margin-bottom: 14px;
}

.eyebrow {
  margin: 0;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: #67867a;
}

h3 {
  margin: 4px 0;
  color: #1f443d;
}

.head-tip {
  margin: 0;
  color: #5f7a70;
}

.head-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.init-state {
  display: grid;
  grid-template-columns: 1.1fr 0.9fr;
  gap: 14px;
}

.upload-card,
.history-mini {
  background: #fff;
  border: 1px solid #d9e8e1;
  border-radius: 14px;
  padding: 12px;
}

.upload-sub {
  margin: 6px 0 0;
  color: #6f877d;
}

.file-brief {
  margin-top: 10px;
  display: flex;
  justify-content: space-between;
  background: #eef6f2;
  border-radius: 10px;
  padding: 8px 10px;
}

.parse-options {
  margin-top: 10px;
  display: grid;
  gap: 8px;
}

.option-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.parse-btn {
  margin-top: 12px;
  width: 100%;
}

.parse-steps {
  margin: 12px 0;
}

.feature-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
}

.feature-item {
  background: linear-gradient(130deg, #eef8f3 0%, #f9fcfb 100%);
  border: 1px solid #d8e7df;
  border-radius: 12px;
  padding: 10px;
  color: #31544a;
}

.block-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.history-list {
  display: grid;
  gap: 8px;
}

.history-card {
  border: 1px solid #dce9e3;
  border-radius: 10px;
  padding: 10px;
  background: #f9fcfb;
}

.history-card p {
  margin: 6px 0;
  color: #6f857b;
}

.parsed-state {
  display: grid;
  grid-template-columns: 34% 66%;
  gap: 14px;
}

.left-pane,
.right-pane {
  background: #fff;
  border: 1px solid #d8e7df;
  border-radius: 14px;
  padding: 12px;
}

.left-tools {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 8px;
  margin-bottom: 8px;
}

.knowledge-tree {
  max-height: 560px;
  overflow: auto;
  border: 1px solid #e0ece6;
  border-radius: 10px;
  padding: 8px;
}

.tree-node-row {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.node-title {
  color: #1f463f;
}

.node-badges {
  display: flex;
  align-items: center;
  gap: 6px;
}

.node-counter {
  font-size: 12px;
  color: #617b71;
}

.left-foot {
  margin-top: 8px;
  display: flex;
  justify-content: space-between;
  color: #6b857b;
  font-size: 12px;
}

.detail-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.detail-head h4 {
  margin: 4px 0 0;
  color: #1d443d;
}

.detail-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.inline-actions {
  margin-top: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.qa-list,
.note-list,
.practice-list {
  display: grid;
  gap: 10px;
  margin-bottom: 10px;
}

.qa-item,
.note-item,
.practice-item {
  border: 1px solid #dde9e3;
  border-radius: 10px;
  padding: 10px;
}

.qa-input-row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 8px;
}

.practice-item h5 {
  margin: 0 0 8px;
}

.q-result {
  margin-top: 8px;
  font-weight: 600;
}

.q-result.ok {
  color: #1f8f57;
}

.q-result.bad {
  color: #cc4c43;
}

.data-strip {
  margin-top: 10px;
  border-top: 1px dashed #dae8e2;
  padding-top: 10px;
}

.data-strip p {
  margin: 8px 0 0;
  color: #607a70;
}

.floating-actions {
  position: fixed;
  right: 24px;
  bottom: 24px;
  display: grid;
  gap: 10px;
  z-index: 20;
}

.history-dialog-list {
  max-height: 360px;
  overflow: auto;
  display: grid;
  gap: 8px;
}

.history-dialog-card {
  border: 1px solid #dce9e3;
  border-radius: 10px;
  padding: 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.history-dialog-card p {
  margin: 4px 0 0;
  color: #698378;
}

@media (max-width: 1100px) {
  .init-state,
  .parsed-state {
    grid-template-columns: 1fr;
  }

  .left-tools {
    grid-template-columns: 1fr 1fr;
  }

  .floating-actions {
    right: 12px;
    bottom: 12px;
    grid-auto-flow: column;
  }
}
</style>