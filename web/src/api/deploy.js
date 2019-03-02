import {get, post} from '@/lib/fetch.js'

export function applyProjectDetailApi(params) {
    return get('/deploy/apply/project/detail', params)
}

export function applySubmitApi(data) {
    return post('/deploy/apply/submit', data)
}

export function applyProjectAllApi() {
    return get('/deploy/apply/project/all')
}

export function applyListApi(params) {
    return get('/deploy/apply/list', params)
}

export function applyDetailApi(params) {
    return get('/deploy/apply/detail', params)
}

export function applyAuditApi(data) {
    return post('/deploy/apply/audit', data)
}

export function applyUpdateApi(data) {
    return post('/deploy/apply/update', data)
}

export function applyDropApi(data) {
    return post('/deploy/apply/drop', data)
}

export function applyRollbackListApi(params) {
    return get('/deploy/apply/rollbacklist', params)
}

export function buildStartApi(data) {
    return post('/deploy/build/start', data)
}

export function buildStatusApi(params) {
    return get('/deploy/build/status', params)
}

export function buildStopApi(data) {
    return post('/deploy/build/stop', data)
}

export function deployStart(data) {
    return post('/deploy/deploy/start', data)
}

export function deployStatusApi(params) {
    return get('/deploy/deploy/status', params)
}

export function deployStopApi(data) {
    return post('/deploy/deploy/stop', data)
}

export function deployRollbackApi(data) {
    return post('/deploy/deploy/rollback', data)
}