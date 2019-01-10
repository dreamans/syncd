import md5 from 'blueimp-md5'

const state = {
    user_id: 0,
    group_name: '',
    name: '',
    email: '',
    mobile: '',
    true_name: '',
    priv: [],
}

const getters = {
    isLogin(state) {
        return state.userId > 0
    },
    getUserName(state) {
        return state.name
    },
    getUserId(state) {
        return state.user_id
    },
    getAvatar(state) {
        return 'https://www.gravatar.com/avatar/' + md5(state.email.toLowerCase()) + '?s=512'
    },
    getPriv(state) {
        return state.priv ? state.priv : []
    },
    getGroupName(state) {
        return state.group_name
    },
    getMobile(state) {
        return state.mobile
    },
    getTrueName(state) {
        return state.true_name
    },
    getEmail(state) {
        return state.email
    },
}

const actions = {
    login({ commit }, userInfo) {
        commit('setUserInfo', {
            user_id: userInfo.user_id,
            name: userInfo.name,
            email: userInfo.email,
            priv: userInfo.priv,
            group_name: userInfo.group_name,
            mobile: userInfo.mobile,
            true_name: userInfo.true_name,
        })
    },
}

const mutations = {
    setUserInfo(state, userInfo) {
        state.user_id = userInfo.user_id
        state.name = userInfo.name
        state.email = userInfo.email
        state.priv = userInfo.priv
        state.group_name = userInfo.group_name
        state.mobile = userInfo.mobile
        state.true_name = userInfo.true_name
    },
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
