<template>
    <div>
        <el-card shadow="never">
            <el-row class="app-btn-group" :gutter="10">
                <el-col :span="3">
                    <el-select
                    :placeholder="$t('submit_time')"
                    size="medium" 
                    clearable style="width: 100%" 
                    v-model="searchTime">
                        <el-option
                            v-for="s in timeList"
                            :key="s.time"
                            :label="s.label"
                            :value="s.time">
                        </el-option>
                    </el-select>
                </el-col>
                <el-col :span="3">
                    <el-select
                    :placeholder="$t('audit_status')"
                    size="medium" 
                    clearable 
                    style="width: 100%" 
                    v-model="searchAuditStatus">
                        <el-option
                            v-for="s in auditStatusList"
                            :key="s.status"
                            :label="s.label"
                            :value="s.status">
                        </el-option>
                    </el-select>
                </el-col>
                <el-col :span="3">
                    <el-select
                    :placeholder="$t('deploy_status')"
                    size="medium" 
                    clearable 
                    style="width: 100%" 
                    v-model="searchStatus">
                        <el-option
                            v-for="s in statusList"
                            :key="s.status"
                            :label="s.label"
                            :value="s.status">
                        </el-option>
                    </el-select>
                </el-col>
                <el-col :span="9">
                    <el-select
                    :placeholder="$t('select_project')"
                    size="medium" 
                    clearable 
                    style="width: 100%" 
                    v-model="searchProjectId">
                        <el-option
                        v-for="p in projectList"
                        :key="p.project_id"
                        :label="p.project_name"
                        :value="p.project_id">
                            {{p.space_name}} <i class="el-icon-arrow-right"></i> {{p.project_name}}
                        </el-option>
                    </el-select>
                </el-col>
                <el-col :span="6">
                    <el-input @keyup.enter.native="searchHandler" v-model="searchInput" size="medium" :placeholder="$t('please_input_keyword_id_or_name')">
                        <el-button @click="searchHandler" slot="append" icon="el-icon-search"></el-button>
                    </el-input>
                </el-col>
            </el-row>
            <el-table
                class="app-table"
                size="medium"
                v-loading="tableLoading"
                :data="tableData">
                <el-table-column prop="id" label="ID" width="80"></el-table-column>
                <el-table-column prop="name" :label="$t('name')"></el-table-column>
                <el-table-column label="项目名称">
                    <template slot-scope="scope">
                         {{ scope.row.project_name }}
                         <el-tooltip effect="dark" :content="'所属空间: ' + scope.row.space_name" placement="top">
                            <span class="app-cursor"><i class="iconfont icon-space"></i></span>
                        </el-tooltip>
                    </template>
                </el-table-column>
                <el-table-column prop="ssh_port" width="100" label="提交时间">
                    <template slot-scope="scope">
                        <el-tooltip effect="dark" :content="$root.FormatDateTime(scope.row.ctime)" placement="top">
                            <span class="app-cursor">{{ $root.FormatDateFromNow(scope.row.ctime) }}</span>
                        </el-tooltip>
                    </template>
                </el-table-column>
                <el-table-column prop="ssh_port" width="100" label="提交者">
                    <template slot-scope="scope">
                        <el-tooltip effect="dark" :content="scope.row.email" placement="top">
                            <span class="app-cursor">{{ scope.row.username }}</span>
                        </el-tooltip>
                    </template>
                </el-table-column>
                <el-table-column prop="audit_status" width="100" label="审核">
                    <template slot-scope="scope">
                        <span class="app-color-warning" v-if="scope.row.audit_status == 1">待审核</span>
                        <span class="app-color-success" v-else-if="scope.row.audit_status == 2">通过</span>
                        <span class="app-color-error" v-else-if="scope.row.audit_status == 3">拒绝</span>
                        <span v-else>--</span>
                    </template>
                </el-table-column>
                <el-table-column prop="ssh_port" width="100" label="状态">
                    <template slot-scope="scope">
                        <span v-if="scope.row.status == 1"><i class="iconfont small left icon-wait"></i>待上线</span>
                        <span v-else-if="scope.row.status == 2"><i class="iconfont small left icon-coffee"></i>上线中</span>
                        <span class="app-color-success" v-else-if="scope.row.status == 3"><i class="iconfont small left icon-success"></i>成功</span>
                        <span class="app-color-error" v-else-if="scope.row.status == 4"><i class="iconfont small left icon-failed"></i>失败</span>
                        <span class="app-color-gray" v-else-if="scope.row.status == 5"><i class="iconfont small left icon-drop"></i>废弃</span>
                        <span v-else>--</span>
                    </template>
                </el-table-column>
                <el-table-column :label="$t('operate')" width="100" align="right">
                    <template slot-scope="scope">
                        <el-dropdown trigger="click" @command="operateHandler($event, scope.row)">
                            <el-button size="small">
                                操作<i class="el-icon-arrow-down el-icon--right"></i>
                            </el-button>
                            <el-dropdown-menu class="app-op-dropdown" slot="dropdown">
                                <el-dropdown-item command="view">
                                    <i class="iconfont left small icon-view"></i>查看
                                </el-dropdown-item>
                                <el-dropdown-item command="edit" 
                                v-if="scope.row.status == 1">
                                    <i class="iconfont left small icon-edit"></i>编辑
                                </el-dropdown-item>
                                <el-dropdown-item command="audit"
                                v-if="scope.row.audit_status == 1">
                                    <i class="iconfont left small icon-audit"></i>审核
                                </el-dropdown-item>
                                <el-dropdown-item command="deploy"
                                v-if="scope.row.audit_status == 2 && (scope.row.status == 1 || scope.row.status == 4)">
                                    <i class="iconfont left small icon-coffee"></i>上线
                                </el-dropdown-item>
                                <el-dropdown-item command="drop"
                                v-if="scope.row.status != 2 && scope.row.status != 5">
                                    <i class="iconfont left small icon-drop"></i>废弃
                                </el-dropdown-item>
                            </el-dropdown-menu>
                        </el-dropdown>
                    </template>
                </el-table-column>
            </el-table>
            <el-pagination
                background
                layout="prev, pager, next"
                class="app-pagination"
                @current-change="currentChangeHandler"
                :current-page.sync="$root.Page"
                :page-size="$root.PageSize"
                :total="$root.Total">
            </el-pagination>
        </el-card>

        <el-dialog
        :width="$root.DialogNormalWidth"
        :title="dialogTitle"
        :visible.sync="dialogVisible"
        @close="closeDialogHandler">
            <div class="app-dialog" v-loading="dialogLoading">
                <el-form size="medium" label-width="120px">
                    <el-form-item label="空间名称">
                        {{ dialogDetail.space_name }}
                    </el-form-item>
                    <el-form-item label="项目名称">
                        {{ dialogDetail.project_name }}
                    </el-form-item>
                    <el-form-item label="上线单">
                        {{ dialogDetail.name }}
                    </el-form-item>
                    <el-form-item label="上线模式">
                        <div v-if="dialogDetail.deploy_mode == 1">
                            <i class="iconfont icon-branch"></i> {{ this.$t('branch_deploy') }} - 分支名: {{ dialogDetail.branch_name }} - 版本: {{ dialogDetail.commit_version ? dialogDetail.commit_version :  'HEAD'}}
                        </div>
                        <div v-else>
                            <i class="iconfont icon-branch"></i> {{ this.$t('tag_deploy') }} - {{ dialogDetail.branch_name }}
                        </div>
                    </el-form-item>
                    <el-form-item label="上线说明">
                        {{ dialogDetail.description }}
                    </el-form-item>
                    <el-form-item label="审核状态">
                        {{ this.auditStatusTitle(dialogDetail.audit_status) }}
                    </el-form-item>
                    <el-form-item label="提交者">
                        {{ dialogDetail.username }} - {{ dialogDetail.email }}
                    </el-form-item>
                    <el-form-item label="提交时间">
                        {{ this.$root.FormatDateTime(dialogDetail.ctime) }}
                    </el-form-item>
                    <template v-if="dialogDetail.cmd == 'audit'">
                        <el-form-item label="审核" v-if="dialogDetail.status == 1 && dialogDetail.audit_status == 1">
                            <div>
                                <el-radio v-model="auditStatus" :label="2"><span class="app-color-success">审核通过</span></el-radio>
                                <el-radio v-model="auditStatus" :label="3"><span class="app-color-error">审核拒绝</span></el-radio>
                            </div>
                        </el-form-item>
                        <el-form-item label="拒绝原因" v-if="dialogDetail.status == 1 && dialogDetail.audit_status == 1 && auditStatus == 3">
                            <el-input type="textarea" :autosize="{ minRows: 2 }" v-model="auditRefusalReason"></el-input>
                        </el-form-item>
                        <el-form-item>
                            <el-button size="small" type="primary" @click="auditStatusHandler">审核</el-button>
                            <el-button size="small" @click="closeDialogHandler">关闭</el-button>
                        </el-form-item>
                    </template>
                </el-form>
            </div>
        </el-dialog>

        <el-dialog
        :width="$root.DialogNormalWidth"
        title="编辑上线单"
        :visible.sync="dialogEditVisible"
        @close="closeEditDialogHandler">
            <div class="app-dialog" v-loading="dialogLoading">
                <el-form 
                class="app-form" 
                ref="dialogRef" 
                :model="dialogForm" 
                size="medium" 
                label-width="130px">
                    <el-form-item
                    :label="$t('project_name')">
                    </el-form-item>
                </el-form>
            </div>
        </el-dialog>
    </div>
</template>

<script>
import { applyProjectAllApi, applyListApi, applyDetailApi, applyProjectDetailApi, applyAuditApi } from '@/api/deploy'
import { resolve, reject } from 'q';
export default {
    data() {
        return {
            searchInput: '',
            searchTime: undefined,
            searchAuditStatus: undefined,
            searchStatus: undefined,
            searchProjectId: undefined,

            tableLoading: false,
            tableData: [],

            timeList: [
                {time: 1, label: '今天'},
                {time: 7, label: '7天内'},
                {time: 30, label: '一个月内'},
                {time: 90, label: '3个月内'},
                {time: 365, label: '一年内'},
                {time: 0, label: '时间不限'},
            ],
            statusList: [
                {status: 1, label: '未上线'},
                {status: 2, label: '上线中'},
                {status: 3, label: '上线成功'},
                {status: 4, label: '上线失败'},
                {status: 5, label: '已废弃'},
            ],
            auditStatusList: [
                {status: 1, label: '待审核'},
                {status: 2, label: '审核通过'},
                {status: 3, label: '审核拒绝'},
            ],
            projectList: [],

            dialogTitle: '',
            dialogVisible: false,
            dialogEditVisible: false,
            dialogLoading: false,
            dialogBtnLoading: false,
            dialogDetail: {},
            dialogForm: {},

            auditStatus: 2,
            auditRefusalReason: '',
        }
    },
    methods: {
        operateHandler(cmd, row) {
            switch(cmd) {
                case 'view':
                case 'audit':
                    this.viewHandler(cmd, row)
                    break
                case 'edit':
                    this.editHandler(cmd, row)
                    break
            }
        },
        closeDialogHandler() {
            this.dialogVisible = false
        },
        openDialogHandler(title) {
            this.dialogTitle = title
            this.dialogVisible = true
        },
        closeEditDialogHandler() {
            this.dialogEditVisible = false
        },
        openEditDialogHandler() {
            this.dialogEditVisible = true
        },
        auditStatusHandler() {
            let postData = {
                id: this.dialogDetail.id,
                audit_status: this.auditStatus,
                audit_refusal_reasion: this.auditRefusalReason,
            }
            this.dialogBtnLoading = true
            applyAuditApi(postData).then(res => {
                this.$message({
                    message: '审核成功',
                    type: 'success',
                    duration: 1000,
                    onClose: () => {
                        this.closeDialogHandler()
                        this.loadTableData()
                        this.dialogBtnLoading = false
                    },
                })
            }).catch(err => {
                this.dialogBtnLoading = false
            })
        },
        editHandler(cmd, row) {
            this.dialogLoading = true
            this.getApplyDetail(cmd, row).then(detail => {
                this.dialogLoading = false
                this.dialogDetail = detail
                this.openEditDialogHandler()
            }).catch(err => {
                this.dialogLoading = false
            })
        },
        viewHandler(cmd, row) {
            this.dialogLoading = true
            this.getApplyDetail(cmd, row).then(detail => {
                this.dialogLoading = false
                this.dialogDetail = detail
                this.openDialogHandler(cmd == 'view' ? '查看' : '审核')
            }).catch(err => {
                this.dialogLoading = false
            })
        },
        getApplyDetail(cmd, row) {
            let promise = new Promise((resolve, reject) => {
                let projDetailPromise = new Promise((resolve, reject) => {
                    applyProjectDetailApi({id: row.project_id}).then(res => {
                        resolve(res)
                    }).catch(err => {
                        reject(err)
                    })
                })
                let applyDetailPromise = new Promise((resolve, reject) => {
                    applyDetailApi({id: row.id}).then(res => {
                        resolve(res)
                    }).catch(err => {
                        reject(err)
                    })
                })
                Promise.all([projDetailPromise, applyDetailPromise]).then(res => {
                    let projDetail = res[0]
                    let applyDetail = res[1]
                    let dialogDetail = {
                        id: row.id,
                        space_name: row.space_name,
                        project_name: row.project_name,
                        name: row.name,
                        deploy_mode: projDetail.deploy_mode,
                        branch_name: applyDetail.branch_name,
                        commit_version: applyDetail.commit_version,
                        description: applyDetail.description,
                        audit_status: row.audit_status,
                        email: row.email,
                        username: row.username,
                        ctime: row.ctime,
                        audit_status: applyDetail.audit_status,
                        status: applyDetail.status,
                        cmd: cmd,
                    }
                    resolve(dialogDetail)
                }).catch(err => {
                    reject(err)
                })
            })
            return promise
        },
        searchHandler() {
            this.$root.PageInit()
            this.loadTableData()
        },
        currentChangeHandler() {
            this.loadTableData()
        },
        auditStatusTitle(auditStatus) {
            let auditTitle = ''
            this.auditStatusList.forEach(item => {
                if (auditStatus == item.status) {
                    auditTitle = item.label
                }
            })
            return auditTitle
        },
        loadTableData() {
            this.tableLoading = true
            let query = {
                keyword : this.searchInput,
                time: this.searchTime,
                audit_status: this.searchAuditStatus,
                status: this.searchStatus,
                project_id: this.searchProjectId,
                offset: this.$root.PageOffset(),
                limit: this.$root.PageSize,
            }
            applyListApi(query).then(res => {
                this.tableData = res.list
                this.$root.Total = res.total
                this.tableLoading = false
            }).catch(err => {
                this.tableLoading = false
            })
        },
        loadProjectAll() {
            applyProjectAllApi().then(res => {
                if (res && res.length > 0 ) {
                    this.projectList = res
                }
            })
        },
    },
    mounted() {
        this.$root.PageInit()
        this.loadTableData()
        this.loadProjectAll()
    },
}
</script>