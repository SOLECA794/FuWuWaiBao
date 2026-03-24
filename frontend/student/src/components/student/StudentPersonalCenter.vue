<template>
  <div class="personal-center">
    <div class="center-header">
      <h2>个人学习中心</h2>
      <div class="tab-buttons">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          :class="['tab-btn', { active: activeTab === tab.key }]"
          @click="activeTab = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <div class="center-content">
      <!-- 提醒通知 -->
      <div v-if="dueReminders.length > 0" class="reminders-section">
        <h3>📅 复习提醒</h3>
        <div class="reminders-list">
          <div v-for="reminder in dueReminders" :key="reminder.plan.id" class="reminder-item">
            <span>{{ reminder.message }}</span>
            <button @click="viewPlanItems(reminder.plan)">去复习</button>
          </div>
        </div>
      </div>
      <!-- 笔记管理 -->
      <div v-if="activeTab === 'notes'" class="tab-content">
        <div class="section-header">
          <h3>我的笔记</h3>
          <button class="add-btn" @click="showNoteDialog = true">添加笔记</button>
        </div>
        <div class="notes-list">
          <div v-for="note in notes" :key="note.id" class="note-item">
            <div class="note-header">
              <span class="note-title">{{ note.title || `第${note.pageNum}页笔记` }}</span>
              <div class="note-actions">
                <button @click="editNote(note)">编辑</button>
                <button @click="deleteNote(note.id)">删除</button>
              </div>
            </div>
            <div class="note-content">{{ note.note }}</div>
            <div class="note-meta">
              <span>课程: {{ getCourseName(note.courseId) }}</span>
              <span>页码: {{ note.pageNum }}</span>
              <span>创建时间: {{ formatDate(note.createdAt) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 收藏管理 -->
      <div v-if="activeTab === 'favorites'" class="tab-content">
        <div class="section-header">
          <h3>我的收藏</h3>
          <div class="filter-tags">
            <button
              v-for="tag in availableTags"
              :key="tag"
              :class="['tag-filter', { active: selectedTag === tag }]"
              @click="selectedTag = selectedTag === tag ? '' : tag"
            >
              {{ tag }}
            </button>
          </div>
        </div>
        <div class="favorites-list">
          <div v-for="fav in filteredFavorites" :key="fav.id" class="favorite-item">
            <div class="favorite-header">
              <span class="favorite-title">{{ fav.title }}</span>
              <div class="favorite-actions">
                <button @click="editFavorite(fav)">编辑标签</button>
                <button @click="deleteFavorite(fav.id)">删除</button>
              </div>
            </div>
            <div class="favorite-tags">
              <span v-for="tag in fav.tags" :key="tag" class="tag">{{ tag }}</span>
            </div>
            <div class="favorite-meta">
              <span>课程: {{ getCourseName(fav.courseId) }}</span>
              <span>节点: {{ fav.nodeId }}</span>
              <span>页码: {{ fav.pageNum }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 复习计划 -->
      <div v-if="activeTab === 'plans'" class="tab-content">
        <div class="section-header">
          <h3>复习计划</h3>
          <button class="add-btn" @click="showPlanDialog = true">创建计划</button>
        </div>
        <div class="plans-list">
          <div v-for="plan in reviewPlans" :key="plan.id" class="plan-item">
            <div class="plan-header">
              <h4>{{ plan.name }}</h4>
              <div class="plan-status" :class="plan.status">
                {{ plan.status === 'active' ? '进行中' : plan.status === 'paused' ? '暂停' : '完成' }}
              </div>
            </div>
            <p class="plan-desc">{{ plan.description }}</p>
            <div class="plan-meta">
              <span>频率: {{ plan.frequency === 'daily' ? '每日' : plan.frequency === 'weekly' ? '每周' : '每月' }}</span>
              <span>下次复习: {{ plan.nextReviewDate ? formatDate(plan.nextReviewDate) : '未设置' }}</span>
            </div>
            <div class="plan-actions">
              <button @click="viewPlanItems(plan)">查看内容</button>
              <button @click="editPlan(plan)">编辑</button>
              <button @click="deletePlan(plan.id)">删除</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 笔记对话框 -->
    <div v-if="showNoteDialog" class="dialog-overlay" @click="showNoteDialog = false">
      <div class="dialog" @click.stop>
        <h3>{{ editingNote ? '编辑笔记' : '添加笔记' }}</h3>
        <form @submit.prevent="saveNote">
          <div class="form-group">
            <label>课程</label>
            <select v-model="noteForm.courseId" required>
              <option v-for="course in courses" :key="course.id" :value="course.id">
                {{ course.title }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>页码</label>
            <input v-model.number="noteForm.pageNum" type="number" min="1" required>
          </div>
          <div class="form-group">
            <label>内容</label>
            <textarea v-model="noteForm.content" required></textarea>
          </div>
          <div class="dialog-actions">
            <button type="button" @click="showNoteDialog = false">取消</button>
            <button type="submit">保存</button>
          </div>
        </form>
      </div>
    </div>

    <!-- 收藏标签编辑对话框 -->
    <div v-if="showTagDialog" class="dialog-overlay" @click="showTagDialog = false">
      <div class="dialog" @click.stop>
        <h3>编辑标签</h3>
        <div class="tag-editor">
          <div class="current-tags">
            <span v-for="tag in editingFavorite.tags" :key="tag" class="tag removable" @click="removeTag(tag)">
              {{ tag }} ×
            </span>
          </div>
          <input
            v-model="newTag"
            @keyup.enter="addTag"
            placeholder="输入新标签，按回车添加"
          >
        </div>
        <div class="dialog-actions">
          <button @click="showTagDialog = false">取消</button>
          <button @click="saveTags">保存</button>
        </div>
      </div>
    </div>

    <!-- 复习计划对话框 -->
    <div v-if="showPlanDialog" class="dialog-overlay" @click="showPlanDialog = false">
      <div class="dialog" @click.stop>
        <h3>{{ editingPlan ? '编辑计划' : '创建复习计划' }}</h3>
        <form @submit.prevent="savePlan">
          <div class="form-group">
            <label>计划名称</label>
            <input v-model="planForm.name" required>
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea v-model="planForm.description"></textarea>
          </div>
          <div class="form-group">
            <label>复习频率</label>
            <select v-model="planForm.frequency" required>
              <option value="daily">每日</option>
              <option value="weekly">每周</option>
              <option value="monthly">每月</option>
            </select>
          </div>
          <div class="dialog-actions">
            <button type="button" @click="showPlanDialog = false">取消</button>
            <button type="submit">保存</button>
          </div>
        </form>
      </div>
    </div>

    <!-- 计划内容对话框 -->
    <div v-if="showPlanItemsDialog" class="dialog-overlay" @click="showPlanItemsDialog = false">
      <div class="dialog large" @click.stop>
        <h3>{{ currentPlan.name }} - 复习内容</h3>
        <div class="plan-items">
          <div v-for="item in planItems" :key="item.id" class="plan-item-detail">
            <div class="item-header">
              <span class="item-type">{{ item.itemType === 'note' ? '笔记' : '收藏' }}</span>
              <div class="item-actions">
                <button @click="markReviewed(item)">标记已复习</button>
                <button @click="removePlanItem(item.id)">移除</button>
              </div>
            </div>
            <div class="item-content">
              <p v-if="item.itemType === 'note'">{{ getNoteContent(item.itemId) }}</p>
              <p v-else>{{ getFavoriteTitle(item.itemId) }}</p>
            </div>
            <div class="item-meta">
              <span>优先级: {{ item.priority }}</span>
              <span>复习次数: {{ item.reviewCount }}</span>
              <span>上次复习: {{ item.lastReviewedAt ? formatDate(item.lastReviewedAt) : '未复习' }}</span>
            </div>
          </div>
        </div>
        <div class="add-items-section">
          <h4>添加复习内容</h4>
          <div class="add-item-form">
            <select v-model="addItemForm.type">
              <option value="note">笔记</option>
              <option value="favorite">收藏</option>
            </select>
            <select v-model="addItemForm.itemId">
              <option v-for="item in availablePlanItems" :key="item.id" :value="item.id">
                {{ addItemForm.type === 'note' ? `第${item.pageNum}页笔记` : item.title }}
              </option>
            </select>
            <input v-model.number="addItemForm.priority" type="number" min="1" max="5" placeholder="优先级">
            <button @click="addPlanItem">添加</button>
          </div>
        </div>
        <div class="dialog-actions">
          <button @click="showPlanItemsDialog = false">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { studentCoursewareApi } from '@/services/v1/coursewareApi'

export default {
  name: 'StudentPersonalCenter',
  props: {
    studentId: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      activeTab: 'notes',
      tabs: [
        { key: 'notes', label: '笔记' },
        { key: 'favorites', label: '收藏' },
        { key: 'plans', label: '复习计划' }
      ],
      notes: [],
      favorites: [],
      reviewPlans: [],
      courses: [],
      availableTags: [],
      selectedTag: '',
      showNoteDialog: false,
      showTagDialog: false,
      showPlanDialog: false,
      showPlanItemsDialog: false,
      editingNote: null,
      editingFavorite: null,
      editingPlan: null,
      currentPlan: null,
      planItems: [],
      noteForm: {
        courseId: '',
        pageNum: 1,
        content: ''
      },
      planForm: {
        name: '',
        description: '',
        frequency: 'daily'
      },
      addItemForm: {
        type: 'note',
        itemId: '',
        priority: 3
      },
      newTag: ''
    }
  },
  computed: {
    filteredFavorites() {
      if (!this.selectedTag) return this.favorites
      return this.favorites.filter(fav => fav.tags && fav.tags.includes(this.selectedTag))
    },
    availablePlanItems() {
      if (this.addItemForm.type === 'note') {
        return this.notes
      } else if (this.addItemForm.type === 'favorite') {
        return this.favorites
      }
      return []
    },
    dueReminders() {
      const now = new Date()
      const dueItems = []
      this.reviewPlans.forEach(plan => {
        if (plan.status === 'active' && plan.nextReviewDate) {
          const nextReview = new Date(plan.nextReviewDate)
          if (nextReview <= now) {
            dueItems.push({
              type: 'plan',
              plan: plan,
              message: `复习计划 "${plan.name}" 到期`
            })
          }
        }
      })
      return dueItems
    }
  },
  async created() {
    await this.loadData()
  },
  methods: {
    async loadData() {
      try {
        // Load courses
        const courseRes = await studentCoursewareApi.list()
        this.courses = courseRes.data || []

        // Load notes
        const notesRes = await studentCoursewareApi.listNotes({
          studentId: this.studentId,
          pageSize: 100
        })
        this.notes = notesRes.data.items || []

        // Load favorites
        const favRes = await studentCoursewareApi.listFavorites({
          studentId: this.studentId,
          pageSize: 100
        })
        this.favorites = (favRes.data.items || []).map(fav => ({
          ...fav,
          tags: fav.tags ? JSON.parse(fav.tags) : []
        }))

        // Extract available tags
        this.availableTags = [...new Set(this.favorites.flatMap(fav => fav.tags))]

        // Load review plans
        const plansRes = await studentCoursewareApi.listReviewPlans(this.studentId)
        this.reviewPlans = plansRes.data || []

      } catch (error) {
        console.error('加载数据失败:', error)
      }
    },

    // 笔记相关方法
    editNote(note) {
      this.editingNote = note
      this.noteForm = {
        courseId: note.courseId,
        pageNum: note.pageNum,
        content: note.note
      }
      this.showNoteDialog = true
    },

    async saveNote() {
      try {
        if (this.editingNote) {
          // Update existing note - assuming we have an update API
          // For now, delete and recreate
          await this.deleteNote(this.editingNote.id)
        }
        await studentCoursewareApi.saveNote({
          studentId: this.studentId,
          courseId: this.noteForm.courseId,
          pageNum: this.noteForm.pageNum,
          content: this.noteForm.content
        })
        this.showNoteDialog = false
        this.editingNote = null
        this.noteForm = { courseId: '', pageNum: 1, content: '' }
        await this.loadData()
      } catch (error) {
        console.error('保存笔记失败:', error)
      }
    },

    async deleteNote(noteId) {
      if (confirm('确定要删除这条笔记吗？')) {
        try {
          await studentCoursewareApi.deleteNote(noteId, this.studentId)
          await this.loadData()
        } catch (error) {
          console.error('删除笔记失败:', error)
        }
      }
    },

    // 收藏相关方法
    editFavorite(fav) {
      this.editingFavorite = { ...fav }
      this.showTagDialog = true
    },

    addTag() {
      if (this.newTag && !this.editingFavorite.tags.includes(this.newTag)) {
        this.editingFavorite.tags.push(this.newTag)
        this.newTag = ''
      }
    },

    removeTag(tag) {
      const index = this.editingFavorite.tags.indexOf(tag)
      if (index > -1) {
        this.editingFavorite.tags.splice(index, 1)
      }
    },

    async saveTags() {
      try {
        await studentCoursewareApi.addFavorite({
          studentId: this.studentId,
          courseId: this.editingFavorite.courseId,
          nodeId: this.editingFavorite.nodeId,
          pageNum: this.editingFavorite.pageNum,
          title: this.editingFavorite.title,
          tags: this.editingFavorite.tags
        })
        this.showTagDialog = false
        this.editingFavorite = null
        await this.loadData()
      } catch (error) {
        console.error('保存标签失败:', error)
      }
    },

    async deleteFavorite(favId) {
      try {
        await studentCoursewareApi.deleteFavorite(favId)
        await this.loadData()
      } catch (error) {
        console.error('删除收藏失败:', error)
      }
    },

    // 复习计划相关方法
    editPlan(plan) {
      this.editingPlan = plan
      this.planForm = {
        name: plan.name,
        description: plan.description,
        frequency: plan.frequency
      }
      this.showPlanDialog = true
    },

    async savePlan() {
      try {
        if (this.editingPlan) {
          await studentCoursewareApi.updateReviewPlan(this.editingPlan.id, this.planForm)
        } else {
          await studentCoursewareApi.createReviewPlan({
            studentId: this.studentId,
            ...this.planForm
          })
        }
        this.showPlanDialog = false
        this.editingPlan = null
        this.planForm = { name: '', description: '', frequency: 'daily' }
        await this.loadData()
      } catch (error) {
        console.error('保存计划失败:', error)
      }
    },

    async deletePlan(planId) {
      try {
        await studentCoursewareApi.deleteReviewPlan(planId)
        await this.loadData()
      } catch (error) {
        console.error('删除计划失败:', error)
      }
    },

    async viewPlanItems(plan) {
      this.currentPlan = plan
      try {
        const res = await studentCoursewareApi.listReviewPlanItems(plan.id)
        this.planItems = res.data || []
        this.showPlanItemsDialog = true
      } catch (error) {
        console.error('加载计划内容失败:', error)
      }
    },

    async addPlanItem() {
      try {
        await studentCoursewareApi.addReviewPlanItem({
          reviewPlanId: this.currentPlan.id,
          itemType: this.addItemForm.type,
          itemId: this.addItemForm.itemId,
          priority: this.addItemForm.priority
        })
        this.addItemForm = { type: 'note', itemId: '', priority: 3 }
        await this.viewPlanItems(this.currentPlan)
      } catch (error) {
        console.error('添加复习项失败:', error)
      }
    },

    async markReviewed(item) {
      try {
        await studentCoursewareApi.updateReviewPlanItem(item.id, {
          lastReviewedAt: new Date().toISOString(),
          reviewCount: item.reviewCount + 1
        })
        await this.viewPlanItems(this.currentPlan)
      } catch (error) {
        console.error('标记复习失败:', error)
      }
    },

    async removePlanItem(itemId) {
      try {
        await studentCoursewareApi.deleteReviewPlanItem(itemId)
        await this.viewPlanItems(this.currentPlan)
      } catch (error) {
        console.error('移除复习项失败:', error)
      }
    },

    // 辅助方法
    getCourseName(courseId) {
      const course = this.courses.find(c => c.id === courseId)
      return course ? course.title : courseId
    },

    getNoteContent(noteId) {
      const note = this.notes.find(n => n.id === noteId)
      return note ? note.note : '笔记内容'
    },

    getFavoriteTitle(favId) {
      const fav = this.favorites.find(f => f.id === favId)
      return fav ? fav.title : '收藏标题'
    },

    formatDate(dateStr) {
      return new Date(dateStr).toLocaleDateString('zh-CN')
    }
  }
}
</script>

<style scoped>
.personal-center {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.center-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.center-header h2 {
  margin: 0;
  color: #333;
}

.tab-buttons {
  display: flex;
  gap: 10px;
}

.tab-btn {
  padding: 10px 20px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.3s;
}

.tab-btn.active {
  background: #007bff;
  color: white;
  border-color: #007bff;
}

.center-content {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.reminders-section {
  background: #fff3cd;
  border: 1px solid #ffeaa7;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 20px;
}

.reminders-section h3 {
  margin-top: 0;
  color: #856404;
}

.reminders-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.reminder-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  padding: 10px;
  border-radius: 6px;
  border-left: 4px solid #ffc107;
}

.reminder-item button {
  background: #ffc107;
  color: #212529;
  border: none;
  padding: 5px 10px;
  border-radius: 4px;
  cursor: pointer;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h3 {
  margin: 0;
  color: #333;
}

.add-btn {
  padding: 8px 16px;
  background: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.notes-list, .favorites-list, .plans-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.note-item, .favorite-item, .plan-item {
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 15px;
  background: #fafafa;
}

.note-header, .favorite-header, .plan-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.note-title, .favorite-title {
  font-weight: bold;
  color: #333;
}

.note-actions, .favorite-actions, .plan-actions {
  display: flex;
  gap: 10px;
}

.note-actions button, .favorite-actions button, .plan-actions button {
  padding: 5px 10px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
}

.note-content {
  margin-bottom: 10px;
  line-height: 1.5;
}

.note-meta, .favorite-meta, .plan-meta {
  display: flex;
  gap: 15px;
  font-size: 0.9em;
  color: #666;
}

.favorite-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  margin-bottom: 10px;
}

.tag {
  background: #e9ecef;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.8em;
}

.tag.removable {
  cursor: pointer;
}

.tag.removable:hover {
  background: #dc3545;
  color: white;
}

.filter-tags {
  display: flex;
  gap: 5px;
}

.tag-filter {
  padding: 5px 10px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 15px;
  cursor: pointer;
  font-size: 0.9em;
}

.tag-filter.active {
  background: #007bff;
  color: white;
}

.plan-status {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.8em;
  font-weight: bold;
}

.plan-status.active {
  background: #d4edda;
  color: #155724;
}

.plan-status.paused {
  background: #fff3cd;
  color: #856404;
}

.plan-status.completed {
  background: #d1ecf1;
  color: #0c5460;
}

.plan-desc {
  margin: 10px 0;
  color: #666;
}

/* 对话框样式 */
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  background: white;
  border-radius: 8px;
  padding: 20px;
  max-width: 500px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
}

.dialog.large {
  max-width: 800px;
}

.dialog h3 {
  margin-top: 0;
  color: #333;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-group textarea {
  min-height: 100px;
  resize: vertical;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.dialog-actions button {
  padding: 8px 16px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
}

.dialog-actions button[type="submit"] {
  background: #007bff;
  color: white;
  border-color: #007bff;
}

.tag-editor {
  margin: 20px 0;
}

.current-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  margin-bottom: 10px;
}

.plan-items {
  margin: 20px 0;
}

.plan-item-detail {
  border: 1px solid #eee;
  border-radius: 6px;
  padding: 10px;
  margin-bottom: 10px;
  background: #f8f9fa;
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.item-type {
  font-weight: bold;
  color: #007bff;
}

.item-content {
  margin-bottom: 8px;
}

.item-meta {
  display: flex;
  gap: 15px;
  font-size: 0.9em;
  color: #666;
}

.add-items-section {
  border-top: 1px solid #eee;
  padding-top: 20px;
  margin-top: 20px;
}

.add-items-section h4 {
  margin-bottom: 15px;
  color: #333;
}

.add-item-form {
  display: flex;
  gap: 10px;
  align-items: center;
}

.add-item-form select,
.add-item-form input {
  padding: 6px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.add-item-form button {
  padding: 6px 12px;
  background: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
</style>