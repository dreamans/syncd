<template>
    <div class="app-login">
        <div class="app-login-inner" :style="{backgroundColor: bgColor}">
            <div class="login-container">
                <el-card class="login-box">
                    <div class="login-title">欢迎登录 Syncd</div>
                    <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" size="medium" class="login-form">
                        <el-form-item  prop="username">
                            <el-input placeholder="用户名或邮箱" prefix-icon="iconfont icon-user" v-model="loginForm.username" autocomplete="off"></el-input>
                        </el-form-item>
                        <el-form-item prop="password">
                            <el-input placeholder="密码" prefix-icon="iconfont icon-lock" type="password" v-model="loginForm.password" autocomplete="off"></el-input>
                        </el-form-item>
                        <el-form-item>
                            <el-button @click="loginHandler" type="primary" style="width: 100%;">登录</el-button>
                        </el-form-item>
                    </el-form>
                </el-card>
            </div>
            <div class="login-cpy">
                © {{ new Date().getFullYear() }} <a href="https://github.com/dreamans/syncd" target="_blank">Syncd</a> 版权所有, MIT License.
            </div>
        </div>
    </div>
</template>

<script>
export default {
    data () {
        return {
            loginLoadding: false,
            loginForm: {
                username: "",
                password: "",
            },
            loginRules: {
                username: [{ required: true, message: '请输入登录名(用户名或邮箱)', trigger: 'blur' }],
                password: [
                    { required: true, message: '请输入密码', trigger: 'blur' },
                    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' },
                ],
            },
        }
    },
    computed: {
        bgColor() {
            let alpha = Math.min(Math.abs(new Date().getHours() - 12) / 10, 0.80)
            return `rgba(0, 0, 0, ${alpha})`
        },
    },
    methods: {
        loginHandler() {
            this.$refs.loginFormRef.validate((valid) => {
                if (!valid) {
                    return false
                }
            });
        },
    },
}
</script>

<style lang="scss" scope>
.app-login {
    width: 100%;
    height: 100%;
    background-image: url('../asset/login_bg.jpg');
    background-position: 50%;
    background-size: cover;
    .app-login-inner {
        width: 100%;
        height: 100%;
    }
    .login-container {
        display: flex;
        justify-content: center;
        .login-box {
            margin-top: 20vh;
            width: 30vw;
            .login-title {
                font-weight: 500;
                text-align: center;
                font-size: 14px;
                margin-bottom: 20px;
            }
        }
    }
    .login-cpy {
        display: flex;
        justify-content: center;
        margin-top: 30px;
        color: #fff;
        a {
            margin: 0 5px;
            color: #1890ff;
        }
    }
}
</style>