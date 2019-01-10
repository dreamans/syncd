import { Form } from 'ant-design-vue'
const ViewApply = {
    render() {
        const { getFieldDecorator, getFieldValue } = this.form
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 16 },
        }
        getFieldDecorator('id', {
            initialValue: this.detail.id,
        })
        return (
            <a-form>
                <a-form-item
                {...{ props: formItemLayout }}
                label='项目空间'>
                    {this.detail.space_name}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='项目名称'>
                    {this.detail.project_name}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='上线单名称'>
                    {this.detail.name}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='上线模式'>
                    {this.detail.repo_mode == 1 ? (
                        <div>
                            <div>
                                <a-icon type="branches" /> 分支上线 - <strong>{ this.detail.repo_branch }</strong> 分支
                            </div>
                            <div>
                                版本: {this.detail.repo_commit}
                            </div>
                        </div>
                    ) : (
                        <span>
                            <a-icon type="tag" /> Tag上线 - Tag: { this.detail.repo_tag }
                        </span>
                    )}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='描述信息'>
                    {this.detail.description}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='申请人'>
                    {this.detail.user_name} ({this.detail.user_email}) - {this.$root.FormatDateTime(this.detail.ctime)}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='审核'>
                    {getFieldDecorator('audit', {
                        initialValue: 1,
                    })(
                        <a-radio-group>
                            <a-radio value={1}><span class="app-color-success">审核通过</span></a-radio>
                            <a-radio value={0}><span class="app-color-error">审核拒绝</span></a-radio>
                        </a-radio-group>
                    )}
                </a-form-item>
                {getFieldValue('audit') == 0 ? (
                    <a-form-item
                    {...{ props: formItemLayout }}
                    label='拒绝原因'>
                        {getFieldDecorator('reject_reason', {
                            initialValue: '',
                        })(
                            <a-textarea placeholder="请输入拒绝原因" rows={3} />
                        )}
                    </a-form-item>
                ) : ''}
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
    methods: {

    },
}
export default Form.create()(ViewApply)
