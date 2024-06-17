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
	"log"
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

func New(opt ...configx.Option) *Appx {
	conf := configx.New(opt...)
	return &Appx{
		Configx: conf,
		Logger:  logx.New(conf),
	}
}

func WithConfigPath(path string) configx.Option {
	return configx.WithPath(path)
}

func WithConfigMode(mode string) configx.Option {
	return configx.WithMode(mode)
}

func (app *Appx) addResource(res resourcex.Resource) {
	resourcex.SetLogger(app.Logger)
	if app.ResourcexGroup == nil {
		app.ResourcexGroup = resourcex.NewResourcexGroup()
	}
	app.ResourcexGroup.AddResource(res)
}

func (app *Appx) EnableDB(ops ...ormx.Option) *Appx {
	if app.DBs == nil {
		orm, err := ormx.New(app.Configx, app.Logger, ops...)
		if err != nil {
			log.Fatalf("orm new failed, %+v", err)
		}
		app.DBs = orm
		app.addResource(app.DBs)
		app.Logger.Info("enable orm success")
	} else {
		// for _, name := range projectName {
		// 	if _, err := app.DBs.AddDB(name); err != nil {
		// 		log.Fatalf("add db %s failed, %+v", projectName, err)
		// 	}
		// }
		for _, op := range ops {
			if err := op(app.DBs); err != nil {
				log.Fatalf("add db failed, %+v", err)
			}
		}
	}

	return app
}

func (app *Appx) EnableCache() *Appx {
	ca := cachex.New(cachex.WithExpiration(time.Duration(app.Configx.CachexDefaultExpiration())*time.Second), cachex.WithCleanupInterval(time.Duration(app.Configx.CachexCleanupInterval())*time.Second))
	app.Cachex = ca

	app.addResource(app.Cachex)
	app.Logger.Info("enable cache success")
	return app
}

func (app *Appx) MigratTables(fns ...func() error) *Appx {
	for _, fn := range fns {
		if err := fn(); err != nil {
			log.Fatalf("init tables failed, %+v", err)
		}
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
	taskx.SetLogger(app.Logger)
	app.TaskxGroup = taskx.NewTaskxGroup()
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

	// 等待中断信号以优雅地关闭服务器
	<-notifyCtx.Done()
	app.Logger.Info("shutting down gracefully in 15 seconds..., press Ctrl+C again to force")
	timeOutCtx, timeOutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer timeOutCancel()

	if app.Serverx != nil {
		// http带超时关闭
		if err := app.Serverx.GracefullyDown(timeOutCtx); err != nil {
			app.Logger.Error("http server Shutdown error: %v \n", err)
		}
		app.Logger.Info("http server close")
	}

	{
		if app.ResourcexGroup != nil {
			app.ResourcexGroup.CloseAll(timeOutCtx)
		}
	}

	app.Logger.Info("App exiting")
	time.Sleep(2 * time.Second)
	app.Logger.Info("App exited")
}
