/*
 * @Author: hugo
 * @Date: 2024-05-11 15:05
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-11 17:12
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

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/ormx"
	"github.com/hugo2lee/gotox/resourcex"
	"github.com/hugo2lee/gotox/serverx"
	"github.com/hugo2lee/gotox/taskx"
	"github.com/hugo2lee/gotox/webx"
	"gorm.io/gorm"
)

type App struct {
	*configx.ConfigCli
	logx.Logger
	*gorm.DB
	*serverx.Server
	*resourcex.Resourcer
	*taskx.Tasker
}

func NewApp(opt ...configx.Option) *App {
	conf := configx.New(opt...)
	return &App{
		ConfigCli: conf,
		Logger:    logx.New(conf),
	}
}

func WithConfigPath(path string) configx.Option {
	return configx.WithPath(path)
}

func WithConfigMode(mode string) configx.Option {
	return configx.WithMode(mode)
}

func (app *App) addResource(res resourcex.Resource) {
	if app.Resourcer == nil {
		app.Resourcer = resourcex.NewResourcer()
	}
	app.Resourcer.AddResource(res)
}

func (app *App) EnableDB() *App {
	orm, err := ormx.New(app.ConfigCli, app.Logger)
	if err != nil {
		log.Fatalf("orm new failed, %+v", err)
	}
	app.DB = orm.DB()
	app.Logger.Info("init InitDependency success")

	app.addResource(orm)
	return app
}

func (app *App) InitTables(fns ...func(*gorm.DB) error) *App {
	for _, fn := range fns {
		if err := fn(app.DB); err != nil {
			log.Fatalf("init tables failed, %+v", err)
		}
	}
	return app
}

func (app *App) RegisterServies(fns ...webx.Handler) *App {
	for _, fn := range fns {
		fn.RegisterRouter(app.Engine)
	}
	app.Logger.Info("bind router success")
	return app
}

func (app *App) EnableWebServer() *App {
	srv := serverx.New(app.ConfigCli, app.Logger).
		EnableAccessLog().
		EnableWrapLog().
		EnableAuth()
	app.Server = srv
	app.Logger.Info("init server success")

	return app
}

func (app *App) EnableTasks(tasks ...taskx.Task) *App {
	app.Tasker = taskx.NewTasker()
	for _, task := range tasks {
		app.Tasker.AddTask(task)
	}
	app.Logger.Info("enable task success")
	return app
}

func (app *App) Run() {
	notifyCtx, notifyStop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT, syscall.SIGTERM)
	defer notifyStop()

	if app.Tasker != nil {
		app.Tasker.Run(notifyCtx)
	}

	if app.Server != nil {
		go app.Server.GracefullyUp(notifyStop)
	}

	// 等待中断信号以优雅地关闭服务器
	<-notifyCtx.Done()
	app.Logger.Info("shutting down gracefully in 15 seconds..., press Ctrl+C again to force")
	timeOutCtx, timeOutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer timeOutCancel()

	if app.Server != nil {
		// http带超时关闭
		if err := app.Server.GracefullyDown(timeOutCtx); err != nil {
			app.Logger.Error("http server Shutdown error: %v \n", err)
		}
		app.Logger.Info("http server close")
	}

	{
		if app.Resourcer != nil {
			app.Resourcer.CloseAll(timeOutCtx)
		}
	}

	app.Logger.Info("App exiting")
	time.Sleep(2 * time.Second)
	app.Logger.Info("App exited")
}
