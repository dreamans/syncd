import { Form } from 'ant-design-vue'
import { getGroupListApi, checkUserExistsApi } from '@/api/user.js'
const UserUpdate = {
    render() {
        const { getFieldDecorator, getFieldValue } = this.form
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 16 },
        }
        getFieldDecorator('id', {
            initialValue: this.detail.id,
        })
        let groupOpts = []
        this.groupList.forEach(g => {
            groupOpts.push(
                <a-select-option value={g.id}>{g.name}</a-select-option>
            )
        })
        let passwordOpts = {
            validateTrigger: 'blur',
        }
        let passwordExtra = ''
        if (this.detail.id) {
            passwordOpts.rules = [
                { type:"string", min: 6, max: 20, message: '密码必须为6-20个字符' },
            ]
            passwordExtra = '不修改请留空'
        } else {
            passwordOpts.rules = [
                { required: true, message: '用户密码不能为空' },
                { type:"string", min: 6, max: 20, message: '密码必须为6-20个字符' },
            ]
        }
        let vm = this
        let userExistsCb = function(type, errmsg) {
            return function(rule, value, callback) {
                if (!value) {
                    callback()
                    return
                }
                let title = ''
                switch (type) {
                    case 'email':
                        title = '邮箱'
                        break
                    case 'name':
                        title = '用户名'
                        break;
                }
                vm.$set(vm.helps, type, `正在验证 ${value} 是否被占用...`)
                checkUserExistsApi({id: vm.detail.id, keyword: value, type: type}).then(res => {
                    if (!res.exists) {
                        vm.$set(vm.helps, type, undefined)
                        callback()
                    } else {
                        vm.$set(vm.helps, type, `抱歉！${title}已被占用，请重新输入`)
                        callback(errmsg)
                    }
                }).catch(err => {
                    vm.$set(vm.helps, type, '网络错误, 校验失败')
                    callback('网络错误, 校验失败')
                })
            }
        }

        return (
            <a-form>
                <a-form-item
                {...{ props: formItemLayout }}
                label='角色'>
                    {getFieldDecorator('group_id', {
                        rules: [
                            { required: true, message: '角色不能为空' },
                        ],
                        initialValue: this.detail.group_id,
                    })(
                        <a-select style="width: 100%">
                            {groupOpts}
                        </a-select>
                    )}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                help={this.helps.name}
                label='用户名'>
                    {getFieldDecorator('name', {
                        rules: [
                            { required: true, message: '用户名不能为空' },
                            { validator: userExistsCb('name', '用户名已经存在')},
                        ],
                        initialValue: this.detail.name,
                        validateTrigger: 'blur',
                    })(
                        <a-input autocomplete="off" placeholder='请输入用户名' />
                    )}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                extra={passwordExtra}
                label='密码'>
                    {getFieldDecorator('password', passwordOpts)(
                        <a-input type="password" autocomplete="off" placeholder='请输入用户密码 6-20个字符' />
                    )}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                help={this.helps.email}
                label='邮箱'>
                    {getFieldDecorator('email', {
                        rules: [
                            { required: true, message: '邮箱地址不能为空' },
                            { type:"email", message: '邮箱格式错误' },
                            { validator: userExistsCb('email', '邮箱已经存在')},
                        ],
                        initialValue: this.detail.email,
                        validateTrigger: 'blur',
                    })(
                        <a-input autocomplete="off" placeholder='请输入邮箱地址' />
                    )}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='真实姓名'>
                    {getFieldDecorator('true_name', {
                        initialValue: this.detail.true_name,
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
                        initialValue: this.detail.mobile,
                        validateTrigger: 'blur',
                    })(
                        <a-input autocomplete="off" placeholder='11位数字，目前只支持中国(+86)手机号' />
                    )}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                help='禁止后用户将无法登录'
                label='允许登录'>
                    {getFieldDecorator('lock_status', {
                        initialValue: this.detail.lock_status ? true: false,
                        valuePropName: 'checked',
                    })(
                        <a-switch>
                            <a-icon type="unlock" slot="checkedChildren"/>
                            <a-icon type="lock" slot="unCheckedChildren"/>
                        </a-switch>
                    )}
                </a-form-item>
            </a-form>
        )
    },
    data() {
        return {
            groupList: [],
            helps: {},
        }
    },
    props: {
        detail: {
            default: () => {
                return {}
            },
            type: Object,
        }
    },
    methods: {
        getUserGroupList() {
            getGroupListApi({offset: 0, limit: 9999}).then(res => {
                this.groupList = res.list
            })
        },
    },
    mounted() {
        this.getUserGroupList()
    },
}
export default Form.create()(UserUpdate)
