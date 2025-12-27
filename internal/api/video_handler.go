package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// getVideo godoc
//
//	@Summary		Get video detail by slug
//	@Description	Get video detail by slug
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string	true	"Video slug"
//	@Success		200		{object}	store.Video
//	@Failure		500		{string}	error
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

func (app *Application) getVideoList(w http.ResponseWriter, r *http.Request) {

}
