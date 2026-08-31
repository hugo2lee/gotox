/*
 * @Author: hugo
 * @Date: 2024-04-23 15:41
 * @LastEditors: hugo2lee
 * @LastEditTime: 2025-04-22 21:21
 * @FilePath: \gotox\webx\middleware\middleware_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package middleware_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/webx/middleware/accesslog"
	"github.com/hugo2lee/gotox/webx/middleware/auth"
	"github.com/hugo2lee/gotox/webx/middleware/hashresponse"
	"github.com/stretchr/testify/assert"
)

func Test_AccessLog(t *testing.T) {
	md := accesslog.NewBuilder(func(ctx context.Context, al accesslog.AccessLog) {}).
		AllowTrace().
		AllowStamp().
		AllowQuery().AllowReqBody().AllowRespBody().Build()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/ping?name=hugo&age=18&gender=male", io.NopCloser(bytes.NewBufferString("hello")))
	req.Header.Set("Authorization", "test-auth")
	req.Header.Set(accesslog.TraceIdName, "traceid-xxxx123")
	req.Header.Set(accesslog.SpanIdName, "trace-this-span-xxxx123")
	req.Header.Set(accesslog.ParentSpanIdName, "trace-parent-span-xxxx123")

	svr := gin.New()
	svr.Use(md)
	svr.POST("/ping", func(c *gin.Context) {
		c.Set("sn", "client-xx-sn")
		c.Set("guid", "client-xx-guid")
		c.String(http.StatusOK, "pong")
	})

	svr.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "pong", recorder.Body.String())
	assert.Equal(t, "traceid-xxxx123", recorder.Header().Get(accesslog.TraceIdName))
}

func Test_Auth(t *testing.T) {
	authList := auth.AuthPair{
		auth.AUTH("client-auth"): auth.NAME("client"),
	}
	md := auth.NewBuilder(authList).Build()

	svr := gin.New()
	svr.Use(md)
	svr.POST("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "%v", c.Keys["auth"])
	})

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/ping", io.NopCloser(bytes.NewBufferString("hello")))
	req.Header.Set("Authorization", "client-auth")
	svr.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "client", recorder.Body.String())
}

func calculateHash(t *testing.T, hasher hash.Hash, data string) string {
	_, err := hasher.Write([]byte(data))
	assert.NoError(t, err, "Failed to write data to hasher")
	return hex.EncodeToString(hasher.Sum(nil))
}

func Test_HashResponse(t *testing.T) {
	hashMiddle := hashresponse.NewBuilder().WithMd5().WithSha1().WithSha256().Build()
	expectBody := "Hello, World!11122"

	svr := gin.New()
	svr.Use(hashMiddle)
	uri := "/ping"
	svr.GET(uri, func(c *gin.Context) {
		c.String(http.StatusOK, expectBody)
	})
	recorder := httptest.NewRecorder()
	svr.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, uri, nil))

	resp := recorder.Result()
	respBody := recorder.Body.String()
	assert.Equal(t, expectBody, respBody)
	assert.Equal(t, resp.Header.Get("Content-Md5"), calculateHash(t, md5.New(), expectBody), "MD5 hash mismatch")
	assert.Equal(t, resp.Header.Get("Content-Sha1"), calculateHash(t, sha1.New(), expectBody), "SHA1 hash mismatch")
	assert.Equal(t, resp.Header.Get("Content-Sha256"), calculateHash(t, sha256.New(), expectBody), "SHA256 hash mismatch")
}
