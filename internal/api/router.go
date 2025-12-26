package api

import (
	"log"
	"net/http"
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

	//health
	r.Get("/", app.healthHandler)
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

	log.Printf("Starting server on %s", app.Config.Port)
	return srv.ListenAndServe()
}
