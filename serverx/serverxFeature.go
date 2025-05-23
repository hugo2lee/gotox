/*
 * @Author: hugo
 * @Date: 2024-04-19 18:02
 * @LastEditors: hugo2lee
 * @LastEditTime: 2025-04-22 21:29
 * @FilePath: \gotox\serverx\serverxFeature.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package serverx

import (
	"context"

	"github.com/hugo2lee/gotox/webx"
	"github.com/hugo2lee/gotox/webx/middleware/accesslog"
	"github.com/hugo2lee/gotox/webx/middleware/auth"
	"github.com/hugo2lee/gotox/webx/middleware/hashresponse"

	"github.com/gin-gonic/gin"
)

func (s *Serverx) EnableAccessLog() *Serverx {
	accesslog.SetLogger(s.logger)
	md := accesslog.NewBuilder(func(ctx context.Context, al accesslog.AccessLog) {
		s.logger.Info("ACCESS %v", al)
	}).AllowTrace().AllowStamp().AllowQuery().AllowReqBody().AllowRespBody().Build()
	s.Engine.Use(md)
	return s
}

func (s *Serverx) EnableAuth() *Serverx {
	auth.SetLogger(s.logger)
	aus := s.config.Auths()
	authList := make(auth.AuthPair, len(aus))
	for name, au := range aus {
		authList[auth.AUTH(au)] = auth.NAME(name)
	}
	md := auth.NewBuilder(authList).Build()
	// s.Engine.Use(md)
	s.AuthMiddle = md
	return s
}

func (s *Serverx) EnableWrapLog() *Serverx {
	webx.SetLogger(s.logger)
	return s
}

func (s *Serverx) EnableMd5Response() *Serverx {
	hashresponse.SetLogger(s.logger)
	md := hashresponse.NewBuilder().WithMd5().Build()
	// s.Engine.Use(md)
	s.HashMiddle = md
	return s
}

func (s *Serverx) EnableSha1Response() *Serverx {
	hashresponse.SetLogger(s.logger)
	md := hashresponse.NewBuilder().WithSha1().Build()
	// s.Engine.Use(md)
	s.HashMiddle = md
	return s
}

func (s *Serverx) EnableSha256Response() *Serverx {
	hashresponse.SetLogger(s.logger)
	md := hashresponse.NewBuilder().WithSha256().Build()
	// s.Engine.Use(md)
	s.HashMiddle = md
	return s
}

func (s *Serverx) LivenessCheck() *Serverx {
	s.Engine.GET("/live", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "live",
		})
	})
	return s
}

func (s *Serverx) ReadinessCheck() *Serverx {
	s.Engine.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ready",
		})
	})
	return s
}

func (s *Serverx) StarupCheck() *Serverx {
	s.Engine.GET("/startup", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "startup",
		})
	})
	return s
}
