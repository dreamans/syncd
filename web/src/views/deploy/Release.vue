<template>
    <div>
        <a-card
        class="app-card app-release"
        title="部署上线单"
        :bordered="false">
            <a-row :gutter="20" class="app-mt-20">
                <a-col :span="8"><span class="sp-title">上线单:</span> {{ detail.name }}</a-col>
                <a-col :span="16"><span class="sp-title">所属项目:</span> {{ detail.space_name }} {{ detail.project_name }}</a-col>
            </a-row>
            <a-row :gutter="20" class="app-mt-20">
                <a-col :span="8"><span class="sp-title">仓库类型:</span> Git</a-col>
                <a-col :span="16" v-if="detail.repo_mode == 1"><span class="sp-title">上线分支:</span> {{ detail.repo_branch }} - {{ detail.repo_commit }}</a-col>
                <a-col :span="16" v-if="detail.repo_mode == 2"><span class="sp-title">上线Tag:</span> {{ detail.repo_tag }}</a-col>
            </a-row>
            <a-row :gutter="20" class="app-mt-20">
                <a-col :span="8"><span class="sp-title">提交者:</span> {{ detail.user_name }} - {{ detail.user_email }}</a-col>
                <a-col :span="16"><span class="sp-title">描述:</span> {{ detail.description }}</a-col>
            </a-row>
            <a-divider />
            <a-button v-if="applyStatus == 3" @click="handleStartDeploy" type="primary" size="large"><a-icon type="cloud-upload" /> 开始上线</a-button>
            <a-button v-else-if="applyStatus ==4" @click="handleStopDeploy" type="danger" size="large"><a-icon type="close-circle" /> 终止上线</a-button>
            <a-button v-else-if="applyStatus == 6 && stopingStatus" type="danger" size="large"><a-icon type="loading" /> 正在终止</a-button>
            <a-button v-else-if="applyStatus == 6" @click="handleStartDeploy" type="primary" size="large"><a-icon type="cloud-upload" /> 再次上线</a-button>

            <div class="shell">
                <div class="shell-title">
                    syncd@localhost:~$
                    <template>
                        <span class="app-color-error" v-if="stopingStatus">正在终止上线...</span>
                        <span class="app-color-info" v-else-if="applyStatus == 4">上线中</span>
                        <span class="app-color-success" v-else-if="applyStatus == 5">上线成功</span>
                        <span class="app-color-error" v-else-if="applyStatus == 6">上线失败</span>
                        <span v-else>等待上线...</span>
                    </template>
                </div>
                <div v-if="this.errorLog" class="shell-body app-color-error">
                    Error Output >>>
                    <pre>{{this.errorLog}}</pre>
                </div>
                <template v-for="d in deployList">
                    <div class="shell-sub-title">
                        <div>>>> {{ d.level == 4 ? '部署到> ' :''}}{{ $root.T(d.name) }}
                            <a-icon v-if="d.status == 2" type="loading-3-quarters" spin/>
                            <a-icon v-if="d.status == 3" type="check"/>
                            <a-icon v-if="d.status == 4 || d.status == 5" type="warning"/>
                        </div>
                    </div>
                    <div class="shell-body">
                        <div v-if="d.status == 1 && stopingStatus">
                            正在拦截并终止...
                        </div>
                        <div v-else-if="d.status == 1">
                            未开始
                        </div>
                        <div v-else-if="d.status == 2">
                            <span class="app-color-info">执行中</span>
                        </div>
                        <div v-else-if="d.status == 2">
                            <span class="app-color-info">执行中</span>
                        </div>
                        <div v-else-if="d.status == 3">
                            <span class="app-color-success">完成</span>
                        </div>
                        <div v-else-if="d.status == 4">
                            <span class="app-color-error">失败</span>
                        </div>
                        <div v-else-if="d.status == 5">
                            <span class="app-color-error">终止</span>
                        </div>
                        <pre>{{d.output}}</pre>
                    </div>
                </template>
            </div>
        </a-card>
    </div>
</template>

<script>
import { getApplyDetailApi, startDeployApi, statusDeployApi, stopDeployApi } from '@/api/deploy.js'
export default {
    data() {
        return {
            id: 0,
            detail: {},
            loopStatus: null,
            applyStatus: 0,
            errorLog: '',
            deployList: [],
        }
    },
    computed: {
        stopingStatus() {
            let isDeploying = false
            if (this.deployList) {
                this.deployList.forEach(item => {
                    if (item.status == 2) {
                        isDeploying = true
                    }
                })
            }
            if (this.applyStatus == 6 && isDeploying) {
                return true
            }
            return false
        },
    },
    methods: {
        handleStartDeploy() {
            startDeployApi({id: this.id}).then(res => {
                this.loopDeployStatus()
            })
        },
        handleStopDeploy() {
            stopDeployApi({id: this.id}).then(res => {
                this.getApplyStatus()
            })
        },
        loopDeployStatus() {
            let vm = this
            vm.loopStatus = setInterval(function() {
                vm.getApplyStatus()
            }, 3000)
            this.getApplyStatus()
        },
        getApplyStatus() {
            statusDeployApi({id: this.id}).then(res => {
                this.applyStatus = res.apply_status
                this.deployList = res.deploy_list
                this.errorLog = res.error_log
                if (this.applyStatus != 4 && !(this.stopingStatus && this.applyStatus == 6)) {
                    clearInterval(this.loopStatus)
                }
            })
        },
        loadApplyDetail() {
            getApplyDetailApi({id: this.id}).then(res => {
                this.detail = res
            })
        },
    },
    mounted() {
        this.id = this.$route.query.id
        if (!this.id) {
            this.$root.GotoRouter('deployDeploy')
        }
        this.loadApplyDetail()
        this.loopDeployStatus()
    },
}
</script>
