package api

import (
	"net/http"

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

func (app *Application) getVideoList(w http.ResponseWriter, r *http.Request) {
	videos, err := app.Store.Video.GetList(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, videos); err != nil {
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
