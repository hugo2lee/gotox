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
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/webx/middleware/accesslog"
	"github.com/stretchr/testify/assert"
)

func Test_AccessLog(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var got accesslog.AccessLog
	md := accesslog.NewBuilder(func(_ context.Context, al accesslog.AccessLog) {
		got = al
	}).AllowTrace().AllowStamp().AllowQuery().AllowReqBody().AllowRespBody().Build()

	engine := gin.New()
	engine.Use(md)
	engine.POST("/accesslog", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	req := httptest.NewRequest(http.MethodPost, "/accesslog?a=1", strings.NewReader("hello"))
	req.Header.Set(accesslog.TraceIdName, "trace-id")
	req.Header.Set(accesslog.Auth, "Bearer secret-auth-token")
	req.Header.Set(accesslog.Token, "secret-token")
	recorder := httptest.NewRecorder()

	engine.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "pong", recorder.Body.String())
	assert.Equal(t, "trace-id", got.TraceId)
	assert.Equal(t, "Bearer [REDACTED]", got.Auth)
	assert.Equal(t, "[REDACTED]", got.Token)
	assert.Equal(t, "a=1", got.Query)
	assert.Equal(t, "hello", got.ReqBody)
	assert.Equal(t, "pong", got.RespBody)
}

func Test_AccessLogSensitiveCredentialsRequireExplicitOptIn(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var got accesslog.AccessLog
	md := accesslog.NewBuilder(func(_ context.Context, al accesslog.AccessLog) {
		got = al
	}).AllowTrace().AllowSensitiveCredentials().Build()

	engine := gin.New()
	engine.Use(md)
	engine.GET("/accesslog", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/accesslog", nil)
	req.Header.Set(accesslog.Auth, "Bearer secret-auth-token")
	req.Header.Set(accesslog.Token, "secret-token")
	recorder := httptest.NewRecorder()

	engine.ServeHTTP(recorder, req)

	assert.Equal(t, "Bearer secret-auth-token", got.Auth)
	assert.Equal(t, "secret-token", got.Token)
}
