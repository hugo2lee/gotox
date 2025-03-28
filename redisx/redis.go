/*
 * @Author: hugo
 * @Date: 2024-04-02 15:09
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-17 15:13
 * @FilePath: \gotox\redisx\redis.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */

package redisx

import (
	"context"
	"sync"

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/resourcex"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var _ resourcex.Resource = (*Redisx)(nil)

type Redisx struct {
	rds    *redis.Client
	logger logx.Logger
}

func New(conf *configx.Configx, logCli logx.Logger) (*Redisx, error) {
	url := conf.RedisUrl()

	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opt)

	result, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrap(err, "redis connect error")
	}
	if result != "PONG" {
		return nil, errors.New("redis ping error")
	}

	return &Redisx{rdb, logCli}, nil
}

func (c *Redisx) Name() string {
	return "redis"
}

func (c *Redisx) DB() *redis.Client {
	return c.rds
}

func (c *Redisx) Close(ctx context.Context, wg *sync.WaitGroup) {
	if err := c.DB().Close(); err != nil {
		c.logger.Error("redis close error %v", err)
		return
	}
	wg.Done()
	c.logger.Info("%s close", c.Name())
}
