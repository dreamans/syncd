import Vue from 'vue'
import VueI18n from 'vue-i18n'

Vue.use(VueI18n)

import zh from './zh.js'
import en from './en.js'
const messages = {
    'zh': zh,
    'en': en,
}

let localeLang = 'en'
if (global.navigator.language) {
    localeLang = global.navigator.language.substr(0, 2)
    localeLang = localeLang.toLowerCase()
}

const i18n = new VueI18n({
    locale: localeLang,
    messages,
})

export default i18n