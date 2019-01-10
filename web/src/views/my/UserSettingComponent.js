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
                label='角色'>
                    { this.$store.getters['account/getGroupName'] }
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='头像'>
                    <a-avatar shape="square" size={64} src={this.$store.getters['account/getAvatar']} />
                    <div class="app-font-small"><a href="http://cn.gravatar.com/support/activating-your-account/" target="_blank"><a-icon type="question-circle" /> 如何修改头像</a></div>
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='用户名'>
                    <template slot="help">
                        <span class="app-font-small"><a-icon type="info-circle" /> 修改用户名请联系管理员</span>
                    </template>
                    <a-input value={this.$store.getters['account/getUserName']} disabled />
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='邮箱'>
                    <template slot="help">
                        <span class="app-font-small"><a-icon type="info-circle" /> 修改邮箱请联系管理员</span>
                    </template>
                    <a-input value={this.$store.getters['account/getEmail']} disabled />
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='真实姓名'>
                    {getFieldDecorator('true_name', {
                        initialValue: this.$store.getters['account/getTrueName'],
                    })(
                        <a-input autocomplete="off" placeholder='请输入真实姓名' />
                    )}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='手机号'>
                    {getFieldDecorator('mobile', {
                        rules: [
                            { validator: function(rule, value, callback) {
                                if (!value) {
                                    callback()
                                }
                                if (/^1[3456789][0-9]{9}$/.test(value)) {
                                    callback()
                                }
                                callback('手机号格式错误')
                            }},
                        ],
                        initialValue: this.$store.getters['account/getMobile'],
                        validateTrigger: 'blur',
                    })(
                        <a-input autocomplete="off" placeholder='11位数字，目前只支持中国(+86)手机号' />
                    )}
                </a-form-item>
            </a-form>
        )
    },
}
export default Form.create()(UserUpdate)
