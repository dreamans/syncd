<template>
    <div>
        <el-card shadow="never">
            <el-row class="app-btn-group">
                <el-col :span="4">
                    <el-button @click="openAddDialogHandler" type="primary" size="medium" icon="iconfont left small icon-add">{{ $t('add_server') }}</el-button>
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
                <el-table-column prop="name" label="名称"></el-table-column>
                <el-table-column prop="group_name" width="180" label="集群"></el-table-column>
                <el-table-column prop="ip" width="180" label="IP/HOST"></el-table-column>
                <el-table-column prop="ssh_port" width="180" label="SSH Port"></el-table-column>
                <el-table-column label="操作" width="180" align="right">
                    <template slot-scope="scope">
                        <el-button
                            icon="el-icon-edit"
                            type="text"
                            @click="openEditDialogHandler(scope.row)">编辑</el-button>
                        
                        <el-button
                            type="text"
                            icon="el-icon-delete"
                            class="app-danger"
                            @click="deleteHandler(scope.row)">删除</el-button>
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

        <el-dialog width="500px" :title="dialogTitle" :visible.sync="dialogVisible">
            <div class="app-dialog" v-loading="dialogLoading">
                <el-form class="app-form" ref="dialogRef" :model="dialogForm" size="medium" label-width="120px">
                    <el-form-item 
                        label="所属集群"
                        prop="group_id"
                        :rules="[
                            { required: true, message: '所属集群不能为空'},
                        ]">
                        <el-select filterable placeholder="关键字搜索" style="width: 100%" v-model="dialogForm.group_id">
                            <el-option
                                v-for="g in serverGroupList"
                                :key="g.id"
                                :label="g.name"
                                :value="g.id">
                            </el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item 
                        label="服务器名称"
                        prop="name"
                        :rules="[
                            { required: true, message: '名称不能为空'},
                        ]">
                        <el-input placeholder="请输入服务器名称" v-model="dialogForm.name" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item 
                        label="IP/HOST"
                        prop="ip"
                        :rules="[
                            { required: true, message: 'IP/HOST不能为空'},
                        ]">
                        <el-input placeholder="请输入服务器IP/HOST" v-model="dialogForm.ip" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item 
                        label="SSH端口"
                        prop="ssh_port"
                        :rules="[
                            { required: true, message: 'SSH端口不能为空'},
                        ]">
                        <el-input maxlength=5 class="app-input-mini" v-model="dialogForm.ssh_port" autocomplete="off"></el-input>
                    </el-form-item>
                </el-form>
                <div slot="footer" class="dialog-footer">
                    <el-button size="small" @click="dialogCloseHandler">取 消</el-button>
                    <el-button :loading="btnLoading" size="small" type="primary" @click="dialogSubmitHandler">确 定</el-button>
                </div>
            </div>
        </el-dialog>

    </div>
</template>

<script>
import { listGroupApi, newServerApi, updateServerApi, listServerApi, deleteServerApi } from '@/api/server'
export default {
    data() {
        return {
            searchInput: '',
            dialogVisible: false,
            dialogTitle: '',
            dialogForm: {
                id: 0,
                group_id: undefined,
                name: '',
                ip: '',
                ssh_port: 22,
            },
            dialogLoading: false,
            btnLoading: false,

            tableData: [],
            tableLoading: false,
            deletePopover: false,

            serverGroupList: [],
        }
    },
    methods: {
        searchHandler() {
            this.$root.PageInit()
            this.loadTableData()
        },
        openAddDialogHandler() {
            this.dialogVisible = true
            this.dialogTitle = '新增服务器'
        },
        openEditDialogHandler(row) {
            this.dialogVisible = true
            this.dialogTitle = '编辑服务器信息'
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
                    opFn = updateServerApi
                } else {
                    opFn = newServerApi
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
            listServerApi({keyword: this.searchInput, offset: this.$root.PageOffset(), limit: this.$root.PageSize}).then(res => {
                this.tableData = res.list
                this.$root.Total = res.total
                this.tableLoading = false
            }).catch(err => {
                this.tableLoading = false
            })
        },
        loadServerGroupList() {
            listGroupApi({offset: 0, limit: 999}).then(res => {
                if (res.list && res.list.length > 0) {
                    this.serverGroupList = res.list
                }
            })
        },
    },
    mounted() {
        this.$root.PageInit()
        this.loadTableData()
        this.loadServerGroupList()
    }
}
</script>
