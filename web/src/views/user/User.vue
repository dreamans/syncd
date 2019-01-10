<template>
    <div>
        <a-card
        class="app-card"
        :bordered="false">
            <div class="app-btn-group">
                <a-row :gutter="10">
                    <a-col :span="4">
                        <a-button v-if="$root.CheckPriv($root.Priv.USER_NEW)" @click="handleOpenAddDialog" type="primary" icon="plus">新增用户</a-button>
                    </a-col>
                    <a-col :span="12"></a-col>
                    <a-col :span="8">
                        <a-input-search v-model="search.keyword" placeholder="关键词搜索，ID、手机号、用户名" @search="handleSearch" enterButton/>
                    </a-col>
                </a-row>
            </div>
            <a-table
            :columns="tableColumns"
            :dataSource="tableSource"
            :pagination="pagination"
            @change="handleTableChange"
            :loading="tableLoading">
                <span slot="last_login_time" slot-scope="text, record">
                    <a-tooltip placement="top">
                        <template slot="title">
                            <div>Time {{ $root.FormatDateTime(text) }}</div>
                            <div>IP {{ record.last_login_ip }}</div>
                        </template>
                        {{$root.FormatDateFromNow(text)}}
                    </a-tooltip>
                </span>
                <span slot="lock_status" slot-scope="text, record">
                    <template v-if="text == 1">
                        <a-tooltip placement="top" >
                            <template slot="title">
                                <span>用户可正常登陆</span>
                            </template>
                            <span class="app-color-success"><a-icon type="unlock" /> 正常</span>
                        </a-tooltip>
                    </template>
                    <template v-if="text == 0">
                        <a-tooltip placement="top">
                            <template slot="title">
                                <span>用户被锁定，禁止登录</span>
                            </template>
                            <span class="app-color-gray"><a-icon type="lock" /> 锁定</span>
                        </a-tooltip>
                    </template>
                </span>
                <span slot="op" slot-scope="text, record">
                    <span v-if="$root.CheckPriv($root.Priv.USER_EDIT)" @click="handleOpenEditDialog(record.id)" class="app-link app-op"><a-icon type="edit" />编辑</span>
                    <a-popconfirm v-if="$root.CheckPriv($root.Priv.USER_DEL)" title="确定要删除此用户吗？" @confirm="handleDeleteUser(record.id)" okText="删除" cancelText="取消">
                        <span class="app-link app-op app-remove"><a-icon type="delete" />删除</span>
                    </a-popconfirm>
                </span>
            </a-table>
        </a-card>

        <a-modal
        :title="dialogTitle"
        :visible="dialogVisible"
        @ok="dialogSubmit"
        :confirmLoading="dialogConfirmLoading"
        :keyboard="false"
        :maskClosable="false"
        okText="确定"
        cancelText="取消"
        :destroyOnClose="true"
        @cancel="dialogCancel">
            <a-spin :spinning="dialogLoading">
                <user-update-component :detail="dialogDetail" ref="updateRef"></user-update-component>
            </a-spin>
        </a-modal>
    </div>
</template>

<script>
import { newUserApi, updateUserApi, getUserListApi, getUserDetailApi, deleteUserApi } from '@/api/user.js'
import UserUpdateComponent from './UserUpdateComponent.js'
export default {
    data () {
        return {
            search: {},
            tableColumns: [
                {dataIndex: "name", title: '用户名'},
                {dataIndex: "email", title: '邮箱'},
                {dataIndex: "group_name", title: '角色', width: '15%'},
                {dataIndex: "lock_status", title: '状态', width: '10%', scopedSlots: { customRender: 'lock_status' }},
                {dataIndex: "last_login_time", title: '上次登录', width: '15%', scopedSlots: { customRender: 'last_login_time' }},
                {dataIndex: "op", title: '操作', width: '15%', align: 'right', scopedSlots: { customRender: 'op' }},
            ],
            tableSource: [],
            pagination: {
                pageSize: 10,
                current: 1,
                total: 0,
            },
            tableLoading: false,
            userGroupList: {},

            dialogTitle: '',
            dialogVisible: false,
            dialogConfirmLoading: false,
            dialogLoading: false,
            dialogDetail: {},
        }
    },
    components: {
        UserUpdateComponent,
    },
    methods: {
        handleTableChange(pagination) {
            this.pagination.current = pagination.current
            this.getDataList({
                page: pagination.current,
                pageSize: pagination.pageSize,
            })
        },
        handleOpenAddDialog() {
            this.dialogTitle = '新增用户'
            this.dialogVisible = true
            this.dialogDetail = {}
        },
        handleOpenEditDialog(id) {
            this.dialogTitle = '编辑用户'
            this.dialogVisible = true
            this.dialogDetail = {}
            this.getDataDetail(id)
        },
        handleSearch(value) {
            this.search.keyword = value
            this.pagination.current = 1
            this.handleTableChange(this.pagination)
        },
        handleDeleteUser(id) {
            deleteUserApi({id}).then(res => {
                this.$message.success('删除成功', 1)
                this.$root.ResetPagination(this.pagination)
                this.handleTableChange(this.pagination)
            })
        },
        dialogSubmit() {
            this.$refs.updateRef.validateFields((err, values) => {
                if (err) {
                    return
                }
                this.dialogConfirmLoading = true
                if (values.password) {
                    values.password = this.$root.Md5Sum(values.password)
                }
                values.lock_status = values.lock_status ? 1: 0
                if (this.dialogDetail.id) {
                    updateUserApi(values).then(res => {
                        this.$message.success('用户信息更新成功', 1, () => {
                            this.dialogCancel()
                            this.dialogConfirmLoading = false
                            this.handleTableChange(this.pagination)
                        })
                    }).catch(err => {
                        this.dialogConfirmLoading = false
                    })
                } else {
                    newUserApi(values).then(res => {
                        this.$message.success('用户创建成功', 1, () => {
                            this.dialogCancel()
                            this.dialogConfirmLoading = false
                            this.handleTableChange(this.pagination)
                        })
                    }).catch(err => {
                        this.dialogConfirmLoading = false
                    })
                }
            })
        },
        dialogCancel() {
            this.dialogVisible = false
        },
        getDataList(params) {
            this.tableLoading = true
            let offset = (params.page - 1) * params.pageSize
            let keyword = this.search.keyword
            getUserListApi({keyword: keyword, offset: offset, limit: params.pageSize}).then(res => {
                this.tableLoading = false
                this.pagination.total = res.total
                this.tableSource = res.list
            }).catch(err => {
                this.tableLoading = false
            })
        },
        getDataDetail(id) {
            this.dialogLoading = true
            getUserDetailApi({id}).then(res => {
                this.dialogLoading = false
                this.dialogDetail = res
            }).catch( err => {
                this.dialogLoading = false
            })
        },
    },
    mounted() {
        this.handleTableChange(this.pagination)
    },
}
</script>
