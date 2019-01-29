import {get, post} from '@/lib/fetch.js'

export function applyProjectDetailApi(params) {
    return get('/deploy/apply/project/detail', params)
}

export function applySubmitApi(data) {
    return post('/deploy/apply/submit', data)
}