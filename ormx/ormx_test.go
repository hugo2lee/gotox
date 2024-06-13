/*
 * @Author: hugo
 * @Date: 2024-04-19 16:24
 * @LastEditors: hugo
 * @LastEditTime: 2024-06-13 20:42
 * @FilePath: \gotox\ormx\ormx_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package ormx_test

import (
	"testing"

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/ormx"
	"github.com/stretchr/testify/assert"
)

func TestMysql(t *testing.T) {
	t.Parallel()
	conf := configx.New(configx.WithPath("../conf"))
	dbGorm, err := ormx.New(conf, logx.Log, "test")
	assert.NoError(t, err)

	db := dbGorm.GetDB("test")

	type User struct {
		Name string
	}
	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)

	err = db.Model(&User{}).Create(&User{Name: "hugo"}).Error
	assert.NoError(t, err)
}
