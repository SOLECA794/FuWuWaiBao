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
