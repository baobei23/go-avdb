package api

import (
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

type Config struct {
	Port    string
	ApiURL  string
	DB      DBConfig
	Crawler crawler.Config
}

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}
