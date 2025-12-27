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
	query := `
        SELECT 
            v.id, v.category, v.name, v.slug, v.origin_name, v.poster_url, v.thumb_url, 
            v.description, v.link_embed, v.created_at, v.updated_at,
            COALESCE(a.actors, '{}'),
            COALESCE(t.tags, '{}'),
            COALESCE(s.studios, '{}'),
            COALESCE(d.directors, '{}')
        FROM video v
        LEFT JOIN LATERAL (
            SELECT ARRAY_AGG(a.name ORDER BY a.name) AS actors
            FROM video_actor va JOIN actor a ON a.id = va.actor_id
            WHERE va.video_id = v.id
        ) a ON TRUE
        LEFT JOIN LATERAL (
            SELECT ARRAY_AGG(t.name ORDER BY t.name) AS tags
            FROM video_tag vt JOIN tag t ON t.id = vt.tag_id
            WHERE vt.video_id = v.id
        ) t ON TRUE
        LEFT JOIN LATERAL (
            SELECT ARRAY_AGG(s.name ORDER BY s.name) AS studios
            FROM video_studio vs JOIN studio s ON s.id = vs.studio_id
            WHERE vs.video_id = v.id
        ) s ON TRUE
        LEFT JOIN LATERAL (
            SELECT ARRAY_AGG(d.name ORDER BY d.name) AS directors
            FROM video_director vd JOIN director d ON d.id = vd.director_id
            WHERE vd.video_id = v.id
        ) d ON TRUE
        WHERE v.slug = $1
    `

	var v Video
	err := s.db.QueryRow(ctx, query, slug).Scan(
		&v.ID, &v.Category, &v.Name, &v.Slug, &v.OriginName,
		&v.PosterURL, &v.ThumbURL, &v.Description, &v.LinkEmbed,
		&v.CreatedAt, &v.UpdatedAt,
		&v.Actor, &v.Tag, &v.Studio, &v.Director,
	)
	if err != nil {
		return nil, err
	}
	return &v, nil

}

func (s *videoStore) GetVideo(ctx context.Context, limit, offset int, search string) ([]Video, int, error) {
	return nil, 0, nil
}
