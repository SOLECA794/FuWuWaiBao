<template>
  <div class="courseware-manage-section" :class="{ 'overlay-open': isOverlayOpen, 'compact-rail': !isOverlayOpen }">
    <div class="section-header">
      <h3 v-show="isOverlayOpen">课件管理</h3>
      <div class="header-actions">
        <button v-show="isOverlayOpen" @click="$emit('open-publish')" class="publish-btn" :disabled="!currentCourseId">发布</button>
        <button v-show="isOverlayOpen" @click="$emit('open-upload')" class="upload-btn">+ 上传</button>
        <button class="toggle-right-btn" @click="isOverlayOpen = !isOverlayOpen" :title="isOverlayOpen ? '恢复侧边栏' : '展开侧边栏'">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12H3"></path><path d="M21 6H3"></path><path d="M21 18H3"></path></svg>
        </button>
      </div>
    </div>

    <div class="courseware-list" v-show="isOverlayOpen">
      <div v-if="courseListLoading" class="list-loading-tip">
        课件列表加载中...
      </div>
      <template v-else>
        <div
          v-for="course in coursewareList"
          :key="course.id"
          class="course-item"
          :class="{ active: course.id === currentCourseId }"
          @click="$emit('select-course', course)"
        >
          <div class="course-main">
            <span class="course-name" :title="course.name">{{ course.name }}</span>
            <div class="course-meta-row">
              <span class="course-type-badge" :class="fileTypeClass(course.fileType)">{{ fileTypeBadge(course.fileType) }}</span>
              <span class="course-meta-text">共{{ Number(course.knowledgePointCount) || 0 }}个知识点</span>
            </div>
          </div>
          <div class="course-actions">
            <span v-if="course.published" class="published-tag">已发布</span>
            <button @click.stop="$emit('delete-course', course.id)" class="del-btn" :disabled="courseListLoading">删除</button>
          </div>
        </div>

        <div v-if="coursewareList.length === 0" class="empty-list-tip">
          <div class="empty-icon">📂</div>
          <p>暂无课件</p>
          <span>点击上方按钮上传第一个课件</span>
        </div>
      </template>
    </div>

    <div class="page-selector" v-if="currentCourseId && isOverlayOpen">
      <div class="page-header-row">
        <h4>{{ currentCourseName }}</h4>
        <span class="page-indicator">{{ currentEditPage }} / {{ currentCourseTotalPages }}</span>
      </div>
      <div class="page-nav-controls">
        <button 
          class="nav-icon-btn" 
          :disabled="currentEditPage <= 1"
          @click="$emit('select-page', currentEditPage - 1)"
          title="上一页"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"></polyline></svg>
        </button>
        <div class="page-nav-center">切换页码</div>
        <button 
          class="nav-icon-btn" 
          :disabled="currentEditPage >= currentCourseTotalPages"
          @click="$emit('select-page', currentEditPage + 1)"
          title="下一页"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"></polyline></svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const isOverlayOpen = ref(false)

const props = defineProps({
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

const emit = defineEmits(['open-publish', 'open-upload', 'select-course', 'delete-course', 'select-page'])

const normalizedFileType = (fileType) => String(fileType || '').trim().toLowerCase()

const fileTypeBadge = (fileType) => {
  const type = normalizedFileType(fileType)
  if (type.includes('pdf')) return 'PDF'
  if (type.includes('ppt') || type.includes('pptx')) return 'PPT'
  if (type.includes('doc')) return 'DOC'
  return 'FILE'
}

const fileTypeClass = (fileType) => {
  const type = normalizedFileType(fileType)
  if (type.includes('pdf')) return 'pdf'
  if (type.includes('ppt')) return 'ppt'
  if (type.includes('doc')) return 'doc'
  return 'file'
}
</script>

<style scoped>
.courseware-manage-section {
  --ui-bg: #ffffff;
  --ui-surface: #f8fcfa;
  --ui-surface-soft: #f2f8f5;
  --ui-border: rgba(120, 156, 140, 0.22);
  --ui-border-strong: rgba(86, 130, 112, 0.34);
  --ui-text: #0f172a;
  --ui-text-muted: #5f7467;
  --ui-accent: #5ca68f;
  --ui-accent-soft: rgba(92, 166, 143, 0.16);
  --ui-shadow-soft: 0 8px 20px rgba(50, 88, 72, 0.1);
  --sidebar-rail-width: 44px;
  --sidebar-full-width: 286px;
  --sidebar-overlay-width: 360px;

  width: var(--sidebar-rail-width);
  height: 100%;
  position: relative;
  z-index: 2;
  background: linear-gradient(180deg, #f7fcf9 0%, #eef6f2 100%);
  border: 1px solid rgba(120, 156, 140, 0.28);
  border-radius: 16px;
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.08);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transform: translateX(0);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  transition: box-shadow 0.28s ease, opacity 0.22s ease, transform 0.38s cubic-bezier(0.22, 1, 0.36, 1);
}

.courseware-manage-section.overlay-open {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: var(--sidebar-overlay-width);
  z-index: 48;
  background: var(--ui-bg);
  border: 0;
  border-radius: 22px;
  box-shadow: 0 22px 48px rgba(15, 23, 42, 0.16);
  animation: drawer-slide-in 0.34s cubic-bezier(0.22, 1, 0.36, 1);
}

.section-header {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 10px 6px;
  border-bottom: 0;
  min-height: 56px;
}

.courseware-manage-section.overlay-open .section-header {
  justify-content: space-between;
  padding: 16px 14px 12px;
  border-bottom: 1px solid var(--ui-border);
}

.section-header h3 {
  font-size: 16px;
  color: var(--ui-text);
  margin: 0;
  white-space: nowrap;
}

.header-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.toggle-right-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--ui-text-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
  border-radius: 12px;
  transition: background 0.2s, color 0.2s;
}

.toggle-right-btn:hover {
  background: var(--ui-surface-soft);
  color: var(--ui-text);
}

.toggle-right-btn svg {
  width: 18px;
  height: 18px;
  transition: transform 0.28s ease;
}

.courseware-manage-section.overlay-open .toggle-right-btn svg {
  transform: rotate(90deg);
}

.publish-btn,
.upload-btn,
.del-btn,
.page-btn {
  border: none;
  border-radius: 12px;
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
  background: linear-gradient(180deg, #79c3ab 0%, #5ca68f 100%);
  box-shadow: 0 6px 14px rgba(92, 166, 143, 0.26);
}

.publish-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.upload-btn {
  background: linear-gradient(180deg, #8ecfbb 0%, #6cb59e 100%);
}

/* 课件列表：占满剩余高度，独立滚动 */
.courseware-list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 12px 12px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  transition: opacity 0.24s ease, transform 0.38s cubic-bezier(0.22, 1, 0.36, 1);
}

.collapsed-course-item {
  position: relative;
  width: 56px;
  height: 64px;
  border: 0;
  border-radius: 16px;
  background: #ffffff;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 5px;
  cursor: pointer;
  transition: border-color 0.2s ease, transform 0.2s ease, box-shadow 0.2s ease;
}

.collapsed-course-item:hover {
  border-color: transparent;
  box-shadow: var(--ui-shadow-soft);
  transform: translateY(-1px);
}

.collapsed-course-item.active {
  background: rgba(226, 244, 237, 0.88);
}

.collapsed-doc-icon {
  min-width: 28px;
  padding: 2px 6px;
  border-radius: 8px;
  text-align: center;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.4px;
  border: 1px solid #e2e8f0;
  color: #334155;
  background: #f8fafc;
}

.collapsed-initial {
  font-size: 13px;
  font-weight: 700;
  color: #334155;
  line-height: 1;
}

.collapsed-empty {
  margin-top: 10px;
  width: 42px;
  height: 42px;
  border-radius: 12px;
  border: 1px dashed #cbd5e1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
}

.collapsed-tooltip {
  position: absolute;
  left: calc(100% + 10px);
  top: 50%;
  transform: translateY(-50%) scale(0.98);
  opacity: 0;
  pointer-events: none;
  width: 250px;
  border-radius: 12px;
  border: 1px solid rgba(148, 163, 184, 0.35);
  background: rgba(15, 23, 42, 0.84);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  color: #f8fafc;
  padding: 10px 12px;
  text-align: left;
  box-shadow: 0 12px 30px rgba(2, 6, 23, 0.26);
  transition: opacity 0.2s ease, transform 0.2s ease;
  z-index: 30;
}

.collapsed-tooltip::before {
  content: '';
  position: absolute;
  left: -7px;
  top: calc(50% - 6px);
  width: 12px;
  height: 12px;
  transform: rotate(45deg);
  background: rgba(46, 55, 64, 0.94);
}

.collapsed-tooltip strong {
  display: block;
  font-size: 13px;
  font-weight: 600;
  line-height: 1.45;
  color: #ffffff;
}

.collapsed-tooltip span {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: #dbe7df;
}

.collapsed-course-item:hover .collapsed-tooltip {
  opacity: 1;
  transform: translateY(-50%) scale(1);
}

.course-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border: 0;
  border-radius: 14px;
  cursor: pointer;
  background: #ffffff;
  transition: border-color 0.2s ease, background 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
}

.course-item:hover {
  border-color: transparent;
  background: #f2f9f5;
  box-shadow: var(--ui-shadow-soft);
  transform: translateY(-1px);
}

.course-item.active {
  border-color: transparent;
  background: rgba(227, 245, 238, 0.92);
  box-shadow: 0 8px 20px rgba(92, 166, 143, 0.18);
}

.course-main {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.course-name {
  font-size: 15px;
  color: var(--ui-text);
  font-weight: 650;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 150px;
}

.course-meta-row {
  display: flex;
  align-items: center;
  gap: 6px;
}

.course-type-badge {
  min-width: 34px;
  padding: 1px 6px;
  border-radius: 999px;
  border: 1px solid #cbd5e1;
  text-align: center;
  font-size: 10px;
  font-weight: 700;
  color: #334155;
  background: #f8fafc;
}

.course-meta-text {
  font-size: 12px;
  color: var(--ui-text-muted);
}

.course-type-badge.pdf,
.collapsed-doc-icon.pdf {
  color: #b45309;
  border-color: #f5d0a8;
  background: #fff7ed;
}

.course-type-badge.ppt,
.collapsed-doc-icon.ppt {
  color: #1d4ed8;
  border-color: #bfdbfe;
  background: #eff6ff;
}

.course-type-badge.doc,
.collapsed-doc-icon.doc {
  color: #065f46;
  border-color: #a7f3d0;
  background: #ecfdf5;
}

.course-type-badge.file,
.collapsed-doc-icon.file {
  color: #334155;
  border-color: #cbd5e1;
  background: #f8fafc;
}

.course-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.published-tag {
  font-size: 11px;
  color: var(--ui-accent);
  background: var(--ui-accent-soft);
  border: 1px solid rgba(92, 166, 143, 0.35);
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
  background: var(--ui-surface);
  border: 1px dashed var(--ui-border-strong);
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
  color: var(--ui-text-muted);
  font-size: 13px;
  font-style: italic;
}

/* 页码选择器：固定展示在底部 */
.page-selector {
  flex-shrink: 0;
  padding: 14px 12px;
  border-top: 1px solid var(--ui-border);
  background: #ffffff;
  position: relative;
  transition: opacity 0.22s ease, transform 0.28s ease;
}

.page-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.page-header-row h4 {
  font-size: 13px;
  font-weight: 600;
  margin: 0;
  color: #334155;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 160px;
}

.page-indicator {
  font-size: 12px;
  font-weight: 500;
  color: var(--ui-text-muted);
  background: var(--ui-surface-soft);
  padding: 2px 8px;
  border-radius: 12px;
}

.page-nav-controls {
  display: flex;
  gap: 8px;
  align-items: center;
  justify-content: space-between;
}

.nav-icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  border: 0;
  background: var(--ui-surface-soft);
  color: var(--ui-text-muted);
  border-radius: 10px;
  cursor: pointer;
  height: 36px;
  transition: all 0.2s ease;
}

.nav-icon-btn {
  width: 44px;
}

.nav-icon-btn svg {
  width: 18px;
  height: 18px;
}

.page-nav-center {
  flex: 1;
  text-align: center;
  font-size: 12px;
  color: var(--ui-text-muted);
  letter-spacing: 0.02em;
}

.nav-icon-btn:hover:not(:disabled) {
  background: #e9f5ef;
  color: #2d5f52;
  box-shadow: 0 4px 10px rgba(42, 78, 64, 0.14);
}

.nav-icon-btn:active:not(:disabled) {
  background: #dbefe6;
}

.nav-icon-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  background: var(--ui-surface-soft);
}

@keyframes drawer-slide-in {
  from {
    transform: translateX(42px);
    opacity: 0.85;
  }

  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@media (max-width: 1280px) {
  .courseware-manage-section {
    width: 42px;
  }

  .courseware-manage-section.overlay-open {
    width: min(330px, calc(100vw - 120px));
  }
}

@media (max-width: 768px) {
  .courseware-manage-section {
    width: min(86vw, 286px);
    max-height: 100%;
    transform: translateX(0);
  }

  .courseware-manage-section.overlay-open {
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    width: min(92vw, 340px);
  }
}
</style>