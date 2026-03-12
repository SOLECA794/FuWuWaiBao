<template>
  <div class="courseware-manage-section">
    <div class="section-header">
      <div>
        <div class="header-caption">课程工作台</div>
        <h3>课件管理</h3>
      </div>
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
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.96) 0%, rgba(247, 250, 252, 0.98) 100%);
  border-right: 1px solid rgba(226, 232, 240, 0.92);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 14px 0 30px rgba(15, 23, 42, 0.05);
}

.section-header {
  flex-shrink: 0;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 18px 18px 16px;
  border-bottom: 1px solid rgba(226, 232, 240, 0.92);
  background: linear-gradient(180deg, rgba(240, 249, 255, 0.9) 0%, rgba(255, 255, 255, 0.72) 100%);
}

.section-header h3 {
  font-size: 18px;
  color: #0f172a;
  margin: 4px 0 0;
}

.header-caption {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #0f766e;
}

.header-actions {
  display: flex;
  gap: 6px;
  flex-direction: column;
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
  padding: 8px 12px;
  font-size: 12px;
  font-weight: 700;
  color: #fff;
  background: linear-gradient(90deg, #0f766e 0%, #0284c7 100%);
  box-shadow: 0 10px 20px rgba(2, 132, 199, 0.18);
}

.publish-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.upload-btn {
  background: linear-gradient(90deg, #0ea5e9 0%, #2563eb 100%);
}

/* 课件列表：占满剩余高度，独立滚动 */
.courseware-list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.course-item {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: flex-start;
  gap: 10px;
  padding: 14px 14px 12px;
  border: 1px solid rgba(226, 232, 240, 0.9);
  border-radius: 16px;
  cursor: pointer;
  background: linear-gradient(180deg, rgba(248, 250, 252, 0.92) 0%, rgba(255, 255, 255, 0.98) 100%);
  transition: border-color 0.18s, background 0.18s, box-shadow 0.18s, transform 0.18s;
}

.course-item:hover {
  border-color: rgba(56, 189, 248, 0.5);
  background: linear-gradient(180deg, rgba(240, 249, 255, 0.95) 0%, rgba(255, 255, 255, 0.98) 100%);
  transform: translateY(-1px);
}

.course-item.active {
  border-color: #0ea5e9;
  background: linear-gradient(180deg, rgba(224, 242, 254, 0.98) 0%, rgba(255, 255, 255, 0.98) 100%);
  box-shadow: 0 14px 28px rgba(14, 165, 233, 0.14);
}

.course-name {
  font-size: 14px;
  color: #0f172a;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.course-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
  flex-shrink: 0;
  width: 100%;
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
  border-radius: 999px;
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
  padding: 14px;
  border-top: 1px solid rgba(226, 232, 240, 0.92);
  background: linear-gradient(180deg, rgba(248, 250, 252, 0.98) 0%, rgba(239, 246, 255, 0.95) 100%);
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
  padding: 6px 10px;
  font-size: 12px;
  background: #e2e8f0;
  color: #334155;
  border-radius: 999px;
}

.page-btn.active {
  background: linear-gradient(90deg, #0284c7 0%, #2563eb 100%);
  color: #fff;
}

.page-btn:hover:not(.active) {
  background: #cbd5e1;
}
</style>