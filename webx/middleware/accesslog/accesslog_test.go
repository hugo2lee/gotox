/*
 * @Author: hugo
 * @Date: 2024-04-23 15:40
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-23 19:16
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
	md := accesslog.NewMiddlewareBuilder(func(ctx context.Context, al accesslog.AccessLog) {
		log.Printf("ACCESS %v \n", al)
	}).AllowReqBody().AllowRespBody().Build()

	ctx := &gin.Context{
		Request: &http.Request{
			// Header: http.Header{},
			Body: io.NopCloser(bytes.NewBufferString("hello")),
			URL: &url.URL{
				Path: "/accesslog",
			},
			Method: "GET",
		},
	}

	md(ctx)
}
