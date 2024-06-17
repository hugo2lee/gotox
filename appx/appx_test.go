/*
 * @Author: hugo
 * @Date: 2024-05-11 15:05
 * @LastEditors: hugo
 * @LastEditTime: 2024-06-17 20:47
 * @FilePath: \gotox\appx\appx_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */

package appx_test

import (
	"testing"

	"github.com/hugo2lee/gotox/appx"
	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/ormx"
)

func TestNewApp(t *testing.T) {
	a := appx.New(configx.WithPath("../conf"))
	a.EnableDB(ormx.WithPostgres("test5"))
	a.EnableDB(ormx.WithMysql("test6"))
	t.Log(a.DBs)
}
