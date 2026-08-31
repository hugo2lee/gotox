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
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoDatabaseWiring(t *testing.T) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	require.NoError(t, err)
	database := client.Database("test")

	wrapped := mongox.NewWithDatabase(database, nil)
	assert.Same(t, database, wrapped.DB())
}

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
