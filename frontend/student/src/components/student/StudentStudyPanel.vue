<template>
  <div class="panel-box">
    <div class="panel-head">
      <div class="eyebrow">学习画像</div>
      <div class="title">学习数据与薄弱点闭环</div>
    </div>

    <div class="data-grid">
      <div class="data-item">
        <div class="num">{{ learningStats.focusScore }}%</div>
        <div class="label">专注度</div>
      </div>
      <div class="data-item">
        <div class="num">{{ learningStats.totalQuestions }} 次</div>
        <div class="label">本课提问</div>
      </div>
      <div class="data-item">
        <div class="num">{{ learningStats.weakPointCount }} 个</div>
        <div class="label">薄弱点</div>
      </div>
      <div class="data-item">
        <div class="num">{{ learningStats.masteryRate }}%</div>
        <div class="label">掌握率</div>
      </div>
    </div>

    <div class="weak-point">
      <div class="title-row">
        <div class="title">AI 诊断薄弱点（点击进入讲解）</div>
        <span class="tag-pager" v-if="weakPointTags.length">{{ weakTagStart }}-{{ weakTagEnd }} / {{ weakPointTags.length }}</span>
      </div>
      <el-tag
        v-for="point in pagedWeakPointTags"
        :key="point.id"
        type="warning"
        size="small"
        style="cursor: pointer; margin: 3px"
        @click="$emit('start-weak-point', point)"
      >
        {{ point.name }}
      </el-tag>
      <div class="weak-pagination" v-if="weakPointTags.length > weakTagPageSize">
        <el-pagination
          small
          background
          layout="prev, pager, next"
          :current-page="weakTagPage"
          :page-size="weakTagPageSize"
          :total="weakPointTags.length"
          @current-change="weakTagPage = $event"
        />
      </div>
    </div>

    <div v-if="currentExplain" class="explain-card">
      <h4>{{ currentWeakPoint }} · 知识点讲解</h4>
      <p>{{ currentExplain }}</p>
      <el-button type="primary" @click="$emit('generate-test')" style="margin-top: 10px">
        已学会，开始习题检测
      </el-button>
    </div>

    <div v-if="currentTest" class="test-card">
      <h4>练习题：{{ currentWeakPoint }}</h4>
      <p>{{ currentTest.question }}</p>
      <el-button
        v-for="(opt, idx) in currentTest.options"
        :key="idx"
        style="margin: 5px"
        @click="$emit('check-answer', opt)"
      >
        {{ opt }}
      </el-button>
      <div v-if="testResult" :style="{ color: testResult.correct ? 'green' : 'red', marginTop: '10px' }">
        {{ testResult.msg }}
      </div>
      <div v-if="testResult && testResult.analysis">
        <small style="display: block; margin-top: 5px">解析：{{ testResult.analysis }}</small>
      </div>
    </div>
  </div>
</template>

<script setup>
/* eslint-disable no-undef */
import { computed, ref, watch } from 'vue'

const weakTagPage = ref(1)
const weakTagPageSize = 8

const props = defineProps({
  learningStats: {
    type: Object,
    default: () => ({})
  },
  weakPointTags: {
    type: Array,
    default: () => []
  },
  currentExplain: {
    type: String,
    default: ''
  },
  currentWeakPoint: {
    type: String,
    default: ''
  },
  currentTest: {
    type: Object,
    default: null
  },
  testResult: {
    type: Object,
    default: null
  }
})

defineEmits(['start-weak-point', 'generate-test', 'check-answer'])

const pagedWeakPointTags = computed(() => {
  const start = (weakTagPage.value - 1) * weakTagPageSize
  return props.weakPointTags.slice(start, start + weakTagPageSize)
})

const weakTagStart = computed(() => {
  if (!props.weakPointTags.length) return 0
  return (weakTagPage.value - 1) * weakTagPageSize + 1
})

const weakTagEnd = computed(() => {
  if (!props.weakPointTags.length) return 0
  return Math.min(weakTagStart.value + pagedWeakPointTags.value.length - 1, props.weakPointTags.length)
})

watch(() => props.weakPointTags.length, (len) => {
  const pageCount = Math.max(1, Math.ceil(len / weakTagPageSize))
  if (weakTagPage.value > pageCount) {
    weakTagPage.value = pageCount
  }
})
</script>

<style scoped>
.panel-box {
  background: linear-gradient(180deg, #ffffff 0%, #f7faf8 100%);
  border-radius: 20px;
  padding: 16px;
  border: 1px solid #d9e7df;
  box-shadow: 0 16px 30px rgba(33, 61, 54, 0.08);
}

.panel-head {
  margin-bottom: 12px;
}

.eyebrow {
  font-size: 12px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #6f867d;
  font-weight: 700;
}

.panel-head .title {
  font-size: 18px;
  font-weight: 700;
  color: #23463f;
  margin-top: 4px;
}

.data-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 16px;
}
.data-item {
  background: #f3f8f5;
  padding: 16px;
  border-radius: 12px;
  text-align: center;
  border: 1px solid #d8e7df;
}

.data-item .num {
  font-size: 20px;
  font-weight: 700;
  color: #2f605a;
}

.data-item .label {
  font-size: 13px;
  color: #667d73;
  margin-top: 4px;
}

.weak-point {
  margin-top: 10px;
}

.title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.weak-point .title {
  font-size: 14px;
  margin-bottom: 8px;
  font-weight: 700;
  color: #2b4d45;
}

.tag-pager {
  font-size: 12px;
  color: #6f867d;
}

.weak-pagination {
  margin-top: 8px;
  display: flex;
  justify-content: flex-end;
}

.explain-card {
  margin-top: 16px;
  padding: 14px;
  background: #f4f8f6;
  border-radius: 12px;
  border-left: 4px solid #2f605a;
}

.explain-card h4 {
  color: #23463f;
}

.explain-card p {
  margin-top: 8px;
  line-height: 1.7;
  color: #495f56;
}

.test-card {
  margin-top: 16px;
  padding: 14px;
  background: #fffaef;
  border-radius: 12px;
  border-left: 4px solid #cc9a2c;
}

.test-card h4 {
  color: #7a5a1a;
}

.test-card p {
  margin-top: 8px;
  margin-bottom: 8px;
  line-height: 1.6;
}

@media (max-width: 720px) {
  .data-grid {
    grid-template-columns: 1fr;
  }

  .title-row {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>