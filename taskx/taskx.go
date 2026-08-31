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

type Task interface {
	Name() string
	Run(context.Context)
}

type Taskx struct {
	name string
	fn   func(ctx context.Context)
}

func NewTaskx(name string, fn func(ctx context.Context)) *Taskx {
	return &Taskx{name: name, fn: fn}
}

func (t *Taskx) Name() string { return t.name }
func (t *Taskx) Run(ctx context.Context) { t.fn(ctx) }

type GroupOption func(*TaskxGroup)

func WithLogger(logger logx.Logger) GroupOption {
	return func(g *TaskxGroup) {
		if logger != nil {
			g.logger = logger
		}
	}
}

type TaskxGroup struct {
	tasks  map[string]Task
	logger logx.Logger
}

func NewTaskxGroup(opts ...GroupOption) *TaskxGroup {
	g := &TaskxGroup{
		tasks:  make(map[string]Task),
		logger: logx.NewNoOpLogger(),
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

func (t *TaskxGroup) AddTask(tasks ...Task) {
	for _, f := range tasks {
		t.tasks[f.Name()] = f
	}
}

func (t *TaskxGroup) Run(ctx context.Context) {
	for _, f := range t.tasks {
		go f.Run(ctx)
		t.logger.Info("task \"%s\" is running \n", f.Name())
	}
}
