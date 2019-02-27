// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gofile

import (
    "io/ioutil"
    "os"
)

func CreateFile(filePath string, data []byte, perm os.FileMode) error {
    return ioutil.WriteFile(filePath, data, perm)
}
