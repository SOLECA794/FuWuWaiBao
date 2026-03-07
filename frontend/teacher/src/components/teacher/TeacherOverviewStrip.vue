<template>
  <div class="overview-strip">
    <div class="overview-item">
      <span class="overview-label">当前课件</span>
      <span class="overview-value">{{ currentCourseName || '未选择' }}</span>
    </div>
    <div class="overview-item">
      <span class="overview-label">当前页</span>
      <span class="overview-value">第{{ currentEditPage }}页 / 共{{ currentCourseTotalPages || 0 }}页</span>
    </div>
    <div class="overview-item">
      <span class="overview-label">发布状态</span>
      <span class="overview-value" :class="currentCoursePublished ? 'success' : 'pending'">
        {{ currentCoursePublished ? '已发布' : '未发布' }}
      </span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  currentCourseName: {
    type: String,
    default: ''
  },
  currentEditPage: {
    type: Number,
    default: 1
  },
  currentCourseTotalPages: {
    type: Number,
    default: 0
  },
  currentCoursePublished: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
.overview-strip {
  height: 48px;
  background: rgba(255, 255, 255, 0.72);
  border-bottom: 1px solid #e6ecf5;
  display: flex;
  align-items: center;
  gap: 28px;
  padding: 0 24px;
}
.overview-item {
  display: flex;
  align-items: center;
  gap: 10px;
}
.overview-label {
  color: #64748b;
  font-size: 13px;
}
.overview-value {
  font-size: 14px;
  color: #0f172a;
  font-weight: 500;
}
.overview-value.success {
  color: #16a34a;
}
.overview-value.pending {
  color: #d97706;
}
</style>