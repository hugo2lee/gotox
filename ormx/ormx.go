/*
 * @Author: hugo
 * @Date: 2024-04-19 16:18
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-17 14:18
 * @FilePath: \gotox\ormx\ormx.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package ormx

import (
	"context"
	"sync"

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/resourcex"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var _ resourcex.Resource = (*ormx)(nil)

type ormx struct {
	gorm   *gorm.DB
	logger logx.Logger
}

func New(conf *configx.Configx, logCli logx.Logger) (*ormx, error) {
	dsn := conf.MysqlDsn()
	if dsn == "" {
		return nil, errors.New("mysql dsn is empty")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 使用 DEBUG 来打印
		// Logger: glogger.New(gormLoggerFunc(logCli.Debug),
		// 	glogger.Config{
		// 		SlowThreshold: 1 * time.Millisecond,
		// 		LogLevel:      glogger.Info,
		// 	}),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "mysql connect error")
	}

	switch conf.Mode() {
	case configx.RUNDEV:
		db = db.Debug()
	case configx.RUNTEST:
		db = db.Debug()
	default:
		db = db.Debug()
	}

	return &ormx{db, logCli}, nil
}

func (c *ormx) DB() *gorm.DB {
	return c.gorm
}

func (c *ormx) Name() string {
	return "orm"
}

func (c *ormx) Close(ctx context.Context, wg *sync.WaitGroup) {
	db, err := c.gorm.DB()
	if err != nil {
		c.logger.Error("gorm DB get %v", err)
		return
	}
	if err := db.Close(); err != nil {
		c.logger.Error("gorm close %v", err)
		return
	}
	wg.Done()
	c.logger.Info("%s close", c.Name())
}

type gormLoggerFunc func(msg string, args ...any)

func (g gormLoggerFunc) Printf(msg string, args ...any) {
	g(msg, args)
}
