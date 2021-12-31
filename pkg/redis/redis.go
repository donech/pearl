package redis

import (
	"context"
	"sync"

	"github.com/donech/pearl/pkg/log"
	"github.com/go-redis/redis"
)

const Default = "default"

var (
	redisMap          sync.Map
	defultRedisClient *redis.Client
)

func New(cfg Config) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if cfg.Name == Default {
		defultRedisClient = client
	}
	redisMap.Store(cfg.Name, client)
}

func Redis(name ...string) *redis.Client {
	if len(name) == 0 || name[0] == Default {
		if defultRedisClient == nil {
			log.S(context.Background()).Panicf("unknown redis.%s (forgotten configure?)", Default)
		}

		return defultRedisClient
	}

	v, ok := redisMap.Load(name[0])

	if !ok {
		log.S(context.Background()).Panicf("unknown db.%s (forgotten configure?)", name[0])
	}

	return v.(*redis.Client)
}
