/*
 * @Author: hugo
 * @Date: 2024-04-17 17:16
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-19 16:51
 * @FilePath: \gotox\redisx\redis_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package redisx_test

import (
	"context"
	"log"
	"testing"

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/redisx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisClientWiring(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"})
	wrapped := redisx.NewWithClient(client, nil)

	assert.Same(t, client, wrapped.DB())
	require.NoError(t, wrapped.Close(context.Background()))
}

func TestRedis(t *testing.T) {
	t.Parallel()
	conf := configx.New(configx.WithPath("../conf"))
	logger := logx.New(conf)
	cli, err := redisx.New(conf, logger)

	assert.NoError(t, err)
	cmd := cli.DB()
	r, err := cmd.Ping(context.Background()).Result()
	assert.NoError(t, err)
	assert.Equal(t, "PONG", r)
	log.Println(r, err)

	r, err = cmd.Set(context.Background(), "name", "hugo", 0).Result()
	assert.NoError(t, err)
	assert.Equal(t, "OK", r)
}
