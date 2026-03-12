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
        <div class="stat-value">{{ studentStats.activeSessions || 0 }}</div>
        <div class="stat-label">活跃追问会话</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ studentStats.hotPages.join('、') || '暂无' }}</div>
        <div class="stat-label">高频提问页码</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ studentStats.reteachCount || 0 }}</div>
        <div class="stat-label">触发重讲次数</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ studentStats.avgTurnsPerSession || 0 }}</div>
        <div class="stat-label">会话平均轮次</div>
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
    default: () => ({ totalQuestions: 0, hotPages: [], keyDifficulties: '暂无', activeSessions: 0, reteachCount: 0, avgTurnsPerSession: 0 })
  }
})
</script>

<style scoped>
.tab-content {
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.96) 0%, rgba(248, 250, 252, 0.92) 100%);
  border: 1px solid rgba(148, 163, 184, 0.18);
  border-radius: 20px;
  padding: 20px;
  box-shadow: 0 18px 34px rgba(15, 23, 42, 0.06);
}
.stats-header {
  margin-bottom: 16px;
}
.stats-header h4 {
  margin: 0;
  color: #0f172a;
  font-size: 18px;
}
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}
.stat-card {
  background: linear-gradient(180deg, rgba(255, 253, 248, 0.84) 0%, rgba(240, 249, 255, 0.82) 100%);
  border: 1px solid rgba(148, 163, 184, 0.18);
  border-radius: 18px;
  padding: 20px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.78);
}
.stat-card.full-width {
  grid-column: span 3;
}
.stat-value {
  font-size: 24px;
  color: #0f172a;
  font-weight: 700;
  margin-bottom: 8px;
}
.stat-label,
.empty-tip {
  color: #64748b;
}

@media (max-width: 960px) {
  .stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .stat-card.full-width {
    grid-column: span 2;
  }
}
</style>