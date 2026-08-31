/*
 * @Author: hugo
 * @Date: 2024-04-28 16:51
 * @LastEditors: hugo
 * @LastEditTime: 2024-09-10 16:02
 * @FilePath: \gotox\webx\middleware\auth\auth.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/logx"
)

type (
	AUTH string
	NAME string
)

type AuthPair map[AUTH]NAME

type Auth struct {
	authList AuthPair
	logger   logx.Logger
}

func NewBuilder(list AuthPair) *Auth {
	return &Auth{
		authList: list,
		logger:   logx.NewNoOpLogger(),
	}
}

func (b *Auth) WithLogger(logger logx.Logger) *Auth {
	if logger != nil {
		b.logger = logger
	}
	return b
}

func (b *Auth) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		au := c.GetHeader("Authorization")

		if val, ok := b.authList[AUTH(au)]; !ok {
			b.logger.Warn("Unauthorized %v", au)
			c.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "Unauthorized",
			})
			return
		} else {
			if c.Keys == nil {
				c.Keys = make(map[string]any)
			}
			c.Keys["auth"] = val
		}

		c.Next()
	}
}

func NoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
