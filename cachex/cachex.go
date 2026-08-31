/*
 * @Author: hugo
 * @Date: 2024-05-30 20:22
 * @LastEditors: hugo
 * @LastEditTime: 2024-06-12 15:07
 * @FilePath: \gotox\cachex\cachex.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package cachex

import (
	"context"
	"time"

	"github.com/hugo2lee/gotox/resourcex"
	"github.com/patrickmn/go-cache"
)

var (
	_ resourcex.Resource = (*Cachex)(nil)
	_ Cachexer           = (*Cachex)(nil)
)

// Cachexer is kept for source compatibility with existing users.
//
// Deprecated: consumers should define the smallest cache interface they need
// at the consuming boundary. New returns *Cachex directly.
type Cachexer interface {
	Name() string
	Set(key string, value any)
	Get(key string) (any, bool)
	Delete(key string)
	Flush()
	Close(ctx context.Context) error
}

type Cachex struct {
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	*cache.Cache
}

type Option func(*Cachex)

func WithExpiration(expiration time.Duration) Option {
	return func(c *Cachex) {
		c.defaultExpiration = expiration
	}
}

func WithCleanupInterval(cleanupInterval time.Duration) Option {
	return func(c *Cachex) {
		c.cleanupInterval = cleanupInterval
	}
}

func New(opts ...Option) *Cachex {
	ca := &Cachex{}
	for _, opt := range opts {
		opt(ca)
	}
	ca.Cache = cache.New(ca.defaultExpiration, ca.cleanupInterval)
	return ca
}

func (c *Cachex) Name() string {
	return "cachex"
}

func (c *Cachex) Set(key string, value any) {
	c.Cache.Set(key, value, cache.DefaultExpiration)
}

func (c *Cachex) Close(ctx context.Context) error {
	c.Cache.Flush()
	return nil
}
