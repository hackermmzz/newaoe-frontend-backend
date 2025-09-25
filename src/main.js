import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import 'tailwindcss/tailwind.css'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css' // 必须导入样式
import { ElMessage } from 'element-plus' // 导入ElMessage

const app = createApp(App)
app.use(ElementPlus)
app.config.globalProperties.$message = ElMessage // 挂载到全局
app.use(router).mount('#app')
