<<<<<<< HEAD:frontend/student/src/main.js
import { createApp } from 'vue'
import App from './App.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import { ElMessage } from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

const app = createApp(App)
app.use(ElementPlus, { locale: zhCn })
app.config.globalProperties.$message = ElMessage
window.ElMessage = ElMessage

<<<<<<< Updated upstream
app.mount('#app')
=======
import { createApp } from 'vue'
import App from './App.vue'
// 引入Element Plus核心库和样式
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
// 引入Element Plus图标库（全局注册，避免导入错误）
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// 创建Vue实例
const app = createApp(App)
// 注册Element Plus
app.use(ElementPlus)
// 全局注册所有Element Plus图标（彻底解决图标导出错误）
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
// 挂载到根节点
app.mount('#app')
>>>>>>> b79727f64ad1860d8e9dc554eec4fdaef2859d48:src/main.js
=======
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.mount('#app')
>>>>>>> Stashed changes
