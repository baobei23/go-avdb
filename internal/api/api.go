package api

import (
	"time"

	"github.com/baobei23/go-avdb/internal/crawler"
	"github.com/baobei23/go-avdb/internal/ratelimiter"
	"github.com/baobei23/go-avdb/internal/store"
	"github.com/baobei23/go-avdb/internal/store/cache"
	"go.uber.org/zap"
)

type Application struct {
	Config       Config
	Store        store.Storage
	CacheStorage cache.Storage
	Crawler      crawler.Crawler
	Logger       *zap.Logger
	RateLimiter  ratelimiter.Limiter
}

type Config struct {
	Port        string
	ApiURL      string
	DB          DBConfig
	Crawler     crawler.Config
	Env         string
	ApiVersion  string
	RedisCfg    RedisConfig
	RateLimiter ratelimiter.Config
	Auth        AuthConfig
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Enabled  bool
}

type DBConfig struct {
	Addr            string
	MaxOpenConns    int
	MaxIdleConns    int
	MaxIdleTime     time.Duration
	MaxConnLifetime time.Duration
}

type AuthConfig struct {
	User string
	Pass string
}
