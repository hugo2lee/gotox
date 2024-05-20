/*
 * @Author: hugo
 * @Date: 2024-05-17 14:04
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-17 15:15
 * @FilePath: \gotox\serverx\serverx.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package serverx

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
)

type Serverx struct {
	config     *configx.Configx
	logger     logx.Logger
	httpSrv    *http.Server
	Engine     *gin.Engine
	AuthMiddle gin.HandlerFunc
}

func New(conf *configx.Configx, log logx.Logger) *Serverx {
	engine := gin.Default()
	return &Serverx{
		config: conf,
		logger: log,
		Engine: engine,
		httpSrv: &http.Server{
			Addr:    conf.Addr(),
			Handler: engine,
		},
	}
}

// 启用http服务
func (s *Serverx) Run() error {
	return s.Engine.Run(s.config.Addr())
}

func (s *Serverx) GracefullyUp(notifyStop func()) {
	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("http server error: %s\n", err)
		notifyStop()
		return
	}
}

func (s *Serverx) GracefullyDown(notifyCtx context.Context) error {
	return s.httpSrv.Shutdown(notifyCtx)
}
