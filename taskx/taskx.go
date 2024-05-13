/*
 * @Author: hugo
 * @Date: 2024-05-10 15:06
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-13 19:58
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

type TaskCli struct {
	name string
	fn   func(ctx context.Context)
}

func NewTaskCli(name string, fn func(ctx context.Context)) *TaskCli {
	return &TaskCli{
		name: name,
		fn:   fn,
	}
}

func (t *TaskCli) Name() string {
	return t.name
}

func (t *TaskCli) Run(ctx context.Context) {
	t.fn(ctx)
}

type Tasker struct {
	tasks map[string]Task
}

func NewTasker() *Tasker {
	return &Tasker{
		tasks: make(map[string]Task),
	}
}

func (t *Tasker) AddTask(tasks ...Task) {
	for _, f := range tasks {
		t.tasks[f.Name()] = f
	}
}

func (t *Tasker) Run(ctx context.Context) {
	for _, f := range t.tasks {
		go f.Run(ctx)
		logg.Info("task \"%s\" is running \n", f.Name())
	}
}
