package store

import "time"

type Video struct {
	ID          int64      `json:"id"`
	Category    string     `json:"category"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	OriginName  string     `json:"origin_name"`
	PosterURL   string     `json:"poster_url"`
	ThumbURL    string     `json:"thumb_url"`
	Description string     `json:"description,omitempty"`
	LinkEmbed   string     `json:"link_embed,omitempty"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Actor       []string   `json:"actor"`
	Tag         []string   `json:"tag"`
	Studio      []string   `json:"studio"`
	Director    []string   `json:"director"`
}

type VideoList struct {
	ID         int64      `json:"id"`
	Category   string     `json:"category"`
	Name       string     `json:"name"`
	Slug       string     `json:"slug"`
	OriginName string     `json:"origin_name"`
	PosterURL  string     `json:"poster_url"`
	ThumbURL   string     `json:"thumb_url"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type Actor struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Studio struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Director struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
