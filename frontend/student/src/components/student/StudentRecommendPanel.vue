<template>
  <section class="recommend-panel">
    <div class="recommend-header">
      <div>
        <div class="recommend-kicker">学习推荐</div>
        <h3>关键词资源推荐（B站 / 51教习）</h3>
        <p>基于当前课程与节点关键词，推荐可直接跳转学习的视频与配套资源。</p>
      </div>
    </div>

    <div class="recommend-toolbar">
      <el-input
        v-model.trim="keyword"
        clearable
        placeholder="输入关键词，如：快速排序 分区思想"
        @keyup.enter="searchResources"
      />
      <el-button type="primary" :loading="loading" @click="searchResources">搜索推荐</el-button>
    </div>

    <div class="keyword-hint" v-if="defaultKeyword">
      建议关键词：{{ defaultKeyword }}
      <el-button text size="small" @click="useDefaultKeyword">使用</el-button>
    </div>

    <div v-if="errorMsg" class="recommend-error">
      {{ errorMsg }}
      <el-button text size="small" @click="searchResources">重试</el-button>
    </div>

    <div v-if="loading && !resourceList.length" class="recommend-loading">正在为你匹配资源...</div>

    <div v-else-if="resourceList.length" class="recommend-list">
      <article class="recommend-item" v-for="item in resourceList" :key="item.id">
        <header>
          <h4>{{ item.title }}</h4>
          <span class="source-tag">{{ item.source || 'Bilibili' }}</span>
        </header>
        <p class="reason">{{ item.recommend_reason || '该资源与当前学习节点相关，可用于补充理解。' }}</p>
        <div class="actions">
          <el-button type="primary" size="small" :disabled="!item.url" @click="openResource(item)">去查看</el-button>
          <span class="meta" v-if="item.duration">时长：{{ item.duration }}</span>
        </div>
      </article>
    </div>

    <div v-else-if="!loading" class="recommend-empty">
      暂无推荐结果，换个关键词试试。
    </div>
  </section>
</template>

<script setup>
/* eslint-disable no-undef */
import { computed, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { studentV1Api } from '../../services/v1'

const props = defineProps({
  courseName: {
    type: String,
    default: ''
  },
  currentNodeTitle: {
    type: String,
    default: ''
  },
  currentPage: {
    type: Number,
    default: 1
  }
})

const loading = ref(false)
const errorMsg = ref('')
const keyword = ref('')
const resourceList = ref([])

const defaultKeyword = computed(() => {
  const parts = [
    String(props.currentNodeTitle || '').trim(),
    String(props.courseName || '').trim(),
    props.currentPage ? `第${props.currentPage}页` : ''
  ].filter(Boolean)
  return parts.join(' ')
})

const normalizeResource = (item, index) => ({
  id: item?.id || item?.ID || `recommend_${Date.now()}_${index}`,
  title: item?.title || item?.Title || '未命名资源',
  source: item?.source || item?.Source || 'Bilibili',
  duration: item?.duration || '',
  recommend_reason: item?.fit_reason || item?.recommend_reason || item?.reason || item?.Reason || '',
  url: item?.url || item?.link || item?.Link || ''
})

const useDefaultKeyword = () => {
  if (!defaultKeyword.value) return
  keyword.value = defaultKeyword.value
}

const searchResources = async () => {
  const query = String(keyword.value || '').trim() || defaultKeyword.value
  if (!query) {
    ElMessage.warning('请输入关键词后再搜索')
    return
  }

  loading.value = true
  errorMsg.value = ''

  try {
    const data = await studentV1Api.recommend.fetchRecommendedResources({
      keyword: query,
      type: '网课',
      source_preference: ['Bilibili', '51jiaoxi'],
      page: 1,
      pageSize: 8
    })

    const normalized = (data.list || []).map(normalizeResource)
    resourceList.value = normalized.slice(0, 8)
  } catch (error) {
    resourceList.value = []
    errorMsg.value = error?.message || '推荐资源拉取失败，请稍后重试。'
  } finally {
    loading.value = false
  }
}

const openResource = (item) => {
  if (!item?.url) {
    ElMessage.warning('该资源缺少可访问链接')
    return
  }
  window.open(item.url, '_blank', 'noopener,noreferrer')
}

watch(defaultKeyword, (next) => {
  if (!keyword.value && next) {
    keyword.value = next
  }
}, { immediate: true })
</script>

<style scoped>
.recommend-panel {
  border: 1px solid #d8e5de;
  border-radius: 18px;
  background: linear-gradient(180deg, #ffffff 0%, #f6faf8 100%);
  padding: 14px;
}

.recommend-kicker {
  font-size: 11px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  font-weight: 700;
  color: #6a8278;
}

.recommend-header h3 {
  margin-top: 4px;
  font-size: 18px;
  color: #23463f;
}

.recommend-header p {
  margin-top: 6px;
  font-size: 13px;
  color: #6f867d;
}

.recommend-toolbar {
  margin-top: 12px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 10px;
}

.keyword-hint {
  margin-top: 10px;
  font-size: 12px;
  color: #5f7970;
  display: flex;
  align-items: center;
  gap: 6px;
}

.recommend-error {
  margin-top: 10px;
  border: 1px solid #f2c1c1;
  background: #fff7f7;
  color: #a94442;
  border-radius: 10px;
  padding: 8px 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.recommend-loading,
.recommend-empty {
  margin-top: 12px;
  border: 1px dashed #cfddd5;
  border-radius: 12px;
  padding: 12px;
  color: #70857c;
  font-size: 13px;
}

.recommend-list {
  margin-top: 12px;
  display: grid;
  gap: 10px;
}

.recommend-item {
  border: 1px solid #d7e5dd;
  border-radius: 12px;
  background: #fff;
  padding: 10px;
}

.recommend-item header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}

.recommend-item h4 {
  font-size: 15px;
  color: #274d46;
}

.source-tag {
  font-size: 11px;
  border-radius: 999px;
  border: 1px solid #d1e2da;
  padding: 2px 8px;
  color: #42665d;
  background: #eef5f1;
}

.reason {
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.55;
  color: #577068;
}

.actions {
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.meta {
  font-size: 12px;
  color: #6d847b;
}

@media (max-width: 768px) {
  .recommend-toolbar {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
