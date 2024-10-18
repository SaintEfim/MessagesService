package interfaces

import "context"

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
}
