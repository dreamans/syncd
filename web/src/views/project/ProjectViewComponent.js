import { Form } from 'ant-design-vue'
import { getProjectApi } from '@/api/project.js'
const ViewProject = {
    render() {
        const formItemLayout = {
            labelCol: { span: 4 },
            wrapperCol: { span: 18 },
        }
        return (
            <a-spin spinning={this.loading}>
                <a-form>
                    <a-form-item
                    label='项目ID'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
                        {this.detail.id}
                    </a-form-item>
                    <a-form-item
                    label='项目名称'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
                        {this.detail.name}
                    </a-form-item>
                    <a-form-item
                    label='项目描述'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
                        {this.detail.description}
                    </a-form-item>
                    <a-form-item
                    label='开启审核'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
                        {this.detail.need_audit ? '需要审核': '不需要审核'}
                    </a-form-item>
                    <a-form-item
                    label='项目启用'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
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
                    label='代码仓库地址'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
                        {this.detail.repo_url}
                    </a-form-item>
                    <a-form-item
                    label='上线模式'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
                        {this.detail.repo_mode ? (
                            <span>
                                <a-icon type="branches" /> 分支上线 - <strong>{ this.detail.repo_branch }</strong> 分支
                            </span>
                        ) : (
                            <span>
                                <a-icon type="tag" /> Tag上线
                            </span>
                        )}
                    </a-form-item>
                    <a-divider></a-divider>
                    <a-form-item
                    label='上线集群'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
                        { function(server_group){
                            let r = []
                            if (server_group) {
                                server_group.forEach(g => {
                                    r.push(
                                        <a-tooltip placement="top" >
                                            <template slot="title">
                                                <span>集群ID: {g.id}, 集群名称: {g.name}</span>
                                            </template>
                                            <a-tag><a-icon type="cluster" /> {g.name}</a-tag>
                                        </a-tooltip>
                                    )
                                })
                            }
                            return r
                        }(this.detail.deploy_servers) }
                    </a-form-item>
                    <a-form-item
                    label='用户'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
                        { this.detail.deploy_user }
                    </a-form-item>
                    <a-form-item
                    label='目录'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
                        { this.detail.deploy_path }
                    </a-form-item>
                    <a-form-item
                    label='部署前运行命令'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
                        <pre class="app-shell">{ this.detail.pre_deploy_cmd }</pre>
                    </a-form-item>
                    <a-form-item
                    label='部署后运行命令'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
                        <pre class="app-shell">{ this.detail.post_deploy_cmd }</pre>
                    </a-form-item>
                    <a-form-item
                    label='部署超时时间(秒)'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
                        { this.detail.deploy_timeout } 秒
                    </a-form-item>
                    <a-divider></a-divider>
                    <a-form-item
                    label='审核通知'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
                        { this.detail.audit_notice_email }
                    </a-form-item>
                    <a-form-item
                    label='上线通知'
                    labelCol={formItemLayout.labelCol}
                    wrapperCol={formItemLayout.wrapperCol}>
                        { this.detail.deploy_notice_email }
                    </a-form-item>
                </a-form>
            </a-spin>
        )
    },
    props: {
        projectId: Number,
    },
    data() {
        return {
            loading: false,
            detail: {},
        }
    },
    methods: {
        getDetail(id) {
            this.loading = true
            getProjectApi({id}).then(res => {
                this.loading = false
                this.detail = res
            })
        }
    },
    mounted() {
        this.getDetail(this.projectId)
    },
}
export default Form.create()(ViewProject)
