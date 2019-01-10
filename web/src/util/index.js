import moment from 'moment'
import md5 from 'blueimp-md5'
import Cookies from 'js-cookie'

let loginTokenKey = 'SYD_AUTH_TOKEN'

export default {

    GotoRouter(name, query) {
        this.$router.push({name, query})
    },

    ResolveFormError(err, values) {
        let key = Object.getOwnPropertyNames(err).shift()
        if (key && err[key] && err[key].errors && err[key].errors.length) {
            let msg = err[key].errors[0].message
            let field = err[key].errors[0].field
            this.$message.error(msg);

            let targetTop = this.$jQuery("#" + field).offset().top
            this.$jQuery('body,html').animate({scrollTop: targetTop - 20}, 1000);
        }
    },

    FormatDateTime(unixtime, format) {
        if (!unixtime) {
            return '--'
        }
        if (!format) {
            format = 'YYYY-MM-DD HH:mm:SS'
        }
        return moment.unix(unixtime).format(format)
    },

    FormatDateFromNow(unixtime) {
        if (!unixtime) {
            return '--'
        }
        return moment.unix(unixtime).fromNow()
    },

    Md5Sum(str) {
        return md5(str);
    },

    ResetPagination(pagination, adjustNum) {
        if (!adjustNum) {
            adjustNum = -1
        }
        pagination.total += adjustNum
        let maxPage = Math.ceil(pagination.total/pagination.pageSize)
        if (pagination.current > maxPage) {
            pagination.current = maxPage
        }
    },

    SetLoginToken(token) {
        return Cookies.set(loginTokenKey, token)
    },

    GetLoginToken() {
        return Cookies.get(loginTokenKey)
    },

    DeleteLoginToken() {
        return Cookies.remove(loginTokenKey)
    },

    CheckPriv(code) {
        return this.$store.getters['account/getPriv'].indexOf(code) > -1
    },

    CheckPrivs(codeArr) {
        if (!codeArr || !codeArr.length) {
            return false
        }
        let checked = false
        codeArr.forEach(code => {
            if (this.CheckPriv(code)) {
                checked = true
            }
        })
        return checked
    },

    T(lang) {
        if (this.$root.Lang[lang]) {
            return this.$root.Lang[lang]
        }
        return lang
    },
}
