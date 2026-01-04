package store

import "context"

func NewMockStore() Storage {
	return Storage{
		Video:    &MockVideoStore{},
		Actor:    &MockActorStore{},
		Director: &MockDirectorStore{},
		Tag:      &MockTagStore{},
		Studio:   &MockStudioStore{},
	}
}

type MockVideoStore struct{}

func (m *MockVideoStore) Upsert(ctx context.Context, video *Video) error {
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

func (m *MockVideoStore) GetBySlug(ctx context.Context, slug string) (*Video, error) {
	return &Video{}, nil
}

func (m *MockVideoStore) GetByActor(ctx context.Context, actor string) ([]Video, error) {
	return []Video{}, nil
}

func (m *MockVideoStore) GetList(ctx context.Context, limit, offset int, search string) ([]VideoList, int, error) {
	return []VideoList{}, 0, nil
}

type MockActorStore struct{}

func (m *MockActorStore) Create(ctx context.Context, actor *Actor) error {
	return nil
}

func (m *MockActorStore) Update(ctx context.Context, actor *Actor) error {
	return nil
}

func (m *MockActorStore) GetList(ctx context.Context) ([]Actor, error) {
	return []Actor{}, nil
}

func (m *MockActorStore) Delete(ctx context.Context, id int64) error {
	return nil
}

type MockDirectorStore struct{}

func (m *MockDirectorStore) Create(ctx context.Context, director *Director) error {
	return nil
}

func (m *MockDirectorStore) GetByID(ctx context.Context, id int64) (*Director, error) {
	return &Director{}, nil
}

func (m *MockDirectorStore) GetList(ctx context.Context) ([]Director, error) {
	return []Director{}, nil
}

type MockTagStore struct{}

func (m *MockTagStore) Create(ctx context.Context, tag *Tag) error {
	return nil
}

func (m *MockTagStore) GetByID(ctx context.Context, id int64) (*Tag, error) {
	return &Tag{}, nil
}

func (m *MockTagStore) GetList(ctx context.Context) ([]Tag, error) {
	return []Tag{}, nil
}

type MockStudioStore struct{}

func (m *MockStudioStore) Create(ctx context.Context, studio *Studio) error {
	return nil
}

func (m *MockStudioStore) GetList(ctx context.Context) ([]Studio, error) {
	return []Studio{}, nil
}
