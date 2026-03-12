<template>
  <div class="overview-strip">
    <div class="overview-item">
      <span class="overview-label">当前课件</span>
      <span class="overview-value course-name" :title="currentCourseName || '未选择'">
        {{ currentCourseName || '未选择' }}
      </span>
    </div>
    <div class="overview-divider"></div>
    <div class="overview-item">
      <span class="overview-label">页码</span>
      <span class="overview-value">第 {{ currentEditPage }} 页 / 共 {{ currentCourseTotalPages || 0 }} 页</span>
    </div>
    <div class="overview-divider"></div>
    <div class="overview-item">
      <span class="overview-label">发布状态</span>
      <span class="overview-value status-val" :class="currentCoursePublished ? 'success' : 'pending'">
        <span class="status-dot"></span>
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
  height: 44px;
  background: #fff;
  border-bottom: 1px solid #e6ecf5;
  display: flex;
  align-items: center;
  padding: 0 20px;
  gap: 0;
  flex-shrink: 0;
}

.overview-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 20px;
}

.overview-item:first-child {
  padding-left: 4px;
}

.overview-divider {
  width: 1px;
  height: 18px;
  background: #e6ecf5;
  flex-shrink: 0;
}

.overview-label {
  color: #94a3b8;
  font-size: 12px;
  white-space: nowrap;
}

.overview-value {
  font-size: 13px;
  color: #1e293b;
  font-weight: 500;
  white-space: nowrap;
}

.overview-value.course-name {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.status-val {
  display: inline-flex;
  align-items: center;
  gap: 5px;
}

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}

.overview-value.success {
  color: #2F605A;
}

.overview-value.success .status-dot {
  background: #2F605A;
}

.overview-value.pending {
  color: #d97706;
}

.overview-value.pending .status-dot {
  background: #d97706;
}
</style>