import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'

import App from './App.vue'
import router from './router'
import { setupRouterGuard } from './router/guard'
import { permission } from './directives/permission'

import './styles/variables.scss'
import './styles/reset.scss'
import './styles/dark.scss'
import './styles/transitions.scss'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ElementPlus, { locale: zhCn })

app.directive('permission', permission)

setupRouterGuard(router)

app.mount('#app')
