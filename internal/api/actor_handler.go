package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/baobei23/go-avdb/internal/store"
	"github.com/go-chi/chi/v5"
)

type actorPayload struct {
	Name string `json:"name"`
}

// getActorList godoc
//
//	@Summary		Get actor list
//	@Description	Get actor list
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]store.Actor
//	@Failure		500	{string}	error
//	@Router			/actor [get]
func (app *Application) getActorList(w http.ResponseWriter, r *http.Request) {
	actors, err := app.Store.Actor.GetList(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, actors); err != nil {
		app.internalServerError(w, r, err)
	}
}

// createActor godoc
//
//	@Summary		Create actor
//	@Description	Create actor
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			request	body		actorPayload	true	"Actor request"
//	@Success		201		{object}	store.Actor
//	@Failure		400		{string}	error
//	@Failure		500		{string}	error
//	@Router			/actor [post]
func (app *Application) createActor(w http.ResponseWriter, r *http.Request) {
	var payload actorPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Name == "" {
		app.badRequestResponse(w, r, fmt.Errorf("name is required"))
		return
	}

	actor := store.Actor{
		Name: payload.Name,
	}

	if err := app.Store.Actor.Create(r.Context(), &actor); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, actor); err != nil {
		app.internalServerError(w, r, err)
	}
}

// updateActor godoc
//
//	@Summary		Update actor
//	@Description	Update actor
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"Actor ID"
//	@Param			request	body		actorPayload	true	"Actor request"
//	@Success		200		{object}	store.Actor
//	@Failure		400		{string}	error
//	@Failure		404		{string}	error
//	@Failure		500		{string}	error
//	@Router			/actor/{id} [put]
func (app *Application) updateActor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var payload actorPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Name == "" {
		app.badRequestResponse(w, r, fmt.Errorf("name is required"))
		return
	}

	actor := store.Actor{
		ID:   id,
		Name: payload.Name,
	}

	if err := app.Store.Actor.Update(r.Context(), &actor); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		case errors.Is(err, store.ErrConflict):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, actor); err != nil {
		app.internalServerError(w, r, err)
	}
}

// deleteActor godoc
//
//	@Summary		Delete actor
//	@Description	Delete actor
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Actor ID"
//	@Success		204	{string}	string
//	@Failure		404	{string}	error
//	@Failure		500	{string}	error
//	@Router			/actor/{id} [delete]
func (app *Application) deleteActor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.Store.Actor.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
