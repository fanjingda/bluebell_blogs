package redis

import (
	"bluebell_blogs/settings"
	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.Poolsize,
	})
	_, err = client.Ping().Result()
	if err != nil {
		zap.L().Error("connect redis failed,%v\n", zap.Error(err))
		return
	}
	return
}

func Close() {
	_ = client.Close()
}
