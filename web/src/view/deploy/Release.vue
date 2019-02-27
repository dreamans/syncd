<template>
    <div class="app-release">
        <el-card shadow="never">
            <div slot="header" class="clearfix">
                <span>部署发布单</span>
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
                <el-button :loading="buildDetail.status == 1" @click="startBuildHandler" size="medium" icon="iconfont small left icon-build" type="primary">构建</el-button>
                <el-button v-if="buildDetail.status == 1" @click="stopBuildHandler" size="medium" type="warning" icon="iconfont small left icon-stop">强制终止</el-button>
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
            <template v-if="preCluster.id">
            <div class="app-divider"></div>
            <div>
                <el-button size="small" icon="iconfont small left icon-send" type="primary">部署-预发布环境</el-button>
                <el-card shadow="never" class="app-mt-20 app-cluster-group">
                    <el-collapse>
                        <el-collapse-item :name="preCluster.id">
                            <template slot="title">
                                <i class="iconfont small left icon-cluster"></i>
                                {{ preCluster.name }}
                                <i class="app-color-success iconfont small left icon-success"></i>
                            </template>
                            <div class="app-item" v-for="s in preCluster.servers" :key="s.id">
                                <i class="iconfont small left icon-wait"></i>
                                {{ s.ip }} - {{ s.name }} <span @click="openDialogBuildHandler" class="app-link">查看</span>
                            </div>
                            <!--div class="app-item app-color-info">
                                <i class="el-icon-loading"></i>
                                1.2.3.4 - local.xyz.001 <span @click="openDialogBuildHandler" class="app-link">查看</span>
                            </div>
                            <div class="app-item app-color-error">
                                <i class="iconfont small left icon-failed"></i>
                                1.2.3.4 - local.xyz.001 <span @click="openDialogBuildHandler" class="app-link">查看</span>
                            </div>
                            <div class="app-item app-color-success">
                                <i class="iconfont small left icon-success"></i>
                                1.2.3.4 - local.xyz.001 <span @click="openDialogBuildHandler" class="app-link">查看</span>
                            </div-->
                        </el-collapse-item>
                    </el-collapse>
                </el-card>
            </div>
            </template>
            <div class="app-divider"></div>
            <div>
                <el-button size="small" icon="iconfont small left icon-send" type="primary">部署-生产环境</el-button>
                <el-card shadow="never" class="app-mt-20 app-cluster-group">
                    <el-collapse>
                        <el-collapse-item v-for="c in onlineCluster" :name="c.id" :key="c.id">
                            <template slot="title">
                                <i class="iconfont small left icon-cluster"></i>
                                {{ c.name }}
                                <i class="app-color-success iconfont small left icon-success"></i>
                            </template>
                            <div class="app-item" v-for="s in c.servers" :key="s.id">
                                <i class="iconfont small left icon-wait"></i>
                                {{ s.ip }} - {{ s.name }} <span @click="openDialogBuildHandler" class="app-link">查看</span>
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
import { applyDetailApi, applyProjectDetailApi, buildStartApi, buildStatusApi, buildStopApi } from '@/api/deploy'
export default {
    data(){
        return {
            id: 0,
            applyDetail: {},
            projectDetail: {},
            buildDetail: {},
            dialogBuildVisible: false,
            isStopBuildLoading: false,

            preCluster: {},
            onlineCluster: {},
        }
    },
    watch: {
        projectDetail() {
            let preCluster = {}
            if ( this.projectDetail.pre_cluster_id 
                && this.projectDetail.cluster_list[this.projectDetail.pre_cluster_id]) {
                let cluster = this.projectDetail.cluster_list[this.projectDetail.pre_cluster_id]
                preCluster.id = cluster.id
                preCluster.name = cluster.name
                preCluster.servers = []
                if (this.projectDetail.server_list[preCluster.id]) {
                    preCluster.servers = this.projectDetail.server_list[preCluster.id]
                    this.preCluster = preCluster
                }
            }

            let clusters = []
            this.projectDetail.online_cluster_ids.forEach(id => {
                if (this.projectDetail.cluster_list[id]) {
                    let cluster = this.projectDetail.cluster_list[id]
                    let servers = []
                    if (this.projectDetail.server_list[id]) {
                        servers = this.projectDetail.server_list[id]
                    }
                    if (servers.length) {
                        clusters.push({
                            id: cluster.id,
                            name: cluster.name,
                            servers: servers,
                        })
                    }
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
        loadBuildStatus() {
            buildStatusApi({id: this.id}).then(res => {
                this.buildDetail = res
                if (res.status == 1) {
                    let vm = this
                    setTimeout(function() {
                        vm.loadBuildStatus()
                    }, 3000)
                } else {
                    this.isStopBuildLoading = false
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
    },
}
</script>
