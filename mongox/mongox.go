/*
 * @Author: hugo
 * @Date: 2024-04-19 17:19
 * @LastEditors: hugo
 * @LastEditTime: 2024-05-17 15:02
 * @FilePath: \gotox\mongox\mongox.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package mongox

import (
	"context"
	"sync"

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/hugo2lee/gotox/resourcex"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var _ resourcex.Resource = (*Mongox)(nil)

type Mongox struct {
	mongo  *mongo.Database
	logger logx.Logger
}

func New(conf *configx.Configx, logCli logx.Logger) (*Mongox, error) {
	uri := conf.MongoUri()
	if uri == "" {
		return nil, errors.New("mongo uri is empty")
	}

	dbName := conf.MongoDb()
	if dbName == "" {
		return nil, errors.New("mongo dbName is empty")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errors.Wrap(err, "mongo connect error")
	}
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, errors.Wrap(err, "mongo ping error")
	}
	return &Mongox{mongo: client.Database(dbName), logger: logCli}, nil
}

func (c *Mongox) Name() string {
	return "mongo"
}

func (c *Mongox) DB() *mongo.Database {
	return c.mongo
}

func (c *Mongox) Close(ctx context.Context, wg *sync.WaitGroup) {
	if err := c.mongo.Client().Disconnect(ctx); err != nil {
		c.logger.Error("mongo close error %v", err)
		return
	}
	wg.Done()
	c.logger.Info("%s close", c.Name())
}
