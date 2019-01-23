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
                    <el-dropdown trigger="click">
                        <span class="item app-cursor">
                            <img class="avatar" src="https://www.gravatar.com/avatar/a8cb36173b0e2e1489c66965be0026cf" />
                            <i class="iconfont small icon-arrow-down"></i>
                        </span>
                        <el-dropdown-menu slot="dropdown" class="app-header-dropdown">
                            <el-dropdown-item class="text"><i class="iconfont small left icon-user"></i>admin</el-dropdown-item>
                            <el-dropdown-item divided><i class="iconfont small left icon-setting"></i>{{ $t('personal_setting') }}</el-dropdown-item>
                            <el-dropdown-item><i class="iconfont small left icon-key"></i>{{ $t('change_password') }}</el-dropdown-item>
                            <el-dropdown-item divided><i class="iconfont small left icon-logout"></i>{{ $t('sign_out') }}</el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                </span>
            </div>
        </header>
        <section class="layer-container">
            <aside class="layer-aside">
                <ScrollBar>
                    <el-menu class="aside-menu" :router="true" :unique-opened="true">
                        <template v-for="menu in AppMenu">
                            <el-submenu v-if="menu.children && menu.children.length > 1" :index="menu.name" :key="menu.name">
                                <template slot="title">
                                    <span v-if="menu.meta.icon" class="iconfont left" :class="menu.meta.icon"></span><span>{{ menu.meta.title }}</span>
                                </template>
                                <el-menu-item :route="{name: childMenu.name}" v-for="childMenu in menu.children" :index="childMenu.name" :key="childMenu.name">
                                    <i class="iconfont small left">
                                        <svg viewBox="0 0 1024 1024" width="1em" height="1em" fill="currentColor" aria-hidden="true"><path d="M384.023552 384.083968l256.016384 0 0 256.016384-256.016384 0 0-256.016384Z"></path></svg>
                                    </i>
                                    <span>{{ childMenu.meta.title }}</span>
                                </el-menu-item>
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
                    <el-breadcrumb-item><i class="iconfont icon-breadcrumbs"></i></el-breadcrumb-item>
                    <el-breadcrumb-item v-for="b in breadcrumb" :key="b">{{ b }}</el-breadcrumb-item>
                </el-breadcrumb>
                <div class="container">
                    <router-view/>
                </div>
            </main>
        </section>
    </div>
</template>

<script>
import ScrollBar from '@/component/ScrollBar';
import { routerMap } from '@/router'
export default {
    data() {
        return {
            breadcrumb: [],
        }
    },
    computed: {
        AppMenu() {
            return routerMap
        }
    },
    watch: {
        '$route.name'() {
            this.breadcrumbItems()
        }
    },
    components: {
        ScrollBar,
    },
    methods: {
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
        }
    },
    mounted() {
        this.breadcrumbItems()
    },
}
</script>

<style lang="scss" scoped>
.layer-global {
    height: 100%;
    .layer-header {
        z-index: 1024;
        position: fixed;
        width: 100%;
        height: 50px;
        background: #3f51b5;
        color: #fff;
        display: flex;
        justify-content: space-between;
        align-items: center;
        .header-left {
            .logo {
                height: 25px;
                margin-left: 30px;
            }
        }
        .header-right {
            display: flex;
            align-items: center;
            .r-item {
                margin-right: 25px;
                font-size: 0;
                .item {
                    color: #fff;
                    display: flex;
                    align-items: center;
                    height: 30px;
                    &:focus {
                        outline: none;
                    }
                    .icon-arrow-down {
                        margin-left: 3px;
                    }
                }
                .avatar {
                    width: 26px;
                    border-radius: 4px;
                }
            }
        }
    }
    .layer-container {
        margin-left: 200px;
        overflow: hidden;
        overflow-y: auto;
        height: 100%;
        .layer-aside {
            border-right: solid 1px #e6e6e6;
            position: fixed;
            left: 0;
            top: 50px;
            bottom: 0;
            width: 200px;
            .aside-menu {
                border-right: none;
                .iconfont {
                    &.left {
                        margin-right: 6px;
                    }
                }
            }
        }
        .layer-main {
            padding-top: 50px;
            height: 100%;
            background: #f0f2f5;
            box-sizing: border-box;
            .bread-crumb {
                background: #fff;
                padding: 16px;
            }
            .container {
                padding: 15px 20px;
            }
        }
    }
}
</style>
