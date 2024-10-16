package redis

import (
	"MessagesService/config"
	"MessagesService/internal/models/interfaces"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
	cfg    *config.Config
}

func NewRedisRepository(client *redis.Client, cfg *config.Config) interfaces.RedisRepository {
	return &RedisRepository{
		client: client,
		cfg:    cfg,
	}
}

func (repo *RedisRepository) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	timeout := repo.cfg.Redis.Timeout
	if timeout == 0 {
		timeout = 5 * time.Hour
	}
	return context.WithTimeout(ctx, timeout)
}

func (repo *RedisRepository) Set(ctx context.Context, key string, value string) error {
	ctxNew := context.Background()

	if err := repo.client.Set(ctxNew, key, value, repo.cfg.Redis.Expiration*time.Second).Err(); err != nil {
		return err
	}

	return nil
}

func (repo *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	ctxNew := context.Background()

	val, err := repo.client.Get(ctxNew, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
