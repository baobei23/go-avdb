package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *Application) crawlHandler(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	result, err := app.Crawler.CrawlPage(r.Context(), pageInt)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, result)
}
