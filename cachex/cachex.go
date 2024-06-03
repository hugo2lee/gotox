/*
 * @Author: hugo
 * @Date: 2024-05-30 20:22
 * @LastEditors: hugo
 * @LastEditTime: 2024-06-03 16:45
 * @FilePath: \gotox\cachex\cachex.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package cachex

import (
	"context"
	"sync"
	"time"

	"github.com/hugo2lee/gotox/resourcex"
	"github.com/patrickmn/go-cache"
)

var (
	_ resourcex.Resource = (*Cachex)(nil)
	_ Cachexer           = (*Cachex)(nil)
)

type Cachexer interface {
	Name() string
	Set(key string, value any)
	Get(key string) (any, bool)
	Close(ctx context.Context, wg *sync.WaitGroup)
}

type Cachex struct {
	*cache.Cache
}

func New(expiration time.Duration) Cachexer {
	return &Cachex{cache.New(expiration, 10*time.Minute)}
}

func (c *Cachex) Name() string {
	return "cachex"
}

func (c *Cachex) Set(key string, value any) {
	c.Cache.Set(key, value, cache.DefaultExpiration)
}

func (c *Cachex) Close(ctx context.Context, wg *sync.WaitGroup) {
	c.Cache.Flush()
	wg.Done()
}
