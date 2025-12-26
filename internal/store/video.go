package store

import (
	"context"
	"database/sql"
)

type videoStore struct {
	db *sql.DB
}

func (s *videoStore) GetList(ctx context.Context) ([]Video, error) {

	rows, err := s.db.QueryContext(ctx, "SELECT id, name, slug, poster_url FROM video ORDER BY created_at DESC LIMIT 20")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var v Video
		if err := rows.Scan(&v.ID, &v.Name, &v.Slug, &v.PosterURL); err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	return videos, nil
}

func (s *videoStore) GetBySlug(ctx context.Context, slug string) (*Video, error) {
	var v Video
	query := "SELECT id, name, slug, poster_url, description FROM video WHERE slug = $1"
	err := s.db.QueryRowContext(ctx, query, slug).Scan(&v.ID, &v.Name, &v.Slug, &v.PosterURL, &v.Description)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
