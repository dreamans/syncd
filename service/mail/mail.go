// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mail

import (
    "github.com/tinystack/goutil/gostring"
    "github.com/dreamans/syncd"
)

type SendMail struct {
    To          []string
    Cc          []string
    Subject     string
    Body        string
}

func (m *SendMail) Send() error {
    return m.send()
}

func (m *SendMail) AsyncSend() {
    go func() {
        m.send()
    }()
}

func (m *SendMail) send() error {
    err := syncd.Mail.Send(&syncd.SendMailMessage{
        To: m.To,
        Cc: m.Cc,
        Subject: m.Subject,
        Body: m.Body,
    })
    if err != nil {
        syncd.Logger.Error("send mail failed, err [%s], to[%s], cc[%s], subject[%s]", err.Error(), gostring.JoinSepStrings(",", m.To...), gostring.JoinSepStrings(",", m.Cc...), m.Subject)
    }
    return err
}
