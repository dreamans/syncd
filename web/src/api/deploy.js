import {get, post} from './fetch.js'

export function submitApplyApi(data) {
    return post('/deploy/apply/submit', data)
}

export function getApplyListApi(params) {
    return get('/deploy/apply/list', params)
}

export function getApplySpaceList(params) {
    return get('/deploy/apply/space/list', params)
}

export function  getApplyProjectList(params) {
    return get('/deploy/apply/project/list', params)
}

export function  getApplyProjectAll(params) {
    return get('/deploy/apply/project/all', params)
}

export function getApplyProjectDetailApi(params) {
    return get('/deploy/apply/project/detail', params)
}

export function getRepoTagListApi(params) {
    return get('/deploy/apply/repo/taglist', params)
}

export function getRepoCommitListApi(params) {
    return get('/deploy/apply/repo/commitlist', params)
}

export function getApplyDetailApi(params) {
    return get('/deploy/apply/detail', params)
}

export function auditApplyApi(data) {
    return post('/deploy/apply/audit', data)
}

export function unAuditApplyApi(data) {
    return post('/deploy/apply/unaudit', data)
}

export function discardApplyApi(data) {
    return post('/deploy/apply/discard', data)
}

export function updateApplyApi(data) {
    return post('/deploy/apply/update', data)
}

export function getOperateLogApi(params) {
    return get('/deploy/apply/log', params)
}

export function startDeployApi(data) {
    return post('/deploy/deploy/start', data)
}

export function statusDeployApi(params) {
    return get('/deploy/deploy/status', params)
}

export function stopDeployApi(data) {
    return post('/deploy/deploy/stop', data)
}
