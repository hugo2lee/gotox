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
	"fmt"
	"hash"
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
	"github.com/hugo2lee/gotox/webx/middleware/hashresponse"
	"github.com/stretchr/testify/assert"
)

func Test_AccessLog(t *testing.T) {
	md := accesslog.NewBuilder(func(ctx context.Context, al accesslog.AccessLog) {
		log.Printf("ACCESS %v \n", al)
	}).
		AllowTrace().
		AllowStamp().
		AllowQuery().AllowReqBody().AllowRespBody().Build()

	recorder := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/ping?name=hugo&age=18&gender=male", io.NopCloser(bytes.NewBufferString("hello")))
	req.Header.Set("Authorization", "MTI6ZmRiNWMxMWQtYzc2OC00MzgzLTgyNjItZTY0NmFhNTE1YjU4")
	req.Header.Set(accesslog.TraceIdName, "traceid-xxxx123")
	req.Header.Set(accesslog.SpanIdName, "trace-this-span-xxxx123")
	req.Header.Set(accesslog.ParentSpanIdName, "trace-parent-span-xxxx123")

	svr := gin.Default()
	svr.Use(md)
	svr.POST("/ping", func(c *gin.Context) {
		if c.Keys == nil {
			c.Keys = make(map[string]any)
		}
		c.Keys["sn"] = "client-xx-sn"
		c.Keys["guid"] = "client-xx-guid"
		time.Sleep(1 * time.Second)
		c.String(200, "pong")
	})

	svr.ServeHTTP(recorder, req)

	log.Printf("resp %v \n", recorder.Body.String())
}

func Test_Auth(t *testing.T) {
	conf := configx.New(configx.WithPath("../../conf"))
	aus := conf.Auths()
	authList := make(auth.AuthPair, len(aus))
	for name, au := range aus {
		authList[auth.AUTH(au)] = auth.NAME(name)
	}

	md := auth.NewBuilder(authList).Build()

	svr := gin.Default()
	svr.Use(md)
	svr.POST("/ping", func(c *gin.Context) {
		fmt.Printf("ACCESS client %v \n", c.Keys["auth"])
		time.Sleep(1 * time.Second)
		c.String(200, "pong")
	})

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/ping", io.NopCloser(bytes.NewBufferString("hello")))
	req.Header.Set("Authorization", "MTI6ZmRiNWMxMWQtYzc2OC00MzgzLTgyNjItZTY0NmFhNTE1YjU4")

	svr.ServeHTTP(recorder, req)

	log.Printf("resp %v \n", recorder.Body.String())
}

// calculateHash 计算哈希值
func calculateHash(t *testing.T, hasher hash.Hash, data string) string {
	_, err := hasher.Write([]byte(data))
	assert.NoError(t, err, "Failed to write data to hasher")
	return hex.EncodeToString(hasher.Sum(nil))
}

func Test_HashResponse(t *testing.T) {
	// 创建 ResponseHashBuilder 中间件
	hashMiddle := hashresponse.NewBuilder().WithMd5().WithSha1().WithSha256().Build()

	// 设置响应
	expectBody := "Hello, World!11122"

	// 模拟 Gin 服务和请求
	svr := gin.Default()
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
	// 验证响应头中的哈希值
	expectedMd5 := resp.Header.Get("Content-Md5")
	expectedSha1 := resp.Header.Get("Content-Sha1")
	expectedSha256 := resp.Header.Get("Content-Sha256")

	assert.Equal(t, expectedMd5, calculateHash(t, md5.New(), expectBody), "MD5 hash mismatch")
	assert.Equal(t, expectedSha1, calculateHash(t, sha1.New(), expectBody), "SHA1 hash mismatch")
	assert.Equal(t, expectedSha256, calculateHash(t, sha256.New(), expectBody), "SHA256 hash mismatch")
}
