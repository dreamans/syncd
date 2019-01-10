<template>
    <div class="app-dashboard">
        <div class="dashboard-container">
            <a-card
            :bordered="false">
                <div class="meta">
                    <a-avatar :size="32" :src="$store.getters['account/getAvatar']">{{ $store.getters['account/getAvatar'] }}</a-avatar>
                    <div class="content">早安，{{$store.getters['account/getUserName']}}，祝你开心每一天！</div>
                </div>
                <a-divider />
                <div class="quick-start">
                    <a-row :gutter="20">
                        <a-col :span="4" v-if="$root.CheckPriv($root.Priv.DEPLOY_APPLY)">
                            <span class="quick-item" @click="$root.GotoRouter('deployApply')">
                                <a-avatar class="item-icon" style="color: #fff; backgroundColor: #f56a00" size="small" icon="edit" />提交上线申请</span>
                            </span>
                        </a-col>
                        <a-col :span="4" v-if="$root.CheckPrivs([$root.Priv.DEPLOY_VIEW_MY, $root.Priv.DEPLOY_VIEW_ALL])">
                            <span class="quick-item" @click="$root.GotoRouter('deployDeploy')">
                                <a-avatar class="item-icon" style="color: #fff; backgroundColor: #7265e6" size="small" icon="edit" />上线单管理</span>
                            </span>
                        </a-col>
                        <a-col :span="4" v-if="$root.CheckPriv($root.Priv.PROJECT_SPACE_VIEW)">
                            <span class="quick-item" @click="$root.GotoRouter('projectSpace')">
                                <a-avatar class="item-icon" style="color: #fff; backgroundColor: #ffbf00" size="small" icon="edit" />项目管理</span>
                            </span>
                        </a-col>
                        <a-col :span="4" v-if="$root.CheckPriv($root.Priv.SERVER_GROUP_VIEW)">
                            <span class="quick-item" @click="$root.GotoRouter('serverGroup')">
                                <a-avatar class="item-icon" style="color: #fff; backgroundColor: #00a2ae" size="small" icon="edit" />集群管理</span>
                            </span>
                        </a-col>
                        <a-col :span="4" v-if="$root.CheckPriv($root.Priv.USER_VIEW)">
                            <span class="quick-item" @click="$root.GotoRouter('userList')">
                                <a-avatar class="item-icon" style="color: #fff; backgroundColor: #268bd2" size="small" icon="edit" />用户管理</span>
                            </span>
                        </a-col>
                    </a-row>
                </div>
                <a-divider />
                <div class="go">
                    <img class="logo" src="../assets/syncd.png" />
                </div>
            </a-card>
        </div>
    </div>
</template>

<script>
export default {

}
</script>
