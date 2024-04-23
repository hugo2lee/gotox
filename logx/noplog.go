/*
 * @Author: hugo
 * @Date: 2024-03-25 19:37
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-02 14:14
 * @FilePath: \gotox\logx\noplog.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */

// 无任何依赖的包变量logger，适用于直接导包的场景，解决实例化结构体没logger的问题
package logx

import "log"

var _ Logger = (*NoOpLogger)(nil)

type NoOpLogger struct{}

func NewNoOpLogger() Logger {
	return &NoOpLogger{}
}

func (n *NoOpLogger) Debug(msg string, args ...any) {
	log.Println(msg, args)
}

func (n *NoOpLogger) Info(msg string, args ...any) {
	log.Println(msg, args)
}

func (n *NoOpLogger) Warn(msg string, args ...any) {
	log.Println(msg, args)
}

func (n *NoOpLogger) Error(msg string, args ...any) {
	log.Println(msg, args)
}
