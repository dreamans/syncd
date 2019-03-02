import md5 from 'blueimp-md5'

const state = {
    user_id: 0,
    username: '',
    email: '',
    mobile: '',
    privilege: [],
    role_name: '',
    truename: '',
}

const getters = {
    isLogin(state) {
        return state.user_id > 0
    },
    getEmail(state) {
        return state.email
    },
    getUserName(state) {
        return state.username
    },
    getUserId(state) {
        return state.user_id
    },
    getAvatar(state) {
        return 'https://www.gravatar.com/avatar/' + md5(state.email.toLowerCase()) + '? s=512'
    },
    getPriv(state) {
        return state.privilege ? state.privilege : []
    },
    getRoleName(state) {
        return state.role_name
    },
    getMobile(state) {
        return state.mobile
    },
    getTrueName(state) {
        return state.truename
    },
}

const actions = {
    status({ commit }, userInfo) {
        commit('setUserInfo', {
            user_id: userInfo.user_id,
            username: userInfo.username,
            email: userInfo.email,
            mobile: userInfo.mobile,
            privilege: userInfo.privilege,
            role_name: userInfo.role_name,
            truename: userInfo.truename,
        })
    },
    userSetting({ commit }, userInfo) {
        commit('userSetting', {
            mobile: userInfo.mobile,
            truename: userInfo.truename,
        })
    },
}

const mutations = {
    setUserInfo(state, userInfo) {
        state.user_id = userInfo.user_id
        state.username = userInfo.username
        state.email = userInfo.email
        state.privilege = userInfo.privilege
        state.role_name = userInfo.role_name
        state.mobile = userInfo.mobile
        state.truename = userInfo.truename
    },
    userSetting(state, userInfo) {
        state.mobile = userInfo.mobile
        state.truename = userInfo.truename
    },
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
