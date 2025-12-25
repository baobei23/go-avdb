package crawler

import (
	"context"
	"time"

	"github.com/baobei23/go-avdb/internal/store"
)

type Crawler interface {
	CrawlPage(ctx context.Context, page int) (*CrawlResult, error)
	CrawlRange(ctx context.Context, start, end int) error
	CrawlAll(ctx context.Context) error
}

// Config untuk crawler
type Config struct {
	BaseURLProvide  string
	BaseURLProvide1 string
	Timeout         time.Duration
	MaxRetries      int
	PageDelay       time.Duration
}

// Service adalah crawler service dengan single responsibility
type Service struct {
	config     Config
	httpClient HTTPClient
	store      store.Storage
}

func NewService(config Config, store store.Storage) *Service {
	return &Service{
		config:     config,
		httpClient: NewHTTPClient(config.Timeout, config.MaxRetries),
		store:      store,
	}
}
