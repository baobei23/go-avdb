package api

import (
	"context"
	"errors"
	"expvar"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baobei23/go-avdb/internal/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *Application) Mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5174")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	if app.Config.RateLimiter.Enabled {
		r.Use(app.RateLimiterMiddleware)
	}

	//health
	r.Get("/", app.health)

	//metrics
	r.Get("/metrics", expvar.Handler().ServeHTTP)

	//swagger
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("doc.json")))

	//crawler
	r.Route("/crawl", func(r chi.Router) {
		r.Get("/{page}", app.crawlPage)
		r.Get("/all", app.crawlAll)
		r.Get("/range", app.crawlRange)
	})

	//video
	r.Route("/video", func(r chi.Router) {
		r.Get("/", app.getVideoList)
		r.Get("/{slug}", app.getVideo)
		r.Get("/actor/{actor}", app.getVideoListByActor)
		r.Get("/director/{director}", app.getVideoListByDirector)
		r.Get("/studio/{studio}", app.getVideoListByStudio)
		r.Get("/tag/{tag}", app.getVideoListByTag)
	})

	// Actor
	r.Route("/actor", func(r chi.Router) {
		r.Post("/", app.createActor)
		r.Get("/", app.getActorList)
		r.Get("/{name}", app.health)
		r.Put("/{id}", app.updateActor)
		r.Delete("/{id}", app.deleteActor)
		//TODO: merging router
	})

	// Director
	r.Route("/director", func(r chi.Router) {
		r.Post("/", app.health)
		r.Get("/", app.health)
		r.Get("/{name}", app.health)
		r.Put("/{id}", app.health)
	})

	// Tag
	r.Route("/tag", func(r chi.Router) {
		r.Post("/", app.health)
		r.Get("/", app.health)
		r.Get("/{name}", app.health)
		r.Put("/{id}", app.health)
	})

	// Studio
	r.Route("/studio", func(r chi.Router) {
		r.Post("/", app.health)
		r.Get("/", app.health)
		r.Get("/{name}", app.health)
		r.Put("/{id}", app.health)
	})

	return r
}

func (app *Application) Run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.Config.Port,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.Logger.Info("signal caught")

		shutdown <- srv.Shutdown(ctx)
	}()

	app.Logger.Info("Server started")
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.Logger.Info("server has stopped")

	return nil
}
