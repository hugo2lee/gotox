/*
 * @Author: hugo
 * @Date: 2024-04-19 17:17
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-11 16:21
 * @FilePath: \gotox\serverx\serverx_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package serverx_test

import (
	"context"
	"fmt"
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
	svr.Engine.GET("/", func(c *gin.Context) {
		log.Printf("client %v \n", c.Keys["auth"])
		time.Sleep(1 * time.Second)
		c.String(200, "pong")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "MTI6ZmRiNWMxMWQtYzc2OC00MzgzLTgyNjItZTY0NmFhNTE1YjU4")
	svr.Engine.ServeHTTP(recorder, req)

	log.Printf("resp %v \n", recorder.Body.String())
}
