import axios from 'axios'
import qs from 'qs'
import Vue from 'vue'
import i18n from '@/lang'

let API_URL = '/api'
let CancelToken = axios.CancelToken

Vue.prototype.$CancelAjaxRequet = function() {}
Vue.prototype.$IsCancel = axios.isCancel

const service = axios.create({
    baseURL: API_URL + '/',
    timeout: 60000,
    withCredentials: true,
})

service.interceptors.request.use(config => {
    return config
}, error => {
    Promise.reject(error)
})

service.interceptors.response.use(response => {
    let res = response.data
    if (!res) {
        res = {
            code: -1,
            message: i18n.t("network_error"),
        }
    }
    if (res.code != 0) {
        return Promise.reject(res)
    }
    return res.data
}, error => {
    if (!axios.isCancel(error)) {
        let res = {
            code: -1,
            message: error.message ? error.message : i18n.t("unknown_error"),
        }
        return Promise.reject(res)
    }
    return Promise.reject(error)
})

export function post(url, data, params, headers) {
    if (!params) {
        params = {}
    }
    params._t = new Date().getTime()
    let config = {
        method: 'post',
        url: url,
        params,
    }
    if (data) {
        if (headers && headers['Content-Type'] == 'multipart/form-data') {
            config.data = data
        } else {
            config.data = qs.stringify(data, { indices: false })
        }
    }
    if (headers) {
        config.headers = headers
    }

    config.cancelToken = new CancelToken(function(cancel) {
        Vue.prototype.$CancelAjaxRequet = function() {
            cancel()
        }
    })

    return service(config)
}

export function get(url, params, headers) {
    if (!params) {
        params = {}
    }
    params._t = new Date().getTime()
    let config = {
        method: 'get',
        url: url,
        params,
    }
    if (headers) {
        config.headers = headers
    }

    config.cancelToken = new CancelToken(function(cancel) {
        Vue.prototype.$CancelAjaxRequet = function() {
            cancel()
        }
    })

    return service(config)
}

export default service;