/*
 * @Author: hugo
 * @Date: 2024-04-24 14:56
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-24 14:56
 * @FilePath: \gotox\webx\middleware\auth\auth_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */

package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/webx/middleware/auth"
	"github.com/stretchr/testify/assert"
)

func Test_Auth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authValue := "MTI6ZmRiNWMxMWQtYzc2OC00MzgzLTgyNjItZTY0NmFhNTE1YjU4"
	authList := auth.AuthPair{
		auth.AUTH(authValue): auth.NAME("LS-cloud-config"),
	}

	engine := gin.New()
	engine.Use(auth.NewBuilder(authList).Build())
	engine.GET("/", func(c *gin.Context) {
		name, exists := c.Get("auth")
		assert.True(t, exists)
		assert.Equal(t, auth.NAME("LS-cloud-config"), name)
		c.String(http.StatusOK, "ok")
	})

	t.Run("authorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", authValue)
		recorder := httptest.NewRecorder()

		engine.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "ok", recorder.Body.String())
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()

		engine.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})
}
