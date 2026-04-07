<template>
  <div class="graph-panel-shell">
    <div class="graph-panel-head">
      <div>
        <h4>知识图谱维护</h4>
        <p>同步关系边、扫描引用健康、一键修复孤儿节点引用。</p>
      </div>
      <button class="close-btn" @click="$emit('close')">关闭</button>
    </div>

    <div class="action-row">
      <button class="action-btn primary" :disabled="syncLoading" @click="handleSync">
        {{ syncLoading ? '同步中...' : '同步图谱' }}
      </button>
      <button class="action-btn" :disabled="scanLoading" @click="handleScan">
        {{ scanLoading ? '扫描中...' : '扫描健康' }}
      </button>
      <button class="action-btn warn" :disabled="repairLoading || !scanReport?.hasOrphans" @click="handleRepair">
        {{ repairLoading ? '修复中...' : '修复引用' }}
      </button>
    </div>

    <transition name="fade-slide">
      <div v-if="summaryVisible" class="summary-grid">
        <div class="metric-card">
          <span>当前课件</span>
          <strong :title="courseName">{{ courseName || '未命名课件' }}</strong>
        </div>
        <div class="metric-card">
          <span>图谱关系边</span>
          <strong>{{ edgeCount }}</strong>
        </div>
        <div class="metric-card" :class="{ danger: orphanCount > 0 }">
          <span>孤儿节点引用</span>
          <strong>{{ orphanCount }}</strong>
        </div>
        <div class="metric-card">
          <span>涉及来源表</span>
          <strong>{{ bucketCount }}</strong>
        </div>
      </div>
    </transition>

    <transition-group name="chip-pop" tag="div" class="orphan-chip-list" v-if="scanReport?.unionOrphanNodeIds?.length">
      <span v-for="id in scanReport.unionOrphanNodeIds" :key="id" class="orphan-chip">{{ id }}</span>
    </transition-group>

    <div class="result-box" v-if="lastMessage">
      {{ lastMessage }}
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { teacherCoursewareApi } from '../../services/v1/coursewareApi'

const props = defineProps({
  courseId: {
    type: String,
    required: true
  },
  courseName: {
    type: String,
    default: ''
  }
})

defineEmits(['close'])

const syncLoading = ref(false)
const scanLoading = ref(false)
const repairLoading = ref(false)

const syncPayload = ref(null)
const scanReport = ref(null)
const lastMessage = ref('')

const summaryVisible = computed(() => Boolean(syncPayload.value || scanReport.value))
const edgeCount = computed(() => Number(syncPayload.value?.edgeCount || 0))
const orphanCount = computed(() => Number(scanReport.value?.unionOrphanNodeIds?.length || 0))
const bucketCount = computed(() => Number(scanReport.value?.buckets?.length || 0))

const handleSync = async () => {
  if (!props.courseId || syncLoading.value) return
  syncLoading.value = true
  lastMessage.value = ''
  try {
    const resp = await teacherCoursewareApi.syncKnowledgeGraph(props.courseId)
    syncPayload.value = resp?.data || {}
    lastMessage.value = `同步完成：共 ${edgeCount.value} 条关系边。`
  } catch (err) {
    lastMessage.value = `同步失败：${err.message || err}`
  } finally {
    syncLoading.value = false
  }
}

const handleScan = async () => {
  if (!props.courseId || scanLoading.value) return
  scanLoading.value = true
  lastMessage.value = ''
  try {
    const resp = await teacherCoursewareApi.getKnowledgeGraphReferenceHealth(props.courseId)
    scanReport.value = resp?.data || null
    if (scanReport.value?.hasOrphans) {
      lastMessage.value = `发现 ${orphanCount.value} 个孤儿引用，建议修复。`
    } else {
      lastMessage.value = '扫描完成：未发现孤儿引用。'
    }
  } catch (err) {
    lastMessage.value = `扫描失败：${err.message || err}`
  } finally {
    scanLoading.value = false
  }
}

const handleRepair = async () => {
  if (!props.courseId || repairLoading.value || !scanReport.value?.hasOrphans) return
  repairLoading.value = true
  lastMessage.value = ''
  try {
    await teacherCoursewareApi.repairKnowledgeGraphReferences(props.courseId, {
      confirm: true,
      nodeIds: scanReport.value.unionOrphanNodeIds || []
    })
    lastMessage.value = '修复完成，正在自动重新扫描...'
    await handleScan()
  } catch (err) {
    lastMessage.value = `修复失败：${err.message || err}`
  } finally {
    repairLoading.value = false
  }
}
</script>

<style scoped>
.graph-panel-shell {
  margin: 10px 12px 12px;
  padding: 12px;
  border-radius: 14px;
  border: 1px solid rgba(92, 166, 143, 0.28);
  background: linear-gradient(180deg, #f8fcfa 0%, #f1f8f4 100%);
  box-shadow: 0 12px 24px rgba(46, 89, 74, 0.12);
}

.graph-panel-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}

.graph-panel-head h4 {
  margin: 0;
  font-size: 15px;
  color: #1f473d;
}

.graph-panel-head p {
  margin: 5px 0 0;
  font-size: 12px;
  line-height: 1.55;
  color: #5f7a70;
}

.close-btn {
  border: 1px solid #cfe3da;
  background: #ffffff;
  color: #2f605a;
  border-radius: 999px;
  padding: 5px 10px;
  font-size: 12px;
  cursor: pointer;
}

.action-row {
  margin-top: 10px;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.action-btn {
  border: 1px solid #b9d7ca;
  border-radius: 10px;
  padding: 8px 10px;
  background: #ffffff;
  color: #2f605a;
  font-size: 12px;
  cursor: pointer;
  transition: transform 0.18s ease, box-shadow 0.18s ease;
}

.action-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 8px 16px rgba(46, 89, 74, 0.12);
}

.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.action-btn.primary {
  background: linear-gradient(180deg, #7dc3ad 0%, #5ca68f 100%);
  color: #ffffff;
  border-color: transparent;
}

.action-btn.warn {
  background: linear-gradient(180deg, #fff8ef 0%, #fff0dc 100%);
  border-color: #f5d5a5;
  color: #925900;
}

.summary-grid {
  margin-top: 10px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}

.metric-card {
  border: 1px solid rgba(92, 166, 143, 0.2);
  border-radius: 10px;
  background: #ffffff;
  padding: 8px;
  min-width: 0;
}

.metric-card span {
  font-size: 11px;
  color: #6a847a;
}

.metric-card strong {
  margin-top: 4px;
  display: block;
  font-size: 14px;
  color: #21483e;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.metric-card.danger {
  border-color: rgba(228, 92, 92, 0.28);
  background: linear-gradient(180deg, #fffaf9 0%, #fff1f0 100%);
}

.orphan-chip-list {
  margin-top: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  max-height: 96px;
  overflow: auto;
}

.orphan-chip {
  padding: 3px 8px;
  border-radius: 999px;
  border: 1px solid rgba(228, 92, 92, 0.3);
  background: #fff5f4;
  color: #b14a4a;
  font-size: 11px;
}

.result-box {
  margin-top: 10px;
  border-radius: 9px;
  padding: 8px 10px;
  font-size: 12px;
  line-height: 1.5;
  color: #365b52;
  border: 1px solid #d2e5dc;
  background: #ffffff;
}

.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.24s ease;
}

.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(5px);
}

.chip-pop-enter-active,
.chip-pop-leave-active {
  transition: all 0.22s ease;
}

.chip-pop-enter-from,
.chip-pop-leave-to {
  opacity: 0;
  transform: scale(0.92);
}
</style>
