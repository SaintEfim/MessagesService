package redis

import (
	"MessagesService/config"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
	cfg    *config.Config
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

func (repo *RedisRepository) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	timeout := repo.cfg.Redis.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return context.WithTimeout(ctx, timeout)
}

func (repo *RedisRepository) Set(ctx context.Context, key string, value string) error {
	ctx, cancel := repo.withTimeout(ctx)
	defer cancel()

	if err := repo.client.Set(ctx, key, value, repo.cfg.Redis.Expiration*time.Second).Err(); err != nil {
		return err
	}

	return nil
}

func (repo *RedisRepository) Get(ctx context.Context, key string) (interface{}, error) {
	ctx, cancel := repo.withTimeout(ctx)
	defer cancel()

	val, err := repo.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return val, nil
}
