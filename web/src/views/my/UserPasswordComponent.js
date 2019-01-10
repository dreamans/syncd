import { Form } from 'ant-design-vue'
const UserUpdate = {
    render() {
        const { getFieldDecorator, getFieldValue } = this.form
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 16 },
        }
        return (
            <a-form>
                <a-form-item
                {...{ props: formItemLayout }}
                label='当前密码'>
                    {getFieldDecorator('password', {
                        rules: [
                            { required: true, message: '当前密码不能为空' },
                            { type:"string", min: 6, max: 20, message: '密码必须为6-20个字符' },
                        ],
                        validateTrigger: 'blur',
                    })(
                        <a-input type="password" autocomplete="off" placeholder='请输入当前密码' />
                    )}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='新密码'>
                    {getFieldDecorator('newpassword', {
                        rules: [
                            { required: true, message: '新密码不能为空' },
                            { type:"string", min: 6, max: 20, message: '密码必须为6-20个字符' },
                        ],
                        validateTrigger: 'blur',
                    })(
                        <a-input type="password" autocomplete="off" placeholder='请输入新密码 6-20个字符' />
                    )}
                </a-form-item>
            </a-form>
        )
    },
}
export default Form.create()(UserUpdate)
