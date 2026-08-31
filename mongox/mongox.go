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

// NewWithDatabase wires an existing Mongo database into Mongox without network I/O.
func NewWithDatabase(database *mongo.Database, logger logx.Logger) *Mongox {
	if logger == nil {
		logger = logx.NewNoOpLogger()
	}
	return &Mongox{mongo: database, logger: logger}
}

// Dial creates, verifies, and selects a Mongo database from configuration.
func Dial(ctx context.Context, conf *configx.Configx, logger logx.Logger) (*Mongox, error) {
	uri := conf.MongoUri()
	if uri == "" {
		return nil, errors.New("mongo uri is empty")
	}

	dbName := conf.MongoDb()
	if dbName == "" {
		return nil, errors.New("mongo dbName is empty")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errors.Wrap(err, "mongo connect error")
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, errors.Wrap(err, "mongo ping error")
	}

	return NewWithDatabase(client.Database(dbName), logger), nil
}

// New is the compatibility constructor that dials Mongo immediately.
//
// Deprecated: use Dial when connection ownership belongs here, or
// NewWithDatabase when the caller/composition root already owns the client.
func New(conf *configx.Configx, logger logx.Logger) (*Mongox, error) {
	return Dial(context.Background(), conf, logger)
}

func (c *Mongox) Name() string {
	return "mongo"
}

func (c *Mongox) DB() *mongo.Database {
	return c.mongo
}

func (c *Mongox) Close(ctx context.Context) error {
	if err := c.mongo.Client().Disconnect(ctx); err != nil {
		c.logger.Error("mongo close error %v", err)
		return err
	}
	c.logger.Info("%s close", c.Name())
	return nil
}
