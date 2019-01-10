<template>
    <div class="app-login">
        <div class="app-login-inner" :style="{backgroundColor: bgColor}">
            <div class="login-container">
                <a-card
                :bordered="false"
                class="login-box">
                    <div class="login-title">欢迎登录 Syncd</div>
                    <a-form :form="loginForm" @submit="handleLogin">
                        <a-form-item>
                            <a-input
                            placeholder='用户名或邮箱'
                            autocomplete="off"
                            v-decorator="[
                                'name',
                                {
                                    rules: [{ required: true, message: '请输入登录名，支持用户名、邮箱登录' }],
                                    validateTrigger: 'blur',
                                }
                            ]">
                                <a-icon slot="prefix" type='user' style="color: rgba(0,0,0,.25)" />
                            </a-input>
                        </a-form-item>
                        <a-form-item>
                            <a-input
                            placeholder='密码'
                            autocomplete="off"
                            type='password'
                            v-decorator="[
                                'pass',
                                {
                                    rules: [{ required: true, message: '请输入密码' }],
                                    validateTrigger: 'blur',
                                }
                            ]">
                                <a-icon slot="prefix" type='lock' style="color: rgba(0,0,0,.25)" />
                            </a-input>
                        </a-form-item>
                        <a-form-item>
                            <a-button style="width: 100%;" type='primary' htmlType='submit'>登录</a-button>
                        </a-form-item>
                    </a-form>
                </a-card>
            </div>
            <div class="login-cpy">
                © {{ new Date().getFullYear() }} <a href="https://github.com/tinystack/syncd" target="_blank">Syncd</a> 版权所有, MIT License.
            </div>
        </div>
    </div>
</template>

<script>
import { loginApi, loginStatusApi } from '@/api/user.js'
export default {
    beforeCreate() {
        this.loginForm = this.$form.createForm(this)
    },
    data () {
        return {
            loginLoadding: false,
        }
    },
    computed: {
        bgColor() {
            let alpha = Math.min(Math.abs(new Date().getHours() - 12) / 10, 0.80)
            return `rgba(0, 0, 0, ${alpha})`
        },
    },
    methods: {
        handleLogin(e) {
            e.preventDefault()
            this.loginForm.validateFields((err, values) => {
                if (err) {
                    this.$root.ResolveFormError(err, values)
                    return
                }
                let postData = {
                    name: values.name,
                    pass: this.$root.Md5Sum(values.pass)
                }
                loginApi(postData).then(res => {
                    this.$store.dispatch('account/login', res)
                    this.$root.SetLoginToken(res.token)
                    this.$root.GotoRouter('dashboard')
                }).catch(err => {
                    this.$message.error(`登录失败, ${err.message}`)
                })
            })
        },
    },
    mounted() {
        loginStatusApi().then(res => {
            if (res.is_login) {
                this.$root.GotoRouter('dashboard')
            }
        }).catch(err => {})
    },
}
</script>
