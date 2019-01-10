import { Form } from 'ant-design-vue'
const GroupUpdate = {
    render() {
        const { getFieldDecorator, getFieldValue } = this.form
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 15 },
        }
        getFieldDecorator('id', {
            initialValue: this.detail.id,
        })
        return (
            <a-form>
                <a-form-item
                {...{ props: formItemLayout }}
                label='集群名称'>
                    {getFieldDecorator('name', {
                        rules: [
                            { required: true, message: '集群名称不能为空' },
                        ],
                        initialValue: this.detail.name,
                    })(
                        <a-input autocomplete="off" placeholder='请输入集群名称' />
                    )}
                </a-form-item>
            </a-form>
        )
    },
    props: {
        detail: {
            default: () => {
                return {}
            },
            type: Object,
        }
    },
}
export default Form.create()(GroupUpdate)
