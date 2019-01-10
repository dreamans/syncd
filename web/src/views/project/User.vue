<template>
    <div>
        <a-card
        class="app-card"
        title="成员管理"
        :bordered="false">
            <template v-if="$root.CheckPriv($root.Priv.PROJECT_USER_NEW)">
                <div class="app-btn-group">
                    <div style="margin-bottom: 15px;">添加新成员到 <strong>{{ this.spaceDetail.name }}</strong> 项目空间</div>
                    <a-row :gutter="10">
                        <a-col :span="12">
                            <a-select
                            labelInValue
                            showSearch
                            allowClear
                            :defaultActiveFirstOption="false"
                            :showArrow="false"
                            :filterOption="false"
                            style="width:100%;"
                            :notFoundContent="searchFetching ? undefined : '未找到用户'"
                            placeholder="通过用户名、邮箱搜索用户"
                            v-model="selectedUser"
                            @search="handleSearchUser">
                                <a-spin v-if="searchFetching" slot="notFoundContent" size="small"/>
                                <a-select-option v-for="u in searchUserList"
                                :value="`${u.id}`">
                                    {{ u.name }} - {{ u.email }} - {{u.group_name}}
                                </a-select-option>
                            </a-select>
                        </a-col>
                        <a-col :span="8">
                            <a-button @click="handleAppendUser" type="primary">添加</a-button>
                        </a-col>
                    </a-row>
                </div>
                <a-divider></a-divider>
            </template>
            <a-table
            :columns="tableColumns"
            :dataSource="tableSource"
            :pagination="pagination"
            @change="handleTableChange"
            :loading="tableLoading">
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
                    <a-popconfirm v-if="$root.CheckPriv($root.Priv.PROJECT_USER_DEL)" title="确定要移除吗？" @confirm="handleRemoveUser(record.id)" okText="移除" cancelText="取消">
                        <span class="app-link app-op app-remove"><a-icon type="delete" />移除</span>
                    </a-popconfirm>
                </span>
            </a-table>
        </a-card>
    </div>
</template>

<script>
import { getSpaceDetailApi, addSpaceUserApi, getSpaceUserListApi, removeSpaceUserApi, searchSpaceUserApi } from '@/api/project.js'
export default {
    data () {
        return {
            spaceDetail: {},
            selectedUser: {},
            searchUserTimer: null,
            searchFetching: false,
            searchUserList: [],

            tableLoading: false,
            tableColumns: [
                {dataIndex: "name", title: '用户名'},
                {dataIndex: "email", title: '邮箱'},
                {dataIndex: "group_name", title: '角色', width: '15%'},
                {dataIndex: "lock_status", title: '状态', width: '10%', scopedSlots: { customRender: 'lock_status' }},
                {dataIndex: "op", title: '操作', width: '10%', align: 'right', scopedSlots: { customRender: 'op' }},
            ],
            tableSource: [],
            pagination: {
                pageSize: 10,
                current: 1,
                total: 0,
            },
        }
    },
    methods: {
        handleSearchUser(value) {
            clearTimeout(this.searchUserTimer)
            let vm = this
            vm.searchFetching = true
            vm.searchUserList = []
            vm.searchUserTimer = setTimeout(function() {
                searchSpaceUserApi({keyword: value}).then(res => {
                    vm.searchUserList = res.list
                    vm.searchFetching = false
                }).catch(err => {
                    vm.searchFetching = false
                })
            }, 500)
        },
        handleAppendUser() {
            if (!this.selectedUser) {
                this.$message.error('请先选择用户再添加');
            }
            addSpaceUserApi({space_id: this.spaceId, user_id: this.selectedUser.key}).then(res => {
                this.$message.success('用户添加成功');
                this.selectedUser = undefined
                this.handleTableChange(this.pagination)
            }).catch(err => {
                this.$message.error('用户添加失败, ' + err.message);
                this.selectedUser = undefined
            })
            this.searchUserList = []
        },
        handleTableChange(pagination) {
            this.pagination.current = pagination.current
            this.getDataList({
                page: pagination.current,
                pageSize: pagination.pageSize,
            })
        },
        handleRemoveUser(id) {
            removeSpaceUserApi({id}).then(res => {
                this.$message.success('移除成功', 1)
                this.$root.ResetPagination(this.pagination)
                this.handleTableChange(this.pagination)
            })
        },
        getDataList(params) {
            this.tableLoading = true
            let offset = (params.page - 1) * params.pageSize
            getSpaceUserListApi({offset: offset, limit: params.pageSize}).then(res => {
                this.tableLoading = false
                this.pagination.total = res.total
                this.tableSource = res.list
            }).catch(err => {
                this.tableLoading = false
            })
        },
        getSpaceDetail(spaceId) {
            getSpaceDetailApi({id: spaceId}).then(res => {
                this.spaceDetail = res
            })
        },
    },
    mounted() {
        let spaceId = parseInt(this.$route.query.space)
        if (!spaceId) {
            this.$root.GotoRouter('projectSpace')
        }
        this.spaceId = spaceId
        this.getSpaceDetail(this.spaceId)
        this.handleTableChange(this.pagination)
    },
}
</script>
