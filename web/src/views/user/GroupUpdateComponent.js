import { Form } from 'ant-design-vue'
import { getPrivListApi, checkGroupExistsApi } from '@/api/user.js'
const GroupUpdate = {
    render() {
        const { getFieldDecorator, getFieldValue } = this.form
        const formItemLayout = {
            labelCol: { span: 4 },
            wrapperCol: { span: 19 },
        }
        getFieldDecorator('id', {
            initialValue: this.detail.id,
        })
        getFieldDecorator('priv', {
            initialValue: this.detail.priv ? this.detail.priv : [],
        })

        let plainPrivCheckRender = []
        this.plainPrivCheckList.forEach(privGroup => {
            let privGroupRender = []
            privGroup.items.forEach(priv => {
                privGroupRender.push(
                    <a-checkbox value={priv.value}>{priv.label}</a-checkbox>
                )
            })
            plainPrivCheckRender.push(
                <a-row class="item">
                    <a-col span={3}>{privGroup.label}</a-col>
                    <a-col span={21}>{privGroupRender}</a-col>
                </a-row>
            )
        })

        let vm = this
        let groupExistsCb = function(type, errmsg) {
            return function(rule, value, callback) {
                if (!value) {
                    callback()
                    return
                }
                let title = '角色名称'
                vm.$set(vm.helps, type, `正在验证 ${value} 是否被占用...`)
                checkGroupExistsApi({id: vm.detail.id, keyword: value}).then(res => {
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
                help={this.helps.name}
                label='角色名称'>
                    {getFieldDecorator('name', {
                        rules: [
                            { required: true, message: '角色名称不能为空' },
                            { validator: groupExistsCb('name', '角色名称已经存在')},
                        ],
                        initialValue: this.detail.name,
                        validateTrigger: 'blur',
                    })(
                        <a-input autocomplete="off" placeholder='请输入角色名称' />
                    )}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='权限设置'>
                    <a-checkbox
                    indeterminate={this.indeterminate}
                    onChange={this.handleCheckAllPrivChange}
                    checked={this.checkPrivAll}>全选</a-checkbox>
                    <div class="app-check-box">
                        <a-checkbox-group
                        value={this.checkedPrivList}
                        onChange={this.handleCheckPrivChange}>
                            {plainPrivCheckRender}
                        </a-checkbox-group>
                    </div>
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
    watch: {
        detail() {
            this.handleCheckPrivChange(this.detail.priv)
        },
        plainPrivCheckList() {
            this.handleCheckPrivChange(this.detail.priv)
        },
    },
    data() {
        return {
            checkPrivAll: false,
            indeterminate: false,
            checkedPrivList: [],
            plainPrivCheckList: [],
            helps: {},
        }
    },
    methods: {
        handleCheckAllPrivChange(e) {
            this.checkPrivAll = e.target.checked
            this.indeterminate = false

            let checkedPrivList = []
            if (this.checkPrivAll) {
                this.plainPrivCheckList.forEach(g => {
                    g.items.forEach(p => {
                        checkedPrivList.push(p.value)
                    })
                })
            }

            this.checkedPrivList = checkedPrivList
            this.form.setFieldsValue({ priv: checkedPrivList})
        },
        handleCheckPrivChange(checkList) {
            if (!checkList) {
                return false
            }
            let allChecked = true
            this.plainPrivCheckList.forEach(g => {
                g.items.forEach(p => {
                    if (checkList.indexOf(p.value) == -1) {
                        allChecked = false
                    }
                })
            })
            this.checkPrivAll = allChecked
            this.indeterminate = false
            if (checkList.length > 0 && !allChecked) {
                this.indeterminate = true
            }

            this.checkedPrivList = checkList
            this.form.setFieldsValue({ priv: checkList})
        },
        loadPrivListData() {
            getPrivListApi().then(res => {
                this.plainPrivCheckList = res.list
            })
        },
    },
    mounted() {
        this.loadPrivListData()
    },
}
export default Form.create()(GroupUpdate)
