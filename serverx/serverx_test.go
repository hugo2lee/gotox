/*
 * @Author: hugo
 * @Date: 2024-04-19 17:17
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-25 19:28
 * @FilePath: \gotox\serverx\serverx_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package serverx_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

func Test_ServerEnableAccessLog(t *testing.T) {
	conf := configx.New(configx.WithPath("../conf"))
	logger := logx.New(conf)

	recorder := httptest.NewRecorder()

	svr := serverx.New(conf, logger).EnableAccessLog()
	svr.Engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	log.Printf("resp %v \n", recorder.Body.String())
}

func Test_ServerEnableAuth(t *testing.T) {
	conf := configx.New(configx.WithPath("../conf"))
	logger := logx.New(conf)

	recorder := httptest.NewRecorder()

	svr := serverx.New(conf, logger).EnableAccessLog().EnableAuth()
	svr.Engine.GET("/", func(c *gin.Context) {
		log.Printf("client %v \n", c.Keys["auth"])
		time.Sleep(1 * time.Second)
		c.String(200, "pong")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "MTI6ZmRiNWMxMWQtYzc2OC00MzgzLTgyNjItZTY0NmFhNTE1YjU4")
	svr.Engine.ServeHTTP(recorder, req)

	log.Printf("resp %v \n", recorder.Body.String())
}
