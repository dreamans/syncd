<template>
    <div>
        <el-card shadow="never">
            <el-row class="app-btn-group">
                <el-col :span="4">
                    <el-button @click="openAddDialogHandler" type="primary" size="medium" icon="iconfont left small icon-add">{{ $t('add_user') }}</el-button>
                </el-col>
                <el-col :span="6" :offset="14">
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
                <el-table-column prop="username" :label="$t('username')"></el-table-column>
                <el-table-column prop="role_name" width="200" :label="$t('role')"></el-table-column>
                <el-table-column prop="email" width="240" :label="$t('email')"></el-table-column>
                <el-table-column prop="status" width="100" :label="$t('status')">
                    <template slot-scope="scope">
                        <span class="app-color-success" v-if="scope.row.status == '1'">
                            <i class="iconfont icon-unlock"></i> 正常
                        </span>
                        <span class="app-color-error" v-else>
                            <i class="iconfont icon-lock"></i> 锁定
                        </span>
                    </template>
                </el-table-column>
                <el-table-column width="150" :label="$t('last_login')">
                    <template slot-scope="scope">
                        <el-tooltip placement="top">
                            <div slot="content">
                                上次登录时间: {{ $root.FormatDateTime(scope.row.last_login_time) }}
                                <br/><br/>
                                上次登录IP: {{ scope.row.last_login_ip }}
                            </div>
                            <span>{{ $root.FormatDateFromNow(scope.row.last_login_time) }}</span>
                        </el-tooltip>
                    </template>
                </el-table-column>
                <el-table-column :label="$t('operate')" width="150" align="right">
                    <template slot-scope="scope">
                        <el-button
                            icon="el-icon-edit"
                            type="text"
                            @click="openEditDialogHandler(scope.row)">{{ $t('edit') }}</el-button>
                        
                        <el-button
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

        <el-dialog width="500px" :title="dialogTitle" :visible.sync="dialogVisible" @close="dialogCloseHandler">
            <div class="app-dialog" v-loading="dialogLoading">
                <el-form class="app-form" ref="dialogRef" :model="dialogForm" size="medium" label-width="80px">
                    <el-form-item 
                    :label="$t('role')"
                    prop="role_id"
                    :rules="[
                        { required: true, message: this.$t('role_cannot_empty')},
                    ]">
                        <el-select filterable :placeholder="$t('keyword_search')" style="width: 100%" v-model="dialogForm.role_id">
                            <el-option
                                v-for="g in roleList"
                                :key="g.id"
                                :label="g.name"
                                :value="g.id">
                            </el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item 
                    :label="$t('username')"
                    prop="username"
                    :rules="[
                        { required: true, message: this.$t('username_cannot_empty'), trigger: 'blur'},
                        { validator: this.userExistsValid('username'), trigger: 'blur' }
                    ]">
                        <el-input :placeholder="$t('please_input_username')" v-model="dialogForm.username" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item 
                    :label="$t('password')"
                    prop="password"
                    :rules="[
                        { required: true, message: this.$t('password_cannot_empty'), trigger: 'blur'},
                        { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur'}
                    ]">
                        <el-input type="password" :placeholder="$t('please_input_password_length_limit')" v-model="dialogForm.password" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item 
                    :label="$t('email')"
                    prop="email"
                    :rules="[
                        { required: true, message: this.$t('email_cannot_empty'), trigger: 'blur'},
                        { type:'email', message: '邮箱格式错误', trigger: 'blur'},
                        { validator: this.userExistsValid('email'), trigger: 'blur' },
                    ]">
                        <el-input :placeholder="$t('please_input_email')" v-model="dialogForm.email" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item 
                    :label="$t('truename')"
                    prop="truename"
                    >
                        <el-input v-model="dialogForm.truename" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item 
                    :label="$t('mobile')"
                    prop="mobile"
                    >
                        <el-input v-model="dialogForm.mobile" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item 
                    :label="$t('allow_login')"
                    prop="mobile"
                    >
                        <div>
                            <el-switch
                            v-model="dialogForm.status"
                            active-value="1"
                            inactive-value="0"
                            active-color="#13ce66">
                            </el-switch>
                            <span style="line-height: 20px; display: inline-flex; align-items: center; vertical-align: middle; margin-left: 5px;">
                                <i v-if="dialogForm.status == 1" class="iconfont icon-unlock"></i>
                                <i v-else class="iconfont icon-lock"></i>
                            </span>
                        </div>
                        <div class="app-form-explain">禁止后用户将无法登录</div>
                    </el-form-item>
                </el-form>
                <div slot="footer" class="dialog-footer">
                    <el-button size="small" @click="dialogCloseHandler">{{ $t('cancel')}}</el-button>
                    <el-button :loading="btnLoading" size="small" type="primary" @click="dialogSubmitHandler">{{ $t('enter')}}</el-button>
                </div>
            </div>
        </el-dialog>

    </div>
</template>

<script>
import { listRoleApi, newUserApi, updateUserApi, listUserApi, existsUserApi } from '@/api/user'
export default {
    data() {
        return {
            searchInput: '',
            dialogVisible: false,
            dialogTitle: '',
            dialogForm: {
                id: 0,
                role_id: undefined,
                username: '',
                password: '',
                email: '',
                truename: '',
                mobile: '',
                status: '1',
            },
            dialogLoading: false,
            btnLoading: false,

            tableData: [],
            tableLoading: false,

            roleList: [],
        }
    },
    methods: {
        userExistsValid(field) {
            let vm = this
            return function(rule, value, callback) {
                if (!value) {
                    callback()
                    return
                }
                let errmsg = '数据重复，请重新输入'
                let loadingmsg = '验证中，请稍后'
                let query = {id: vm.dialogForm.id}
                switch (field) {
                    case 'username':
                        errmsg = '用户名已经存在，请重新输入'
                        loadingmsg = '正在验证用户名是否被占用，请稍等'
                        query.username = value
                        break
                    case 'email':
                        errmsg = '邮箱已经存在，请重新输入'
                        loadingmsg = '正在验证邮箱是否被占用，请稍等'
                        query.email = value
                        break
                }
                let modal = vm.$message({
                    iconClass: 'el-icon-loading',
                    message: loadingmsg,
                    duration: 0,
                });
                existsUserApi(query).then(res => {
                    modal.close()
                    if (!res.exists) {
                        callback()
                    } else {
                        callback(errmsg)
                    }
                }).catch(err => {
                    modal.close()
                    callback('网络错误, 校验失败')
                })
            }
        },
        searchHandler() {
            this.$root.PageInit()
            this.loadTableData()
        },
        openAddDialogHandler() {
            this.dialogVisible = true
            this.dialogTitle = this.$t('add_user')
        },
        openEditDialogHandler(row) {
            this.dialogVisible = true
            this.dialogTitle = this.$t('edit_server_info')
            this.dialogLoading = true
            detailServerApi({id: row.id}).then(res => {
                this.dialogLoading = false

                this.dialogForm = res
            }).catch(err => {
                this.dialogCloseHandler()
            })
        },
        dialogCloseHandler() {
            this.dialogVisible = false
            this.dialogLoading = false
            this.btnLoading = false
            this.$refs.dialogRef.resetFields();
            this.dialogForm.ssh_port = 22
        },
        deleteHandler(row) {
            this.$root.ConfirmDelete(() => {
                deleteServerApi({id: row.id}).then(res => {
                    this.$root.MessageSuccess()
                    this.$root.PageReset()
                    this.loadTableData()
                })
            })
        },
        currentChangeHandler() {
            this.loadTableData()
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
                    opFn = updateUserApi
                } else {
                    opFn = newUserApi
                }

                let postData = Object.assign({}, this.dialogForm)
                if (postData.password) {
                    postData.password = this.$root.Md5Sum(postData.password)
                }
                opFn(postData).then(res => {
                    this.$root.MessageSuccess(() => {
                        this.dialogCloseHandler()
                        this.btnLoading = false
                        this.loadTableData()
                    })
                }).catch(err => {
                    this.$message.error(err.message)
                    this.btnLoading = false
                })
            });
        },
        loadTableData() {
            this.tableLoading = true
            listUserApi({keyword: this.searchInput, offset: this.$root.PageOffset(), limit: this.$root.PageSize}).then(res => {
                this.tableData = res.list
                this.$root.Total = res.total
                this.tableLoading = false
            }).catch(err => {
                this.tableLoading = false
            })
        },
        loadRoleList() {
            listRoleApi({offset: 0, limit: 999}).then(res => {
                if (res.list && res.list.length > 0) {
                    this.roleList = res.list
                }
            })
        },
    },
    mounted() {
        this.$root.PageInit()
        this.loadTableData()
        this.loadRoleList()
    }
}
</script>
