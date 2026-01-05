package main

import (
	"expvar"
	"log"
	"runtime"
	"time"

	_ "github.com/baobei23/go-avdb/docs"
	"github.com/baobei23/go-avdb/internal/api"
	"github.com/baobei23/go-avdb/internal/crawler"
	"github.com/baobei23/go-avdb/internal/db"
	"github.com/baobei23/go-avdb/internal/env"
	"github.com/baobei23/go-avdb/internal/ratelimiter"
	"github.com/baobei23/go-avdb/internal/store/cache"
	"github.com/redis/go-redis/v9"

	"github.com/baobei23/go-avdb/internal/store"
	"go.uber.org/zap"
)

//	@title			Go - AVDB
//	@description	AVDB crawler and API
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath	/
func main() {
	cfg := api.Config{
		Port:   env.GetString("PORT", ":8080"),
		ApiURL: env.GetString("API_URL", "http://localhost:8080"),
		DB: api.DBConfig{
			Addr:            env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/avdb?sslmode=disable"),
			MaxOpenConns:    env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns:    env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:     time.Duration(env.GetInt("DB_MAX_IDLE_TIME", 15)) * time.Minute,
			MaxConnLifetime: time.Duration(env.GetInt("DB_MAX_CONN_LIFETIME", 1)) * time.Hour,
		},
		Crawler: crawler.Config{
			BaseURLProvide:  env.GetString("BASE_URL_PROVIDE", "https://avdbapi.com/api.php/provide/vod/"),
			BaseURLProvide1: env.GetString("BASE_URL_PROVIDE1", "https://avdbapi.com/api.php/provide1/vod/"),
			Timeout:         time.Duration(env.GetInt("CRAWLER_TIMEOUT", 30)) * time.Second,
			MaxRetries:      env.GetInt("CRAWLER_MAX_RETRIES", 3),
			PageDelay:       time.Duration(env.GetInt("CRAWLER_PAGE_DELAY", 2)) * time.Second,
			WorkerCount:     env.GetInt("CRAWLER_WORKER_COUNT", 3),
		},
		Env:        env.GetString("ENV", "development"),
		ApiVersion: "",
		RedisCfg: api.RedisConfig{
			Addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			Password: env.GetString("REDIS_PASSWORD", ""),
			DB:       env.GetInt("REDIS_DB", 0),
			Enabled:  env.GetBool("REDIS_ENABLED", false),
		},
		RateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATE_LIMITER_REQUESTS_PER_TIME_FRAME", 1000),
			TimeFrame:            time.Duration(env.GetInt("RATE_LIMITER_TIME_FRAME", 1)) * time.Minute,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", false),
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
		cfg.DB.MaxConnLifetime,
	)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("database connection pool established")

	// cache
	var redis *redis.Client
	if cfg.RedisCfg.Enabled {
		redis = cache.NewRedisClient(cfg.RedisCfg.Addr, cfg.RedisCfg.Password, cfg.RedisCfg.DB)
		logger.Info("redis connection pool established")
	}

	// rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.RateLimiter.RequestsPerTimeFrame,
		cfg.RateLimiter.TimeFrame,
	)

	// storage
	storage := store.NewStorage(db)

	// redis
	cacheStorage := cache.NewRedisStorage(redis)

	// crawler
	crawlerService := crawler.NewService(cfg.Crawler, storage)

	// application
	app := &api.Application{
		Config:       cfg,
		Store:        storage,
		CacheStorage: cacheStorage,
		Logger:       logger,
		Crawler:      crawlerService,
		RateLimiter:  rateLimiter,
	}

	// metrics collected
	expvar.NewString("version").Set(cfg.ApiVersion)
	expvar.Publish("database", expvar.Func(func() any {
		stats := db.Stat()
		return map[string]interface{}{
			"total_conns":        stats.TotalConns(),
			"idle_conns":         stats.IdleConns(),
			"acquired_conns":     stats.AcquiredConns(),
			"constructing_conns": stats.ConstructingConns(),
			"max_conns":          stats.MaxConns(),
		}
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.Mount()

	log.Fatal(app.Run(mux))
}
