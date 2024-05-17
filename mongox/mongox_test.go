/*
 * @Author: hugo
 * @Date: 2024-04-19 17:23
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-19 17:23
 * @FilePath: \gotox\mongox\mongox_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package mongox_test

import (
	"context"
	"log"
	"testing"

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/mongox"
	"github.com/stretchr/testify/assert"
)

func TestMongo(t *testing.T) {
	t.Parallel()
	conf := configx.New(configx.WithPath("../conf"))
	logger := logx.New(conf)
	db, err := mongox.New(conf, logger)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	reslut, err := db.DB().Collection("user").InsertOne(context.TODO(), map[string]string{"name": "hugo"})
	assert.NoError(t, err)
	log.Println(reslut)
}
