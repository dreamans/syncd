// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ssh

import (
    "os"
    "io"

    "github.com/pkg/sftp"
)

type SftpClient struct {
    client  *sftp.Client
}

func (c SftpClient) SendFile(localFile, dstFile string) error {
    var (
        srcFileHandle   *os.File
        dstFileHandle   *sftp.File
        err             error
    )
    if srcFileHandle, err = os.Open(localFile); err != nil {
        return err
    }
    defer srcFileHandle.Close()

    //dstPath := path.Base(dstFile)
    dstFileHandle, err = c.client.Create(dstFile)
    if err != nil {
        return err
    }
    defer dstFileHandle.Close()

    buf := make([]byte, 1024)
    for {
        n, err := srcFileHandle.Read(buf)
        if err != nil && err != io.EOF {
            return err
        }
        if n == 0 {
            break
        }
        if _, err := dstFileHandle.Write(buf[:n]); err != nil {
            return err
        }
    }
    return nil
}

func (c SftpClient) Close() {
    if c.client != nil {
        c.Close()
    }
}
