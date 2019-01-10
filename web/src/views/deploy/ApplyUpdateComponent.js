import { Form } from 'ant-design-vue'
import { getRepoTagListApi, getRepoCommitListApi, getApplyProjectDetailApi } from '@/api/deploy.js'
const UpdateApply = {
    render() {
        const { getFieldDecorator, getFieldValue } = this.form
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 16 },
        }
        const renderTagListOpts = () => {
            let list = []
            this.tagList.forEach(t => {
                list.push(
                    <a-select-option value={`${t}`}>{`${t}`}</a-select-option>
                )
            })
            return list
        }
        const renderCommitListOpts = () => {
            let list = []
            this.commitList.forEach(t => {
                list.push(
                    <a-select-option value={`${t}`}>{`${t}`}</a-select-option>
                )
            })
            return list
        }
        getFieldDecorator('project_id', {
            initialValue: this.detail.project_id,
        })
        getFieldDecorator('id', {
            initialValue: this.detail.id,
        })
        return (
            <a-form>
                <a-form-item
                {...{ props: formItemLayout }}
                label='项目名称'>
                    {this.projectDetail.name}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='上线单名称'>
                    {getFieldDecorator('name', {
                        rules: [
                            { required: true, message: '上线单名称不能为空' },
                        ],
                        validateTrigger: 'blur',
                        initialValue: this.detail.name,
                    })(
                        <a-input autocomplete="off" placeholder='请输入上线单名称' />
                    )}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='上线模式'>
                    {this.projectDetail.repo_mode == 1 ? (
                        <span>
                            <a-icon type="branches" /> 分支上线 - <strong>{ this.projectDetail.repo_branch }</strong> 分支
                        </span>
                    ) : (
                        <span>
                            <a-icon type="tag" /> Tag上线
                        </span>
                    )}
                </a-form-item>
                {this.projectDetail.repo_mode == 1 ? (
                    <a-form-item
                    {...{ props: formItemLayout }}
                    label='选择上线版本'>
                        {getFieldDecorator('commit', {
                            rules: [
                                { required: true, message: '选择上线版本' },
                            ],
                            validateTrigger: 'blur',
                        })(<div></div>)}
                        <div>
                            <a-select
                            showSearch
                            placeholder="关键词搜索"
                            notFoundContent={this.fetching ? undefined : '未找到数据'}
                            onChange={this.handleSelectCommit}
                            defaultActiveFirstOption={false}
                            style={{ width: '100%' }}>
                                { this.fetching ? (
                                    <a-spin slot="notFoundContent" size="small"/>
                                ): '' }
                                {renderCommitListOpts()}
                            </a-select>
                            <a-button onClick={this.handleFetchCommitList}><a-icon type="sync" spin={this.fetching} /> { this.fetchBtnTitle }</a-button>
                            {
                                getFieldValue('commit') ? (
                                    <div style="line-height: 1; margin-top: 10px;">选择版本 - { getFieldValue('commit') }</div>
                                ) : ''
                            }
                            { this.detail.repo_branch && this.detail.repo_commit ? (
                                <div style="line-height: 1; margin-top: 10px;">
                                    <p>当前版本 - {this.detail.repo_branch} - {this.detail.repo_commit}</p>
                                </div>
                            ): ''}
                        </div>
                    </a-form-item>
                ) : ''}
                {this.projectDetail.repo_mode == 2 ? (
                    <a-form-item
                    {...{ props: formItemLayout }}
                    label='选择上线Tag'>
                        {getFieldDecorator('tag', {
                            rules: [
                                { required: true, message: '选择上线Tag' },
                            ],
                            validateTrigger: 'blur',
                        })(<div></div>)}
                        <a-row gutter={10}>
                            <a-col span={16}>
                                <a-select
                                showSearch
                                placeholder="关键词搜索"
                                notFoundContent={this.fetching ? undefined : '未找到数据'}
                                onChange={this.handleSelectTags}
                                defaultActiveFirstOption={false}
                                style={{ width: '100%' }}>
                                    { this.fetching ? (
                                        <a-spin slot="notFoundContent" size="small"/>
                                    ): '' }
                                    {renderTagListOpts()}
                                </a-select>
                            </a-col>
                            <a-col span={8}>
                                <a-button onClick={this.handleFetchTagList}><a-icon type="sync" spin={this.fetching} /> { this.fetchBtnTitle }</a-button>
                            </a-col>
                        </a-row>
                        { this.detail.repo_tag != "" ?  (<div>{this.detail.repo_tag}</div>): '' }
                    </a-form-item>
                ) : ''}
                <a-form-item
                {...{ props: formItemLayout }}
                help={this.nameHelp}
                label='上线说明'>
                    {getFieldDecorator('description', {
                        rules: [
                            { required: true, message: '请填写上线说明' },
                        ],
                        validateTrigger: 'blur',
                        initialValue: this.detail.description,
                    })(
                        <a-textarea placeholder="请详细填写上线说明" rows={3} />
                    )}
                </a-form-item>
            </a-form>
        )
    },
    props: {
        detail: {
            type: Object,
            default: () => {
                return {}
            },
        },
    },
    data() {
        return {
            projectDetail: {},
            fetching: false,
            commitList: [],
            tagList: [],
            fetchBtnTitle: '拉取列表',
        }
    },
    methods: {
        handleSelectCommit(val) {
            this.form.setFieldsValue({commit: val})
        },
        handleSelectTags(val) {
            this.form.setFieldsValue({tag: val})
        },
        handleFetchCommitList() {
            this.fetching = true
            this.fetchBtnTitle = '列表拉取中...'
            getRepoCommitListApi({id: this.detail.project_id}).then(res => {
                this.fetching = false
                if (res.list) {
                    this.commitList = res.list
                }
                this.fetchBtnTitle = '列表拉取成功'
            }).catch(err => {
                this.fetching = false
                this.fetchBtnTitle = '列表拉取失败'
                this.$error({
                    title: '列表拉取失败',
                    content: err.message,
                    okText: '知道了',
                });
            })
        },
        handleFetchTagList() {
            this.fetching = true
            this.fetchBtnTitle = '列表拉取中...'
            getRepoTagListApi({id: this.detail.project_id}).then(res => {
                this.fetching = false
                if (res.list) {
                    this.tagList = res.list
                }
                this.fetchBtnTitle = '列表拉取成功'
            }).catch(err => {
                this.fetching = false
                this.fetchBtnTitle = '列表拉取失败'
                this.$error({
                    title: '列表拉取失败',
                    content: err.message,
                    okText: '知道了',
                });
            })
        },
        loadProjectDetail() {
            getApplyProjectDetailApi({id: this.detail.project_id}).then(res => {
                this.projectDetail = res
            })
        },
    },
    mounted() {
        this.loadProjectDetail()
    },
}
export default Form.create()(UpdateApply)
