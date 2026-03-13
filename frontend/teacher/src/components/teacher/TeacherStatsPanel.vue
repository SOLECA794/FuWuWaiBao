<template>
  <div class="tab-content">
    <div class="stats-header" v-if="currentCourseId">
      <h4>学情分析 - {{ currentCourseName }}</h4>
    </div>
    <div class="stats-grid" v-if="currentCourseId">
      <div class="stat-card">
        <div class="stat-value">{{ studentStats.totalQuestions }}</div>
        <div class="stat-label">总提问数</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ studentStats.hotPages.join('、') || '暂无' }}</div>
        <div class="stat-label">高频提问页码</div>
      </div>
      <div class="stat-card full-width">
        <div class="stat-value">{{ studentStats.keyDifficulties }}</div>
        <div class="stat-label">重点难点</div>
      </div>
    </div>
    <div v-else class="empty-tip">请先选择一个课件查看学情数据</div>
  </div>
</template>

<script setup>
defineProps({
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
    default: () => ({ totalQuestions: 0, hotPages: [], keyDifficulties: '暂无' })
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
.stats-header {
  margin-bottom: 16px;
}
.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}
.stat-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 20px;
}
.stat-card.full-width {
  grid-column: span 2;
}
.stat-value {
  font-size: 24px;
  color: #0f172a;
  font-weight: 600;
  margin-bottom: 8px;
}
.stat-label,
.empty-tip {
  color: #64748b;
}
</style>