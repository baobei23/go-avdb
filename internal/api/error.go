package api

import (
	"net/http"

	"go.uber.org/zap"
)

// 500
func (app *Application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Error("internal error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

// 403
func (app *Application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.Logger.Warn("forbidden", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	writeJSONError(w, http.StatusForbidden, "forbidden")
}

// 400
func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Warn("bad request", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

// 409
func (app *Application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Error("conflict response", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))

	writeJSONError(w, http.StatusConflict, err.Error())
}

// 404
func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Warn("not found error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))

	writeJSONError(w, http.StatusNotFound, "not found")
}

// 401
func (app *Application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Warn("unauthorized error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

// 401 Basic
func (app *Application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Warn("unauthorized basic error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

// 429
func (app *Application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.Logger.Warn("rate limit exceeded", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	w.Header().Set("Retry-After", retryAfter)

	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
