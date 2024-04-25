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

type MiddlewareBuilder struct {
	authList AuthPair
}

func NewMiddlewareBuilder(list AuthPair) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		authList: list,
	}
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		au := c.GetHeader("Authorization")

		if val, ok := b.authList[AUTH(au)]; !ok {
			logg.Warn("Unauthorized", au)
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
