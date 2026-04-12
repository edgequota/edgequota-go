// Package cache provides SDK helpers for EdgeQuota's admin cache purge API.
//
// Usage:
//
//	c, _ := cache.NewClient("http://edgequota-admin:9090")
//
//	// Purge response cache entries by surrogate-key tags:
//	_ = c.PurgeTags(ctx, "menu-demo", "published-menus")
//
//	// Purge a single cached response by URL:
//	_ = c.PurgeURL(ctx, "/v1/images/tenant/shoro/image/01KNJ7182C90DH3K1KT1VB3HZB")
//
//	// Purge auth cache entries by tags:
//	_ = c.PurgeAuthTags(ctx, "table-t123")
package cache

import (
	"context"
	"fmt"
	"net/http"
	"time"

	admin "github.com/edgequota/edgequota-go/gen/http/admin/v1"
)

// Client wraps the generated admin API client with convenience methods.
type Client struct {
	inner *admin.ClientWithResponses
}

// Option configures a Client.
type Option func(*options)

type options struct {
	httpClient *http.Client
}

// WithHTTPClient sets a custom http.Client for requests.
func WithHTTPClient(c *http.Client) Option {
	return func(o *options) { o.httpClient = c }
}

// NewClient creates a cache purge client targeting the EdgeQuota admin API.
// adminURL is the base URL of the admin server (e.g. "http://edgequota-admin:9090").
func NewClient(adminURL string, opts ...Option) (*Client, error) {
	o := &options{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
	for _, opt := range opts {
		opt(o)
	}
	inner, err := admin.NewClientWithResponses(adminURL, admin.WithHTTPClient(o.httpClient))
	if err != nil {
		return nil, fmt.Errorf("edgequota cache: create client: %w", err)
	}
	return &Client{inner: inner}, nil
}

// PurgeTags invalidates all cached HTTP responses tagged with any of the given
// surrogate-key tags.
func (c *Client) PurgeTags(ctx context.Context, tags ...string) error {
	resp, err := c.inner.PurgeResponseCacheTagsWithResponse(ctx, admin.PurgeTagsRequest{Tags: tags})
	if err != nil {
		return fmt.Errorf("edgequota cache: purge tags: %w", err)
	}
	if resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("edgequota cache: purge tags: unexpected status %d", resp.StatusCode())
	}
	return nil
}

// PurgeURL invalidates a single cached GET response by URL path.
func (c *Client) PurgeURL(ctx context.Context, urlPath string) error {
	return c.PurgeURLWithMethod(ctx, http.MethodGet, urlPath)
}

// PurgeURLWithMethod invalidates a single cached response by HTTP method and URL path.
func (c *Client) PurgeURLWithMethod(ctx context.Context, method, urlPath string) error {
	resp, err := c.inner.PurgeResponseCacheURLWithResponse(ctx, admin.PurgeURLRequest{Url: urlPath, Method: &method})
	if err != nil {
		return fmt.Errorf("edgequota cache: purge url: %w", err)
	}
	if resp.StatusCode() != http.StatusNoContent && resp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("edgequota cache: purge url: unexpected status %d", resp.StatusCode())
	}
	return nil
}

// PurgeAuthTags invalidates all cached auth decisions tagged with any of the
// given tags. Tags correspond to cache_tags values returned by the auth service
// in CheckResponse.
func (c *Client) PurgeAuthTags(ctx context.Context, tags ...string) error {
	resp, err := c.inner.PurgeAuthCacheTagsWithResponse(ctx, admin.PurgeTagsRequest{Tags: tags})
	if err != nil {
		return fmt.Errorf("edgequota cache: purge auth tags: %w", err)
	}
	if resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("edgequota cache: purge auth tags: unexpected status %d", resp.StatusCode())
	}
	return nil
}
