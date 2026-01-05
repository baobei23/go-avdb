package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/baobei23/go-avdb/internal/store"
	"github.com/redis/go-redis/v9"
)

type VideoStore struct {
	redis *redis.Client
}

const VideoExpireTime = time.Hour

func (s *VideoStore) Get(ctx context.Context, slug string) (*store.Video, error) {
	key := fmt.Sprintf("video:%s", slug)

	val, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var video store.Video
	if err := json.Unmarshal([]byte(val), &video); err != nil {
		return nil, err
	}

	return &video, nil
}

func (s *VideoStore) Set(ctx context.Context, slug string, video *store.Video) error {
	key := fmt.Sprintf("video:slug:%s", video.Slug)

	data, err := json.Marshal(video)
	if err != nil {
		return err
	}

	return s.redis.Set(ctx, key, data, VideoExpireTime).Err()
}

func (s *VideoStore) GetList(ctx context.Context, limit, offset int, search string) ([]store.VideoList, int, error) {
	key := fmt.Sprintf("video:list:%d:%d:%s", limit, offset, search)

	val, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}

	type cachedList struct {
		Videos []store.VideoList `json:"videos"`
		Total  int               `json:"total"`
	}

	var cached cachedList
	if err := json.Unmarshal([]byte(val), &cached); err != nil {
		return nil, 0, err
	}

	return cached.Videos, cached.Total, nil
}

func (s *VideoStore) SetList(ctx context.Context, limit, offset int, search string, videos []store.VideoList, total int) error {
	key := fmt.Sprintf("video:list:%d:%d:%s", limit, offset, search)

	type cachedList struct {
		Videos []store.VideoList `json:"videos"`
		Total  int               `json:"total"`
	}

	payload := cachedList{
		Videos: videos,
		Total:  total,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return s.redis.Set(ctx, key, data, VideoExpireTime).Err()
}
