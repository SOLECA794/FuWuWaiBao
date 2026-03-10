<template>
  <div class="course-card">
    <div class="course-header">
      <h3>{{ currentCourseName || '暂无课程' }} - 第{{ currentPage }}页</h3>
      <div class="course-tag">
        <el-tag size="small">课件学习</el-tag>
        <span class="page-count">{{ currentPage }}/{{ totalPage }}</span>
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
          <strong>{{ node.title || node.node_id }}</strong>
          <span>{{ node.text }}</span>
        </button>
      </div>
    </div>
    <div class="control-tip">支持自动记录断点，切页后将实时同步学习进度</div>
  </div>
</template>

<script setup>
/* eslint-disable no-undef */
defineProps({
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
</script>

<style scoped>
.course-card {
  background: rgba(255, 255, 255, 0.96);
  border-radius: 14px;
  padding: 20px;
  height: 100%;
  border: 1px solid #e6ecf5;
  box-shadow: 0 10px 24px rgba(15, 23, 42, 0.08);
  display: flex;
  flex-direction: column;
}
.course-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.course-header h3 {
  font-size: 16px;
  color: #333;
}
.course-tag {
  display: flex;
  gap: 8px;
  align-items: center;
}
.page-count {
  font-size: 13px;
  color: #666;
}
.course-content {
  flex: 1;
  background: #f8fbff;
  border-radius: 8px;
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
  background: #f8fbff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
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
  background: #fff;
  border-radius: 10px;
  padding: 10px 12px;
  text-align: left;
  cursor: pointer;
}
.node-chip strong {
  display: block;
  margin-bottom: 4px;
  font-size: 12px;
  color: #2563eb;
}
.node-chip span {
  display: block;
  font-size: 13px;
  color: #334155;
  line-height: 1.6;
}
.node-chip.active {
  border-color: #2563eb;
  box-shadow: 0 0 0 1px rgba(37, 99, 235, 0.15);
  background: #eff6ff;
}
</style>