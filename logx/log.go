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

type Logx struct {
	logger *zap.Logger
}

// Open constructs the configured logger and returns filesystem/setup errors to
// the caller.
func Open(conf *configx.Configx) (*Logx, error) {
	core, err := zapLoggerBuilder(conf.LogDir(), conf.Mode())
	if err != nil {
		return nil, err
	}
	zaplog := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	cli := &Logx{logger: zaplog}
	cli.Info("Logger is ready")
	return cli, nil
}

// New is the fail-fast compatibility constructor.
//
// Prefer Open when the caller needs explicit error handling. New panics on
// initialization failure rather than terminating the process.
func New(conf *configx.Configx) *Logx {
	cli, err := Open(conf)
	if err != nil {
		panic(err)
	}
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
