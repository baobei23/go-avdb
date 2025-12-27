package api

import (
	"time"

	"github.com/baobei23/go-avdb/internal/crawler"
	"github.com/baobei23/go-avdb/internal/store"
	"go.uber.org/zap"
)

type Application struct {
	Config  Config
	Store   store.Storage
	Crawler crawler.Crawler
	Logger  *zap.Logger
}
type AppLogger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)

	// optional helpers
	With(fields ...zap.Field) AppLogger
	Sync() error
}

type Config struct {
	Port    string
	ApiURL  string
	DB      DBConfig
	Crawler crawler.Config
	Env     string
	Version string
}

type DBConfig struct {
	Addr            string
	MaxOpenConns    int
	MaxIdleConns    int
	MaxIdleTime     time.Duration
	MaxConnLifetime time.Duration
}
