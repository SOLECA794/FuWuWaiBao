<template>
  <div class="course-card">
    <div class="course-header">
      <div>
        <div class="header-label">课堂内容</div>
        <h3>{{ currentCourseName || '暂无课程' }} · 第{{ currentPage }}页</h3>
      </div>
      <div class="course-tag">
        <el-tag size="small" effect="dark">学习中</el-tag>
        <span class="page-count">{{ currentPage }}/{{ totalPage }}</span>
      </div>
    </div>

    <div class="status-strip">
      <div class="status-row">
        <span class="status-pill">进度 {{ progressPercent }}%</span>
        <span class="status-pill">{{ isPlay ? '正在讲解' : '已暂停' }}</span>
        <span class="status-pill" v-if="currentNodeTitle">节点 {{ currentNodeTitle }}</span>
        <span class="status-pill" v-if="pageTimelineDuration > 0">{{ formatTime(currentTimelineSec) }} / {{ formatTime(pageTimelineDuration) }}</span>
      </div>
      <div class="status-track" v-if="pageTimelineDuration > 0">
        <div class="status-fill" :style="{ width: timelinePercent + '%' }"></div>
      </div>
      <div class="status-track" v-else>
        <div class="progress-fill" :style="{ width: progressPercent + '%' }"></div>
      </div>
      <div class="status-note" v-if="audioStatusText || activeNodeDuration > 0">
        <span v-if="activeNodeDuration > 0">节点 {{ formatTime(activeNodeElapsedSec) }} / {{ formatTime(activeNodeDuration) }}</span>
        <span>{{ activeNodeTypeLabel }}</span>
        <span v-if="audioStatusText">{{ audioStatusText }}</span>
      </div>
    </div>

    <div class="course-content">
      <!-- 脚本和图片切换 -->
      <div v-if="scriptContent || courseImg" class="content-switcher">
        <el-radio-group v-model="viewMode" size="small" @change="handleViewModeChange">
          <el-radio-button v-if="scriptContent" value="script">讲稿</el-radio-button>
          <el-radio-button v-if="courseImg" value="image">课件</el-radio-button>
        </el-radio-group>
      </div>

      <!-- 脚本视图 -->
      <StudentScriptViewer
        v-if="viewMode === 'script' && scriptContent"
        :script-content="scriptContent"
        :is-loading="isScriptLoading"
        @error="handleScriptError"
        @warning="handleScriptWarning"
      />

      <!-- 课件图片视图 -->
      <div v-if="viewMode === 'image' || !scriptContent" class="course-image-view">
        <img v-if="courseImg" :src="courseImg" alt="课件内容" class="course-img" />
        <div v-else class="no-courseware">当前没有可预览课件，请联系教师先发布课件</div>
        <div
          v-if="tracePoint"
          class="trace-highlight"
          :style="{ top: traceTop + 'px', left: traceLeft + 'px' }"
        ></div>
      </div>
    </div>

    <div class="course-control">
      <el-button @click="$emit('prev-page')" icon="el-icon-arrow-left" size="small">上一页</el-button>
      <el-button @click="$emit('toggle-play')" :icon="isPlay ? 'el-icon-pause' : 'el-icon-play'" size="small">
        {{ isPlay ? '暂停' : '播放' }}
      </el-button>
      <el-button @click="$emit('toggle-tts')" :type="ttsEnabled ? 'primary' : 'default'" plain size="small">
        {{ ttsEnabled ? '语音已开' : '语音已关' }}
      </el-button>
      <el-button @click="$emit('speak-current-node')" icon="el-icon-microphone" size="small" plain>
        朗读当前节点
      </el-button>
      <el-button @click="$emit('next-page')" icon="el-icon-arrow-right" size="small">下一页</el-button>
    </div>
    <div class="control-tip">系统会自动记录到当前页，下次可直接续学。</div>
  </div>
</template>

<script setup>
/* eslint-disable no-undef */
import { computed, ref } from 'vue'
import StudentScriptViewer from './StudentScriptViewer.vue'

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
  ttsEnabled: {
    type: Boolean,
    default: true
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
  },
  scriptContent: {
    type: String,
    default: ''
  },
  isScriptLoading: {
    type: Boolean,
    default: false
  }
})

defineEmits(['prev-page', 'toggle-play', 'next-page', 'select-node', 'toggle-tts', 'speak-current-node', 'script-error', 'script-warning'])

// 内容切换模式
const viewMode = ref('script') // 'script' 或 'image'

/**
 * 处理内容模式切换
 */
function handleViewModeChange() {
  // 模式已通过 v-model 更新
}

/**
 * 处理脚本渲染错误
 */
function handleScriptError(errors) {
  console.error('[StudentCoursePanel] 脚本错误:', errors)
  // 可选：发出错误事件给父组件
}

/**
 * 处理脚本渲染警告
 */
function handleScriptWarning(warnings) {
  console.warn('[StudentCoursePanel] 脚本警告:', warnings)
  // 可选：发出警告事件给父组件
}

const timelinePercent = computed(() => {
  if (!props.pageTimelineDuration) return 0
  return Math.min(100, Math.max(0, Math.round((props.currentTimelineSec / props.pageTimelineDuration) * 100)))
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
.status-strip {
  margin-bottom: 14px;
  padding: 10px 12px;
  border-radius: 12px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.95) 0%, rgba(244, 250, 247, 0.96) 100%);
  border: 1px solid rgba(148, 163, 184, 0.16);
}
.status-row {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 7px;
}
.status-pill {
  display: inline-flex;
  align-items: center;
  padding: 3px 9px;
  border-radius: 999px;
  font-size: 12px;
  color: #3b5d54;
  background: #edf4f0;
  border: 1px solid #d5e4dc;
}

.status-track,
.progress-track {
  height: 8px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.7);
  overflow: hidden;
}
.status-fill,
.progress-fill {
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #0f766e 0%, #0284c7 100%);
}
.status-note {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 7px;
  font-size: 12px;
  color: #628075;
}
.course-content {
  flex: 1;
  background: linear-gradient(180deg, #ffffff 0%, #f8fbff 100%);
  border-radius: 16px;
  border: 1px solid rgba(226, 232, 240, 0.9);
  display: flex;
  flex-direction: column;
}

.course-image-view {
  position: relative;
  flex: 1;
  min-height: 0;
  padding: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: auto;
}

.course-img {
  display: block;
  width: 100%;
  max-width: 100%;
  height: auto;
  max-height: min(72vh, calc(100vh - 360px));
  object-fit: contain;
  border-radius: 10px;
  border: 1px solid rgba(148, 163, 184, 0.25);
  background: #ffffff;
}

.no-courseware {
  margin: auto;
  color: #64748b;
  font-size: 14px;
  text-align: center;
}

.trace-highlight {
  position: absolute;
  width: 22px;
  height: 22px;
  border-radius: 999px;
  border: 2px solid rgba(239, 68, 68, 0.95);
  background: rgba(239, 68, 68, 0.2);
  box-shadow: 0 0 0 6px rgba(239, 68, 68, 0.12);
  transform: translate(-50%, -50%);
  pointer-events: none;
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