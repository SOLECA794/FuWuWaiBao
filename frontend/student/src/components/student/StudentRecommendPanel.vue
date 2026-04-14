<template>
  <section class="recommend-panel">
    <header class="recommend-header">
      <div>
        <div class="recommend-kicker">学习推荐</div>
        <h3>遗传算法课件配套推荐中心</h3>
        <p>已切换为前端演示模式：全流程可跑通，无需后端返回结果。</p>
      </div>
      <el-button type="primary" plain @click="openPlanDialog">查看今日学习单</el-button>
    </header>

    <div class="demo-director-strip">
      <div>
        <strong>演示路线：推荐资源 → 提交练习 → 个人中心复盘</strong>
        <p>先加入学习单，再点击“一键演示闭环”，可直接推进到随堂练习页面。</p>
      </div>
      <div class="director-actions">
        <el-button size="small" plain @click="emit('navigate-section', 'classroom')">回课堂学习</el-button>
        <el-button size="small" type="primary" plain @click="emit('navigate-section', 'practice')">去随堂练习</el-button>
      </div>
    </div>

    <div class="recommend-body">
      <aside class="recommend-left-pane">
        <div class="recommend-toolbar">
          <el-input
            v-model.trim="keyword"
            clearable
            placeholder="输入关键词，如：遗传算法 选择复制"
            @keyup.enter="searchResources"
          />
          <el-button type="primary" @click="searchResources">筛选资源</el-button>
        </div>

        <div class="hint-row" v-if="defaultKeyword">
          建议关键词：{{ defaultKeyword }}
          <el-button text size="small" @click="useDefaultKeyword">一键应用</el-button>
        </div>

        <div class="filter-grid">
          <el-select v-model="selectedSource" size="small">
            <el-option label="全部来源" value="all" />
            <el-option label="B站" value="Bilibili" />
            <el-option label="51教习" value="51教习" />
            <el-option label="慕课" value="MOOC" />
          </el-select>
          <el-select v-model="selectedDifficulty" size="small">
            <el-option label="全部难度" value="all" />
            <el-option label="入门" value="入门" />
            <el-option label="进阶" value="进阶" />
            <el-option label="强化" value="强化" />
          </el-select>
        </div>

        <div class="quick-actions">
          <el-button size="small" plain @click="openFlowDialog('导学')">导学流程</el-button>
          <el-button size="small" plain @click="openFlowDialog('练习')">练习流程</el-button>
          <el-button size="small" plain @click="openFlowDialog('冲刺')">冲刺流程</el-button>
        </div>

        <div class="queue-card">
          <h4>今日学习单</h4>
          <div v-if="studyQueue.length" class="queue-list">
            <article v-for="item in studyQueue" :key="item.id" class="queue-item">
              <strong>{{ item.title }}</strong>
              <span>{{ item.duration }} · {{ item.source }}</span>
            </article>
          </div>
          <div v-else class="queue-empty">还未加入资源，点击右侧“加入学习单”即可。</div>
        </div>
      </aside>

      <section class="recommend-right-pane">
        <div class="recommend-count">共 {{ filteredResources.length }} 条推荐资源</div>
        <div v-if="filteredResources.length" class="recommend-grid">
          <article class="recommend-item" v-for="item in filteredResources" :key="item.id">
            <header>
              <h4>{{ item.title }}</h4>
              <span class="source-tag">{{ item.source }}</span>
            </header>
            <p class="reason">{{ item.reason }}</p>
            <div class="meta-row">
              <span>{{ item.duration }}</span>
              <span>{{ item.difficulty }}</span>
              <span>{{ item.fitNode }}</span>
            </div>
            <div class="actions">
              <el-button size="small" type="primary" @click="openDetail(item)">详情</el-button>
              <el-button size="small" plain @click="enqueueResource(item)">加入学习单</el-button>
              <el-button size="small" plain @click="simulateLearning(item)">模拟学习</el-button>
              <el-button size="small" type="success" plain @click="runDemoFlow(item)">一键演示闭环</el-button>
            </div>
          </article>
        </div>
        <div v-else class="recommend-empty">当前筛选条件下暂无资源，建议切换关键词或难度。</div>
      </section>
    </div>

    <el-dialog v-model="detailDialogVisible" title="资源详情" width="560px">
      <div v-if="activeResource" class="detail-dialog">
        <h4>{{ activeResource.title }}</h4>
        <p>{{ activeResource.reason }}</p>
        <ul>
          <li>来源：{{ activeResource.source }}</li>
          <li>时长：{{ activeResource.duration }}</li>
          <li>适配节点：{{ activeResource.fitNode }}</li>
          <li>推荐策略：{{ activeResource.strategy }}</li>
        </ul>
      </div>
    </el-dialog>

    <el-dialog v-model="planDialogVisible" title="今日学习单" width="620px">
      <div class="plan-dialog">
        <article v-for="(item, index) in studyQueue" :key="item.id" class="plan-step">
          <span>{{ index + 1 }}</span>
          <div>
            <strong>{{ item.title }}</strong>
            <p>{{ item.duration }} · {{ item.source }} · {{ item.fitNode }}</p>
          </div>
        </article>
        <div v-if="!studyQueue.length" class="queue-empty">学习单为空，请先加入推荐资源。</div>
      </div>
    </el-dialog>

    <el-dialog v-model="flowDialogVisible" :title="flowDialogTitle" width="560px">
      <div class="flow-dialog">
        <p v-for="(line, index) in flowDialogLines" :key="`${flowDialogTitle}-${index}`">{{ line }}</p>
      </div>
    </el-dialog>
  </section>
</template>

<script setup>
/* eslint-disable no-undef */
import { computed, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  courseName: {
    type: String,
    default: ''
  },
  currentNodeTitle: {
    type: String,
    default: ''
  },
  currentPage: {
    type: Number,
    default: 1
  }
})
const emit = defineEmits(['navigate-section'])

const keyword = ref('')
const selectedSource = ref('all')
const selectedDifficulty = ref('all')
const querySnapshot = ref('')
const activeResource = ref(null)
const detailDialogVisible = ref(false)
const planDialogVisible = ref(false)
const flowDialogVisible = ref(false)
const flowDialogTitle = ref('')
const flowDialogLines = ref([])
const studyQueue = ref([])

const mockResources = ref([
  {
    id: 'ga-1',
    title: '遗传算法入门：编码、适应度与种群演化',
    source: 'Bilibili',
    duration: '18分钟',
    difficulty: '入门',
    fitNode: '基础概念',
    strategy: '先建立概念图，再看例题',
    reason: '适合第1页到第3页的术语理解，覆盖“编码-选择-交叉-变异”完整主线。'
  },
  {
    id: 'ga-2',
    title: '轮盘赌选择与锦标赛选择对比实验',
    source: '51教习',
    duration: '25分钟',
    difficulty: '进阶',
    fitNode: '选择复制',
    strategy: '先看对比图，再做参数实验',
    reason: '针对“选择复制”薄弱点，给出两类策略的收敛速度和多样性差异。'
  },
  {
    id: 'ga-3',
    title: '单点交叉、多点交叉与变异率调参实战',
    source: 'MOOC',
    duration: '32分钟',
    difficulty: '强化',
    fitNode: '交叉与变异',
    strategy: '先固定编码长度，再逐步调变异率',
    reason: '用于课堂练习后的强化复盘，能直接映射你当前课件中的应用题。'
  },
  {
    id: 'ga-4',
    title: '遗传算法在路径规划中的可视化演示',
    source: 'Bilibili',
    duration: '21分钟',
    difficulty: '进阶',
    fitNode: '应用案例',
    strategy: '先看动画，再回到公式推导',
    reason: '可将抽象流程映射到具体问题，减少“会背不会用”的情况。'
  },
  {
    id: 'ga-5',
    title: 'GA 课堂冲刺：高频题型与易错点清单',
    source: '51教习',
    duration: '16分钟',
    difficulty: '强化',
    fitNode: '期末冲刺',
    strategy: '错因分类 + 速记卡片',
    reason: '适合考前 24 小时快速过一遍，重点覆盖概率计算与交叉结果判断。'
  }
])

const addResourceToQueue = (item) => {
  const exists = studyQueue.value.some((record) => record.id === item.id)
  if (exists) return false
  studyQueue.value = [...studyQueue.value, item]
  return true
}

const defaultKeyword = computed(() => {
  const parts = [
    String(props.currentNodeTitle || '').trim(),
    String(props.courseName || '').trim(),
    props.currentPage ? `第${props.currentPage}页` : ''
  ].filter(Boolean)
  return parts.join(' ')
})

const filteredResources = computed(() => {
  const query = String(querySnapshot.value || keyword.value || '').trim().toLowerCase()
  return mockResources.value.filter((item) => {
    const sourceMatch = selectedSource.value === 'all' || item.source === selectedSource.value
    const diffMatch = selectedDifficulty.value === 'all' || item.difficulty === selectedDifficulty.value
    const queryMatch = !query || [item.title, item.reason, item.fitNode, item.strategy].join(' ').toLowerCase().includes(query)
    return sourceMatch && diffMatch && queryMatch
  })
})

const useDefaultKeyword = () => {
  if (!defaultKeyword.value) return
  keyword.value = defaultKeyword.value
}

const searchResources = () => {
  querySnapshot.value = String(keyword.value || '').trim() || defaultKeyword.value
  if (!querySnapshot.value) {
    ElMessage.warning('请输入关键词后再筛选')
    return
  }
  ElMessage.success('已按当前条件更新推荐列表（演示模式）')
}

const openDetail = (item) => {
  activeResource.value = item
  detailDialogVisible.value = true
}

const enqueueResource = (item) => {
  addResourceToQueue(item)
  planDialogVisible.value = true
}

const simulateLearning = (item) => {
  flowDialogTitle.value = `模拟学习流程：${item.title}`
  flowDialogLines.value = [
    '1. 进入资源并完成导学预览（约 3 分钟）',
    '2. 结合课件当前页做 2 道自测题（约 6 分钟）',
    '3. AI 助手回顾错因并输出复习提纲（约 4 分钟）',
    '4. 自动写入今日学习单，等待课堂回放联动（已模拟）'
  ]
  flowDialogVisible.value = true
}

const runDemoFlow = (item) => {
  const inserted = addResourceToQueue(item)
  if (inserted) {
    ElMessage.success('资源已加入学习单，正在进入随堂练习')
  } else {
    ElMessage.success('学习单已存在该资源，正在进入随堂练习')
  }
  emit('navigate-section', 'practice')
}

const openFlowDialog = (mode) => {
  const map = {
    导学: [
      '导学模式：先理解概念定义，再看流程图。',
      '建议顺序：基础概念 -> 选择复制 -> 交叉变异。'
    ],
    练习: [
      '练习模式：每学完 1 个节点立即做 2 题。',
      '错题自动加入“薄弱强化”队列。'
    ],
    冲刺: [
      '冲刺模式：优先做概率与交叉计算题。',
      '建议在 30 分钟内完成 1 轮专题复习。'
    ]
  }
  flowDialogTitle.value = `${mode}流程说明`
  flowDialogLines.value = map[mode] || ['当前模式暂无说明。']
  flowDialogVisible.value = true
}

const openPlanDialog = () => {
  planDialogVisible.value = true
}

watch(defaultKeyword, (next) => {
  if (!keyword.value && next) {
    keyword.value = next
    querySnapshot.value = next
  }
}, { immediate: true })
</script>

<style scoped>
.recommend-panel {
  height: 100%;
  min-height: 0;
  border: 1px solid #d8e5de;
  border-radius: 18px;
  background: linear-gradient(180deg, #ffffff 0%, #f6faf8 100%);
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  overflow: hidden;
}

.recommend-header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: flex-start;
}

.demo-director-strip {
  border: 1px solid #d6e5dd;
  border-radius: 12px;
  background: linear-gradient(145deg, #f3f9f5 0%, #eaf4ef 100%);
  padding: 10px 12px;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}

.demo-director-strip strong {
  font-size: 13px;
  color: #21483f;
}

.demo-director-strip p {
  margin: 4px 0 0;
  font-size: 12px;
  color: #55756b;
}

.director-actions {
  display: inline-flex;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.recommend-kicker {
  font-size: 11px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  font-weight: 700;
  color: #6a8278;
}

.recommend-header h3 {
  margin-top: 4px;
  font-size: 18px;
  color: #23463f;
}

.recommend-header p {
  margin-top: 6px;
  font-size: 13px;
  color: #6f867d;
}

.recommend-body {
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: 310px minmax(0, 1fr);
  gap: 10px;
  overflow: hidden;
}

.recommend-left-pane,
.recommend-right-pane {
  min-height: 0;
  border: 1px solid #dbe8e1;
  border-radius: 14px;
  background: #fff;
  padding: 10px;
}

.recommend-left-pane {
  display: flex;
  flex-direction: column;
  gap: 10px;
  overflow: auto;
}

.recommend-toolbar {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 8px;
}

.hint-row {
  font-size: 12px;
  color: #58736a;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.filter-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 8px;
}

.quick-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.queue-card {
  border: 1px solid #dce9e2;
  border-radius: 12px;
  background: #f8fcf9;
  padding: 10px;
}

.queue-card h4 {
  margin: 0 0 8px;
  font-size: 14px;
  color: #274d46;
}

.queue-list {
  display: grid;
  gap: 8px;
}

.queue-item {
  border: 1px solid #d6e6de;
  border-radius: 10px;
  background: #fff;
  padding: 8px;
  display: grid;
  gap: 4px;
}

.queue-item strong {
  font-size: 13px;
  color: #24493f;
}

.queue-item span,
.queue-empty {
  font-size: 12px;
  color: #678077;
}

.recommend-right-pane {
  display: flex;
  flex-direction: column;
  gap: 8px;
  overflow: hidden;
}

.recommend-count {
  font-size: 12px;
  color: #58746a;
}

.recommend-grid {
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  overflow: auto;
  padding-right: 2px;
}

.recommend-item {
  border: 1px solid #d7e5dd;
  border-radius: 12px;
  background: linear-gradient(180deg, #ffffff 0%, #f8fcfa 100%);
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.recommend-item header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}

.recommend-item h4 {
  margin: 0;
  font-size: 15px;
  color: #274d46;
}

.source-tag {
  font-size: 11px;
  border-radius: 999px;
  border: 1px solid #d1e2da;
  padding: 2px 8px;
  color: #42665d;
  background: #eef5f1;
}

.reason {
  margin: 0;
  font-size: 12px;
  line-height: 1.55;
  color: #577068;
}

.meta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  font-size: 11px;
  color: #5f7a70;
}

.actions {
  margin-top: auto;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.recommend-empty {
  flex: 1;
  min-height: 0;
  border: 1px dashed #cfddd5;
  border-radius: 12px;
  padding: 12px;
  color: #70857c;
  font-size: 13px;
}

.detail-dialog h4 {
  margin: 0 0 8px;
  color: #254a41;
}

.detail-dialog p {
  margin: 0 0 10px;
  color: #5f7b71;
}

.detail-dialog ul {
  margin: 0;
  padding-left: 18px;
  color: #4c665d;
  line-height: 1.8;
}

.plan-dialog {
  display: grid;
  gap: 10px;
}

.plan-step {
  display: flex;
  gap: 10px;
  border: 1px solid #dbe8e2;
  border-radius: 10px;
  padding: 10px;
}

.plan-step span {
  width: 24px;
  height: 24px;
  border-radius: 999px;
  background: #e7f2ec;
  color: #2f605a;
  font-size: 12px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.plan-step p {
  margin: 4px 0 0;
  font-size: 12px;
  color: #668177;
}

.flow-dialog {
  display: grid;
  gap: 8px;
}

.flow-dialog p {
  margin: 0;
  border: 1px solid #dbe8e2;
  background: #f7fbf9;
  border-radius: 10px;
  padding: 8px 10px;
  color: #4f6d63;
  font-size: 13px;
}

@media (max-width: 1080px) {
  .recommend-body {
    grid-template-columns: 1fr;
  }

  .recommend-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .recommend-toolbar {
    grid-template-columns: minmax(0, 1fr);
  }

  .recommend-header {
    flex-direction: column;
  }

  .demo-director-strip {
    flex-direction: column;
  }
}
</style>
