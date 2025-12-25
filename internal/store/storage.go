package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Video interface {
		// Video operations
		Upsert(ctx context.Context, video Video) error
		//GetByID(ctx context.Context, id int64) (*Video, error)

		// Relationship operations
		UpsertActor(ctx context.Context, videoID int64, actor []string) error
		UpsertTag(ctx context.Context, videoID int64, tag []string) error
		UpsertDirector(ctx context.Context, videoID int64, director []string) error
		UpsertStudio(ctx context.Context, videoID int64, studio string) error
	}
	Actor interface {
		Upsert(ctx context.Context, actor Actor) error
		GetByID(ctx context.Context, id int64) (*Actor, error)
		GetList(ctx context.Context) ([]Actor, error)
	}
	Tag interface {
		Upsert(ctx context.Context, tag Tag) error
		GetByID(ctx context.Context, id int64) (*Tag, error)
		GetList(ctx context.Context) ([]Tag, error)
	}
	Studio interface {
		Upsert(ctx context.Context, studio Studio) error
		GetByID(ctx context.Context, id int64) (*Studio, error)
		GetList(ctx context.Context) ([]Studio, error)
	}
	Director interface {
		Upsert(ctx context.Context, director Director) error
		GetByID(ctx context.Context, id int64) (*Director, error)
		GetList(ctx context.Context) ([]Director, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Video: &videoStore{db: db},
	}
}
