<template>
    <div>
        <el-card shadow="never">
            <el-row class="app-btn-group">
                <el-col :span="4">
                    <el-button v-if="$root.CheckPriv($root.Priv.USER_ROLE_NEW)" @click="openAddDialogHandler" type="primary" size="medium" icon="iconfont left small icon-add">{{ $t('add_role') }}</el-button>&nbsp;
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
                <el-table-column prop="name" :label="$t('name')"></el-table-column>
                <el-table-column :label="$t('operate')" width="180" align="right">
                    <template slot-scope="scope">
                        <el-button
                        v-if="$root.CheckPriv($root.Priv.USER_ROLE_EDIT)"
                        icon="el-icon-edit"
                        type="text"
                        @click="openEditDialogHandler(scope.row)">{{ $t('edit') }}</el-button>
                        <el-button
                        v-if="$root.CheckPriv($root.Priv.USER_ROLE_DEL)"
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

        <el-dialog :width="$root.DialogNormalWidth" :title="dialogTitle" :visible.sync="dialogVisible" @close="dialogCloseHandler">
            <div class="app-dialog" v-loading="dialogLoading">
                <el-form ref="dialogRef" :model="dialogForm" size="medium" label-width="80px">
                    <el-form-item 
                    :label="$t('role_name')"
                    prop="name"
                    :rules="[
                        { required: true, message: $t('name_cannot_empty'), trigger: 'blur'},
                    ]">
                        <el-input v-model="dialogForm.name" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item 
                    :label="$t('privilege_setting')"
                    prop="privilege">
                        <el-checkbox :indeterminate="isIndeterminate" v-model="privCheckAll" @change="checkAllCheckHandler">{{ $t('check_all') }}</el-checkbox>
                        <el-checkbox-group class="app-checkbox-group" v-model="dialogForm.privilege">
                            <template v-for="(privGroup, index) in privilegeList">
                                <el-row :key="index" class="app-mt-line">
                                    <el-col :span="3">
                                        <span class="app-label">{{ privGroup.label }}</span>
                                    </el-col>
                                    <el-col :span="21">
                                        <el-checkbox v-for="priv in privGroup.items" :label="priv.value" :key="priv.value">{{priv.label}}</el-checkbox>
                                    </el-col>
                                </el-row>
                            </template>
                        </el-checkbox-group>
                    </el-form-item>
                </el-form>
                <div slot="footer" class="dialog-footer">
                    <el-button size="small" @click="dialogCloseHandler">{{ $t('cancel') }}</el-button>
                    <el-button :loading="btnLoading" size="small" type="primary" @click="dialogSubmitHandler">{{ $t('enter') }}</el-button>
                </div>
            </div>
        </el-dialog>

    </div>
</template>

<script>
import { privListApi, newRoleApi, listRoleApi, detailRoleApi, updateRoleApi, deleteRoleApi } from '@/api/user'
export default {
    data() {
        return {
            searchInput: '',
            dialogVisible: false,
            dialogTitle: '',
            dialogForm: {
                id: 0,
                name: '',
                privilege: [],
            },
            dialogLoading: false,
            btnLoading: false,

            tableData: [],
            tableLoading: false,

            isIndeterminate: false,
            privCheckAll: false,

            privilegeList: [],
        }
    },
    watch: {
        'dialogForm.privilege'(val) {
            this.checkedChange(val)
        },
    },
    methods: {
        checkAllCheckHandler(val) {
            let privList = []
            if (val) {
                this.privilegeList.forEach(privGroup => {
                    privGroup.items.forEach(p => {
                        privList.push(p.value)
                    })
                })
            } else {
                privList = []
            }
            this.dialogForm.privilege = privList
            this.isIndeterminate = false
        },
        checkedChange(val) {
            this.isIndeterminate = false
            let checkAll = true
            this.privilegeList.forEach(privGroup => {
                privGroup.items.forEach(p => {
                    if (val.indexOf(p.value) == -1) {
                        checkAll = false
                    }
                })
            })
            if (checkAll) {
                this.privCheckAll = true
            } else {
                this.privCheckAll = false
                if (val && val.length > 0) {
                    this.isIndeterminate = true
                } else {
                    this.isIndeterminate = false
                }
            }
        },
        searchHandler() {
            this.$root.PageInit()
            this.loadTableData()
        },
        openAddDialogHandler() {
            this.dialogVisible = true
            this.loadPrivList()
            this.dialogTitle = this.$t('add_role')
        },
        openEditDialogHandler(row) {
            this.dialogVisible = true
            this.loadPrivList()
            this.dialogTitle = this.$t('edit_role_info')
            this.dialogLoading = true
            detailRoleApi({id: row.id}).then(res => {
                this.dialogLoading = false
                this.dialogForm = {
                    id: res.id,
                    name: res.name,
                    privilege: res.privilege ? res.privilege : [],
                }
            }).catch(err => {
                this.dialogCloseHandler()
            })
        },
        dialogCloseHandler() {
            this.dialogVisible = false
            this.dialogLoading = false
            this.btnLoading = false
            this.$refs.dialogRef.resetFields();
            this.dialogForm = {
                id: 0,
                name: '',
                privilege: [],
            }
        },
        deleteHandler(row) {
            this.$root.ConfirmDelete(() => {
                deleteRoleApi({id: row.id}).then(res => {
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
                    opFn = updateRoleApi
                } else {
                    opFn = newRoleApi
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
        loadTableData() {
            this.tableLoading = true
            listRoleApi({keyword: this.searchInput, offset: this.$root.PageOffset(), limit: this.$root.PageSize}).then(res => {
                this.tableData = res.list
                this.$root.Total = res.total
                this.tableLoading = false
            }).catch(err => {
                this.tableLoading = false
            })
        },
        loadPrivList() {
            privListApi().then(res => {
                this.privilegeList = res
            })
        },
    },
    mounted() {
        this.$root.PageInit()
        this.loadTableData()
    }
}
</script>
