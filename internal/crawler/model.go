package crawler

// APIResponse untuk endpoint provide (base data)
type APIResponse struct {
	Code      interface{} `json:"code"`
	Msg       string      `json:"msg"`
	Page      interface{} `json:"page"`
	PageCount interface{} `json:"pagecount"`
	Limit     string      `json:"limit"`
	Total     interface{} `json:"total"`
	List      []VideoItem `json:"list"`
}

// VideoItem dari API provide
type VideoItem struct {
	ID          int64       `json:"vod_id"`
	TypeName    string      `json:"type_name"`
	Name        string      `json:"vod_name"`
	Slug        string      `json:"vod_en"`
	OriginName  string      `json:"vod_sub"`
	PosterURL   string      `json:"vod_pic"`
	ThumbURL    string      `json:"vod_pic_thumb"`
	Description string      `json:"vod_content"`
	Actor       []string    `json:"vod_actor"`
	Category    []string    `json:"vod_class"`
	Director    []string    `json:"vod_director"`
	CreatedAt   string      `json:"vod_time"`
	Episodes    EpisodeData `json:"vod_play_url"`
}

// EpisodeData struktur nested
type EpisodeData struct {
	ServerData map[string]ServerInfo `json:"server_data"`
}

type ServerInfo struct {
	LinkEmbed string `json:"url"`
}

// APIResponseProvide1 untuk endpoint provide1 (supplemental data)
type APIResponseProvide1 struct {
	Code      interface{}         `json:"code"`
	Msg       string              `json:"msg"`
	Page      interface{}         `json:"page"`
	PageCount interface{}         `json:"pagecount"`
	Limit     string              `json:"limit"`
	Total     interface{}         `json:"total"`
	List      []VideoItemProvide1 `json:"list"`
}

// VideoItemProvide1 dari API provide1 (hanya untuk studio)
type VideoItemProvide1 struct {
	VodID     int64  `json:"vod_id"`
	VodWriter string `json:"vod_writer"` // Studio name
}
