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
                    <el-button @click="openAddDialogHandler" icon="iconfont left small icon-add" size="medium" type="primary">{{ $t('add_project') }}</el-button>
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
                <el-table-column prop="need_audit" width="180" :label="$t('open_audit')"></el-table-column>
                <el-table-column prop="status" width="180" :label="$t('project_enable')">
                    <template slot-scope="scope">
                        <span class="app-color-success" v-if="scope.row.status == '1'">
                            <i class="iconfont icon-unlock"></i> {{ $t('normal') }}
                        </span>
                        <span class="app-color-error" v-else>
                            <i class="iconfont icon-lock"></i> {{ $t('locking') }}
                        </span>
                    </template>
                </el-table-column>
                <el-table-column :label="$t('operate')" width="260" align="right">
                    <template slot-scope="scope">
                        <el-button
                        type="text"
                        icon="el-icon-delete"
                        class="app-danger"
                        @click="deleteHandler(scope.row)">{{ $t('remove') }}</el-button>
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

        <el-dialog :width="$root.DialogNormalWidth" :title="dialogTitle" :visible.sync="dialogVisible" @close="dialogCloseHandler">
            <div class="app-dialog" v-loading="dialogLoading">
                <el-form class="app-form" ref="dialogRef" :model="dialogForm" size="medium" label-width="130px">
                    <h4 class="app-form-subtitle">基本设置</h4>
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
                                <span v-if="dialogForm.need_audit">开启</span>
                                <span v-else>关闭</span>
                            </span>
                        </div>
                        <div class="app-form-explain">{{ $t('if_open_apply_need_audit') }}</div>
                    </el-form-item>

                    <div class="app-divider"></div>
                    <h4 class="app-form-subtitle">仓库设置</h4>

                    <el-form-item
                    :label="$t('repo_url')"
                    prop="repo_url"
                    :rules="[
                        { required: true, message: $t('repo_url_cannot_empty'), trigger: 'blur'},
                    ]">
                        <el-input :placeholder="$t('please_input_repo_url')" v-model="dialogForm.repo_url" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item 
                    :label="$t('repo_branch')"
                    prop="repo_branch">
                        <el-input class="app-input-mini" v-model="dialogForm.repo_branch" autocomplete="off"></el-input>
                        <div class="app-form-explain">若不指定，需要在发起上线时手动填写分支(或Tag)名称</div>
                    </el-form-item>
                    <div class="app-divider"></div>
                    <h4 class="app-form-subtitle">部署设置</h4>

                    <el-form-item 
                    :label="$t('pre_release_cluster')"
                    prop="pre_release_cluster">
                        <el-select 
                        class="app-input-mini"
                        v-model="dialogForm.pre_release_cluster" 
                        filterable 
                        clearable 
                        placeholder="关键词搜索">
                            <el-option
                            v-for="cluster in clusterList"
                            :key="cluster.id"
                            :label="cluster.name"
                            :value="cluster.id">
                            </el-option>
                        </el-select>
                    </el-form-item>

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
                        placeholder="关键词搜索">
                            <el-option
                            v-for="cluster in clusterList"
                            :key="cluster.id"
                            :label="cluster.name"
                            :value="cluster.id">
                            </el-option>
                        </el-select>
                        <div v-if="dialogForm.online_cluster && dialogForm.online_cluster.length">
                            <span>已选集群列表</span>
                            <ul class="app-form-box">
                                <li class="item" v-for="id in dialogForm.online_cluster" :key="id">
                                    <span><i class="iconfont small left icon-cluster"></i>{{ formatClusterName(id) }}</span>
                                    <span>
                                        <el-button @click="removeClusterHandler(id)" icon="el-icon-delete" type="text">移除</el-button>
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

                    <el-form-item 
                    :label="$t('deploy_timeout')"
                    prop="deploy_timeout"
                    :rules="[
                        { required: true, message: $t('deploy_timeout_tips'), trigger: 'blur'},
                    ]">
                        <el-input class="app-input-mini" :placeholder="$t('deploy_timeout')" v-model="dialogForm.deploy_timeout" autocomplete="off"></el-input>
                    </el-form-item>

                </el-form>
                <div slot="footer" class="dialog-footer">
                    <el-button class="app-input-small" size="small" @click="dialogCloseHandler">{{ $t('cancel') }}</el-button>
                    <el-button :loading="btnLoading" size="small" type="primary" @click="dialogSubmitHandler">{{ $t('enter') }}</el-button>
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
import { listSpaceApi, detailSpaceApi, newProjectApi, updateProjectApi, listProjectApi } from '@/api/project'
import { listGroupApi } from '@/api/server' 
export default {
    data() {
        return {
            searchInput: '',
            tableLoading: false,
            tableData: [],

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
                pre_release_cluster: undefined,
                online_cluster: [],
                deploy_user: '',
                deploy_path: '',
                pre_deploy_cmd: '',
                after_deploy_cmd: '',
                deploy_timeout: 120,

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
        searchHandler() {
            this.$root.PageInit()
            this.loadTableData()
        },
        openAddDialogHandler() {
            this.dialogVisible = true
            this.dialogTitle = this.$t('add_project')
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
            listProjectApi({keyword: this.searchInput, offset: this.$root.PageOffset(), limit: this.$root.PageSize}).then(res => {
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
            })
        },
        loadClusterList() {
            listGroupApi({offset: 0, limit: 999}).then(res => {
                if (res.list) {
                    this.clusterList = res.list
                }
            })
        },
    },
    mounted() {
        this.loadSpaceList()
        this.loadClusterList()
        this.spaceId = 7
    }
}
</script>
