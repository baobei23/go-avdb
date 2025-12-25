package crawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClient interface {
	Get(ctx context.Context, url string) ([]byte, error)
}

type DefaultHTTPClient struct {
	client     *http.Client
	maxRetries int
	retryDelay time.Duration
}

func NewHTTPClient(timeout time.Duration, maxRetries int) *DefaultHTTPClient {
	return &DefaultHTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		maxRetries: maxRetries,
		retryDelay: 2 * time.Second,
	}
}

func (c *DefaultHTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.retryDelay):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("unexpected status: %d", resp.StatusCode)
			continue
		}

		return body, nil
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", c.maxRetries, lastErr)
}
