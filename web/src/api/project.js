import {get, post} from '@/lib/fetch.js'

export function newSpaceApi(data) {
    return post('/project/space/add', data)
}

export function updateSpaceApi(data) {
    return post('/project/space/update', data)
}

export function listSpaceApi(params) {
    return get('/project/space/list', params)
}

export function detailSpaceApi(params) {
    return get('/project/space/detail', params)
}

export function deleteSpaceApi(data) {
    return post('/project/space/delete', data)
}

export function searchMemberApi(params) {
    return get('/project/member/search', params)
}

export function addMemberApi(data) {
    return post('/project/member/add', data)
}

export function listMemberApi(params) {
    return get('/project/member/list', params)
}

export function removeMemberApi(data) {
    return post('/project/member/remove', data)
}

export function newProjectApi(data) {
    return post('/project/add', data)
}

export function updateProjectApi(data) {
    return post('/project/update', data)
}

export function listProjectApi(params) {
    return get('/project/list', params)
}

export function switchStatusProjectApi(data) {
    return post('/project/switchstatus', data)
}

export function detailProjectApi(params) {
    return get('/project/detail', params)
}

export function deleteProjectApi(data) {
    return post('/project/delete', data)
}

export function updateBuildScriptApi(data) {
    return post('/project/buildscript', data)
}

export function updateHookScriptApi(data) {
    return post('/project/hookscript', data)
}
