package redisx_test

import (
	"context"
	"testing"

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
