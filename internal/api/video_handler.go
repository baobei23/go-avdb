package api

import (
	"math"
	"net/http"
	"strconv"

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
//	@Param			sort	query		string	false	"Sort"
//	@Success		200		{object}	getVideoList
//	@Failure		500		{object}	error
//	@Router			/video [get]
func (app *Application) getVideoList(w http.ResponseWriter, r *http.Request) {

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	videos, total, err := app.Store.Video.GetList(r.Context(), limit, offset, search)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, getVideoList{
		Data:      videos,
		Total:     total,
		Page:      page,
		PageCount: int(math.Ceil(float64(total) / float64(limit))),
		Limit:     limit,
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}
