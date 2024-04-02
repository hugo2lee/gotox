package redisx_test

import (
	"context"
	"log"
	"testing"

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/redisx"
	"github.com/stretchr/testify/assert"
)

func TestRedis(t *testing.T) {
	t.Parallel()
	conf := configx.New(configx.WithPath("../configx/conf"))
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
