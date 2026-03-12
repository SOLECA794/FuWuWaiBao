<template>
  <div class="tab-content">
    <div class="chart-header" v-if="currentCourseId">
      <h4>学习卡点分析 - {{ currentCourseName }}</h4>
      <div class="chart-type">
        <span>图表类型：</span>
        <button
          v-for="type in chartTypes"
          :key="type.value"
          class="chart-btn"
          :class="{ active: chartType === type.value }"
          @click="$emit('update:chartType', type.value)"
        >
          {{ type.label }}
        </button>
      </div>
    </div>
    <div v-if="currentCourseId" class="chart-container">
      <div ref="chartRef" class="chart"></div>
      <div class="chart-tip">
        <p>数据说明：</p>
        <ul>
          <li>提问量：该页面学生发起的提问总数</li>
          <li>停留时长：学生平均停留时长（秒）</li>
          <li>卡点指数：综合提问量+停留时长计算的卡点程度（0-10）</li>
        </ul>
      </div>
    </div>
    <div v-else class="empty-tip">请先选择一个课件查看卡点分析</div>
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
  chartType: {
    type: String,
    default: 'bar'
  },
  cardData: {
    type: Array,
    default: () => []
  }
})

defineEmits(['update:chartType'])

const chartTypes = [
  { value: 'bar', label: '柱状图' },
  { value: 'line', label: '折线图' },
  { value: 'pie', label: '饼图' }
]

const chartRef = ref(null)
let chartInstance = null

const chartOption = computed(() => {
  const pages = props.cardData.map(item => `第${item.page}页`)
  const questionCounts = props.cardData.map(item => item.提问量)
  const stayTimes = props.cardData.map(item => item.停留时长)
  const cardScores = props.cardData.map(item => item.卡点指数)

  if (props.chartType === 'line') {
    return {
      title: { text: '各页面学习卡点趋势' },
      tooltip: { trigger: 'axis' },
      legend: { data: ['提问量', '停留时长(秒)', '卡点指数'] },
      xAxis: { type: 'category', data: pages },
      yAxis: { type: 'value' },
      series: [
        { name: '提问量', type: 'line', data: questionCounts },
        { name: '停留时长(秒)', type: 'line', data: stayTimes },
        { name: '卡点指数', type: 'line', data: cardScores, lineStyle: { color: '#ff4d4f' }, itemStyle: { color: '#ff4d4f' } }
      ]
    }
  }

  if (props.chartType === 'pie') {
    const top5CardData = [...props.cardData].sort((a, b) => b.卡点指数 - a.卡点指数).slice(0, 5)
    return {
      title: { text: 'TOP5 卡点页面占比' },
      tooltip: { trigger: 'item' },
      legend: { orient: 'vertical', left: 'left', data: top5CardData.map(item => `第${item.page}页`) },
      series: [
        {
          name: '卡点指数',
          type: 'pie',
          radius: ['40%', '70%'],
          data: top5CardData.map(item => ({
            name: `第${item.page}页`,
            value: item.卡点指数
          })),
          label: { formatter: '{b}: {c} ({d}%)' }
        }
      ]
    }
  }

  return {
    title: { text: '各页面学习卡点数据' },
    tooltip: { trigger: 'axis' },
    legend: { data: ['提问量', '停留时长(秒)', '卡点指数'] },
    xAxis: { type: 'category', data: pages },
    yAxis: { type: 'value' },
    series: [
      { name: '提问量', type: 'bar', data: questionCounts },
      { name: '停留时长(秒)', type: 'bar', data: stayTimes },
      { name: '卡点指数', type: 'bar', data: cardScores, itemStyle: { color: '#ff4d4f' } }
    ]
  }
})

const renderChart = async () => {
  if (!props.currentCourseId || !chartRef.value) return
  await nextTick()
  if (!chartInstance) {
    chartInstance = echarts.init(chartRef.value)
  }
  chartInstance.setOption(chartOption.value)
  chartInstance.resize()
}

const resizeChart = () => {
  if (chartInstance) {
    chartInstance.resize()
  }
}

watch(() => [props.currentCourseId, props.chartType, props.cardData], () => {
  renderChart()
}, { deep: true })

onMounted(() => {
  renderChart()
  window.addEventListener('resize', resizeChart)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', resizeChart)
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
})
</script>

<style scoped>
.tab-content {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 18px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}
.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.chart-type {
  display: flex;
  align-items: center;
  gap: 8px;
}
.chart-btn {
  border: none;
  border-radius: 8px;
  padding: 8px 12px;
  background: #e2e8f0;
  color: #334155;
  cursor: pointer;
}
.chart-btn.active {
  background: #2F605A;
  color: #fff;
}
.chart-container {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.chart {
  width: 100%;
  height: 400px;
}
.chart-tip {
  background: #F4F7F7;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 12px 14px;
  color: #475569;
}
.chart-tip ul {
  padding-left: 18px;
  margin-top: 8px;
}
.empty-tip {
  text-align: center;
  color: #64748b;
}
</style>