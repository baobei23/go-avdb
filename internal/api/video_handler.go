package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) getVideo(w http.ResponseWriter, r *http.Request) {
	video, err := app.Store.Video.GetVideoBySlug(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, video); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *Application) getVideoList(w http.ResponseWriter, r *http.Request) {
	videos, err := app.Store.Video.GetVideoList(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, videos); err != nil {
		app.internalServerError(w, r, err)
	}
}
