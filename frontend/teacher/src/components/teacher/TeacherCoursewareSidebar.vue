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
      <div v-if="courseListLoading" class="empty-tip">课件列表加载中...</div>
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
        <div v-if="coursewareList.length === 0" class="empty-tip">暂无课件，请点击上方按钮上传</div>
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
  flex: 3;
  background: rgba(255, 255, 255, 0.96);
  border-right: 1px solid #e6ecf5;
  padding: 20px;
  overflow: auto;
}
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.section-header h3 {
  font-size: 18px;
  color: #0f172a;
}
.header-actions {
  display: flex;
  gap: 8px;
}
.publish-btn,
.upload-btn,
.del-btn,
.page-btn {
  border: none;
  border-radius: 8px;
  cursor: pointer;
}
.publish-btn,
.upload-btn {
  padding: 8px 12px;
  color: #fff;
  background: #2563eb;
}
.publish-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.upload-btn {
  background: #0ea5e9;
}
.courseware-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.course-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 14px;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  cursor: pointer;
  background: #f8fbff;
}
.course-item.active {
  border-color: #2563eb;
  box-shadow: 0 8px 20px rgba(37, 99, 235, 0.12);
}
.course-name {
  font-size: 14px;
  color: #0f172a;
  font-weight: 500;
}
.course-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}
.published-tag {
  font-size: 12px;
  color: #16a34a;
  background: #f0fdf4;
  padding: 2px 8px;
  border-radius: 999px;
}
.del-btn {
  padding: 6px 10px;
  background: #fee2e2;
  color: #dc2626;
}
.empty-tip {
  color: #64748b;
  font-size: 13px;
  text-align: center;
  padding: 20px 0;
}
.page-selector {
  margin-top: 20px;
}
.page-selector h4 {
  font-size: 14px;
  margin-bottom: 12px;
  color: #334155;
}
.page-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.page-btn {
  padding: 8px 12px;
  background: #e2e8f0;
  color: #334155;
}
.page-btn.active {
  background: #2563eb;
  color: #fff;
}
</style>