<template>
  <div class="analytics-workbench">
    <header class="head-area">
      <div>
        <p class="eyebrow">我的学习画像</p>
        <h3>AI 学情诊断・精准学习指引</h3>
        <p class="head-tip">基于你的学习行为，AI 自动诊断学习状态，帮你精准提升</p>
      </div>
      <div class="head-controls">
        <el-select v-model="selectedCourse" size="small" placeholder="切换课程" style="width: 180px">
          <el-option label="全课程数据" value="all" />
          <el-option label="本课程数据" value="current" />
        </el-select>
        <el-button size="small" @click="refreshData">刷新数据</el-button>
        <el-button size="small" @click="exportReport">导出报告</el-button>
      </div>
    </header>

    <section class="overview-cards">
      <article
        v-for="card in overviewCards"
        :key="card.id"
        class="card"
        :class="{ highlight: card.highlight }"
        @click="activeTab = card.tabName"
      >
        <div class="card-head">
          <span class="card-emoji">{{ card.emoji }}</span>
          <span class="card-label">{{ card.label }}</span>
        </div>
        <div class="card-number">{{ card.number }}</div>
        <div class="card-sub">{{ card.sub }}</div>
        <div class="card-trend" :class="[card.trend > 0 ? 'up' : 'down']">
          {{ card.trend > 0 ? '↑' : '↓' }}{{ Math.abs(card.trend) }}%
        </div>
      </article>
    </section>

    <section class="diagnostic-area">
      <el-tabs v-model="activeTab" animated>
        <el-tab-pane label="📈 知识点掌握" name="mastery">
          <div class="tab-content">
            <div class="mastery-grid">
              <div class="mastery-chart">
                <h4>章节掌握度热力图</h4>
                <div class="heatmap">
                  <div v-for="ch in chapters" :key="ch.id" class="heatmap-row">
                    <span class="ch-name">{{ ch.name }}</span>
                    <div class="heat-items">
                      <div
                        v-for="point in ch.points"
                        :key="point.id"
                        class="heat-item"
                        :style="{ backgroundColor: getMasteryColor(point.mastery) }"
                        :title="`${point.name}: ${point.mastery}%`"
                        @click="showPoiDetail(point)"
                      ></div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="mastery-pie">
                <h4>掌握度分布</h4>
                <div class="pie-visual">
                  <div class="pie-item mastered">
                    <strong>{{ masteredCount }}</strong>
                    <small>已掌握</small>
                  </div>
                  <div class="pie-item basic">
                    <strong>{{ basicCount }}</strong>
                    <small>基本掌握</small>
                  </div>
                  <div class="pie-item weak">
                    <strong>{{ weakCount }}</strong>
                    <small>未掌握</small>
                  </div>
                </div>
              </div>
            </div>
            <div class="poi-list">
              <h4>知识点详情</h4>
              <table class="poi-table">
                <thead>
                  <tr>
                    <th>知识点</th>
                    <th>章节</th>
                    <th>掌握度</th>
                    <th>习题正确率</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="point in poiList.slice(0, 8)" :key="point.id">
                    <td><strong>{{ point.name }}</strong></td>
                    <td>{{ point.chapter }}</td>
                    <td>
                      <el-progress :percentage="point.mastery" :status="masteryStatus(point.mastery)" :show-text="false" />
                    </td>
                    <td>{{ point.qCorrectRate }}%</td>
                    <td>
                      <el-button size="small" text type="primary" @click="jumpToLearn(point)">学习</el-button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div class="ai-diagnosis">
              <h4>AI 诊断·知识点掌握情况</h4>
              <p>{{ aiDiagnosisText.mastery }}</p>
              <el-button type="primary" @click="generateReviewPlan">生成复习计划</el-button>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="⚠️ 薄弱点诊断" name="weakness">
          <div class="tab-content">
            <div class="weakness-top5">
              <h4>薄弱知识点 TOP5</h4>
              <div class="top5-cards">
                <article v-for="(item, idx) in weaknessTop5" :key="item.id" class="weakness-card">
                  <div class="rank-badge">{{ idx + 1 }}</div>
                  <div class="card-body">
                    <div>
                      <strong>{{ item.name }}</strong>
                      <p class="cause">原因：{{ item.cause }}</p>
                    </div>
                    <div class="score">{{ item.mastery }}%</div>
                  </div>
                  <div class="card-actions">
                    <el-button size="small" text @click="askAbout(item)">AI讲解</el-button>
                    <el-button size="small" text @click="practiceProblem(item)">练习</el-button>
                  </div>
                </article>
              </div>
            </div>
            <div class="weakness-radar">
              <h4>薄弱原因分析</h4>
              <div class="radar-sim">
                <div class="radar-item">
                  <span>习题正确率低</span>
                  <div class="radar-bar" style="width: 65%"></div>
                </div>
                <div class="radar-item">
                  <span>学习时长不足</span>
                  <div class="radar-bar" style="width: 48%"></div>
                </div>
                <div class="radar-item">
                  <span>提问频繁</span>
                  <div class="radar-bar" style="width: 72%"></div>
                </div>
                <div class="radar-item">
                  <span>错题重做率低</span>
                  <div class="radar-bar" style="width: 54%"></div>
                </div>
                <div class="radar-item">
                  <span>需要重讲</span>
                  <div class="radar-bar" style="width: 38%"></div>
                </div>
              </div>
            </div>
            <div class="weakness-plan">
              <h4>AI 提升方案</h4>
              <div class="plan-step" v-for="step in improvementPlan" :key="step.id">
                <span class="step-no">{{ step.step }}</span>
                <div>
                  <strong>{{ step.title }}</strong>
                  <p>{{ step.desc }}</p>
                </div>
              </div>
              <el-button type="primary" @click="addToTodo">加入待办任务</el-button>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="📝 错题分析" name="mistake">
          <div class="tab-content">
            <div class="practice-overview">
              <div class="stat-item">
                <span class="stat-num">{{ totalPractices }}</span>
                <span class="stat-label">总练习数</span>
              </div>
              <div class="stat-item">
                <span class="stat-num">{{ avgCorrectRate }}%</span>
                <span class="stat-label">平均正确率</span>
              </div>
              <div class="stat-item">
                <span class="stat-num">{{ mistakeCount }}</span>
                <span class="stat-label">错题总数</span>
              </div>
              <div class="stat-item">
                <span class="stat-num">{{ retakeRate }}%</span>
                <span class="stat-label">错题重做率</span>
              </div>
            </div>
            <div class="trend-chart">
              <h4>正确率趋势</h4>
              <div class="mini-chart">
                <div v-for="day in 7" :key="`trend-${day}`" class="trend-bar" :style="{ height: (50 + Math.random() * 30) + '%' }"></div>
              </div>
            </div>
            <div class="mistake-list">
              <h4>错题详情（最近 5 题）</h4>
              <article v-for="mistake in mistakeList.slice(0, 5)" :key="mistake.id" class="mistake-item">
                <div class="mistake-head">
                  <strong>{{ mistake.stem }}</strong>
                  <span class="mistake-poi">{{ mistake.poi }}</span>
                </div>
                <p class="mistake-cause">错因：{{ mistake.cause }}</p>
                <p class="mistake-answer">正确答案：{{ mistake.answer }}</p>
                <div class="mistake-actions">
                  <el-button size="small" text @click="redoMistake(mistake)">重做</el-button>
                  <el-button size="small" text @click="addToReview(mistake)">加入复习</el-button>
                </div>
              </article>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="⏱️ 学习行为" name="behavior">
          <div class="tab-content">
            <div class="behavior-grid">
              <div class="behavior-chart">
                <h4>学习时长分布</h4>
                <div class="behavior-bars">
                  <div v-for="t of 7" :key="`time-${t}`" class="bar-item">
                    <div class="bar" :style="{ height: (30 + Math.random() * 40) + '%' }"></div>
                    <span>周{{ t }}</span>
                  </div>
                </div>
              </div>
              <div class="activity-dist">
                <h4>学习活动占比</h4>
                <div class="activity-items">
                  <div class="activity-item video">
                    <span>视频学习</span>
                    <div class="bar"><div class="fill" style="width: 42%"></div></div>
                    <span>42%</span>
                  </div>
                  <div class="activity-item practice">
                    <span>习题练习</span>
                    <div class="bar"><div class="fill" style="width: 28%"></div></div>
                    <span>28%</span>
                  </div>
                  <div class="activity-item qa">
                    <span>问答互动</span>
                    <div class="bar"><div class="fill" style="width: 18%"></div></div>
                    <span>18%</span>
                  </div>
                  <div class="activity-item note">
                    <span>记笔记</span>
                    <div class="bar"><div class="fill" style="width: 12%"></div></div>
                    <span>12%</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="ai-behavior-tip">
              <h4>AI 学习习惯建议</h4>
              <p>{{ aiDiagnosisText.behavior }}</p>
              <el-button type="primary" @click="optimizeLearning">优化学习计划</el-button>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="💬 问答详情" name="qa">
          <div class="tab-content">
            <div class="qa-overview">
              <div class="qa-stat">
                <span class="qa-num">{{ totalQa }}</span>
                <span class="qa-label">总提问数</span>
              </div>
              <div class="qa-stat">
                <span class="qa-num">{{ qaResolveRate }}%</span>
                <span class="qa-label">问题解决率</span>
              </div>
              <div class="qa-stat">
                <span class="qa-num">2.3m</span>
                <span class="qa-label">平均响应</span>
              </div>
            </div>
            <div class="qa-cloud">
              <h4>高频提问知识点</h4>
              <div class="word-cloud">
                <span class="word" style="font-size: 24px; font-weight: bold">正应力</span>
                <span class="word" style="font-size: 18px">切应力</span>
                <span class="word" style="font-size: 20px; font-weight: bold">应变</span>
                <span class="word" style="font-size: 14px">弹性模量</span>
                <span class="word" style="font-size: 19px">剪应力</span>
              </div>
            </div>
            <div class="qa-history">
              <h4>历史问答（最近 5 条）</h4>
              <article v-for="qa in qaHistory.slice(0, 5)" :key="qa.id" class="qa-item">
                <p><strong>Q:</strong> {{ qa.question }}</p>
                <p><strong>A:</strong> {{ qa.answer }}</p>
                <small>{{ qa.time }}</small>
              </article>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="📄 诊断报告" name="report">
          <div class="tab-content">
            <div class="report-section">
              <h4>学期学情诊断报告</h4>
              <div class="report-content">
                <p><strong>学习状态总评：</strong> 学习态度认真，掌握度良好。近期专注度提升 5%，在重点知识点上投入了更多时间，习题正确率趋于稳定。建议继续加强薄弱知识点的学习。</p>
                <p><strong>核心成果：</strong></p>
                <ul>
                  <li>综合掌握度：78%（较上周↑5%）</li>
                  <li>薄弱知识点数：3 个（正在逐步改善）</li>
                  <li>累计学习时长：32.5 小时</li>
                  <li>习题正确率：72%（较上周↑3%）</li>
                </ul>
                <p><strong>优势总结：</strong> 对基础知识掌握扎实，学习态度积极，互动频繁。</p>
                <p><strong>待改进方向：</strong> 建议加强 "应变" "弹性模量" 等重点难点的学习，增加相关配套练习的投入。</p>
              </div>
            </div>
            <div class="report-actions">
              <el-button type="primary" @click="exportPdfReport">导出 PDF 报告</el-button>
              <el-button @click="shareReport">分享报告</el-button>
              <el-button @click="syncReviewPlan">同步复习计划</el-button>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </section>

    <section class="action-guide">
      <h4>🎯 AI 推荐学习行动</h4>
      <div class="action-cards">
        <article class="action-item action-1">
          <div class="action-icon">1️⃣</div>
          <div>
            <strong>薄弱点专项学习：正应力与切应力</strong>
            <p>掌握度 45%，建议优先学习</p>
          </div>
          <el-button type="primary" size="small" @click="actionLearn">立即学习</el-button>
        </article>
        <article class="action-item action-2">
          <div class="action-icon">2️⃣</div>
          <div>
            <strong>错题重做：3 道关键习题</strong>
            <p>提升正确率到 75% 以上</p>
          </div>
          <el-button type="primary" size="small" @click="actionPractice">立即做题</el-button>
        </article>
        <article class="action-item action-3">
          <div class="action-icon">3️⃣</div>
          <div>
            <strong>生成本周复习计划</strong>
            <p>智能推荐，个性化复习方案</p>
          </div>
          <el-button type="primary" size="small" @click="generateReviewPlan">立即生成</el-button>
        </article>
      </div>
    </section>

    <div class="floating-quick-ops">
      <el-button type="primary" circle @click="aiDiagnose" title="AI 诊断">诊</el-button>
      <el-button circle @click="refreshData" title="刷新数据">刷</el-button>
      <el-button circle @click="exportReport" title="导出报告">导</el-button>
      <el-button circle @click="scrollToTop" title="回到顶部">顶</el-button>
    </div>
  </div>
</template>

<script setup>
/* eslint-disable no-undef */
import { computed, onMounted, ref, unref } from 'vue'
import { ElMessage } from 'element-plus'

const selectedCourse = ref('current')
const activeTab = ref('mastery')

defineProps({
  learningStats: {
    type: Object,
    default: () => ({})
  },
  weakPointTags: {
    type: Array,
    default: () => []
  },
  currentExplain: String,
  currentWeakPoint: String,
  currentTest: Object,
  testResult: Object
})

defineEmits(['start-weak-point', 'generate-test', 'check-answer'])

// 预制数据：章节与知识点
const chapters = ref([
  {
    id: 'ch1',
    name: '第1章 应力分析',
    points: [
      { id: 'p1-1', name: '正应力', mastery: 82 },
      { id: 'p1-2', name: '切应力', mastery: 45 },
      { id: 'p1-3', name: '应力转换', mastery: 67 },
      { id: 'p1-4', name: '主应力', mastery: 78 }
    ]
  },
  {
    id: 'ch2',
    name: '第2章 应变',
    points: [
      { id: 'p2-1', name: '线应变', mastery: 85 },
      { id: 'p2-2', name: '角应变', mastery: 52 },
      { id: 'p2-3', name: '应变张量', mastery: 71 },
      { id: 'p2-4', name: '主应变', mastery: 88 }
    ]
  },
  {
    id: 'ch3',
    name: '第3章 物理性质',
    points: [
      { id: 'p3-1', name: '弹性模量', mastery: 38 },
      { id: 'p3-2', name: '泊松比', mastery: 62 },
      { id: 'p3-3', name: '剪模量', mastery: 75 },
      { id: 'p3-4', name: '体积模量', mastery: 80 }
    ]
  }
])

const poiList = computed(() => {
  const all = []
  chapters.value.forEach((ch) => {
    ch.points.forEach((p) => {
      all.push({
        id: p.id,
        name: p.name,
        chapter: ch.name,
        mastery: p.mastery,
        qCorrectRate: Math.max(30, p.mastery - 10 + Math.random() * 20)
      })
    })
  })
  return all
})

// 统计数据
const masteredCount = computed(() => poiList.value.filter(p => p.mastery >= 80).length)
const basicCount = computed(() => poiList.value.filter(p => p.mastery >= 60 && p.mastery < 80).length)
const weakCount = computed(() => poiList.value.filter(p => p.mastery < 60).length)
const totalPoints = computed(() => poiList.value?.length || 0)
const avgMastery = computed(() => {
  const sum = poiList.value.reduce((acc, p) => acc + p.mastery, 0)
  return Math.round(sum / Math.max(1, totalPoints.value))
})

const totalPractices = 86
const avgCorrectRate = 72
const mistakeCount = 24
const retakeRate = 58
const totalQa = 12
const qaResolveRate = 92

const weaknessTop5 = ref([
  { id: 1, name: '正应力与切应力', cause: '习题正确率低', mastery: 45 },
  { id: 2, name: '弹性模量', cause: '多次重讲', mastery: 38 },
  { id: 3, name: '角应变', cause: '提问频繁', mastery: 52 },
  { id: 4, name: '应变张量', cause: '练习不足', mastery: 61 },
  { id: 5, name: '主应力概念', cause: '理解偏差', mastery: 58 }
])

const improvementPlan = ref([
  { id: 1, step: '1', title: '基础讲解', desc: '回归教材，重点理解正应力、切应力的定义与物理意义' },
  { id: 2, step: '2', title: '配套练习', desc: '生成 10 道基础应力计算题，逐步巩固' },
  { id: 3, step: '3', title: '专项测试', desc: '一周后进行小测试，检验学习效果' }
])

const mistakeList = ref([
  { id: 1, stem: '受力面积为 10cm²，作用力 100N，求正应力', poi: '正应力', cause: '单位换算错误', answer: '10 MPa' },
  { id: 2, stem: '材料直径 20mm，拉伸力 30kN，求应力', poi: '应力计算', cause: '公式理解有误', answer: '95.5 MPa' },
  { id: 3, stem: '应变 0.001，原长 200mm，求伸长量', poi: '应变定义', cause: '混淆了相对与绝对变形', answer: '0.2 mm' },
  { id: 4, stem: '弹性模量是什么？', poi: '弹性模量', cause: '概念记忆不清', answer: '应力与应变的比值' },
  { id: 5, stem: '剪应力的计算公式', poi: '切应力', cause: '记错了公式', answer: 'τ = F/A' }
])

const qaHistory = ref([
  { id: 1, question: '正应力和切应力的区别是什么？', answer: '正应力垂直于截面，切应力平行于截面。', time: '今天 14:30' },
  { id: 2, question: '如何计算应变？', answer: '应变 = 变形量 / 原始长度。', time: '今天 12:15' },
  { id: 3, question: '弹性模量是常数吗？', answer: '在弹性范围内，对于某种材料是常数。', time: '昨天 16:45' },
  { id: 4, question: '主应力如何求？', answer: '通过应力摩尔圆或特征方程求解。', time: '昨天 10:20' },
  { id: 5, question: '泊松比的范围？', answer: '一般在 0 到 0.5 之间。', time: '前天 15:00' }
])

const aiDiagnosisText = {
  mastery: '综合掌握度为 78%，相比上周提升 5%，学习进度稳定向前。其中 "基础概念" 掌握最好（85%+），但 "非线性计算" 和 "复杂应力状态分析" 仍需加强。建议重点关注这两个领域，配合针对性练习。',
  behavior: '最近一周学习时长 6.5 小时，平均每天 50 分钟，专注度 85%。你的学习习惯较好，主要集中在evening（晚上7-10点效率最高）。建议保持现有学习节奏，继续加强 "错题复盘" 和 "配套练习" 的投入比例。'
}

// eslint-disable-next-line vue/no-ref-as-operand
const overviewCards = computed(() => {
  const masterValue = unref(avgMastery)
  const correctValue = avgCorrectRate
  const weakValue = unref(weakCount)
  const qaValue = unref(totalQa)
  const rateValue = unref(qaResolveRate)
  return [
    {
      id: 1,
      emoji: '📊',
      label: '综合掌握度',
      number: `${masterValue}%`,
      sub: '所有知识点平均掌握程度',
      trend: 5,
      tabName: 'mastery'
    },
    {
      id: 2,
      emoji: '⚠️',
      label: '薄弱知识点',
      number: `${weakValue}个`,
      sub: '掌握度＜60% 的知识点',
      trend: -2,
      tabName: 'weakness'
    },
    {
      id: 3,
      emoji: '🎯',
      label: '学习专注度',
      number: '85%',
      sub: '有效学习时长占比',
      trend: 3,
      tabName: 'behavior'
    },
    {
      id: 4,
      emoji: '✏️',
      label: '习题正确率',
      number: `${correctValue}%`,
    sub: '所有练习平均正确率',
    trend: 2,
    tabName: 'mistake'
  },
  {
    id: 5,
    emoji: '💬',
    label: '互动提问数',
    number: `${qaValue}次`,
    sub: `解决率 ${rateValue}%`,
    trend: 1,
    tabName: 'qa'
  },
  {
    id: 6,
    emoji: '📅',
    label: '累计学习时长',
    number: '32.5h',
    sub: '本周新增 6.5 小时',
    trend: 8,
    tabName: 'behavior'
  }
  ]
})

const getMasteryColor = (mastery) => {
  if (mastery >= 80) return '#52c41a'
  if (mastery >= 60) return '#faad14'
  return '#f5222d'
}

const masteryStatus = (mastery) => {
  if (mastery >= 80) return 'success'
  if (mastery >= 60) return 'warning'
  return 'exception'
}

const showPoiDetail = (point) => {
  ElMessage.info(`${point.name}：掌握度 ${point.mastery}%，建议查看详情`)
}

const jumpToLearn = (point) => {
  ElMessage.success(`已推荐你学习 "${point.name}"`)
}

const askAbout = (item) => {
  ElMessage.success(`AI 讲解 "${item.name}"（演示模式）`)
}

const practiceProblem = (item) => {
  ElMessage.success(`生成 "${item.name}" 专项练习题（演示模式）`)
}

const addToTodo = () => {
  ElMessage.success('已添加到待办任务（演示模式）')
}

const redoMistake = (mistake) => {
  ElMessage.success(`已为你生成 "${mistake.stem}" 的重做版本（演示模式）`)
}

const addToReview = () => {
  ElMessage.success('已加入复习计划（演示模式）')
}

const generateReviewPlan = () => {
  ElMessage.success('本周复习计划已生成并同步到个人中心（演示模式）')
}

const optimizeLearning = () => {
  ElMessage.success('学习计划已优化（演示模式）')
}

const aiDiagnose = () => {
  ElMessage.success('AI 诊断已启动，请提出你的学习疑问（演示模式）')
}

const refreshData = () => {
  ElMessage.success('数据已刷新')
}

const exportReport = () => {
  ElMessage.success('诊断报告已导出为 PDF（演示模式）')
}

const exportPdfReport = () => {
  ElMessage.success('PDF 报告已导出（演示模式）')
}

const shareReport = () => {
  ElMessage.success('报告分享链接已生成（演示模式）')
}

const syncReviewPlan = () => {
  ElMessage.success('复习计划已同步到个人中心（演示模式）')
}

const actionLearn = () => {
  ElMessage.success('已跳转到知识点学习页面（演示模式）')
}

const actionPractice = () => {
  ElMessage.success('已跳转到错题重做页面（演示模式）')
}

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  ElMessage.info('学习分析数据已加载（前端模拟数据）')
})
</script>

<style scoped>
.analytics-workbench {
  background: linear-gradient(180deg, #ffffff 0%, #f6fbf8 100%);
  border-radius: 22px;
  border: 1px solid #d7e6de;
  padding: 16px;
  box-shadow: 0 18px 36px rgba(24, 55, 46, 0.08);
}

.head-area {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 14px;
}

.eyebrow {
  margin: 0;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: #67867a;
}

h3 {
  margin: 4px 0 0;
  color: #1f443d;
  font-size: 20px;
}

.head-tip {
  margin: 6px 0 0;
  color: #5f7a70;
  font-size: 13px;
}

.head-controls {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.overview-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 10px;
  margin-bottom: 12px;
}

.card {
  background: #fff;
  border: 1px solid #d8e7df;
  border-radius: 12px;
  padding: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.08);
  border-color: #2f605a;
}

.card.highlight {
  background: linear-gradient(135deg, #f0f9f6 0%, #e8f4f1 100%);
}

.card-head {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
}

.card-emoji {
  font-size: 20px;
}

.card-label {
  font-size: 12px;
  color: #5a7a6f;
}

.card-number {
  font-size: 22px;
  font-weight: bold;
  color: #1f443d;
  margin: 4px 0;
}

.card-sub {
  font-size: 11px;
  color: #6f857b;
  line-height: 1.4;
}

.card-trend {
  margin-top: 6px;
  font-size: 12px;
  font-weight: 600;
}

.card-trend.up {
  color: #1f8f57;
}

.card-trend.down {
  color: #cc4c43;
}

.diagnostic-area {
  background: #fff;
  border: 1px solid #d8e7df;
  border-radius: 14px;
  padding: 14px;
  margin-bottom: 12px;
}

.tab-content {
  padding: 8px 0;
}

.mastery-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 12px;
  margin-bottom: 12px;
}

.mastery-chart,
.mastery-pie {
  background: #f9fcfb;
  border: 1px solid #dce9e3;
  border-radius: 10px;
  padding: 10px;
}

.heatmap {
  display: grid;
  gap: 6px;
}

.heatmap-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.ch-name {
  min-width: 70px;
  font-size: 12px;
  color: #3e6860;
}

.heat-items {
  display: flex;
  gap: 4px;
}

.heat-item {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.heat-item:hover {
  transform: scale(1.2);
}

.pie-visual {
  display: grid;
  grid-template-columns: 1fr;
  gap: 8px;
}

.pie-item {
  text-align: center;
  padding: 8px;
  border-radius: 8px;
  border: 1px solid #dce9e3;
}

.pie-item.mastered {
  background: #f0f8f3;
  border-color: #1f8f57;
}

.pie-item.basic {
  background: #fffbf0;
  border-color: #faad14;
}

.pie-item.weak {
  background: #fef2f0;
  border-color: #cc4c43;
}

.pie-item strong {
  display: block;
  font-size: 18px;
}

.pie-item small {
  color: #667d73;
}

.poi-list,
.ai-diagnosis,
.weakness-top5,
.weakness-radar,
.weakness-plan,
.practice-overview,
.trend-chart,
.mistake-list,
.behavior-grid,
.ai-behavior-tip,
.qa-overview,
.qa-cloud,
.qa-history,
.report-section {
  background: #f9fcfb;
  border: 1px solid #dce9e3;
  border-radius: 10px;
  padding: 10px;
  margin-bottom: 10px;
}

.poi-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}

.poi-table thead {
  background: #eef6f2;
  border-bottom: 2px solid #d8e7df;
}

.poi-table th,
.poi-table td {
  padding: 8px;
  text-align: left;
}

.top5-cards {
  display: grid;
  gap: 8px;
}

.weakness-card {
  background: #fff;
  border: 1px solid #dce9e3;
  border-radius: 8px;
  padding: 10px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.rank-badge {
  min-width: 32px;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #faad14;
  color: #fff;
  border-radius: 50%;
  font-weight: bold;
  font-size: 14px;
}

.card-body {
  flex: 1;
}

.card-body strong {
  display: block;
}

.cause {
  margin: 4px 0 0;
  font-size: 11px;
  color: #6f857b;
}

.score {
  font-size: 16px;
  font-weight: bold;
  color: #1f443d;
}

.card-actions {
  display: flex;
  gap: 4px;
}

.radar-sim {
  display: grid;
  gap: 8px;
}

.radar-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.radar-item span {
  min-width: 80px;
  font-size: 12px;
}

.radar-bar {
  flex: 1;
  height: 8px;
  background: #dce9e3;
  border-radius: 4px;
  overflow: hidden;
}

.radar-bar::after {
  content: '';
  display: block;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, #1f8f57, #52c41a);
  border-radius: 4px;
}

.plan-step {
  display: flex;
  gap: 10px;
  margin-bottom: 8px;
  padding: 8px;
  border-radius: 8px;
  background: #fff;
}

.step-no {
  min-width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #e8f0ec;
  border-radius: 50%;
  font-weight: bold;
  color: #2f605a;
}

.stat-item {
  text-align: center;
  padding: 10px;
  border-radius: 8px;
  background: #fff;
  border: 1px solid #dce9e3;
}

.stat-num {
  display: block;
  font-size: 18px;
  font-weight: bold;
  color: #1f443d;
}

.stat-label {
  display: block;
  font-size: 11px;
  color: #6f857b;
  margin-top: 4px;
}

.mini-chart {
  display: flex;
  gap: 6px;
  align-items: flex-end;
  height: 120px;
}

.trend-bar {
  flex: 1;
  background: linear-gradient(180deg, #52c41a 0%, #1f8f57 100%);
  border-radius: 4px;
}

.mistake-item {
  background: #fff;
  border: 1px solid #dce9e3;
  border-radius: 8px;
  padding: 10px;
  margin-bottom: 8px;
}

.mistake-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.mistake-poi {
  font-size: 11px;
  background: #e8f0ec;
  color: #2f605a;
  padding: 2px 6px;
  border-radius: 4px;
}

.mistake-cause,
.mistake-answer {
  font-size: 12px;
  margin: 4px 0;
  line-height: 1.5;
}

.mistake-actions {
  margin-top: 6px;
  display: flex;
  gap: 6px;
}

.behavior-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 10px !important;
}

.behavior-chart,
.activity-dist {
  background: #fff;
  border: 1px solid #dce9e3;
  border-radius: 8px;
  padding: 10px;
}

.behavior-bars {
  display: flex;
  gap: 6px;
  align-items: flex-end;
  height: 120px;
}

.bar-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
}

.bar {
  flex: 1;
  width: 100%;
  background: #dce9e3;
  border-radius: 4px;
  max-width: 24px;
  background: linear-gradient(180deg, #52c41a 0%, #1f8f57 100%);
}

.activity-items {
  display: grid;
  gap: 8px;
}

.activity-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.activity-item span:first-child {
  min-width: 60px;
}

.activity-item .bar {
  flex: 1;
  height: 12px;
  left: 0;
  background: #dce9e3;
  border-radius: 6px;
  overflow: hidden;
}

.activity-item .fill {
  display: block;
  height: 100%;
  background: #1f8f57;
}

.word-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  justify-content: center;
  padding: 10px;
  background: #fff;
  border-radius: 8px;
}

.word {
  display: inline-block;
  padding: 4px 10px;
  background: linear-gradient(135deg, #e8f0ec 0%, #f0f8f3 100%);
  border: 1px solid #dce9e3;
  border-radius: 6px;
  cursor: pointer;
  color: #2f605a;
  transition: all 0.2s;
}

.word:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.qa-item {
  background: #fff;
  border: 1px solid #dce9e3;
  border-radius: 8px;
  padding: 10px;
  margin-bottom: 8px;
}

.qa-item p {
  margin: 4px 0;
  font-size: 12px;
  line-height: 1.5;
}

.qa-item small {
  color: #6f857b;
}

.report-content {
  background: #fff;
  border: 1px solid #dce9e3;
  border-radius: 8px;
  padding: 10px;
  font-size: 13px;
  line-height: 1.6;
}

.report-content p {
  margin: 6px 0;
}

.report-content ul {
  margin: 6px 0 6px 20px;
  padding: 0;
}

.report-content li {
  margin: 3px 0;
}

.report-actions {
  margin-top: 10px;
  display: flex;
  gap: 8px;
}

.action-guide {
  background: linear-gradient(135deg, #f0f9f6 0%, #e8f4f1 100%);
  border: 1px solid #d8e7df;
  border-radius: 12px;
  padding: 12px;
  margin-bottom: 12px;
}

.action-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 10px;
  margin-top: 10px;
}

.action-item {
  background: #fff;
  border: 1px solid #d8e7df;
  border-radius: 10px;
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.action-icon {
  font-size: 24px;
}

.action-item strong {
  display: block;
  font-size: 13px;
}

.action-item p {
  margin: 0;
  font-size: 11px;
  color: #5f7a70;
}

.floating-quick-ops {
  position: fixed;
  right: 24px;
  bottom: 24px;
  display: grid;
  gap: 10px;
  z-index: 20;
}

@media (max-width: 1100px) {
  .overview-cards {
    grid-template-columns: repeat(3, 1fr);
  }

  .mastery-grid {
    grid-template-columns: 1fr;
  }

  .behavior-grid {
    grid-template-columns: 1fr;
  }

  .action-cards {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 600px) {
  .head-area {
    flex-direction: column;
  }

  .overview-cards {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>