/*
 * @Author: hugo
 * @Date: 2024-04-23 15:40
 * @LastEditors: hugo
 * @LastEditTime: 2024-10-22 15:29
 * @FilePath: \gotox\webx\middleware\accesslog\accesslog_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package accesslog_test

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/webx/middleware/accesslog"
)

func Test_AccessLog(t *testing.T) {
	// accesslog.SetLogger(logx.NewNoOpLogger())
	md := accesslog.NewBuilder(func(ctx context.Context, al accesslog.AccessLog) {
		log.Printf("ACCESS %v \n", al)
	}).AllowTrace().AllowStamp().AllowQuery().AllowReqBody().AllowRespBody().Build()

	ctx := &gin.Context{
		Request: &http.Request{
			// Header: http.Header{},
			Body: io.NopCloser(bytes.NewBufferString("hello")),
			URL: &url.URL{
				Path: "/accesslog?a=1",
			},
			Method: "GET",
		},
	}
	// ctx.Keys = make(map[string]any)
	// ctx.Keys["sn"] = "client-sn"
	// ctx.Keys["guid"] = "client-guid"

	md(ctx)
}
