package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type videoStore struct {
	db *pgxpool.Pool
}

func (s *videoStore) GetList(ctx context.Context) ([]Video, error) {

	rows, err := s.db.Query(ctx, "SELECT id, name, slug, poster_url FROM video ORDER BY created_at DESC LIMIT 20")
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
	err := s.db.QueryRow(ctx, query, slug).Scan(&v.ID, &v.Name, &v.Slug, &v.PosterURL, &v.Description)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *videoStore) GetVideo(
	ctx context.Context,
	limit, offset int,
	search string,
) ([]Video, int, error) {

	if search != "" {
		search = "%" + search + "%"
	}

	// count
	var total int
	err := s.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM video
		WHERE ($1 = '' OR name ILIKE $1)
	`, search).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// main query (NO JOIN)
	rows, err := s.db.Query(ctx, `
		SELECT
			id, category, name, slug, origin_name,
			poster_url, thumb_url, description,
			link_embed, created_at, updated_at
		FROM video
		WHERE ($1 = '' OR name ILIKE $1)
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
	`, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var v Video
		if err := rows.Scan(
			&v.ID, &v.Category, &v.Name, &v.Slug, &v.OriginName,
			&v.PosterURL, &v.ThumbURL, &v.Description,
			&v.LinkEmbed, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		videos = append(videos, v)
	}

	return videos, total, nil
}
