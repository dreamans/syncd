// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mail

import (
    "bytes"

    "github.com/tinystack/goutil/gostring"
    deployService "github.com/dreamans/syncd/service/deploy"
)

func DeploySend(emailList []string, apply *deployService.Apply, taskList []deployService.DeployTask) {
    var subjectBuff bytes.Buffer
    subjectBuff.WriteString("syncd-上线通知:")
    subjectBuff.WriteString(apply.Name)
    if apply.Status == deployService.APPLY_STATUS_DEPLOY_SUCCESS {
        subjectBuff.WriteString("[上线成功]")
    } else {
        subjectBuff.WriteString("[上线失败]")
    }
    subject := subjectBuff.String()

    var bodyBuff bytes.Buffer
    bodyBuff.WriteString(`
<style>
    .syncd-main {
        font-family: "Chinese Quote", BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", "Helvetica Neue", Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
        background: #f0f2f5;
        font-size: 14px;
        padding: 0px;
        border: 0px;
        overflow: hidden;
    }
    .syncd-main .syncd-card {
        font-size: 14px;
        margin: 20px;
        padding: 20px;
        background: #fff;
        line-height: 1.8;
    }
    .syncd-main .syncd-head .title {
        font-weight: 500;
        width: 100px;
        display: inline-block;
    }
    .syncd-main .syncd-body .head-title {
        font-weight: 500;
    }
    .syncd-shell {
        background-color: #282e34;
        color: #b8c0cc;
        padding: 10px;
        border-radius: 4px;
        margin-top: 20px;
        margin-bottom: 20px;
    }
    .syncd-shell .shell-sub-title {
        color: #fff;
    }
    .syncd-shell .syncd-shell-item {
        margin-bottom: 15px;
    }
    .syncd-shell .shell-sub-body pre {
        margin: 0;
    }
    .syncd-cpy {
        padding: 0 20px 10px;
        color: rgba(0, 0, 0, 0.65);
        font-size: 12px;
    }
    .syncd-link {
        text-decoration: none;
        color: #1890ff;
    }
    .syncd-success {
        color: #52c41a;
    }
    .syncd-failed {
        color: #f5222d;
    }
</style>
<div class="syncd-main">
    <div class="syncd-card syncd-head">
    `)
    bodyBuff.WriteString("<div><span class=\"title\">上线单ID:</span>")
    bodyBuff.WriteString(gostring.Int2Str(apply.ID))
    bodyBuff.WriteString("</div><div><span class=\"title\">上线单名称:</span>")
    bodyBuff.WriteString(apply.Name)
    bodyBuff.WriteString("</div><div><span class=\"title\">状态:</span>")
    if apply.Status == deployService.APPLY_STATUS_DEPLOY_SUCCESS {
        bodyBuff.WriteString("<span class=\"syncd-success\">成功</span>")
    } else {
        bodyBuff.WriteString("<span class=\"syncd-failed\">失败</span>")
    }
    bodyBuff.WriteString("</div></div>")
    bodyBuff.WriteString(`<div class="syncd-card syncd-body">
        <div class="head-title">上线信息</div>
        <div class="syncd-shell">`)

    if apply.ErrorLog != "" {
        bodyBuff.WriteString("<div class=\"syncd-shell-item\"><div class=\"shell-sub-body syncd-failed\">错误信息 >>><pre>")
        bodyBuff.WriteString(apply.ErrorLog)
        bodyBuff.WriteString("</pre></div></div>")
    }

    for _, t := range taskList {
        bodyBuff.WriteString("<div class=\"syncd-shell-item\">")
        bodyBuff.WriteString("<div class=\"shell-sub-title\">")
        bodyBuff.WriteString(">>> ")
        bodyBuff.WriteString(t.Name)
        bodyBuff.WriteString("</div><div class=\"shell-sub-body\">")
        if t.Status == deployService.DEPLOY_STATUS_END {
            bodyBuff.WriteString("<span class=\"syncd-success\">完成</span>")
        } else {
            bodyBuff.WriteString("<span class=\"syncd-failed\">失败</span>")
        }
        if t.Output != "" {
            bodyBuff.WriteString("<pre>")
            bodyBuff.WriteString(t.Output)
            bodyBuff.WriteString("</pre>")
        }
        bodyBuff.WriteString("</div></div>")
    }

    bodyBuff.WriteString(`</div>
    </div>
    <div class="syncd-cpy">
        ©️ 2019 <a class="syncd-link" href="https://github.com/tinystack/syncd" target="_blank">Syncd</a> 版权所有, MIT License.
    </div>
</div>`)

    body := bodyBuff.String()

    mail := &SendMail{
        To: emailList,
        Subject: subject,
        Body: body,
    }
    mail.AsyncSend()
}
