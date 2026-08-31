/*
 * @Author: hugo
 * @Date: 2024-05-11 14:00
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-11 15:11
 * @FilePath: \gotox\resourcex\resourcex_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package resourcex_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/redisx"
	"github.com/hugo2lee/gotox/resourcex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubResource struct {
	name    string
	closeFn func(context.Context) error
}

func (s stubResource) Name() string { return s.name }

func (s stubResource) Close(ctx context.Context) error {
	return s.closeFn(ctx)
}

func TestResourceGroupClosesResources(t *testing.T) {
	closed := make(chan struct{})
	group := resourcex.NewResourcexGroup()
	group.AddResource(resourcex.NewResourcex("memory", func(context.Context) {
		close(closed)
	}))

	err := group.CloseAll(context.Background())
	require.NoError(t, err)

	select {
	case <-closed:
	default:
		t.Fatal("resource was not closed")
	}
}

func TestResourceGroupReturnsCloseError(t *testing.T) {
	expected := errors.New("close failed")
	group := resourcex.NewResourcexGroup()
	group.AddResource(stubResource{
		name: "failing",
		closeFn: func(context.Context) error {
			return expected
		},
	})

	err := group.CloseAll(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, expected)
}

func TestResourceGroupReturnsContextErrorOnTimeout(t *testing.T) {
	unblock := make(chan struct{})
	defer close(unblock)

	group := resourcex.NewResourcexGroup()
	group.AddResource(stubResource{
		name: "blocked",
		closeFn: func(context.Context) error {
			<-unblock
			return nil
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := group.CloseAll(ctx)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestRedisResource(t *testing.T) {
	t.Log("integration test: requires ../conf and a real redis service")

	conf := configx.New(configx.WithPath("../conf"))
	rds, err := redisx.New(conf, logx.New(conf))
	assert.NoError(t, err)

	group := resourcex.NewResourcexGroup()
	group.AddResource(rds)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	assert.NoError(t, group.CloseAll(ctx))
}
