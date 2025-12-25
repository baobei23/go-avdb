package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClient interface {
	FetchProvide(url string) (*APIResponse, error)
	FetchProvide1(url string) (*APIResponseProvide1, error)
}

type httpClient struct {
	client *http.Client
}

func NewHTTPClient(timeout time.Duration, retries int) HTTPClient {
	return &httpClient{
		client: &http.Client{
			// Simple timeout, retries would need a custom Transport or loop in the caller
			Timeout: timeout,
		},
	}
}

func (c *httpClient) FetchProvide(url string) (*APIResponse, error) {
	var resp APIResponse
	if err := c.fetch(url, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *httpClient) FetchProvide1(url string) (*APIResponseProvide1, error) {
	var resp APIResponseProvide1
	if err := c.fetch(url, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *httpClient) fetch(url string, target any) error {
	res, err := c.client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}
	return nil
}
