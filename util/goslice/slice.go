// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goslice

func FilterSliceInt(sl []int) []int {
    var newSlice []int
    for _, s := range sl {
        if s != 0 {
            newSlice = append(newSlice, s)
        }
    }
    return newSlice
}

func InSlice(v interface{}, sl []interface{}) bool {
    for _, vv := range sl {
        if vv == v {
            return true
        }
    }
    return false
}

func InSliceInt(v int, sl []int) bool {
    for _, vv := range sl {
        if vv == v {
            return true
        }
    }
    return false
}

func InSliceString(v string, sl []string) bool {
    for _, vv := range sl {
        if vv == v {
            return true
        }
    }
    return false
}

func SliceIntersect(slice1, slice2 []interface{}) (diffslice []interface{}) {
    for _, v := range slice1 {
        if InSlice(v, slice2) {
            diffslice = append(diffslice, v)
        }
    }
    return
}

func SliceIntersectInt(slice1, slice2 []int) (diffslice []int) {
    for _, v := range slice1 {
        if InSliceInt(v, slice2) {
            diffslice = append(diffslice, v)
        }
    }
    return
}

func SliceComplementaryInt(slice1, slice2 []int) (diffslice []int) {
    for _, v := range slice1 {
        if !InSliceInt(v, slice2) {
            diffslice = append(diffslice, v)
        }
    }
    return
}