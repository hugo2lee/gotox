/*
 * @Author: hugo
 * @Date: 2024-03-20 17:17
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-02 15:02
 * @FilePath: \gotox\logx\zapBuilder.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package logx

import (
	"fmt"
	"os"
	"path"

	"github.com/hugo2lee/gotox/configx"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DEFAULTLOGFILE  = "./log/dev.log"
	DEFAULTLOGLEVEl = "dev"
)

func zapLoggerBuilder(filePath, env string) (zapcore.Core, error) {
	if filePath == "" {
		filePath = DEFAULTLOGFILE
	}
	if env == "" {
		env = DEFAULTLOGLEVEl
	}

	fileWriter, err := getLogWriter(filePath)
	if err != nil {
		return nil, err
	}
	writer := zapcore.NewMultiWriteSyncer(fileWriter, zapcore.AddSync(os.Stdout))
	encoder := getLogEncoder(env)
	level := getLogLevel(env)
	return zapcore.NewCore(encoder, writer, level), nil
}

func getLogWriter(filePath string) (zapcore.WriteSyncer, error) {
	if err := os.MkdirAll(path.Dir(filePath), 0o755); err != nil {
		return nil, fmt.Errorf("create log directory %s: %w", path.Dir(filePath), err)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    100,
		MaxBackups: 50,
		MaxAge:     30,
		Compress:   false,
	}

	return zapcore.AddSync(lumberJackLogger), nil
}

func getLogEncoder(env string) zapcore.Encoder {
	var encoderConfig zapcore.EncoderConfig
	switch env {
	case configx.RUNDEV:
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	case configx.RUNTEST:
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	case configx.RUNPROD:
		encoderConfig = zap.NewProductionEncoderConfig()
	default:
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogLevel(env string) zapcore.Level {
	switch env {
	case configx.RUNDEV:
		return zapcore.DebugLevel
	case configx.RUNTEST:
		return zapcore.DebugLevel
	case configx.RUNPROD:
		return zapcore.InfoLevel
	default:
		return zapcore.DebugLevel
	}
}
