/*
 * @Author: hugo
 * @Date: 2024-04-19 17:17
 * @LastEditors: hugo2lee
 * @LastEditTime: 2025-04-22 22:01
 * @FilePath: \gotox\serverx\serverx_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package serverx_test

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/serverx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestConfig(t *testing.T) *configx.Configx {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "dev.toml"), []byte(`
[server]
addr = "127.0.0.1:0"

[auths]
client = "client-auth"
`), 0o600))
	return configx.New(configx.WithPath(dir), configx.WithMode(configx.RUNDEV))
}

func Test_ServerEnableAccessLog(t *testing.T) {
	recorder := httptest.NewRecorder()
	svr := serverx.New(newTestConfig(t), logx.NewNoOpLogger()).EnableAccessLog()
	svr.Engine.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	svr.Engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "ok", recorder.Body.String())
}

func Test_ServerEnableAuth(t *testing.T) {
	recorder := httptest.NewRecorder()
	svr := serverx.New(newTestConfig(t), logx.NewNoOpLogger()).EnableAuth()
	svr.Engine.Use(svr.AuthMiddle)
	svr.Engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "%v", c.Keys["auth"])
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "client-auth")
	svr.Engine.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "client", recorder.Body.String())
}

func calculateHash(t *testing.T, hasher hash.Hash, data string) string {
	_, err := hasher.Write([]byte(data))
	assert.NoError(t, err, "Failed to write data to hasher")
	return hex.EncodeToString(hasher.Sum(nil))
}

func Test_ServerEnableMd5Response(t *testing.T) {
	recorder := httptest.NewRecorder()
	responseStr := "Hello, World!11122"

	svr := serverx.New(newTestConfig(t), logx.NewNoOpLogger()).EnableMd5Response()
	svr.Engine.Use(svr.HashMiddle)
	svr.Engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, responseStr)
	})
	svr.Engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, responseStr, recorder.Body.String())
	assert.Equal(t, calculateHash(t, md5.New(), responseStr), recorder.Header().Get("content-MD5"))
}
