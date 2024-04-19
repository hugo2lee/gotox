/*
 * @Author: hugo
 * @Date: 2024-03-19 19:57
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-19 16:44
 * @FilePath: \gotox\logx\log_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package logx_test

import (
	"testing"

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
)

func TestLogger(t *testing.T) {
	conf := configx.New(configx.WithPath("../conf"))
	logger := logx.New(conf)

	// build random log for check max size
	for i := 0; i < 10; i++ {
		logger.Debug("debug %v", i)
		logger.Info("info %v", i)
	}
	// logger.Debug("debug %v", "debug")
	// logger.Info("info %v", "info")
	// logger.Warn("warn %v", "warn")
	// logger.Error("error %v", errors.New("test"))
}
