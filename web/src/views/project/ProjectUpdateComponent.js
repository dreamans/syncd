import { Form } from 'ant-design-vue'
import { newProjectApi, updateProjectApi, getProjectApi, checkProjectNameExistsApi } from '@/api/project.js'
import { getGroupListApi } from '@/api/server.js'
const UpdateProject = {
    render() {
        const { getFieldDecorator, getFieldValue, setFieldsValue } = this.form
        const formItemLayout = {
            labelCol: { span: 4 },
            wrapperCol: { span: 18 },
        }
        const formItemSmallLayout = {
            labelCol: { span: 4 },
            wrapperCol: { span: 6 },
        }
        const renderServerGroupOpts = () => {
            let list = []
            this.serverGroupList.forEach(ser => {
                list.push(
                    <a-select-option value={ser.id}>{ser.name}</a-select-option>
                )
            })
            return list
        }
        const renderServerGroupList = () => {
            let srvList = getFieldValue('deploy_server')
            let renderList = []
            if (srvList) {
                srvList.forEach((srv, index) => {
                    let server = this.findServerItem(srv)
                    if (!server.id) {
                        return
                    }
                    renderList.push(
                        <div class="ant-list-item">
                            <div class="ant-list-item-content ant-list-item-content-single">{index + 1}. {server.name} (ID:{server.id})</div>
                            <ul class="ant-list-item-action">
                                <li>
                                    <a onClick={this.removeServerItem.bind(this, srv)} class="oper-delete"><icon-delete />移除</a>
                                </li>
                            </ul>
                        </div>
                    )
                })
            }

            let renderTpl = ''
            if (renderList.length) {
                renderTpl = (
                    <div class="app-server-group-list">
                        <div class="title">已选集群列表</div>
                        <div class="ant-list ant-list-split ant-list-bordered">
                            <div class="ant-spin-nested-loading">
                                <div class="ant-spin-container" style="max-height: 200px; overflow-y: auto;">
                                    {renderList}
                                </div>
                            </div>
                        </div>
                    </div>
                )
            }
            return (
                renderTpl
            )
        }

        let vm = this
        let nameExistsCb = function(rule, value, callback) {
            if (!value) {
                callback()
                return
            }
            vm.nameHelp = `正在验证名称是否被占用...`
            checkProjectNameExistsApi({id: vm.detail.id, space_id: vm.spaceId, keyword: value}).then(res => {
                if (!res.exists) {
                    vm.nameHelp = undefined
                    callback()
                } else {
                    vm.nameHelp = `抱歉！该名称已经存在，请重新输入`
                    callback(`抱歉！该名称已经存在，请重新输入`)
                }
            }).catch(err => {
                vm.nameHelp = '网络错误, 校验失败'
                callback('网络错误, 校验失败')
            })
        }

        return (
            <a-form onSubmit={this.handleSubmit}>
                <a-spin spinning={this.detailLoading}>
                    <a-form-item
                    {...{ props: formItemLayout }}
                    help={this.nameHelp}
                    label='项目名称'>
                        {getFieldDecorator('name', {
                            rules: [
                                { required: true, message: '项目名称不能为空' },
                                { validator: nameExistsCb},
                            ],
                            initialValue: this.detail.name,
                            validateTrigger: 'blur',
                        })(
                            <a-input autocomplete="off" placeholder='请输入项目名称' />
                        )}
                    </a-form-item>
                    <a-form-item
                    {...{ props: formItemLayout }}
                    label='项目描述'>
                        {getFieldDecorator('description', {
                            rules: [
                                { required: true, message: '项目名称不能为空' },
                            ],
                            initialValue: this.detail.description,
                            validateTrigger: 'blur',
                        })(
                            <a-textarea placeholder="请输入项目描述信息" rows={3} />
                        )}
                    </a-form-item>
                    <a-form-item
                    {...{ props: formItemLayout }}
                    label='开启审核'
                    help='开启后，上线单需要审核通过后才能发起上线'>
                        {getFieldDecorator('need_audit', {
                            initialValue: this.detail.need_audit ? true: false,
                            valuePropName: 'checked',
                        })(
                            <a-switch checkedChildren="开启" unCheckedChildren="关闭"/>
                        )}
                    </a-form-item>
                    <a-form-item
                    {...{ props: formItemLayout }}
                    label='项目状态'>
                        {this.detail.status == 1 ? (
                            <span class="app-color-success">
                                <a-icon type="check" /> 已启用
                            </span>
                        ): (
                            <span class="app-color-error">
                                <a-icon type="close" /> 未启用
                            </span>
                        )}
                    </a-form-item>
                    <a-divider></a-divider>
                    <a-form-item
                    {...{ props: formItemLayout }}
                    label='Git仓库地址'>
                        {getFieldDecorator('repo_url', {
                            rules: [
                                { required: true, message: '代码仓库地址不能为空' },
                            ],
                            initialValue: this.detail.repo_url,
                            validateTrigger: 'blur',
                        })(
                            <a-input autocomplete="off" placeholder='请输入代码仓库地址' />
                        )}
                    </a-form-item>
                    <a-form-item
                    {...{ props: formItemLayout }}
                    help="测试环境推荐分支上线，生产环境推荐tag上线"
                    label='上线模式'>
                        {getFieldDecorator('repo_mode', {
                            rules: [
                                { required: true, message: '请选择上线模式' },
                            ],
                            initialValue: this.detail.repo_mode ? this.detail.repo_mode : 1,
                        })(
                            <a-radio-group>
                                <a-radio value={1}>分支上线</a-radio>
                                <a-radio value={2}>tag上线</a-radio>
                            </a-radio-group>
                        )}
                    </a-form-item>
                    { getFieldValue('repo_mode') == 1 ? (
                        <a-form-item
                        {...{ props: formItemSmallLayout }}
                        label='指定上线分支'>
                            {getFieldDecorator('repo_branch', {
                                rules: [
                                    { required: true, message: '请输入指定上线分支' },
                                ],
                                initialValue: this.detail.repo_branch,
                                validateTrigger: 'blur',
                            })(
                                <a-input autocomplete="off" placeholder='上线分支' />
                            )}
                        </a-form-item>
                    ) : ''}

                    <a-form-item
                    {...{ props: formItemLayout }}
                    label='排除文件'>
                        {getFieldDecorator('exclude_files', {
                            initialValue: this.detail.exclude_files,
                        })(
                            <a-textarea placeholder="请输入要排除的文件" rows={3} />
                        )}
                        <div>排除不上线的文件，一行一个，支持通配符 `*`</div>
                    </a-form-item>

                    <a-divider></a-divider>
                    <a-form-item
                    {...{ props: formItemLayout }}
                    label='选择上线集群'>
                        <a-select
                        allowClear={true}
                        showSearch
                        placeholder="关键词搜索"
                        notFoundContent="无数据"
                        style={{ width: '200px' }}
                        optionFilterProp="children"
                        onSelect={this.handleServerSelect}>
                            {renderServerGroupOpts()}
                        </a-select>
                        {getFieldDecorator('deploy_server', {
                            rules: [
                                { required: true, message: '请选择上线集群' },
                            ],
                            initialValue: this.detail.deploy_server ? this.detail.deploy_server: [],
                        })(
                            <div>
                                {renderServerGroupList()}
                            </div>
                        )}
                    </a-form-item>
                    <a-form-item
                    {...{ props: formItemSmallLayout }}
                    label='用户'>
                        {getFieldDecorator('deploy_user', {
                            rules: [
                                { required: true, message: '请选择上线集群' },
                            ],
                            initialValue: this.detail.deploy_user,
                        })(
                            <a-input autocomplete="off" placeholder='目标机用户' />
                        )}
                    </a-form-item>
                    <a-form-item
                    {...{ props: formItemLayout }}
                    label='目录'>
                        {getFieldDecorator('deploy_path', {
                            rules: [
                                { required: true, message: '请设置代码部署目录' },
                            ],
                            initialValue: this.detail.deploy_path,
                            validateTrigger: 'blur',
                        })(
                            <a-input autocomplete="off" placeholder='代码/包部署的目录' />
                        )}
                    </a-form-item>
                    <a-form-item
                    {...{ props: formItemLayout }}
                    label='部署前运行命令'>
                        {getFieldDecorator('pre_deploy_cmd', {
                            initialValue: this.detail.pre_deploy_cmd
                        })(
                            <a-textarea placeholder="代码部署之前运行的命令, 每行一个命令" rows={3} />
                        )}
                    </a-form-item>
                    <a-form-item
                    {...{ props: formItemLayout }}
                    label='部署后运行命令'>
                        {getFieldDecorator('post_deploy_cmd', {
                            initialValue: this.detail.post_deploy_cmd
                        })(
                            <a-textarea placeholder="代码部署之后运行的命令, 每行一个命令" rows={3} />
                        )}
                    </a-form-item>
                    <a-form-item
                    {...{ props: formItemSmallLayout }}
                    label='部署超时时间(秒)'>
                        {getFieldDecorator('deploy_timeout', {
                            rules: [
                                { required: true, message: '请设置部署超时时间' },
                                { validator: function(rule, value, callback) {
                                    if (!value) {
                                        callback()
                                        return
                                    }
                                    let num = Number(value)
                                    if (isNaN(num) || parseInt(num) != num) {
                                        callback('请输入有效正整数')
                                        return
                                    }
                                    if (num < 1) {
                                        callback('请输入有效正整数')
                                        return
                                    }
                                    callback()
                                }},
                            ],
                            initialValue: this.detail.deploy_timeout ? this.detail.deploy_timeout: 120,
                            validateTrigger: 'blur',
                        })(
                            <a-input autocomplete="off" placeholder='部署超时时间' />
                        )}
                    </a-form-item>
                    <a-divider></a-divider>
                    <a-form-item
                    {...{ props: formItemLayout }}
                    label='审核通知'>
                        <div slot="help">用户提交待审核上线单时, 系统会通过邮件通知相关leader及时审核, 多个邮箱地址请用 `,` 相隔</div>
                        {getFieldDecorator('audit_notice_email', {
                            initialValue: this.detail.audit_notice_email,
                            validateTrigger: 'blur',
                        })(
                            <a-input autocomplete="off" placeholder='请输入接收审核通知的邮箱地址' />
                        )}
                    </a-form-item>
                    <a-form-item
                    {...{ props: formItemLayout }}
                    label='上线通知'>
                        <div slot="help">接收上线通知的邮箱地址, 多个邮箱地址请用 `,` 相隔</div>
                        {getFieldDecorator('deploy_notice_email', {
                            initialValue: this.detail.deploy_notice_email,
                            validateTrigger: 'blur',
                        })(
                            <a-input autocomplete="off" placeholder='请输入接收上线通知的邮箱地址' />
                        )}
                    </a-form-item>
                    <a-divider></a-divider>
                    <div style="text-align: right">
                        <a-button  type="primary" htmlType='submit'>提交</a-button>
                    </div>
                </a-spin>
            </a-form>
        )
    },
    props: {
        projectId: {
            type: Number,
            default: 0,
        },
        spaceId: {
            type: Number,
            default: 0,
        },
    },
    data () {
        return {
            id: 0,
            detail: {},
            serverGroupList: [],
            detailLoading: false,
            nameHelp: undefined,
        }
    },
    methods: {
        handleSubmit(e) {
            e.preventDefault()
            this.form.validateFields((err, values) => {
                if (err) {
                    this.$root.ResolveFormError(err, values)
                    return
                }
                let postData = {... values}
                if (!this.spaceId) {
                    this.$error({
                        title: '参数错误',
                        content: (
                            <div>
                                项目空间ID丢失，操作失败，请重试!
                            </div>
                        ),
                    })
                    return
                }
                postData.space_id = this.spaceId
                postData.id = this.projectId
                postData.need_audit = postData.need_audit ? 1: 0
                postData.deploy_server = this.filterInvalidServerGroup(postData.deploy_server)
                if (this.projectId) {
                    updateProjectApi(postData).then(res => {
                        this.$success({
                            title: '更新成功',
                            okText: "确定",
                            content: (
                                <div>恭喜，项目更新成功，点击确定返回项目列表</div>
                            ),
                            onOk: () => {
                                this.$emit('close')
                            }
                        });
                    })
                } else {
                    newProjectApi(postData).then(res => {
                        this.$success({
                            title: '新增成功',
                            okText: "确定",
                            content: (
                                <div>恭喜，项目新增成功，点击确定返回项目列表</div>
                            ),
                            onOk: () => {
                                this.$emit('close')
                            }
                        });
                    })
                }
            })
        },
        handleServerSelect(value) {
            let list = this.form.getFieldValue('deploy_server')
            if (!list) {
                list = []
            }
            if (list.indexOf(value) == -1) {
                list.push(value)
            }
            this.form.setFieldsValue({ deploy_server: list})
        },
        findServerItem(id) {
            let server = {}
            this.serverGroupList.forEach(srv => {
                if (srv.id == id) {
                    server = srv
                }
            })
            return server
        },
        removeServerItem(srvId) {
            let list = this.form.getFieldValue('deploy_server')
            if (!list) {
                list = []
            }
            let index = list.indexOf(srvId)
            if (index > -1) {
                list.splice(index, 1)
            }
            this.form.setFieldsValue({ deploy_server: list})
        },
        filterInvalidServerGroup(groupList) {
            let newGroupList = []
            if (groupList) {
                groupList.forEach(id => {
                    if (this.findServerItem(id).id) {
                        newGroupList.push(id)
                    }
                })
            }
            return newGroupList
        },
        getDataDetail(id) {
            this.detailLoading = true
            getProjectApi({id}).then(res => {
                this.detail = res
                this.detailLoading = false
            })
        },
        getDataGroupList() {
            getGroupListApi({offset: 0, limit: 9999}).then(res => {
                this.serverGroupList = res.list
            })
        },
    },
    mounted() {
        if (this.projectId) {
            this.getDataDetail(this.projectId)
        }
        this.getDataGroupList()
    },
}
export default Form.create()(UpdateProject)
