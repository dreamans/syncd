import { Form } from 'ant-design-vue'
const GroupUpdate = {
    render() {
        const { getFieldDecorator, getFieldValue } = this.form
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 15 },
        }
        let groupListOptions = []
        this.groupList.forEach((group) => {
            groupListOptions.push(
                <a-select-option value={group.id}>{ group.name }</a-select-option>
            )
        })
        getFieldDecorator('id', {
            initialValue: this.detail.id,
        })
        return (
            <a-form>
                <a-form-item
                {...{ props: formItemLayout }}
                label='所属集群'>
                    {getFieldDecorator('group_id', {
                        rules: [
                            { required: true, message: '请选择服务器所属分组' },
                        ],
                        initialValue: this.detail.group_id,
                    })(
                        <a-select
                        allowClear={true}
                        showSearch
                        placeholder="支持关键词搜索"
                        notFoundContent="无分组数据"
                        optionFilterProp="children">
                            {groupListOptions}
                        </a-select>
                    )}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='服务器名称'>
                    {getFieldDecorator('name', {
                        rules: [
                            { required: true, message: '服务器名称不能为空' },
                        ],
                        initialValue: this.detail.name,
                    })(
                        <a-input autocomplete="off" placeholder='请输入服务器名称' />
                    )}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='IP'>
                    {getFieldDecorator('ip', {
                        rules: [
                            { required: true, message: '服务器IP不能为空' },
                            { validator: function(rule, value, callback) {
                                if (!value) {
                                    callback('请输入ip地址')
                                    return
                                }
                                let arrIp = value.split('.')
                                if (arrIp.length != 4) {
                                    callback('请输入有效ip地址')
                                }
                                arrIp.forEach(i => {
                                    let num = parseInt(i)
                                    if (isNaN(num)) {
                                        callback('请输入有效ip地址')
                                    }
                                    if (num < 0 || num > 255) {
                                        callback('请输入有效ip地址')
                                    }
                                })
                                callback()
                            }},
                        ],
                        initialValue: this.detail.ip,
                    })(
                        <a-input autocomplete="off" placeholder='请输入服务器IP' />
                    )}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='ssh端口'>
                    {getFieldDecorator('ssh_port', {
                        rules: [
                            { required: true, type: 'integer', min: 1, max: 65535, message: '请输入正确ssh端口号，1-65535之间的数字' },
                        ],
                        initialValue: this.detail.ssh_port,
                        normalize: (arg) => {
                            let num = parseInt(arg)
                            if (isNaN(num)) {
                                return arg
                            }
                            return num
                        },
                    })(
                        <a-input autocomplete="off" placeholder='请输入ssh端口' />
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
        groupList: {
            default: () => {
                return []
            },
            type: Array,
        },
    },
}
export default Form.create()(GroupUpdate)
