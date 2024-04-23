/*
 * @Author: hugo
 * @Date: 2024-04-23 15:41
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-23 18:58
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
	"github.com/hugo2lee/gotox/webx/middleware/accesslog"
)

func Test_AccessLog(t *testing.T) {
	md := accesslog.NewMiddlewareBuilder(func(ctx context.Context, al accesslog.AccessLog) {
		log.Printf("ACCESS %v \n", al)
	}).AllowReqBody().AllowRespBody().Build()

	recorder := httptest.NewRecorder()

	svr := gin.Default()
	svr.Use(md)
	svr.POST("/ping", func(c *gin.Context) {
		time.Sleep(1 * time.Second)
		c.String(200, "pong")
	})
	svr.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/ping", io.NopCloser(bytes.NewBufferString("hello"))))

	log.Printf("resp %v \n", recorder.Body.String())
}
