//go:build integration

package resourcex_test

import (
	"context"
	"testing"
	"time"

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/redisx"
	"github.com/hugo2lee/gotox/resourcex"
	"github.com/stretchr/testify/require"
)

func TestRedisResourceIntegration(t *testing.T) {
	conf := configx.New(configx.WithPath("../conf"))
	rds, err := redisx.Dial(context.Background(), conf, logx.New(conf))
	require.NoError(t, err)

	group := resourcex.NewResourcexGroup()
	group.AddResource(rds)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	require.NoError(t, group.CloseAll(ctx))
}
