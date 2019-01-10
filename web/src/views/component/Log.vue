<template>
    <div>
        <a-card
        title="操作日志"
        :bodyStyle="{padding: '0px'}"
        :bordered="false">
            <a-table
            :columns="tableColumns"
            :dataSource="logData"
            :bordered="false"
            :pagination="false"
            :locale="{emptyText:'无日志'}"
            :showHeader="false">
                <span slot="ctime" class="app-cursor" slot-scope="text, record">
                    <a-tooltip placement="top">
                        <template slot="title">
                            <span>{{ $root.FormatDateTime(text) }}</span>
                        </template>
                        {{ $root.FormatDateFromNow(text) }}
                    </a-tooltip>
                </span>
                <span slot="user" slot-scope="text, record">
                    <template v-if="record.user_id">
                        {{record.user_name}}(UID:{{record.user_id}})
                    </template>
                    <template v-else>
                        robot
                    </template>

                </span>
                <span slot="op_content" slot-scope="text, record">
                    <span>{{$root.T(record.op_name)}}</span>
                    <span v-if="record.op_content">: {{record.op_content}}</span>
                </span>
            </a-table>
        </a-card>
    </div>
</template>

<script>
export default {
    data() {
        return {
            tableColumns: [
                {dataIndex: "op_content", scopedSlots: { customRender: 'op_content' }},
                {dataIndex: "user", width: '25%', scopedSlots: { customRender: 'user' }},
                {dataIndex: "ctime", width: '20%', scopedSlots: { customRender: 'ctime' }},
            ],
            dataSource: [],
        }
    },
    props: {
        logData: {
            type: Array,
            default: () => {
                return []
            },
        },
        divider: {
            type: Boolean,
            default: false,
        },
    },
}
</script>
