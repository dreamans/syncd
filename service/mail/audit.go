// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mail

import (
    "bytes"
    "fmt"
    "time"

    "github.com/tinystack/goutil/gostring"
    deployService "github.com/tinystack/syncd/service/deploy"
    projectService "github.com/tinystack/syncd/service/project"
)

func AuditSend(emailList []string, apply *deployService.Apply, project *projectService.Project, userId int, userName, userEmail string) {
    var subjectBuff bytes.Buffer
    subjectBuff.WriteString("syncd-审核通知:")
    subjectBuff.WriteString(apply.Name)

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
        margin: 20px;
        overflow: hidden;
    }
    .syncd-main .syncd-card {
        font-size: 14px;
        margin: 20px;
        padding: 20px;
        background: #fff;
        line-height: 1.8;
    }
    .syncd-main .item {
        display: flex;
        margin-bottom: 20px;
    }
    .syncd-main .title {
        font-weight: 500;
        width: 100px;
    }
    .syncd-main .content {
        width: 70%;
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
</style>
<div class="syncd-main">
    <div class="syncd-card">
        <p>您好，</p>
        <p>目前有上线单 “我的测试上线单” 需要审核，请尽快登录系统进行审核</p>
    </div>
    <div class="syncd-card">
        <div class="item">
            <div class="title">上线单ID:</div>
            <div class="content">`)
    bodyBuff.WriteString(gostring.Int2Str(apply.ID))
    bodyBuff.WriteString(`</div>
        </div>
        <div class="item">
            <div class="title">名称:</div>
            <div class="content">`)
    bodyBuff.WriteString(apply.Name)
    bodyBuff.WriteString(`</div>
        </div>
        <div class="item">
            <div class="title">描述:</div>
            <div class="content">`)
    bodyBuff.WriteString(apply.Description)
    bodyBuff.WriteString(`</div>
        </div>
        <div class="item">
            <div class="title">所属项目:</div>
            <div class="content">`)
    bodyBuff.WriteString(project.Name)
    bodyBuff.WriteString(`</div>
        </div>
        <div class="item">
            <div class="title">提交者:</div>
            <div class="content">`)
            bodyBuff.WriteString(fmt.Sprintf("%s(ID:%d) - %s", userName, userId, userEmail))
    bodyBuff.WriteString(`</div>
        </div>
        <div class="item">
            <div class="title">提交时间:</div>
            <div class="content">`)
    bodyBuff.WriteString(time.Unix(int64(apply.Ctime), 0).Format("2006-01-02 15:04:05"))
    bodyBuff.WriteString(`</div>
        </div>
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
