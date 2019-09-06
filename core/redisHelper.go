package core

import (
	"fmt"
	//"github.com/gomodule/redigo/redis"
	"gopkg.in/redis.v3"
	"time"
)

var dbClient *redis.Client

func InitRedis() error {
	c := GetConfig()
	dbClient = redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort),
		MaxRetries:  3,
		ReadTimeout: time.Millisecond * 1000,
		PoolSize:    100,
		PoolTimeout: time.Millisecond * 300,
		Password:    c.RedisAuth,
	})

	_, err := dbClient.Ping().Result()
	return err
}

func Get(key string) (string, error) {
	return dbClient.Get(key).Result()
}

func Set(key, value string, expiration time.Duration) (string, error) {
	return dbClient.Set(key, value, expiration).Result()
}

func Expire(key string, expire time.Duration) (bool, error) {
	return dbClient.Expire(key, expire).Result()
}

func Delete(key string) (int64, error) {
	return dbClient.Del(key).Result()
}

func CloseRedis() error {
	return dbClient.Close()
}
