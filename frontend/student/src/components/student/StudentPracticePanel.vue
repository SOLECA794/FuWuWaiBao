<template>
  <div class="practice-panel">
    <header class="practice-header">
      <div>
        <p class="kicker">随堂练习</p>
        <h2>{{ courseName || '当前课程' }} · 课堂测验</h2>
        <p class="desc">完成后可查看批改结果、正确率和错题解析，并一键加入错题本。</p>
      </div>
      <div class="meta">
        <span>题量 {{ questions.length }}</span>
        <span>当前节点 {{ currentNodeTitle || '未定位' }}</span>
      </div>
    </header>

    <section class="question-list" v-if="!submitted">
      <div v-if="loading" class="loading-tip">正在生成本页练习题...</div>
      <article class="question-card" v-for="(q, idx) in questions" :key="q.id">
        <h3>{{ idx + 1 }}. {{ q.stem }}</h3>
        <div class="options">
          <button
            v-for="opt in q.options"
            :key="opt.key"
            class="option-btn"
            :class="{ active: answers[q.id] === opt.key }"
            @click="selectAnswer(q.id, opt.key)"
          >
            {{ opt.key }}. {{ opt.text }}
          </button>
        </div>
      </article>

      <div class="submit-row">
        <button class="btn ghost" :disabled="loading" @click="reset">重置作答</button>
        <button class="btn primary" :disabled="loading" @click="submit">提交批改</button>
      </div>
    </section>

    <section class="result-panel" v-else>
      <div class="score-grid">
        <div class="score-card">
          <span>总题数</span>
          <strong>{{ result.total }}</strong>
        </div>
        <div class="score-card">
          <span>答对题数</span>
          <strong>{{ result.correct }}</strong>
        </div>
        <div class="score-card">
          <span>正确率</span>
          <strong>{{ result.accuracy }}%</strong>
        </div>
      </div>

      <div class="wrong-list" v-if="result.wrongs.length">
        <h3>错题解析</h3>
        <article class="wrong-item" v-for="item in result.wrongs" :key="item.id">
          <div class="wrong-head">
            <strong>{{ item.stem }}</strong>
            <span>你的答案：{{ item.userAnswer || '未作答' }}</span>
          </div>
          <p>正确答案：{{ item.correctAnswer }}</p>
          <p>解析：{{ item.analysis }}</p>
        </article>
        <div class="submit-row">
          <button class="btn primary" @click="addWrongSet">加入错题本</button>
          <button class="btn ghost" @click="goToWrongSet">去个人中心查看错题</button>
          <button class="btn ghost" @click="retry">再练一组</button>
        </div>
      </div>

      <div class="all-right" v-else>
        <h3>全部正确，表现很好！</h3>
        <div class="submit-row">
          <button class="btn ghost" @click="goToWrongSet">去个人中心看练习记录</button>
          <button class="btn ghost" @click="retry">再练一组</button>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
/* eslint-disable no-undef */
import { computed, onMounted, ref } from 'vue'
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
  studentId: {
    type: String,
    default: ''
  },
  courseId: {
    type: String,
    default: ''
  },
  nodeId: {
    type: String,
    default: ''
  },
  pageNum: {
    type: Number,
    default: 1
  }
})
const emit = defineEmits(['jump-personal-practice'])

const questions = ref([
  {
    id: 'q1',
    stem: '正应力的常用判定核心是以下哪一项？',
    options: [
      { key: 'A', text: '仅看材料名称' },
      { key: 'B', text: '先看受力方向与截面法向关系' },
      { key: 'C', text: '只看题目图形大小' },
      { key: 'D', text: '先求剪应力再猜测' }
    ],
    answer: 'B',
    analysis: '正应力由法向分量引起，应先判断受力方向与截面法向关系。'
  },
  {
    id: 'q2',
    stem: '构件受拉时，正应力符号通常记为？',
    options: [
      { key: 'A', text: '正值' },
      { key: 'B', text: '负值' },
      { key: 'C', text: '恒为0' },
      { key: 'D', text: '与面积无关' }
    ],
    answer: 'A',
    analysis: '常用工程符号约定中，拉应力为正，压应力为负。'
  },
  {
    id: 'q3',
    stem: '以下哪项最能提高变形题中的判定速度？',
    options: [
      { key: 'A', text: '背答案' },
      { key: 'B', text: '每题都完整推导一遍' },
      { key: 'C', text: '建立“受力-变形-判定”快速检查清单' },
      { key: 'D', text: '随机选择计算公式' }
    ],
    answer: 'C',
    analysis: '建立固定判定清单可显著降低认知切换成本并减少漏判。'
  }
])

const taskId = ref('')
const answers = ref({})
const submitted = ref(false)
const loading = ref(false)
const useRemoteData = ref(false)
const submitResult = ref(null)

const localResult = computed(() => {
  const wrongs = questions.value
    .filter((q) => answers.value[q.id] !== q.answer)
    .map((q) => ({
      id: q.id,
      stem: q.stem,
      userAnswer: answers.value[q.id] || '',
      correctAnswer: q.answer,
      analysis: q.analysis
    }))
  const total = questions.value.length
  const correct = total - wrongs.length
  const accuracy = total > 0 ? Math.round((correct / total) * 100) : 0
  return { total, correct, accuracy, wrongs }
})

const result = computed(() => {
  if (submitResult.value) {
    const total = Number(submitResult.value.totalCount || 0)
    const correct = Number(submitResult.value.correctCount || 0)
    const score = Number(submitResult.value.score || 0)
    const details = Array.isArray(submitResult.value.details) ? submitResult.value.details : []
    const wrongs = details
      .filter((item) => item && item.correct === false)
      .map((item) => ({
        id: item.questionId,
        stem: item.content || '题目内容',
        userAnswer: item.userAnswer || '',
        correctAnswer: item.correctAnswer || '',
        analysis: item.explanation || item.aiComment || '暂无解析'
      }))
    return {
      total,
      correct,
      accuracy: total > 0 ? Math.round((correct / total) * 100) : score,
      wrongs
    }
  }
  return localResult.value
})

const selectAnswer = (id, key) => {
  answers.value = { ...answers.value, [id]: key }
}

const reset = () => {
  answers.value = {}
}

const buildSubmitPayload = () => {
  return questions.value.map((q) => ({
    questionId: q.id,
    userAnswer: answers.value[q.id] || ''
  }))
}

const submit = async () => {
  if (!questions.value.length) {
    ElMessage.warning('暂无可提交题目')
    return
  }
  if (useRemoteData.value && taskId.value) {
    try {
      loading.value = true
      const payload = await studentV1Api.coursewares.submitPractice({
        taskId: taskId.value,
        studentId: props.studentId,
        answers: buildSubmitPayload()
      })
      submitResult.value = payload?.data || null
      submitted.value = true
      ElMessage.success('已完成批改并同步练习记录')
      return
    } catch (error) {
      ElMessage.warning(`接口批改失败，已切换本地批改：${error.message}`)
    } finally {
      loading.value = false
    }
  }
  submitted.value = true
  ElMessage.success('已完成批改')
}

const addWrongSet = () => {
  if (useRemoteData.value) {
    ElMessage.success('错题已进入系统错题本，可在个人中心-练习记录查看')
    return
  }
  ElMessage.success(`已加入错题本（演示模式 ${result.value.wrongs.length}题）`)
}

const goToWrongSet = () => {
  emit('jump-personal-practice')
}

const mapQuestionFromApi = (q, idx) => {
  const options = Array.isArray(q?.options) ? q.options : []
  return {
    id: String(q?.questionId || `q${idx + 1}`),
    stem: String(q?.content || `题目 ${idx + 1}`),
    options: options.map((text, i) => ({
      key: String.fromCharCode(65 + i),
      text: String(text || '')
    })),
    answer: '',
    analysis: ''
  }
}

const loadPracticeTask = async () => {
  if (!props.studentId || !props.courseId) return
  try {
    loading.value = true
    const resp = await studentV1Api.coursewares.generatePractice({
      studentId: props.studentId,
      courseId: props.courseId,
      nodeId: props.nodeId,
      pageNum: props.pageNum || 1,
      difficulty: 2,
      count: 3
    })
    const data = resp?.data || {}
    const remoteQuestions = Array.isArray(data.questions) ? data.questions : []
    if (!remoteQuestions.length) {
      useRemoteData.value = false
      ElMessage.warning('后端暂未返回题目，已使用演示题目')
      return
    }
    taskId.value = String(data.taskId || '')
    questions.value = remoteQuestions.map(mapQuestionFromApi)
    useRemoteData.value = taskId.value !== ''
  } catch (error) {
    useRemoteData.value = false
    ElMessage.warning(`练习题加载失败，已使用演示题目：${error.message}`)
  } finally {
    loading.value = false
  }
}

const retry = async () => {
  submitted.value = false
  answers.value = {}
  submitResult.value = null
  if (useRemoteData.value) {
    await loadPracticeTask()
  }
}

onMounted(async () => {
  await loadPracticeTask()
})
</script>

<style scoped>
.practice-panel { border: 1px solid #d9e7df; border-radius: 18px; background: #fff; padding: 14px; }
.practice-header { display: flex; justify-content: space-between; gap: 12px; }
.kicker { margin: 0; color: #6a8278; font-size: 12px; }
.practice-header h2 { margin: 4px 0; font-size: 20px; color: #23463f; }
.desc { margin: 0; color: #688379; font-size: 13px; }
.meta { display: grid; gap: 6px; font-size: 12px; color: #46655d; align-content: start; }
.question-list { margin-top: 12px; display: grid; gap: 10px; }
.loading-tip { font-size: 12px; color: #6d857b; }
.question-card { border: 1px solid #e0ebe5; border-radius: 12px; padding: 10px; background: #f9fcfa; }
.question-card h3 { margin: 0 0 8px; color: #24453f; font-size: 15px; }
.options { display: grid; gap: 8px; }
.option-btn { border: 1px solid #d2e4db; background: #fff; border-radius: 10px; padding: 8px; text-align: left; cursor: pointer; color: #3f5f56; }
.option-btn.active { border-color: #2f605a; background: #eaf4ef; color: #1f433b; font-weight: 700; }
.submit-row { margin-top: 8px; display: flex; gap: 8px; }
.btn { border: 1px solid #cfe1d8; border-radius: 10px; padding: 8px 12px; cursor: pointer; background: #fff; color: #35564d; }
.btn:disabled { opacity: 0.6; cursor: not-allowed; }
.btn.primary { border-color: #2f605a; background: #2f605a; color: #fff; }
.score-grid { margin-top: 12px; display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 10px; }
.score-card { border: 1px solid #deebe4; border-radius: 12px; padding: 10px; background: #f8fcfa; }
.score-card span { font-size: 12px; color: #6d857b; }
.score-card strong { display: block; margin-top: 6px; font-size: 24px; color: #23463f; }
.wrong-list { margin-top: 12px; display: grid; gap: 10px; }
.wrong-item { border: 1px solid #f0d7d7; border-radius: 12px; background: #fff8f8; padding: 10px; color: #5b4343; }
.wrong-head { display: flex; justify-content: space-between; gap: 8px; font-size: 13px; }
.all-right { margin-top: 12px; border: 1px solid #d5e7de; border-radius: 12px; background: #f4faf7; padding: 12px; color: #2f605a; }
</style>
