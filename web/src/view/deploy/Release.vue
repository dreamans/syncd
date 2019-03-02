<template>
    <div class="app-release">
        <el-card shadow="never">
            <div slot="header" class="clearfix">
                <span>
                    部署发布单 - 
                    <span v-if="deployDetail.status == 1">待上线</span>
                    <span v-if="deployDetail.status == 2" class="app-color-info">{{ $t('onlineing') }}</span>
                    <span v-if="deployDetail.status == 3" class="app-color-success">{{ $t('have_onlined') }}</span>
                    <span v-if="deployDetail.status == 4" class="app-color-error">{{ $t('online_failed') }}</span>
                    <span v-if="deployDetail.status == 6" class="app-color-error">{{ $t('rollback') }}</span>
                </span>
            </div>
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
                    <span v-if="this.projectDetail.deploy_mode == 1">
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
                    <span class="sp-title">{{ $t('description') }}:</span>
                    <span>{{ applyDetail.description }}</span>
                </el-col>
            </el-row>
            <div class="app-divider"></div>
            <div>
                <template v-if="deployDetail.status == 1 || deployDetail.status == 2 || deployDetail.status == 4">
                    <el-button :loading="buildDetail.status == 1" @click="startBuildHandler" size="medium" icon="iconfont small left icon-build" type="primary">{{ $t('build') }}</el-button>
                    <el-button v-if="buildDetail.status == 1" @click="stopBuildHandler" size="medium" type="warning" icon="iconfont small left icon-stop">{{ $t('forced_termination') }}</el-button>
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
                            <span v-if="isStopBuildLoading && buildDetail.status == 1" class="app-color-warning">
                                {{ $t('stopping') }}...
                            </span>
                            <span v-else-if="buildDetail.status == 1" class="app-color-info">{{ $t('building') }}...</span>
                            <span v-else-if="buildDetail.status == 2" class="app-color-success">{{ $t('build_finish') }}</span>
                            <span v-else-if="buildDetail.status == 3" class="app-color-error">{{ $t('build_failed') }}</span>
                            <span v-else>
                                {{ $t('unbuild') }}
                            </span>
                        </span>
                    </el-col>
                </el-row>
                <el-row class="app-mt-20" :gutter="20">
                    <el-col :span="10">
                        <span class="sp-title">{{ $t('build_log') }}:</span>
                        <span v-if="buildDetail.status == 2 || buildDetail.status == 3">
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
            <div>
                <el-button v-if="deployDetail.status == 1 || deployDetail.status == 2" :loading="deployDetail.status == 2" @click="startDeployHandler" size="medium" icon="iconfont small left icon-send" type="primary">{{ $t('deploy') }}</el-button>
                <el-button v-if="deployDetail.status == 4" :loading="deployDetail.status == 2" @click="startDeployHandler" size="medium" icon="iconfont small left icon-send" type="primary">{{ $t('rebuild') }}</el-button>
                <el-button v-if="deployDetail.status == 2" @click="stopDeployHandler" size="medium" icon="iconfont small left icon-stop" type="warning">{{ $t('forced_termination') }}</el-button>
                <div v-if="(deployDetail.status == 3 || deployDetail.status == 4) && applyDetail.rollback_id">
                    <el-button icon="iconfont small left icon-rollback" size="medium" @click="rollbackDeployHandler" type="danger">{{ $t('rollback') }}</el-button>
                    <el-alert class="app-mt-10" :title="$t('rollback_apply_order_tips')" type="warning"></el-alert>
                </div>
                <div v-if="deployDetail.status == 6">
                    <el-alert type="warning" show-icon :closable="false">
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
                            <a class="app-link" :href="'/deploy/release?id=' + applyDetail.rollback_apply_id" target="_blank">{{ $t('click_to_view_rollback_order') }}</a>
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
                                <span v-if="deployDetail.groupStatus[c.id] == 0" class="app-color-gray"><i class="iconfont small left icon-wait"></i>{{ $t('wait_deploy') }}</span>
                                <span v-else-if="deployDetail.groupStatus[c.id] == 1" class="app-color-info"><i class="iconfont el-icon-loading"></i>{{ $t('deploying') }}</span>
                                <span v-else-if="deployDetail.groupStatus[c.id] == 2" class="app-color-success"><i class="iconfont small left icon-success"></i>{{ $t('deploy_success') }}</span>
                                <span v-else-if="deployDetail.groupStatus[c.id] == 3" class="app-color-error"><i class="iconfont small left icon-failed"></i>{{ $t('deploy_failed') }}</span>
                                <span v-else-if="deployDetail.groupStatus[c.id] == 4" class="app-color-warning"><i class="iconfont small left icon-stop"></i>{{ $t('be_deined') }}</span>
                            </template>
                            <div class="app-item" v-for="s in c.servers" :key="s.id">
                                <i class="iconfont small left icon-server"></i>
                                {{ s.ip }} - {{ s.name }}
                                <span v-if="deployDetail.servers[s.id] == undefined"></span>
                                <span v-else-if="deployDetail.servers[s.id].status == 0" class="app-color-gray"><i class="iconfont small left icon-wait"></i></span>
                                <span v-else-if="deployDetail.servers[s.id].status == 1" class="app-color-info"><i class="iconfont el-icon-loading"></i></span>
                                <span v-else-if="deployDetail.servers[s.id].status == 2" class="app-color-success"><i class="iconfont small left icon-success"></i> <span @click="openDialogDeployHandler(deployDetail.servers[s.id])" class="app-link">{{ $t('view') }}</span></span>
                                <span v-else-if="deployDetail.servers[s.id].status == 3" class="app-color-error"><i class="iconfont small left icon-failed"></i> <span @click="openDialogDeployHandler(deployDetail.servers[s.id])" class="app-link">{{ $t('view') }}</span></span>
                                <span v-else-if="deployDetail.servers[s.id].status == 4" class="app-color-warning"><i class="iconfont small left icon-stop"></i></span>
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
            <div v-if="buildDetail.status == 3"><i class="app-color-error el-icon-warning"></i> {{ $t('build_failed') }}: <span v-if="buildDetail.errmsg" class="app-color-error">{{ buildDetail.errmsg }}</span></div>
            <div v-if="buildDetail.status == 2"><i class="app-color-success el-icon-success"></i> {{ $t('build_finish') }}</div>
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
            </div>
        </el-dialog>

        <el-dialog
        :width="$root.DialogNormalWidth"
        :title="$t('deploy_log')"
        :visible.sync="dialogDeployVisible"
        @close="closeDialogDeployHandler">
            <div v-if="serverDeployDetail.status == 3"><i class="app-color-error el-icon-warning"></i> {{ $t('deploy_failed') }}: <span v-if="serverDeployDetail.errmsg" class="app-color-error">{{ serverDeployDetail.errmsg }}</span></div>
            <div v-if="serverDeployDetail.status == 2"><i class="app-color-success el-icon-success"></i> {{ $t('deploy_success') }}</div>
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
            </div>
        </el-dialog>
    </div>
</template>

<script>
import { applyDetailApi, applyProjectDetailApi, buildStartApi, buildStatusApi, buildStopApi, deployStart, deployStatusApi, deployStopApi, deployRollbackApi } from '@/api/deploy'
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
            serverDeployDetail: {},
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
    },
    methods: {
        rollbackDeployHandler() {
            let vm = this
            this.$root.ConfirmDelete(function(){
                deployRollbackApi({id: vm.id}).then(res => {
                    
                })
            }, this.$t('makesure_rollback_order'))
        },
        openDialogDeployHandler(srv) {
            console.log(srv)
            let arrOutput = []
            try {
                arrOutput = JSON.parse(srv.output)
            } catch(err) { }
            this.serverDeployDetail = {
                status: srv.status,
                errmsg: srv.errmsg,
                output: arrOutput,
            }
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

            })
        },
        deployStatusDetail(detail) {
            let groupStatus = {}
            let servers = {}
            if (detail.task_list) {
                let groupStatusList = {}
                let groupIds = []
                detail.task_list.forEach(srv => {
                    if (!groupStatusList[srv.group_id]) {
                        groupStatusList[srv.group_id] = []
                    }
                    groupIds.push(srv.group_id)
                    groupStatusList[srv.group_id].push(srv.status)
                    servers[srv.server_id] = srv
                })
                groupIds.forEach(gid => {
                    let none = groupStatusList[gid].indexOf(0) > -1
                    let ing = groupStatusList[gid].indexOf(1) > -1
                    let done = groupStatusList[gid].indexOf(2) > -1
                    let failed = groupStatusList[gid].indexOf(3) > -1
                    let terminate = groupStatusList[gid].indexOf(4) > -1
                    let status = 2 // 0 - unstart, 1 - starting, 2 - done , 3 - failed, 4 - terminate
                    if (none) {
                        if (ing || done || failed) {
                            status = 1
                        } else if (terminate) {
                            status = 4
                        } else {
                            status = 0
                        }
                    } else if (ing) {
                        status = 1
                    } else if (done) {
                        if (terminate) {
                            status = 4
                        } else if (failed) {
                            status = 3
                        } else {
                            status = 2
                        }
                    } else if (failed) {
                        if (terminate) {
                            status = 4
                        } else {
                            status = 3
                        }
                    } else {
                        status = 4
                    }
                    groupStatus[gid] = status
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
                if (res.status == 2) {
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
    },
    mounted() {
        this.id = this.$route.query.id
        this.initApplyDetail()
        this.loadBuildStatus()
        this.loadDeployStatus()
    },
}
</script>
