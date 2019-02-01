<template>
    <div>
        <el-card shadow="never">
            <el-row class="app-btn-group" :gutter="10">
                <el-col :span="3">
                    <el-select
                    :placeholder="$t('submit_time')"
                    size="medium" 
                    clearable style="width: 100%" 
                    v-model="searchTime">
                        <el-option
                            v-for="s in timeList"
                            :key="s.time"
                            :label="s.label"
                            :value="s.time">
                        </el-option>
                    </el-select>
                </el-col>
                <el-col :span="3">
                    <el-select
                    :placeholder="$t('audit_status')"
                    size="medium" 
                    clearable 
                    style="width: 100%" 
                    v-model="searchStatus">
                        <el-option
                            v-for="s in auditStatusList"
                            :key="s.status"
                            :label="s.label"
                            :value="s.status">
                        </el-option>
                    </el-select>
                </el-col>
                <el-col :span="3">
                    <el-select
                    :placeholder="$t('deploy_status')"
                    size="medium" 
                    clearable 
                    style="width: 100%" 
                    v-model="searchStatus">
                        <el-option
                            v-for="s in statusList"
                            :key="s.status"
                            :label="s.label"
                            :value="s.status">
                        </el-option>
                    </el-select>
                </el-col>
                <el-col :span="9">
                    <el-select
                    :placeholder="$t('select_project')"
                    size="medium" 
                    clearable 
                    style="width: 100%" 
                    v-model="searchProjectId">
                        <el-option
                        v-for="p in projectList"
                        :key="p.project_id"
                        :label="p.project_name"
                        :value="p.project_id">
                            {{p.space_name}} <i class="el-icon-arrow-right"></i> {{p.project_name}}
                        </el-option>
                    </el-select>
                </el-col>
                <el-col :span="6">
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
                <el-table-column prop="name" :label="$t('name')"></el-table-column>
                <el-table-column label="空间/项目名称">
                    <template slot-scope="scope">
                        
                    </template>
                </el-table-column>
                <el-table-column prop="ssh_port" width="100" label="提交时间"></el-table-column>
                <el-table-column prop="ssh_port" width="100" label="提交者"></el-table-column>
                <el-table-column prop="ssh_port" width="100" label="状态"></el-table-column>
                <el-table-column :label="$t('operate')" width="180" align="right">
                    <template slot-scope="scope">
                        
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
    </div>
</template>

<script>
import { applyProjectAllApi } from '@/api/deploy'
export default {
    data() {
        return {
            searchInput: '',
            searchTime: undefined,
            searchStatus: undefined,
            searchProjectId: undefined,

            tableLoading: false,
            tableData: [],

            timeList: [
                {time: 1, label: '今天'},
                {time: 7, label: '7天内'},
                {time: 30, label: '一个月内'},
                {time: 90, label: '3个月内'},
                {time: 365, label: '一年内'},
                {time: 0, label: '时间不限'},
            ],
            statusList: [
                {status: 1, label: '未上线'},
                {status: 2, label: '上线中'},
                {status: 3, label: '上线成功'},
                {status: 4, label: '上线失败'},
                {status: 5, label: '已废弃'},
            ],
            auditStatusList: [
                {status: 1, label: '待审核'},
                {status: 2, label: '审核通过'},
                {status: 3, label: '审核拒绝'},
            ],
            projectList: [],
        }
    },
    methods: {
        searchHandler() {

        },
        currentChangeHandler() {
            this.loadTableData()
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
        loadProjectAll() {
            applyProjectAllApi().then(res => {
                if (res && res.length > 0 ) {
                    this.projectList = res
                }
            })
        },
    },
    mounted() {
        this.loadProjectAll()
    },
}
</script>