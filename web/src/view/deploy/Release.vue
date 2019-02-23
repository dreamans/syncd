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
                <el-button @click="startBuildHandler" size="small" icon="iconfont small left icon-build" type="primary">构建</el-button>
                <el-row class="app-mt-20" :gutter="20">
                    <el-col :span="10">
                        <span class="sp-title">上次构建时间:</span>
                        <span>2019-02-12 11:30:50</span>
                    </el-col>
                    <el-col :span="14">
                        <span class="sp-title">状态:</span>
                        <span>
                            未构建
                            <span class="app-color-info">构建中</span>
                            <span class="app-color-success">构建成功</span>
                            <span class="app-color-error">构建失败</span>
                        </span>
                    </el-col>
                </el-row>
                <el-row class="app-mt-20" :gutter="20">
                    <el-col :span="14">
                        <span class="sp-title">Tar包位置:</span>
                        <span>/usr/local/var/syncd_build/74b4358ed86b2d873.tar.gz</span>
                    </el-col>
                </el-row>
            </div>
            <div class="app-divider"></div>
            <div>
                <el-button size="small" icon="iconfont small left icon-send" type="primary">部署</el-button>
                <el-card shadow="never" class="app-mt-20">
                    <el-collapse>
                        <el-collapse-item name="3">
                            <template slot="title">
                                <i class="iconfont small left icon-cluster"></i>
                                阿里云集群.huabei.01
                                <i class="app-color-success iconfont small left icon-success"></i>
                            </template>
                            <div class="app-item">
                                <i class="iconfont small left icon-wait"></i>
                                1.2.3.4 - local.xyz.001 <span @click="openDialogDeployHandler" class="app-link">查看</span>
                            </div>
                            <div class="app-item app-color-info">
                                <i class="el-icon-loading"></i>
                                1.2.3.4 - local.xyz.001 <span @click="openDialogDeployHandler" class="app-link">查看</span>
                            </div>
                            <div class="app-item app-color-error">
                                <i class="iconfont small left icon-failed"></i>
                                1.2.3.4 - local.xyz.001 <span @click="openDialogDeployHandler" class="app-link">查看</span>
                            </div>
                            <div class="app-item app-color-success">
                                <i class="iconfont small left icon-success"></i>
                                1.2.3.4 - local.xyz.001 <span @click="openDialogDeployHandler" class="app-link">查看</span>
                            </div>
                        </el-collapse-item>
                    </el-collapse>
                </el-card>
            </div>
        </el-card>
        <el-dialog
        :width="$root.DialogNormalWidth"
        title="部署详情"
        :visible.sync="dialogDeployVisible"
        @close="closeDialogDeployHandler">
            <div>
                服务器: 1.2.3.4 - local.xyz.001
            </div>
            <div class="app-terminal-log">
<pre>
<span class="terminal-cmd">$ ./scp a git@hello</span>
syncd@localhost:~$ 上线成功
>>> 更新代码仓库文件 
完成
>>> 拉取代码文件 
完成
>>> 打包 
完成
>>> 部署到> asdf-www.yizhuanlan.com 
完成
/usr/local/nginx/sbin/nginx -s stop

/usr/local/nginx/sbin/nginx
</pre>
                </div>
        </el-dialog>
    </div>
</template>

<script>
import { applyDetailApi, applyProjectDetailApi, buildStartApi } from '@/api/deploy'
export default {
    data(){
        return {
            id: 0,
            applyDetail: {},
            projectDetail: {},
            dialogDeployVisible: false,
        }
    },
    methods: {
        closeDialogDeployHandler() {
            this.dialogDeployVisible = false
        },
        openDialogDeployHandler() {
            this.dialogDeployVisible = true
        },
        startBuildHandler() {
            buildStartApi({id: this.id}).then(res => {
                
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
    },
}
</script>
