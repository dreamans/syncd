import { Form } from 'ant-design-vue'
import { checkSpaceNameExistsApi } from '@/api/project.js'
const UpdateSpace = {
    render() {
        const { getFieldDecorator, getFieldValue } = this.form
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 16 },
        }
        getFieldDecorator('id', {
            initialValue: this.detail.id,
        })
        let vm = this
        let nameExistsCb = function(rule, value, callback) {
            if (!value) {
                callback()
                return
            }
            vm.nameHelp = `正在验证名称是否被占用...`
            checkSpaceNameExistsApi({id: vm.detail.id, keyword: value}).then(res => {
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
            <a-form>
                <a-form-item
                {...{ props: formItemLayout }}
                help={this.nameHelp}
                label='项目空间名称'>
                    {getFieldDecorator('name', {
                        rules: [
                            { required: true, message: '项目空间名称不能为空' },
                            { validator: nameExistsCb},
                        ],
                        initialValue: this.detail.name,
                        validateTrigger: 'blur',
                    })(
                        <a-input autocomplete="off" placeholder='请输入项目空间名称' />
                    )}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='描述信息'>
                    {getFieldDecorator('description', {
                        initialValue: this.detail.description,
                    })(
                        <a-textarea placeholder="请输入描述信息, 500字以内" rows={3} />
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
        },
    },
    data() {
        return {
            nameHelp: undefined,
        }
    },
}
export default Form.create()(UpdateSpace)
