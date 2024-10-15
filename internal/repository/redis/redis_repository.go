package redis

import (
	"MessagesService/config"
	"context"
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

func (repo *RedisRepository) Set(ctx context.Context, key string, value string) error {
	if err := repo.client.Set(ctx, key, value, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (repo *RedisRepository) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := repo.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return val, nil
}
