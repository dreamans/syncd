import {get, post} from '@/lib/fetch.js'

export function privListApi() {
    return get('/user/role/privlist')
}

export function newRoleApi(data) {
    return post('/user/role/add', data)
}

export function listRoleApi(params) {
    return get('/user/role/list', params)
}

export function detailRoleApi(params) {
    return get('/user/role/detail', params)
}

export function updateRoleApi(data) {
    return post('/user/role/update', data)
}

export function deleteRoleApi(data) {
    return post('/user/role/delete', data)
}

export function newUserApi(data) {
    return post('/user/add', data)
}

export function updateUserApi(data) {
    return post('/user/update', data)
}

export function listUserApi(params) {
    return get('/user/list', params)
}

export function existsUserApi(params) {
    return get('/user/exists', params)
}

export function detailUserApi(params) {
    return get('/user/detail', params)
}

export function deleteUserApi(data) {
    return post('/user/delete', data)
}

export function userSettingApi(data) {
    return post('/user/my/setting', data)
}

export function userPasswordApi(data) {
    return post('/user/my/password', data)
}
