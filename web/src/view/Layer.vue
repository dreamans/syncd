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
                                <el-dropdown-item><i class="iconfont small left icon-help"></i>帮助</el-dropdown-item>
                            </a>
                            <a class="app-dropdown-link" href="https://github.com/dreamans/syncd" target="_blank">
                                <el-dropdown-item><i class="iconfont small left icon-pull-request"></i>为 Syncd 做贡献</el-dropdown-item>
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
                            <el-dropdown-item divided><i class="iconfont small left icon-setting"></i>个人设置</el-dropdown-item>
                            <el-dropdown-item><i class="iconfont small left icon-key"></i>修改密码</el-dropdown-item>
                            <el-dropdown-item divided><i class="iconfont small left icon-logout"></i>退出登录</el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                </span>
            </div>
        </header>
        <section class="layer-container">
            <aside class="layer-aside">
                <ScrollBar>
                    <el-menu >
                        <el-submenu index="1">
                            <template slot="title">
                                <i class="el-icon-location"></i><span>导航一</span>
                            </template>
                            <el-menu-item index="1-1">选项1</el-menu-item>
                        </el-submenu>
                        <el-menu-item index="2">
                            <i class="el-icon-menu"></i>
                            <span slot="title">导航二</span>
                        </el-menu-item>
                        <el-menu-item index="3" disabled>
                            <i class="el-icon-document"></i>
                            <span slot="title">导航三</span>
                        </el-menu-item>
                        <el-menu-item index="4">
                            <i class="el-icon-setting"></i>
                            <span slot="title">导航四</span>
                        </el-menu-item>
                    </el-menu>
                </ScrollBar>
            </aside>
            <main class="layer-main">
               
            </main>
        </section>
    </div>
</template>

<script>
import ScrollBar from '@/component/ScrollBar';
export default {
    components: {
        ScrollBar,
    },
}
</script>


<style lang="scss" scoped>
.layer-global {
    .layer-header {
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
        margin-left: 180px;
        overflow: hidden;
        overflow-y: auto;
        .layer-aside {
            background: #0000cc;
            position: fixed;
            left: 0;
            top: 50px;
            bottom: 0;
            width: 180px;
        }
        .layer-main {
            margin-top: 50px;
        }
    }
}
</style>
