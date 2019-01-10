import { Form } from 'ant-design-vue'
import { getOperateLogApi, statusDeployApi } from '@/api/deploy.js'
import LogComponent from '../component/Log.vue'
const ViewApply = {
    render() {
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 16 },
        }
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
                label='状态'>
                    {this.applyStatusShow(this.detail.status)}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='申请时间'>
                    {this.$root.FormatDateTime(this.detail.ctime)}
                </a-form-item>
                <a-form-item
                {...{ props: formItemLayout }}
                label='申请人'>
                    {this.detail.user_name} ({this.detail.user_email})
                </a-form-item>
                <LogComponent log-data={this.operateLog}></LogComponent>
            </a-form>
        )
    },
    data() {
        return {
            operateLog: [],
        }
    },
    components: {
        LogComponent,
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
        applyStatusShow(status) {
            let tit = '未知'
            switch (status) {
                case 1:
                    tit = '待审核'
                    break
                case 2:
                    tit = '审核不通过'
                    break
                case 3:
                    tit = '审核通过待上线'
                    break
                case 4:
                    tit = '上线中'
                    break
                case 5:
                    tit = '上线成功'
                    break
                case 6:
                    tit = '上线失败'
                    break
                case 7:
                    tit = '废弃'
                    break
            }
            return tit
        },
        loadOperateLog() {
            getOperateLogApi({id: this.detail.id}).then(res => {
                this.operateLog = res.list
            })
        },
    },
    mounted() {
        this.loadOperateLog()
    },
}
export default Form.create()(ViewApply)
