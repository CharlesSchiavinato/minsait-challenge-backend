package cache

import (
	"context"
	"time"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/service/cache"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/util"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client     *redis.Client
	Expiration time.Duration
}

func NewRedis(config *util.Config) (cache.Cache, error) {
	opt, err := redis.ParseURL(config.CacheURL)

	if err != nil {
		return nil, err
	}

	expiration, err := time.ParseDuration(config.CacheExpiration)

	if err != nil {
		return nil, err
	}

	cacheRedis := &Redis{
		Client:     redis.NewClient(opt),
		Expiration: expiration,
	}

	err = cacheRedis.Client.Ping(context.Background()).Err()

	return cacheRedis, err
}

func (redis *Redis) Close() error {
	return redis.Client.Close()
}

func (redis *Redis) Check() error {
	return redis.Client.Ping(context.Background()).Err()
}
