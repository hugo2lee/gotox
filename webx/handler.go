/*
 * @Author: hugo
 * @Date: 2024-04-19 17:54
 * @LastEditors: hugo
 * @LastEditTime: 2024-07-31 17:50
 * @FilePath: \gotox\webx\handler.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package webx

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/webx/middleware/accesslog"
)

// 受制于泛型，这里只能使用包变量，如无任何实例赋予就用这个
var logg logx.Logger = logx.Log

// // 自定义的logger，建议实例化赋予
func SetLogger(l logx.Logger) {
	logg = l
}

type Handler interface {
	// PublicAPI(server *gin.Engine)
	// PrivateAPI(server *gin.Engine)
	RegisterRouter(*gin.Engine)
}

func Wrap(fn func(ctx *gin.Context) (Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := fn(ctx)
		if err != nil {
			// 打印日志
			logg.Error("traceid: %s, biz error: %v", traceIDFromContext(ctx), err)
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapBind[T any](fn func(ctx *gin.Context, req T) (Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var t T
		if err := ctx.Bind(&t); err != nil {
			// 打印日志
			logg.Error("Bind Error %v", err)
			ctx.JSON(http.StatusBadRequest, Response{Message: err.Error()})
			return
		}
		res, err := fn(ctx, t)
		if err != nil {
			// 打印日志
			logg.Error("traceid: %s, biz error: %v", traceIDFromContext(ctx), err)
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapPage(fn func(ctx *gin.Context, page, pageSize int) (Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page, _ := strconv.Atoi(ctx.Query("page"))
		pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))
		res, err := fn(ctx, page, pageSize)
		if err != nil {
			// 打印日志
			logg.Error("traceid: %s, biz error: %v", traceIDFromContext(ctx), err)
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapBindQueryAndBody[Q any, B any](fn func(ctx *gin.Context, query Q, body B) (Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var q Q
		var b B
		if err := ctx.BindQuery(&q); err != nil {
			// 打印日志
			logg.Error("query Bind Error %v", err)
			ctx.JSON(http.StatusBadRequest, Response{Message: err.Error()})
			return
		}
		if err := ctx.BindJSON(&b); err != nil {
			// 打印日志
			logg.Error("body Bind Error %v", err)
			ctx.JSON(http.StatusBadRequest, Response{Message: err.Error()})
			return
		}
		res, err := fn(ctx, q, b)
		if err != nil {
			// 打印日志
			logg.Error("traceid: %s, biz error: %v", traceIDFromContext(ctx), err)
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func traceIDFromContext(ctx *gin.Context) string {
	if val, ok := ctx.Keys[accesslog.GinKeyTraceName]; ok {
		if traceID, ok := val.(string); ok {
			return traceID
		}
	}
	return "not found"
}
