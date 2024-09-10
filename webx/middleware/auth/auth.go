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

// 受制于泛型，这里只能使用包变量，如无任何实例赋予就用这个
var logg logx.Logger = logx.Log

// // 自定义的logger，建议实例化赋予
func SetLogger(l logx.Logger) {
	logg = l
}

type (
	AUTH string
	NAME string
)

type AuthPair map[AUTH]NAME

type Auth struct {
	authList AuthPair
}

func NewBuilder(list AuthPair) *Auth {
	return &Auth{
		authList: list,
	}
}

func (b *Auth) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		au := c.GetHeader("Authorization")

		if val, ok := b.authList[AUTH(au)]; !ok {
			logg.Warn("Unauthorized %v", au)
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

		// 这里会执行到业务代码
		c.Next()
	}
}

func NoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
