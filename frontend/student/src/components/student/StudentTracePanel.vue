<template>
  <div class="panel-box">
    <p>点击课件任意位置 → 圈选知识点 → 精准提问</p>
    <el-button type="primary" plain size="small" @click="$emit('open-trace-mode')">开启溯源定位</el-button>
    <div v-if="traceLog" class="trace-log">{{ traceLog }}</div>
  </div>
</template>

<script setup>
/* eslint-disable no-undef */
defineProps({
  traceLog: {
    type: String,
    default: ''
  }
})

defineEmits(['open-trace-mode'])
</script>

<style scoped>
.panel-box {
  background: rgba(255, 255, 255, 0.96);
  border-radius: 14px;
  padding: 16px;
  border: 1px solid #e6ecf5;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.06);
}
.trace-log {
  margin-top: 12px;
  padding: 10px;
  background: #fff7e6;
  border-radius: 6px;
  font-size: 13px;
  color: #d48806;
}
</style>