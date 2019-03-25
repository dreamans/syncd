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
                <el-table-column :label="$t('project_name')">
                    <template slot-scope="scope">
                         {{ scope.row.project_name }}
                         <el-tooltip effect="dark" :content="$t('belong_to_space') + ': ' + scope.row.space_name" placement="top">
                            <span class="app-cursor"><i class="iconfont icon-space"></i></span>
                        </el-tooltip>
                    </template>
                </el-table-column>
                <el-table-column prop="ssh_port" width="100" :label="$t('submit_time')">
                    <template slot-scope="scope">
                        <el-tooltip effect="dark" :content="$root.FormatDateTime(scope.row.ctime)" placement="top">
                            <span class="app-cursor">{{ $root.FormatDateFromNow(scope.row.ctime) }}</span>
                        </el-tooltip>
                    </template>
                </el-table-column>
                <el-table-column prop="ssh_port" width="100" :label="$t('submiter')">
                    <template slot-scope="scope">
                        <el-tooltip effect="dark" :content="scope.row.email" placement="top">
                            <span class="app-cursor">{{ scope.row.username }}</span>
                        </el-tooltip>
                    </template>
                </el-table-column>
                <el-table-column prop="audit_status" width="100" :label="$t('audit')">
                    <template slot-scope="scope">
                        <span class="app-color-warning" v-if="scope.row.audit_status == $root.AuditStatusPending">{{ $t('unaudit') }}</span>
                        <span class="app-color-success" v-else-if="scope.row.audit_status == $root.AuditStatusOk">{{ $t('pass') }}</span>
                        <span class="app-color-error" v-else-if="scope.row.audit_status == $root.AuditStatusRefuse">{{ $t('denied') }}</span>
                        <span v-else>--</span>
                    </template>
                </el-table-column>
                <el-table-column prop="ssh_port" width="100" :label="$t('status')">
                    <template slot-scope="scope">
                        <span v-if="scope.row.status == $root.ApplyStatusNone"><i class="iconfont small left icon-wait"></i>{{ $t('wait_online') }}</span>
                        <span v-else-if="scope.row.status == $root.ApplyStatusIng"><i class="iconfont small left icon-coffee"></i>{{ $t('onlineing') }}</span>
                        <span class="app-color-success" v-else-if="scope.row.status == $root.ApplyStatusSuccess"><i class="iconfont small left icon-success"></i>{{ $t('success') }}</span>
                        <span class="app-color-error" v-else-if="scope.row.status == $root.ApplyStatusFailed"><i class="iconfont small left icon-failed"></i>{{ $t('failed') }}</span>
                        <span class="app-color-gray" v-else-if="scope.row.status == $root.ApplyStatusDrop"><i class="iconfont small left icon-drop"></i>{{ $t('drop') }}</span>
                        <span class="app-color-error" v-else-if="scope.row.status == $root.ApplyStatusRollback"><i class="iconfont small left icon-rollback"></i>{{ $t('rollback') }}</span>
                        <span v-else>--</span>
                    </template>
                </el-table-column>
                <el-table-column :label="$t('operate')" width="100" align="right">
                    <template slot-scope="scope">
                        <el-dropdown trigger="click" @command="operateHandler($event, scope.row)">
                            <el-button size="small">
                                {{ $t('operate') }}<i class="el-icon-arrow-down el-icon--right"></i>
                            </el-button>
                            <el-dropdown-menu class="app-op-dropdown" slot="dropdown">
                                <el-dropdown-item command="view" v-if="$root.CheckPriv($root.Priv.DEPLOY_VIEW)">
                                    <i class="iconfont left small icon-view"></i>{{ $t('view') }}
                                </el-dropdown-item>
                                <el-dropdown-item command="edit" 
                                v-if="scope.row.status == $root.ApplyStatusNone && (scope.row.audit_status == $root.AuditStatusPending || scope.row.audit_status == $root.AuditStatusRefuse) && $root.CheckPriv($root.Priv.DEPLOY_EDIT)">
                                    <i class="iconfont left small icon-edit"></i>{{ $t('edit') }}
                                </el-dropdown-item>
                                <el-dropdown-item command="audit"
                                v-if="scope.row.audit_status == $root.AuditStatusPending && scope.row.status == $root.ApplyStatusNone && $root.CheckPriv($root.Priv.DEPLOY_AUDIT)">
                                    <i class="iconfont left small icon-audit"></i>{{ $t('audit') }}
                                </el-dropdown-item>
                                <el-dropdown-item command="deploy"
                                v-if="scope.row.audit_status == $root.AuditStatusOk && (scope.row.status == $root.ApplyStatusNone || scope.row.status == $root.ApplyStatusIng || scope.row.status == $root.ApplyStatusSuccess || scope.row.status == $root.ApplyStatusFailed || scope.row.status == $root.ApplyStatusRollback) && $root.CheckPriv($root.Priv.DEPLOY_DEPLOY)">
                                    <i class="iconfont left small icon-coffee"></i>{{ $t('online') }}
                                </el-dropdown-item>
                                <el-dropdown-item command="drop"
                                v-if="scope.row.status != $root.ApplyStatusIng && scope.row.status != $root.ApplyStatusDrop && $root.CheckPriv($root.Priv.DEPLOY_DROP)">
                                    <i class="iconfont left small icon-drop"></i>{{ $t('drop') }}
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
                    <el-form-item :label="$t('space_name')">
                        {{ dialogDetail.space_name }}
                    </el-form-item>
                    <el-form-item :label="$t('project_name')">
                        {{ dialogDetail.project_name }}
                    </el-form-item>
                    <el-form-item :label="$t('apply_order')">
                        {{ dialogDetail.name }}
                    </el-form-item>
                    <el-form-item :label="$t('deploy_mode')">
                        <div v-if="dialogDetail.deploy_mode == 1">
                            <i class="iconfont icon-branch"></i> {{ this.$t('branch_deploy') }} - {{ $t('branch_name') }}: {{ dialogDetail.branch_name }} - {{ $t('version') }}: {{ dialogDetail.commit_version ? dialogDetail.commit_version :  'HEAD'}}
                        </div>
                        <div v-else>
                            <i class="iconfont icon-tag"></i> {{ this.$t('tag_deploy') }} - {{ dialogDetail.branch_name }}
                        </div>
                    </el-form-item>
                    <el-form-item :label="$t('deploy_illustrate')">
                        {{ dialogDetail.description }}
                    </el-form-item>
                    <el-form-item :label="$t('audit_status')">
                        {{ this.auditStatusTitle(dialogDetail.audit_status) }}
                    </el-form-item>
                    <el-form-item :label="$t('submiter')">
                        {{ dialogDetail.username }} - {{ dialogDetail.email }}
                    </el-form-item>
                    <el-form-item :label="$t('submit_time')">
                        {{ this.$root.FormatDateTime(dialogDetail.ctime) }}
                    </el-form-item>
                    <template v-if="dialogDetail.cmd == 'audit'">
                        <el-form-item :label="$t('audit')" v-if="dialogDetail.status == $root.ApplyStatusNone && dialogDetail.audit_status == $root.AuditStatusPending">
                            <div>
                                <el-radio v-model="auditStatus" :label="$root.AuditStatusOk"><span class="app-color-success">{{ $t('audit_pass') }}</span></el-radio>
                                <el-radio v-model="auditStatus" :label="$root.AuditStatusRefuse"><span class="app-color-error">{{ $t('audit_denied') }}</span></el-radio>
                            </div>
                        </el-form-item>
                        <el-form-item :label="$t('deined_reason')" v-if="dialogDetail.status == $root.ApplyStatusNone && dialogDetail.audit_status == $root.AuditStatusPending && auditStatus == $root.AuditStatusRefuse">
                            <el-input type="textarea" :autosize="{ minRows: 2 }" v-model="auditRefusalReason"></el-input>
                        </el-form-item>
                        <el-form-item>
                            <el-button size="small" type="primary" @click="dialogSubmitAuditStatusHandler">{{ $t('audit') }}</el-button>
                            <el-button size="small" @click="closeDialogHandler">{{ $t('close') }}</el-button>
                        </el-form-item>
                    </template>
                </el-form>
            </div>
        </el-dialog>

        <el-dialog
        :width="$root.DialogNormalWidth"
        :title="$t('edit_apply_order')"
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
                        {{ dialogDetail.project_name }}
                    </el-form-item>
                    <el-form-item 
                    :label="$t('apply_name')">
                        {{ dialogDetail.name}}
                    </el-form-item>
                    <el-form-item :label="$t('deploy_mode')">
                        <span v-if="dialogDetail.deploy_mode == $root.DeployModeBranch">
                            <i class="iconfont icon-branch"></i> - {{ $t('branch_deploy') }}<template v-if="dialogDetail.repo_branch"> - <strong>{{ dialogDetail.repo_branch }}</strong> {{ $t('branch') }}</template>
                        </span>
                        <span v-if="dialogDetail.deploy_mode == $root.DeployModelTag">
                            <i class="iconfont icon-tag"></i> {{ $t('tag_deploy') }}
                        </span>
                    </el-form-item>

                    <el-form-item 
                    v-if="dialogDetail.deploy_mode == $root.DeployModelTag"
                    :label="$t('tag_name')"
                    prop="branch_name"
                    :rules="[
                        { required: true, message: $t('tag_name_cannot_empty'), trigger: 'blur'},
                    ]">
                        <el-input class="app-input-mini" :placeholder="$t('please_input_tag_name')" v-model="dialogForm.branch_name" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item 
                    v-if="dialogDetail.deploy_mode == $root.DeployModeBranch && dialogDetail.repo_branch == ''"
                    :label="$t('branch_name')"
                    prop="branch_name"
                    :rules="[
                        { required: true, message: $t('branch_name_cannot_empty'), trigger: 'blur'},
                    ]">
                        <el-input class="app-input-mini" :placeholder="$t('please_input_branch_name')" v-model="dialogForm.branch_name" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item 
                    v-if="dialogDetail.deploy_mode == $root.DeployModeBranch"
                    :label="$t('commit_version')"
                    prop="commit_version">
                        <el-input class="app-input-normal" :placeholder="$t('please_input_commit_version')" v-model="dialogForm.commit_version" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item 
                    :label="$t('deploy_illustrate')"
                    prop="description"
                    :rules="[
                        { required: true, message: $t('deploy_illustrate_cannot_empty'), trigger: 'blur'},
                    ]">
                        <el-input :rows="4" type="textarea" :placeholder="$t('please_input_deploy_illustrate')" v-model="dialogForm.description" autocomplete="off"></el-input>
                    </el-form-item>
                </el-form>
                <div slot="footer" class="dialog-footer">
                    <el-button size="small" @click="closeEditDialogHandler">{{ $t('cancel') }}</el-button>
                    <el-button :loading="dialogBtnLoading" size="small" type="primary" @click="dialogSubmitEditHandler">{{ $t('enter') }}</el-button>
                </div>
            </div>
        </el-dialog>
    </div>
</template>

<script>
import { applyProjectAllApi, applyListApi, applyDetailApi, applyProjectDetailApi, applyAuditApi, applyUpdateApi, applyDropApi } from '@/api/deploy'
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
                {time: 1, label: this.$t('today')},
                {time: 7, label: this.$t('7day')},
                {time: 30, label: this.$t('within_one_month')},
                {time: 90, label: this.$t('within_three_months')},
                {time: 365, label: this.$t('within_a_year')},
                {time: 0, label: this.$t('any_time')},
            ],
            statusList: [
                {status: this.$root.ApplyStatusNone, label: this.$t('not_online')},
                {status: this.$root.ApplyStatusIng, label: this.$t('onlineing')},
                {status: this.$root.ApplyStatusSuccess, label: this.$t('online_success')},
                {status: this.$root.ApplyStatusFailed, label: this.$t('online_failed')},
                {status: this.$root.ApplyStatusDrop, label: this.$t('deprecated')},
            ],
            auditStatusList: [
                {status: this.$root.AuditStatusPending, label: this.$t('unaudit')},
                {status: this.$root.AuditStatusOk, label: this.$t('audit_pass')},
                {status: this.$root.AuditStatusRefuse, label: this.$t('audit_denied')},
            ],
            projectList: [],

            dialogTitle: '',
            dialogVisible: false,
            dialogEditVisible: false,
            dialogLoading: false,
            dialogBtnLoading: false,
            dialogDetail: {},
            dialogForm: {
                id: 0,
                branch_name: '',
                commit_version: '',
                description: '',
            },

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
                case 'drop':
                    this.dropHandler(cmd, row)
                    break
                case 'deploy':
                    this.deployHandler(cmd, row)
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
        dialogSubmitEditHandler() {
            this.$refs.dialogRef.validate((valid) => {
                if (!valid) {
                    return false;
                }
                this.dialogBtnLoading = true
                applyUpdateApi(this.dialogForm).then(res => {
                    this.$message({
                        message: this.$t('update_success'),
                        type: 'success',
                        duration: 1000,
                        onClose: () => {
                            this.closeEditDialogHandler()
                            this.loadTableData()
                            this.dialogBtnLoading = false
                        },
                    })
                }).catch(err => {
                    this.dialogBtnLoading = false
                })
            })
        },
        dialogSubmitAuditStatusHandler() {
            let postData = {
                id: this.dialogDetail.id,
                audit_status: this.auditStatus,
                audit_refusal_reasion: this.auditRefusalReason,
            }
            this.dialogBtnLoading = true
            applyAuditApi(postData).then(res => {
                this.$message({
                    message: this.$t('audit_success'),
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
        deployHandler(cmd, row) {
            this.$router.push({name: 'deployRelease', query: { id: row.id}})
        },
        dropHandler(cmd, row) {
            this.$root.ConfirmDelete(() => {
                applyDropApi({id: row.id}).then(res => {
                    this.loadTableData()
                })
            }, this.$t('drop_deploy_apply_tips'))
        },
        editHandler(cmd, row) {
            this.dialogLoading = true
            this.getApplyDetail(cmd, row).then(detail => {
                this.dialogLoading = false
                this.dialogDetail = detail
                this.dialogForm = {
                    id: row.id,
                    branch_name: detail.branch_name,
                    commit_version: detail.commit_version,
                    description: detail.description,
                }
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
                this.openDialogHandler(cmd == 'view' ? this.$t('view') : this.$t('audit'))
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
                        repo_branch: projDetail.repo_branch,
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
        if (this.$route.query.id) {
            this.searchInput = this.$route.query.id
        }
        this.$root.PageInit()
        this.loadTableData()
        this.loadProjectAll()
    },
}
</script>