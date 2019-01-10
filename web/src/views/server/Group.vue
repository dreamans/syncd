<template>
    <div>
        <a-card
        class="app-card"
        :bordered="false">
            <div class="app-btn-group">
                <a-row :gutter="10">
                    <a-col :span="4">
                        <a-button v-if="$root.CheckPriv($root.Priv.SERVER_GROUP_NEW)" @click="handleOpenAddGroupDialog" type="primary" icon="plus">新增集群</a-button>
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
                    <span v-if="$root.CheckPriv($root.Priv.SERVER_VIEW)" @click="handleShowServerList(record.id)" class="app-link app-op"><a-icon type="bars" />服务器列表</span>
                    <span v-if="$root.CheckPriv($root.Priv.SERVER_GROUP_EDIT)" @click="handleOpenEditGroupDialog(record.id)" class="app-link app-op"><a-icon type="edit" />编辑</span>
                    <a-popconfirm v-if="$root.CheckPriv($root.Priv.SERVER_GROUP_DEL)" title="确定要删除此分组吗？" @confirm="handleDeleteGroup(record.id)" okText="删除" cancelText="取消">
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
                <group-update-component :detail="dialogDetail" ref="groupUpdateRef"></group-update-component>
            </a-spin>
        </a-modal>
    </div>
</template>

<script>
import GroupUpdateComponent from './GroupUpdateComponent.js'
import { newGroupApi, updateGroupApi, getGroupListApi, getGroupDetailApi, deleteGroupApi, getServerListApi } from '@/api/server.js'
export default {
    data() {
        return {
            tableColumns: [
                {dataIndex: "id", title: '集群ID', width: '10%'},
                {dataIndex: "name", title: '集群名称'},
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
            dialogDetail: {},
            dialogLoading: false,

            search: {
                keyword: '',
            },
        }
    },
    components: {
        GroupUpdateComponent,
    },
    methods: {
        handleSearch(value) {
            this.search.keyword = value
            this.pagination.current = 1
            this.handleTableChange(this.pagination)
        },
        handleTableChange(pagination) {
            this.pagination.current = pagination.current
            this.getDataList({
                page: pagination.current,
                pageSize: pagination.pageSize,
            })
        },
        handleOpenAddGroupDialog() {
            this.dialogTitle = '新增集群'
            this.dialogVisible = true
            this.dialogDetail = {}
        },
        handleOpenEditGroupDialog(id) {
            this.dialogTitle = '编辑集群信息'
            this.dialogVisible = true
            this.dialogDetail = {}
            this.getDataDetail(id)
        },
        handleDeleteGroup(id) {
            deleteGroupApi({id}).then(res => {
                this.$message.success('删除成功', 1)
                this.$root.ResetPagination(this.pagination)
                this.handleTableChange(this.pagination)
            })
        },
        handleShowServerList(id) {
            getServerListApi({group_id: id}).then(res => {
                let serverList = res.list ? res.list: []
                let renderServerList = []
                serverList.forEach(s => {
                    let link = '/server/list?op=edit&id=' + s.id
                    renderServerList.push(
                        <div class="item">
                        <span style="display:inline-block; width: 50%">
                            <icon-server /> {s.ip}:[{s.ssh_port}]
                            {this.$root.CheckPriv(this.$root.Priv.SERVER_EDIT) ? (<a class="op" target="_blank" href={link}>编辑</a>) : ''}

                        </span>
                        <span style="display:inline-block; width: 50%">{s.name}</span></div>
                    )
                })
                if (renderServerList.length == 0) {
                    renderServerList = (
                        <div>暂无服务器</div>
                    )
                }
                this.$info({
                    title: '服务器列表',
                    width: "50vw",
                    content: (
                        <div class="app-modal-list">{renderServerList}</div>
                    ),
                })
            })
        },
        dialogCancel() {
            this.dialogVisible = false
        },
        dialogSubmit() {
            this.$refs.groupUpdateRef.validateFields((err, values) => {
                if (err) {
                    return
                }
                this.dialogConfirmLoading = true
                if (this.dialogDetail.id) {
                    updateGroupApi(values).then(res => {
                        this.$message.success('集群更新成功', 1, () => {
                            this.dialogCancel()
                            this.dialogConfirmLoading = false
                            this.handleTableChange(this.pagination)
                        })
                    }).catch(err => {
                        this.dialogConfirmLoading = false
                    })
                } else {
                    newGroupApi(values).then(res => {
                        this.$message.success('集群创建成功', 1, () => {
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
        getDataDetail(id) {
            this.dialogLoading = true
            getGroupDetailApi({id}).then(res => {
                this.dialogLoading = false
                this.dialogDetail = res
            }).catch( err => {
                this.dialogLoading = false
            })
        },
        getDataList(params) {
            this.tableLoading = true
            let offset = (params.page - 1) * params.pageSize
            getGroupListApi({keyword: this.search.keyword, offset: offset, limit: params.pageSize}).then(res => {
                this.tableLoading = false
                this.pagination.total = res.total
                this.tableSource = res.list
            }).catch(err => {
                this.tableLoading = false
            })
        },
    },
    mounted() {
        this.handleTableChange(this.pagination)
    },
}
</script>
