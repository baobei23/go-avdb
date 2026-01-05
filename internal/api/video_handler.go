package api

import (
	"math"
	"net/http"

	"github.com/baobei23/go-avdb/internal/store"
	"github.com/go-chi/chi/v5"
)

// TODO: implement cache
type getVideoList struct {
	Data      []store.VideoList `json:"data"`
	Total     int               `json:"total"`
	Page      int               `json:"page"`
	PageCount int               `json:"page_count"`
	Limit     int               `json:"limit"`
}

// getVideo godoc
//
//	@Summary		Get video detail by slug
//	@Description	Get video detail by slug
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string	true	"Video slug"
//	@Success		200		{object}	[]store.VideoList
//	@Failure		500		{object}	error
//	@Router			/video/{slug} [get]
func (app *Application) getVideo(w http.ResponseWriter, r *http.Request) {
	video, err := app.Store.Video.GetBySlug(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, video); err != nil {
		app.internalServerError(w, r, err)
	}
}

// getVideoList godoc
//
//	@Summary		Get video list
//	@Description	Get video list
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int		false	"Page"
//	@Param			limit	query		int		false	"Limit"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	getVideoList
//	@Failure		500		{object}	error
//	@Router			/video [get]
func (app *Application) getVideoList(w http.ResponseWriter, r *http.Request) {
	pq := store.PaginationQuery{}
	pq, err := pq.Parse(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	videos, total, err := app.Store.Video.GetList(r.Context(), pq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, getVideoList{
		Data:      videos,
		Total:     total,
		Page:      pq.Page,
		PageCount: int(math.Ceil(float64(total) / float64(pq.Limit))),
		Limit:     pq.Limit,
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}
