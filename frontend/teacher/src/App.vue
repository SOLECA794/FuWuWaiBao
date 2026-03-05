<template>
  <div class="teacher-app">
    <!-- 顶部导航栏 -->
    <div class="top-nav">
      <div class="nav-left">
        <span class="app-title">智能互动教学平台 · 教师端</span>
      </div>
      <div class="nav-right">
        <div class="teacher-info">
          <div class="avatar">
            <img src="https://picsum.photos/id/1005/40/40" alt="头像" />
          </div>
          <span class="teacher-name">教师 2025T001</span>
        </div>
      </div>
    </div>

    <!-- 主体内容 -->
    <div class="main-content">
      <!-- 左侧：课件管理区 -->
      <div class="courseware-manage-section">
        <div class="section-header">
          <h3>课件管理</h3>
          <div class="header-actions">
            <button @click="showPublishModal = true" class="publish-btn" :disabled="!currentCourseId">发布课件</button>
            <button @click="showUploadModal = true" class="upload-btn">+ 上传课件</button>
          </div>
        </div>

        <div class="courseware-list">
          <div
            v-for="course in coursewareList"
            :key="course.id"
            class="course-item"
            :class="{ active: course.id === currentCourseId }"
            @click="selectCourse(course)"
          >
            <span class="course-name">{{ course.name }}</span>
            <div class="course-actions">
              <span v-if="course.published" class="published-tag">已发布</span>
              <button @click.stop="deleteCourse(course.id)" class="del-btn">删除</button>
            </div>
          </div>
          <div v-if="coursewareList.length === 0" class="empty-tip">暂无课件，请点击上方按钮上传</div>
        </div>

        <div class="page-selector" v-if="currentCourseId">
          <h4>当前课件：{{ currentCourseName }}</h4>
          <div class="page-buttons">
            <button
              v-for="page in currentCourseTotalPages"
              :key="page"
              class="page-btn"
              :class="{ active: page === currentEditPage }"
              @click="selectEditPage(page)"
            >
              第{{ page }}页
            </button>
          </div>
        </div>
      </div>

      <!-- 右侧：讲稿编辑与分析区 -->
      <div class="editor-section">
        <div class="tabs">
          <button
            class="tab-btn"
            :class="{ active: activeTab === 'script' }"
            @click="activeTab = 'script'"
          >
            讲稿编辑
          </button>
          <button
            class="tab-btn"
            :class="{ active: activeTab === 'stats' }"
            @click="activeTab = 'stats'"
          >
            学情分析
          </button>
          <button
            class="tab-btn"
            :class="{ active: activeTab === 'questions' }"
            @click="activeTab = 'questions'"
          >
            提问统计
          </button>
          <button
            class="tab-btn"
            :class="{ active: activeTab === '卡点分析' }"
            @click="activeTab = '卡点分析'; renderChart()"
          >
            学习卡点可视化
          </button>
        </div>

        <!-- 讲稿编辑 -->
        <div class="tab-content" v-if="activeTab === 'script'">
          <div class="course-preview">
            <h4>第{{ currentEditPage }}页课件预览</h4>
            <img
              :src="`http://localhost:3000/api/courseware/${currentCourseId}/page/${currentEditPage}`"
              alt="课件预览"
              class="preview-img"
              v-if="currentCourseId"
            />
            <div class="no-preview" v-else>请先上传并选择课件</div>
          </div>

          <div class="script-editor">
            <div class="editor-actions">
              <button @click="generateAIScript" class="ai-btn">AI 生成讲稿</button>
              <button @click="saveScript" class="save-btn">保存讲稿</button>
            </div>
            <textarea
              v-model="currentScript"
              placeholder="请输入本页讲稿内容，支持AI生成..."
              class="script-textarea"
            ></textarea>
          </div>
        </div>

        <!-- 学情分析 -->
        <div class="tab-content" v-if="activeTab === 'stats'">
          <div class="stats-header" v-if="currentCourseId">
            <h4>学情分析 - {{ currentCourseName }}</h4>
          </div>
          <div class="stats-grid" v-if="currentCourseId">
            <div class="stat-card">
              <div class="stat-value">{{ studentStats.totalQuestions }}</div>
              <div class="stat-label">总提问数</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">{{ studentStats.hotPages.join('、') || '暂无' }}</div>
              <div class="stat-label">高频提问页码</div>
            </div>
            <div class="stat-card" style="grid-column: span 2">
              <div class="stat-value">{{ studentStats.keyDifficulties }}</div>
              <div class="stat-label">重点难点</div>
            </div>
          </div>
          <div v-else class="empty-tip">请先选择一个课件查看学情数据</div>
        </div>

        <!-- 提问统计 -->
        <div class="tab-content" v-if="activeTab === 'questions'">
          <div class="questions-header" v-if="currentCourseId">
            <h4>提问统计 - {{ currentCourseName }}</h4>
            <div class="filter-bar">
              <span>按页码筛选：</span>
              <select v-model="filterPage" class="page-select">
                <option value="">全部</option>
                <option v-for="page in currentCourseTotalPages" :key="page" :value="page">第{{ page }}页</option>
              </select>
            </div>
          </div>
          <div class="questions-list" v-if="currentCourseId">
            <div
              v-for="q in filteredQuestions"
              :key="q.id"
              class="question-item"
            >
              <div class="question-meta">
                <span class="student-id">学生 {{ q.studentId }}</span>
                <span class="page-tag">第{{ q.page }}页</span>
                <span class="time">{{ q.time }}</span>
              </div>
              <div class="question-content">{{ q.content }}</div>
              <div class="answer-content" v-if="q.answer">
                <span class="answer-label">AI 回复：</span>{{ q.answer }}
              </div>
            </div>
            <div v-if="filteredQuestions.length === 0" class="empty-tip">暂无提问记录</div>
          </div>
          <div v-else class="empty-tip">请先选择一个课件查看提问统计</div>
        </div>

        <!-- 学习卡点可视化 -->
        <div class="tab-content" v-if="activeTab === '卡点分析'">
          <div class="chart-header" v-if="currentCourseId">
            <h4>学习卡点分析 - {{ currentCourseName }}</h4>
            <div class="chart-type">
              <span>图表类型：</span>
              <button 
                class="chart-btn" 
                :class="{ active: chartType === 'bar' }"
                @click="chartType = 'bar'; renderChart()"
              >柱状图</button>
              <button 
                class="chart-btn" 
                :class="{ active: chartType === 'line' }"
                @click="chartType = 'line'; renderChart()"
              >折线图</button>
              <button 
                class="chart-btn" 
                :class="{ active: chartType === 'pie' }"
                @click="chartType = 'pie'; renderChart()"
              >饼图</button>
            </div>
          </div>
          <div v-if="currentCourseId" class="chart-container">
            <div id="卡点图表" class="chart" style="width:100%;height:400px;"></div>
            <div class="chart-tip">
              <p>数据说明：</p>
              <ul>
                <li>提问量：该页面学生发起的提问总数</li>
                <li>停留时长：学生平均停留时长（秒）</li>
                <li>卡点指数：综合提问量+停留时长计算的卡点程度（0-10）</li>
              </ul>
            </div>
          </div>
          <div v-else class="empty-tip">请先选择一个课件查看卡点分析</div>
        </div>
      </div>
    </div>

    <!-- 课件上传弹窗 -->
    <div class="modal-overlay" v-if="showUploadModal" @click="showUploadModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>上传课件（PPT/PDF）</h3>
          <button @click="showUploadModal = false" class="close-btn">×</button>
        </div>
        <div class="upload-form">
          <input
            type="file"
            ref="fileInput"
            accept=".ppt,.pptx,.pdf"
            @change="handleFileSelect"
            class="file-input"
          />
          <div class="file-name" v-if="selectedFileName">{{ selectedFileName }}</div>
          <button @click="uploadCourseware" class="upload-submit" :disabled="!selectedFileName">
            {{ uploadLoading ? '上传中...' : '上传并解析' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 发布课件弹窗 -->
    <div class="modal-overlay" v-if="showPublishModal" @click="showPublishModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>发布课件给学生端</h3>
          <button @click="showPublishModal = false" class="close-btn">×</button>
        </div>
        <div class="publish-form">
          <div class="form-item">
            <label>当前课件：</label>
            <span>{{ currentCourseName }}</span>
          </div>
          <div class="form-item">
            <label>发布范围：</label>
            <select v-model="publishScope" class="scope-select">
              <option value="all">全部学生</option>
              <option value="class1">班级1</option>
              <option value="class2">班级2</option>
            </select>
          </div>
          <div class="form-actions">
            <button @click="publishCourseware" class="confirm-btn">确认发布</button>
            <button @click="showPublishModal = false" class="cancel-btn">取消</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue';
import * as echarts from 'echarts';

// ========== 核心变量 ==========
const coursewareList = ref([
  { id: '1', name: 'Python数据分析·第3课：缺失值处理', totalPages: 10, published: true },
  { id: '2', name: 'Python数据分析·第2课：数据清洗', totalPages: 8, published: false }
]);
const currentCourseId = ref('1');
const currentCourseName = ref('Python数据分析·第3课：缺失值处理');
const currentCourseTotalPages = ref(10);
const currentEditPage = ref(1);
const currentScript = ref('这是第1页的讲稿内容...');

const showUploadModal = ref(false);
const fileInput = ref(null);
const selectedFileName = ref('');
const uploadLoading = ref(false);

const showPublishModal = ref(false);
const publishScope = ref('all');

const activeTab = ref('script');
const studentStats = ref({
  totalQuestions: 12,
  hotPages: [4, 6, 8],
  keyDifficulties: '缺失值填充、异常值识别'
});

// 提问统计
const questionRecords = ref([
  { id: 1, studentId: '2025001', page: 4, content: '这一页的缺失值填充方法有哪些？', answer: '常用方法有均值填充、中位数填充、KNN填充等。', time: '2026-03-01 20:15' },
  { id: 2, studentId: '2025002', page: 4, content: '为什么要处理缺失值？不处理会有什么影响？', answer: '缺失值会导致统计偏差，影响模型训练效果。', time: '2026-03-01 20:20' },
  { id: 3, studentId: '2025001', page: 6, content: '异常值识别的常用方法是什么？', answer: '常用方法有3σ原则、箱线图法、Z-score等。', time: '2026-03-01 20:30' }
]);
const filterPage = ref('');
const filteredQuestions = computed(() => {
  if (!filterPage.value) return questionRecords.value;
  return questionRecords.value.filter(q => q.page === Number(filterPage.value));
});

// 学习卡点可视化
const chartType = ref('bar');
let chartInstance = null;
const cardData = ref([
  { page: 1, 提问量: 1, 停留时长: 20, 卡点指数: 1.2 },
  { page: 2, 提问量: 0, 停留时长: 15, 卡点指数: 0.8 },
  { page: 3, 提问量: 2, 停留时长: 30, 卡点指数: 2.5 },
  { page: 4, 提问量: 8, 停留时长: 90, 卡点指数: 8.5 },
  { page: 5, 提问量: 3, 停留时长: 45, 卡点指数: 3.8 },
  { page: 6, 提问量: 7, 停留时长: 80, 卡点指数: 7.2 },
  { page: 7, 提问量: 2, 停留时长: 25, 卡点指数: 2.1 },
  { page: 8, 提问量: 6, 停留时长: 75, 卡点指数: 6.8 },
  { page: 9, 提问量: 1, 停留时长: 18, 卡点指数: 1.0 },
  { page: 10, 提问量: 0, 停留时长: 12, 卡点指数: 0.5 }
]);

// ========== 初始化 ==========
onMounted(async () => {
  await loadCoursewareList();
  await loadStudentStats(currentCourseId.value);
  await loadCardData(currentCourseId.value);
});

onUnmounted(() => {
  if (chartInstance) {
    chartInstance.dispose();
    chartInstance = null;
  }
});

// ========== 课件管理 ==========
const loadCoursewareList = async () => {
  try {
    const res = await fetch('http://localhost:3000/api/teacher/courseware-list');
    const data = await res.json();
    coursewareList.value = data.list || coursewareList.value;
  } catch (err) {
    console.error('加载课件列表失败', err);
  }
};

const selectCourse = async (course) => {
  currentCourseId.value = course.id;
  currentCourseName.value = course.name;
  currentCourseTotalPages.value = course.totalPages;
  currentEditPage.value = 1;
  await loadScript(course.id, 1);
  await loadStudentStats(course.id);
  await loadCardData(course.id);
};

const deleteCourse = async (courseId) => {
  if (!confirm('确定删除该课件吗？')) return;
  try {
    await fetch(`http://localhost:3000/api/teacher/courseware/${courseId}`, { method: 'DELETE' });
    await loadCoursewareList();
    if (currentCourseId.value === courseId) {
      currentCourseId.value = '';
      currentCourseName.value = '';
      currentCourseTotalPages.value = 0;
    }
  } catch (err) {
    alert('删除课件失败：' + err.message);
  }
};

const selectEditPage = async (page) => {
  currentEditPage.value = page;
  await loadScript(currentCourseId.value, page);
};

// ========== 讲稿相关 ==========
const loadScript = async (courseId, page) => {
  try {
    const res = await fetch(`http://localhost:3000/api/teacher/script/${courseId}/${page}`);
    const data = await res.json();
    currentScript.value = data.content || currentScript.value;
  } catch (err) {
    currentScript.value = '';
  }
};

const saveScript = async () => {
  if (!currentScript.value.trim()) {
    alert('请输入讲稿内容！');
    return;
  }
  try {
    await fetch('http://localhost:3000/api/teacher/script/save', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        courseId: currentCourseId.value,
        page: currentEditPage.value,
        content: currentScript.value
      })
    });
    alert('讲稿保存成功！');
  } catch (err) {
    alert('保存讲稿失败：' + err.message);
  }
};

const generateAIScript = async () => {
  try {
    currentScript.value = 'AI正在生成讲稿...';
    const res = await fetch('http://localhost:3000/api/teacher/ai-generate-script', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        courseId: currentCourseId.value,
        page: currentEditPage.value,
        courseName: currentCourseName.value
      })
    });
    const data = await res.json();
    currentScript.value = data.content || 'AI生成失败，请重试';
  } catch (err) {
    currentScript.value = '生成失败：' + err.message;
  }
};

// ========== 上传课件 ==========
const handleFileSelect = (e) => {
  const file = e.target.files[0];
  if (file) {
    selectedFileName.value = file.name;
  }
};

const uploadCourseware = async () => {
  const file = fileInput.value.files[0];
  if (!file) return;

  uploadLoading.value = true;
  const formData = new FormData();
  formData.append('courseware', file);

  try {
    const res = await fetch('http://localhost:3000/api/teacher/upload-courseware', {
      method: 'POST',
      body: formData
    });
    const data = await res.json();
    alert('课件上传并解析成功！');
    showUploadModal.value = false;
    selectedFileName.value = '';
    await loadCoursewareList();
  } catch (err) {
    alert('上传失败：' + err.message);
  } finally {
    uploadLoading.value = false;
    fileInput.value.value = '';
  }
};

// ========== 发布课件 ==========
const publishCourseware = async () => {
  try {
    await fetch('http://localhost:3000/api/teacher/publish-courseware', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        courseId: currentCourseId.value,
        scope: publishScope.value
      })
    });
    const course = coursewareList.value.find(c => c.id === currentCourseId.value);
    if (course) course.published = true;
    alert('课件发布成功！学生端已可查看。');
    showPublishModal.value = false;
  } catch (err) {
    alert('发布失败：' + err.message);
  }
};

// ========== 学情分析 ==========
const loadStudentStats = async (courseId) => {
  try {
    const res = await fetch(`http://localhost:3000/api/teacher/student-stats/${courseId}`);
    const data = await res.json();
    studentStats.value = {
      totalQuestions: data.totalQuestions || 12,
      hotPages: data.hotPages || [4, 6, 8],
      keyDifficulties: data.keyDifficulties || '缺失值填充、异常值识别'
    };
  } catch (err) {
    studentStats.value = {
      totalQuestions: 0,
      hotPages: [],
      keyDifficulties: '加载失败'
    };
  }
};

// ========== 学习卡点可视化 ==========
const loadCardData = async (courseId) => {
  try {
    const res = await fetch(`http://localhost:3000/api/teacher/card-data/${courseId}`);
    const data = await res.json();
    cardData.value = data.list || cardData.value;
  } catch (err) {
    console.error('加载卡点数据失败，使用模拟数据', err);
  }
};

const renderChart = () => {
  if (chartInstance) {
    chartInstance.dispose();
  }

  chartInstance = echarts.init(document.getElementById('卡点图表'));

  const pages = cardData.value.map(item => `第${item.page}页`);
  const questionCounts = cardData.value.map(item => item.提问量);
  const stayTimes = cardData.value.map(item => item.停留时长);
  const cardScores = cardData.value.map(item => item.卡点指数);

  let option = {};

  switch (chartType.value) {
    case 'bar':
      option = {
        title: { text: '各页面学习卡点数据' },
        tooltip: { trigger: 'axis' },
        legend: { data: ['提问量', '停留时长(秒)', '卡点指数'] },
        xAxis: { type: 'category', data: pages },
        yAxis: { type: 'value' },
        series: [
          { name: '提问量', type: 'bar', data: questionCounts },
          { name: '停留时长(秒)', type: 'bar', data: stayTimes },
          { name: '卡点指数', type: 'bar', data: cardScores, itemStyle: { color: '#ff4d4f' } }
        ]
      };
      break;
    case 'line':
      option = {
        title: { text: '各页面学习卡点趋势' },
        tooltip: { trigger: 'axis' },
        legend: { data: ['提问量', '停留时长(秒)', '卡点指数'] },
        xAxis: { type: 'category', data: pages },
        yAxis: { type: 'value' },
        series: [
          { name: '提问量', type: 'line', data: questionCounts },
          { name: '停留时长(秒)', type: 'line', data: stayTimes },
          { name: '卡点指数', type: 'line', data: cardScores, lineStyle: { color: '#ff4d4f' }, itemStyle: { color: '#ff4d4f' } }
        ]
      };
      break;
    case 'pie':
      const top5CardData = [...cardData.value].sort((a, b) => b.卡点指数 - a.卡点指数).slice(0, 5);
      option = {
        title: { text: 'TOP5 卡点页面占比' },
        tooltip: { trigger: 'item' },
        legend: { orient: 'vertical', left: 'left', data: top5CardData.map(item => `第${item.page}页`) },
        series: [
          {
            name: '卡点指数',
            type: 'pie',
            radius: ['40%', '70%'],
            data: top5CardData.map(item => ({
              name: `第${item.page}页`,
              value: item.卡点指数
            })),
            label: {
              formatter: '{b}: {c} ({d}%)'
            }
          }
        ]
      };
      break;
  }

  chartInstance.setOption(option);
  window.addEventListener('resize', () => {
    if (chartInstance) {
      chartInstance.resize();
    }
  });
};
</script>

<style scoped>
/* 全局 */
.teacher-app {
  width: 100%;
  height: 100vh;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background: #f5f7fa;
}

/* 顶部导航 */
.top-nav {
  height: 56px;
  background: #1677ff;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  color: white;
}
.app-title {
  font-size: 16px;
  font-weight: 500;
}
.teacher-info {
  display: flex;
  align-items: center;
  gap: 10px;
}
.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
}
.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.teacher-name {
  font-size: 14px;
}

/* 主体内容 */
.main-content {
  display: flex;
  height: calc(100vh - 56px);
}

/* 左侧课件管理区 */
.courseware-manage-section {
  width: 320px;
  background: white;
  border-right: 1px solid #e8e8e8;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.section-header h3 {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin: 0;
}
.header-actions {
  display: flex;
  gap: 8px;
}
.upload-btn {
  padding: 6px 12px;
  background: #1677ff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
}
.publish-btn {
  padding: 6px 12px;
  background: #52c41a;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
}
.publish-btn:disabled {
  background: #999;
  cursor: not-allowed;
}
.courseware-list {
  flex: 1;
  overflow-y: auto;
}
.course-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: #fafafa;
  border-radius: 6px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.2s;
}
.course-item:hover {
  background: #f0f7ff;
}
.course-item.active {
  background: #e6f7ff;
  border-left: 3px solid #1677ff;
}
.course-name {
  font-size: 14px;
  color: #333;
}
.course-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}
.published-tag {
  font-size: 12px;
  color: #52c41a;
  background: #f6ffed;
  padding: 2px 6px;
  border-radius: 4px;
}
.del-btn {
  padding: 4px 8px;
  background: #ff4d4f;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}
.page-selector h4 {
  font-size: 14px;
  color: #666;
  margin-bottom: 12px;
}
.page-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.page-btn {
  padding: 6px 10px;
  border: 1px solid #d9d9d9;
  background: white;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}
.page-btn.active {
  background: #1677ff;
  color: white;
  border-color: #1677ff;
}
.empty-tip {
  text-align: center;
  color: #999;
  font-size: 14px;
  padding: 20px;
}

/* 右侧编辑与分析区 */
.editor-section {
  flex: 1;
  background: white;
  display: flex;
  flex-direction: column;
}
.tabs {
  display: flex;
  border-bottom: 1px solid #e8e8e8;
}
.tab-btn {
  flex: 1;
  padding: 14px 0;
  background: white;
  border: none;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  border-bottom: 2px solid transparent;
}
.tab-btn.active {
  color: #1677ff;
  border-bottom-color: #1677ff;
}
.tab-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

/* 讲稿编辑 */
.course-preview {
  margin-bottom: 20px;
}
.course-preview h4 {
  font-size: 14px;
  color: #333;
  margin-bottom: 12px;
}
.preview-img {
  max-width: 100%;
  max-height: 300px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}
.no-preview {
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  border: 1px dashed #d9d9d9;
  border-radius: 8px;
  background: #fafafa;
}
.script-editor {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.editor-actions {
  display: flex;
  gap: 12px;
}
.ai-btn {
  padding: 8px 16px;
  background: #52c41a;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.save-btn {
  padding: 8px 16px;
  background: #1677ff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.script-textarea {
  width: 100%;
  height: 300px;
  padding: 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  resize: vertical;
  font-size: 14px;
  line-height: 1.6;
}

/* 学情分析 */
.stats-header h4, .questions-header h4, .chart-header h4 {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin-bottom: 20px;
}
.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}
.stat-card {
  background: #fafafa;
  padding: 20px;
  border-radius: 8px;
  text-align: center;
}
.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: #1677ff;
  margin-bottom: 8px;
}
.stat-label {
  font-size: 14px;
  color: #666;
}

/* 提问统计 */
.filter-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  font-size: 14px;
  color: #666;
}
.page-select {
  padding: 6px 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 14px;
}
.questions-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.question-item {
  background: #fafafa;
  padding: 16px;
  border-radius: 8px;
}
.question-meta {
  display: flex;
  gap: 12px;
  margin-bottom: 8px;
  font-size: 12px;
  color: #666;
}
.page-tag {
  color: #1677ff;
  background: #e6f7ff;
  padding: 2px 6px;
  border-radius: 4px;
}
.question-content {
  font-size: 14px;
  color: #333;
  margin-bottom: 8px;
}
.answer-content {
  font-size: 14px;
  color: #666;
}
.answer-label {
  color: #1677ff;
  font-weight: 500;
}

/* 学习卡点可视化 */
.chart-type {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  font-size: 14px;
  color: #666;
}
.chart-btn {
  padding: 6px 12px;
  border: 1px solid #d9d9d9;
  background: white;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
}
.chart-btn.active {
  background: #1677ff;
  color: white;
  border-color: #1677ff;
}
.chart-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}
.chart {
  background: #fafafa;
  padding: 16px;
  border-radius: 8px;
}
.chart-tip {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
}
.chart-tip ul {
  margin: 8px 0 0 20px;
  padding: 0;
}

/* 上传/发布弹窗 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 100;
}
.modal-content {
  width: 500px;
  background: white;
  border-radius: 8px;
  padding: 24px;
}
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e8e8e8;
}
.modal-header h3 {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin: 0;
}
.close-btn {
  background: transparent;
  border: none;
  font-size: 20px;
  color: #666;
  cursor: pointer;
}
.upload-form, .publish-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.file-input, .scope-select {
  padding: 10px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  background: #fafafa;
}
.file-name {
  padding: 10px;
  background: #fafafa;
  border-radius: 4px;
  color: #333;
  font-size: 14px;
}
.upload-submit {
  padding: 10px;
  background: #1677ff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.upload-submit:disabled {
  background: #999;
  cursor: not-allowed;
}
.form-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #333;
}
.form-item label {
  width: 80px;
  font-weight: 500;
}
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
}
.confirm-btn {
  padding: 8px 16px;
  background: #1677ff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.cancel-btn {
  padding: 8px 16px;
  background: white;
  color: #666;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  cursor: pointer;
}
</style>