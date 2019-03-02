<template>
    <div>
        <el-card shadow="never">
            <el-form class="app-form" ref="formRef" :model="form" size="medium" label-width="130px">
                <el-form-item
                :label="$t('select_project_space')"
                prop="space_id"
                :rules="[
                    { required: true, message: $t('project_space_cannot_empty'), trigger: 'blur'},
                ]">
                    <el-select
                    @change="changeSpaceHandler"
                    v-model="form.space_id" 
                    filterable 
                    :placeholder="$t('please_input_keyword')">
                        <el-option
                        v-for="space in spaceList"
                        :key="space.id"
                        :label="space.name"
                        :value="space.id">
                        </el-option>
                    </el-select>
                </el-form-item>

                <el-form-item 
                :label="$t('select_project')"
                prop="project_id"
                :rules="[
                    { required: true, message: $t('project_cannot_empty'), trigger: 'blur'},
                ]">
                    <el-select 
                    v-model="form.project_id" 
                    filterable 
                    :placeholder="$t('please_input_keyword')">
                        <el-option
                        v-for="proj in projectList"
                        :key="proj.id"
                        :label="proj.name"
                        :value="proj.id">
                        </el-option>
                    </el-select>
                </el-form-item>

                <el-form-item>
                    <el-button  icon="el-icon-edit-outline
" size="small" type="primary" @click="openDialogHandler">{{ $t('input_apply_order') }}</el-button>
                </el-form-item>

            </el-form>
        </el-card>
        <el-dialog
        :width="$root.DialogNormalWidth"
        :title="$t('input_deploy_apply')"
        :visible.sync="dialogVisible"
        @close="closeDialogHandler">
            <div class="app-dialog" v-loading="dialogLoading">
                <el-form class="app-form" ref="dialogRef" :model="dialogForm" size="medium" label-width="130px">
                    <el-form-item 
                    :label="$t('project_name')">
                        {{ projectDetail.name }}
                    </el-form-item>
               
                    <el-form-item 
                    :label="$t('apply_name')"
                    prop="name"
                    :rules="[
                        { required: true, message: $t('name_cannot_empty'), trigger: 'blur'},
                    ]">
                        <el-input :placeholder="$t('please_input_apply_name')" v-model="dialogForm.name" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item :label="$t('deploy_mode')">
                        <span v-if="projectDetail.deploy_mode == 1">
                            <i class="iconfont icon-branch"></i> - {{ $t('branch_deploy') }}<template v-if="projectDetail.repo_branch"> - <strong>{{ projectDetail.repo_branch }}</strong> {{ $t('branch') }}</template>
                        </span>
                        <span v-if="projectDetail.deploy_mode == 2">
                            <i class="iconfont icon-tag"></i> {{ $t('tag_deploy') }}
                        </span>
                    </el-form-item>

                    <el-form-item 
                    v-if="projectDetail.deploy_mode == 2"
                    :label="$t('tag_name')"
                    prop="branch_name"
                    :rules="[
                        { required: true, message: $t('tag_name_cannot_empty'), trigger: 'blur'},
                    ]">
                        <el-input class="app-input-mini" :placeholder="$t('please_input_tag_name')" v-model="dialogForm.branch_name" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item 
                    v-if="projectDetail.deploy_mode == 1 && projectDetail.repo_branch == ''"
                    :label="$t('branch_name')"
                    prop="branch_name"
                    :rules="[
                        { required: true, message: $t('branch_name_cannot_empty'), trigger: 'blur'},
                    ]">
                        <el-input class="app-input-mini" :placeholder="$t('please_input_branch_name')" v-model="dialogForm.branch_name" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item 
                    v-if="projectDetail.deploy_mode == 1"
                    :label="$t('commit_version')"
                    prop="commit_version">
                        <el-input class="app-input-normal" :placeholder="$t('please_input_commit_version')" v-model="dialogForm.commit_version" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item 
                    :label="$t('deploy_illustrate')"
                    prop="description"
                    :rules="[
                        { required: true, message: $t('deploy_illustrate_cannot_empty'), trigger: 'blur'},
                    ]">
                        <el-input :rows="4" type="textarea" :placeholder="$t('please_input_deploy_illustrate')" v-model="dialogForm.description" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item
                    :label="$t('rollback_apply')"
                    prop="rollback_id">
                        <el-select 
                        v-model="dialogForm.rollback_id" 
                        filterable 
                        style="width: 100%;"
                        :placeholder="$t('please_input_keyword')">
                            <el-option
                            v-for="r in rollbackList"
                            :key="r.id"
                            :label="'(ID:' + r.id + ') ' + r.name"
                            :value="r.id">
                            </el-option>
                        </el-select>
                    </el-form-item>

                 </el-form>
                 <div slot="footer" class="dialog-footer">
                    <el-button size="small" @click="closeDialogHandler">{{ $t('cancel') }}</el-button>
                    <el-button :loading="btnLoading" size="small" type="primary" @click="dialogSubmitHandler">{{ $t('enter') }}</el-button>
                </div>
            </div>
        </el-dialog>
    </div>
</template>

<script>
import { listSpaceApi, listProjectApi } from '@/api/project'
import { applyProjectDetailApi, applySubmitApi, applyRollbackListApi } from '@/api/deploy'
export default {
    data() {
        return {
            form: {
                space_id: undefined,
                project_id: undefined,
            },
            spaceList: [],
            projectList: [],
            rollbackList: [],

            dialogVisible: false,
            dialogLoading: false,
            dialogForm: {},
            btnLoading: false,
            projectDetail: {},
        }
    },
    methods: {
        dialogSubmitHandler() {
            this.$refs.dialogRef.validate((valid) => {
                if (!valid) {
                    return false;
                }
                this.btnLoading = true
                let postData = {
                    project_id: this.form.project_id,
                    space_id: this.form.space_id,
                    name: this.dialogForm.name,
                    branch_name: this.dialogForm.branch_name,
                    description: this.dialogForm.description,
                    commit_version: this.dialogForm.commit_version,
                    rollback_id: this.dialogForm.rollback_id,
                }
                applySubmitApi(postData).then(res => {
                    this.$alert(this.$t('deploy_apply_submit_success'), this.$t('submit_success'), {
                        type: 'success',
                        confirmButtonText: this.$t('enter'),
                        callback: action => {
                            this.closeDialogHandler()
                            this.$router.push({name: 'deployDeploy'})
                        }
                    })
                    this.btnLoading = false
                }).catch(err => {
                    this.btnLoading = false
                })
            })
        },
        openDialogHandler() {
            this.$refs.formRef.validate((valid) => {
                if (!valid) {
                    return false;
                }
                applyProjectDetailApi({id: this.form.project_id}).then(res => {
                    this.projectDetail = res
                    this.dialogVisible = true
                })
                applyRollbackListApi({id: this.form.project_id}).then(res => {
                    if (res) {
                        this.rollbackList = res
                    }
                })
            })
        },
        closeDialogHandler() {
            this.dialogVisible = false
        },
        changeSpaceHandler(spaceId) {
            this.form.project_id = undefined
            this.loadProjectList(spaceId)
        },
        loadSpaceList() {
            listSpaceApi({offset: 0, limit : 999}).then(res => {
                if (res.list) {
                    this.spaceList = res.list
                }
            })
        },
        loadProjectList(spaceId) {
            listProjectApi({space_id: spaceId, offset: 0, limit : 999}).then(res => {
                if (res.list) {
                    this.projectList = res.list
                }
            })
        },
    },
    mounted() {
        this.loadSpaceList()
    },
}
</script>
