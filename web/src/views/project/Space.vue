<template>
    <div>
        <a-card
        class="app-card"
        title="项目空间"
        :bordered="false">
            <div class="app-btn-group">
                <a-row :gutter="10">
                    <a-col :span="4">
                        <a-button v-if="$root.CheckPriv($root.Priv.PROJECT_SPACE_NEW)" @click="handleOpenAddDialog" type="primary" icon="plus">新增</a-button>
                    </a-col>
                    <a-col :span="14"></a-col>
                    <a-col :span="6">
                        <a-input-search v-model="search.keyword" placeholder="按名称搜索" @search="handleSearch" enterButton/>
                    </a-col>
                </a-row>
            </div>
            <a-table
            :columns="tableColumns"
            :dataSource="tableSource"
            :pagination="pagination"
            @change="handleTableChange"
            :loading="tableLoading">
                <span class="app-content-list" slot="name" slot-scope="text, record">
                    <div class="title">{{ record.name }}</div>
                    <div v-if="record.description" class="description">{{ record.description }}</div>
                </span>
                <span slot="op" slot-scope="text, record">
                    <span v-if="$root.CheckPriv($root.Priv.PROJECT_VIEW)" @click="$root.GotoRouter('projectProject', {space: record.id})" class="app-link app-op"><icon-project />项目管理</span>
                    <span v-if="$root.CheckPriv($root.Priv.PROJECT_USER_VIEW)" @click="$root.GotoRouter('projectUser', {space: record.id})" class="app-link app-op"><a-icon type="team" />成员管理</span>
                    <span v-if="$root.CheckPriv($root.Priv.PROJECT_SPACE_EDIT)" @click="handleOpenUpdateDialog(record.id)" class="app-link app-op"><a-icon type="edit" />编辑</span>
                    <template v-if="$root.CheckPriv($root.Priv.PROJECT_SPACE_DEL)">
                        <template v-if="record.have_project">
                            <a-tooltip placement="topRight" >
                                <template slot="title">
                                    <span>项目列表不为空，禁止删除</span>
                                </template>
                                <span class="app-op app-color-gray app-no-allow"><a-icon type="delete" />删除</span>
                            </a-tooltip>
                        </template>
                        <template v-else>
                            <a-popconfirm title="确定要删除此分组吗？" @confirm="handleDelete(record.id)" okText="删除" cancelText="取消">
                                <span class="app-link app-op app-remove"><a-icon type="delete" />删除</span>
                            </a-popconfirm>
                        </template>
                    </template>
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
                <space-update-component ref="updateRef" :detail="dialogDetail"></space-update-component>
            </a-spin>
        </a-modal>

    </div>
</template>

<script>
import { getSpaceListApi, updateSpaceApi, newSpaceApi, getSpaceDetailApi, deleteSpaceApi } from '@/api/project.js'
import SpaceUpdateComponent from './SpaceUpdateComponent.js'
export default {
    components: {
        SpaceUpdateComponent,
    },
    data() {
        return {
            tableLoading: false,
            tableColumns: [
                {dataIndex: "name", title: '项目空间', scopedSlots: { customRender: 'name' }},
                {dataIndex: "op", title: '操作', width: '30%', align: 'right', scopedSlots: { customRender: 'op' }},
            ],
            tableSource: [],
            search: {
                keyword: '',
            },
            pagination: {
                pageSize: 10,
                current: 1,
                total: 0,
            },

            dialogTitle: '',
            dialogDetail: {},
            dialogVisible: false,
            dialogConfirmLoading: false,
            dialogLoading: false,
        }
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
        handleOpenAddDialog() {
            this.dialogTitle = '新增项目空间'
            this.dialogDetail = {}
            this.dialogVisible = true
        },
        handleOpenUpdateDialog(id) {
            this.dialogTitle = '编辑项目空间'
            this.dialogVisible = true
            this.dialogDetail = {}
            this.getDataDetail(id)
        },
        handleDelete(id) {
            deleteSpaceApi({id}).then(res => {
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
                    updateSpaceApi(values).then(res => {
                        this.$message.success('更新项目空间成功', 1, () => {
                            this.dialogCancel()
                            this.dialogConfirmLoading = false
                            this.handleTableChange(this.pagination)
                        })
                    }).catch(err => {
                        this.dialogConfirmLoading = false
                    })
                } else {
                    newSpaceApi(values).then(res => {
                        this.$message.success('新增项目空间成功', 1, () => {
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
        getDataDetail(id) {
            this.dialogLoading = true
            getSpaceDetailApi({id}).then(res => {
                this.dialogLoading = false
                this.dialogDetail = res
            }).catch( err => {
                this.dialogLoading = false
            })
        },
        getDataList(params) {
            this.tableLoading = true
            let offset = (params.page - 1) * params.pageSize
            getSpaceListApi({keyword: this.search.keyword, offset: offset, limit: params.pageSize}).then(res => {
                this.tableLoading = false
                this.pagination.total = res.total
                this.tableSource = res.list
            }).catch(err => {
                this.tableLoading = false
            })
        }
    },
    mounted() {
        this.handleTableChange(this.pagination)
    },
}
</script>
