<template>
  <div>
    <header class="header">
      <div class="header-left">
        <div class="logo">学</div>
        <div class="system-name">智能学习课堂系统</div>
      </div>
      <div class="header-right">
        <span>用户名：{{ studentId }}</span>
        <span>当前课程：{{ currentCourseName || '未分配课程' }}</span>
        <span class="backend-status" :class="backendStatus">后端：{{ backendStatusLabel }}</span>
      </div>
    </header>

    <section class="learning-overview">
      <div class="overview-card">
        <div class="overview-title">学习进度</div>
        <el-progress :percentage="progressPercent" :stroke-width="10" :show-text="true" />
      </div>
      <div class="overview-card compact">
        <div class="overview-title">学习状态</div>
        <div class="overview-meta">
          <span>当前页：第{{ currentPage }}页</span>
          <span>播放：{{ isPlay ? '进行中' : '已暂停' }}</span>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
/* eslint-disable no-undef */
defineProps({
  studentId: {
    type: String,
    required: true
  },
  currentCourseName: {
    type: String,
    default: ''
  },
  backendStatus: {
    type: String,
    default: 'checking'
  },
  backendStatusLabel: {
    type: String,
    default: '检测中'
  },
  progressPercent: {
    type: Number,
    default: 0
  },
  currentPage: {
    type: Number,
    default: 1
  },
  isPlay: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
.header {
  height: 60px;
  background: linear-gradient(90deg, #2563eb 0%, #3b82f6 100%);
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  box-shadow: 0 10px 24px rgba(37, 99, 235, 0.24);
}
.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}
.logo {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.18);
  border: 1px solid rgba(255, 255, 255, 0.38);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
}
.system-name {
  font-size: 16px;
  font-weight: 500;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
}
.backend-status {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
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
.learning-overview {
  display: grid;
  grid-template-columns: 1.4fr 1fr;
  gap: 12px;
  padding: 12px 20px 0;
}
.overview-card {
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid #e6ecf5;
  padding: 12px 16px;
  box-shadow: 0 8px 18px rgba(15, 23, 42, 0.06);
}
.overview-card.compact {
  display: flex;
  flex-direction: column;
  justify-content: center;
}
.overview-title {
  font-size: 13px;
  color: #475569;
  margin-bottom: 8px;
}
.overview-meta {
  display: flex;
  gap: 18px;
  font-size: 13px;
  color: #1e293b;
}
</style>