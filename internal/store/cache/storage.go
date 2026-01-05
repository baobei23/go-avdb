package cache

import (
	"context"

	"github.com/baobei23/go-avdb/internal/store"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Video interface {
		Get(ctx context.Context, slug string) (*store.Video, error)
		Set(ctx context.Context, slug string, video *store.Video) error
		GetList(ctx context.Context, limit, offset int, search string) ([]store.VideoList, int, error)
		SetList(ctx context.Context, limit, offset int, search string, videos []store.VideoList, total int) error
	}
}

func NewRedisStorage(redis *redis.Client) Storage {
	return Storage{
		Video: &VideoStore{redis: redis},
	}
}
