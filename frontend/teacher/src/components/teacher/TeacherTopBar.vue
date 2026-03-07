<template>
  <div class="top-nav">
    <div class="nav-left">
      <span class="app-title">智能互动教学平台 · 教师端</span>
    </div>
    <div class="nav-right">
      <div class="backend-status" :class="backendStatusClass">
        后端：{{ backendStatusText }}
      </div>
      <div class="teacher-info">
        <div class="avatar">
          <span>教</span>
        </div>
        <span class="teacher-name">教师 2025T001</span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  backendStatusClass: {
    type: String,
    default: 'checking'
  },
  backendStatusText: {
    type: String,
    default: '检测中'
  }
})
</script>

<style scoped>
.top-nav {
  height: 60px;
  background: linear-gradient(90deg, #2563eb 0%, #3b82f6 100%);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  color: white;
  box-shadow: 0 10px 24px rgba(37, 99, 235, 0.25);
}
.nav-right {
  display: flex;
  align-items: center;
  gap: 12px;
}
.app-title {
  font-size: 16px;
  font-weight: 500;
}
.backend-status {
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.2);
}
.backend-status.online {
  background: rgba(82, 196, 26, 0.2);
  color: #f6ffed;
}
.backend-status.offline {
  background: rgba(255, 77, 79, 0.2);
  color: #fff1f0;
}
.backend-status.checking {
  background: rgba(250, 173, 20, 0.2);
  color: #fffbe6;
}
.teacher-info {
  display: flex;
  align-items: center;
  gap: 10px;
}
.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
}
.avatar span {
  width: 100%;
  height: 100%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  background: rgba(255, 255, 255, 0.2);
}
.teacher-name {
  font-size: 14px;
}
</style>