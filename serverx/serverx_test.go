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
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/serverx"
	"github.com/stretchr/testify/assert"
)

func Test_ServerUp(t *testing.T) {
	conf := configx.New(configx.WithPath("../conf"))
	log := logx.New(conf)
	svr := serverx.New(conf, log)

	_, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT, syscall.SIGTERM)
	defer cancel()

	go svr.GracefullyUp(cancel)

	time.Sleep(3 * time.Second)
	cancel()

	timeOut, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := svr.GracefullyDown(timeOut)
	assert.NoError(t, err)

	time.Sleep(5 * time.Second)

	fmt.Println("done")
}

func Test_ServerEnableAccessLog(t *testing.T) {
	conf := configx.New(configx.WithPath("../conf"))
	logger := logx.New(conf)

	recorder := httptest.NewRecorder()

	svr := serverx.New(conf, logger).EnableAccessLog()
	svr.Engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	log.Printf("resp %v \n", recorder.Body.String())
}

func Test_ServerEnableAuth(t *testing.T) {
	conf := configx.New(configx.WithPath("../conf"))
	logger := logx.New(conf)

	recorder := httptest.NewRecorder()

	svr := serverx.New(conf, logger).EnableAccessLog().EnableAuth()
	svr.Engine.Use(svr.AuthMiddle)
	svr.Engine.GET("/", func(c *gin.Context) {
		log.Printf("client %v \n", c.Keys["auth"])
		time.Sleep(1 * time.Second)
		c.String(200, "pong")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "client-auth")
	svr.Engine.ServeHTTP(recorder, req)

	log.Printf("resp %v \n", recorder.Body.String())
}

// calculateHash 计算哈希值
func calculateHash(t *testing.T, hasher hash.Hash, data string) string {
	_, err := hasher.Write([]byte(data))
	assert.NoError(t, err, "Failed to write data to hasher")
	// 计算哈希值并转换为十六进制字符串
	// 注意：这里的 Sum(nil) 会返回一个新的切片
	return hex.EncodeToString(hasher.Sum(nil))
}

func Test_ServerEnableMd5Response(t *testing.T) {
	conf := configx.New(configx.WithPath("../conf"))
	logger := logx.New(conf)

	recorder := httptest.NewRecorder()
	responseStr := "Hello, World!11122"

	svr := serverx.New(conf, logger).EnableAccessLog().EnableMd5Response()
	svr.Engine.Use(svr.HashMiddle)
	svr.Engine.GET("/", func(c *gin.Context) {
		c.String(200, responseStr)
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	svr.Engine.ServeHTTP(recorder, req)

	log.Printf("resp %v \n", recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, calculateHash(t, md5.New(), responseStr), recorder.Header().Get("content-MD5"))
}
