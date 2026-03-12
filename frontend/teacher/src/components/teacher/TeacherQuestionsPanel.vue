<template>
  <div class="tab-content">
    <div class="questions-header" v-if="currentCourseId">
      <h4>提问统计 - {{ currentCourseName }}</h4>
      <div class="filter-bar">
        <span>按页码筛选：</span>
        <select :value="filterPage" class="page-select" @change="$emit('update:filterPage', $event.target.value)">
          <option value="">全部</option>
          <option v-for="page in currentCourseTotalPages" :key="page" :value="page">第{{ page }}页</option>
        </select>
      </div>
    </div>
    <div class="questions-list" v-if="currentCourseId">
      <div
        v-for="q in filteredQuestions"
        :key="q.id"
        class="question-item"
      >
        <div class="question-meta">
          <span class="student-id">学生 {{ q.studentId }}</span>
          <span class="page-tag">第{{ q.page }}页</span>
          <span class="time">{{ q.time }}</span>
        </div>
        <div class="question-content">{{ q.content }}</div>
        <div class="answer-content" v-if="q.answer">
          <span class="answer-label">AI 回复：</span>{{ q.answer }}
        </div>
        <div class="question-flags" v-if="q.needReteach">
          <span class="flag danger">触发重讲</span>
        </div>
      </div>
      <div v-if="filteredQuestions.length === 0" class="empty-tip">暂无提问记录</div>
    </div>
    <div v-else class="empty-tip">请先选择一个课件查看提问统计</div>
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
  currentCourseTotalPages: {
    type: Number,
    default: 0
  },
  filterPage: {
    type: [String, Number],
    default: ''
  },
  filteredQuestions: {
    type: Array,
    default: () => []
  }
})

defineEmits(['update:filterPage'])
</script>

<style scoped>
.tab-content {
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.96) 0%, rgba(248, 250, 252, 0.92) 100%);
  border: 1px solid rgba(148, 163, 184, 0.18);
  border-radius: 20px;
  padding: 20px;
  box-shadow: 0 18px 34px rgba(15, 23, 42, 0.06);
}
.questions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.questions-header h4 {
  margin: 0;
  color: #0f172a;
  font-size: 18px;
}
.filter-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #64748b;
}
.page-select {
  border: 1px solid rgba(148, 163, 184, 0.28);
  border-radius: 12px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.82);
}
.questions-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.question-item {
  background: linear-gradient(180deg, rgba(255, 253, 248, 0.84) 0%, rgba(240, 249, 255, 0.82) 100%);
  border: 1px solid rgba(148, 163, 184, 0.18);
  border-radius: 18px;
  padding: 14px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.78);
}
.question-meta {
  display: flex;
  gap: 12px;
  color: #64748b;
  font-size: 12px;
  margin-bottom: 8px;
}
.question-content {
  color: #0f172a;
  line-height: 1.7;
}
.answer-content {
  margin-top: 8px;
  color: #334155;
  line-height: 1.6;
}
.question-flags {
  margin-top: 10px;
}
.flag {
  display: inline-flex;
  align-items: center;
  padding: 3px 8px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
}
.flag.danger {
  background: #fff1f2;
  color: #be123c;
}
.answer-label {
  color: #2563eb;
  font-weight: 600;
}
.empty-tip {
  text-align: center;
  color: #64748b;
}
</style>