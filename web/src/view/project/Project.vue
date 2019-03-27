<template>
    <div>
        <el-card shadow="never" v-if="!spaceId">
            <el-alert
            class="app-btn-group"
            :title="$t('prompt_message')"
            type="info"
            :closable="false"
            :description="$t('project_select_space_tips')"
            show-icon></el-alert>
            <div class="app-btn-group">
                <el-button type="primary" @click="switchSpaceHandler" icon="el-icon-refresh" size="small">{{ $t('select_project_space') }}</el-button>
            </div>
            <div class="app-divider"></div>
        </el-card>
        <el-card shadow="never" v-else>
            <div class="app-btn-group">
                <el-button type="primary" @click="switchSpaceHandler" icon="el-icon-refresh" size="small">{{ $t('switch_project_space') }}</el-button>
            </div>
            <div class="app-btn-group" style="min-height: 30px">
                <div v-loading="spaceLoading">
                    <h4>{{ this.spaceDetail.name }}</h4>
                    <p class="app-description">{{ this.spaceDetail.description }}</p>
                </div>
            </div>
            <div class="app-divider"></div>
            <el-row class="app-btn-group">
                <el-col :span="4">
                    <el-button v-if="$root.CheckPriv($root.Priv.PROJECT_NEW)" @click="openAddDialogHandler" icon="iconfont left small icon-add" size="medium" type="primary">{{ $t('add_project') }}</el-button>&nbsp;
                </el-col>
                <el-col :span="6" :offset="14">
                    <el-input @keyup.enter.native="searchHandler" v-model="searchInput" size="medium" :placeholder="$t('please_input_keyword')">
                        <el-button @click="searchHandler" slot="append" icon="el-icon-search"></el-button>
                    </el-input>
                </el-col>
            </el-row>
            
            <el-table
                class="app-table"
                size="medium"
                v-loading="tableLoading"
                :data="tableData">
                <el-table-column prop="name" :label="$t('project_name')"></el-table-column>
                <el-table-column align="center" prop="need_audit" width="150" :label="$t('open_audit')">
                    <template slot-scope="scope">
                        <span v-if="scope.row.need_audit == 1">{{ $t('yes') }}</span>
                        <span v-else>{{ $t('no') }}</span>
                    </template>
                </el-table-column>
                <el-table-column prop="status" width="120" :label="$t('project_enable')">
                    <template slot-scope="scope">
                        <el-switch
                        v-if="$root.CheckPriv($root.Priv.PROJECT_AUDIT)"
                        @change="enableSwitchHandler($event, scope.row)"
                        v-model="scope.row.status"
                        :active-value="1"
                        :inactive-value="0"
                        active-color="#13ce66">
                        </el-switch>
                        <span style="margin-left: 5px;" v-if="scope.row.status">{{ $t('have_enabled') }}</span>
                        <span style="margin-left: 5px;" v-else>{{ $t('not_enable') }}</span>
                    </template>
                </el-table-column>
                <el-table-column :label="$t('operate')" width="380" align="right">
                    <template slot-scope="scope">
                        <el-button
                        v-if="$root.CheckPriv($root.Priv.PROJECT_HOOK)"
                        icon="iconfont small left icon-webhook"
                        type="text"
                        @click="openHookDialogHandler(scope.row)">{{ $t('hook') }}</el-button>
                        <el-button
                        v-if="$root.CheckPriv($root.Priv.PROJECT_BUILD)"
                        icon="iconfont small left icon-build"
                        type="text"
                        @click="openBuildDialogHandler(scope.row)">{{ $t('build_setting') }}</el-button>
                        <el-button
                        v-if="$root.CheckPriv($root.Priv.PROJECT_VIEW)"
                        icon="el-icon-view"
                        type="text"
                        @click="openViewDialogHandler(scope.row)">{{ $t('view') }}</el-button>
                        <el-button
                        v-if="$root.CheckPriv($root.Priv.PROJECT_EDIT)"
                        icon="el-icon-edit"
                        type="text"
                        @click="openEditDialogHandler(scope.row)">{{ $t('edit') }}</el-button>
                        <el-button
                        v-if="$root.CheckPriv($root.Priv.PROJECT_DEL)"
                        type="text"
                        icon="el-icon-delete"
                        class="app-danger"
                        @click="deleteHandler(scope.row)">{{ $t('delete') }}</el-button>
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

        <el-dialog :top="$root.DialogNormalTop" :width="$root.DialogLargeWidth" :title="dialogTitle" :visible.sync="dialogVisible" @close="dialogCloseHandler">
            <div class="app-dialog" v-loading="dialogLoading">
                <el-form class="app-form" ref="dialogRef" :model="dialogForm" size="medium" label-width="130px">
                    <h4 class="app-form-subtitle">{{ $t('base_setting') }}</h4>
                    <el-form-item 
                    :label="$t('project_name')"
                    prop="name"
                    :rules="[
                        { required: true, message: $t('name_cannot_empty'), trigger: 'blur'},
                    ]">
                        <el-input :placeholder="$t('please_input_project_name')" v-model="dialogForm.name" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item
                    :label="$t('description')"
                    prop="description">
                        <el-input :placeholder="$t('please_input_project_description')" type="textarea" :rows="3" v-model="dialogForm.description" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item
                    :label="$t('open_audit')"
                    prop="need_audit"
                    >
                        <div>
                            <el-switch
                            v-model="dialogForm.need_audit"
                            :active-value="1"
                            :inactive-value="0"
                            active-color="#13ce66">
                            </el-switch>
                            <span style="margin-left: 8px;">
                                <span v-if="dialogForm.need_audit">{{ $t('open') }}</span>
                                <span v-else>{{ $t('close') }}</span>
                            </span>
                        </div>
                        <div class="app-form-explain">{{ $t('if_open_apply_need_audit') }}</div>
                    </el-form-item>

                    <div class="app-divider"></div>
                    <h4 class="app-form-subtitle">{{ $t('repo_setting') }}</h4>

                    <el-form-item
                    :label="$t('repo_url')"
                    prop="repo_url"
                    :rules="[
                        { required: true, message: $t('repo_url_cannot_empty'), trigger: 'blur'},
                    ]">
                        <el-input :placeholder="$t('please_input_repo_url')" v-model="dialogForm.repo_url" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item
                    :label="$t('deploy_mode')"
                    prop="deploy_mode"
                    :rules="[
                        { required: true, message: $t('deploy_mode_cannot_empty'), trigger: 'blur'},
                    ]">
                        <el-radio-group v-model="dialogForm.deploy_mode">
                            <el-radio :label="1">{{ $t('branch_deploy') }}</el-radio>
                            <el-radio :label="2">{{ $t('tag_deploy') }}</el-radio>
                        </el-radio-group>
                        <div class="app-form-explain">{{ $t('deploy_mode_tips') }}</div>
                    </el-form-item>

                    <el-form-item 
                    v-if="dialogForm.deploy_mode == 1"
                    :label="$t('repo_branch')"
                    prop="repo_branch">
                        <el-input class="app-input-mini" v-model="dialogForm.repo_branch" autocomplete="off"></el-input>
                        <div class="app-form-explain">{{ $t('if_not_need_to_assign_branch_name') }}</div>
                    </el-form-item>
                    <div class="app-divider"></div>
                    <h4 class="app-form-subtitle">{{ $t('deploy_setting') }}</h4>

                    <el-form-item 
                    :label="$t('online_cluster')"
                    prop="online_cluster"
                    :rules="[
                        { required: true, message: $t('online_cluster_cannot_empty'), trigger: 'blur' },
                    ]">
                        <el-select
                        class="app-input-mini"
                        v-model="selectOnlineCluster"
                        filterable
                        clearable
                        @change="selectClusterHandler"
                        :placeholder="$t('please_input_keyword')">
                            <el-option
                            v-for="cluster in clusterList"
                            :key="cluster.id"
                            :label="cluster.name"
                            :value="cluster.id">
                            </el-option>
                        </el-select>
                        <div v-if="dialogForm.online_cluster && dialogForm.online_cluster.length">
                            <span>{{ $t('selected_cluster_list') }}</span>
                            <ul class="app-form-box">
                                <li class="item" v-for="id in dialogForm.online_cluster" :key="id">
                                    <span><i class="iconfont small left icon-cluster"></i>{{ formatClusterName(id) }}</span>
                                    <span>
                                        <el-button @click="removeClusterHandler(id)" icon="el-icon-delete" type="text">{{ $t('remove') }}</el-button>
                                    </span>
                                </li>
                            </ul>
                        </div>
                    </el-form-item>

                    <el-form-item 
                    :label="$t('user')"
                    prop="deploy_user"
                    :rules="[
                        { required: true, message: $t('user_cannot_empty'), trigger: 'blur'},
                    ]">
                        <el-input class="app-input-mini" :placeholder="$t('deploy_user')" v-model="dialogForm.deploy_user" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item 
                    :label="$t('path')"
                    prop="deploy_path"
                    :rules="[
                        { required: true, message: $t('deploy_path_cannot_empty'), trigger: 'blur'},
                    ]">
                        <el-input :placeholder="$t('deploy_path')" v-model="dialogForm.deploy_path" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item 
                    :label="$t('pre_deploy_cmd')"
                    prop="pre_deploy_cmd">
                        <el-input :placeholder="$t('pre_deploy_cmd_tips')" type="textarea" :rows="3" v-model="dialogForm.pre_deploy_cmd" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item 
                    :label="$t('after_deploy_cmd')"
                    prop="after_deploy_cmd">
                        <el-input :placeholder="$t('after_deploy_cmd_tips')" type="textarea" :rows="3" v-model="dialogForm.after_deploy_cmd" autocomplete="off"></el-input>
                    </el-form-item>

                    <div class="app-divider"></div>
                    <h4 class="app-form-subtitle">{{ $t('email_setting') }}</h4>
                    <el-form-item 
                    :label="$t('audit_notice')"
                    prop="audit_notice">
                        <el-input :placeholder="$t('audit_notice_tips')" v-model="dialogForm.audit_notice" autocomplete="off"></el-input>
                        <div class="app-form-explain" v-html="$t('audit_notice_explain')"></div>
                    </el-form-item>
                    <el-form-item 
                    :label="$t('deploy_notice')"
                    prop="deploy_notice">
                        <el-input :placeholder="$t('deploy_notice_tips')" v-model="dialogForm.deploy_notice" autocomplete="off"></el-input>
                        <div class="app-form-explain" v-html="$t('deploy_notice_explain')"></div>
                    </el-form-item>

                </el-form>
                <div slot="footer" class="dialog-footer">
                    <el-button size="small" @click="dialogCloseHandler">{{ $t('cancel') }}</el-button>
                    <el-button :loading="btnLoading" size="small" type="primary" @click="dialogSubmitHandler">{{ $t('enter') }}</el-button>
                </div>
            </div>
        </el-dialog>

        <el-dialog :top="$root.DialogNormalTop" :width="$root.DialogLargeWidth" :title="$t('view_project_info')" :visible.sync="dialogViewVisible" @close="dialogViewVisible = false">
            <div class="app-dialog" v-loading="dialogViewLoading">
                <el-form size="medium" label-width="130px">
                    <h4 class="app-form-subtitle">{{ $t('base_setting') }}</h4>
                    <el-form-item 
                    :label="$t('project_id')">
                        {{ dialogViewForm.id }}
                    </el-form-item>

                    <el-form-item  :label="$t('project_name')">
                        {{ dialogViewForm.name }}
                    </el-form-item>

                    <el-form-item :label="$t('description')">
                        {{ dialogViewForm.description }}
                    </el-form-item>

                    <el-form-item :label="$t('open_audit')">
                        <span v-if="dialogViewForm.need_audit">{{ $t('need_audit') }}</span>
                        <span v-else>{{ $t('not_audit') }}</span>
                    </el-form-item>

                    <div class="app-divider"></div>
                    <h4 class="app-form-subtitle">{{ $t('repo_setting') }}</h4>

                    <el-form-item :label="$t('repo_url')">
                        {{ dialogViewForm.repo_url }}
                    </el-form-item>

                    <el-form-item :label="$t('deploy_mode')">
                        <span v-if="dialogViewForm.deploy_mode == 1">
                            <i class="iconfont icon-branch"></i> - {{ $t('branch_deploy') }} - <strong>{{ dialogViewForm.repo_branch }}</strong> {{ $t('branch') }}
                        </span>
                        <span v-if="dialogViewForm.deploy_mode == 2">
                            <i class="iconfont icon-branch"></i> {{ $t('tag_deploy') }}
                        </span>
                    </el-form-item>

                    <el-form-item :label="$t('repo_branch')">
                        {{ dialogViewForm.repo_branch }}
                    </el-form-item>
                
                    <div class="app-divider"></div>
                    <h4 class="app-form-subtitle">{{ $t('deploy_setting') }}</h4>

                    <el-form-item :label="$t('online_cluster')">
                        <div v-if="dialogViewForm.online_cluster && dialogViewForm.online_cluster.length">
                            <ul class="app-form-box">
                                <li class="item" v-for="id in dialogViewForm.online_cluster" :key="id">
                                    <span><i class="iconfont small left icon-cluster"></i>{{ formatClusterName(id) }}</span>
                                </li>
                            </ul>
                        </div>
                    </el-form-item>

                    <el-form-item :label="$t('user')">
                        {{ dialogViewForm.deploy_user }}
                    </el-form-item>

                    <el-form-item :label="$t('path')">
                        {{ dialogViewForm.deploy_path }}
                    </el-form-item>

                    <el-form-item :label="$t('pre_deploy_cmd')">
                        <el-input type="textarea" :rows="3" :value="dialogViewForm.pre_deploy_cmd" readonly="readonly"></el-input>
                    </el-form-item>

                    <el-form-item :label="$t('after_deploy_cmd')">
                        <el-input type="textarea" :rows="3" :value="dialogViewForm.after_deploy_cmd" readonly="readonly"></el-input>
                    </el-form-item>

                    <el-form-item :label="$t('audit_notice')">
                        {{ dialogViewForm.audit_notice }}
                    </el-form-item>

                    <el-form-item :label="$t('deploy_notice')">
                        {{ dialogViewForm.deploy_notice }}
                    </el-form-item>

                </el-form>
            </div>
        </el-dialog>

        <el-dialog :top="$root.DialogNormalTop" :width="$root.DialogNormalWidth" :title="$t('edit_build_script')" :visible.sync="dialogBuildVisible" @close="dialogBuildVisible = false">
            <div class="app-dialog" v-loading="dialogBuildLoading">
                <div class="app-shell-editor">
                    <textarea id="editor-textarea"></textarea>
                </div>
                <h4 class="app-form-subtitle">{{ $t('illustrate') }}</h4>
                <div class="app-form-notice">
                    <p>{{ $t('build_script_tips') }}:</p>
                    <p>
                        <i class="iconfont icon-dot"></i><span class="code">${env_workspace}</span> - {{ $t('build_script_env_workspace') }}
                    </p>
                    <p>
                        <i class="iconfont icon-dot"></i><span class="code">${env_pack_file}</span> - <span v-html="$t('build_script_env_pack_file')"></span>
                    </p>
                    <p>
                        <a href="https://github.com/dreamans/syncd" class="app-link" target="_blank">{{ $t('view_build_script_eg') }}</a>
                    </p>
                </div>
                <div slot="footer" class="dialog-footer">
                    <el-button size="small" @click="dialogBuildVisible = false">{{ $t('cancel') }}</el-button>
                    <el-button :loading="btnLoading" size="small" type="primary" @click="dialogSubmitBuildHandler">{{ $t('enter') }}</el-button>
                </div>
            </div>
        </el-dialog>

        <el-dialog :top="$root.DialogNormalTop" :width="$root.DialogNormalWidth" :title="$t('edit_hook_script')" :visible.sync="dialogHookVisible" @close="dialogHookVisible = false">
            <div class="app-dialog" v-loading="dialogHookLoading">
                <el-form label-position="top" size="medium" label-width="130px">
                    <el-form-item :label="$t('build_hook_script')">
                        <el-input type="textarea" v-model="dialogHookForm.build_hook_script" :autosize="{ minRows: 3, maxRows: 10}"></el-input>
                        <div class="app-form-explain">
                            <p>{{ $t('build_hook_script_tips') }}:</p>
                            <p><i class="iconfont icon-dot"></i><span class="code">${env_apply_id}</span> - 申请单ID</p>
                            <p><i class="iconfont icon-dot"></i><span class="code">${env_apply_name}</span> - 申请单名称</p>
                            <p><i class="iconfont icon-dot"></i><span class="code">${env_pack_file}</span> - 打包的文件绝对路径</p>
                            <p><i class="iconfont icon-dot"></i><span class="code">${env_build_output}</span> - 构建脚本的原始输出 (JSON字符串)</p>
                            <p><i class="iconfont icon-dot"></i><span class="code">${env_build_errmsg}</span> - 构建错误信息</p>
                            <p><i class="iconfont icon-dot"></i><span class="code">${env_build_status}</span> - 构建结果状态，1 - 成功，0 - 失败</p>
                        </div>
                    </el-form-item>
                    <el-form-item :label="$t('deploy_hook_script')">
                        <el-input type="textarea" v-model="dialogHookForm.deploy_hook_script" :autosize="{ minRows: 3, maxRows: 10}"></el-input>
                        <div class="app-form-explain">
                            <p>{{ $t('deploy_hook_script_tips') }}:</p>
                            <p><i class="iconfont icon-dot"></i><span class="code">${env_apply_id}</span> - 申请单ID</p>
                            <p><i class="iconfont icon-dot"></i><span class="code">${env_apply_name}</span> - 申请单名称</p>
                            <p><i class="iconfont icon-dot"></i><span class="code">${env_deploy_status}</span> - 部署状态</p>
                        </div>
                    </el-form-item>
                </el-form>                
                <div slot="footer" class="dialog-footer">
                    <el-button size="small" @click="dialogHookVisible = false">{{ $t('cancel') }}</el-button>
                    <el-button :loading="btnLoading" size="small" type="primary" @click="dialogSubmitHookHandler">{{ $t('enter') }}</el-button>
                </div>
            </div>
        </el-dialog>

        <el-dialog 
        :width="$root.DialogSmallWidth"
        :visible.sync="dialogSpaceVisible"
        :title="$t('select_project_space')">
            <div style="padding-bottom: 50px;">
                <el-select filterable @change="selectSpaceHandler" v-model="spaceId" style="width: 100%;" size="small" :placeholder="$t('select_project_space')">
                    <el-option v-for="s in spaceList" :key="s.id" :value="s.id" :label="$root.Substr(s.name, 50)" ></el-option>
                </el-select>
            </div>
        </el-dialog>

    </div>
</template>

<script>
import { 
    listSpaceApi, 
    detailSpaceApi, 
    newProjectApi, 
    updateProjectApi, 
    listProjectApi, 
    switchStatusProjectApi, 
    detailProjectApi, 
    deleteProjectApi, 
    updateBuildScriptApi,
    updateHookScriptApi
} from '@/api/project'
import { listGroupApi } from '@/api/server' 
import codeMirror from 'codemirror/lib/codemirror.js'
import 'codemirror/lib/codemirror.css'
import 'codemirror/mode/shell/shell.js'
import 'codemirror/theme/dracula.css'
import 'codemirror/addon/scroll/simplescrollbars.css'
import 'codemirror/addon/scroll/simplescrollbars.js'

export default {
    data() {
        return {
            editorInstance: null,
            dialogBuildVisible: false,
            dialogBuildLoading: false,
            dialogBuildForm: {
                id: 0,
                build_script: '',
            },

            dialogHookVisible: false,
            dialogHookLoading: false,
            dialogHookForm: {
                id: 0,
                build_hook_script: '',
                deploy_hook_script: '',
            },

            searchInput: '',
            tableLoading: false,
            tableData: [],

            dialogViewVisible: false,
            dialogViewLoading: false,
            dialogViewForm: {},

            dialogVisible: false,
            dialogTitle: '',
            dialogLoading: false,
            dialogForm: {
                id: 0,
                name: '',
                description: '',
                need_audit: 0,
                repo_url: '',
                repo_branch: '',
                deploy_mode: 0,
                online_cluster: [],
                deploy_user: '',
                deploy_path: '',
                pre_deploy_cmd: '',
                after_deploy_cmd: '',
                audit_notice: '',
                deploy_notice: '',
                status: 0,
            },
            btnLoading: false,

            dialogSpaceVisible: false,
            spaceId: undefined,
            spaceLoading: false,
            spaceDetail: {},
            spaceList: [],

            clusterList: [],
            selectOnlineCluster: undefined,
        }
    },
    watch: {
        spaceId() {
            if (this.spaceId) {
                this.spaceLoading = true
                detailSpaceApi({id: this.spaceId}).then(res => {
                    this.spaceDetail = res
                    this.spaceLoading = false
                }).catch(err => {
                    this.spaceLoading = false
                })
                this.$root.PageInit()
                this.loadTableData()
            }
        }
    },
    methods: {
        openHookDialogHandler(row) {
            this.dialogHookVisible = true
            this.dialogHookLoading = true
            detailProjectApi({id: row.id}).then(res => {
                this.dialogHookForm = {
                    id: res.id,
                    build_hook_script: res.build_hook_script,
                    deploy_hook_script: res.deploy_hook_script,
                }
                this.dialogHookLoading = false
            })
        },
        dialogSubmitHookHandler() {
            updateHookScriptApi(this.dialogHookForm).then(res => {
                this.$root.MessageSuccess(() => {
                    this.dialogHookVisible = false
                })
            })
        },

        dialogSubmitBuildHandler() {
            this.dialogBuildForm.build_script = this.getBuildEditorValue()
            updateBuildScriptApi(this.dialogBuildForm).then(res => {
                this.$root.MessageSuccess(() => {
                    this.dialogBuildVisible = false
                })
            })
        },
        openBuildDialogHandler(row) {
            this.dialogBuildVisible = true
            this.dialogBuildLoading = true
            detailProjectApi({id: row.id}).then(res => {
                this.dialogBuildForm = {
                    id: res.id,
                    build_script: res.build_script,
                }
                this.dialogBuildLoading = false
                this.$nextTick(() => {
                    this.createBuildEditor(this.dialogBuildForm.build_script)
                })
            })
        },
        createBuildEditor(content) {
            if (!this.editorInstance) {
                this.editorInstance = codeMirror.fromTextArea(
                    document.getElementById('editor-textarea'),
                    {
                        theme: "dracula",
                        mode: 'shell',
                        tabSize: 4,
                        indentUnit: 4,
                        lineWrapping: 'wrap',
                        lineNumbers: true,
                        matchBrackets: true,
                        scrollbarStyle: 'simple',
                    }
                )
            }
            if (!content) {
                content = ''
            }
            this.editorInstance.setValue(content)
        },
        getBuildEditorValue() {
            if (!this.editorInstance) {        
                return ''
            }
            return this.editorInstance.getValue()
        },
        openViewDialogHandler(row) {
            this.dialogViewVisible = true
            this.dialogViewLoading = true
            detailProjectApi({id: row.id}).then(res => {
                this.dialogViewForm = res
                this.dialogViewLoading = false
            })
        },
        deleteHandler(row) {
            this.$root.ConfirmDelete(() => {
                deleteProjectApi({id: row.id}).then(res => {
                    this.$root.MessageSuccess()
                    this.$root.PageReset()
                    this.loadTableData()
                })
            })
        },
        enableSwitchHandler(val, row) {
            switchStatusProjectApi({status: val, id: row.id}).then(res => {
            }).catch(err => {
                row.status = Number(!val)
            })
        },
        searchHandler() {
            this.$root.PageInit()
            this.loadTableData()
        },
        openAddDialogHandler() {
            this.dialogVisible = true
            this.dialogTitle = this.$t('add_project')
        },
        openEditDialogHandler(row) {
            this.dialogVisible = true
            this.dialogTitle = this.$t('edit_project')
            this.dialogLoading = true
            detailProjectApi({id: row.id}).then(res => {
                this.dialogLoading = false
                this.dialogForm = res
            }).catch(err => {
                this.dialogCloseHandler()
            })
        },
        dialogSubmitHandler() {
            let vm = this
            this.$refs.dialogRef.validate((valid) => {
                if (!valid) {
                    return false;
                }
                this.btnLoading = true
                let opFn
                if (this.dialogForm.id) {
                    opFn = updateProjectApi
                } else {
                    opFn = newProjectApi
                }
                this.dialogForm.space_id = this.spaceId
                opFn(this.dialogForm).then(res => {
                    this.$root.MessageSuccess(() => {
                        this.dialogCloseHandler()
                        this.btnLoading = false
                        this.loadTableData()
                    })
                }).catch(err => {
                    this.btnLoading = false
                })
            });
        },
        dialogCloseHandler() {
            this.dialogVisible = false
            this.dialogLoading = false
            this.btnLoading = false
            this.$refs.dialogRef.resetFields();
            this.dialogForm = {}
        },
        selectClusterHandler(clusterId) {
            this.selectOnlineCluster = undefined
            if (!clusterId) {
                return
            }
            if (!this.dialogForm.online_cluster) {
                this.dialogForm.online_cluster = []
            }
            if (this.dialogForm.online_cluster.indexOf(clusterId) == -1) {
                this.dialogForm.online_cluster.push(clusterId)
            }
        },
        removeClusterHandler(clusterId) {
            if (this.dialogForm.online_cluster) {
                this.dialogForm.online_cluster.forEach((id, index) => {
                    if (id == clusterId) {
                        this.dialogForm.online_cluster.splice(index, 1)
                    }
                })
            }
        },
        selectSpaceHandler() {
            this.dialogSpaceVisible = false
        },
        switchSpaceHandler() {
            this.dialogSpaceVisible = true
        },
        formatClusterName(id) {
            let name = '--'
            this.clusterList.forEach(cluster => {
                if (id == cluster.id) {
                    name = cluster.name
                }
            })
            return name
        },
        currentChangeHandler() {
            this.loadTableData()
        },
        loadTableData() {
            this.tableLoading = true
            listProjectApi({space_id: this.spaceId, keyword: this.searchInput, offset: this.$root.PageOffset(), limit: this.$root.PageSize}).then(res => {
                this.tableData = res.list
                this.$root.Total = res.total
                this.tableLoading = false
            }).catch(err => {
                this.tableLoading = false
            })
        },
        loadSpaceList() {
            listSpaceApi({offset: 0, limit: 999}).then(res => {
                if (res.list) {
                    this.spaceList = res.list
                }
                this.initSpaceId()
            })
        },
        loadClusterList() {
            listGroupApi({offset: 0, limit: 999}).then(res => {
                if (res.list) {
                    this.clusterList = res.list
                }
            })
        },
        initSpaceId() {
            if (this.spaceList.length && !this.spaceId) {
                this.spaceId = this.spaceList[0].id
            }
        },
    },
    mounted() {
        this.loadSpaceList()
        this.loadClusterList()
    }
}
</script>
