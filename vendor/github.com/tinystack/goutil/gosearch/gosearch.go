// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gosearch

import (
    "sort"
)

func SearchIntSlice(a []int, x int) int {
    sort.Ints(a)
    l := sort.SearchInts(a, x)
    if l == len(a) {
        return -1
    }
    return l
}

func SearchStringSlice(s []string, x string) int {
    sort.Strings(s)
    l := sort.SearchStrings(s, x)
    if l == len(s) {
        return -1
    }
    return l
}
