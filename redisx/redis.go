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

// NewWithClient wires an existing Redis client into Redisx without network I/O.
func NewWithClient(client *redis.Client, logger logx.Logger) *Redisx {
	if logger == nil {
		logger = logx.NewNoOpLogger()
	}
	return &Redisx{rds: client, logger: logger}
}

// Dial creates and verifies a Redis client from configuration.
func Dial(ctx context.Context, conf *configx.Configx, logger logx.Logger) (*Redisx, error) {
	url := conf.RedisUrl()

	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	result, err := client.Ping(ctx).Result()
	if err != nil {
		_ = client.Close()
		return nil, errors.Wrap(err, "redis connect error")
	}
	if result != "PONG" {
		_ = client.Close()
		return nil, errors.New("redis ping error")
	}

	return NewWithClient(client, logger), nil
}

// New is the compatibility constructor that dials Redis immediately.
//
// Deprecated: use Dial when connection ownership belongs here, or
// NewWithClient when the caller/composition root already owns the client.
func New(conf *configx.Configx, logger logx.Logger) (*Redisx, error) {
	return Dial(context.Background(), conf, logger)
}

func (c *Redisx) Name() string {
	return "redis"
}

func (c *Redisx) DB() *redis.Client {
	return c.rds
}

func (c *Redisx) Close(ctx context.Context) error {
	if err := c.DB().Close(); err != nil {
		c.logger.Error("redis close error %v", err)
		return err
	}
	c.logger.Info("%s close", c.Name())
	return nil
}
