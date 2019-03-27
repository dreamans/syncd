// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "html/template"
    "bytes"
    "time"
    "fmt"

    "github.com/dreamans/syncd"
    "github.com/dreamans/syncd/util/gostring"
)

const (
    MAIL_MODE_AUDIT_NOTICE = 1
    MAIL_MODE_AUDIT_RESULT = 2
    MAIL_MODE_DEPLOY = 3
)

const (
    MAIL_STATUS_SUCCESS = 1
    MAIL_STATUS_FAILED = 0
)

type MailMessage struct {
    Mail        string
    ApplyId     int
    Mode        int
    Status      int
    Title       string
}

func MailSend(msg *MailMessage) {
    mails, ok := mailSendToMails(msg.Mail)
    if !ok {
        return
    }
    syncd.App.Mail.AsyncSend(&syncd.SendMailMessage{
        To: mails,
        Subject: gostring.JoinStrings(mailSubjectPrefix(msg.Mode), msg.Title),
        Body: mailBodyTemplate(msg),
    })
}

func mailSubjectPrefix(mode int) string {
    prefix := "Syncd邮件通知:"
    switch mode {
    case MAIL_MODE_AUDIT_NOTICE:
        fallthrough
    case MAIL_MODE_AUDIT_RESULT:
        prefix = "Syncd审核通知:"
    case MAIL_MODE_DEPLOY:
        prefix = "Syncd部署通知:"
    }
    return prefix
}

func mailApplyLink(applyId int) string {
    return fmt.Sprintf("%s/deploy/deploy?id=%d", syncd.App.AppHost, applyId)
}

func mailSendToMails(mail string) ([]string, bool) {
    if mail != "" {
        mails := gostring.Str2StrSlice(mail, ",")
        if len(mails) > 0 {
            return mails, true
        }
    }
    return nil, false
}

func mailBodyTemplate(msg *MailMessage) string {
    link := mailApplyLink(msg.ApplyId)
    tpl := `
    <style>
        .syncd-main {
            font-family: "Chinese Quote", BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", "Helvetica Neue", Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
            background: #fff;
            font-size: 16px;
            padding: 0px;
            border: 0px;
            overflow: hidden;
        }
        .syncd-main .mt50 {
            margin-top: 50px;
        }
        .syncd-main a {
            cursor: pointer;
            text-decoration: none;
            color: #1890ff;
        }
        .syncd-main a:hover {
            color: #40a9ff;
        }
        .syncd-main .strong {
            font-weight: 500;
        }
        .syncd-main .syncd-card {
            background: #fff;
            line-height: 1.8;
            border: 8px dashed #EFEFEF;
            padding: 5px 20px 10px;
            width: 80%;
            margin: 10px auto 0;
            box-sizing: border-box;
        }
        .syncd-main .btn {
            line-height: 1;
            cursor: pointer;
            background: #fff;
            border: 1px solid #dcdfe6;
            color: #606266;
            text-align: center;
            box-sizing: border-box;
            outline: none;
            margin: 0;
            transition: .1s;
            font-weight: 500;
            font-size: 14px;
            padding: 9px 15px;
            border-radius: 4px;
            text-decoration: none;
        }
        .syncd-main .btn-primary {
            color: #fff;
            background-color: #409eff;
            border-color: #409eff;
        }
        .syncd-main .btn-primary:hover {
            background: #66b1ff;
            border-color: #66b1ff;
            color: #fff;
        }
        .syncd-main .tips {
            font-size: 14px;
            color: #909399;
        }
        .syncd-main .underline {
            text-decoration:underline;
        }
        .syncd-main .syncd-cpy {
            padding: 0 20px 10px;
            color: rgba(0, 0, 0, 0.65);
            font-size: 12px;
            width: 80%;
            margin: 0 auto;
            margin-top: 20px;
        }
        .syncd-main .syncd-success {
            color: #52c41a;
        }
        .syncd-main .syncd-failed {
            color: #f5222d;
        }
    </style>
    <div class="syncd-main">
        <div class="syncd-card">
            <p>您好:</p>
            {{ if eq .Mode 1 }}
            <p>上线申请单 <a target="_blank" href="{{ .Link }}">“{{ .Title }}(ID:{{ .Id }})”</a> 需要您审核，请尽快登录系统进行操作。</p>
            {{ else if eq .Mode 2 }}
            <p>
                您提交的上线申请单 <a target="_blank" href="{{ .Link }}">“{{ .Title }}(ID:{{ .Id }})”</a> 
                {{ if eq .Status 1 }}
                <span class="syncd-success">审核通过</span>，可登录系统进行后续操作。
                {{ else }}
                <span class="syncd-failed">审核不通过</span>，可登录系统查看原因。
                {{ end }}
            </p>
            {{ else if eq .Mode 3 }}
            <p>
                发布单 <a target="_blank" href="{{ .Link }}">“{{ .Title }}(ID:{{ .Id }})”</a> 
                {{ if eq .Status 1 }}
                <span class="syncd-success">部署成功</span>
                {{ else }}
                <span class="syncd-failed">部署失败</span>
                {{ end }}
                ，可登录系统查看详细日志。
            </p>
            {{ end }}
            <p class="mt50"><a href="{{ .Link }}" class="btn btn-primary" target="_blank">登录Syncd</a></p>
            <p class="tips">或者复制此链接到浏览器进行访问: <span class="underline">{{ .Link }}</span></p>
        </div>
        <div class="syncd-cpy">©️ {{ .Year }} <a target="_blank" href="https://github.com/dreamans/syncd">Syncd</a>. All Rights Reserved. MIT License.</div>
    </div>
    `
    t, err := template.New("mail").Parse(tpl)
    if err != nil {
        syncd.App.Logger.Error("sendmail body template parse failed, err[%s], mode[%d], apply_id[%d]", err.Error(), msg.Mode, msg.ApplyId)
        return ""
    }
    buf := new(bytes.Buffer)
    data := map[string]interface{}{
        "Title": msg.Title,
        "Link": link,
        "Mode": msg.Mode,
        "Status": msg.Status,
        "Id": msg.ApplyId,
        "Year": time.Now().Year(),
    }
    t.Execute(buf, data)
    return buf.String()
}