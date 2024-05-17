/*
 * @Author: hugo
 * @Date: 2024-03-14 15:44
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-17 15:01
 * @FilePath: \gotox\logx\log.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */

package logx

import (
	"github.com/hugo2lee/gotox/configx"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

var _ Logger = (*Logx)(nil)

var Log Logger

type Logx struct {
	logger *zap.Logger
}

func New(conf *configx.Configx) *Logx {
	zaplog := zap.New(zapLoggerBuilder(conf.LogDir(), conf.Mode()), zap.AddCaller(), zap.AddCallerSkip(1))

	// var err error
	// var zaplog *zap.Logger
	// switch conf.RunMode() {
	// case config.RUNDEV:
	// 	zaplog, err = zap.NewDevelopment(zap.AddCallerSkip(1))
	// case config.RUNTEST:
	// 	zaplog, err = zap.NewDevelopment(zap.AddCallerSkip(1))
	// case config.RUNPROD:
	// 	zaplog, err = zap.NewProduction(zap.AddCallerSkip(1))
	// default:
	// 	zaplog, err = zap.NewProduction(zap.AddCallerSkip(1))
	// }
	// if err != nil {
	// 	log.Fatalf("zap.NewProduction error: %s \n", err)
	// }

	// zap.ReplaceGlobals(zaplog)

	cli := &Logx{
		logger: zaplog,
	}

	Log = cli

	Log.Info("Logger is ready")

	return cli
}

func (l *Logx) Debug(msg string, args ...any) {
	l.logger.Sugar().Debugf(msg, args...)
}

func (l *Logx) Info(msg string, args ...any) {
	l.logger.Sugar().Infof(msg, args...)
}

func (l *Logx) Warn(msg string, args ...any) {
	l.logger.Sugar().Warnf(msg, args...)
}

func (l *Logx) Error(msg string, args ...any) {
	l.logger.Sugar().Errorf(msg, args...)
}

// 这里是为了包变量Log初始化
func init() {
	Log = NewNoOpLogger()
}
