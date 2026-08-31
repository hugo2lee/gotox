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
	"errors"
	"fmt"
	"sync"

	"github.com/hugo2lee/gotox/logx"
)

type Resource interface {
	Name() string
	Close(context.Context) error
}

type Resourcex struct {
	name    string
	closeFn func(ctx context.Context)
}

func NewResourcex(name string, fn func(ctx context.Context)) *Resourcex {
	return &Resourcex{name: name, closeFn: fn}
}

func (r *Resourcex) Name() string { return r.name }

func (r *Resourcex) Close(ctx context.Context) error {
	r.closeFn(ctx)
	return nil
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

func (r *ResourcexGroup) CloseAll(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(r.resources))

	wg.Add(len(r.resources))
	for name, resource := range r.resources {
		name := name
		resource := resource
		go func() {
			defer wg.Done()
			if err := resource.Close(ctx); err != nil {
				errCh <- fmt.Errorf("close resource %q: %w", name, err)
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		r.logger.Info("close resource timeout")
		return ctx.Err()
	case <-done:
		close(errCh)
	}

	var closeErrs []error
	for err := range errCh {
		closeErrs = append(closeErrs, err)
		r.logger.Error("%v", err)
	}
	if err := errors.Join(closeErrs...); err != nil {
		return err
	}

	r.logger.Info("all resource closed")
	return nil
}
