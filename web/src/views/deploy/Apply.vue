<template>
    <div>
        <a-card
        class="app-card"
        title="提交上线申请"
        :bordered="false">
            <div class="app-btn-group">
                <a-row :gutter="10">
                    <a-col :span="8">
                        <a-select
                        showSearch
                        placeholder="请选择项目空间"
                        v-model="selectedSpaceId"
                        :notFoundContent="fetchingSpace ? undefined : '未找到数据'"
                        style="width: 100%">
                            <a-spin v-if="fetchingSpace" slot="notFoundContent" size="small"/>
                            <a-select-option v-for="s in spaceList" :key="`${s.id}`">{{ s.name }}</a-select-option>
                        </a-select>
                    </a-col>
                    <a-col :span="8" v-if="selectedSpaceId">
                        <a-select
                        showSearch
                        placeholder="请选择项目"
                        v-model="selectedProjectId"
                        :notFoundContent="fetchingProject ? undefined : '未找到数据'"
                        style="width: 100%">
                            <a-spin v-if="fetchingProject" slot="notFoundContent" size="small"/>
                            <a-select-option v-for="p in projectList" :key="`${p.id}`">{{ p.name }}</a-select-option>
                        </a-select>
                    </a-col>
                    <a-col :span="8" v-if="selectedProjectId">
                        <a-button @click="handleOpenDialog" type="primary"><a-icon type="file" /> 填写上线单</a-button>
                    </a-col>
                </a-row>
            </div>
        </a-card>

        <a-modal
        title="填写上线申请"
        :visible="dialogVisible"
        @ok="dialogSubmit"
        :confirmLoading="dialogConfirmLoading"
        :keyboard="false"
        :maskClosable="false"
        okText="提交申请"
        cancelText="放弃"
        :destroyOnClose="true"
        width="50%"
        @cancel="dialogCancel">
            <apply-update-component ref="updateRef" :detail="{project_id: this.projectId}"></apply-update-component>
        </a-modal>

    </div>
</template>

<script>
import { getApplySpaceList, submitApplyApi, getApplyProjectList } from '@/api/deploy.js'
import ApplyUpdateComponent from './ApplyUpdateComponent.js'
export default {
    data() {
        return {
            fetchingSpace: false,
            spaceList: [],
            selectedSpaceId: undefined,
            fetchingProject: false,
            projectList: [],
            selectedProjectId: undefined,
            dialogVisible: false,
            dialogConfirmLoading: false,
        }
    },
    components: {
        ApplyUpdateComponent,
    },
    watch: {
        selectedSpaceId(spaceId) {
            this.loadProjectList(spaceId)
            this.projectList = []
            this.selectedProjectId = undefined
        },
    },
    computed: {
        projectId() {
            let spaceId = parseInt(this.selectedProjectId)
            if (isNaN(spaceId)) {
                return 0
            }
            return spaceId
        },
    },
    methods: {
        handleOpenDialog() {
            this.dialogVisible = true
        },
        dialogSubmit() {
            this.$refs.updateRef.validateFields((err, values) => {
                if (err) {
                    return
                }
                this.dialogConfirmLoading = true
                submitApplyApi(values).then(res => {
                    this.dialogConfirmLoading = false
                    this.$success({
                        title: '提交成功',
                        okText: "确定",
                        content: (
                            <div>恭喜，上线申请提交成功</div>
                        ),
                        onOk: () => {
                            if (this.$root.CheckPrivs([this.$root.Priv.DEPLOY_VIEW_MY, this.$root.Priv.DEPLOY_VIEW_ALL])) {
                                this.$root.GotoRouter('deployDeploy')
                            }
                            this.dialogCancel()
                        }
                    });
                }).catch(err => {
                    this.dialogConfirmLoading = false
                })
            })
        },
        dialogCancel() {
            this.dialogVisible = false
        },
        loadProjectList(spaceId) {
            this.fetchingProject = true
            getApplyProjectList({space_id: spaceId, status: 1}).then(res => {
                this.fetchingProject = false
                this.projectList = res.list
                if (!res.list) {
                    this.$message.warning('该空间下无可用项目')
                }
            }).catch(err => {
                this.fetchingProject = false
            })
        },
        loadSpaceList() {
            this.fetchingSpace = true
            getApplySpaceList().then(res => {
                this.fetchingSpace = false
                this.spaceList = res.list
            }).catch(err => {
                this.fetchingSpace = false
            })
        },
    },
    mounted() {
        this.loadSpaceList()
    },
}
</script>
