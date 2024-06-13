/*
 * @Author: hugo
 * @Date: 2024-04-19 16:18
 * @LastEditors: hugo
 * @LastEditTime: 2024-06-13 16:54
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

var _ resourcex.Resource = (*Ormx)(nil)

const DefaultProjectName = ""

type Ormx struct {
	conf   *configx.Configx
	logger logx.Logger
	gorms  map[string]*gorm.DB
}

func New(conf *configx.Configx, logCli logx.Logger) (*Ormx, error) {
	or := &Ormx{
		conf,
		logCli,
		make(map[string]*gorm.DB),
	}
	return or.Add(DefaultProjectName)
}

func (o *Ormx) Add(projectName string) (*Ormx, error) {
	dsn := o.conf.MysqlDsn()
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
			SingularTable: true,        // 使用单数表名
			TablePrefix:   projectName, // 表名前缀
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "mysql connect error")
	}

	switch o.conf.Mode() {
	case configx.RUNDEV:
		db = db.Debug()
	case configx.RUNTEST:
		db = db.Debug()
	default:
		db = db.Debug()
	}

	if o.gorms == nil {
		o.gorms = make(map[string]*gorm.DB)
	}

	o.gorms[projectName] = db

	return o, nil
}

func (c *Ormx) DB(projectName ...string) *gorm.DB {
	if len(projectName) == 0 {
		return c.gorms[DefaultProjectName]
	}
	return c.gorms[projectName[0]]
}

func (c *Ormx) Name() string {
	return "orm"
}

func (c *Ormx) Close(ctx context.Context, wg *sync.WaitGroup) {
	for name, gor := range c.gorms {
		if name == DefaultProjectName {
			name = "default"
		}
		db, err := gor.DB()
		if err != nil {
			c.logger.Error("gorm %s DB get %v", name, err)
			return
		}
		if err := db.Close(); err != nil {
			c.logger.Error("gorm %s close %v", name, err)
			return
		}
	}

	wg.Done()
	c.logger.Info("%s close", c.Name())
}

type gormLoggerFunc func(msg string, args ...any)

func (g gormLoggerFunc) Printf(msg string, args ...any) {
	g(msg, args)
}
