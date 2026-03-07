<template>
  <div class="modal-overlay" v-if="visible" @click="$emit('close')">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>发布课件给学生端</h3>
        <button @click="$emit('close')" class="close-btn">×</button>
      </div>
      <div class="publish-form">
        <div class="form-item">
          <label>当前课件：</label>
          <span>{{ currentCourseName }}</span>
        </div>
        <div class="form-item">
          <label>发布范围：</label>
          <select :value="publishScope" class="scope-select" @change="$emit('update:publishScope', $event.target.value)">
            <option value="all">全部学生</option>
            <option value="class1">班级1</option>
            <option value="class2">班级2</option>
          </select>
        </div>
        <div class="form-actions">
          <button @click="$emit('submit')" class="confirm-btn">确认发布</button>
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
  }
})

defineEmits(['close', 'submit', 'update:publishScope'])
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
.scope-select {
  border: 1px solid #dbe3ef;
  border-radius: 8px;
  padding: 6px 10px;
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
  background: #2563eb;
  color: #fff;
}
.cancel-btn {
  background: #e2e8f0;
  color: #334155;
}
</style>