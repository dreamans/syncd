<template>
    <div>
        <el-card shadow="never">
            <el-row class="app-btn-group">
                <el-col :span="4">
                    <el-button v-if="$root.CheckPriv($root.Priv.SERVER_GROUP_NEW)" @click="openAddDialogHandler" type="primary" size="medium" icon="iconfont left small icon-add">{{ $t('add_cluster') }}</el-button>&nbsp;
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
                        v-if="$root.CheckPriv($root.Priv.SERVER_GROUP_EDIT)"
                        icon="el-icon-edit"
                        type="text"
                        @click="openEditDialogHandler(scope.row)">{{ $t('edit') }}</el-button>
                        <el-button
                        v-if="$root.CheckPriv($root.Priv.SERVER_GROUP_DEL)"
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

        <el-dialog :width="$root.DialogSmallWidth" :title="dialogTitle" :visible.sync="dialogVisible" @close="dialogCloseHandler">
            <div class="app-dialog" v-loading="dialogLoading">
                <el-form ref="dialogRef" :model="dialogForm" size="medium" label-width="80px">
                    <el-form-item 
                        :label="$t('server_group_name')"
                        prop="name"
                        :rules="[
                            { required: true, message: $t('name_cannot_empty'), trigger: 'blur'},
                        ]">
                        <el-input v-model="dialogForm.name" autocomplete="off"></el-input>
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
import { newGroupApi, updateGroupApi, listGroupApi, deleteGroupApi, detailGroupApi } from '@/api/server'
export default {
    data() {
        return {
            searchInput: '',
            dialogVisible: false,
            dialogTitle: '',
            dialogForm: {},
            dialogLoading: false,
            btnLoading: false,

            tableData: [],
            tableLoading: false,
        }
    },
    methods: {
        searchHandler() {
            this.$root.PageInit()
            this.loadTableData()
        },
        openAddDialogHandler() {
            this.dialogVisible = true
            this.dialogTitle = this.$t('add_cluster')
        },
        openEditDialogHandler(row) {
            this.dialogVisible = true
            this.dialogTitle = this.$t('edit_cluster')
            this.dialogLoading = true
            detailGroupApi({id: row.id}).then(res => {
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
            this.dialogForm = {}
        },
        deleteHandler(row) {
            this.$root.ConfirmDelete(() => {
                deleteGroupApi({id: row.id}).then(res => {
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
                    opFn = updateGroupApi
                } else {
                    opFn = newGroupApi
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
            listGroupApi({keyword: this.searchInput, offset: this.$root.PageOffset(), limit: this.$root.PageSize}).then(res => {
                this.tableData = res.list
                this.$root.Total = res.total
                this.tableLoading = false
            }).catch(err => {
                this.tableLoading = false
            })
        }
    },
    mounted() {
        this.$root.PageInit()
        this.loadTableData()
    }
}
</script>
