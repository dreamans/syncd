<template>
    <div class="app-release">
        <el-card shadow="never">
            <div slot="header" class="clearfix">
                <span>
                    部署发布单 - 
                    <span v-if="deployDetail.status == 1">待上线</span>
                    <span v-if="deployDetail.status == 2" class="app-color-info">上线中</span>
                    <span v-if="deployDetail.status == 3" class="app-color-success">已上线</span>
                    <span v-if="deployDetail.status == 4" class="app-color-error">上线失败</span>
                </span>
            </div>
            <el-row :gutter="20">
                <el-col :span="10">
                    <span class="sp-title">发布单:</span>
                    <span>{{ this.applyDetail.name }}</span>
                </el-col>
                <el-col :span="14">
                    <span class="sp-title">所属项目:</span>
                    <span>{{ this.projectDetail.name }}</span>
                </el-col>
            </el-row>
            <el-row class="app-mt-20" :gutter="20">
                <el-col :span="10">
                    <span class="sp-title">上线模式:</span>
                    <span v-if="this.projectDetail.deploy_mode == 1">
                        <i class="iconfont icon-branch"></i> 分支上线 - {{ applyDetail.branch_name }} - commit:<template v-if="applyDetail.commit_version != ''">{{ applyDetail.commit_version }}</template><template v-else>HEAD</template>
                    </span>
                    <span v-else>
                        <i class="iconfont icon-tag"></i> TAG上线 - {{ applyDetail.branch_name }}
                    </span>
                </el-col>
                <el-col :span="14">
                    <span class="sp-title">提交者:</span>
                    <span>{{ applyDetail.username }} - {{ applyDetail.email }} - {{ $root.FormatDateTime(applyDetail.ctime) }}</span>
                </el-col>
            </el-row>
            <el-row class="app-mt-20" :gutter="20">
                <el-col :span="16">
                    <span class="sp-title">描述:</span>
                    <span>{{ applyDetail.description }}</span>
                </el-col>
            </el-row>
            <div class="app-divider"></div>
            <div>
                <template v-if="deployDetail.status == 1 || deployDetail.status == 2 || deployDetail.status == 4">
                    <el-button :loading="buildDetail.status == 1" @click="startBuildHandler" size="medium" icon="iconfont small left icon-build" type="primary">构建</el-button>
                    <el-button v-if="buildDetail.status == 1" @click="stopBuildHandler" size="medium" type="warning" icon="iconfont small left icon-stop">强制终止</el-button>
                </template>
                <el-row class="app-mt-20" :gutter="20">
                    <el-col :span="10">
                        <span class="sp-title">上次构建时间:</span>
                        <span v-if="buildDetail.start_time">{{ $root.FormatDateTime(buildDetail.start_time) }}</span>
                        <span v-if="buildDetail.finish_time"> - 耗时: {{ $root.FormatDateDuration((buildDetail.finish_time-buildDetail.start_time) * 1000) }}</span>
                    </el-col>
                    <el-col :span="14">
                        <span class="sp-title">状态:</span>
                        <span>
                            <span v-if="isStopBuildLoading && buildDetail.status == 1" class="app-color-warning">
                                正在终止...
                            </span>
                            <span v-else-if="buildDetail.status == 1" class="app-color-info">构建中...</span>
                            <span v-else-if="buildDetail.status == 2" class="app-color-success">构建成功</span>
                            <span v-else-if="buildDetail.status == 3" class="app-color-error">构建失败</span>
                            <span v-else>
                                未构建
                            </span>
                        </span>
                    </el-col>
                </el-row>
                <el-row class="app-mt-20" :gutter="20">
                    <el-col :span="10">
                        <span class="sp-title">构建日志:</span>
                        <span v-if="buildDetail.status == 2 || buildDetail.status == 3">
                            <span @click="openDialogBuildHandler" class="app-link">查看</span>
                        </span>
                    </el-col>
                    <el-col :span="14">
                        <span class="sp-title">Tar包位置:</span>
                        <span v-if="buildDetail.tar">{{ buildDetail.tar }}</span>
                        <span v-else>未生成</span>
                    </el-col>
                </el-row>
            </div>
            <div class="app-divider"></div>
            <div>
                <template v-if="deployDetail.status == 1 || deployDetail.status == 2 || deployDetail.status == 4">
                    <el-button v-if="deployDetail.status == 1 || deployDetail.status == 2" :loading="deployDetail.status == 2" @click="startDeployHandler" size="small" icon="iconfont small left icon-send" type="primary">部署</el-button>
                    <el-button v-if="deployDetail.status == 4" :loading="deployDetail.status == 2" @click="startDeployHandler" size="small" icon="iconfont small left icon-send" type="primary">重新部署</el-button>
                </template>
                <el-card shadow="never" class="app-mt-20 app-cluster-group">
                    <div slot="header">集群列表</div>
                    <el-collapse>
                        <el-collapse-item class="app-cluster-item" v-for="c in onlineCluster" :name="c.id" :key="c.id">
                            <template slot="title">
                                <i class="iconfont small left icon-cluster"></i>
                                {{ c.name }}&nbsp;&nbsp;
                                <span v-if="deployDetail.groupStatus[c.id] == 0" class="app-color-gray"><i class="iconfont small left icon-wait"></i>等待部署</span>
                                <span v-else-if="deployDetail.groupStatus[c.id] == 1" class="app-color-info"><i class="iconfont el-icon-loading"></i>部署中</span>
                                <span v-else-if="deployDetail.groupStatus[c.id] == 2" class="app-color-success"><i class="iconfont small left icon-success"></i>部署成功</span>
                                <span v-else-if="deployDetail.groupStatus[c.id] == 3" class="app-color-error"><i class="iconfont small left icon-failed"></i>部署失败</span>
                                <span v-else class="app-color-warning"><i class="iconfont small left icon-question"></i>未知集群</span>
                            </template>
                            <div class="app-item" v-for="s in c.servers" :key="s.id">
                                <i class="iconfont small left icon-server"></i>
                                {{ s.ip }} - {{ s.name }}
                                <span v-if="deployDetail.servers[s.id] == undefined" class="app-color-warning"><i class="app-color-gray iconfont small left icon-question"></i></span>
                                <span v-else-if="deployDetail.servers[s.id].status == 0" class="app-color-gray"><i class="iconfont small left icon-wait"></i></span>
                                <span v-else-if="deployDetail.servers[s.id].status == 1" class="app-color-info"><i class="iconfont el-icon-loading"></i></span>
                                <span v-else-if="deployDetail.servers[s.id].status == 2" class="app-color-success"><i class="iconfont small left icon-success"></i> <span @click="openDialogBuildHandler" class="app-link">查看</span></span>
                                <span v-else-if="deployDetail.servers[s.id].status == 3" class="app-color-error"><i class="iconfont small left icon-failed"></i> <span @click="openDialogBuildHandler" class="app-link">查看</span></span>
                            </div>
                        </el-collapse-item>
                    </el-collapse>
                </el-card>
            </div>
        </el-card>
        <el-dialog
        :width="$root.DialogNormalWidth"
        title="构建日志"
        :visible.sync="dialogBuildVisible"
        @close="closeDialogBuildHandler">
            <div v-if="buildDetail.status == 3"><i class="app-color-error el-icon-warning"></i> 构建失败: <span v-if="buildDetail.errmsg" class="app-color-error">{{ buildDetail.errmsg }}</span></div>
            <div v-if="buildDetail.status == 2"><i class="app-color-success el-icon-success"></i> 构建成功</div>
            <div class="app-terminal-log">
                <template v-for="(cmd, index) in buildDetail.output">
                    <div :key="index">
                        <div class="terminal-cmd" :class="{ 'app-color-success': cmd.success, 'app-color-error': !cmd.success}">
                            [cmd ] $ {{ cmd.cmd }}
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
import { applyDetailApi, applyProjectDetailApi, buildStartApi, buildStatusApi, buildStopApi, deployStart, deployStatusApi } from '@/api/deploy'
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
                    let status = 2 // 0 - unstart, 1 - starting, 2 - done , 3 - failed
                    if (none) {
                        if (ing || done || failed) {
                            status = 1
                        } else {
                            status = 0
                        }
                    } else if (ing) {
                        status = 1
                    } else {
                        if (failed) {
                            status = 3
                        } else {
                            status = 2
                        }
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
