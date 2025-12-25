package api

import (
	"net/http"

	"go.uber.org/zap"
)

func (app *Application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Error("internal error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *Application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.Logger.Warn("forbidden", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	writeJSONError(w, http.StatusForbidden, "forbidden")
}

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Warn("bad request", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *Application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Error("conflict response", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))

	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Warn("not found error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))

	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *Application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Warn("unauthorized error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *Application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Warn("unauthorized basic error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *Application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.Logger.Warn("rate limit exceeded", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	w.Header().Set("Retry-After", retryAfter)

	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
