<template>
    <div>
        <el-card shadow="never" v-if="!spaceId">
            <el-alert
            class="app-btn-group"
            :title="$t('prompt_message')"
            type="info"
            :closable="false"
            :description="$t('project_member_select_space_tips')"
            show-icon></el-alert>
            <div class="app-btn-group">
                <el-button type="primary" @click="switchSpaceHandler" icon="el-icon-refresh" size="small">{{ $t('select_project_space') }}</el-button>
            </div>
            <div class="app-divider"></div>
        </el-card>
        <el-card shadow="never" v-else>
            <div class="app-btn-group">
                <el-button type="primary" @click="switchSpaceHandler" icon="el-icon-refresh" size="small">{{ $t('switch_project_space') }}</el-button>
            </div>
            <div class="app-btn-group" style="min-height: 30px">
                <div v-loading="spaceLoading">
                    <h4>{{ this.spaceDetail.name }}</h4>
                    <p class="app-description">{{ this.spaceDetail.description }}</p>
                </div>
            </div>
            <div class="app-divider"></div>
            <template v-if="$root.CheckPriv($root.Priv.PROJECT_USER_NEW)">
                <div class="app-btn-title">{{ $t('add_new_member') }}</div>
                <el-row class="app-btn-group">
                    <el-col :span="10">
                        <el-select 
                        v-model="memberId"
                        style="width: 100%;"
                        size="small"
                        clearable
                        filterable
                        remote
                        :loading-text="$t('looking_for_users')"
                        :no-match-text="$t('user_not_found')"
                        :remote-method="searchMemberHandler"
                        :loading="memberLoading"
                        :placeholder="$t('search_for_users_by_username_and_email')">
                            <el-option v-for="m in memberSearchList" :key="m.id" :value="m.id" :label="m.username + ' - ' + m.email +  ' - ' + m.role_name"></el-option>
                        </el-select>
                    </el-col>
                    <el-col :span="8">
                        <el-button style="margin-left: 10px;" size="small" type="primary" @click="addMemberHandler">{{ $t('add') }}</el-button>
                    </el-col>
                </el-row>
            </template>
            <el-table
                class="app-table"
                size="medium"
                v-loading="tableLoading"
                :data="tableData">
                <el-table-column prop="username" width="220" :label="$t('username')">
                    <template slot-scope="scope">
                        <span v-if="scope.row.username">{{ scope.row.username }}</span>
                        <span class="app-line-through" v-else>已删除</span>
                    </template>
                </el-table-column>
                <el-table-column prop="email" width="260" :label="$t('email')"></el-table-column>
                <el-table-column prop="role_name" :label="$t('role')"></el-table-column>
                <el-table-column prop="status" width="100" :label="$t('status')">
                    <template slot-scope="scope">
                        <span class="app-color-success" v-if="scope.row.status == '1'">
                            <i class="iconfont icon-unlock"></i> {{ $t('normal') }}
                        </span>
                        <span class="app-color-error" v-else>
                            <i class="iconfont icon-lock"></i> {{ $t('locking') }}
                        </span>
                    </template>
                </el-table-column>
                <el-table-column :label="$t('operate')" width="100" align="right">
                    <template slot-scope="scope">
                        <el-button
                        v-if="$root.CheckPriv($root.Priv.PROJECT_USER_DEL)"
                        type="text"
                        icon="el-icon-delete"
                        class="app-danger"
                        @click="deleteHandler(scope.row)">{{ $t('remove') }}</el-button>
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

        <el-dialog 
        :width="$root.DialogSmallWidth"
        :visible.sync="dialogSpaceVisible"
        :title="$t('select_project_space')">
            <div style="padding-bottom: 50px;">
                <el-select filterable @change="selectSpaceHandler" v-model="spaceId" style="width: 100%;" size="small" :placeholder="$t('select_project_space')">
                    <el-option v-for="s in spaceList" :key="s.id" :value="s.id" :label="$root.Substr(s.name, 50)" ></el-option>
                </el-select>
            </div>
        </el-dialog>
    </div>
</template>

<script>
import { listSpaceApi, detailSpaceApi, searchMemberApi, addMemberApi, listMemberApi, removeMemberApi } from '@/api/project'
import Code from '@/lib/code'
export default {
    data() {
        return {
            dialogSpaceVisible: false,

            tableData: [],
            tableLoading: false,

            spaceId: undefined,
            spaceLoading: false,
            spaceList: [],
            spaceDetail: {},

            memberId: undefined,
            memberLoading: false,
            memberSearchList: [],
            lastSearchTimer: null,
        }
    },
    watch: {
        spaceId() {
            if (this.spaceId) {
                this.spaceLoading = true
                detailSpaceApi({id: this.spaceId}).then(res => {
                    this.spaceDetail = res
                    this.spaceLoading = false
                }).catch(err => {
                    this.spaceLoading = false
                })
                this.$root.PageInit()
                this.loadTableData()
            }
        }
    },
    methods: {
        searchMemberHandler(query) {
            let vm = this
            this.memberLoading = true
            clearTimeout(this.lastSearchTimer)
            this.lastSearchTimer = setTimeout(() => {
                searchMemberApi({keyword: query}).then(res => {
                    if (res) {
                        this.memberSearchList = res
                    }
                    this.memberLoading = false
                }).catch(err => {
                    this.memberLoading = false
                })
            }, 1000);
        },
        addMemberHandler() {
            if (!this.memberId) {
                this.$message.error(this.$t('search_and_select_users_before_adding'));
                return
            }
            addMemberApi({member_id: this.memberId, space_id: this.spaceId}).then(res => {
                this.memberId = undefined
                this.memberSearchList = []
                this.$message.success(this.$t('member_added_successfully'));
                this.$root.PageInit()
                this.loadTableData()
            }).catch(err => {
                if (Code.CODE_ERR_DATA_REPEAT == err.code) {
                    this.$message.error(this.$t('member_already_exists_do_not_add_repeat'));
                }
            })
        },
        selectSpaceHandler() {
            this.dialogSpaceVisible = false
        },
        switchSpaceHandler() {
            this.dialogSpaceVisible = true
        },
        deleteHandler(row) {
            this.$root.ConfirmDelete(() => {
                removeMemberApi({id: row.id}).then(res => {
                    this.$root.MessageSuccess()
                    this.$root.PageReset()
                    this.loadTableData()
                })
            })
        },
        currentChangeHandler() {
            this.loadTableData()
        },
        loadTableData() {
            this.tableLoading = true
            listMemberApi({space_id: this.spaceId, offset: this.$root.PageOffset(), limit: this.$root.PageSize}).then(res => {
                this.tableData = res.list
                this.$root.Total = res.total
                this.tableLoading = false
            }).catch(err => {
                this.tableLoading = false
            })
        },
        loadSpaceList() {
            listSpaceApi({offset: 0, limit: 999}).then(res => {
                if (res.list) {
                    this.spaceList = res.list
                }
                this.initSpaceId()
            })
        },
        initSpaceId() {
            if (this.spaceList.length && !this.spaceId) {
                this.spaceId = this.spaceList[0].id
            }
        },
    },
    mounted() {
        this.loadSpaceList()
    }
}
</script>
