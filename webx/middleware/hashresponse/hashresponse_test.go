/*
 * @Author: hugo2lee
 * @Date: 2025-04-21 15:56
 * @LastEditors: hugo2lee
 * @LastEditTime: 2025-04-22 21:16
 * @FilePath: \gotox\webx\middleware\hashresponse\hashresponse_test.go
 * @Description:
 *
 * Copyright (c) 2025 by hugo, All Rights Reserved.
 */

package hashresponse_test

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/webx/middleware/hashresponse"
	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	// 创建 ResponseHashBuilder 中间件
	hashMiddle := hashresponse.NewBuilder().WithMd5().WithSha1().WithSha256().Build()

	recorder := httptest.NewRecorder()
	// 设置响应
	expectBody := "Hello, World!223355"

	// 模拟 Gin 上下文和请求
	gin.SetMode(gin.TestMode)
	c, engine := gin.CreateTestContext(recorder)
	engine.Use(hashMiddle)
	pathh := "/"
	engine.GET(pathh, func(c *gin.Context) {
		c.String(http.StatusOK, expectBody)
	})
	c.Request = httptest.NewRequest(http.MethodGet, pathh, nil)

	// 处理请求
	engine.HandleContext(c)

	resp := recorder.Result()
	respBody := recorder.Body.String()
	assert.Equal(t, expectBody, respBody)
	// 验证响应头中的哈希值
	expectedMd5 := resp.Header.Get("Content-Md5")
	expectedSha1 := resp.Header.Get("Content-Sha1")
	expectedSha256 := resp.Header.Get("Content-Sha256")

	assert.Equal(t, expectedMd5, calculateHash(t, md5.New(), respBody), "MD5 hash mismatch")
	assert.Equal(t, expectedSha1, calculateHash(t, sha1.New(), respBody), "SHA1 hash mismatch")
	assert.Equal(t, expectedSha256, calculateHash(t, sha256.New(), respBody), "SHA256 hash mismatch")
}

// calculateHash 计算哈希值
func calculateHash(t *testing.T, hasher hash.Hash, data string) string {
	_, err := hasher.Write([]byte(data))
	assert.NoError(t, err, "Failed to write data to hasher")
	// 计算哈希值并转换为十六进制字符串
	// 注意：这里的 Sum(nil) 会返回一个新的切片
	return hex.EncodeToString(hasher.Sum(nil))
}
