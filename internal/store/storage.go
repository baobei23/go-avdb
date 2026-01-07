package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Video interface {
		// Video operations
		Upsert(ctx context.Context, video *Video) error

		// Read
		GetBySlug(ctx context.Context, slug string) (*Video, error)
		GetList(ctx context.Context, pq PaginationQuery) ([]VideoList, int, error)
		GetListByActor(ctx context.Context, actor string, pq PaginationQuery) ([]VideoList, int, error)
		GetListByDirector(ctx context.Context, director string, pq PaginationQuery) ([]VideoList, int, error)
		GetListByStudio(ctx context.Context, studio string, pq PaginationQuery) ([]VideoList, int, error)
		GetListByTag(ctx context.Context, tag string, pq PaginationQuery) ([]VideoList, int, error)

		// Relationship operations
		UpsertActor(ctx context.Context, videoID int64, actor []string) error
		UpsertTag(ctx context.Context, videoID int64, tag []string) error
		UpsertDirector(ctx context.Context, videoID int64, director []string) error
		UpsertStudio(ctx context.Context, videoID int64, studio string) error
	}
	Actor interface {
		Create(ctx context.Context, actor *Actor) error
		Update(ctx context.Context, actor *Actor) error
		GetList(ctx context.Context) ([]Actor, error)
		Delete(ctx context.Context, id int64) error
	}
	Tag interface {
		Create(ctx context.Context, tag *Tag) error
		GetList(ctx context.Context) ([]Tag, error)
	}
	Studio interface {
		Create(ctx context.Context, studio *Studio) error
		GetList(ctx context.Context) ([]Studio, error)
	}
	Director interface {
		Create(ctx context.Context, director *Director) error
		GetList(ctx context.Context) ([]Director, error)
	}
}

func NewStorage(db *pgxpool.Pool) Storage {
	return Storage{
		Video: &videoStore{db: db},
		Actor: &actorStore{db: db},
	}
}
