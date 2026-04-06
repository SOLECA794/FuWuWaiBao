<template>
  <div class="courseware-manage-section" :class="{ 'collapsed': isCollapsed }">
    <div class="section-header">
      <h3 v-show="!isCollapsed">课件管理</h3>
      <div class="header-actions">
        <button v-show="!isCollapsed" @click="$emit('open-publish')" class="publish-btn" :disabled="!currentCourseId">发布</button>
        <button v-show="!isCollapsed" @click="$emit('open-upload')" class="upload-btn">+ 上传</button>
        <button class="toggle-right-btn" @click="isCollapsed = !isCollapsed" :title="isCollapsed ? '展开侧边栏' : '收起侧边栏'">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12H3"></path><path d="M21 6H3"></path><path d="M21 18H3"></path></svg>
        </button>
      </div>
    </div>

    <div class="courseware-list" v-show="!isCollapsed">
      <div v-if="courseListLoading" class="list-loading-tip">课件列表加载中...</div>
      <template v-else>
        <div
          v-for="course in coursewareList"
          :key="course.id"
          class="course-item"
          :class="{ active: course.id === currentCourseId }"
          @click="$emit('select-course', course)"
        >
          <span class="course-name" :title="course.name">{{ course.name }}</span>
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

    <div class="page-selector" v-if="currentCourseId" v-show="!isCollapsed">
      <div class="page-header-row">
        <h4>{{ currentCourseName }}</h4>
        <div class="page-header-meta">
          <span class="page-indicator">{{ currentEditPage }} / {{ currentCourseTotalPages }}</span>
          <button
            v-if="currentCourseId"
            type="button"
            class="ref-health-btn"
            title="扫描各表对讲授节点 node_id 的引用是否与当前节点一致"
            @click="runReferenceHealth"
          >
            引用健康
          </button>
        </div>
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
        <button 
          class="catalog-btn" 
          @click="showCatalog = !showCatalog"
          :class="{ active: showCatalog }"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="8" y1="6" x2="21" y2="6"></line><line x1="8" y1="12" x2="21" y2="12"></line><line x1="8" y1="18" x2="21" y2="18"></line><line x1="3" y1="6" x2="3.01" y2="6"></line><line x1="3" y1="12" x2="3.01" y2="12"></line><line x1="3" y1="18" x2="3.01" y2="18"></line></svg>
          目录
        </button>
        <button 
          class="nav-icon-btn" 
          :disabled="currentEditPage >= currentCourseTotalPages"
          @click="$emit('select-page', currentEditPage + 1)"
          title="下一页"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"></polyline></svg>
        </button>
      </div>

      <!-- 弹出式目录选择栏 -->
      <transition name="slide-up">
        <div class="catalog-popup" v-if="showCatalog">
          <div class="catalog-header">
            <span>选择页码</span>
            <button class="close-catalog" @click="showCatalog = false">&times;</button>
          </div>
          <div class="catalog-content">
            <button
              v-for="page in currentCourseTotalPages"
              :key="page"
              class="catalog-page-btn"
              :class="{ active: page === currentEditPage }"
              @click="handleSelectPage(page)"
            >
              第{{ page }}页
            </button>
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { teacherCoursewareApi } from '../../services/v1/coursewareApi'

const isCollapsed = ref(window.innerWidth <= 1200)
const showCatalog = ref(false)

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

watch(() => props.currentCourseId, () => {
  showCatalog.value = false
})

const handleSelectPage = (page) => {
  emit('select-page', page)
  showCatalog.value = false
}

const refHealthBusy = ref(false)

async function runReferenceHealth() {
  if (!props.currentCourseId || refHealthBusy.value) return
  refHealthBusy.value = true
  try {
    const res = await teacherCoursewareApi.getKnowledgeGraphReferenceHealth(props.currentCourseId)
    const d = res.data || {}
    const orphans = d.unionOrphanNodeIds || []
    if (!d.hasOrphans) {
      window.alert('未发现孤儿节点引用（各表 node 与当前 teaching_nodes 一致）。')
      return
    }
    const preview = orphans.length > 12
      ? `${orphans.slice(0, 12).join(', ')} … 共 ${orphans.length} 个`
      : orphans.join(', ')
    const ok = window.confirm(
      `发现不在讲授节点中的引用：${preview}\n\n是否执行修复？将清空相关 node 字段、软删除指向孤儿节点的收藏，并重建图谱边。`
    )
    if (!ok) return
    const fix = await teacherCoursewareApi.repairKnowledgeGraphReferences(props.currentCourseId, { confirm: true })
    window.alert(
      fix.message +
        '\n\n' +
        JSON.stringify(fix.data || fix, null, 2).slice(0, 3500)
    )
  } catch (e) {
    window.alert('引用健康检查失败：' + (e?.message || e))
  } finally {
    refHealthBusy.value = false
  }
}
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
  transition: flex-basis 0.3s ease, width 0.3s ease;
}

.courseware-manage-section.collapsed {
  flex: 0 0 56px;
  width: 56px;
}

.section-header {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #e6ecf5;
  min-height: 56px;
}

.courseware-manage-section.collapsed .section-header {
  padding: 16px 0;
  justify-content: center;
}

.section-header h3 {
  font-size: 16px;
  color: #0f172a;
  margin: 0;
  white-space: nowrap;
}

.header-actions {
  display: flex;
  gap: 6px;
  align-items: center;
}

.toggle-right-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  color: #94a3b8;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
  border-radius: 4px;
  transition: background 0.2s, color 0.2s;
}

.toggle-right-btn:hover {
  background: #f1f5f9;
  color: #1e293b;
}

.toggle-right-btn svg {
  width: 18px;
  height: 18px;
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
  background: #2F605A;
}

.publish-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.upload-btn {
  background: #356F68;
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
  background: #ffffff;
  transition: border-color 0.18s, background 0.18s, box-shadow 0.18s;
}

.course-item:hover {
  border-color: #8FC1B5;
  background: #F4F7F7;
}

.course-item.active {
  border-color: #2F605A;
  background: #F4F7F7;
  box-shadow: 0 4px 12px rgba(47, 96, 90, 0.12);
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
  color: #2F605A;
  background: #E8F0EF;
  border: 1px solid #8FC1B5;
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
  background: #ffffff;
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
  padding: 16px;
  border-top: 1px solid #e6ecf5;
  background: #ffffff;
  position: relative;
}

.page-header-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 6px;
}

.ref-health-btn {
  padding: 4px 8px;
  font-size: 11px;
  border-radius: 6px;
  border: 1px solid #cbd5e1;
  background: #f8fafc;
  color: #334155;
  cursor: pointer;
}

.ref-health-btn:hover {
  border-color: #2F605A;
  color: #1e293b;
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
  color: #64748b;
  background: #f1f5f9;
  padding: 2px 8px;
  border-radius: 12px;
}

.page-nav-controls {
  display: flex;
  gap: 8px;
  align-items: center;
  justify-content: space-between;
}

.nav-icon-btn, .catalog-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #e2e8f0;
  background: #fff;
  color: #64748b;
  border-radius: 8px;
  cursor: pointer;
  height: 36px;
  transition: all 0.2s;
}

.nav-icon-btn {
  width: 44px;
}

.catalog-btn {
  flex: 1;
  gap: 8px;
  font-size: 13px;
  font-weight: 500;
}

.nav-icon-btn svg, .catalog-btn svg {
  width: 18px;
  height: 18px;
}

.nav-icon-btn:hover:not(:disabled), .catalog-btn:hover:not(.active) {
  background: #f8fafc;
  border-color: #cbd5e1;
  color: #334155;
}

.nav-icon-btn:active:not(:disabled), .catalog-btn:active {
  background: #f1f5f9;
}

.nav-icon-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  background: #f8fafc;
}

.catalog-btn.active {
  background: #2F605A;
  color: #fff;
  border-color: #2F605A;
}

/* 弹出式目录 */
.catalog-popup {
  position: absolute;
  bottom: 100%;
  left: 0;
  right: 0;
  background: #fff;
  border-top: 1px solid #e6ecf5;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.08);
  border-radius: 16px 16px 0 0;
  display: flex;
  flex-direction: column;
  max-height: 350px;
  z-index: 10;
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(100%);
  opacity: 0;
}

.catalog-header {
  padding: 12px 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #f1f5f9;
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
}

.close-catalog {
  background: transparent;
  border: none;
  font-size: 20px;
  line-height: 1;
  color: #94a3b8;
  cursor: pointer;
  padding: 0 4px;
}

.close-catalog:hover {
  color: #475569;
}

.catalog-content {
  padding: 12px 16px;
  overflow-y: auto;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.catalog-page-btn {
  padding: 8px 0;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  color: #475569;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.catalog-page-btn.active {
  background: #2F605A;
  color: #fff;
  border-color: #2F605A;
  font-weight: 500;
}

.catalog-page-btn:hover:not(.active) {
  background: #e2e8f0;
}
</style>