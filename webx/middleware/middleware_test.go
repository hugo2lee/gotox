/*
 * @Author: hugo
 * @Date: 2024-04-23 15:41
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-29 15:01
 * @FilePath: \gotox\webx\middleware\middleware_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package middleware_test

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/webx/middleware/accesslog"
	"github.com/hugo2lee/gotox/webx/middleware/auth"
)

func Test_AccessLog(t *testing.T) {
	md := accesslog.NewMiddlewareBuilder(func(ctx context.Context, al accesslog.AccessLog) {
		log.Printf("ACCESS %v \n", al)
	}).AllowQuery().AllowReqBody().AllowRespBody().Build()

	recorder := httptest.NewRecorder()

	svr := gin.Default()
	svr.Use(md)
	svr.POST("/ping", func(c *gin.Context) {
		time.Sleep(1 * time.Second)
		c.String(200, "pong")
	})
	svr.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/ping?name=hugo&age=18&gender=male", io.NopCloser(bytes.NewBufferString("hello"))))

	log.Printf("resp %v \n", recorder.Body.String())
}

func Test_Auth(t *testing.T) {
	conf := configx.New(configx.WithPath("../../conf"))
	aus := conf.Auths()
	authList := make(auth.AuthPair, len(aus))
	for name, au := range aus {
		authList[auth.AUTH(au)] = auth.NAME(name)
	}

	md := auth.NewMiddlewareBuilder(authList).Build()

	svr := gin.Default()
	svr.Use(md)
	svr.POST("/ping", func(c *gin.Context) {
		log.Printf("ACCESS client %v \n", c.Keys["auth"])
		time.Sleep(1 * time.Second)
		c.String(200, "pong")
	})

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/ping", io.NopCloser(bytes.NewBufferString("hello")))
	req.Header.Set("Authorization", "MTI6ZmRiNWMxMWQtYzc2OC00MzgzLTgyNjItZTY0NmFhNTE1YjU4")

	svr.ServeHTTP(recorder, req)

	log.Printf("resp %v \n", recorder.Body.String())
}
