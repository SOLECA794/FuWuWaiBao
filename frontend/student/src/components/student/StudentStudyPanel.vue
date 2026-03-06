<template>
  <div class="panel-box">
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
      <div class="title">AI 诊断薄弱点（点击可学习）</div>
      <el-tag
        v-for="point in weakPointTags"
        :key="point.id"
        type="danger"
        size="small"
        style="cursor: pointer; margin: 3px"
        @click="$emit('start-weak-point', point)"
      >
        {{ point.name }}
      </el-tag>
    </div>

    <div v-if="currentExplain" class="explain-card">
      <h4>📘 {{ currentWeakPoint }} · 知识点讲解</h4>
      <p>{{ currentExplain }}</p>
      <el-button type="primary" @click="$emit('generate-test')" style="margin-top: 10px">
        已学会，开始习题检测
      </el-button>
    </div>

    <div v-if="currentTest" class="test-card">
      <h4>📝 练习题：{{ currentWeakPoint }}</h4>
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
defineProps({
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
</script>

<style scoped>
.panel-box {
  background: rgba(255, 255, 255, 0.96);
  border-radius: 14px;
  padding: 16px;
  border: 1px solid #e6ecf5;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.06);
}
.data-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 16px;
}
.data-item {
  background: #f8fbff;
  padding: 16px;
  border-radius: 8px;
  text-align: center;
  border: 1px solid #e6ecf5;
}
.data-item .num {
  font-size: 20px;
  font-weight: bold;
  color: #1989fa;
}
.data-item .label {
  font-size: 13px;
  color: #666;
  margin-top: 4px;
}
.weak-point {
  margin-top: 10px;
}
.weak-point .title {
  font-size: 14px;
  margin-bottom: 8px;
  font-weight: 500;
}
.explain-card {
  margin-top: 16px;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
  border-left: 4px solid #1989fa;
}
.test-card {
  margin-top: 16px;
  padding: 12px;
  background: #fff7e6;
  border-radius: 8px;
  border-left: 4px solid #faad14;
}
</style>