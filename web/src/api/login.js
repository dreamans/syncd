import {get, post} from '@/lib/fetch.js'

export function loginApi(data) {
    return post('/login', data)
}

export function loginStatusApi() {
    return get("/login/status")
}

export function logoutApi() {
    return post('/logout')
}