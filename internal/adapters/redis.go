package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/BeatEcoprove/identityService/config"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
	"github.com/redis/go-redis/v9"
)

var redisConnection *RedisConnection

type RedisConnection struct {
	client *redis.Client
	ctx    context.Context
}

func GetRedis() *RedisConnection {
	if redisConnection == nil {
		redisConnection = newRedisConnection()
	}

	return redisConnection
}

func newRedisConnection() *RedisConnection {
	config := config.GetConfig()

	return &RedisConnection{
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", config.REDIS_HOST, config.REDIS_PORT),
			DB:   config.REDIS_DB,
		}),
		ctx: context.Background(),
	}
}

func (r *RedisConnection) GetValue(key interfaces.RedisKey) (string, error) {
	return r.client.Get(r.ctx, key.Key).Result()
}

func (r *RedisConnection) SetValue(key interfaces.RedisKey, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key.Key, value, expiration).Err()
}

func (r *RedisConnection) GetAndDelValue(key interfaces.RedisKey) (string, error) {
	return r.client.GetDel(r.ctx, key.Key).Result()
}

func (r *RedisConnection) Close() error {
	return r.client.Close()
}
