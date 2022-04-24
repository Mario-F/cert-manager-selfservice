import { createApp } from 'vue'
import { createPinia } from 'pinia'
import '@coreui/coreui/dist/css/coreui.min.css'
import '@coreui/icons/css/free.min.css'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
