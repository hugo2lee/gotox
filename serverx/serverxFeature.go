/*
 * @Author: hugo
 * @Date: 2024-04-19 18:02
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-19 22:01
 * @FilePath: \gotox\serverx\serverxFeature.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package serverx

import (
	"context"

	"github.com/hugo2lee/gotox/webx/middleware/accesslog"

	"github.com/gin-gonic/gin"
)

func (s *Server) EnableAccessLog() *Server {
	accesslog.SetLogger(s.logger)
	md := accesslog.NewMiddlewareBuilder(func(ctx context.Context, al accesslog.AccessLog) {
		s.logger.Info("ACCESS %v", al)
	}).AllowReqBody().AllowRespBody().Build()
	s.Engine.Use(md)
	return s
}

func (s *Server) LivenessCheck() *Server {
	s.Engine.GET("/live", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "live",
		})
	})
	return s
}

func (s *Server) ReadinessCheck() *Server {
	s.Engine.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ready",
		})
	})
	return s
}

func (s *Server) StarupCheck() *Server {
	s.Engine.GET("/startup", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "startup",
		})
	})
	return s
}
