package crawler

// APIResponse represents the generic response structure from the API
type APIResponse struct {
	Code      any         `json:"code"`
	Msg       string      `json:"msg"`
	Page      any         `json:"page"`
	PageCount any         `json:"pagecount"`
	Limit     string      `json:"limit"`
	Total     any         `json:"total"`
	List      []VideoItem `json:"list"`
}

// VideoItem represents a video object in the `provide` API response
type VideoItem struct {
	ID          int64    `json:"id"`
	TypeName    string   `json:"type_name"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	OriginName  string   `json:"origin_name"`
	PosterURL   string   `json:"poster_url"`
	ThumbURL    string   `json:"thumb_url"`
	Category    []string `json:"category"`
	Actor       []string `json:"actor"`
	Director    []string `json:"director"`
	CreatedAt   string   `json:"created_at"`
	Time        string   `json:"time"`
	Description string   `json:"description"`
	Episodes    struct {
		ServerData map[string]struct {
			LinkEmbed string `json:"link_embed"`
		} `json:"server_data"`
	} `json:"episodes"`
}

// APIResponseProvide1 represents the response from `provide1` API
type APIResponseProvide1 struct {
	Code      any                 `json:"code"`
	Msg       string              `json:"msg"`
	Page      any                 `json:"page"`
	PageCount any                 `json:"pagecount"`
	List      []VideoItemProvide1 `json:"list"`
}

// VideoItemProvide1 represents a video object in the `provide1` API response
type VideoItemProvide1 struct {
	VodID     int64  `json:"vod_id"`
	VodWriter string `json:"vod_writer"`
}
