<template>
  <div class="top-nav">
    <div class="nav-left">
      <div class="brand-mark">T</div>
      <div class="brand-copy">
        <span class="app-title">智能互动教学平台</span>
        <span class="app-subtitle">教师编排台</span>
      </div>
    </div>
    <div class="nav-right">
      <div class="backend-status" :class="backendStatusClass">
        <span class="status-dot"></span>
        服务 {{ backendStatusText }}
      </div>
      <div class="teacher-info">
        <div class="avatar">
          <span>教</span>
        </div>
        <div class="teacher-copy">
          <span class="teacher-name">教师 2025T001</span>
          <span class="teacher-role">课程编排与发布</span>
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
  }
})
</script>

<style scoped>
.top-nav {
  min-height: 72px;
  background:
    radial-gradient(circle at left top, rgba(255, 255, 255, 0.18), transparent 30%),
    linear-gradient(105deg, #0f766e 0%, #0369a1 42%, #1d4ed8 100%);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  color: white;
  box-shadow: 0 18px 40px rgba(14, 116, 144, 0.2);
  flex-shrink: 0;
}

.nav-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-mark {
  width: 42px;
  height: 42px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 800;
  background: rgba(255, 255, 255, 0.16);
  border: 1px solid rgba(255, 255, 255, 0.22);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.18);
}

.brand-copy {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.app-title {
  font-size: 17px;
  font-weight: 700;
  letter-spacing: 0.03em;
}

.app-subtitle {
  font-size: 12px;
  opacity: 0.78;
  letter-spacing: 0.12em;
  text-transform: uppercase;
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
  font-weight: 600;
  padding: 8px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.14);
  border: 1px solid rgba(255, 255, 255, 0.12);
  backdrop-filter: blur(10px);
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
  gap: 10px;
  padding: 8px 12px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.12);
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
}

.teacher-copy {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.teacher-name {
  font-size: 13px;
  font-weight: 700;
}

.teacher-role {
  font-size: 11px;
  opacity: 0.78;
}
</style>