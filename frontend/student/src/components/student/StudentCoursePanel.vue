<template>
  <div class="course-card">
    <div class="course-header">
      <div>
        <div class="header-label">当前讲授</div>
        <h3>{{ currentCourseName || '暂无课程' }} - 第{{ currentPage }}页</h3>
      </div>
      <div class="course-tag">
        <el-tag size="small" effect="dark">课件学习</el-tag>
        <span class="page-count">{{ currentPage }}/{{ totalPage }}</span>
      </div>
    </div>

    <div class="progress-strip">
      <div class="progress-meta">
        <span>学习进度 {{ progressPercent }}%</span>
        <span>{{ isPlay ? '讲授进行中' : '等待继续' }}</span>
      </div>
      <div class="progress-track">
        <div class="progress-fill" :style="{ width: progressPercent + '%' }"></div>
      </div>
    </div>

    <div class="playback-strip" v-if="pageTimelineDuration > 0">
      <div class="playback-head">
        <div>
          <div class="playback-title">讲授时间轴</div>
          <div class="playback-subtitle">当前节点 {{ currentNodeTitle || '未定位节点' }}</div>
        </div>
        <div class="playback-meta-group">
          <span class="playback-mode-badge" :class="playbackModeClass">{{ playbackModeLabel }}</span>
          <div class="playback-time">{{ formatTime(currentTimelineSec) }} / {{ formatTime(pageTimelineDuration) }}</div>
        </div>
      </div>
      <div class="playback-track">
        <div class="playback-fill" :style="{ width: timelinePercent + '%' }"></div>
      </div>
      <div class="playback-node-meta" v-if="activeNodeDuration > 0">
        <span>节点进度 {{ formatTime(activeNodeElapsedSec) }} / {{ formatTime(activeNodeDuration) }}</span>
        <span>{{ activeNodeTypeLabel }}</span>
        <span v-if="audioStatusText">{{ audioStatusText }}</span>
      </div>
    </div>

    <div class="course-content">
      <img v-if="courseImg" :src="courseImg" alt="课件内容" class="course-img" />
      <div v-else class="no-courseware">当前没有可预览课件，请联系教师先发布课件</div>
      <div
        v-if="tracePoint"
        class="trace-highlight"
        :style="{ top: traceTop + 'px', left: traceLeft + 'px' }"
      ></div>
    </div>

    <div class="course-control">
      <el-button @click="$emit('prev-page')" icon="el-icon-arrow-left" size="small">上一页</el-button>
      <el-button @click="$emit('toggle-play')" :icon="isPlay ? 'el-icon-pause' : 'el-icon-play'" size="small">
        {{ isPlay ? '暂停' : '播放' }}
      </el-button>
      <el-button @click="$emit('next-page')" icon="el-icon-arrow-right" size="small">下一页</el-button>
    </div>
    <div class="page-summary" v-if="pageSummary">
      <h4>本页摘要</h4>
      <p>{{ pageSummary }}</p>
    </div>
    <div class="node-panel" v-if="playbackNodes.length">
      <div class="node-panel-header">
        <h4>讲授节点</h4>
        <span>{{ playbackNodes.length }} 段</span>
      </div>
      <div class="node-list">
        <button
          v-for="node in playbackNodes"
          :key="node.node_id"
          class="node-chip"
          :class="{ active: node.node_id === currentNodeId }"
          @click="$emit('select-node', node.node_id)"
        >
          <div class="node-chip-head">
            <strong>{{ node.title || node.node_id }}</strong>
            <span class="node-type" :class="node.type">{{ nodeTypeLabel(node.type) }}</span>
          </div>
          <span>{{ node.text }}</span>
        </button>
      </div>
    </div>
    <div class="control-tip">支持自动记录断点，切页后将实时同步学习进度</div>
  </div>
</template>

<script setup>
/* eslint-disable no-undef */
import { computed } from 'vue'

const props = defineProps({
  currentCourseName: {
    type: String,
    default: ''
  },
  currentPage: {
    type: Number,
    default: 1
  },
  totalPage: {
    type: Number,
    default: 1
  },
  pageTimelineDuration: {
    type: Number,
    default: 0
  },
  currentTimelineSec: {
    type: Number,
    default: 0
  },
  activeNodeElapsedSec: {
    type: Number,
    default: 0
  },
  activeNodeDuration: {
    type: Number,
    default: 0
  },
  currentNodeTitle: {
    type: String,
    default: ''
  },
  activeNodeTypeLabel: {
    type: String,
    default: ''
  },
  playbackMode: {
    type: String,
    default: 'duration_timeline'
  },
  playbackAudioMeta: {
    type: Object,
    default: null
  },
  progressPercent: {
    type: Number,
    default: 0
  },
  courseImg: {
    type: String,
    default: ''
  },
  playbackNodes: {
    type: Array,
    default: () => []
  },
  currentNodeId: {
    type: String,
    default: ''
  },
  pageSummary: {
    type: String,
    default: ''
  },
  tracePoint: {
    type: Boolean,
    default: false
  },
  traceTop: {
    type: Number,
    default: 0
  },
  traceLeft: {
    type: Number,
    default: 0
  },
  isPlay: {
    type: Boolean,
    default: false
  }
})

defineEmits(['prev-page', 'toggle-play', 'next-page', 'select-node'])

const timelinePercent = computed(() => {
  if (!props.pageTimelineDuration) return 0
  return Math.min(100, Math.max(0, Math.round((props.currentTimelineSec / props.pageTimelineDuration) * 100)))
})

const playbackModeLabel = computed(() => {
  if (props.playbackMode === 'audio_timeline') return '音频轨'
  return '时间轴'
})

const playbackModeClass = computed(() => {
  return props.playbackMode === 'audio_timeline' ? 'audio' : 'duration'
})

const audioStatusText = computed(() => {
  const status = props.playbackAudioMeta?.audio_status
  const duration = Number(props.playbackAudioMeta?.audio_duration_sec || 0)
  if (!status) return ''
  if (status === 'ready' && duration > 0) {
    return `音频已生成 ${formatTime(duration)}`
  }
  if (status === 'processing') return '音频生成中'
  return '使用时长驱动讲解'
})

const nodeTypeLabel = (type) => {
  if (type === 'opening') return '开场'
  if (type === 'transition') return '过渡'
  return '讲解'
}

const formatTime = (seconds) => {
  const normalized = Math.max(0, Math.floor(seconds || 0))
  const mins = String(Math.floor(normalized / 60)).padStart(2, '0')
  const secs = String(normalized % 60).padStart(2, '0')
  return `${mins}:${secs}`
}
</script>

<style scoped>
.course-card {
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98) 0%, rgba(248, 250, 252, 0.96) 100%);
  border-radius: 24px;
  padding: 20px;
  height: 100%;
  border: 1px solid rgba(148, 163, 184, 0.16);
  box-shadow: 0 20px 40px rgba(15, 23, 42, 0.08);
  display: flex;
  flex-direction: column;
}
.course-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.header-label {
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #0f766e;
  margin-bottom: 4px;
}
.course-header h3 {
  font-size: 20px;
  color: #0f172a;
}
.course-tag {
  display: flex;
  gap: 8px;
  align-items: center;
}
.page-count {
  font-size: 13px;
  color: #475569;
}
.progress-strip {
  margin-bottom: 14px;
  padding: 12px 14px;
  border-radius: 16px;
  background: linear-gradient(90deg, rgba(15, 118, 110, 0.08) 0%, rgba(2, 132, 199, 0.08) 100%);
  border: 1px solid rgba(14, 165, 233, 0.14);
}
.progress-meta {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
  font-size: 12px;
  color: #0f172a;
  font-weight: 600;
}
.progress-track {
  height: 8px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.7);
  overflow: hidden;
}
.progress-fill {
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #0f766e 0%, #0284c7 100%);
}
.playback-strip {
  margin-bottom: 14px;
  padding: 12px 14px;
  border-radius: 16px;
  background: linear-gradient(180deg, rgba(15, 23, 42, 0.03) 0%, rgba(255, 255, 255, 0.78) 100%);
  border: 1px solid rgba(148, 163, 184, 0.18);
}
.playback-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}
.playback-title {
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #0369a1;
}
.playback-subtitle {
  margin-top: 4px;
  font-size: 14px;
  color: #0f172a;
  font-weight: 600;
}
.playback-time {
  font-size: 13px;
  font-weight: 700;
  color: #0f172a;
}
.playback-meta-group {
  display: flex;
  align-items: center;
  gap: 10px;
}
.playback-mode-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 58px;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.04em;
}
.playback-mode-badge.duration {
  color: #075985;
  background: rgba(14, 165, 233, 0.12);
}
.playback-mode-badge.audio {
  color: #065f46;
  background: rgba(16, 185, 129, 0.14);
}
.playback-track {
  height: 10px;
  border-radius: 999px;
  background: rgba(226, 232, 240, 0.9);
  overflow: hidden;
}
.playback-fill {
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #0ea5e9 0%, #2563eb 100%);
}
.playback-node-meta {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-top: 8px;
  font-size: 12px;
  color: #475569;
}
.course-content {
  flex: 1;
  background: linear-gradient(180deg, #ffffff 0%, #f8fbff 100%);
  border-radius: 16px;
  border: 1px solid rgba(226, 232, 240, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}
.course-img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}
.no-courseware {
  font-size: 14px;
  color: #64748b;
}
.trace-highlight {
  position: absolute;
  width: 180px;
  height: 100px;
  border: 3px solid #ff6633;
  background: rgba(255, 102, 51, 0.1);
  pointer-events: none;
  border-radius: 6px;
  animation: flash 1.2s infinite;
}
@keyframes flash {
  0% { opacity: 0.4; }
  50% { opacity: 0.8; }
  100% { opacity: 0.4; }
}
.course-control {
  margin-top: 16px;
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 12px;
}
.control-tip {
  margin-top: 10px;
  text-align: center;
  font-size: 12px;
  color: #64748b;
}
.page-summary,
.node-panel {
  margin-top: 14px;
  background: rgba(248, 250, 252, 0.86);
  border: 1px solid rgba(226, 232, 240, 0.9);
  border-radius: 16px;
  padding: 12px;
}
.page-summary h4,
.node-panel-header h4 {
  margin: 0 0 8px;
  font-size: 14px;
  color: #0f172a;
}
.page-summary p {
  margin: 0;
  font-size: 13px;
  color: #475569;
  line-height: 1.7;
}
.node-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}
.node-panel-header span {
  font-size: 12px;
  color: #64748b;
}
.node-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 220px;
  overflow: auto;
}
.node-chip {
  border: 1px solid #dbe3ef;
  background: linear-gradient(180deg, #ffffff 0%, #f8fafc 100%);
  border-radius: 14px;
  padding: 10px 12px;
  text-align: left;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}
.node-chip:hover {
  transform: translateY(-1px);
  box-shadow: 0 8px 16px rgba(15, 23, 42, 0.06);
}
.node-chip-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 4px;
}
.node-chip strong {
  font-size: 12px;
  color: #2563eb;
}
.node-type {
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 700;
  color: #475569;
  background: #e2e8f0;
}
.node-type.opening {
  background: rgba(16, 185, 129, 0.14);
  color: #047857;
}
.node-type.explain {
  background: rgba(14, 165, 233, 0.14);
  color: #0369a1;
}
.node-type.transition {
  background: rgba(245, 158, 11, 0.14);
  color: #b45309;
}
.node-chip span {
  display: block;
  font-size: 13px;
  color: #334155;
  line-height: 1.6;
}
.node-chip.active {
  border-color: #2563eb;
  box-shadow: 0 0 0 1px rgba(37, 99, 235, 0.15), 0 10px 20px rgba(37, 99, 235, 0.08);
  background: linear-gradient(180deg, #eff6ff 0%, #f8fbff 100%);
}

@media (max-width: 720px) {
  .course-header,
  .progress-meta,
  .playback-head,
  .playback-node-meta,
  .node-chip-head {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>