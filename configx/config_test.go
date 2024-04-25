/*
 * @Author: hugo
 * @Date: 2024-03-12 15:28
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-25 16:55
 * @FilePath: \gotox\configx\config_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package configx_test

import (
	"log"
	"testing"

	"github.com/hugo2lee/gotox/configx"
)

type custom struct {
	*configx.ConfigCli
}

func (c custom) LogDir() string {
	return "custom log dir"
}

func TestConfigExample(t *testing.T) {
	c := configx.New()
	log.Println(c.LogDir())

	cu := custom{c}
	log.Println(cu.LogDir())

	aus := cu.Auths()
	log.Println(aus)
}
