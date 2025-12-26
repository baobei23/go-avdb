package crawler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/baobei23/go-avdb/internal/store"
	"github.com/baobei23/go-avdb/internal/util"
)

type Crawler interface {
	CrawlPage(ctx context.Context, page int) error
	CrawlRange(ctx context.Context, start, end int) error
	CrawlAll(ctx context.Context) error
}

// Config for crawler
type Config struct {
	BaseURLProvide  string
	BaseURLProvide1 string
	Timeout         time.Duration
	MaxRetries      int
	PageDelay       time.Duration
}

// service is crawler service with single responsibility
type service struct {
	config     Config
	httpClient HTTPClient
	store      store.Storage
}

func NewService(config Config, store store.Storage) *service {
	return &service{
		config:     config,
		httpClient: NewHTTPClient(config.Timeout, config.MaxRetries),
		store:      store,
	}
}

// CrawlPage crawls a specific page from both endpoints
func (s *service) CrawlPage(ctx context.Context, page int) error {
	// 1. Crawl Provide (Base Data) first to ensure videos exist
	if err := s.crawlProvide(ctx, page); err != nil {
		return fmt.Errorf("provide error: %w", err)
	}

	// 2. Crawl Provide1 (Supplemental Data - Studio)
	if err := s.crawlProvide1(ctx, page); err != nil {
		return fmt.Errorf("provide1 error: %w", err)
	}

	return nil
}

// CrawlRange crawls a range of pages
func (s *service) CrawlRange(ctx context.Context, start, end int) error {
	step := 1
	if start > end {
		step = -1
	}

	for i := start; i != end+step; i += step {
		if err := s.CrawlPage(ctx, i); err != nil {
			log.Printf("Error crawling page %d: %v", i, err)
		} else {
			log.Printf("Crawled page %d", i)
		}

		if i != end {
			time.Sleep(s.config.PageDelay)
		}
	}
	return nil
}

// CrawlAll crawls all pages
func (s *service) CrawlAll(ctx context.Context) error {
	// Fetch page 1 to get total pages
	url := fmt.Sprintf("%s?ac=detail&pg=1", s.config.BaseURLProvide)
	resp, err := s.httpClient.FetchProvide(url)
	if err != nil {
		return err
	}

	pageCount := 1
	switch v := resp.PageCount.(type) {
	case float64:
		pageCount = int(v)
	case int:
		pageCount = v
	case string:
		pageCount, _ = strconv.Atoi(v)
	}

	log.Printf("Total pages to crawl: %d", pageCount)
	return s.CrawlRange(ctx, pageCount, 1)
}

func (s *service) crawlProvide(ctx context.Context, page int) error {
	url := fmt.Sprintf("%s?ac=detail&pg=%d", s.config.BaseURLProvide, page)

	resp, err := s.httpClient.FetchProvide(url)
	if err != nil {
		return err
	}

	for _, item := range resp.List {
		if err := s.processVideoItem(ctx, item); err != nil {
			log.Printf("Error processing video %d: %v", item.ID, err)
		}
	}
	return nil
}

func (s *service) crawlProvide1(ctx context.Context, page int) error {
	url := fmt.Sprintf("%s?ac=detail&pg=%d", s.config.BaseURLProvide1, page)

	resp, err := s.httpClient.FetchProvide1(url)
	if err != nil {
		return err
	}

	for _, item := range resp.List {
		if err := s.processVideoItemProvide1(ctx, item); err != nil {
			log.Printf("Error processing provide1 video %d: %v", item.VodID, err)
		}
	}
	return nil
}

func (s *service) processVideoItem(ctx context.Context, item VideoItem) error {
	createdAt := util.ParseDate(item.CreatedAt)

	video := store.Video{
		ID:          item.ID,
		Category:    item.TypeName,
		Name:        item.Name,
		Slug:        item.Slug,
		OriginName:  item.OriginName,
		PosterURL:   item.PosterURL,
		ThumbURL:    item.ThumbURL,
		Description: item.Description,
		LinkEmbed:   item.Episodes.ServerData["Full"].LinkEmbed,
		CreatedAt:   createdAt,
	}

	if err := s.store.Video.Upsert(ctx, &video); err != nil {
		return err
	}

	if len(item.Actor) > 0 {
		_ = s.store.Video.UpsertActor(ctx, video.ID, trimStrings(item.Actor))
	}
	if len(item.Category) > 0 {
		_ = s.store.Video.UpsertTag(ctx, video.ID, trimStrings(item.Category))
	}
	if len(item.Director) > 0 {
		_ = s.store.Video.UpsertDirector(ctx, video.ID, trimStrings(item.Director))
	}

	return nil
}

func (s *service) processVideoItemProvide1(ctx context.Context, item VideoItemProvide1) error {
	if item.VodWriter != "" {
		return s.store.Video.UpsertStudio(ctx, item.VodID, item.VodWriter)
	}
	return nil
}

func trimStrings(s []string) []string {
	var res []string
	for _, v := range s {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}
