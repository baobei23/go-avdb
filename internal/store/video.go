package store

import (
	"context"
	"database/sql"
)

type videoStore struct {
	db *sql.DB
}

// Upsert video (INSERT ON CONFLICT UPDATE)
func (s *videoStore) Upsert(ctx context.Context, video Video) error {

	return nil
}

// Helper generic upsert for reference tables
func (s *videoStore) upsertRef(ctx context.Context, table string, value string) (int, error) {
	return 0, nil
}

func (s *videoStore) UpsertActor(ctx context.Context, videoID int64, actor []string) error {

	return nil
}

func (s *videoStore) UpsertTag(ctx context.Context, videoID int64, tag []string) error {

	return nil
}

func (s *videoStore) UpsertDirector(ctx context.Context, videoID int64, director []string) error {

	return nil
}

func (s *videoStore) UpsertStudio(ctx context.Context, videoID int64, studioName string) error {
	return nil
}
