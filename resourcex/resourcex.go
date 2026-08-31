/*
 * @Author: hugo
 * @Date: 2024-05-11 17:17
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-17 15:19
 * @FilePath: \gotox\resourcex\resourcex.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package resourcex

import (
	"context"
	"sync"

	"github.com/hugo2lee/gotox/logx"
)

type Resource interface {
	Name() string
	Close(context.Context, *sync.WaitGroup)
}

type Resourcex struct {
	name    string
	closeFn func(ctx context.Context)
}

func NewResourcex(name string, fn func(ctx context.Context)) *Resourcex {
	return &Resourcex{name: name, closeFn: fn}
}

func (r *Resourcex) Name() string { return r.name }

func (r *Resourcex) Close(ctx context.Context, wg *sync.WaitGroup) {
	r.closeFn(ctx)
	wg.Done()
}

type GroupOption func(*ResourcexGroup)

func WithLogger(logger logx.Logger) GroupOption {
	return func(g *ResourcexGroup) {
		if logger != nil {
			g.logger = logger
		}
	}
}

type ResourcexGroup struct {
	resources map[string]Resource
	logger    logx.Logger
}

func NewResourcexGroup(opts ...GroupOption) *ResourcexGroup {
	g := &ResourcexGroup{
		resources: make(map[string]Resource),
		logger:    logx.NewNoOpLogger(),
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

func (r *ResourcexGroup) AddResource(res ...Resource) {
	for _, f := range res {
		r.resources[f.Name()] = f
	}
}

func (r *ResourcexGroup) CloseAll(ctx context.Context) {
	wg := new(sync.WaitGroup)
	wg.Add(len(r.resources))
	for _, f := range r.resources {
		go f.Close(ctx, wg)
	}
	wgChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgChan)
	}()

	select {
	case <-ctx.Done():
		r.logger.Info("close resource timeout")
	case <-wgChan:
		r.logger.Info("all resource closed")
	}
}
