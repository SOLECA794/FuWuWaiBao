<template>
  <div class="tab-content">
    <div class="questions-header" v-if="currentCourseId">
      <h4>提问统计 - {{ currentCourseName }}</h4>
      <div class="filter-bar">
        <span>按页码筛选：</span>
        <select :value="filterPage" class="page-select" @change="$emit('update:filterPage', $event.target.value)">
          <option value="">全部</option>
          <option v-for="page in currentCourseTotalPages" :key="page" :value="page">第{{ page }}页</option>
        </select>
        <button
          v-if="hasUncoveredNodes"
          type="button"
          class="uncovered-filter-btn"
          @click="onlyUncoveredQuestions = !onlyUncoveredQuestions"
        >
          {{ onlyUncoveredQuestions ? '显示全部问题' : '仅看未覆盖节点问题' }}
        </button>
      </div>
    </div>
    <div class="questions-list" v-if="currentCourseId">
      <div
        v-for="q in visibleQuestions"
        :key="q.id"
        class="question-item"
      >
        <div class="question-meta">
          <span class="student-id">学生 {{ q.studentId }}</span>
          <span class="page-tag">第{{ q.page }}页</span>
          <button class="node-tag clickable" :class="{ uncovered: isUncoveredNode(q.nodeId) }" v-if="q.nodeId" @click="$emit('focus-node', q)">
            {{ q.nodeTitle || q.nodeId }}
          </button>
          <span class="uncovered-flag" v-if="isUncoveredNode(q.nodeId)">未覆盖</span>
          <span class="time">{{ q.time }}</span>
        </div>
        <div class="question-content">{{ q.content }}</div>
        <div class="answer-content" v-if="q.answer">
          <span class="answer-label">AI 回复：</span>{{ q.answer }}
        </div>
      </div>
      <div v-if="visibleQuestions.length === 0" class="empty-tip">暂无提问记录</div>
    </div>
    <div v-else class="empty-tip">请先选择一个课件查看提问统计</div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  currentCourseId: {
    type: String,
    default: ''
  },
  currentCourseName: {
    type: String,
    default: ''
  },
  currentCourseTotalPages: {
    type: Number,
    default: 0
  },
  filterPage: {
    type: [String, Number],
    default: ''
  },
  filteredQuestions: {
    type: Array,
    default: () => []
  },
  uncoveredNodeIds: {
    type: Array,
    default: () => []
  }
})

const onlyUncoveredQuestions = ref(false)

const uncoveredNodeIdSet = computed(() => {
  return new Set((props.uncoveredNodeIds || []).map(item => String(item || '').trim()).filter(Boolean))
})

const hasUncoveredNodes = computed(() => uncoveredNodeIdSet.value.size > 0)

const visibleQuestions = computed(() => {
  const base = Array.isArray(props.filteredQuestions) ? props.filteredQuestions : []
  if (!onlyUncoveredQuestions.value || uncoveredNodeIdSet.value.size === 0) {
    return base
  }
  return base.filter(item => uncoveredNodeIdSet.value.has(String(item?.nodeId || '').trim()))
})

const isUncoveredNode = (nodeId) => uncoveredNodeIdSet.value.has(String(nodeId || '').trim())

defineEmits(['update:filterPage', 'focus-node'])
</script>

<style scoped>
.tab-content {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 18px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}
.questions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.filter-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #64748b;
}
.page-select {
  border: 1px solid #dbe3ef;
  border-radius: 8px;
  padding: 6px 10px;
}
.questions-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.question-item {
  background: #F4F7F7;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 14px;
}
.question-meta {
  display: flex;
  gap: 12px;
  color: #64748b;
  font-size: 12px;
  margin-bottom: 8px;
}

.node-tag {
  color: #2f605a;
  background: #e7f2ee;
  border: 1px solid #cde2db;
  padding: 1px 8px;
  border-radius: 999px;
}

.node-tag.clickable {
  cursor: pointer;
}

.node-tag.uncovered {
  border-color: #c0841a;
  border-style: dashed;
  color: #8a4b10;
  background: #fff4e5;
}

.uncovered-flag {
  color: #b45309;
  background: #fff7ed;
  border: 1px solid #fed7aa;
  padding: 1px 6px;
  border-radius: 999px;
}

.uncovered-filter-btn {
  border: 1px solid #cbd5e1;
  background: #f8fafc;
  color: #334155;
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 12px;
  cursor: pointer;
}
.question-content {
  color: #0f172a;
  line-height: 1.7;
}
.answer-content {
  margin-top: 8px;
  color: #334155;
  line-height: 1.6;
}
.answer-label {
  color: #2F605A;
  font-weight: 600;
}
.empty-tip {
  text-align: center;
  color: #64748b;
}
</style>