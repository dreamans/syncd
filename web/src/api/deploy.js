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
