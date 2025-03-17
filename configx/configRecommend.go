/*
 * @Author: hugo
 * @Date: 2024-04-02 10:16
 * @LastEditors: hugo2lee
 * @LastEditTime: 2025-03-17 14:23
 * @FilePath: \gotox\configx\configRecommend.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package configx

import (
	"fmt"
	"time"
)

// log config
func (c *Configx) LogDir() string {
	return c.viper.GetString("log.dir")
}

// server addr config
func (c *Configx) Addr() string {
	return c.viper.GetString("server.addr")
}

// redis config
func (c *Configx) RedisUrl() string {
	return c.viper.GetString("redis.url")
}

// cache config
func (c *Configx) CachexDefaultExpiration() time.Duration {
	return c.viper.GetDuration("cache.defaultExpirationSec")
}

func (c *Configx) CachexCleanupInterval() time.Duration {
	return c.viper.GetDuration("cache.cleanupIntervalSec")
}

// mysql config
func (c *Configx) MysqlDsn() string {
	return c.viper.GetString("mysql.dsn")
}

// mysql multiple config
func (c *Configx) MysqlDsnWithName(name string) string {
	return  c.viper.GetString(fmt.Sprintf("mysql.%s", name))
}

// pg config
func (c *Configx) PostgresDsn() string {
	return c.viper.GetString("postgres.dsn")
}

// mongo config
func (c *Configx) MongoUri() string {
	return c.viper.GetString("mongo.uri")
}

func (c *Configx) MongoDb() string {
	return c.viper.GetString("mongo.db")
}

// kafka config
func (c *Configx) KafkaBrokers() []string {
	return c.viper.GetStringSlice("kafka.brokers")
}

func (c *Configx) KafkaTopic() string {
	return c.viper.GetString("kafka.topic")
}

func (c *Configx) KafkaGroup() string {
	return c.viper.GetString("kafka.group")
}

// auth config
func (c *Configx) Auths() map[string]string {
	return c.viper.GetStringMapString("auths")
}
