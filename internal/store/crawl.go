package store

import (
	"context"
	"fmt"
)

func (s *videoStore) Upsert(ctx context.Context, video *Video) error {
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
	_, err := s.db.Exec(ctx, query,
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

func (s *videoStore) upsertBatchRefs(ctx context.Context, table string, names []string) ([]int64, error) {
	if len(names) == 0 {
		return nil, nil
	}

	queryInsert := fmt.Sprintf("INSERT INTO %s (name) SELECT unnest($1::text[]) ON CONFLICT (name) DO NOTHING", table)
	_, err := s.db.Exec(ctx, queryInsert, names)
	if err != nil {
		return nil, fmt.Errorf("failed to batch insert %s: %w", table, err)
	}

	querySelect := fmt.Sprintf("SELECT id FROM %s WHERE name = ANY($1::text[])", table)
	rows, err := s.db.Query(ctx, querySelect, names)
	if err != nil {
		return nil, fmt.Errorf("failed to select ids %s: %w", table, err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *videoStore) linkBatchRefs(ctx context.Context, table string, refColumn string, videoID int64, refIDs []int64) error {
	if len(refIDs) == 0 {
		return nil
	}
	query := fmt.Sprintf("INSERT INTO %s (video_id, %s) SELECT $1, unnest($2::bigint[]) ON CONFLICT DO NOTHING", table, refColumn)
	_, err := s.db.Exec(ctx, query, videoID, refIDs)
	return err
}

func (s *videoStore) UpsertActor(ctx context.Context, videoID int64, actor []string) error {
	ids, err := s.upsertBatchRefs(ctx, "actor", actor)
	if err != nil {
		return err
	}
	return s.linkBatchRefs(ctx, "video_actor", "actor_id", videoID, ids)
}

func (s *videoStore) UpsertTag(ctx context.Context, videoID int64, tag []string) error {
	ids, err := s.upsertBatchRefs(ctx, "tag", tag)
	if err != nil {
		return err
	}
	return s.linkBatchRefs(ctx, "video_tag", "tag_id", videoID, ids)
}

func (s *videoStore) UpsertDirector(ctx context.Context, videoID int64, director []string) error {
	ids, err := s.upsertBatchRefs(ctx, "director", director)
	if err != nil {
		return err
	}
	return s.linkBatchRefs(ctx, "video_director", "director_id", videoID, ids)
}

func (s *videoStore) UpsertStudio(ctx context.Context, videoID int64, studioName string) error {
	if studioName == "" {
		return nil
	}
	ids, err := s.upsertBatchRefs(ctx, "studio", []string{studioName})
	if err != nil {
		return err
	}
	return s.linkBatchRefs(ctx, "video_studio", "studio_id", videoID, ids)
}
