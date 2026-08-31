/*
 * @Author: hugo
 * @Date: 2024-03-12 15:01
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-17 14:04
 * @FilePath: \gotox\configx\config.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package configx

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const (
	RUNMODESTR = "RUNMODE"
	// RUNMODEKEY = "env.mode"

	RUNDEV  = "dev"
	RUNTEST = "test"
	RUNPROD = "prod"

	DEFAULTPATH       = "./conf"
	DEFAULTCONFIGTYPE = "toml"
	DEFAULTMODE       = RUNDEV
)

type Configx struct {
	mode  string
	path  string
	viper *viper.Viper
	err   error
}

type Option func(*Configx)

func (c *Configx) setError(err error) {
	if c.err == nil {
		c.err = err
	}
}

func WithMode(mode string) Option {
	return func(cli *Configx) {
		if mode != RUNDEV && mode != RUNPROD && mode != RUNTEST {
			cli.setError(fmt.Errorf("invalid mode %q: only support %s, %s, %s", mode, RUNDEV, RUNPROD, RUNTEST))
			return
		}
		cli.mode = mode
	}
}

func WithPath(path string) Option {
	return func(cli *Configx) {
		if path == "" {
			cli.setError(fmt.Errorf("config path must not be empty"))
			return
		}
		cli.path = path
	}
}

// Load reads configuration and returns initialization errors to the caller.
// Prefer this path when the composition root needs to decide whether to exit,
// retry, report, or recover from configuration failures.
func Load(options ...Option) (*Configx, error) {
	cli := &Configx{}
	for _, opt := range options {
		opt(cli)
	}
	if cli.err != nil {
		return nil, cli.err
	}

	v := viper.New()
	v.SetConfigType(DEFAULTCONFIGTYPE)
	v.SetDefault(RUNMODESTR, DEFAULTMODE)

	if err := v.BindEnv(RUNMODESTR); err != nil {
		return nil, fmt.Errorf("bind %s environment variable: %w", RUNMODESTR, err)
	}

	if cli.mode != "" {
		v.Set(RUNMODESTR, cli.mode)
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
		return nil, fmt.Errorf("read %s config from %s: %w", cli.mode, cli.path, err)
	}

	log.Printf("Using config mode: %s, file: %s \n", v.GetString(RUNMODESTR), v.ConfigFileUsed())
	cli.viper = v
	return cli, nil
}

// New is the fail-fast compatibility constructor.
//
// Prefer Load when the caller needs explicit error handling. New panics on
// initialization failure rather than terminating the process with log.Fatal.
func New(options ...Option) *Configx {
	cli, err := Load(options...)
	if err != nil {
		panic(err)
	}
	return cli
}

func (c *Configx) Mode() string {
	return c.viper.GetString(RUNMODESTR)
}

func (c *Configx) Viper() *viper.Viper {
	return c.viper
}
