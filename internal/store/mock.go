package store

import "context"

func NewMockStore() Storage {
	return Storage{
		Video: &MockVideoStore{},
	}
}

type MockVideoStore struct{}

func (m *MockVideoStore) Upsert(ctx context.Context, video Video) error {
	return nil
}

func (m *MockVideoStore) UpsertActor(ctx context.Context, videoID int64, actor []string) error {
	return nil
}

func (m *MockVideoStore) UpsertTag(ctx context.Context, videoID int64, tag []string) error {
	return nil
}

func (m *MockVideoStore) UpsertDirector(ctx context.Context, videoID int64, director []string) error {
	return nil
}

func (m *MockVideoStore) UpsertStudio(ctx context.Context, videoID int64, studio string) error {
	return nil
}

func (m *MockVideoStore) GetVideoBySlug(ctx context.Context, slug string) (*Video, error) {
	return &Video{}, nil
}

func (m *MockVideoStore) GetVideoList(ctx context.Context) ([]Video, error) {
	return []Video{}, nil
}
