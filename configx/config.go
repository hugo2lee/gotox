/*
 * @Author: hugo
 * @Date: 2024-03-12 15:01
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-19 16:43
 * @FilePath: \gotox\configx\config.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package configx

import (
	"log"

	"github.com/spf13/viper"
)

const (
	RUNMODESTR = "RUNMODE"
	// RUNMODEKEY = "env.mode"

	RUNDEV  = "dev"
	RUNTEST = "test"
	RUNPROD = "prod"

	DEFAULTPATH       = "../conf"
	DEFAULTCONFIGTYPE = "toml"
	DEFAULTMODE       = RUNDEV
)

type ConfigCli struct {
	mode  string
	path  string
	viper *viper.Viper
}

type option func(*ConfigCli)

func WithMode(mode string) option {
	return func(cli *ConfigCli) {
		if mode != RUNDEV && mode != RUNPROD && mode != RUNTEST {
			log.Fatalf("invalid mode: %s! only support: %s, %s, %s", mode, RUNDEV, RUNPROD, RUNTEST)
		}
		cli.mode = mode
	}
}

func WithPath(path string) option {
	return func(cli *ConfigCli) {
		if path == "" {
			log.Fatalf("invalid path: %s", path)
		}
		cli.path = path
	}
}

func New(options ...option) *ConfigCli {
	// 初始化配置过程中可以直接panic

	cli := &ConfigCli{}
	for _, opt := range options {
		opt(cli)
	}

	v := viper.New()
	v.SetConfigType(DEFAULTCONFIGTYPE)

	// 先设置默认值
	v.SetDefault(RUNMODESTR, DEFAULTMODE)
	// v.SetEnvPrefix("go")          // 设置环境变量的前缀

	// 再绑定环境变量
	err := v.BindEnv(RUNMODESTR)
	if err != nil {
		log.Fatalf("config BindEnv error: %s \n", err)
	}

	// 手动指定的优先级最高
	if cli.mode != "" {
		v.Set(RUNMODESTR, cli.mode)
		// 配置的文件名
		v.SetConfigName(cli.mode)
	} else {
		v.SetConfigName(v.GetString(RUNMODESTR))
		cli.mode = v.GetString(RUNMODESTR)
	}

	if cli.path != "" {
		v.AddConfigPath(cli.path)
	} else {
		v.AddConfigPath(DEFAULTPATH)
		cli.path = DEFAULTPATH
	}

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	log.Printf("Using config mode: %s, file: %s \n", v.GetString(RUNMODESTR), v.ConfigFileUsed())

	return &ConfigCli{
		viper: v,
	}
}

func (c *ConfigCli) Mode() string {
	return c.viper.GetString(RUNMODESTR)
}
