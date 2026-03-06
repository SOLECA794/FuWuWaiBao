<<<<<<< HEAD:frontend/student/src/main.js
import { createApp } from 'vue'
import App from './App.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
// 关键：导入ElMessage并全局注册
import { ElMessage } from 'element-plus'

const app = createApp(App)
app.use(ElementPlus, { locale: zhCn })
// 关键：把ElMessage挂载到全局，让所有组件都能使用
app.config.globalProperties.$message = ElMessage
// 额外：也可以直接暴露ElMessage到window（兼容新手写法）
window.ElMessage = ElMessage

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
