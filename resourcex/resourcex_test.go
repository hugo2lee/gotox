/*
 * @Author: hugo
 * @Date: 2024-05-11 14:00
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-11 15:11
 * @FilePath: \gotox\resourcex\resourcex_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package resourcex_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/redisx"
	"github.com/hugo2lee/gotox/resourcex"
	"github.com/stretchr/testify/assert"
)

func TestResourcer(t *testing.T) {
	t.Log("hello")

	re := resourcex.NewResourcexGroup()
	re.AddResource(resourcex.NewResourcex("redis", func(ctx context.Context) {
		log.Println("outter call close")
		for {
			select {
			case <-ctx.Done():
				log.Printf("outter timeout \n")
				return
			default:
				log.Printf("ping \n")
				time.Sleep(1 * time.Second)
				return
			}
		}
	}))

	timeOut, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	re.CloseAll(timeOut)

	time.Sleep(5 * time.Second)

	t.Log("finish")
}

func TestRedisResource(t *testing.T) {
	t.Log("hello")

	rds, err := redisx.New(configx.New(configx.WithPath("../conf")), logx.New(configx.New(configx.WithPath("../conf"))))
	assert.NoError(t, err)

	re := resourcex.NewResourcexGroup()
	re.AddResource(rds)

	timeOut, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	re.CloseAll(timeOut)

	time.Sleep(5 * time.Second)

	t.Log("finish")
}
