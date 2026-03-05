import { createApp } from 'vue'
import App from './App.vue'
// 引入 Element Plus
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
// 引入富文本编辑器样式
import '@vueup/vue-quill/dist/vue-quill.snow.css'

const app = createApp(App)
// 全局注册 Element Plus
app.use(ElementPlus)
app.mount('#app')
