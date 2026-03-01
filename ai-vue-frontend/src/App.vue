<template>
  <!-- 1. 整体布局：顶部导航 + 课件+聊天分栏 + 思维导图弹窗 -->
  <div class="student-app">
    <!-- 顶部导航栏 -->
    <div class="top-nav">
      <h2>智能互动讲课系统 - 学生端</h2>
      <div class="top-actions">
        <input type="file" @change="handleFileUpload" ref="fileInput" style="display: none" accept=".pdf,.pptx" />
        <button @click="$refs.fileInput.click()" class="upload-btn">📁 上传课件</button>
        <button @click="fetchMindMap" class="mindmap-btn">📊 思维导图</button>
      </div>
    </div>

    <!-- 主体：课件区 + 聊天区 -->
    <div class="main-layout">
      <!-- 左侧：课件播放器（翻页/播放暂停） -->
      <div class="courseware-panel">
        <h3>课件展示区</h3>
        <!-- 课件讲稿预览 -->
        <div class="page-display">
          <div class="page-preview">
            <p class="page-preview-label">第{{ currentPage }} 页讲稿</p>
            <p v-if="scriptLoading" class="page-preview-loading">正在加载讲稿…</p>
            <p v-else-if="currentScript" class="page-preview-text">{{ currentScript }}</p>
            <p v-else class="page-preview-text">该页目前还没有生成讲稿。</p>
          </div>
          <div class="page-preview-meta">
            <label class="page-preview-label" for="pageRange">页面跳转</label>
            <input id="pageRange" class="page-range" type="range" :min="1" :max="Math.max(totalPages, 1)" :value="currentPage" @input="goToPage(Number($event.target.value))" />
            <div class="page-buttons">
              <button v-for="page in pageRange" :key="page" :class="['page-btn', { active: page === currentPage }]" @click="goToPage(page)">
                {{ page }}
              </button>
            </div>
          </div>
        </div>
        <!-- 课件控制按钮 -->
        <div class="courseware-controls">
          <button @click="prevPage" :disabled="currentPage === 1" class="ctrl-btn">上一页</button>
          <span class="page-info">第 {{ currentPage }} / {{ totalPages }} 页</span>
          <button @click="nextPage" :disabled="currentPage === totalPages" class="ctrl-btn">下一页</button>
          <button @click="togglePlay" class="ctrl-btn play-btn">{{ isPlaying ? '暂停播放' : '开始播放' }}</button>
        </div>
      </div>

      <!-- 右侧：聊天区（带页码、AI交互） -->
      <div class="chat-panel">
        <h3>智能问答（自动携带当前页码）</h3>
        <!-- 聊天记录 -->
        <div class="message-list">
          <div 
            v-for="msg in chatMessages" 
            :key="msg.id"
            class="message"
            :class="{ 'user-msg': msg.isUser }"
          >
            <!-- 页码标签 -->
            <div class="page-tag" v-if="msg.pageTag">第{{ msg.pageTag }}页</div>
            <!-- 消息内容 -->
            <div class="msg-content">{{ msg.content }}</div>
          </div>
        </div>
        <!-- 输入区 -->
        <div class="chat-input">
          <textarea 
            v-model="inputText" 
            placeholder="输入你的问题（比如：这一页的公式怎么理解？）"
            class="input-area"
            @keydown.enter.prevent="sendMessage"
          ></textarea>
          <div class="input-controls">
            <button @click="sendMessage" class="send-btn" :disabled="!inputText.trim()">发送</button>
            <span class="page-tip">当前关联页码：{{ currentPage }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 思维导图弹窗（点击顶部按钮显示） -->
    <div class="modal-overlay" v-if="showMindMap" @click="showMindMap = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>课程思维导图</h3>
          <button @click="showMindMap = false" class="close-btn">×</button>
        </div>
        <div class="mindmap-content">
          <div v-if="mindmapNodes.length" class="mindmap-tree">
            <div
              v-for="(node, idx) in mindmapNodes"
              :key="idx"
              class="mindmap-node"
              :style="{ paddingLeft: `${node.level * 18}px` }"
            >
              <span class="mindmap-bullet">•</span>
              <span>{{ node.text }}</span>
            </div>
          </div>
          <div v-else class="mindmap-empty">当前页还没有可展示的思维导图，请先上传并生成讲稿。</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 2. 功能逻辑代码（所有交互都在这里）
import { computed, onMounted, ref, watch } from 'vue';

// ========== ① 课件相关变量 ==========
const currentDocId = ref(localStorage.getItem('last_doc_id') || ''); 
const currentPage = ref(1); 
const totalPages = ref(1); 
const currentScript = ref('等待课件加载...'); 
const currentMindMap = ref(null); 
const isPlaying = ref(false); 
const showMindMap = ref(false); 
const scriptLoading = ref(false);
const lessonsCache = ref({});

// ========== ② 聊天相关变量 (初始化逻辑) ==========
const inputText = ref('');
// 从 localStorage 加载历史消息，如果没有则显示欢迎语
const savedMessages = localStorage.getItem('chat_history');
const chatMessages = ref(savedMessages ? JSON.parse(savedMessages) : [
  {
    id: 1,
    content: '你好！我是你的智能助教。请先上传 PDF 或 PPTX 课件，我会为你生成详细的讲义和互动内容。',
    isUser: false
  }
]);

// 监听消息变化并持久化
watch(chatMessages, (newMsgs) => {
  localStorage.setItem('chat_history', JSON.stringify(newMsgs));
}, { deep: true });

onMounted(() => {
  // 恢复之前解析过的页数
  const savedTotalPages = localStorage.getItem('last_total_pages');
  if (savedTotalPages) {
    totalPages.value = parseInt(savedTotalPages);
  }

  if (currentDocId.value) {
    fetchLessonContent();
  }
});

const pageRange = computed(() => {
  const pages = [];
  const visible = Math.min(totalPages.value, 7);
  let start = Math.max(1, currentPage.value - Math.floor(visible / 2));
  if (start + visible - 1 > totalPages.value) {
    start = Math.max(1, totalPages.value - visible + 1);
  }
  for (let i = 0; i < visible; i += 1) {
    pages.push(start + i);
  }
  return pages;
});

const mindmapNodes = computed(() => {
  const markdown = currentMindMap.value || '';
  const lines = markdown
    .split('\n')
    .map((line) => line.replace(/\t/g, '    '))
    .map((line) => line.trimEnd())
    .filter((line) => line.trim().length > 0);

  const nodes = [];
  lines.forEach((line) => {
    const bulletMatch = line.match(/^(\s*)([-*+]\s+)(.+)$/);
    if (bulletMatch) {
      const indentSpaces = bulletMatch[1].length;
      nodes.push({
        level: Math.floor(indentSpaces / 2),
        text: bulletMatch[3].trim(),
      });
      return;
    }

    const headingMatch = line.match(/^(#+)\s+(.+)$/);
    if (headingMatch) {
      nodes.push({
        level: Math.max(0, headingMatch[1].length - 1),
        text: headingMatch[2].trim(),
      });
      return;
    }

    nodes.push({ level: 0, text: line.trim() });
  });

  return nodes;
});

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value && page !== currentPage.value) {
    currentPage.value = page;
  }
};

watch(currentPage, () => {
  fetchLessonContent();
});

// ========== ④ 课件控制函数 ==========
// 上一页
const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
};
// 下一页
const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
  }
};
// 播放/暂停切换
const togglePlay = () => {
  isPlaying.value = !isPlaying.value;
  if (isPlaying.value) {
    console.log('开始播放课件，当前页：', currentPage.value);
  } else {
    console.log('暂停播放课件');
  }
};

// ========== ⑤ 核心交互函数：文件上传、聊天、获取详情 ==========

// 1. 上传文件逻辑 (一次上传，持久使用)
const handleFileUpload = async (event) => {
  const file = event.target.files[0];
  if (!file) return;

  const formData = new FormData();
  formData.append('file', file);

  try {
    chatMessages.value.push({ id: Date.now(), content: `正在上传并解析: ${file.name}...`, isUser: false });
    
    const res = await fetch('http://localhost:8000/upload', {
      method: 'POST',
      body: formData
    });
    
    const data = await res.json();
    currentDocId.value = data.doc_id;
    totalPages.value = data.total_pages;
    
    // 持久化 ID 和总页数，供刷新后恢复
    localStorage.setItem('last_doc_id', data.doc_id);
    localStorage.setItem('last_total_pages', data.total_pages.toString());
    
    chatMessages.value.push({ id: Date.now() + 1, content: `文件解析成功！共 ${data.total_pages} 页。你可以开始提阅了。`, isUser: false });
    fetchLessonContent();
  } catch (error) {
    alert('上传失败: ' + error.message);
  }
};

// 2. 获取当前页的讲稿和导图
const fetchLessonContent = async () => {
  if (!currentDocId.value) return;

  const cachedLessons = lessonsCache.value[currentDocId.value];
  if (cachedLessons) {
    const cachedPageData = cachedLessons.find((item) => item.page === currentPage.value);
    if (cachedPageData) {
      currentScript.value = cachedPageData.script || '';
      currentMindMap.value = cachedPageData.mindmap_markdown || '';
    }
    return;
  }

  scriptLoading.value = true;
  try {
    const res = await fetch(`http://localhost:8000/lessons/${currentDocId.value}?mode=llm`);
    const data = await res.json();
    const lessons = Array.isArray(data.lessons) ? data.lessons : [];
    lessonsCache.value[currentDocId.value] = lessons;
    const pageData = lessons.find(l => l.page === currentPage.value);
    if (pageData) {
      currentScript.value = pageData.script;
      currentMindMap.value = pageData.mindmap_markdown;
    }
  } catch (e) {
    console.error('获取讲稿失败', e);
  } finally {
    scriptLoading.value = false;
  }
};

// 3. 发送聊天消息
const sendMessage = async () => {
  if (!inputText.value.trim()) return;

  const userMsg = {
    id: Date.now(),
    content: inputText.value.trim(),
    pageTag: currentPage.value,
    isUser: true
  };
  chatMessages.value.push(userMsg);
  inputText.value = '';

  try {
    chatMessages.value.push({ id: Date.now() + 1, content: 'AI正在思考中...', isUser: false });

    // 修改为你的 Python 后端地址
    const response = await fetch('http://localhost:8000/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        doc_id: currentDocId.value, // 使用保存的 ID
        question: userMsg.content,
        page: currentPage.value,
        mode: 'llm' // 强制 LLM 模式
      })
    });

    const data = await response.json();
    chatMessages.value.pop(); // 删除加载中状态
    
    chatMessages.value.push({
      id: Date.now() + 2,
      content: data.answer || '抱歉，我没有理解这个问题。',
      pageTag: data.source_page,
      isUser: false
    });

    // 如果触发了建议，可以在下方动态显示（可选实现）
    if (data.follow_up_suggestion) {
      chatMessages.value.push({
        id: Date.now() + 3,
        content: `💡 建议：${data.follow_up_suggestion}`,
        isUser: false,
        isSuggestion: true
      });
    }

  } catch (error) {
    chatMessages.value.pop();
    chatMessages.value.push({ id: Date.now() + 2, content: `错误: ${error.message}`, isUser: false });
  }
};

// 获取思维导图
const fetchMindMap = async () => {
  showMindMap.value = true;
  if (!currentMindMap.value && currentDocId.value) {
    await fetchLessonContent();
  }
};
</script>

<style scoped>
/* 3. 样式代码（直接用，不用改） */
.student-app {
  width: 100%;
  height: 100vh;
  overflow: hidden;
  font-family: Arial, sans-serif;
}
/* 顶部导航 */
.top-nav {
  height: 60px;
  line-height: 60px;
  padding: 0 20px;
  background: #4285f4;
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.mindmap-btn {
  background: white;
  color: #4285f4;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
}
/* 主体布局 */
.main-layout {
  display: flex;
  flex-wrap: wrap;
  height: calc(100vh - 60px);
}
/* 课件区 */
.courseware-panel {
  flex: 2;
  min-width: 320px;
  padding: 20px;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.page-display {
  width: 100%;
  min-height: 260px;
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 8px 16px rgba(0,0,0,0.08);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.page-preview {
  flex: 1;
}
.page-preview-label {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}
.page-preview-loading {
  color: #777;
  font-style: italic;
}
.page-preview-text {
  color: #444;
  line-height: 1.6;
  white-space: pre-wrap;
  max-height: 260px;
  overflow: hidden;
}
.page-preview-meta {
  margin-top: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.page-range {
  width: 100%;
}
.page-buttons {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}
.page-btn {
  border: 1px solid #ddd;
  background: white;
  padding: 4px 10px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}
.page-btn.active {
  background: #4285f4;
  color: white;
  border-color: #4285f4;
}
.courseware-controls {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-top: 10px;
  flex-wrap: wrap;
}
.ctrl-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  background: #4285f4;
  color: white;
  cursor: pointer;
}
.ctrl-btn:disabled {
  background: #999;
  cursor: not-allowed;
}
.page-info {
  color: #666;
}
/* 聊天区 */
.chat-panel {
  width: 420px;
  padding: 20px;
  background: white;
  border-left: 1px solid #eee;
  display: flex;
  flex-direction: column;
}
.message-list {
  flex: 1;
  overflow-y: auto;
  margin-bottom: 20px;
  padding: 10px;
}
.message {
  margin-bottom: 15px;
  max-width: 90%;
}
.user-msg {
  margin-left: auto;
}
.page-tag {
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
}
.msg-content {
  padding: 8px 12px;
  border-radius: 8px;
  background: #eee;
}
.user-msg .msg-content {
  background: #4285f4;
  color: white;
}
.chat-input {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.input-area {
  height: 80px;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  resize: none;
}
.input-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.send-btn {
  padding: 8px 20px;
  background: #4285f4;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.send-btn:disabled {
  background: #999;
  cursor: not-allowed;
}
.page-tip {
  font-size: 12px;
  color: #666;
}
/* 思维导图弹窗 */
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
}
.modal-content {
  width: 600px;
  background: white;
  border-radius: 8px;
  padding: 20px;
}
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
}
.close-btn {
  background: transparent;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: #666;
}
.mindmap-tree {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 60vh;
  overflow-y: auto;
  padding-right: 6px;
}
.mindmap-node {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  line-height: 1.5;
  color: #333;
}
.mindmap-bullet {
  color: #4285f4;
  font-weight: 700;
}
.mindmap-empty {
  color: #777;
  font-size: 14px;
}

@media (max-width: 1200px) {
  .main-layout {
    flex-direction: column;
    height: auto;
  }
  .chat-panel {
    width: 100%;
    border-left: none;
    border-top: 1px solid #eee;
  }
  .courseware-panel {
    width: 100%;
  }
  .courseware-controls {
    justify-content: center;
  }
}

@media (max-width: 640px) {
  .top-nav {
    flex-direction: column;
    row-gap: 10px;
    text-align: center;
  }
  .top-actions {
    width: 100%;
    justify-content: center;
    gap: 10px;
  }
  .chat-panel {
    padding: 16px;
  }
  .courseware-controls {
    gap: 8px;
  }
}
</style>