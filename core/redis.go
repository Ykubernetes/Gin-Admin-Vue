package core

import (
	"gitee.com/go-server/global"
	"github.com/go-redis/redis/v8"
)

func RedisConn() *redis.Client {
	RedisCli := redis.NewClient(&redis.Options{
		Addr:         global.Config.Redis.GetHost(),
		Password:     global.Config.Redis.Password,
		Username:     "default",
		DB:           global.Config.Redis.SelectDb,
		PoolSize:     global.Config.Redis.PolSize,
		MinIdleConns: global.Config.Redis.MinIdleConn,
	})
	global.Log.Infof("Redis Server is Running in %s", global.Config.Redis.GetHost())
	return RedisCli
}
