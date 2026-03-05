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
