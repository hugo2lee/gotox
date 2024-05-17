/*
 * @Author: hugo
 * @Date: 2024-05-10 15:06
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-17 16:48
 * @FilePath: \gotox\taskx\taskx.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package taskx

import (
	"context"

	"github.com/hugo2lee/gotox/logx"
)

var logg logx.Logger = logx.Log

// // 自定义的logger，建议实例化赋予
func SetLogger(l logx.Logger) {
	logg = l
}

type Task interface {
	Name() string
	// task with context, run in new goroutine, context is canceled when interrupt
	Run(context.Context)
}

type Taskx struct {
	name string
	fn   func(ctx context.Context)
}

func NewTaskx(name string, fn func(ctx context.Context)) *Taskx {
	return &Taskx{
		name: name,
		fn:   fn,
	}
}

func (t *Taskx) Name() string {
	return t.name
}

func (t *Taskx) Run(ctx context.Context) {
	t.fn(ctx)
}

type TaskxGroup struct {
	tasks map[string]Task
}

func NewTaskxGroup() *TaskxGroup {
	return &TaskxGroup{
		tasks: make(map[string]Task),
	}
}

func (t *TaskxGroup) AddTask(tasks ...Task) {
	for _, f := range tasks {
		t.tasks[f.Name()] = f
	}
}

func (t *TaskxGroup) Run(ctx context.Context) {
	for _, f := range t.tasks {
		go f.Run(ctx)
		logg.Info("task \"%s\" is running \n", f.Name())
	}
}
