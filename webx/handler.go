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

type Handler interface {
	RegisterRouter(*gin.Engine)
}

func fallbackLogger(logger logx.Logger) logx.Logger {
	if logger != nil {
		return logger
	}
	return logx.NewNoOpLogger()
}

func Wrap(fn func(ctx *gin.Context) (Response, error)) gin.HandlerFunc {
	return WrapWithLogger(nil, fn)
}

func WrapWithLogger(logger logx.Logger, fn func(ctx *gin.Context) (Response, error)) gin.HandlerFunc {
	logger = fallbackLogger(logger)
	return func(ctx *gin.Context) {
		res, err := fn(ctx)
		if err != nil {
			logger.Error("traceid: %s, biz error: %v", traceIDFromContext(ctx), err)
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapBind[T any](fn func(ctx *gin.Context, req T) (Response, error)) gin.HandlerFunc {
	return WrapBindWithLogger(nil, fn)
}

func WrapBindWithLogger[T any](logger logx.Logger, fn func(ctx *gin.Context, req T) (Response, error)) gin.HandlerFunc {
	logger = fallbackLogger(logger)
	return func(ctx *gin.Context) {
		var t T
		if err := ctx.Bind(&t); err != nil {
			logger.Error("Bind Error %v", err)
			ctx.JSON(http.StatusBadRequest, Response{Message: err.Error()})
			return
		}
		res, err := fn(ctx, t)
		if err != nil {
			logger.Error("traceid: %s, biz error: %v", traceIDFromContext(ctx), err)
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapPage(fn func(ctx *gin.Context, page, pageSize int) (Response, error)) gin.HandlerFunc {
	return WrapPageWithLogger(nil, fn)
}

func WrapPageWithLogger(logger logx.Logger, fn func(ctx *gin.Context, page, pageSize int) (Response, error)) gin.HandlerFunc {
	logger = fallbackLogger(logger)
	return func(ctx *gin.Context) {
		page, _ := strconv.Atoi(ctx.Query("page"))
		pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))
		res, err := fn(ctx, page, pageSize)
		if err != nil {
			logger.Error("traceid: %s, biz error: %v", traceIDFromContext(ctx), err)
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapBindQueryAndBody[Q any, B any](fn func(ctx *gin.Context, query Q, body B) (Response, error)) gin.HandlerFunc {
	return WrapBindQueryAndBodyWithLogger(nil, fn)
}

func WrapBindQueryAndBodyWithLogger[Q any, B any](logger logx.Logger, fn func(ctx *gin.Context, query Q, body B) (Response, error)) gin.HandlerFunc {
	logger = fallbackLogger(logger)
	return func(ctx *gin.Context) {
		var q Q
		var b B
		if err := ctx.BindQuery(&q); err != nil {
			logger.Error("query Bind Error %v", err)
			ctx.JSON(http.StatusBadRequest, Response{Message: err.Error()})
			return
		}
		if err := ctx.BindJSON(&b); err != nil {
			logger.Error("body Bind Error %v", err)
			ctx.JSON(http.StatusBadRequest, Response{Message: err.Error()})
			return
		}
		res, err := fn(ctx, q, b)
		if err != nil {
			logger.Error("traceid: %s, biz error: %v", traceIDFromContext(ctx), err)
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
