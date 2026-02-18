package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	authv1http "github.com/edgequota/edgequota-go/gen/http/auth/v1"
)

// Allow returns a CheckResponse that allows the request, optionally
// injecting headers into the upstream request.
func Allow(requestHeaders map[string]string) authv1http.CheckResponse {
	resp := authv1http.CheckResponse{
		Allowed:    true,
		StatusCode: 200,
	}
	if len(requestHeaders) > 0 {
		resp.RequestHeaders = &requestHeaders
	}
	return resp
}

// Deny returns a CheckResponse that denies the request.
func Deny(statusCode int, body string, responseHeaders map[string]string) authv1http.CheckResponse {
	resp := authv1http.CheckResponse{
		Allowed:    false,
		StatusCode: int32(statusCode),
	}
	if body != "" {
		resp.DenyBody = &body
	}
	if len(responseHeaders) > 0 {
		resp.ResponseHeaders = &responseHeaders
	}
	return resp
}

// ExtractBearerToken extracts a Bearer token from the Authorization header
// in the CheckRequest.Headers map. Returns empty string if not found.
func ExtractBearerToken(req *authv1http.CheckRequest) string {
	auth := req.Headers["Authorization"]
	if auth == "" {
		auth = req.Headers["authorization"]
	}
	if !strings.HasPrefix(auth, "Bearer ") && !strings.HasPrefix(auth, "bearer ") {
		return ""
	}
	return auth[7:]
}

// JWTValidator validates HMAC-signed JWTs and extracts claims.
type JWTValidator struct {
	secret []byte
}

// NewJWTValidator creates a validator with the given HMAC secret.
func NewJWTValidator(secret string) *JWTValidator {
	return &JWTValidator{secret: []byte(secret)}
}

// ValidateToken parses and validates a JWT, returning its claims.
func (v *JWTValidator) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return v.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

// CreateToken creates a signed JWT with the given claims and expiry.
func (v *JWTValidator) CreateToken(claims map[string]interface{}, expiry time.Duration) (string, error) {
	mc := jwt.MapClaims{}
	for k, val := range claims {
		mc[k] = val
	}
	mc["exp"] = time.Now().Add(expiry).Unix()
	mc["iat"] = time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mc)
	return token.SignedString(v.secret)
}
