<template>
  <div class="tab-content">
    <div class="stats-header" v-if="currentCourseId">
      <h4>学情分析与可视化决策中心 - {{ currentCourseName }}</h4>
      <p>围绕 节点热度 / 掌握度 / 班级趋势 自动生成教学建议，形成教学闭环。</p>
    </div>

    <div v-if="currentCourseId" class="stats-layout">
      <section class="summary-grid">
        <article class="stat-card">
          <div class="stat-value">{{ studentStats.totalQuestions || 0 }}</div>
          <div class="stat-label">总提问数</div>
        </article>
        <article class="stat-card">
          <div class="stat-value">{{ studentStats.activeSessions || 0 }}</div>
          <div class="stat-label">活跃会话数</div>
        </article>
        <article class="stat-card">
          <div class="stat-value">{{ Math.round(Number(studentStats.avgTurnsPerSession || 0) * 10) / 10 }}</div>
          <div class="stat-label">平均追问轮次</div>
        </article>
        <article class="stat-card">
          <div class="stat-value">{{ Math.round(Number(masteryRadar.avgMastery || 0)) }}%</div>
          <div class="stat-label">平均掌握度</div>
        </article>
      </section>

      <section class="decision-board">
        <div class="decision-title">决策建议</div>
        <div class="decision-summary">{{ learningInsights.summary || '暂无建议' }}</div>
        <div class="decision-grid">
          <article class="decision-card">
            <h5>建议重讲节点</h5>
            <ul>
              <li v-for="item in reteachNodes" :key="item.nodeId">
                <strong>{{ item.title }}</strong>
                <span>P{{ item.page }} | 错误率 {{ Math.round((Number(item.errorRate || 0) * 1000)) / 10 }}%</span>
              </li>
              <li v-if="reteachNodes.length === 0" class="empty-row">暂无明显重讲节点</li>
            </ul>
          </article>
          <article class="decision-card">
            <h5>建议补前置知识</h5>
            <ul>
              <li v-for="item in prerequisiteGaps" :key="`${item.nodeId}_${item.suggestedPrereqId}`">
                <strong>{{ item.title }}</strong>
                <span>建议前置：{{ item.suggestedPrereq }} ({{ item.suggestedPrereqId }})</span>
              </li>
              <li v-if="prerequisiteGaps.length === 0" class="empty-row">暂无前置缺口</li>
            </ul>
          </article>
        </div>
      </section>

      <section class="chart-grid">
        <article class="chart-card">
          <header>
            <h5>知识点热力图</h5>
            <span>热力 = 提问 + 对话 + 重讲加权</span>
          </header>
          <div class="chart-view" ref="heatmapChartRef"></div>
        </article>

        <article class="chart-card">
          <header>
            <h5>知识掌握度雷达图</h5>
            <span>按页聚合节点掌握度</span>
          </header>
          <div class="chart-view" ref="radarChartRef"></div>
        </article>

        <article class="chart-card full-width">
          <header>
            <h5>班级学情趋势图</h5>
            <span>按天观察提问与重讲变化</span>
          </header>
          <div class="chart-view trend" ref="trendChartRef"></div>
        </article>
      </section>

      <section class="node-table-card">
        <div class="table-header">
          <h5>节点聚合指标</h5>
          <button
            v-if="hasUncoveredNodes"
            type="button"
            class="tiny-filter-btn"
            @click="onlyUncovered = !onlyUncovered"
          >
            {{ onlyUncovered ? '显示全部节点' : '仅看未覆盖节点' }}
          </button>
        </div>
        <div class="node-table" v-if="filteredNodeStats.length">
          <div class="row head">
            <span>节点</span>
            <span>提问</span>
            <span>停留(秒)</span>
            <span>错误率</span>
            <span>掌握度</span>
          </div>
          <div class="row" v-for="item in filteredNodeStats" :key="item.nodeId">
            <span>{{ item.title || item.nodeId }}</span>
            <span>{{ item.questionCount || item.dialogueCount || 0 }}</span>
            <span>{{ Number(item.stayTime || 0).toFixed(1) }}</span>
            <span>{{ Math.round(Number(item.errorRate || 0) * 1000) / 10 }}%</span>
            <span>{{ Math.round(Number(item.masteryScore || 0)) }}%</span>
          </div>
        </div>
        <div class="empty-tip" v-else>暂无节点聚合数据</div>
      </section>
    </div>

    <div v-else class="empty-tip">请先选择一个课件查看学情数据</div>
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import * as echarts from 'echarts'

const props = defineProps({
  currentCourseId: {
    type: String,
    default: ''
  },
  currentCourseName: {
    type: String,
    default: ''
  },
  studentStats: {
    type: Object,
    default: () => ({
      totalQuestions: 0,
      hotPages: [],
      keyDifficulties: '暂无',
      nodeStats: [],
      mappingCoverage: null,
      nodeHeatmap: [],
      masteryRadar: { indicators: [], values: [], avgMastery: 0 },
      classTrend: [],
      learningInsights: { reteachNodes: [], prerequisiteGaps: [], summary: '' }
    })
  }
})

const onlyUncovered = ref(false)
const heatmapChartRef = ref(null)
const radarChartRef = ref(null)
const trendChartRef = ref(null)

let heatmapChart = null
let radarChart = null
let trendChart = null

const nodeStats = computed(() => (Array.isArray(props.studentStats?.nodeStats) ? props.studentStats.nodeStats : []))
const heatmapData = computed(() => (Array.isArray(props.studentStats?.nodeHeatmap) ? props.studentStats.nodeHeatmap : []))
const masteryRadar = computed(() => props.studentStats?.masteryRadar || { indicators: [], values: [], avgMastery: 0 })
const classTrend = computed(() => (Array.isArray(props.studentStats?.classTrend) ? props.studentStats.classTrend : []))
const learningInsights = computed(() => props.studentStats?.learningInsights || { reteachNodes: [], prerequisiteGaps: [], summary: '' })

const reteachNodes = computed(() => {
  const list = Array.isArray(learningInsights.value?.reteachNodes) ? learningInsights.value.reteachNodes : []
  return list.slice(0, 6)
})

const prerequisiteGaps = computed(() => {
  const list = Array.isArray(learningInsights.value?.prerequisiteGaps) ? learningInsights.value.prerequisiteGaps : []
  return list.slice(0, 6)
})

const uncoveredNodeIdSet = computed(() => {
  const ids = props.studentStats?.mappingCoverage?.uncoveredNodeIds
  if (!Array.isArray(ids) || ids.length === 0) return new Set()
  return new Set(ids.map(item => String(item || '').trim()).filter(Boolean))
})

const hasUncoveredNodes = computed(() => uncoveredNodeIdSet.value.size > 0)

const filteredNodeStats = computed(() => {
  if (!onlyUncovered.value || uncoveredNodeIdSet.value.size === 0) {
    return nodeStats.value
  }
  return nodeStats.value.filter(item => uncoveredNodeIdSet.value.has(String(item?.nodeId || '').trim()))
})

const ensureCharts = () => {
  if (heatmapChartRef.value && !heatmapChart) {
    heatmapChart = echarts.init(heatmapChartRef.value)
  }
  if (radarChartRef.value && !radarChart) {
    radarChart = echarts.init(radarChartRef.value)
  }
  if (trendChartRef.value && !trendChart) {
    trendChart = echarts.init(trendChartRef.value)
  }
}

const renderHeatmap = () => {
  if (!heatmapChart) return
  const labels = heatmapData.value.map(item => item.title || item.nodeId || '-')
  const values = heatmapData.value.map(item => Number(item.heat || 0))
  heatmapChart.setOption({
    grid: { left: 70, right: 18, top: 26, bottom: 24 },
    xAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: '#e2e8f0' } },
      axisLabel: { color: '#64748b' }
    },
    yAxis: {
      type: 'category',
      data: labels,
      axisLabel: { color: '#334155', width: 150, overflow: 'truncate' }
    },
    series: [{
      type: 'bar',
      data: values,
      itemStyle: {
        color: params => {
          const value = Number(params.value || 0)
          if (value >= 15) return '#dc2626'
          if (value >= 8) return '#f97316'
          return '#0f766e'
        },
        borderRadius: [0, 10, 10, 0]
      },
      label: { show: true, position: 'right', color: '#1e293b' }
    }],
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' }
    }
  }, true)
}

const renderRadar = () => {
  if (!radarChart) return
  const indicators = Array.isArray(masteryRadar.value?.indicators) ? masteryRadar.value.indicators : []
  const values = Array.isArray(masteryRadar.value?.values) ? masteryRadar.value.values : []
  radarChart.setOption({
    radar: {
      indicator: indicators.length ? indicators : [{ name: '暂无数据', max: 100 }],
      radius: '60%',
      splitArea: { areaStyle: { color: ['#f8fafc', '#eef2ff'] } },
      axisName: { color: '#334155' }
    },
    series: [{
      type: 'radar',
      data: [{
        value: values.length ? values : [0],
        name: '掌握度',
        areaStyle: { color: 'rgba(15, 118, 110, 0.25)' },
        lineStyle: { color: '#0f766e' },
        itemStyle: { color: '#0f766e' }
      }]
    }],
    tooltip: {}
  }, true)
}

const renderTrend = () => {
  if (!trendChart) return
  const days = classTrend.value.map(item => item.day)
  const question = classTrend.value.map(item => Number(item.questionCount || 0))
  const reteach = classTrend.value.map(item => Number(item.reteachCount || 0))
  const errorRate = classTrend.value.map(item => Math.round(Number(item.errorRate || 0) * 1000) / 10)

  trendChart.setOption({
    legend: { top: 4, textStyle: { color: '#475569' } },
    tooltip: { trigger: 'axis' },
    grid: { left: 42, right: 46, top: 40, bottom: 30 },
    xAxis: { type: 'category', data: days, axisLabel: { color: '#64748b' } },
    yAxis: [
      { type: 'value', name: '次数', axisLabel: { color: '#64748b' }, splitLine: { lineStyle: { color: '#e2e8f0' } } },
      { type: 'value', name: '错误率%', axisLabel: { color: '#64748b' } }
    ],
    series: [
      { name: '提问次数', type: 'line', data: question, smooth: true, itemStyle: { color: '#0f766e' } },
      { name: '重讲次数', type: 'line', data: reteach, smooth: true, itemStyle: { color: '#f97316' } },
      { name: '错误率', type: 'bar', yAxisIndex: 1, data: errorRate, itemStyle: { color: '#dc2626', opacity: 0.6 } }
    ]
  }, true)
}

const renderAll = async () => {
  await nextTick()
  ensureCharts()
  renderHeatmap()
  renderRadar()
  renderTrend()
}

const handleResize = () => {
  if (heatmapChart) heatmapChart.resize()
  if (radarChart) radarChart.resize()
  if (trendChart) trendChart.resize()
}

watch(
  () => [heatmapData.value, masteryRadar.value, classTrend.value, props.currentCourseId],
  () => {
    renderAll()
  },
  { deep: true }
)

onMounted(() => {
  renderAll()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  if (heatmapChart) heatmapChart.dispose()
  if (radarChart) radarChart.dispose()
  if (trendChart) trendChart.dispose()
})
</script>

<style scoped>
.tab-content {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 18px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
  overflow: auto;
}

.stats-header {
  margin-bottom: 16px;
}

.stats-header h4 {
  margin: 0;
  color: #0f172a;
}

.stats-header p {
  margin: 6px 0 0;
  color: #64748b;
  font-size: 13px;
}

.stats-layout {
  display: grid;
  gap: 14px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
}

.stat-card {
  border: 1px solid #dbe5df;
  border-radius: 12px;
  padding: 14px;
  background: linear-gradient(180deg, #fbfdfc 0%, #f4f8f6 100%);
}

.stat-value {
  font-size: 28px;
  color: #0f172a;
  font-weight: 700;
}

.stat-label {
  margin-top: 6px;
  color: #64748b;
  font-size: 12px;
}

.decision-board {
  border: 1px solid #dbe5df;
  border-radius: 12px;
  padding: 14px;
  background: #f8fbfa;
}

.decision-title {
  font-size: 16px;
  font-weight: 700;
  color: #14532d;
}

.decision-summary {
  margin-top: 6px;
  color: #334155;
  font-size: 13px;
}

.decision-grid {
  margin-top: 10px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.decision-card {
  border: 1px solid #cfded5;
  border-radius: 10px;
  background: #ffffff;
  padding: 10px;
}

.decision-card h5 {
  margin: 0 0 8px;
  font-size: 13px;
  color: #0f172a;
}

.decision-card ul {
  margin: 0;
  padding: 0;
  list-style: none;
  display: grid;
  gap: 8px;
}

.decision-card li {
  display: grid;
  gap: 2px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 7px;
  background: #f8fafc;
}

.decision-card strong {
  font-size: 13px;
  color: #1e293b;
}

.decision-card span {
  font-size: 12px;
  color: #64748b;
}

.chart-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.chart-card {
  border: 1px solid #dbe5df;
  border-radius: 12px;
  padding: 10px;
  background: #ffffff;
}

.chart-card.full-width {
  grid-column: span 2;
}

.chart-card header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 8px;
}

.chart-card h5 {
  margin: 0;
  color: #0f172a;
  font-size: 14px;
}

.chart-card span {
  color: #64748b;
  font-size: 12px;
}

.chart-view {
  margin-top: 8px;
  height: 280px;
}

.chart-view.trend {
  height: 300px;
}

.node-table-card {
  border: 1px solid #dbe5df;
  border-radius: 12px;
  background: #ffffff;
  padding: 12px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.table-header h5 {
  margin: 0;
  font-size: 14px;
  color: #0f172a;
}

.node-table {
  display: grid;
  gap: 6px;
}

.row {
  display: grid;
  grid-template-columns: 2.3fr 0.8fr 0.9fr 0.8fr 0.8fr;
  gap: 8px;
  align-items: center;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 8px 10px;
  background: #f8fafc;
  font-size: 12px;
  color: #334155;
}

.row.head {
  background: #eef2f7;
  color: #1e293b;
  font-weight: 700;
}

.tiny-filter-btn {
  border: 1px solid #cbd5e1;
  background: #f8fafc;
  color: #334155;
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 12px;
  cursor: pointer;
}

.empty-tip,
.empty-row {
  color: #64748b;
  font-size: 12px;
}

@media (max-width: 1200px) {
  .summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .decision-grid {
    grid-template-columns: 1fr;
  }

  .chart-grid {
    grid-template-columns: 1fr;
  }

  .chart-card.full-width {
    grid-column: span 1;
  }
}
</style>
