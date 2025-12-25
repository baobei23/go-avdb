package main

import (
	"expvar"
	"log"
	"runtime"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/baobei23/go-avdb/internal/api"
	"github.com/baobei23/go-avdb/internal/crawler"
	"github.com/baobei23/go-avdb/internal/db"
	"github.com/baobei23/go-avdb/internal/env"

	"github.com/baobei23/go-avdb/internal/store"
	"go.uber.org/zap"
)

const version = "1.0.0"

func main() {
	cfg := api.Config{
		Port:   env.GetString("PORT", ":8080"),
		ApiURL: env.GetString("API_URL", "localhost:8080"),
		DB: api.DBConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/avdb?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		Crawler: crawler.Config{
			BaseURLProvide:  env.GetString("BASE_URL_PROVIDE", "https://avdbapi.com/api.php/provide/vod/"),
			BaseURLProvide1: env.GetString("BASE_URL_PROVIDE1", "https://avdbapi.com/api.php/provide1/vod/"),
			Timeout:         time.Duration(env.GetInt("CRAWLER_TIMEOUT", 30)) * time.Second,
			MaxRetries:      env.GetInt("CRAWLER_MAX_RETRIES", 3),
			PageDelay:       time.Duration(env.GetInt("CRAWLER_PAGE_DELAY", 2)) * time.Second,
		},
	}
	// logger
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	// database
	db, err := db.New(
		cfg.DB.Addr,
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("database connection pool established")

	// storage
	storage := store.NewStorage(db)

	// crawler
	crawlerService := crawler.NewService(cfg.Crawler, storage)

	// application
	app := &api.Application{
		Config:  cfg,
		Store:   storage,
		Logger:  logger,
		Crawler: crawlerService,
	}

	// metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.Mount()

	log.Fatal(app.Run(mux))
}
