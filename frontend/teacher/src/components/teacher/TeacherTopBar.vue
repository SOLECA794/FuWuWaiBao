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
  height: 56px;
  background: linear-gradient(90deg, #1d4ed8 0%, #2563eb 60%, #3b82f6 100%);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  color: white;
  box-shadow: 0 2px 12px rgba(37, 99, 235, 0.3);
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
  background: rgba(255, 255, 255, 0.15);
}

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}

.backend-status.online {
  color: #bbf7d0;
}

.backend-status.online .status-dot {
  background: #4ade80;
  animation: pulse-green 2s ease-in-out infinite;
}

@keyframes pulse-green {
  0%, 100% { box-shadow: 0 0 0 0 rgba(74, 222, 128, 0.5); }
  50% { box-shadow: 0 0 0 4px rgba(74, 222, 128, 0); }
}

.backend-status.offline {
  color: #fecaca;
}

.backend-status.offline .status-dot {
  background: #f87171;
}

.backend-status.checking {
  color: #fde68a;
}

.backend-status.checking .status-dot {
  background: #fbbf24;
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
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
}

.teacher-name {
  font-size: 13px;
  opacity: 0.9;
}
</style>