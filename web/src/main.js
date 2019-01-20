import Vue from 'vue'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import App from './App.vue'
import router from './router'
import store from './store'
import util from './lib/util.js'
import data from './lib/data.js'
import './scss/app.scss'

Vue.config.productionTip = false
Vue.use(ElementUI);

new Vue({
    router,
    store,
    methods: util,
    data: data,
    render: h => h(App)
}).$mount('#app')