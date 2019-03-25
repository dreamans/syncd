<template>
    <div class="app-release">
        <el-card shadow="never">
            <div slot="header" class="clearfix">
                <span>
                    部署发布单 - 
                    <span v-if="deployDetail.status == $root.ApplyStatusNone">待上线</span>
                    <span v-if="deployDetail.status == $root.ApplyStatusIng" class="app-color-info">{{ $t('onlineing') }}</span>
                    <span v-if="deployDetail.status == $root.ApplyStatusSuccess" class="app-color-success">{{ $t('have_onlined') }}</span>
                    <span v-if="deployDetail.status == $root.ApplyStatusFailed" class="app-color-error">{{ $t('online_failed') }}</span>
                    <span v-if="deployDetail.status == $root.ApplyStatusRollback" class="app-color-error">{{ $t('rollback') }}</span>
                </span>
            </div>
            <!-- apply -->
            <el-row :gutter="20">
                <el-col :span="10">
                    <span class="sp-title">{{ $t('apply_order') }}:</span>
                    <span>{{ this.applyDetail.name }}</span>
                </el-col>
                <el-col :span="14">
                    <span class="sp-title">{{ $t('belong_project') }}:</span>
                    <span>{{ this.projectDetail.name }}</span>
                </el-col>
            </el-row>
            <el-row class="app-mt-20" :gutter="20">
                <el-col :span="10">
                    <span class="sp-title">{{ $t('deploy_mode') }}:</span>
                    <span v-if="this.projectDetail.deploy_mode == $root.DeployModeBranch">
                        <i class="iconfont icon-branch"></i> {{ $t('branch_deploy') }} - {{ applyDetail.branch_name }} - commit:<template v-if="applyDetail.commit_version != ''">{{ applyDetail.commit_version }}</template><template v-else>HEAD</template>
                    </span>
                    <span v-else>
                        <i class="iconfont icon-tag"></i> {{ $t('tag_deploy') }} - {{ applyDetail.branch_name }}
                    </span>
                </el-col>
                <el-col :span="14">
                    <span class="sp-title">{{ $t('submiter') }}:</span>
                    <span>{{ applyDetail.username }} - {{ applyDetail.email }} - {{ $root.FormatDateTime(applyDetail.ctime) }}</span>
                </el-col>
            </el-row>
            <el-row class="app-mt-20" :gutter="20">
                <el-col :span="16">
                    <span class="sp-title">{{ $t('apply_type') }}:</span>
                    <span>
                        <span v-if="!applyDetail.is_rollback_apply">上线单</span>
                        <span class="app-color-warning" v-else>回滚单</span>
                    </span>
                </el-col>
            </el-row>
            <el-row class="app-mt-20" :gutter="20">
                <el-col :span="16">
                    <span class="sp-title">{{ $t('description') }}:</span>
                    <span>{{ applyDetail.description }}</span>
                </el-col>
            </el-row>
            <div class="app-divider"></div>
            <!-- build -->
            <div>
                <template v-if="(deployDetail.status == $root.ApplyStatusNone || deployDetail.status == $root.ApplyStatusIng || deployDetail.status == $root.ApplyStatusFailed) && !applyDetail.is_rollback_apply">
                    <el-button :disabled="deployDetail.status == $root.ApplyStatusIng" :loading="buildDetail.status == $root.ApplyStatusNone" @click="startBuildHandler" size="medium" icon="iconfont small left icon-build" type="primary">{{ $t('build') }}</el-button>
                    <el-button v-if="buildDetail.status == $root.ApplyStatusNone" @click="stopBuildHandler" size="medium" type="warning" icon="iconfont small left icon-stop">{{ $t('forced_termination') }}</el-button>
                </template>
                <el-row class="app-mt-20" :gutter="20">
                    <el-col :span="10">
                        <span class="sp-title">{{ $t('last_buid_time') }}:</span>
                        <span v-if="buildDetail.start_time">{{ $root.FormatDateTime(buildDetail.start_time) }}</span>
                        <span v-if="buildDetail.finish_time"> - {{ $t('cost_time') }}: {{ $root.FormatDateDuration((buildDetail.finish_time-buildDetail.start_time) * 1000) }}</span>
                    </el-col>
                    <el-col :span="14">
                        <span class="sp-title">{{ $t('status') }}:</span>
                        <span>
                            <span v-if="isStopBuildLoading && buildDetail.status == $root.BuildStatusStart" class="app-color-warning">
                                {{ $t('stopping') }}...
                            </span>
                            <span v-else-if="buildDetail.status == $root.BuildStatusStart" class="app-color-info">{{ $t('building') }}...</span>
                            <span v-else-if="buildDetail.status == $root.BuildStatusSuccess" class="app-color-success">{{ $t('build_finish') }}</span>
                            <span v-else-if="buildDetail.status == $root.BuildStatusFailed" class="app-color-error">{{ $t('build_failed') }}</span>
                            <span v-else>
                                {{ $t('unbuild') }}
                            </span>
                        </span>
                    </el-col>
                </el-row>
                <el-row class="app-mt-20" :gutter="20">
                    <el-col :span="10">
                        <span class="sp-title">{{ $t('build_log') }}:</span>
                        <span>
                            <span @click="openDialogBuildHandler" class="app-link">{{ $t('view') }}</span>
                        </span>
                    </el-col>
                    <el-col :span="14">
                        <span class="sp-title">{{ $t('tar_pack_path') }}:</span>
                        <span v-if="buildDetail.tar">{{ buildDetail.tar }}</span>
                        <span v-else>{{ $t('uncreate') }}</span>
                    </el-col>
                </el-row>
            </div>
            <div class="app-divider"></div>
            <!-- deploy -->
            <div>
                <el-button v-if="deployDetail.status == $root.ApplyStatusNone || deployDetail.status == $root.ApplyStatusIng" :loading="deployDetail.status == $root.ApplyStatusIng" @click="startDeployHandler" size="medium" icon="iconfont small left icon-send" type="primary">{{ $t('deploy') }}</el-button>
                <el-button v-if="deployDetail.status == $root.ApplyStatusFailed" :loading="deployDetail.status == $root.ApplyStatusIng" @click="startDeployHandler" size="medium" icon="iconfont small left icon-send" type="primary">{{ $t('redeploy') }}</el-button>
                <el-button v-if="deployDetail.status == $root.ApplyStatusIng" @click="stopDeployHandler" size="medium" icon="iconfont small left icon-stop" type="warning">{{ $t('forced_termination') }}</el-button>
                <template v-if="(deployDetail.status == $root.ApplyStatusSuccess || deployDetail.status == $root.ApplyStatusFailed) && applyDetail.rollback_id">
                    <el-button icon="iconfont small left icon-rollback" size="medium" @click="rollbackDeployHandler" type="danger">{{ $t('rollback') }}</el-button>
                    <el-alert class="app-mt-10" :title="$t('rollback_apply_order_tips')" type="warning"></el-alert>
                </template>
                <div v-if="deployDetail.status == $root.ApplyStatusRollback">
                    <el-alert class="app-mt-10" type="warning" show-icon :closable="false">
                        <template slot="title">
                            <strong>{{ $t('rollback_created') }}</strong>
                            -
                            <strong v-if="applyDetail.rollback_status == 1">{{ $t('rollback_unstart') }}</strong>
                            <strong v-else-if="applyDetail.rollback_status == 2">{{ $t('rollbacking') }}</strong>
                            <strong v-else-if="applyDetail.rollback_status == 3">{{ $t('rollback_success') }}</strong>
                            <strong v-else-if="applyDetail.rollback_status == 4">{{ $t('rollback_failed') }}</strong>
                            <strong v-else-if="applyDetail.rollback_status == 5">{{ $t('rollback_drop') }}</strong>
                            <strong v-else>{{ $t('rollback_unknow') }}</strong>
                            -
                            <span class="app-link" @click="$router.push({query: {id: applyDetail.rollback_apply_id}})">{{ $t('click_to_view_rollback_order') }}</span>
                        </template>
                    </el-alert>
                </div>
                <div v-if="applyDetail.is_rollback_apply">
                    <el-alert class="app-mt-10" type="warning" show-icon :closable="false">
                        <template slot="title">
                            {{ $t('rollback_tips') }}, <span class="app-link" @click="$router.push({query: {id: applyDetail.rollback_apply_id}})">{{ $t('back_apply') }}</span>
                        </template>
                    </el-alert>
                </div>
                <el-card shadow="never" class="app-mt-20 app-cluster-group">
                    <div slot="header">{{ $t('cluster_list') }}</div>
                    <el-collapse :value="projectDetail.online_cluster_ids">
                        <el-collapse-item class="app-cluster-item" v-for="c in onlineCluster" :name="c.id" :key="c.id">
                            <template slot="title">
                                <i class="iconfont small left icon-cluster"></i>
                                {{ c.name }}&nbsp;&nbsp;
                                <span v-if="deployDetail.groupStatus[c.id] == $root.DeployGroupStatusNone" class="app-color-gray"><i class="iconfont small left icon-wait"></i>{{ $t('wait_deploy') }}</span>
                                <span v-else-if="deployDetail.groupStatus[c.id] == $root.DeployGroupStatusStart" class="app-color-info"><i class="iconfont el-icon-loading"></i>{{ $t('deploying') }}</span>
                                <span v-else-if="deployDetail.groupStatus[c.id] == $root.DeployGroupStatusSuccess" class="app-color-success"><i class="iconfont small left icon-success"></i>{{ $t('deploy_success') }}</span>
                                <span v-else-if="deployDetail.groupStatus[c.id] == $root.DeployGroupStatusFailed" class="app-color-error"><i class="iconfont small left icon-failed"></i>{{ $t('deploy_failed') }}</span>
                            </template>
                            <div class="app-item" v-for="s in c.servers" :key="s.id">
                                <i class="iconfont small left icon-server"></i>
                                {{ s.ip }} - {{ s.name }}
                                <span v-if="deployDetail.servers[s.id] == undefined"></span>
                                <span v-else-if="deployDetail.servers[s.id].status == $root.DeployGroupStatusNone" class="app-color-gray"><i class="iconfont small left icon-wait"></i></span>
                                <span v-else-if="deployDetail.servers[s.id].status == $root.DeployGroupStatusStart" class="app-color-info"><i class="iconfont el-icon-loading"></i> <span @click="openDialogDeployHandler(s.id)" class="app-link">{{ $t('view') }}</span></span>
                                <span v-else-if="deployDetail.servers[s.id].status == $root.DeployGroupStatusSuccess" class="app-color-success"><i class="iconfont small left icon-success"></i> <span @click="openDialogDeployHandler(s.id)" class="app-link">{{ $t('view') }}</span></span>
                                <span v-else-if="deployDetail.servers[s.id].status == $root.DeployGroupStatusFailed" class="app-color-error"><i class="iconfont small left icon-failed"></i> <span @click="openDialogDeployHandler(s.id)" class="app-link">{{ $t('view') }}</span></span>
                            </div>
                        </el-collapse-item>
                    </el-collapse>
                </el-card>
            </div>
        </el-card>
        <el-dialog
        :width="$root.DialogNormalWidth"
        :title="$t('build_log')"
        :visible.sync="dialogBuildVisible"
        @close="closeDialogBuildHandler">
            <div v-if="buildDetail.status == $root.BuildStatusFailed"><i class="app-color-error el-icon-warning"></i> {{ $t('build_failed') }}: <span v-if="buildDetail.errmsg" class="app-color-error">{{ buildDetail.errmsg }}</span></div>
            <div v-if="buildDetail.status == $root.BuildStatusSuccess"><i class="app-color-success el-icon-success"></i> {{ $t('build_finish') }}</div>
            <div v-if="buildDetail.status == $root.BuildStatusStart"><i class="app-color-info el-icon-info"></i> {{ $t('build_ing') }}</div>
            <div class="app-terminal-log">
                <template v-for="(cmd, index) in buildDetail.output">
                    <div :key="index">
                        <div class="terminal-cmd" :class="{ 'app-color-success': cmd.success, 'app-color-error': !cmd.success}">
                            [cmd] $ {{ cmd.cmd }}
                            <span v-if="cmd.success" class="iconfont icon-right"></span>
                            <span v-else class="iconfont icon-wrong"></span>
                        </div>
                        <div><pre>{{ cmd.stdout }}</pre></div>
                        <div><pre>{{ cmd.stderr }}</pre></div>
                    </div>
                </template>
                <i v-if="buildDetail.status == $root.BuildStatusStart" class="el-icon-loading app-color-white"></i>
            </div>
        </el-dialog>

        <el-dialog
        :width="$root.DialogNormalWidth"
        :title="$t('deploy_log')"
        :visible.sync="dialogDeployVisible"
        @close="closeDialogDeployHandler">
            <div v-if="serverDeployDetail.status == $root.DeployGroupStatusFailed"><i class="app-color-error el-icon-warning"></i> {{ $t('deploy_failed') }}: <span v-if="serverDeployDetail.errmsg" class="app-color-error">{{ serverDeployDetail.errmsg }}</span></div>
            <div v-if="serverDeployDetail.status == $root.DeployGroupStatusSuccess"><i class="app-color-success el-icon-success"></i> {{ $t('deploy_success') }}</div>
            <div v-if="serverDeployDetail.status == $root.DeployGroupStatusStart"><i class="app-color-info el-icon-info"></i> {{ $t('deploying') }}</div>
            <div class="app-terminal-log">
                <template v-for="(cmd, index) in serverDeployDetail.output">
                    <div :key="index">
                        <div class="terminal-cmd" :class="{ 'app-color-success': cmd.success, 'app-color-error': !cmd.success}">
                            [cmd] $ {{ cmd.cmd }}
                            <span v-if="cmd.success" class="iconfont icon-right"></span>
                            <span v-else class="iconfont icon-wrong"></span>
                        </div>
                        <div><pre>{{ cmd.stdout }}</pre></div>
                        <div><pre>{{ cmd.stderr }}</pre></div>
                    </div>
                </template>
                <i v-if="serverDeployDetail.status == 1" class="el-icon-loading app-color-white"></i>
            </div>
        </el-dialog>
    </div>
</template>

<script>
import { 
    applyDetailApi, 
    applyProjectDetailApi, 
    buildStartApi, 
    buildStatusApi, 
    buildStopApi, 
    deployStart, 
    deployStatusApi, 
    deployStopApi, 
    deployRollbackApi } from '@/api/deploy'
import Code from '@/lib/code'
export default {
    data(){
        return {
            id: 0,
            applyDetail: {},
            projectDetail: {},
            buildDetail: {},
            deployDetail: {},
            dialogBuildVisible: false,
            isStopBuildLoading: false,
            onlineCluster: {},

            dialogDeployVisible: false,
            serverDeploySid: 0,
        }
    },
    watch: {
        projectDetail() {
            let clusters = []
            this.projectDetail.online_cluster_ids.forEach(id => {
                if (this.projectDetail.cluster_list[id]) {
                    let cluster = this.projectDetail.cluster_list[id]
                    let servers = []
                    this.projectDetail.server_list.forEach(s => {
                        if (s.group_id == id) {
                            servers.push(s)
                        }
                    })
                    clusters.push({
                        id: cluster.id,
                        name: cluster.name,
                        servers: servers,
                    })
                }
            })
            this.onlineCluster = clusters
        },
        '$route.query'() {
            this.initPageLoader()
        },
    },
    computed: {
        serverDeployDetail() {
            let detail = {}
            if (this.deployDetail.servers && this.deployDetail.servers[this.serverDeploySid]) {
                let srv = this.deployDetail.servers[this.serverDeploySid]
                detail = {
                    status: srv.status,
                    errmsg: srv.errmsg,
                    output: srv.task,
                }
            }
            return detail
        },
    },
    methods: {
        rollbackDeployHandler() {
            let vm = this
            this.$root.ConfirmDelete(function(){
                let msg = vm.$message({
                    message: vm.$t('rollback_creating'),
                    type: 'info',
                    iconClass: 'el-icon-loading',
                    duration: 0,
                })
                deployRollbackApi({id: vm.id}).then(res => {
                    setTimeout(() => {
                        vm.loadDeployStatus()
                        vm.initApplyDetail()
                        msg.close()
                    }, 1000)
                })
            }, this.$t('makesure_rollback_order'))
        },
        openDialogDeployHandler(sid) {
            this.serverDeploySid = sid
            this.dialogDeployVisible = true
        },
        closeDialogDeployHandler() {
            this.dialogDeployVisible = false
        },
        closeDialogBuildHandler() {
            this.dialogBuildVisible = false
        },
        openDialogBuildHandler() {
            this.dialogBuildVisible = true
        },
        startBuildHandler() {
            buildStartApi({id: this.id}).then(res => {
                this.loadBuildStatus()
            })
        },
        stopBuildHandler() {
            buildStopApi({id: this.id}).then(res => {
                this.isStopBuildLoading = true
            })
        },
        startDeployHandler() {
            deployStart({id: this.id}).then(res => {
                this.loadDeployStatus()
            }).catch(err => {
                if (err.code == Code.CODE_ERR_NO_DATA) {
                    this.$confirm(this.$t('build_success_and_deploy'), this.$t('deploy_tips'), {
                        confirmButtonText: this.$t('start_build'),
                        cancelButtonText: this.$t('i_known'),
                        type: 'warning',
                    }).then(() => {
                        this.startBuildHandler()
                    }).catch(() => { })
                }
            })
        },
        stopDeployHandler() {
            deployStopApi({id: this.id}).then(res => {
                this.loadDeployStatus()
            })
        },
        deployStatusDetail(detail) {
            let groupStatus = {}
            let servers = {}
            if (detail.task_list) {
                detail.task_list.forEach(task => {
                    groupStatus[task.group_id] = task.status
                    if (task.content) {
                        task.content.forEach(srv => {
                            servers[srv.id] = {
                                status: srv.status,
                                errmsg: srv.error,
                                task: srv.task ? srv.task : [],
                            }
                        })
                    }
                })
            }
            this.deployDetail = {
                status: detail.status,
                groupStatus: groupStatus,
                servers: servers,
            }
        },
        loadBuildStatus() {
            buildStatusApi({id: this.id}).then(res => {
                this.buildDetail = res
                if (res.status == 1) {
                    let vm = this
                    setTimeout(function() {
                        vm.loadBuildStatus()
                    }, 5000)
                } else {
                    this.isStopBuildLoading = false
                }
            })
        },
        loadDeployStatus() {
            deployStatusApi({id: this.id}).then(res => {
                this.deployStatusDetail(res)
                if (res.status == this.$root.ApplyStatusIng) {
                    let vm = this
                    setTimeout(function() {
                        vm.loadDeployStatus()
                    }, 5000)
                }
            })
        },
        initApplyDetail() {
            applyDetailApi({id: this.id}).then(res => {
                this.applyDetail = res
                applyProjectDetailApi({id: this.applyDetail.project_id}).then(res => {
                    this.projectDetail = res
                })
            })
        },
        initPageLoader() {
            this.id = this.$route.query.id
            this.initApplyDetail()
            this.loadBuildStatus()
            this.loadDeployStatus()
        },
    },
    mounted() {
        this.initPageLoader()
    },
}
</script>