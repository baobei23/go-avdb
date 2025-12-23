package store

import (
	"context"
	"database/sql"
)

type VideoStore struct {
	db *sql.DB
}

func (s *VideoStore) Create(ctx context.Context) error {
	return nil
}
