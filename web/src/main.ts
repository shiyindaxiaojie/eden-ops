import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'element-plus/dist/index.css'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import App from './App.vue'
import router from './router'
import './styles/index.scss'
import waves from '@/directive/waves'
import Pagination from '@/components/Pagination/index.vue'

NProgress.configure({ showSpinner: false })

const app = createApp(App)

// 注册Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 使用Pinia替代Vuex
const pinia = createPinia()
app.use(pinia)
app.use(router)
app.use(ElementPlus)

// 注册全局指令
app.directive('waves', waves)

// 注册全局组件
app.component('Pagination', Pagination)

app.mount('#app')