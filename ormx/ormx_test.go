/*
 * @Author: hugo
 * @Date: 2024-04-19 16:24
 * @LastEditors: hugo
 * @LastEditTime: 2024-06-17 20:52
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
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestDefaultDB(t *testing.T) {
	t.Parallel()
	conf := configx.New(configx.WithPath("../conf"))
	dbGorm, err := ormx.New(conf, logx.Log)
	assert.NoError(t, err)

	db := dbGorm.GetDB()

	type User struct {
		Name string
	}
	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)

	err = db.Model(&User{}).Create(&User{Name: "hugo"}).Error
	assert.NoError(t, err)
}

func TestMysqlDB(t *testing.T) {
	t.Parallel()
	conf := configx.New(configx.WithPath("../conf"))

	pj := "test7"
	dbGorm, err := ormx.New(conf, logx.Log, ormx.WithMysql(pj))
	assert.NoError(t, err)

	db := dbGorm.GetDB(pj)

	type User struct {
		Name string
	}
	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)

	err = db.Model(&User{}).Create(&User{Name: "hugo"}).Error
	assert.NoError(t, err)
}

func TestPgDB(t *testing.T) {
	t.Parallel()
	conf := configx.New(configx.WithPath("../conf"))

	pj := "test2"
	dbGorm, err := ormx.New(conf, logx.Log, ormx.WithPostgres(pj))
	assert.NoError(t, err)

	db := dbGorm.GetDB(pj)

	type User struct {
		Name string
	}
	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)

	err = db.Model(&User{}).Create(&User{Name: "hugo"}).Error
	assert.NoError(t, err)
}

func TestGormMysql(t *testing.T) {
	t.Parallel()
	dsn := "root:root@tcp(localhost:3306)/dev?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名
			},
		})
	assert.NoError(t, err)

	type User struct {
		Name string
	}

	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)
}

func TestGormPostgres(t *testing.T) {
	t.Parallel()
	dsn := "postgres://root:root@localhost:5432/dev?search_path=dev-schema&sslmode=disable&timezone=Asia/Shanghai"

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名
			},
		})
	assert.NoError(t, err)

	type User struct {
		Name string
	}

	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)
}
