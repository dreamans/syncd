// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ssh

import (
    "net"
    "time"

    gossh "golang.org/x/crypto/ssh"
    "github.com/pkg/sftp"
)

type Conn struct {
    Addr        string
    User        string
    Password    string
    KeyBytes    []byte
    sshClient   *gossh.Client
}

func (c *Conn) Connect() error {
    var (
        auth    []gossh.AuthMethod
        err     error
        client  *gossh.Client
    )
    auth, err = c.makeAuth()
    if err != nil {
        return err
    }
    clientConfig := &gossh.ClientConfig{
        User: c.User,
        Auth: auth,
        Timeout: 30 * time.Second,
        HostKeyCallback: func(hostname string, remote net.Addr, key gossh.PublicKey) error {
            return nil
        },
    }
    if client, err = gossh.Dial("tcp", c.Addr, clientConfig); err != nil {
        return err
    }
    c.sshClient = client

    return nil
}

func (c Conn) NewSession() (*Session, error) {
    var (
        session *gossh.Session
        err     error
    )
    if session, err = c.sshClient.NewSession(); err != nil {
        return nil, err
    }
    return &Session{
        session: session,
    }, nil
}

func (c Conn) NewSftpClient() (*SftpClient, error) {
    sftpClient, err := sftp.NewClient(c.sshClient);
    if err != nil {
        return nil, err
    }
    return &SftpClient{
        client: sftpClient,
    }, nil
}

func (c Conn) makeAuth() ([]gossh.AuthMethod, error) {
    var (
        signer  gossh.Signer
        err     error
    )
    auth := make([]gossh.AuthMethod, 0)
    if len(c.KeyBytes) == 0 {
        auth = append(auth, gossh.Password(c.Password))
    } else {
        if c.Password == "" {
            signer, err = gossh.ParsePrivateKey(c.KeyBytes)
        } else {
            signer, err = gossh.ParsePrivateKeyWithPassphrase(c.KeyBytes, []byte(c.Password))
        }
        if err != nil {
            return nil, err
        }
        auth = append(auth, gossh.PublicKeys(signer))
    }
    return auth, nil
}
