/*
 * @Author: hugo
 * @Date: 2024-05-11 15:05
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-11 16:54
 * @FilePath: \gotox\appx\appx_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */

package appx_test

import (
	"testing"

	"github.com/hugo2lee/gotox/appx"
)

func TestNewApp(t *testing.T) {
	appx.New().Run()
}
