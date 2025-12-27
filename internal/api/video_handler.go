package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

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
//	@Description	Get video list with pagination and search
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int		false	"Page number"
//	@Param			limit	query		int		false	"Limit per page"
//	@Param			search	query		string	false	"Search by name"
//	@Success		200		{object}	struct{Data []store.Video; Meta map[string]any}
//	@Failure		500		{string}	error
//	@Router			/video [get]
func (app *Application) getVideoList(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10
	search := r.URL.Query().Get("search")

	if p := r.URL.Query().Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	offset := (page - 1) * limit

	videos, total, err := app.Store.Video.GetVideo(r.Context(), limit, offset, search)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	res := struct {
		Data any `json:"data"`
		Meta any `json:"meta"`
	}{
		Data: videos,
		Meta: map[string]any{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	}

	if err := writeJSON(w, http.StatusOK, res); err != nil {
		app.internalServerError(w, r, err)
	}
}

// getVideoByActor godoc
//
//	@Summary		Get videos by actor
//	@Description	Get videos by actor
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			actor	path		string	true	"Actor name"
//	@Success		200		{object}	[]store.Video
//	@Failure		400		{string}	error
//	@Failure		500		{string}	error
//	@Router			/video/actor/{actor} [get]
// func (app *Application) getVideoByActor(w http.ResponseWriter, r *http.Request) {
// 	actor, err := app.Store.Video.GetByActor(r.Context(), chi.URLParam(r, "actor"))
// 	if err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}

// 	if err := app.jsonResponse(w, http.StatusOK, actor); err != nil {
// 		app.internalServerError(w, r, err)
// 	}
// }
