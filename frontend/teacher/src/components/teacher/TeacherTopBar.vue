<template>
  <div class="top-nav">
    <div class="nav-left">
      <span class="app-icon">🎓</span>
      <span class="app-title">智能互动教学平台 · 教师端</span>
    </div>
    <div class="nav-right">
      <div class="backend-status" :class="backendStatusClass">
        <span class="status-dot"></span>
        后端：{{ backendStatusText }}
      </div>
      <div class="teacher-info">
        <div class="avatar">
          <span>教</span>
        </div>
        <div class="account-actions">
          <span class="teacher-name">{{ username }}</span>
          <button class="logout-btn" @click="$emit('logout')">退出登录</button>
        </div>
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
  },
  username: {
    type: String,
    default: '教师'
  }
})

defineEmits(['logout'])
</script>

<style scoped>
.top-nav {
  height: 56px;
  background: #ffffff;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  color: #1e293b;
  border-bottom: 1px solid #e2e8f0;
  flex-shrink: 0;
}

.nav-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.app-icon {
  font-size: 18px;
}

.app-title {
  font-size: 15px;
  font-weight: 600;
  letter-spacing: 0.01em;
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.backend-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 14px;
  background: #F4F7F7;
}

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}

.backend-status.online {
  color: #166534;
}

.backend-status.online .status-dot {
  background: #22c55e;
  animation: pulse-green 2s ease-in-out infinite;
}

@keyframes pulse-green {
  0%, 100% { box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.4); }
  50% { box-shadow: 0 0 0 4px rgba(34, 197, 94, 0); }
}

.backend-status.offline {
  color: #991b1b;
}

.backend-status.offline .status-dot {
  background: #ef4444;
}

.backend-status.checking {
  color: #854d0e;
}

.backend-status.checking .status-dot {
  background: #eab308;
  animation: blink 1.4s ease-in-out infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.25; }
}

.teacher-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #F4F7F7;
  color: #2F605A;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
}

.account-actions {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: flex-start;
}

.teacher-name {
  font-size: 13px;
  font-weight: 600;
  color: #1e293b;
}

.logout-btn {
  background: transparent;
  border: none;
  padding: 0;
  margin-top: 2px;
  font-size: 11px;
  color: #64748b;
  cursor: pointer;
  transition: color 0.2s;
}

.logout-btn:hover {
  color: #ef4444;
  text-decoration: underline;
}
</style>