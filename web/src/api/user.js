import {get, post} from './fetch.js'

export function newGroupApi(data) {
    return post('/user/group/new', data)
}

export function updateGroupApi(data) {
    return post('/user/group/update', data)
}

export function getGroupListApi(params) {
    return get('/user/group/list', params)
}

export function getGroupDetailApi(params) {
    return get('/user/group/detail', params)
}

export function deleteGroupApi(data) {
    return post('/user/group/delete', data)
}

export function getPrivListApi() {
    return get('/user/group/priv')
}

export function newUserApi(data) {
    return post('/user/new', data)
}

export function updateUserApi(data) {
    return post('/user/update', data)
}

export function getUserListApi(params) {
    return get('/user/list', params)
}

export function getUserDetailApi(params) {
    return get('/user/detail', params)
}

export function checkUserExistsApi(params) {
    return get('/user/checkexists', params)
}

export function deleteUserApi(data) {
    return post('/user/delete', data)
}

export function searchUserApi(params) {
    return get('/user/search', params)
}

export function loginApi(data) {
    return post('/user/login', data)
}

export function logoutApi() {
    return post('/user/logout')
}

export function loginStatusApi() {
    return get('/user/login/status')
}

export function updateMyApi(data) {
    return post('/user/my/update', data)
}

export function updatePasswordApi(data) {
    return post('/user/my/password', data)
}
