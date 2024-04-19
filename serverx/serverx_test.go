/*
 * @Author: hugo
 * @Date: 2024-04-19 17:17
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-19 17:49
 * @FilePath: \gotox\serverx\serverx_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package serverx_test

import (
	"context"
	"testing"
	"time"

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/mongox"
	"github.com/hugo2lee/gotox/ormx"
	"github.com/hugo2lee/gotox/redisx"
	"github.com/hugo2lee/gotox/serverx"
	"github.com/stretchr/testify/assert"
)

func Test_IocResource(t *testing.T) {
	conf := configx.New(configx.WithPath("../conf"))
	log := logx.New(conf)

	db, err := ormx.New(conf, log)
	assert.NoError(t, err)

	rds, err := redisx.New(conf, log)
	assert.NoError(t, err)

	mongo, err := mongox.New(conf, log)
	assert.NoError(t, err)

	svr := serverx.New(conf, log)
	svr.AddResource(db, rds, mongo)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	svr.CloseResource(ctx)
}

func Test_ServerUp(t *testing.T) {
	conf := configx.New(configx.WithPath("../conf"))
	log := logx.New(conf)

	db, err := ormx.New(conf, log)
	assert.NoError(t, err)

	rds, err := redisx.New(conf, log)
	assert.NoError(t, err)

	mongo, err := mongox.New(conf, log)
	assert.NoError(t, err)

	svr := serverx.New(conf, log)
	svr.AddResource(db, rds, mongo)

	svr.GracefullyUp()
}
