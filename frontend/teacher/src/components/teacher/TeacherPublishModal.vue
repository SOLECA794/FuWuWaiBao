<template>
  <div class="modal-overlay" v-if="visible" @click="$emit('close')">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>发布课件到课程班级</h3>
        <button @click="$emit('close')" class="close-btn">×</button>
      </div>
      <div class="publish-form">
        <div class="form-item">
          <label>当前课件：</label>
          <span>{{ currentCourseName || '未选择课件' }}</span>
        </div>
        <div class="form-item stacked">
          <label>所属课程：</label>
          <select :value="teachingCourseId" class="scope-select" @change="$emit('update:teachingCourseId', $event.target.value)">
            <option value="">请选择课程</option>
            <option v-for="item in teachingCourseOptions" :key="item.id" :value="item.id">{{ item.name }}</option>
          </select>
        </div>
        <div class="form-item stacked">
          <label>教学班级：</label>
          <select :value="courseClassId" class="scope-select" @change="$emit('update:courseClassId', $event.target.value)">
            <option value="">请选择班级</option>
            <option v-for="item in courseClassOptions" :key="item.id" :value="item.id">{{ item.name }}</option>
          </select>
        </div>
        <div class="form-item">
          <label>发布范围：</label>
          <select :value="publishScope" class="scope-select" @change="$emit('update:publishScope', $event.target.value)">
            <option value="all">全部学生</option>
            <option value="class">仅所选班级</option>
          </select>
        </div>
        <div class="form-actions">
          <button @click="$emit('submit')" class="confirm-btn" :disabled="!teachingCourseId || !courseClassId">确认发布</button>
          <button @click="$emit('close')" class="cancel-btn">取消</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  currentCourseName: {
    type: String,
    default: ''
  },
  publishScope: {
    type: String,
    default: 'all'
  },
  teachingCourseId: {
    type: String,
    default: ''
  },
  courseClassId: {
    type: String,
    default: ''
  },
  teachingCourseOptions: {
    type: Array,
    default: () => []
  },
  courseClassOptions: {
    type: Array,
    default: () => []
  }
})

defineEmits(['close', 'submit', 'update:publishScope', 'update:teachingCourseId', 'update:courseClassId'])
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal-content {
  width: 420px;
  background: #fff;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 24px 50px rgba(15, 23, 42, 0.18);
}
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 18px;
  border-bottom: 1px solid #e6ecf5;
}
.close-btn {
  border: none;
  background: transparent;
  font-size: 24px;
  cursor: pointer;
}
.publish-form {
  padding: 18px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.form-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.form-item.stacked {
  flex-direction: column;
  align-items: stretch;
}

.form-item.stacked label {
  font-size: 13px;
  color: #475569;
}

.scope-select {
  border: 1px solid #dbe3ef;
  border-radius: 8px;
  padding: 6px 10px;
  min-height: 34px;
}
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
.confirm-btn,
.cancel-btn {
  border: none;
  border-radius: 10px;
  padding: 10px 14px;
  cursor: pointer;
}
.confirm-btn {
  background: #2F605A;
  color: #fff;
}

.confirm-btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}
.cancel-btn {
  background: #e2e8f0;
  color: #334155;
}
</style>