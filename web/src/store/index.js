import Vue from 'vue'
import Vuex from 'vuex'

import account from './module/account.js'

Vue.use(Vuex)

export default new Vuex.Store({
    modules: {
        account,
    }
})