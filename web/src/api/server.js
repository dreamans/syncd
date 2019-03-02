import {get, post} from '@/lib/fetch.js'

export function newGroupApi(data) {
    return post('/server/group/add', data)
}

export function listGroupApi(params) {
    return get('/server/group/list', params)
}

export function deleteGroupApi(data) {
    return post('/server/group/delete', data)
}

export function detailGroupApi(params) {
    return get('/server/group/detail', params)
}

export function updateGroupApi(data) {
    return post('/server/group/update', data)
}

export function newServerApi(data) {
    return post('/server/add', data)
}

export function updateServerApi(data) {
    return post('/server/update', data)
}

export function listServerApi(params) {
    return get('/server/list', params)
}

export function deleteServerApi(data) {
    return post('/server/delete', data)
}

export function detailServerApi(params) {
    return get('/server/detail', params)
}