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

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/webx/middleware/accesslog"
	"github.com/hugo2lee/gotox/webx/middleware/auth"
	"github.com/hugo2lee/gotox/webx/middleware/hashresponse"
)

func (s *Serverx) accessLogMiddleware(verbose bool) gin.HandlerFunc {
	builder := accesslog.NewBuilder(func(ctx context.Context, al accesslog.AccessLog) {
		s.logger.Info("ACCESS %v", al)
	}).WithLogger(s.logger).AllowTrace().AllowStamp()

	if verbose {
		builder.AllowQuery().AllowReqBody().AllowRespBody()
	}
	return builder.Build()
}

// EnableAccessLog enables metadata-only access logging. Authorization/token
// values are redacted by the accesslog middleware and request query/body or
// response body are not persisted by default.
func (s *Serverx) EnableAccessLog() *Serverx {
	s.Engine.Use(s.accessLogMiddleware(false))
	return s
}

// EnableVerboseAccessLog opts into query, request-body, and response-body
// logging. Credential headers remain redacted unless callers build the
// accesslog middleware themselves and explicitly opt into sensitive values.
func (s *Serverx) EnableVerboseAccessLog() *Serverx {
	s.Engine.Use(s.accessLogMiddleware(true))
	return s
}

func (s *Serverx) EnableAuth() *Serverx {
	aus := s.config.Auths()
	authList := make(auth.AuthPair, len(aus))
	for name, au := range aus {
		authList[auth.AUTH(au)] = auth.NAME(name)
	}
	md := auth.NewBuilder(authList).WithLogger(s.logger).Build()
	s.AuthMiddle = md
	return s
}

// EnableWrapLog is kept for source compatibility. webx.Wrap* now uses stdout
// fallback by default; callers that need a custom logger should use the
// corresponding webx.Wrap*WithLogger function.
func (s *Serverx) EnableWrapLog() *Serverx {
	return s
}

func (s *Serverx) EnableMd5Response() *Serverx {
	md := hashresponse.NewBuilder().WithLogger(s.logger).WithMd5().Build()
	s.HashMiddle = md
	return s
}

func (s *Serverx) EnableSha1Response() *Serverx {
	md := hashresponse.NewBuilder().WithLogger(s.logger).WithSha1().Build()
	s.HashMiddle = md
	return s
}

func (s *Serverx) EnableSha256Response() *Serverx {
	md := hashresponse.NewBuilder().WithLogger(s.logger).WithSha256().Build()
	s.HashMiddle = md
	return s
}

func (s *Serverx) LivenessCheck() *Serverx {
	s.Engine.GET("/live", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "live"})
	})
	return s
}

func (s *Serverx) ReadinessCheck() *Serverx {
	s.Engine.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ready"})
	})
	return s
}

func (s *Serverx) StarupCheck() *Serverx {
	s.Engine.GET("/startup", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "startup"})
	})
	return s
}
