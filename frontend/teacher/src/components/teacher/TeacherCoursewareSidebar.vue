<template>
  <div class="courseware-manage-section">
    <div class="section-header">
      <h3>课件管理</h3>
      <div class="header-actions">
        <button @click="$emit('open-publish')" class="publish-btn" :disabled="!currentCourseId">发布课件</button>
        <button @click="$emit('open-upload')" class="upload-btn">+ 上传课件</button>
      </div>
    </div>

    <div class="courseware-list">
      <div v-if="courseListLoading" class="list-loading-tip">课件列表加载中...</div>
      <template v-else>
        <div
          v-for="course in coursewareList"
          :key="course.id"
          class="course-item"
          :class="{ active: course.id === currentCourseId }"
          @click="$emit('select-course', course)"
        >
          <span class="course-name">{{ course.name }}</span>
          <div class="course-actions">
            <span v-if="course.published" class="published-tag">已发布</span>
            <button @click.stop="$emit('delete-course', course.id)" class="del-btn" :disabled="courseListLoading">删除</button>
          </div>
        </div>
        <div v-if="coursewareList.length === 0" class="empty-list-tip">
          <div class="empty-icon">📂</div>
          <p>暂无课件</p>
          <span>点击上方按钮上传您的第一个课件</span>
        </div>
      </template>
    </div>

    <div class="page-selector" v-if="currentCourseId">
      <h4>当前课件：{{ currentCourseName }}</h4>
      <div class="page-buttons">
        <button
          v-for="page in currentCourseTotalPages"
          :key="page"
          class="page-btn"
          :class="{ active: page === currentEditPage }"
          @click="$emit('select-page', page)"
        >
          第{{ page }}页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  coursewareList: {
    type: Array,
    default: () => []
  },
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
  currentEditPage: {
    type: Number,
    default: 1
  },
  courseListLoading: {
    type: Boolean,
    default: false
  }
})

defineEmits(['open-publish', 'open-upload', 'select-course', 'delete-course', 'select-page'])
</script>

<style scoped>
.courseware-manage-section {
  flex: 0 0 280px;
  width: 280px;
  background: #fff;
  border-right: 1px solid #e6ecf5;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.section-header {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #e6ecf5;
}

.section-header h3 {
  font-size: 16px;
  color: #0f172a;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 6px;
}

.publish-btn,
.upload-btn,
.del-btn,
.page-btn {
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: opacity 0.2s, transform 0.1s;
}

.publish-btn:active:not(:disabled),
.upload-btn:active:not(:disabled) {
  transform: scale(0.96);
}

.publish-btn,
.upload-btn {
  padding: 6px 11px;
  font-size: 13px;
  color: #fff;
  background: #2563eb;
}

.publish-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.upload-btn {
  background: #0ea5e9;
}

/* 课件列表：占满剩余高度，独立滚动 */
.courseware-list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.course-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  cursor: pointer;
  background: #f8fbff;
  transition: border-color 0.18s, background 0.18s, box-shadow 0.18s;
}

.course-item:hover {
  border-color: #93c5fd;
  background: #eff6ff;
}

.course-item.active {
  border-color: #2563eb;
  background: #eff6ff;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.12);
}

.course-name {
  font-size: 13px;
  color: #0f172a;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 140px;
}

.course-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.published-tag {
  font-size: 11px;
  color: #16a34a;
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  padding: 1px 7px;
  border-radius: 999px;
  white-space: nowrap;
}

.del-btn {
  padding: 4px 8px;
  font-size: 12px;
  background: #fee2e2;
  color: #dc2626;
  border-radius: 6px;
}

.del-btn:hover:not(:disabled) {
  background: #fecaca;
}

/* 暂无课件空状态 */
.empty-list-tip {
  padding: 36px 16px;
  text-align: center;
  background: #f8fbff;
  border: 2px dashed #e2e8f0;
  border-radius: 12px;
}

.empty-icon {
  font-size: 30px;
  margin-bottom: 10px;
}

.empty-list-tip p {
  color: #334155;
  font-weight: 500;
  margin: 0 0 4px 0;
  font-size: 14px;
}

.empty-list-tip span {
  font-size: 12px;
  color: #94a3b8;
}

.list-loading-tip {
  text-align: center;
  padding: 20px;
  color: #94a3b8;
  font-size: 13px;
  font-style: italic;
}

/* 页码选择器：固定展示在底部 */
.page-selector {
  flex-shrink: 0;
  padding: 12px 14px;
  border-top: 1px solid #e6ecf5;
  background: #f8fbff;
}

.page-selector h4 {
  font-size: 12px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.page-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  max-height: 100px;
  overflow-y: auto;
}

.page-btn {
  padding: 4px 10px;
  font-size: 12px;
  background: #e2e8f0;
  color: #334155;
  border-radius: 6px;
}

.page-btn.active {
  background: #2563eb;
  color: #fff;
}

.page-btn:hover:not(.active) {
  background: #cbd5e1;
}
</style>