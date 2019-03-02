<template>
    <div class="layer-global">
        <header class="layer-header">
            <div class="header-left">
                <img class="logo" src="@/asset/logo.png" />
            </div>
            <div class="header-right">
                <span class="r-item">
                    <el-dropdown trigger="click">
                        <span class="item app-cursor">
                            <i class="iconfont icon-question-circle-fill"></i>
                            <i class="iconfont small icon-arrow-down"></i>
                        </span>
                        <el-dropdown-menu slot="dropdown" class="app-header-dropdown">
                            <a class="app-dropdown-link" href="https://github.com/dreamans/syncd/issues" target="_blank">
                                <el-dropdown-item><i class="iconfont small left icon-help"></i>{{ $t('help') }}</el-dropdown-item>
                            </a>
                            <a class="app-dropdown-link" href="https://github.com/dreamans/syncd" target="_blank">
                                <el-dropdown-item><i class="iconfont small left icon-pull-request"></i>{{ $t('contribute_to_syncd') }}</el-dropdown-item>
                            </a>
                        </el-dropdown-menu>
                    </el-dropdown>
                </span>
                <span class="r-item">
                    <el-dropdown trigger="click" @command="userSettingHandler">
                        <span class="item app-cursor">
                            <img class="avatar" :src="$store.getters['account/getAvatar']" />
                            <i class="iconfont small icon-arrow-down"></i>
                        </span>
                        <el-dropdown-menu slot="dropdown" class="app-header-dropdown">
                            <el-dropdown-item class="text"><i class="iconfont small left icon-user"></i>{{ $store.getters['account/getUserName'] }}</el-dropdown-item>
                            <el-dropdown-item command="setting" divided><i class="iconfont small left icon-setting"></i>{{ $t('personal_setting') }}</el-dropdown-item>
                            <el-dropdown-item command="password"><i class="iconfont small left icon-key"></i>{{ $t('change_password') }}</el-dropdown-item>
                            <el-dropdown-item command="logout" divided><i class="iconfont small left icon-logout"></i>{{ $t('sign_out') }}</el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                </span>
            </div>
        </header>
        <section class="layer-container">
            <aside class="layer-aside">
                <ScrollBar>
                    <el-menu class="aside-menu" :default-active="activeMenu" :router="true" :unique-opened="true">
                        <template v-for="menu in AppMenu">
                            <el-submenu v-if="menu.children && menu.children.length > 1" :index="menu.name" :key="menu.name">
                                <template slot="title">
                                    <span v-if="menu.meta.icon" class="iconfont left" :class="menu.meta.icon"></span><span>{{ menu.meta.title }}</span>
                                </template>
                                <template v-for="childMenu in menu.children">
                                    <el-menu-item v-if="!(childMenu.meta && childMenu.meta.hide)" :route="{name: childMenu.name}" :index="childMenu.name" :key="childMenu.name">
                                        <i class="iconfont small left">
                                            <svg viewBox="0 0 1024 1024" width="1em" height="1em" fill="currentColor" aria-hidden="true"><path d="M384.023552 384.083968l256.016384 0 0 256.016384-256.016384 0 0-256.016384Z"></path></svg>
                                        </i>
                                        <span>{{ childMenu.meta.title }}</span>
                                    </el-menu-item>
                                </template>
                            </el-submenu>
                            <el-menu-item :route="{name: menu.children[0].name}" v-else-if="menu.children && menu.children.length == 1" :index="menu.children[0].name" :key="menu.children[0].name">
                                <i v-if="menu.children[0].meta.icon" class="iconfont left" :class="menu.children[0].meta.icon"></i>
                                <span>{{ menu.children[0].meta.title }}</span>
                            </el-menu-item>
                        </template>
                    </el-menu>
                </ScrollBar>
            </aside>
            <main class="layer-main">
                <el-breadcrumb separator="/" class="bread-crumb">
                    <el-breadcrumb-item><i class="iconfont small icon-breadcrumbs"></i></el-breadcrumb-item>
                    <el-breadcrumb-item v-for="b in breadcrumb" :key="b">{{ b }}</el-breadcrumb-item>
                </el-breadcrumb>
                <div class="container">
                    <router-view/>
                    <div class="app-cpy">
                        <p>
                        ©️ {{ new Date().getFullYear() }} <a href="https://github.com/dreamans/syncd" target="_blank">Syncd</a>. All Rights Reserved. MIT License.
                        </p>
                    </div>
                </div>
            </main>
        </section>

        <el-dialog :width="$root.DialogSmallWidth" title="个人设置" :visible.sync="settingDialogVisible" @close="closeUserSettingDialogHandler">
            <el-form class="app-form" ref="settingDialogRef" :model="settingForm" size="medium" label-width="80px">
                <el-form-item 
                :label="$t('role')">
                    {{ $store.getters['account/getRoleName'] }}
                </el-form-item>
                <el-form-item 
                :label="$t('avatar')">
                    <img :src="$store.getters['account/getAvatar']" class="app-avatar normal">
                    &nbsp;<a class="app-link" href="http://cn.gravatar.com/support/activating-your-account/" target="_blank"><i class="iconfont small left icon-question"></i>如何修改头像</a>
                </el-form-item>
                <el-form-item 
                :label="$t('username')">
                    <el-input :value="$store.getters['account/getUserName']" :disabled="true"></el-input>
                    <span class="app-input-help"><i class="el-icon-info"></i> 修改用户名请联系管理员</span>
                </el-form-item>
                <el-form-item 
                :label="$t('email')">
                    <el-input :value="$store.getters['account/getEmail']" :disabled="true"></el-input>
                    <span class="app-input-help"><i class="el-icon-info"></i> 修改邮箱请联系管理员</span>
                </el-form-item>
                <el-form-item 
                :label="$t('truename')"
                prop="truename">
                    <el-input v-model="settingForm.truename"></el-input>
                </el-form-item>
                <el-form-item 
                :label="$t('mobile')"
                prop="mobile">
                    <el-input v-model="settingForm.mobile"></el-input>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button size="small" @click="closeUserSettingDialogHandler">{{ $t('cancel')}}</el-button>
                <el-button :loading="btnLoading" size="small" type="primary" @click="dialogUserSettingSubmitHandler">{{ $t('enter')}}</el-button>
            </div>
        </el-dialog>
    
        <el-dialog :width="$root.DialogSmallWidth" title="修改密码" :visible.sync="passwordDialogVisible" @close="closePasswordSettingDialogHandler">
            <el-form class="app-form" ref="passwordDialogRef" :model="passwordForm" size="medium" label-width="80px">
                <el-form-item 
                :label="$t('current_password')"
                prop="password"
                :rules="[
                    { min: 6, max: 20, message: this.$t('strlen_between', {min: 6, max: 20}), trigger: 'blur'},
                    { required: true, message: this.$t('current_password_cannot_empty'), trigger: 'blur'},
                ]">
                    <el-input type="password" :placeholder="$t('please_input_password_length_limit')" v-model="passwordForm.password" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item 
                :label="$t('new_password')"
                prop="new_password"
                :rules="[
                    { min: 6, max: 20, message: this.$t('strlen_between', {min: 6, max: 20}), trigger: 'blur'},
                    { required: true, message: this.$t('new_password_cannot_empty'), trigger: 'blur'},
                ]">
                    <el-input type="password" :placeholder="$t('please_input_password_length_limit')" v-model="passwordForm.new_password" autocomplete="off"></el-input>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button size="small" @click="closePasswordSettingDialogHandler">{{ $t('cancel')}}</el-button>
                <el-button :loading="btnLoading" size="small" type="primary" @click="dialogPasswordSettingSubmitHandler">{{ $t('enter')}}</el-button>
            </div>
        </el-dialog>

    </div>
</template>

<script>
import ScrollBar from '@/component/ScrollBar';
import { routerMap } from '@/router'
import Code from '@/lib/code'
import { loginStatusApi, logoutApi } from '@/api/login'
import { userSettingApi, userPasswordApi } from '@/api/user'
export default {
    data() {
        return {
            breadcrumb: [],
            activeMenu: '',

            btnLoading: false,
            settingDialogVisible: false,
            settingForm: {},
            passwordDialogVisible: false,
            passwordForm: {},
        }
    },
    computed: {
        AppMenu() {
            let menu = []
            let currName = this.$route.name
            routerMap.forEach(first => {
                let newSecond = []
                let newFirst = Object.assign({}, first)
                first.children.forEach(second => {
                    newSecond.push(second)
                })
                newFirst.children = newSecond
                menu.push(newFirst)
            })
            return routerMap
        }
    },
    watch: {
        '$route.name'() {
            this.breadcrumbItems()
            this.initActiveMenu()
        },
    },
    components: {
        ScrollBar,
    },
    methods: {
        userSettingHandler(cmd) {
            switch (cmd) {
                case 'logout':
                    logoutApi().then(res => {
                        this.$router.push({name: 'login'})
                    })
                    break
                case 'setting':
                    this.showUserSettingDialogHandler()
                    break
                case 'password':
                    this.showPasswordSettingDialogHandler()
                    break
            }
        },

        showPasswordSettingDialogHandler() {
            this.passwordDialogVisible = true
        },
        closePasswordSettingDialogHandler() {
            this.passwordDialogVisible = false
        },
        dialogPasswordSettingSubmitHandler() {
            let vm = this
            this.$refs.passwordDialogRef.validate((valid) => {
                if (!valid) {
                    return false;
                }
                let postData = {
                    password: this.$root.Md5Sum(vm.passwordForm.password),
                    new_password: this.$root.Md5Sum(vm.passwordForm.new_password)
                }
                this.btnLoading = true
                userPasswordApi(postData).then(res => {
                    this.$root.MessageSuccess(() => {
                        this.btnLoading = false
                        this.closePasswordSettingDialogHandler()
                    })
                }).catch(err => {
                    if (err.code == Code.CODE_ERR_USER_OR_PASS_WRONG) {
                        this.$message({
                            message: '当前密码错误，请重新输入',
                            type: 'warning',
                        });
                    }
                    this.btnLoading = false
                })
            })
        },

        showUserSettingDialogHandler() {
            this.settingForm = {
                truename: this.$store.getters['account/getTrueName'],
                mobile: this.$store.getters['account/getMobile'],
            }
            this.settingDialogVisible = true
        },
        closeUserSettingDialogHandler() {
            this.settingDialogVisible = false
        },
        dialogUserSettingSubmitHandler() {
            this.btnLoading = true
            userSettingApi(this.settingForm).then(res => {
                this.$root.MessageSuccess(() => {
                    this.btnLoading = false
                    this.closeUserSettingDialogHandler()
                })
                this.$store.dispatch('account/userSetting', {
                    mobile: this.settingForm.mobile,
                    truename: this.settingForm.truename,
                })
            }).catch(err => {
                this.btnLoading = false
            })
        },

        initActiveMenu() {
            this.activeMenu = this.$route.name
        },
        breadcrumbItems() {
            let breadcrumb = []
            this.AppMenu.forEach(menu => {
                menu.children.forEach(sub => {
                    if (sub.name != this.$route.name) {
                        return
                    }
                    if (menu.meta.title) {
                        breadcrumb.push(menu.meta.title)
                    }
                    breadcrumb.push(sub.meta.title)
                })
            })
            this.breadcrumb = breadcrumb
        },
        initLoginStatus() {
            loginStatusApi().then(res => {
                if (res.is_login) {
                    this.$store.dispatch('account/status', {
                        user_id: res.user_id,
                        username: res.username,
                        email: res.email,
                        mobile: res.mobile,
                        privilege: res.privilege ? res.privilege : [],
                        role_name: res.role_name,
                        truename: res.truename,
                    })
                } else {
                    this.$message.error('用户未登录', 1)
                    this.$router.push({name: 'login'})
                }
            }).catch(err => {
                this.$message.error('用户未登录', 1)
                this.$router.push({name: 'login'})
            })
        },
    },
    mounted() {
        this.initLoginStatus()
        this.breadcrumbItems()
        this.initActiveMenu()
    },
}
</script>
