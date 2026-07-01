import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './assets/main.css'
import 'github-markdown-css/github-markdown-light.css'

createApp(App).use(router).mount('#app')
