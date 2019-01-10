// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goweb

func filterFlags(content string) string {
    for i, char := range content {
        if char == ' ' || char == ';' {
            return content[:i]
        }
    }
    return content
}

