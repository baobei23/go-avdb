package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Video interface {
		Create(ctx context.Context) error
	}
	// Actor interface {
	// 	Create(context.Context) error
	// }
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Video: &VideoStore{db: db},
	}
}
