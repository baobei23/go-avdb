package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func (app *Application) crawlPage(w http.ResponseWriter, r *http.Request) {
	pageStr := chi.URLParam(r, "page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	go func() {
		ctx := context.Background()
		if err := app.Crawler.CrawlPage(ctx, page); err != nil {
			app.Logger.Error("crawl page failed", zap.Int("page", page), zap.Error(err))
		}
	}()

	if err := app.jsonResponse(w, http.StatusAccepted, map[string]string{
		"status": "crawling page started",
		"page":   pageStr,
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *Application) crawlAll(w http.ResponseWriter, r *http.Request) {
	go func() {
		ctx := context.Background()
		if err := app.Crawler.CrawlAll(ctx); err != nil {
			app.Logger.Error("crawl all failed", zap.Error(err))
		}
	}()

	if err := app.jsonResponse(w, http.StatusAccepted, map[string]string{
		"status": "crawling all started",
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *Application) crawlRange(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	start, err := strconv.Atoi(startStr)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	end, err := strconv.Atoi(endStr)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	go func() {
		ctx := context.Background()
		if err := app.Crawler.CrawlRange(ctx, start, end); err != nil {
			app.Logger.Error("crawl range failed", zap.Int("start", start), zap.Int("end", end), zap.Error(err))
		}
	}()

	if err := app.jsonResponse(w, http.StatusAccepted, map[string]string{
		"status": "crawling range started",
		"start":  startStr,
		"end":    endStr,
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}
