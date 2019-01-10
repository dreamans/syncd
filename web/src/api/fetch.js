import axios from 'axios'
import qs from 'qs'
import Vue from 'vue'

import Code from './code.js'
import Router from '@/router'

let API_URL = 'http://localhost:8868/api'
let CancelToken = axios.CancelToken
let source = CancelToken.source()

Vue.prototype.$CancelAjaxRequet = function() {}
Vue.prototype.$IsCancel = axios.isCancel

const service = axios.create({
    baseURL: API_URL + '/',
    timeout: 60000,
    withCredentials: true,
});

service.interceptors.request.use(config => {
    return config
}, error => {
    Promise.reject(error)
});

service.interceptors.response.use(response => {
    let res = response.data
    if (!res) {
        res = {
            code: -1,
            message: '网络错误',
        }
    }
    if (res.code != 0) {
        switch (res.code) {
            case Code.CODE_ERR_SYSTEM:
            case Code.CODE_ERR_APP:
            case Code.CODE_ERR_PARAM:
                Vue.prototype.$message.error(res.message)
                break
            case Code.CODE_ERR_NO_PRIV:
                Vue.prototype.$message.error('无操作权限')
                //Router.push({name: 'dashboard'})
                break
            case Code.CODE_ERR_NO_LOGIN:
                Router.push({name: 'login'})
                break
        }
        return Promise.reject(res)
    }
    return res.data;
}, error => {
    if (!axios.isCancel(error)) {
        let res = {
            code: -1,
            message: error.message ? error.message : '未知错误',
        }
        Vue.prototype.$message.error(res.message)
        return Promise.reject(res)
    }
    return Promise.reject(error)
});

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
};

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
