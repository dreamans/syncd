// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/dreamans/syncd"
	"github.com/dreamans/syncd/router/route"
)

func main() {
	app := syncd.NewApp()

	route.RegisterRoute(app.Gin)

	app.Start()
}