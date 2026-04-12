// Package cache provides a client for EdgeQuota's admin cache purge API.
//
// Usage:
//
//	c := cache.NewClient("http://edgequota-admin:9090")
//	_ = c.PurgeTags(ctx, "menu-demo", "published-menus")
//	_ = c.PurgeURL(ctx, "/v1/images/tenant/shoro/image/01KNJ7182C90DH3K1KT1VB3HZB")
package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client calls the EdgeQuota admin API to invalidate cached responses.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets a custom http.Client for requests.
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.httpClient = c }
}

// NewClient creates a cache purge client targeting the EdgeQuota admin API.
// adminURL is the base URL of the admin server (e.g. "http://edgequota-admin:9090").
func NewClient(adminURL string, opts ...Option) *Client {
	c := &Client{
		baseURL: strings.TrimRight(adminURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// PurgeTags invalidates all cached responses tagged with any of the given
// surrogate-key tags. Tags correspond to Surrogate-Key / Cache-Tag header
// values that backends emit (e.g. "menu-demo", "tenant-shoro").
func (c *Client) PurgeTags(ctx context.Context, tags ...string) error {
	body, err := json.Marshal(purgeTagsRequest{Tags: tags})
	if err != nil {
		return fmt.Errorf("edgequota cache: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/cache/purge/tags", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("edgequota cache: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("edgequota cache: purge tags: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("edgequota cache: purge tags: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// PurgeURL invalidates a single cached GET response by URL path (including
// query string). For example: "/v1/menus/menu-demo/published".
func (c *Client) PurgeURL(ctx context.Context, url string) error {
	return c.PurgeURLWithMethod(ctx, http.MethodGet, url)
}

// PurgeURLWithMethod invalidates a single cached response by HTTP method
// and URL path.
func (c *Client) PurgeURLWithMethod(ctx context.Context, method, url string) error {
	body, err := json.Marshal(purgeRequest{URL: url, Method: method})
	if err != nil {
		return fmt.Errorf("edgequota cache: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/cache/purge", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("edgequota cache: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("edgequota cache: purge url: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("edgequota cache: purge url: unexpected status %d", resp.StatusCode)
	}
	return nil
}

type purgeRequest struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

type purgeTagsRequest struct {
	Tags []string `json:"tags"`
}
