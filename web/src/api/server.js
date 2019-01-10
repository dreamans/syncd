import {get, post} from './fetch.js'

export function newGroupApi(data) {
    return post('/server/group/new', data)
}

export function updateGroupApi(data) {
    return post('/server/group/update', data)
}

export function getGroupListApi(params) {
    return get('/server/group/list', params)
}

export function getGroupMultiApi(params) {
    return get('/server/group/multi', params)
}

export function getGroupDetailApi(params) {
    return get('/server/group/detail', params)
}

export function deleteGroupApi(data) {
    return post('/server/group/delete', data)
}

export function newServerApi(data) {
    return post('/server/new', data)
}

export function updateServerApi(data) {
    return post('/server/update', data)
}

export function getServerListApi(params) {
    return get('/server/list', params)
}

export function getServerDetailApi(params) {
    return get('/server/detail', params)
}

export function deleteServerApi(data) {
    return post('/server/delete', data)
}
