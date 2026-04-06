<template>
  <div class="personal-center">
    <div class="header">
      <div>
        <h2>个人学习中心</h2>
        <p>围绕学生学习记录、复习提醒、练习与通知做集中管理。</p>
      </div>
      <div class="tabs">
        <button v-for="tab in tabs" :key="tab.key" :class="['tab', { active: activeTab === tab.key }]" @click="activeTab = tab.key">
          {{ tab.label }}
        </button>
      </div>
    </div>

    <div v-if="dueReminders.length" class="banner">
      <div v-for="item in dueReminders" :key="item.key" class="banner-item">
        <span>{{ item.message }}</span>
        <button @click="openReminder(item)">立即处理</button>
      </div>
    </div>

    <section v-if="activeTab === 'notes'" class="panel">
      <div class="section-head">
        <h3>我的笔记</h3>
        <button @click="startNoteEdit()">添加笔记</button>
      </div>
      <div v-if="showNoteForm" class="editor">
        <select v-model="noteForm.courseId">
          <option disabled value="">选择课程</option>
          <option v-for="course in courses" :key="course.id" :value="course.id">{{ course.title }}</option>
        </select>
        <input v-model.number="noteForm.pageNum" type="number" min="1" placeholder="页码">
        <textarea v-model="noteForm.content" rows="3" placeholder="输入笔记内容"></textarea>
        <div class="actions">
          <button @click="showNoteForm = false">取消</button>
          <button class="primary" @click="saveNote">保存</button>
        </div>
      </div>
      <div v-if="!notes.length" class="empty">还没有笔记。</div>
      <div v-for="note in notes" :key="note.id" class="card">
        <div class="card-head">
          <strong>{{ note.title || `第${note.pageNum}页笔记` }}</strong>
          <div class="actions">
            <button @click="startNoteEdit(note)">编辑</button>
            <button @click="deleteNote(note.id)">删除</button>
          </div>
        </div>
        <p>{{ note.note }}</p>
        <small>{{ getCourseName(note.courseId) }} · 第 {{ note.pageNum }} 页 · {{ formatDate(note.createdAt) }}</small>
      </div>
    </section>

    <section v-if="activeTab === 'favorites'" class="panel">
      <div class="section-head">
        <h3>我的收藏</h3>
        <div class="tags">
          <button :class="['tag-btn', { active: !selectedTag }]" @click="selectedTag = ''">全部</button>
          <button v-for="tag in availableTags" :key="tag" :class="['tag-btn', { active: selectedTag === tag }]" @click="selectedTag = tag">
            {{ tag }}
          </button>
        </div>
      </div>
      <div v-if="!filteredFavorites.length" class="empty">还没有收藏。</div>
      <div v-for="fav in filteredFavorites" :key="fav.id" class="card">
        <div class="card-head">
          <strong>{{ fav.title || '未命名收藏' }}</strong>
          <div class="actions">
            <button @click="startFavoriteEdit(fav)">编辑标签</button>
            <button @click="deleteFavorite(fav.id)">删除</button>
          </div>
        </div>
        <div class="tags">
          <span v-for="tag in fav.tags" :key="tag" class="chip">{{ tag }}</span>
          <span v-if="!fav.tags.length" class="chip muted">未设置标签</span>
        </div>
        <small>{{ getCourseName(fav.courseId) }} · 节点 {{ fav.nodeId || '-' }} · 第 {{ fav.pageNum || 0 }} 页</small>
        <div v-if="editingFavorite && editingFavorite.id === fav.id" class="editor slim">
          <input v-model="newTag" placeholder="输入标签后点添加">
          <button @click="addTag">添加</button>
          <div class="tags">
            <span v-for="tag in editingFavorite.tags" :key="tag" class="chip removable" @click="removeTag(tag)">{{ tag }} ×</span>
          </div>
          <div class="actions">
            <button @click="cancelFavoriteEdit">取消</button>
            <button class="primary" @click="saveTags">保存标签</button>
          </div>
        </div>
      </div>
    </section>

    <section v-if="activeTab === 'plans'" class="panel">
      <div class="section-head">
        <h3>复习计划</h3>
        <button @click="startPlanEdit()">创建计划</button>
      </div>
      <div v-if="showPlanForm" class="editor">
        <input v-model="planForm.name" placeholder="计划名称">
        <textarea v-model="planForm.description" rows="2" placeholder="计划描述"></textarea>
        <select v-model="planForm.frequency">
          <option value="daily">每日</option>
          <option value="weekly">每周</option>
          <option value="monthly">每月</option>
        </select>
        <div class="actions">
          <button @click="showPlanForm = false">取消</button>
          <button class="primary" @click="savePlan">保存</button>
        </div>
      </div>
      <div v-if="!reviewPlans.length" class="empty">还没有复习计划。</div>
      <div v-for="plan in reviewPlans" :key="plan.id" class="card">
        <div class="card-head">
          <strong>{{ plan.name }}</strong>
          <span :class="['badge', plan.status]">{{ formatPlanStatus(plan.status) }}</span>
        </div>
        <p>{{ plan.description || '暂无描述' }}</p>
        <small>{{ formatFrequency(plan.frequency) }} · 下次复习 {{ formatDate(plan.nextReviewDate) || '未设置' }}</small>
        <div class="actions top-gap">
          <button @click="viewPlanItems(plan)">查看内容</button>
          <button @click="startPlanEdit(plan)">编辑</button>
          <button @click="deletePlan(plan.id)">删除</button>
        </div>
      </div>
      <div v-if="currentPlan" class="subpanel">
        <div class="section-head">
          <h4>{{ currentPlan.name }} 的复习内容</h4>
          <button @click="currentPlan = null">收起</button>
        </div>
        <div v-if="!planItems.length" class="empty">还没有复习项。</div>
        <div v-for="item in planItems" :key="item.id" class="card compact">
          <div class="card-head">
            <strong>{{ item.itemType === 'note' ? getNoteContent(item.itemId) : getFavoriteTitle(item.itemId) }}</strong>
            <div class="actions">
              <button @click="markReviewed(item)">标记已复习</button>
              <button @click="removePlanItem(item.id)">移除</button>
            </div>
          </div>
          <small>{{ item.itemType === 'note' ? '笔记' : '收藏' }} · 优先级 {{ item.priority || 1 }} · 已复习 {{ item.reviewCount || 0 }} 次</small>
        </div>
        <div class="editor slim">
          <select v-model="addItemForm.type">
            <option value="note">笔记</option>
            <option value="favorite">收藏</option>
          </select>
          <select v-model="addItemForm.itemId">
            <option disabled value="">选择内容</option>
            <option v-for="item in availablePlanItems" :key="item.id" :value="item.id">{{ addItemForm.type === 'note' ? item.title : item.title }}</option>
          </select>
          <input v-model.number="addItemForm.priority" type="number" min="1" max="5" placeholder="优先级">
          <button class="primary" @click="addPlanItem">添加到计划</button>
        </div>
      </div>
    </section>

    <section v-if="activeTab === 'practice'" class="panel split">
      <div>
        <div class="section-head"><h3>练习记录</h3><button @click="loadPracticeData">刷新</button></div>
        <div v-if="!practiceHistory.length" class="empty">还没有练习记录。</div>
        <div v-for="item in practiceHistory" :key="item.taskId" class="card compact">
          <strong>{{ getCourseName(item.courseId) }} · 第 {{ item.pageNum || 0 }} 页</strong>
          <small>{{ item.difficulty || 'normal' }} · {{ formatDate(item.createdAt) }}</small>
          <small>得分 {{ item.attempt && item.attempt.score || 0 }} / 正确 {{ item.attempt && item.attempt.correctCount || 0 }} / {{ item.attempt && item.attempt.totalCount || item.questionCnt || 0 }}</small>
        </div>
      </div>
      <div>
        <div class="section-head"><h3>错题重做</h3><button @click="loadPracticeData">刷新</button></div>
        <div v-if="!wrongQuestions.length" class="empty">没有错题。</div>
        <div v-for="item in wrongQuestions" :key="item.recordId" class="card compact">
          <strong>{{ item.content }}</strong>
          <small>你的答案: {{ item.userAnswer || '-' }}</small>
          <small>参考答案: {{ item.correctAnswer || item.referenceAnswer || '-' }}</small>
          <div class="actions top-gap"><button class="primary" @click="retryWrongQuestion(item.questionId)">生成重做任务</button></div>
        </div>
      </div>
    </section>

    <section v-if="activeTab === 'notifications'" class="panel">
      <div class="section-head">
        <h3>通知中心</h3>
        <div class="actions">
          <span class="chip">{{ unreadNotificationCount }} 条未读</span>
          <button @click="markAllNotificationsRead">全部已读</button>
        </div>
      </div>
      <div v-if="!notifications.length" class="empty">还没有通知。</div>
      <div v-for="item in notifications" :key="item.id" class="card compact">
        <div class="card-head">
          <strong>{{ item.title }}</strong>
          <span :class="['badge', item.status]">{{ formatNotificationStatus(item.status) }}</span>
        </div>
        <p>{{ item.content }}</p>
        <small>{{ formatNotificationType(item.type) }} · {{ formatDate(item.createdAt) }}</small>
        <div class="actions top-gap">
          <button @click="markNotificationRead(item)" :disabled="item.status === 'read'">标记已读</button>
          <button @click="deleteNotification(item.id)">删除</button>
        </div>
      </div>
    </section>

    <section v-if="activeTab === 'tasks'" class="panel split">
      <div>
        <div class="section-head"><h3>调度任务</h3><button @click="loadTaskData">刷新</button></div>
        <div v-if="!scheduledTasks.length" class="empty">还没有调度任务。</div>
        <div v-for="task in scheduledTasks" :key="task.id" class="card compact">
          <div class="card-head">
            <strong>{{ formatTaskType(task.taskType) }}</strong>
            <span :class="['badge', task.status]">{{ task.status }}</span>
          </div>
          <small>{{ task.description || '无描述' }}</small>
          <small>计划时间: {{ formatDate(task.scheduledAt || task.nextAttempt) || '-' }}</small>
          <div class="actions top-gap">
            <button @click="executeTask(task.id)">立即执行</button>
            <button @click="deleteTask(task.id)">删除</button>
          </div>
        </div>
      </div>
      <div>
        <div class="section-head"><h3>任务状态</h3><button @click="loadTaskData">刷新</button></div>
        <div v-if="!taskStatuses.length" class="empty">还没有执行记录。</div>
        <div v-for="item in taskStatuses" :key="item.id" class="card compact">
          <div class="card-head">
            <strong>{{ formatTaskType(item.taskType) }}</strong>
            <span :class="['badge', item.status]">{{ item.status }}</span>
          </div>
          <small>进度 {{ item.progress || 0 }}% · {{ item.message || '等待执行结果' }}</small>
          <small>{{ formatDate(item.startTime) || '-' }} 至 {{ formatDate(item.endTime) || '-' }}</small>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import { ElMessage } from 'element-plus'
import { studentCoursewareApi } from '@/services/v1/coursewareApi'
import { studentNotificationApi } from '@/services/v1/notificationApi'
import { studentTaskApi } from '@/services/v1/taskApi'

const emptyNoteForm = () => ({ courseId: '', pageNum: 1, content: '' })
const emptyPlanForm = () => ({ name: '', description: '', frequency: 'daily' })
const emptyAddItemForm = () => ({ type: 'note', itemId: '', priority: 3 })

export default {
  name: 'StudentPersonalCenter',
  props: { studentId: { type: String, required: true } },
  data() {
    return {
      activeTab: 'notes',
      tabs: [{ key: 'notes', label: '笔记' }, { key: 'favorites', label: '收藏' }, { key: 'plans', label: '复习计划' }, { key: 'practice', label: '练习记录' }, { key: 'notifications', label: '通知' }, { key: 'tasks', label: '任务' }],
      courses: [], notes: [], favorites: [], reviewPlans: [], planItems: [], practiceHistory: [], wrongQuestions: [], notifications: [], scheduledTasks: [], taskStatuses: [],
      unreadNotificationCount: 0, selectedTag: '', currentPlan: null, showNoteForm: false, showPlanForm: false, editingNote: null, editingPlan: null, editingFavorite: null, newTag: '',
      noteForm: emptyNoteForm(), planForm: emptyPlanForm(), addItemForm: emptyAddItemForm()
    }
  },
  computed: {
    availableTags() { return [...new Set(this.favorites.flatMap(item => item.tags || []))] },
    filteredFavorites() { return this.selectedTag ? this.favorites.filter(item => (item.tags || []).includes(this.selectedTag)) : this.favorites },
    availablePlanItems() { return this.addItemForm.type === 'favorite' ? this.favorites : this.notes },
    dueReminders() {
      const now = Date.now()
      const planReminders = this.reviewPlans.filter(item => item.status === 'active' && item.nextReviewDate && new Date(item.nextReviewDate).getTime() <= now).map(item => ({ key: `plan-${item.id}`, kind: 'plan', message: `复习计划“${item.name}”已到期`, plan: item }))
      const notificationReminders = this.notifications.filter(item => item.status !== 'read' && String(item.type || '').includes('review')).map(item => ({ key: `notification-${item.id}`, kind: 'notification', message: item.title, notification: item }))
      return [...planReminders, ...notificationReminders].slice(0, 5)
    }
  },
  async created() { await this.loadData() },
  methods: {
    async loadData() { await this.loadCollections(); await this.loadPracticeData(); await this.loadNotificationData(); await this.loadTaskData() },
    async loadCollections() {
      try {
        const [courseRes, notesRes, favRes, plansRes] = await Promise.all([studentCoursewareApi.list(), studentCoursewareApi.listNotes({ studentId: this.studentId, pageSize: 100 }), studentCoursewareApi.listFavorites({ studentId: this.studentId, pageSize: 100 }), studentCoursewareApi.listReviewPlans(this.studentId)])
        this.courses = courseRes.data || []
        this.notes = (notesRes.data && notesRes.data.items) || []
        this.favorites = (((favRes.data && favRes.data.items) || [])).map(item => ({ ...item, tags: this.parseFavoriteTags(item.tags) }))
        this.reviewPlans = plansRes.data || []
      } catch (error) { ElMessage.error(`加载个人资料失败：${error.message}`) }
    },
    async loadPracticeData() {
      try {
        const [historyRes, wrongRes] = await Promise.all([studentCoursewareApi.getPracticeHistory({ studentId: this.studentId, pageSize: 20 }), studentCoursewareApi.getWrongQuestions({ studentId: this.studentId, pageSize: 20 })])
        this.practiceHistory = historyRes.data.items || []
        this.wrongQuestions = wrongRes.data.items || []
      } catch (error) { ElMessage.error(`加载练习数据失败：${error.message}`) }
    },
    async loadNotificationData() {
      try {
        const [listRes, countRes] = await Promise.all([studentNotificationApi.list({ studentId: this.studentId, pageSize: 20 }), studentNotificationApi.getUnreadCount(this.studentId)])
        this.notifications = listRes.data.list || []
        this.unreadNotificationCount = (countRes.data && countRes.data.unreadCount) || 0
      } catch (error) { ElMessage.error(`加载通知失败：${error.message}`) }
    },
    async loadTaskData() {
      try {
        const [taskRes, statusRes] = await Promise.all([studentTaskApi.listScheduled({ studentId: this.studentId, pageSize: 20 }), studentTaskApi.listStatuses({ studentId: this.studentId, pageSize: 20 })])
        this.scheduledTasks = taskRes.data.list || []
        this.taskStatuses = statusRes.data.list || []
      } catch (error) { ElMessage.error(`加载任务失败：${error.message}`) }
    },
    startNoteEdit(note = null) { this.editingNote = note; this.noteForm = note ? { courseId: note.courseId, pageNum: note.pageNum, content: note.note } : emptyNoteForm(); this.showNoteForm = true },
    async saveNote() {
      if (!this.noteForm.courseId || !this.noteForm.content) return ElMessage.warning('请先填写课程和笔记内容')
      try {
        if (this.editingNote) await studentCoursewareApi.deleteNote(this.editingNote.id, this.studentId)
        await studentCoursewareApi.saveNote({ studentId: this.studentId, courseId: this.noteForm.courseId, pageNum: this.noteForm.pageNum, content: this.noteForm.content })
        this.showNoteForm = false; this.editingNote = null; this.noteForm = emptyNoteForm(); await this.loadCollections(); ElMessage.success('笔记已保存')
      } catch (error) { ElMessage.error(`保存笔记失败：${error.message}`) }
    },
    async deleteNote(id) { if (!window.confirm('确认删除这条笔记吗？')) return; try { await studentCoursewareApi.deleteNote(id, this.studentId); await this.loadCollections(); ElMessage.success('笔记已删除') } catch (error) { ElMessage.error(`删除笔记失败：${error.message}`) } },
    startFavoriteEdit(favorite) { this.editingFavorite = { ...favorite, tags: [...(favorite.tags || [])] }; this.newTag = '' },
    cancelFavoriteEdit() { this.editingFavorite = null; this.newTag = '' },
    addTag() { const tag = this.newTag.trim(); if (tag && !this.editingFavorite.tags.includes(tag)) this.editingFavorite.tags.push(tag); this.newTag = '' },
    removeTag(tag) { this.editingFavorite.tags = this.editingFavorite.tags.filter(item => item !== tag) },
    async saveTags() {
      try {
        await studentCoursewareApi.addFavorite({ studentId: this.studentId, courseId: this.editingFavorite.courseId, nodeId: this.editingFavorite.nodeId, pageNum: this.editingFavorite.pageNum, title: this.editingFavorite.title, tags: this.editingFavorite.tags })
        this.cancelFavoriteEdit(); await this.loadCollections(); ElMessage.success('标签已更新')
      } catch (error) { ElMessage.error(`保存标签失败：${error.message}`) }
    },
    async deleteFavorite(id) { if (!window.confirm('确认删除这条收藏吗？')) return; try { await studentCoursewareApi.deleteFavorite(id); await this.loadCollections(); ElMessage.success('收藏已删除') } catch (error) { ElMessage.error(`删除收藏失败：${error.message}`) } },
    startPlanEdit(plan = null) { this.editingPlan = plan; this.planForm = plan ? { name: plan.name, description: plan.description, frequency: plan.frequency } : emptyPlanForm(); this.showPlanForm = true },
    async savePlan() {
      if (!this.planForm.name) return ElMessage.warning('请先填写计划名称')
      try {
        if (this.editingPlan) await studentCoursewareApi.updateReviewPlan(this.editingPlan.id, this.planForm)
        else await studentCoursewareApi.createReviewPlan({ studentId: this.studentId, ...this.planForm })
        this.showPlanForm = false; this.editingPlan = null; this.planForm = emptyPlanForm(); await this.loadCollections(); ElMessage.success('复习计划已保存')
      } catch (error) { ElMessage.error(`保存计划失败：${error.message}`) }
    },
    async deletePlan(id) { if (!window.confirm('确认删除这条复习计划吗？')) return; try { await studentCoursewareApi.deleteReviewPlan(id); if (this.currentPlan && this.currentPlan.id === id) this.currentPlan = null; await this.loadCollections(); ElMessage.success('计划已删除') } catch (error) { ElMessage.error(`删除计划失败：${error.message}`) } },
    async viewPlanItems(plan) { try { const res = await studentCoursewareApi.listReviewPlanItems(plan.id); this.currentPlan = plan; this.planItems = res.data || [] } catch (error) { ElMessage.error(`加载复习项失败：${error.message}`) } },
    async addPlanItem() {
      if (!this.currentPlan || !this.addItemForm.itemId) return ElMessage.warning('请先选择要加入的内容')
      try {
        await studentCoursewareApi.addReviewPlanItem({ reviewPlanId: this.currentPlan.id, itemType: this.addItemForm.type, itemId: this.addItemForm.itemId, priority: this.addItemForm.priority })
        this.addItemForm = emptyAddItemForm(); await this.viewPlanItems(this.currentPlan); ElMessage.success('复习项已添加')
      } catch (error) { ElMessage.error(`添加复习项失败：${error.message}`) }
    },
    async markReviewed(item) { try { await studentCoursewareApi.updateReviewPlanItem(item.id, { lastReviewedAt: new Date().toISOString(), reviewCount: (item.reviewCount || 0) + 1 }); await this.viewPlanItems(this.currentPlan); await this.loadCollections(); ElMessage.success('已标记复习') } catch (error) { ElMessage.error(`标记复习失败：${error.message}`) } },
    async removePlanItem(id) { try { await studentCoursewareApi.deleteReviewPlanItem(id); await this.viewPlanItems(this.currentPlan); ElMessage.success('复习项已移除') } catch (error) { ElMessage.error(`移除复习项失败：${error.message}`) } },
    async retryWrongQuestion(questionId) { try { const res = await studentCoursewareApi.retryWrongQuestion({ questionId, studentId: this.studentId }); ElMessage.success(`已生成重做任务 ${res.data.taskId}`); await this.loadPracticeData() } catch (error) { ElMessage.error(`生成重做任务失败：${error.message}`) } },
    async markNotificationRead(item) { if (item.status === 'read') return; try { await studentNotificationApi.markAsRead(item.id); await this.loadNotificationData(); ElMessage.success('通知已标记为已读') } catch (error) { ElMessage.error(`处理通知失败：${error.message}`) } },
    async markAllNotificationsRead() { try { await studentNotificationApi.markAllAsRead(this.studentId); await this.loadNotificationData(); ElMessage.success('全部通知已读') } catch (error) { ElMessage.error(`批量已读失败：${error.message}`) } },
    async deleteNotification(id) { if (!window.confirm('确认删除这条通知吗？')) return; try { await studentNotificationApi.remove(id); await this.loadNotificationData(); await this.loadTaskData(); ElMessage.success('通知已删除') } catch (error) { ElMessage.error(`删除通知失败：${error.message}`) } },
    async executeTask(id) { try { await studentTaskApi.executeNow(id); await this.loadTaskData(); await this.loadNotificationData(); ElMessage.success('任务已提交执行') } catch (error) { ElMessage.error(`执行任务失败：${error.message}`) } },
    async deleteTask(id) { if (!window.confirm('确认删除这条任务吗？')) return; try { await studentTaskApi.remove(id); await this.loadTaskData(); ElMessage.success('任务已删除') } catch (error) { ElMessage.error(`删除任务失败：${error.message}`) } },
    async openReminder(item) { if (item.kind === 'plan') return this.viewPlanItems(item.plan); await this.markNotificationRead(item.notification) },
    parseFavoriteTags(raw) { if (Array.isArray(raw)) return raw; if (!raw) return []; try { const parsed = JSON.parse(raw); return Array.isArray(parsed) ? parsed : [] } catch (error) { return [] } },
    getCourseName(id) { const target = this.courses.find(item => item.id === id); return (target && target.title) || id || '未关联课程' },
    getNoteContent(id) { const target = this.notes.find(item => item.id === id); return (target && target.note) || '未找到笔记' },
    getFavoriteTitle(id) { const target = this.favorites.find(item => item.id === id); return (target && target.title) || '未找到收藏' },
    formatDate(value) { return value ? new Date(value).toLocaleString('zh-CN') : '' },
    formatFrequency(value) { return { daily: '每日', weekly: '每周', monthly: '每月' }[value] || value || '-' },
    formatPlanStatus(value) { return { active: '进行中', paused: '暂停', completed: '已完成' }[value] || value || '-' },
    formatNotificationType(value) { return { system: '系统通知', review_reminder: '复习提醒', practice: '练习通知' }[value] || value || '普通通知' },
    formatNotificationStatus(value) { return { unread: '未读', read: '已读', scheduled: '待发送', sent: '已发送' }[value] || value || '-' },
    formatTaskType(value) { return { review_reminder: '复习提醒', practice_generation: '练习生成', notification: '通知发送' }[value] || value || '通用任务' }
  }
}
</script>

<style scoped>
.personal-center { display: grid; gap: 16px; padding: 20px; color: #1f2937; }
.header, .panel, .subpanel { background: rgba(255, 255, 255, 0.96); border: 1px solid #e5edf8; border-radius: 16px; padding: 18px; box-shadow: 0 8px 20px rgba(15, 23, 42, 0.06); }
.header h2, .section-head h3, .section-head h4 { margin: 0; }
.header p, p, small { margin: 0; color: #64748b; }
.tabs, .actions, .tags { display: flex; flex-wrap: wrap; gap: 8px; }
.tab, button, select, input, textarea { border-radius: 10px; border: 1px solid #cfd9e8; background: #fff; font: inherit; }
.tab, button { padding: 8px 14px; cursor: pointer; }
.tab.active, .primary { background: #1d4ed8; border-color: #1d4ed8; color: #fff; }
.panel, .subpanel, .editor, .card, .banner { display: grid; gap: 12px; }
.section-head, .card-head, .banner-item { display: flex; justify-content: space-between; gap: 12px; align-items: center; }
.editor { padding: 14px; background: #f8fbff; border: 1px solid #dbe8ff; border-radius: 14px; }
.editor.slim { grid-template-columns: repeat(auto-fit, minmax(140px, 1fr)); align-items: center; }
input, select, textarea { width: 100%; padding: 10px 12px; }
textarea { resize: vertical; }
.card { padding: 14px; background: #f8fafc; border-radius: 14px; border: 1px solid #e5edf8; }
.card.compact small + small { margin-top: 4px; }
.banner { padding: 14px 18px; background: #fff7d6; border: 1px solid #f4db86; border-radius: 16px; }
.split { display: grid; grid-template-columns: repeat(auto-fit, minmax(260px, 1fr)); gap: 16px; }
.chip, .badge { display: inline-flex; align-items: center; border-radius: 999px; padding: 4px 10px; font-size: 12px; background: #e2e8f0; color: #334155; }
.chip.muted, .badge.read, .badge.completed { background: #dbeafe; color: #1d4ed8; }
.badge.active, .badge.unread, .badge.pending { background: #dcfce7; color: #166534; }
.badge.paused, .badge.scheduled { background: #fef3c7; color: #92400e; }
.badge.failed { background: #fee2e2; color: #b91c1c; }
.top-gap { margin-top: 6px; }
.empty { padding: 20px; text-align: center; color: #94a3b8; background: #f8fafc; border-radius: 12px; }
.removable { cursor: pointer; }
@media (max-width: 768px) { .personal-center { padding: 12px; } .header, .panel, .subpanel { padding: 14px; } .section-head, .card-head, .banner-item { flex-direction: column; align-items: flex-start; } }
</style>
