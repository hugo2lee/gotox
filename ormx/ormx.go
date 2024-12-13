/*
 * @Author: hugo
 * @Date: 2024-04-19 16:18
 * @LastEditors: hugo2lee
 * @LastEditTime: 2024-12-13 14:05
 * @FilePath: \gotox\ormx\ormx.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package ormx

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/resourcex"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var _ resourcex.Resource = (*Ormx)(nil)

const (
	DefaultProjectName         = DefaultMysqlProjectName
	MYSQL                      = "mysql"
	DefaultMysqlProjectName    = ""
	POSTGRES                   = "postgres"
	DefaultPostgresProjectName = "public"
)

type BaseModel struct {
	ID      uint  `gorm:"primaryKey;autoIncrement"`
	Created int64 `gorm:"autoCreateTime:milli"`
	Updated int64 `gorm:"autoUpdateTime:milli"`
	Deleted gorm.DeletedAt
	UUID    string `gorm:"size:36;uniqueIndex"`
}

type Option func(*Ormx) error

type Ormx struct {
	conf   *configx.Configx
	logger logx.Logger
	gorms  map[string]*gorm.DB
}

func New(conf *configx.Configx, logCli logx.Logger, ops ...Option) (*Ormx, error) {
	orm := &Ormx{
		conf,
		logCli,
		make(map[string]*gorm.DB),
	}

	if ops == nil {
		if err := WithMysql(DefaultProjectName)(orm); err != nil {
			return nil, err
		}
		return orm, nil
	}

	for _, op := range ops {
		if err := op(orm); err != nil {
			return nil, err
		}
	}

	return orm, nil
}

func (orm *Ormx) dialDB(dialer gorm.Dialector, tablePrefix string) (*gorm.DB, error) {
	db, err := gorm.Open(dialer, &gorm.Config{
		// 使用 DEBUG 来打印
		// Logger: glogger.New(gormLoggerFunc(logCli.Debug),
		// 	glogger.Config{
		// 		SlowThreshold: 1 * time.Millisecond,
		// 		LogLevel:      glogger.Info,
		// 	}),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,        // 使用单数表名
			TablePrefix:   tablePrefix, // 表名前缀
		},
	})
	if err != nil {
		return nil, err
	}
	bareDb, err := db.DB()
	if err != nil {
		return nil, err
	}
	bareDb.SetMaxIdleConns(10)
	bareDb.SetMaxOpenConns(100)
	if err := bareDb.Ping(); err != nil {
		return nil, err
	}

	// make sure schema exists
	if dialer.Name() == POSTGRES {
		ns := strings.SplitN(tablePrefix, ".", 2)
		if len(ns) == 2 {
			var exists bool
			if err := db.Raw("SELECT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = ?)", ns[0]).Scan(&exists).Error; err != nil {
				return nil, err
			}
			orm.logger.Info("schema %s exists %v", ns[0], exists)
			if !exists {
				str := fmt.Sprintf("CREATE SCHEMA %s", ns[0])
				if err := db.Exec(str).Error; err != nil {
					return nil, err
				}
				orm.logger.Info("schema %s created", ns[0])
			}
		}
	}

	return db, nil
}

func WithMysql(projectName ...string) Option {
	return func(o *Ormx) error {
		for _, project := range projectName {
			dsn := o.conf.MysqlDsn()
			if dsn == "" {
				return errors.New("mysql dsn is empty")
			}
			dl := mysql.Open(dsn)

			prefix := ""
			if project != "" {
				prefix = fmt.Sprintf("%s_", project)
			}
			db, err := o.dialDB(dl, prefix)
			if err != nil {
				return errors.Wrapf(err, "dial db %s failed", project)
			}

			switch o.conf.Mode() {
			case configx.RUNDEV:
				db = db.Debug()
			case configx.RUNTEST:
				db = db.Debug()
			case configx.RUNPROD:
				// db = db
			default:
				db = db.Debug()
			}

			o.gorms[project] = db
		}

		return nil
	}
}

func WithPostgres(projectName ...string) Option {
	return func(o *Ormx) error {
		for _, project := range projectName {
			dsn := o.conf.PostgresDsn()
			if dsn == "" {
				return errors.New("postgres dsn is empty")
			}
			dl := postgres.Open(dsn)

			prefix := ""
			if project != "" {
				prefix = fmt.Sprintf("%s.", project)
			}
			db, err := o.dialDB(dl, prefix)
			if err != nil {
				return errors.Wrapf(err, "dial db.schema %s failed", project)
			}

			switch o.conf.Mode() {
			case configx.RUNDEV:
				db = db.Debug()
			case configx.RUNTEST:
				db = db.Debug()
			default:
				db = db.Debug()
			}
			o.gorms[project] = db
		}
		return nil
	}
}

func (c *Ormx) GetDB(projectName ...string) *gorm.DB {
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
		if name == DefaultMysqlProjectName {
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
