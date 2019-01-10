<template>
    <div>
        <a-card
        class="app-card"
        title="上线单管理"
        :bordered="false">

            <div class="app-btn-group">
                <a-row :gutter="10">
                    <a-col :span="3">
                        <a-select
                        v-model="search.time"
                        allowClear
                        placeholder="选择提交时间"
                        style="width: 100%"
                        optionFilterProp="children">
                            <a-select-option v-for="t in searchTime" :value="t.day">{{ t.title }}</a-select-option>
                        </a-select>
                    </a-col>
                    <a-col :span="3">
                        <a-select
                        v-model="search.status"
                        allowClear
                        placeholder="状态"
                        style="width: 100%"
                        optionFilterProp="children">
                            <a-select-option v-for="s in searchStatus" :value="s.id">{{ s.title }}</a-select-option>
                        </a-select>
                    </a-col>
                    <a-col :span="12">
                        <a-select
                        v-model="search.projectId"
                        allowClear
                        showSearch
                        placeholder="选择项目"
                        notFoundContent="无项目数据"
                        style="width: 100%"
                        optionFilterProp="children">
                            <a-select-option v-for="proj in searchProjectList" :value="proj.id">{{ proj.space_name }} <a-icon type="right" /> {{ proj.name }} <span v-if="proj.status == 0">(未启用)</span> </a-select-option>
                        </a-select>
                    </a-col>
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
                <span class="app-cursor" slot="user_name" slot-scope="text, record">
                    <template v-if="record.user_name">
                    <a-tooltip placement="top">
                        <template slot="title">
                            <span>{{ record.user_email }}</span>
                        </template>
                        {{ record.user_name }}
                    </a-tooltip>
                    </template>
                    <template v-else>
                        <span class="app-line-through">已删除</span>
                    </template>
                </span>
                <span slot="project_name" slot-scope="text, record">
                    <template v-if="!record.space_name || !record.project_name">
                        <span class="app-line-through">已删除</span>
                    </template>
                    <template v-else>
                        {{ record.space_name }} <a-icon class="app-color-gray" type="right" /> {{ record.project_name }}
                    </template>
                </span>
                <span :class="rowClassName(record.status)" slot="name" slot-scope="text,record">
                    {{ text }}
                </span>
                <span class="app-cursor" slot="ctime" slot-scope="text, record">
                    <a-tooltip placement="top">
                        <template slot="title">
                            <span>{{ $root.FormatDateTime(text) }}</span>
                        </template>
                        {{ $root.FormatDateFromNow(text) }}
                    </a-tooltip>
                </span>
                <span :class="rowClassName(record.status)" slot="status" slot-scope="text, record">
                    <span v-if="text == 1"><a-icon type="clock-circle" /> 待审核</span>
                    <span v-if="text == 2"><a-icon type="warning" /> 审核拒绝</span>
                    <span v-if="text == 3"><a-icon type="smile" /> 待上线</span>
                    <span v-if="text == 4"><a-icon type="coffee" /> 上线中</span>
                    <span v-if="text == 5"><a-icon type="check-circle" /> 已上线</span>
                    <span v-if="text == 6"><a-icon type="frown" /> 上线失败</span>
                    <span v-if="text == 7"><a-icon type="close-circle" /> 已废弃</span>
                </span>
                <span slot="op" slot-scope="text, record">
                    <a-dropdown :trigger="['click']">
                        <a-menu slot="overlay" @click="handleMenuClick($event, record)">
                            <a-menu-item key="view"><a-icon type="eye" />查看</a-menu-item>

                            <template v-if="
                                $root.CheckPriv($root.Priv.DEPLOY_DEPLOY_ALL)
                                || (
                                    $root.CheckPriv($root.Priv.DEPLOY_DEPLOY_MY)
                                    && record.user_id == $store.getters['account/getUserId']
                                )">
                                <a-menu-item key="deploy" v-if="record.status == 3 || record.status == 4"><a-icon type="coffee" />上线</a-menu-item>
                                <a-menu-item key="deploy" v-if="record.status == 6 && $root.CheckPrivs([$root.Priv.DEPLOY_DEPLOY_MY, $root.Priv.DEPLOY_DEPLOY_ALL])"><a-icon type="coffee" />再次上线</a-menu-item>
                            </template>

                            <template v-if="$root.CheckPriv($root.Priv.DEPLOY_EDIT_MY) && record.user_id == $store.getters['account/getUserId']">
                                <a-menu-item key="edit" v-if="record.status == 1 || record.status == 2"><a-icon type="edit" />编辑</a-menu-item>
                            </template>

                            <template v-if="
                                $root.CheckPriv($root.Priv.DEPLOY_AUDIT_ALL)
                                || (
                                    $root.CheckPriv($root.Priv.DEPLOY_AUDIT_MY)
                                    && record.user_id == $store.getters['account/getUserId']
                                )">
                                <a-menu-item key="audit" v-if="record.status == 1"><a-icon type="audit" />审核</a-menu-item>
                                <a-menu-item key="unaudit" v-if="record.status == 3"><a-icon type="close" />取消审核</a-menu-item>
                            </template>

                            <a-menu-item key="discard" v-if="
                                (record.status == 1 ||
                                record.status == 2 ||
                                record.status == 3 ||
                                record.status == 6) && (
                                    $root.CheckPriv($root.Priv.DEPLOY_DROP_ALL)
                                    || (
                                        $root.CheckPriv($root.Priv.DEPLOY_DROP_MY)
                                        && record.user_id == $store.getters['account/getUserId']
                                    )
                                )
                            "><a-icon type="delete" />废弃</a-menu-item>

                        </a-menu>
                        <a-button style="margin-left: 8px">
                            操作 <a-icon type="down" />
                        </a-button>
                    </a-dropdown>
                </span>
            </a-table>
        </a-card>

        <a-modal
        title="查看上线单"
        :visible="dialogViewVisible"
        :confirmLoading="dialogViewConfirmLoading"
        :destroyOnClose="true"
        width="50%"
        :footer="null"
        @cancel="dialogViewCancel">
            <apply-view-component :detail="applyDetail"></apply-view-component>
        </a-modal>

        <a-modal
        title="审核上线单"
        :visible="dialogAuditVisible"
        :confirmLoading="dialogAuditConfirmLoading"
        :destroyOnClose="true"
        width="50%"
        @ok="dialogAuditSubmit"
        okText="审核"
        cancelText="取消"
        @cancel="dialogAuditCancel">
            <apply-audit-component :detail="applyDetail" ref="auditRef"></apply-audit-component>
        </a-modal>

        <a-modal
        title="编辑上线单"
        :visible="dialogEditVisible"
        :confirmLoading="dialogEditConfirmLoading"
        :destroyOnClose="true"
        width="50%"
        @ok="dialogEditSubmit"
        okText="再次提交审核"
        cancelText="放弃"
        @cancel="dialogEditCancel">
            <apply-update-component :detail="applyDetail" ref="auditRef"></apply-update-component>
        </a-modal>

    </div>
</template>

<script>
import { getApplyListApi, getApplyDetailApi, auditApplyApi, unAuditApplyApi, discardApplyApi, getApplyProjectAll, updateApplyApi } from '@/api/deploy.js'
import ApplyViewComponent from './ApplyViewComponent.js'
import ApplyAuditComponent from './ApplyAuditComponent.js'
import ApplyUpdateComponent from './ApplyUpdateComponent.js'
export default {
    data() {
        return {
            tableColumns: [
                {dataIndex: "id", title: 'ID', width: '6%'},
                {dataIndex: "name", title: '上线单名称', scopedSlots: { customRender: 'name' }},
                {dataIndex: "project_name", title: '空间名称/项目名称', width: '30%', scopedSlots: { customRender: 'project_name' }},
                {dataIndex: "ctime", title: '提交时间', width: '10%', scopedSlots: { customRender: 'ctime' }},
                {dataIndex: "user_name", title: '提交者', width: '10%', scopedSlots: { customRender: 'user_name' }},
                {dataIndex: "status", title: '状态', width: '12%', scopedSlots: { customRender: 'status' }},
                {dataIndex: "op", title: '操作', width: '10%', align: 'right', scopedSlots: { customRender: 'op' }},
            ],
            tableSource: [],
            pagination: {
                pageSize: 10,
                current: 1,
                total: 0,
            },
            tableLoading: false,
            search: {
                keyword: '',
                projectId: undefined,
                time: undefined,
                status: undefined,
            },
            searchProjectList: [],
            searchTime: [
                {
                    day: 1, title: '今天',
                },
                {
                    day: 7, title: '7天内',
                },
                {
                    day: 30, title: '一个月内',
                },
                {
                    day: 90, title: '3个月内',
                },
                {
                    day: 365, title: '一年内',
                },
                {
                    day: 0, title: '时间不限',
                },
            ],
            searchStatus: [
                {
                    id: 0, title: '不限',
                },
                {
                    id: 1, title: '待审核',
                },
                {
                    id: 2, title: '审核拒绝',
                },
                {
                    id: 3, title: '待上线',
                },
                {
                    id: 4, title: '上线中',
                },
                {
                    id: 5, title: '已上线',
                },
                {
                    id: 6, title: '上线失败',
                },
                {
                    id: 7, title: '已废弃',
                },
            ],

            dialogViewVisible: false,
            dialogViewConfirmLoading: false,
            applyDetail: {},

            dialogAuditVisible:false,
            dialogAuditConfirmLoading: false,

            dialogEditVisible: false,
            dialogEditConfirmLoading: false,
        }
    },
    components: {
        ApplyViewComponent, ApplyAuditComponent, ApplyUpdateComponent,
    },
    methods: {
        handleSearch(value) {
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
        handleMenuClick({key}, record) {
            switch (key) {
                case 'view':
                    this.handleShowView(record.id)
                    break
                case 'edit':
                    this.handleShowEdit(record.id)
                    break
                case 'audit':
                    this.handleShowAudit(record.id)
                    break
                case 'unaudit':
                    this.$confirm({
                        title: '取消审核',
                        content: '确定要取消已审核的上线单吗？',
                        onOk: () => {
                            this.handleUnAudit(record.id)
                        },
                        okText: '确定',
                        cancelText: '取消',
                    })
                    break
                case 'discard':
                    this.$confirm({
                        title: '废弃上线单',
                        content: '废弃后将无法恢复, 确定要废弃此上线单吗？',
                        onOk: () => {
                            this.handleDiscard(record.id)
                        },
                        okText: '确定',
                        cancelText: '取消',
                    })
                    break
                case 'deploy':
                    this.$root.GotoRouter('deployRelease', { id: record.id })
                    break
            }
        },

        handleShowView(id) {
            this.dialogViewConfirmLoading = true
            getApplyDetailApi({id}).then(res => {
                this.dialogViewVisible = true
                this.dialogViewConfirmLoading = false
                this.applyDetail = res
            }).catch(err => {
                this.dialogViewConfirmLoading = false
            })
        },
        dialogViewCancel() {
            this.dialogViewVisible = false
        },

        handleShowAudit(id) {
            this.dialogAuditConfirmLoading = true
            getApplyDetailApi({id}).then(res => {
                this.dialogAuditVisible = true
                this.dialogAuditConfirmLoading = false
                this.applyDetail = res
            }).catch(err => {
                this.dialogAuditConfirmLoading = false
            })
        },
        dialogAuditCancel() {
            this.dialogAuditVisible = false
        },
        dialogAuditSubmit() {
            this.$refs.auditRef.validateFields((err, values) => {
                if (err) {
                    return
                }
                this.dialogAuditConfirmLoading = true
                auditApplyApi(values).then(res => {
                    this.dialogAuditConfirmLoading = false
                    this.handleTableChange(this.pagination)
                    this.dialogAuditCancel()
                }).catch(err => {
                    this.dialogAuditConfirmLoading = false
                })
            })
        },

        handleShowEdit(id) {
            this.dialogEditConfirmLoading = true
            getApplyDetailApi({id}).then(res => {
                this.dialogEditVisible = true
                this.dialogEditConfirmLoading = false
                this.applyDetail = res
            }).catch(err => {
                this.dialogEditConfirmLoading = false
            })
        },
        dialogEditCancel() {
            this.dialogEditVisible = false
        },
        dialogEditSubmit() {
            this.$refs.auditRef.validateFields((err, values) => {
                if (err) {
                    return
                }
                this.dialogEditConfirmLoading = true
                updateApplyApi(values).then(res => {
                    this.dialogAuditConfirmLoading = false
                    this.handleTableChange(this.pagination)
                    this.dialogEditCancel()
                }).catch(err => {
                    this.dialogAuditConfirmLoading = false
                })
            })
        },

        handleUnAudit(id) {
            const hideLoading = this.$message.loading('取消审核上线单...', 0);
            unAuditApplyApi({id}).then(res => {
                hideLoading()
                setTimeout(() => {
                    this.$message.success("上线单取消成功")
                }, 500)
                this.handleTableChange(this.pagination)
            }).catch(err => {
                hideLoading()
                this.$message.error("上线单取消失败")
            })
        },
        handleDiscard(id) {
            const hideLoading = this.$message.loading('废弃上线单...', 0);
            discardApplyApi({id}).then(res => {
                hideLoading()
                this.handleTableChange(this.pagination)
            }).catch(err => {
                hideLoading()
                this.$message.error("上线单废弃失败")
            })
        },
        rowClassName(status) {
            let className = ''
            switch (status) {
                case 2:
                    className = 'app-color-warning'
                    break
                case 4:
                    className = 'app-color-info'
                    break
                case 7:
                    className = 'app-color-gray'
                    break
                case 6:
                    className = 'app-color-error'
                    break
            }
            return className
        },
        getDataList(params) {
            this.tableLoading = true
            let offset = (params.page - 1) * params.pageSize
            getApplyListApi({
                keyword: this.search.keyword,
                project_id: this.search.projectId,
                time: this.search.time,
                status: this.search.status,
                offset: offset,
                limit: params.pageSize,
            }).then(res => {
                this.tableLoading = false
                this.pagination.total = res.total
                this.tableSource = res.list
            }).catch(err => {
                this.tableLoading = false
            })
        },
        getProjectList() {
            getApplyProjectAll().then(res => {
                this.searchProjectList = res.list
            })
        },
    },
    mounted() {
        this.handleTableChange(this.pagination)
        this.getProjectList()
    },
}
</script>
