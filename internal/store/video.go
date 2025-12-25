package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type videoStore struct {
	db *sql.DB
}

// Upsert video (INSERT ON CONFLICT UPDATE)
func (s *videoStore) Upsert(ctx context.Context, video Video) error {
	query := `
		INSERT INTO video (id, category, name, slug, origin_name, poster_url, thumb_url, description, link_embed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
		ON CONFLICT (id) DO UPDATE SET
			category = EXCLUDED.category,
			name = EXCLUDED.name,
			slug = EXCLUDED.slug,
			origin_name = EXCLUDED.origin_name,
			poster_url = EXCLUDED.poster_url,
			thumb_url = EXCLUDED.thumb_url,
			description = EXCLUDED.description,
			link_embed = EXCLUDED.link_embed,
			created_at = EXCLUDED.created_at,
			updated_at = NOW();
	`
	_, err := s.db.ExecContext(ctx, query,
		video.ID,
		video.Category,
		video.Name,
		video.Slug,
		video.OriginName,
		video.PosterURL,
		video.ThumbURL,
		video.Description,
		video.LinkEmbed,
		video.CreatedAt,
	)
	return err
}

// Helper generic upsert for reference tables
func (s *videoStore) upsertRef(ctx context.Context, table string, value string) (int64, error) {
	if value == "" {
		return 0, nil
	}

	var id int64
	// Try to find existing
	queryFind := fmt.Sprintf("SELECT id FROM %s WHERE name = $1", table)
	err := s.db.QueryRowContext(ctx, queryFind, value).Scan(&id)
	if err == nil {
		return id, nil
	}

	// Insert new
	queryInsert := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id", table)
	err = s.db.QueryRowContext(ctx, queryInsert, value).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert %s: %w", table, err)
	}
	return id, nil
}

func (s *videoStore) UpsertActor(ctx context.Context, videoID int64, actor []string) error {
	for _, name := range actor {
		id, err := s.upsertRef(ctx, "actor", name)
		if err != nil {
			log.Printf("Error upserting actor %s: %v", name, err)
			continue
		}
		if id == 0 {
			continue
		}
		_, err = s.db.ExecContext(ctx, "INSERT INTO video_actor (video_id, actor_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", videoID, id)
		if err != nil {
			log.Printf("Error linking actor %d to video %d: %v", id, videoID, err)
		}
	}
	return nil
}

func (s *videoStore) UpsertTag(ctx context.Context, videoID int64, tag []string) error {
	for _, name := range tag {
		id, err := s.upsertRef(ctx, "tag", name)
		if err != nil {
			log.Printf("Error upserting tag %s: %v", name, err)
			continue
		}
		if id == 0 {
			continue
		}
		_, err = s.db.ExecContext(ctx, "INSERT INTO video_tag (video_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", videoID, id)
		if err != nil {
			log.Printf("Error linking tag %d to video %d: %v", id, videoID, err)
		}
	}
	return nil
}

func (s *videoStore) UpsertDirector(ctx context.Context, videoID int64, director []string) error {
	for _, name := range director {
		id, err := s.upsertRef(ctx, "director", name)
		if err != nil {
			log.Printf("Error upserting director %s: %v", name, err)
			continue
		}
		if id == 0 {
			continue
		}
		_, err = s.db.ExecContext(ctx, "INSERT INTO video_director (video_id, director_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", videoID, id)
		if err != nil {
			log.Printf("Error linking director %d to video %d: %v", id, videoID, err)
		}
	}
	return nil
}

func (s *videoStore) UpsertStudio(ctx context.Context, videoID int64, studioName string) error {
	id, err := s.upsertRef(ctx, "studio", studioName)
	if err != nil {
		return err
	}
	if id == 0 {
		return nil
	}
	_, err = s.db.ExecContext(ctx, "INSERT INTO video_studio (video_id, studio_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", videoID, id)
	return err
}

func (s *videoStore) GetVideoList(ctx context.Context) ([]Video, error) {
	// Implement simple list for now
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

func (s *videoStore) GetVideoBySlug(ctx context.Context, slug string) (*Video, error) {
	var v Video
	query := "SELECT id, name, slug, poster_url, description FROM video WHERE slug = $1"
	err := s.db.QueryRowContext(ctx, query, slug).Scan(&v.ID, &v.Name, &v.Slug, &v.PosterURL, &v.Description)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
