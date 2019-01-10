<template>
    <div>
        <a-card
        class="app-card"
        :bordered="false">
            <div class="app-btn-group">
                <a-row :gutter="10">
                    <a-col :span="4">
                        <template v-if="$root.CheckPriv($root.Priv.USER_ROLE_NEW)">
                            <a-button @click="handleOpenAddDialog" type="primary" icon="plus">新增角色</a-button>
                        </template>
                    </a-col>
                    <a-col :span="14"></a-col>
                    <a-col :span="6">
                        <a-input-search v-model="search.keyword" placeholder="关键词搜索，ID、名称" @search="handleSearch" enterButton/>
                    </a-col>
                </a-row>
            </div>
            <a-table
            :columns="tableColumns"
            :dataSource="tableSource"
            :pagination="pagination"
            @change="handleTableChange"
            :loading="tableLoading">
                <span slot="op" slot-scope="text, record">
                    <span v-if="$root.CheckPriv($root.Priv.USER_ROLE_EDIT)" @click="handleOpenEditDialog(record.id)" class="app-link app-op"><a-icon type="edit" />编辑</span>
                    <a-popconfirm v-if="$root.CheckPriv($root.Priv.USER_ROLE_DEL)" title="确定要删除此服务器吗？" @confirm="handleDeleteGroup(record.id)" okText="删除" cancelText="取消">
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
        width="55vw"
        :destroyOnClose="true"
        @cancel="dialogCancel">
            <a-spin :spinning="dialogLoading">
                <group-update-component :detail="dialogDetail" ref="updateRef"></group-update-component>
            </a-spin>
        </a-modal>
    </div>
</template>

<script>
import { newGroupApi, updateGroupApi, getGroupListApi, getGroupDetailApi, deleteGroupApi } from '@/api/user.js'
import GroupUpdateComponent from './GroupUpdateComponent.js'
export default {
    data () {
        return {
            search: {},
            tableColumns: [
                {dataIndex: "id", title: 'ID', width: '10%'},
                {dataIndex: "name", title: '名称'},
                {dataIndex: "op", title: '操作', width: '30%', align: 'right', scopedSlots: { customRender: 'op' }},
            ],
            tableSource: [],
            pagination: {
                pageSize: 10,
                current: 1,
                total: 0,
            },
            tableLoading: false,

            dialogTitle: '',
            dialogVisible: false,
            dialogConfirmLoading: false,
            dialogLoading: false,
            dialogDetail: {},
        }
    },
    components: {
        GroupUpdateComponent,
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
            this.dialogTitle = '新增角色'
            this.dialogVisible = true
            this.dialogDetail = {}
        },
        handleOpenEditDialog(id) {
            this.dialogTitle = '编辑服务器信息'
            this.dialogVisible = true
            this.dialogDetail = {}
            this.getDataDetail(id)
        },
        handleSearch(value) {
            this.search.keyword = value
            this.pagination.current = 1
            this.handleTableChange(this.pagination)
        },
        handleDeleteGroup(id) {
            deleteGroupApi({id}).then(res => {
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
                if (this.dialogDetail.id) {
                    updateGroupApi(values).then(res => {
                        this.$message.success('角色信息更新成功', 1, () => {
                            this.dialogCancel()
                            this.dialogConfirmLoading = false
                            this.handleTableChange(this.pagination)
                        })
                    }).catch(err => {
                        this.dialogConfirmLoading = false
                    })
                } else {
                    newGroupApi(values).then(res => {
                        this.$message.success('角色信息创建成功', 1, () => {
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
            getGroupListApi({keyword: keyword, offset: offset, limit: params.pageSize}).then(res => {
                this.tableLoading = false
                this.pagination.total = res.total
                this.tableSource = res.list
            }).catch(err => {
                this.tableLoading = false
            })
        },
        getDataDetail(id) {
            this.dialogLoading = true
            getGroupDetailApi({id}).then(res => {
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
