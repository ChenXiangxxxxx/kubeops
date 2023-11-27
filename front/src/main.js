import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/antd.dark.css'
import * as Icons from '@ant-design/icons-vue'

const app =createApp(App)
for (const i in Icons){
    app.component(i,Icons[i])
}
app.use(Antd)
app.use(router)

app.mount('#app')
