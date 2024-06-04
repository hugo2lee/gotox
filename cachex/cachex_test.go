/*
 * @Author: hugo
 * @Date: 2024-05-30 20:03
 * @LastEditors: hugo
 * @LastEditTime: 2024-06-04 10:27
 * @FilePath: \gotox\cachex\cachex_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package cachex_test

import (
	"testing"
	"time"

	"github.com/hugo2lee/gotox/cachex"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cc := cachex.New(cachex.WithExpiration(time.Minute))
	cc.Set("age", 18)
	v, ok := cc.Get("age")

	assert.True(t, ok)
	assert.Equal(t, 18, v)
}
