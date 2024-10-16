package redis

import (
	"context"
	"time"

	"MessagesService/config"
	"MessagesService/internal/models/interfaces"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisRepository struct {
	client *redis.Client
	cfg    *config.Config
}

func NewRedisClient(ctx context.Context, cfg *config.Config, logger *zap.Logger) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})
}

func NewRedisRepository(client *redis.Client, cfg *config.Config) interfaces.RedisRepository {
	return &RedisRepository{
		client: client,
		cfg:    cfg,
	}
}

func (repo *RedisRepository) Set(ctx context.Context, key string, value string) error {
	ctxNew, cancel := context.WithTimeout(ctx, repo.cfg.Redis.Timeout*time.Second)
	if ctxNew.Err() != nil {
		cancel()
		return ctxNew.Err()
	}

	if err := repo.client.Set(ctxNew, key, value, repo.cfg.Redis.Expiration*time.Second).Err(); err != nil {
		return err
	}
	return nil
}

func (repo *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	ctxNew, cancel := context.WithTimeout(ctx, repo.cfg.Redis.Timeout*time.Second)
	if ctxNew.Err() != nil {
		cancel()
		return "", ctxNew.Err()
	}

	val, err := repo.client.Get(ctxNew, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
