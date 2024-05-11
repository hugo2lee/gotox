/*
 * @Author: hugo
 * @Date: 2024-04-02 15:09
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-11 15:03
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

var _ resourcex.Resource = (*redisCli)(nil)

type redisCli struct {
	rds    *redis.Client
	logger logx.Logger
}

func New(conf *configx.ConfigCli, logCli logx.Logger) (*redisCli, error) {
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

	return &redisCli{rdb, logCli}, nil
}

func (c *redisCli) Name() string {
	return "redis"
}

func (c *redisCli) DB() *redis.Client {
	return c.rds
}

func (c *redisCli) Close(ctx context.Context, wg *sync.WaitGroup) {
	if err := c.DB().Close(); err != nil {
		c.logger.Error("redis close error %v", err)
		return
	}
	wg.Done()
	c.logger.Info("%s close", c.Name())
}
