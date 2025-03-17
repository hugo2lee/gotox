/*
 * @Author: hugo
 * @Date: 2024-04-19 16:24
 * @LastEditors: hugo2lee
 * @LastEditTime: 2024-12-13 13:59
 * @FilePath: \gotox\ormx\ormx_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */

package ormx_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
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

func TestWithMysqlMultipleDb(t *testing.T) {
	t.Parallel()
	conf := configx.New(configx.WithPath("../conf"))

	db1 := "test1"
	db2 := "test2"

	dbGorm, err := ormx.New(conf, logx.Log, ormx.WithMysqlMultipleDb(db1, db2))
	assert.NoError(t, err)

	type User struct {
		Name string
	}

	dbEn1 := dbGorm.GetDB(db1)

	err = dbEn1.AutoMigrate(&User{})
	assert.NoError(t, err)

	err = dbEn1.Model(&User{}).Create(&User{Name: "hugo-db1"}).Error
	assert.NoError(t, err)

	dbEn2 := dbGorm.GetDB(db2)

	err = dbEn2.AutoMigrate(&User{})
	assert.NoError(t, err)

	err = dbEn2.Model(&User{}).Create(&User{Name: "hugo-db2"}).Error
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
		Id   int64
		Name string
	}

	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)
}

type BaseModel struct {
	ID      uint  `gorm:"primaryKey;autoIncrement"`
	Created int64 `gorm:"autoCreateTime:milli"`
	Updated int64 `gorm:"autoUpdateTime:milli"`
	Deleted gorm.DeletedAt
	UUID    string `gorm:"size:36;uniqueIndex"`
}

func (u *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	uu, err := uuid.NewV7()
	if err != nil {
		return err
	}
	u.UUID = uu.String()
	return nil
}

func TestGormBaseModel(t *testing.T) {
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

	db = db.Debug()

	type User struct {
		BaseModel
		Name string
	}

	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)

	err = db.Create(&User{Name: "erica"}).Error
	assert.NoError(t, err)

	err = db.Where("name = ?", "hugo").Delete(&User{}).Error
	assert.NoError(t, err)

	var user User
	err = db.Where("name = ?", "erica").First(&user).Error
	assert.NoError(t, err)

	user.Name = "hugo"
	err = db.Save(&user).Error
	assert.NoError(t, err)

	fmt.Println("=== end ===")
}
