/*
 * @Author: hugo
 * @Date: 2024-05-11 15:05
 * @LastEditors: hugo
 * @LastEditTime: 2024-06-03 16:53
 * @FilePath: \gotox\appx\appx_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */

package appx_test

import (
	"testing"

	"github.com/hugo2lee/gotox/appx"
	"github.com/hugo2lee/gotox/configx"
	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	a := appx.New(configx.WithPath("../conf")).EnableCache()

	a.Cachex.Set("age", 18)
	v, ok := a.Cachex.Get("age")
	assert.True(t, ok)
	assert.Equal(t, 18, v)

	// a.Run()
	t.Log(a)
}
