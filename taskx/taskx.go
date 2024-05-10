/*
 * @Author: hugo
 * @Date: 2024-05-10 15:06
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-10 19:48
 * @FilePath: \gotox\taskx\taskx.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package taskx

import (
	"context"
	"log"
)

type Tasker interface {
	// task with context, run in new goroutine, context is canceled when interrupt
	Run(context.Context)
}

type TaskFunc func(context.Context)

type Task struct {
	Name string
	Func TaskFunc
}

func New(name string, f TaskFunc) *Task {
	return &Task{
		Name: name,
		Func: f,
	}
}

func (t *Task) Run(ctx context.Context) {
	go t.Func(ctx)
	log.Printf("task \"%s\" is running \n", t.Name)
}
