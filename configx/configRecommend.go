/*
 * @Author: hugo
 * @Date: 2024-04-02 10:16
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-02 11:32
 * @FilePath: \gotox\config\configFetch.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package configx

// log config
func (c *ConfigCli) LogDir() string {
	return c.viper.GetString("log.dir")
}

// server addr config
func (c *ConfigCli) Addr() string {
	return c.viper.GetString("server.addr")
}

// redis config
func (c *ConfigCli) RedisAddr() string {
	return c.viper.GetString("redis.addr")
}

func (c *ConfigCli) RedisPassword() string {
	return c.viper.GetString("redis.password")
}

func (c *ConfigCli) RedisDb() int {
	return c.viper.GetInt("redis.db")
}

// mysql config
func (c *ConfigCli) MysqlDsn() string {
	return c.viper.GetString("mysql.dsn")
}

// mongo config
func (c *ConfigCli) MongoUri() string {
	return c.viper.GetString("mongo.uri")
}

func (c *ConfigCli) MongoDb() string {
	return c.viper.GetString("mongo.db")
}
