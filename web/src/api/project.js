import {get, post} from './fetch.js'

export function newProjectApi(data) {
    return post('/project/new', data)
}

export function updateProjectApi(data) {
    return post('/project/update', data)
}

export function listProjectApi(params) {
    return get('/project/list', params)
}

export function getProjectApi(params) {
    return get('/project/detail', params)
}

export function deleteProjectApi(data) {
    return post('/project/delete', data)
}

export function changeProjectStatusApi(data) {
    return post('/project/status/change', data)
}

export function getSpaceListApi(params) {
    return get('/project/space/list', params)
}

export function newSpaceApi(data) {
    return post('/project/space/new', data)
}

export function updateSpaceApi(data) {
    return post('/project/space/update', data)
}

export function getSpaceDetailApi(params) {
    return get('/project/space/detail', params)
}

export function deleteSpaceApi(data) {
    return post('/project/space/delete', data)
}

export function checkSpaceNameExistsApi(params) {
    return get('/project/space/exists', params)
}

export function checkProjectNameExistsApi(params) {
    return get('/project/exists', params)
}

export function addSpaceUserApi(data) {
    return post('/project/user/add', data)
}

export function getSpaceUserListApi(params) {
    return get('/project/user/list')
}

export function removeSpaceUserApi(data) {
    return post('/project/user/remove', data)
}

export function searchSpaceUserApi(data) {
    return get('/project/user/search', data)
}

export function resetRepoApi(data) {
    return post('/project/repo/reset', data)
}

export function checkServerApi(params) {
    return get('/project/server/check', params)
}
