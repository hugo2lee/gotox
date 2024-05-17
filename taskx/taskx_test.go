/*
 * @Author: hugo
 * @Date: 2024-05-10 15:06
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-11 13:50
 * @FilePath: \gotox\taskx\taskx_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */

package taskx_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/hugo2lee/gotox/taskx"
)

func TestTasker(t *testing.T) {
	t.Log("hello")

	timeOut, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tasker := taskx.NewTaskxGroup()
	tasker.AddTask(taskx.NewTaskx("ping", func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				log.Printf("timeout \n")
				return
			default:
				log.Printf("ping \n")
				time.Sleep(1 * time.Second)
			}
		}
	}))
	tasker.Run(timeOut)

	time.Sleep(5 * time.Second)

	t.Log("finish")
}
