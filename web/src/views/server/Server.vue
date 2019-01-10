<template>
    <div>
        <a-card
        class="app-card"
        :bordered="false">
            <div class="app-btn-group">
                <a-row :gutter="10">
                    <a-col :span="4">
                        <a-button v-if="$root.CheckPriv($root.Priv.SERVER_NEW)" @click="handleOpenAddServerDialog" type="primary" icon="plus">新增服务器</a-button>
                    </a-col>
                    <a-col :span="10"></a-col>
                    <a-col :span="4">
                        <a-select
                        v-model="search.groupId"
                        allowClear
                        showSearch
                        placeholder="搜索指定集群"
                        notFoundContent="无集群数据"
                        style="width: 100%"
                        optionFilterProp="children">
                            <a-select-option v-for="group in dialogGroupList" :value="group.id">{{ group.name }}</a-select-option>
                        </a-select>
                    </a-col>
                    <a-col :span="6">
                        <a-input-search v-model="search.keyword" placeholder="关键词搜索，ID、名称、IP" @search="handleSearch" enterButton/>
                    </a-col>
                </a-row>
            </div>
            <a-table
            :columns="tableColumns"
            :dataSource="tableSource"
            :pagination="pagination"
            @change="handleTableChange"
            :loading="tableLoading">
                <template slot="group_id" slot-scope="text, record">
                    <span v-if="groupList[text]">{{ groupList[text].name }}</span>
                    <span v-else><span class="app-line-through">已删除</span></span>
                </template>
                <span slot="op" slot-scope="text, record">
                    <span v-if="$root.CheckPriv($root.Priv.SERVER_EDIT)" @click="handleOpenEditServerDialog(record.id)" class="app-link app-op"><a-icon type="edit" />编辑</span>
                    <a-popconfirm v-if="$root.CheckPriv($root.Priv.SERVER_DEL)" title="确定要删除此服务器吗？" @confirm="handleDeleteServer(record.id)" okText="删除" cancelText="取消">
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
                <server-update-component :group-list="dialogGroupList" :detail="dialogDetail" ref="groupUpdateRef"></server-update-component>
            </a-spin>
        </a-modal>
    </div>
</template>

<script>
import ServerUpdateComponent from './ServerUpdateComponent.js'
import { newServerApi, updateServerApi, getServerListApi, getServerDetailApi, deleteServerApi, getGroupListApi } from '@/api/server.js'
export default {
    data() {
        return {
            tableColumns: [
                {dataIndex: "id", title: 'ID', width: '6%'},
                {dataIndex: "name", title: '名称'},
                {dataIndex: "group_id", title: '集群', width: '20%', scopedSlots: { customRender: 'group_id' }},
                {dataIndex: "ip", title: 'IP', width: '15%'},
                {dataIndex: "ssh_port", title: 'SSH Port', width: '10%'},
                {dataIndex: "op", title: '操作', width: '20%', align: 'right', scopedSlots: { customRender: 'op' }},
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
            dialogGroupList: [],
            dialogLoading: false,

            search: {
                keyword: '',
                groupId: undefined,
            },
        }
    },
    components: {
        ServerUpdateComponent,
    },
    computed: {
        groupList() {
            let newGroupList = {}
            this.dialogGroupList.forEach(g => {
                newGroupList[g.id] = g
            })
            return newGroupList
        },
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
        handleOpenAddServerDialog() {
            this.dialogTitle = '新增服务器'
            this.dialogVisible = true
            this.dialogDetail = {}
        },
        handleOpenEditServerDialog(id) {
            this.dialogTitle = '编辑服务器信息'
            this.dialogVisible = true
            this.dialogDetail = {}
            this.getDataDetail(id)
        },
        handleDeleteServer(id) {
            deleteServerApi({id}).then(res => {
                this.$message.success("删除成功", 1)
                this.$root.ResetPagination(this.pagination)
                this.handleTableChange(this.pagination)
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
                    updateServerApi(values).then(res => {
                        this.$message.success('服务器信息更新成功', 1, () => {
                            this.dialogCancel()
                            this.dialogConfirmLoading = false
                            this.handleTableChange(this.pagination)
                        })
                    }).catch(err => {
                        this.dialogConfirmLoading = false
                    })
                } else {
                    newServerApi(values).then(res => {
                        this.$message.success('服务器信息创建成功', 1, () => {
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
            getServerDetailApi({id}).then(res => {
                this.dialogLoading = false
                this.dialogDetail = res
            }).catch( err => {
                this.dialogLoading = false
            })
        },
        getDataList(params) {
            this.tableLoading = true
            let offset = (params.page - 1) * params.pageSize
            let groupId = this.search.groupId ? this.search.groupId: 0
            let keyword = this.search.keyword
            getServerListApi({group_id: groupId, keyword: keyword, offset: offset, limit: params.pageSize}).then(res => {
                this.tableLoading = false
                this.pagination.total = res.total
                this.tableSource = res.list
            }).catch(err => {
                this.tableLoading = false
            })
        },
        getDataGroupList() {
            getGroupListApi({offset: 0, limit: 9999}).then(res => {
                if (res.list) {
                    this.dialogGroupList = res.list
                }
            })
        },
    },
    mounted() {
        this.getDataGroupList()
        let op = this.$route.query.op
        let id = this.$route.query.id
        if (op == 'edit' && id > 0) {
            this.handleOpenEditServerDialog(id)
            this.handleSearch(id)
        } else {
            this.handleTableChange(this.pagination)
        }
    },
}
</script>
