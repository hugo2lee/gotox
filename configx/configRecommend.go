/*
 * @Author: hugo
 * @Date: 2024-04-02 10:16
 * @LastEditors: hugo
 * @LastEditTime: 2024-06-03 16:34
 * @FilePath: \gotox\configx\configRecommend.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package configx

import "time"

// log config
func (c *Configx) LogDir() string {
	return c.viper.GetString("log.dir")
}

// server addr config
func (c *Configx) Addr() string {
	return c.viper.GetString("server.addr")
}

// cache config
func (c *Configx) RedisUrl() string {
	return c.viper.GetString("redis.url")
}

// redis config
func (c *Configx) CachexDefaultExpiration() time.Duration {
	return c.viper.GetDuration("cache.defaultExpirationSec")
}

// mysql config
func (c *Configx) MysqlDsn() string {
	return c.viper.GetString("mysql.dsn")
}

// mongo config
func (c *Configx) MongoUri() string {
	return c.viper.GetString("mongo.uri")
}

func (c *Configx) MongoDb() string {
	return c.viper.GetString("mongo.db")
}

func (c *Configx) Auths() map[string]string {
	return c.viper.GetStringMapString("auths")
}
