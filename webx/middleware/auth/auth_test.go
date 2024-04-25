/*
 * @Author: hugo
 * @Date: 2024-04-24 14:56
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-24 14:56
 * @FilePath: \gotox\webx\middleware\auth\auth_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */

package auth_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/webx/middleware/auth"
)

func Test_Auth(t *testing.T) {
	// accesslog.SetLogger(logx.NewNoOpLogger())
	authList := map[auth.AUTH]auth.NAME{
		auth.AUTH("MTI6ZmRiNWMxMWQtYzc2OC00MzgzLTgyNjItZTY0NmFhNTE1YjU4"): auth.NAME("LS-cloud-config"),
	}
	md := auth.NewMiddlewareBuilder(authList).Build()

	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{},
			// Body: io.NopCloser(bytes.NewBufferString("hello")),
			// URL: &url.URL{
			// Path: "/accesslog",
			// },
			Method: "GET",
		},
	}

	md(ctx)
}
