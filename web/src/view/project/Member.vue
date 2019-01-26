<template>
    <div>
        <el-card shadow="never" v-if="!spaceId">
            <el-alert
            class="app-btn-group"
            title="提示信息"
            type="info"
            :closable="false"
            description="成员管理需要指定项目空间，请点击 '选择项目空间' 按钮，在弹出窗口中进行选择。"
            show-icon></el-alert>
            <div class="app-btn-group">
                <el-button type="primary" @click="switchSpaceHandler" icon="el-icon-refresh" size="small">选择项目空间</el-button>
            </div>
            <div class="app-divider"></div>
        </el-card>
        <el-card shadow="never" v-else>
            <div class="app-btn-group">
                <el-button type="primary" @click="switchSpaceHandler" icon="el-icon-refresh" size="small">切换项目空间</el-button>
            </div>
            <div class="app-btn-group">
                <h4>{{ this.spaceDetail.name }}</h4>
                <p class="app-description">{{ this.spaceDetail.description }}</p>
            </div>
            <div class="app-divider"></div>
            <div class="app-btn-title">添加新成员</div>
            <el-row class="app-btn-group">
                <el-col :span="10">
                    <el-select 
                    v-model="memberId"
                    style="width: 100%;"
                    size="small"
                    clearable
                    filterable
                    remote
                    loading-text="正在查找用户..."
                    no-match-text="未查找到用户"
                    :remote-method="searchMemberHandler"
                    :loading="memberLoading"
                    placeholder="通过用户名、邮箱搜索用户">
                        <el-option v-for="m in memberSearchList" :key="m.id" :value="m.id" :label="m.username + ' - ' + m.email +  ' - ' + m.role_name"></el-option>
                    </el-select>
                </el-col>
                <el-col :span="8">
                    <el-button style="margin-left: 10px;" size="small" type="primary" @click="addMemberHandler">{{ $t('add') }}</el-button>
                </el-col>
            </el-row>
            <el-table
                class="app-table"
                size="medium"
                v-loading="tableLoading"
                :data="tableData">
                <el-table-column prop="name" :label="$t('project_space')">
                    <div slot-scope="scope" class="app-table-explain">
                        <h4>{{ scope.row.name }}</h4>
                        <p class="app-description">{{ scope.row.description }}</p>
                    </div>
                </el-table-column>
                <el-table-column :label="$t('operate')" width="180" align="right">
                    <template slot-scope="scope">
                        <el-button
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

        <el-dialog 
        :width="$root.DialogSmallWidth"
        :visible.sync="dialogSpaceVisible"
        title="选择项目空间">
            <div style="padding-bottom: 50px;">
                <el-select filterable @change="selectSpaceHandler" v-model="spaceId" style="width: 100%;" size="small" placeholder="选择项目空间">
                    <el-option v-for="s in spaceList" :key="s.id" :value="s.id" :label="$root.Substr(s.name, 50)" ></el-option>
                </el-select>
            </div>
        </el-dialog>
    </div>
</template>

<script>
import { listSpaceApi, detailSpaceApi, searchMemberApi, addMemberApi, listMemberApi } from '@/api/project'
import Code from '@/lib/code'
export default {
    data() {
        return {
            dialogSpaceVisible: false,

            tableData: [],
            tableLoading: false,

            spaceId: undefined,
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
                detailSpaceApi({id: this.spaceId}).then(res => {
                    this.spaceDetail = res
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
                this.$message.error('请先搜索并选择用户后再添加');
                return
            }
            addMemberApi({member_id: this.memberId, space_id: this.spaceId}).then(res => {
                this.memberId = undefined
                this.memberSearchList = []
                this.$message.success('成员添加成功');
            }).catch(err => {
                if (Code.CODE_ERR_DATA_REPEAT == err.code) {
                    this.$message.error('成员已经存在，请勿重复添加');
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
                deleteSpaceApi({id: row.id}).then(res => {
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
            })
        },
    },
    mounted() {
        this.loadSpaceList()
    }
}
</script>
