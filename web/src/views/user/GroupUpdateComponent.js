import { Form } from 'ant-design-vue'
import { getPrivListApi } from '@/api/user.js'
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

        return (
            <a-form>
                <a-form-item
                {...{ props: formItemLayout }}
                label='角色名称'>
                    {getFieldDecorator('name', {
                        rules: [
                            { required: true, message: '角色名称不能为空' },
                        ],
                        initialValue: this.detail.name,
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
