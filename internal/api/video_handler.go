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
	MetaData  string            `json:"metadata,omitempty"`
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

// getVideoListByActor godoc
//
//	@Summary		Get video list by actor
//	@Description	Get video list by actor
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			actor	path		string	true	"Actor"
//	@Param			page	query		int		false	"Page"
//	@Param			limit	query		int		false	"Limit"
//	@Success		200		{object}	getVideoList
//	@Failure		500		{object}	error
//	@Router			/video/actor/{actor} [get]
func (app *Application) getVideoListByActor(w http.ResponseWriter, r *http.Request) {
	actor := chi.URLParam(r, "actor")
	pq := store.PaginationQuery{}
	pq, err := pq.Parse(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	videos, total, err := app.Store.Video.GetListByActor(r.Context(), actor, pq)
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
		MetaData:  actor,
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}

// getVideoListByDirector godoc
//
//	@Summary		Get video list by director
//	@Description	Get video list by director
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			director	path		string	true	"Director"
//	@Param			page		query		int		false	"Page"
//	@Param			limit		query		int		false	"Limit"
//	@Success		200			{object}	getVideoList
//	@Failure		500			{object}	error
//	@Router			/video/director/{director} [get]
func (app *Application) getVideoListByDirector(w http.ResponseWriter, r *http.Request) {
	director := chi.URLParam(r, "director")
	pq := store.PaginationQuery{}
	pq, err := pq.Parse(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	videos, total, err := app.Store.Video.GetListByDirector(r.Context(), director, pq)
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
		MetaData:  director,
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}

// getVideoListByStudio godoc
//
//	@Summary		Get video list by studio
//	@Description	Get video list by studio
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			studio	path		string	true	"Studio"
//	@Param			page	query		int		false	"Page"
//	@Param			limit	query		int		false	"Limit"
//	@Success		200		{object}	getVideoList
//	@Failure		500		{object}	error
//	@Router			/video/studio/{studio} [get]
func (app *Application) getVideoListByStudio(w http.ResponseWriter, r *http.Request) {
	studio := chi.URLParam(r, "studio")
	pq := store.PaginationQuery{}
	pq, err := pq.Parse(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	videos, total, err := app.Store.Video.GetListByStudio(r.Context(), studio, pq)
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
		MetaData:  studio,
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}
