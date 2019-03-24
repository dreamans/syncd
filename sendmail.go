// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package syncd

import (
    "strings"

    "gopkg.in/gomail.v2"
)

type SendMail struct {
    Enable  int
    Smtp    string
    Port    int
    User    string
    Pass    string
    dialer  *gomail.Dialer
}

func NewSendMail(mail *SendMail) *SendMail {
    mail.dialer = gomail.NewPlainDialer(mail.Smtp, mail.Port, mail.User, mail.Pass)
    return mail
}

func (mail *SendMail) send(msg *SendMailMessage) error {
    if mail.Enable == 0 {
        return nil
    }
    msg.mail = mail
    m := msg.NewMessage()
    if err := mail.dialer.DialAndSend(m); err != nil {
        return err
    }
    return nil
}

func (mail *SendMail) Send(msg *SendMailMessage) {
    if err := mail.send(msg); err != nil {
        App.Logger.Error(
            "send mail failed, to[%s], cc[%s], subject[%s]", 
            strings.Join(msg.To, ","),
            strings.Join(msg.Cc, ","),
            msg.Subject,
        )
    }
}

func (mail *SendMail) AsyncSend(msg *SendMailMessage) {
    go mail.Send(msg)
}

type SendMailMessage struct {
    From    string
    To      []string
    Cc      []string
    Subject string
    Body    string
    Attach  string
    mail    *SendMail
}

func (m *SendMailMessage) NewMessage() *gomail.Message {
    mailMsg := gomail.NewMessage()
    if m.From == "" {
        mailMsg.SetHeader("From", m.mail.User)
    } else {
        mailMsg.SetHeader("From", m.From)
    }
    mailMsg.SetHeader("To", m.To...)
    if len(m.Cc) > 0 {
        mailMsg.SetHeader("Cc", m.Cc...)
    }
    mailMsg.SetHeader("Subject", m.Subject)
    mailMsg.SetBody("text/html", m.Body)
    if m.Attach != "" {
        mailMsg.Attach(m.Attach)
    }
    return mailMsg
}