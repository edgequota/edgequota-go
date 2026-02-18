package ratelimit

import (
	rlv1http "github.com/edgequota/edgequota-go/gen/http/ratelimit/v1"
)

// TenantLimits holds rate limit parameters for a tenant.
// This is a convenience type for template implementations; the
// wire type is rlv1http.GetLimitsResponse.
type TenantLimits struct {
	Average int64  `json:"average"`
	Burst   int64  `json:"burst"`
	Period  string `json:"period"`
}

// NewResponse creates a GetLimitsResponse from tenant limits with the
// given tenant key for Redis bucket isolation.
func NewResponse(tenantKey string, limits TenantLimits) rlv1http.GetLimitsResponse {
	resp := rlv1http.GetLimitsResponse{
		Average: limits.Average,
		Burst:   limits.Burst,
		Period:  limits.Period,
	}
	if tenantKey != "" {
		resp.TenantKey = &tenantKey
	}
	return resp
}

// WithCache returns a copy of the response with cache control set.
func WithCache(resp rlv1http.GetLimitsResponse, maxAgeSec int64) rlv1http.GetLimitsResponse {
	resp.CacheMaxAgeSeconds = &maxAgeSec
	return resp
}

// WithNoStore returns a copy of the response with caching disabled.
func WithNoStore(resp rlv1http.GetLimitsResponse) rlv1http.GetLimitsResponse {
	noStore := true
	resp.CacheNoStore = &noStore
	return resp
}
