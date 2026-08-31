/*
 * @Author: hugo
 * @Date: 2024-05-11 15:05
 * @LastEditors: hugo
 * @LastEditTime: 2024-06-17 20:45
 * @FilePath: \gotox\appx\appx.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package appx

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hugo2lee/gotox/cachex"
	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/ormx"
	"github.com/hugo2lee/gotox/resourcex"
	"github.com/hugo2lee/gotox/serverx"
	"github.com/hugo2lee/gotox/taskx"
	"github.com/hugo2lee/gotox/webx"
)

type Appx struct {
	Configx        *configx.Configx
	Logger         logx.Logger
	Cachex         cachex.Cachexer
	DBs            *ormx.Ormx
	Serverx        *serverx.Serverx
	ResourcexGroup *resourcex.ResourcexGroup
	TaskxGroup     *taskx.TaskxGroup
}

// Build constructs Appx and returns configuration/logger initialization errors
// to the caller. Prefer this path when the composition root owns failure policy.
func Build(opt ...configx.Option) (*Appx, error) {
	conf, err := configx.Load(opt...)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	logger, err := logx.Open(conf)
	if err != nil {
		return nil, fmt.Errorf("open logger: %w", err)
	}
	return &Appx{
		Configx: conf,
		Logger:  logger,
	}, nil
}

// New is the fail-fast compatibility constructor.
//
// Prefer Build when the caller needs explicit error handling.
func New(opt ...configx.Option) *Appx {
	app, err := Build(opt...)
	if err != nil {
		panic(err)
	}
	return app
}

func WithConfigPath(path string) configx.Option { return configx.WithPath(path) }
func WithConfigMode(mode string) configx.Option { return configx.WithMode(mode) }

func (app *Appx) addResource(res resourcex.Resource) {
	if app.ResourcexGroup == nil {
		app.ResourcexGroup = resourcex.NewResourcexGroup(resourcex.WithLogger(app.Logger))
	}
	app.ResourcexGroup.AddResource(res)
}

// TryEnableDB enables or extends ORM configuration and returns failures to the
// caller instead of deciding process lifetime inside the library.
func (app *Appx) TryEnableDB(ops ...ormx.Option) error {
	if app.DBs == nil {
		orm, err := ormx.Dial(app.Configx, app.Logger, ops...)
		if err != nil {
			return fmt.Errorf("dial orm: %w", err)
		}
		app.DBs = orm
		app.addResource(app.DBs)
		app.Logger.Info("enable orm success")
		return nil
	}

	for _, op := range ops {
		if err := op(app.DBs); err != nil {
			return fmt.Errorf("add db: %w", err)
		}
	}
	return nil
}

// EnableDB is the fail-fast compatibility fluent API.
//
// Prefer TryEnableDB when the caller needs explicit error handling.
func (app *Appx) EnableDB(ops ...ormx.Option) *Appx {
	if err := app.TryEnableDB(ops...); err != nil {
		panic(err)
	}
	return app
}

func (app *Appx) EnableCache() *Appx {
	ca := cachex.New(
		cachex.WithExpiration(time.Duration(app.Configx.CachexDefaultExpiration())*time.Second),
		cachex.WithCleanupInterval(time.Duration(app.Configx.CachexCleanupInterval())*time.Second),
	)
	app.Cachex = ca
	app.addResource(app.Cachex)
	app.Logger.Info("enable cache success")
	return app
}

// MigrateTables executes table migration callbacks and reports the first
// failure to the caller.
func (app *Appx) MigrateTables(fns ...func() error) error {
	for i, fn := range fns {
		if err := fn(); err != nil {
			return fmt.Errorf("migrate tables callback %d: %w", i, err)
		}
	}
	return nil
}

// MigratTables is kept for source compatibility with the original misspelled
// fluent API.
//
// Deprecated: use MigrateTables for explicit error handling.
func (app *Appx) MigratTables(fns ...func() error) *Appx {
	if err := app.MigrateTables(fns...); err != nil {
		panic(err)
	}
	return app
}

func (app *Appx) RegisterServies(fns ...func() webx.Handler) *Appx {
	for _, fn := range fns {
		fn().RegisterRouter(app.Serverx.Engine)
	}
	app.Logger.Info("bind router success")
	return app
}

func (app *Appx) EnableWebServer() *Appx {
	srv := serverx.New(app.Configx, app.Logger).
		EnableAccessLog().
		EnableWrapLog().
		EnableAuth()
	app.Serverx = srv
	app.Logger.Info("init server success")
	return app
}

func (app *Appx) EnableTasks(taskGenFuncs ...func() taskx.Task) *Appx {
	app.TaskxGroup = taskx.NewTaskxGroup(taskx.WithLogger(app.Logger))
	for _, taskGen := range taskGenFuncs {
		app.TaskxGroup.AddTask(taskGen())
	}
	app.Logger.Info("enable task success")
	return app
}

func (app *Appx) Run() {
	notifyCtx, notifyStop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT, syscall.SIGTERM)
	defer notifyStop()

	if app.TaskxGroup != nil {
		app.TaskxGroup.Run(notifyCtx)
	}

	if app.Serverx != nil {
		go app.Serverx.GracefullyUp(notifyStop)
	}

	<-notifyCtx.Done()
	app.Logger.Info("shutting down gracefully in 15 seconds..., press Ctrl+C again to force")
	timeOutCtx, timeOutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer timeOutCancel()

	if app.Serverx != nil {
		if err := app.Serverx.GracefullyDown(timeOutCtx); err != nil {
			app.Logger.Error("http server Shutdown error: %v \n", err)
		}
		app.Logger.Info("http server close")
	}

	if app.ResourcexGroup != nil {
		if err := app.ResourcexGroup.CloseAll(timeOutCtx); err != nil {
			app.Logger.Error("close resources error: %v", err)
		}
	}

	app.Logger.Info("App exiting")
	time.Sleep(2 * time.Second)
	app.Logger.Info("App exited")
}
