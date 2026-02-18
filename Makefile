.PHONY: all generate generate-grpc generate-http build test lint clean

OPENAPI_BASE ?= https://raw.githubusercontent.com/edgequota/edgequota/refs/heads/main/api/openapi

all: generate build test

# ── Code generation ──────────────────────────────────────────────────

generate: generate-grpc generate-http

generate-grpc:
	buf generate

generate-http:
	oapi-codegen --config oapi-codegen-auth.yaml $(OPENAPI_BASE)/auth/v1/auth.yaml
	oapi-codegen --config oapi-codegen-ratelimit.yaml $(OPENAPI_BASE)/ratelimit/v1/ratelimit.yaml
	oapi-codegen --config oapi-codegen-events.yaml $(OPENAPI_BASE)/events/v1/events.yaml

# ── Build ────────────────────────────────────────────────────────────

build:
	go build ./...

# ── Test ─────────────────────────────────────────────────────────────

test:
	go test ./...

test-race:
	go test -race ./...

# ── Lint ─────────────────────────────────────────────────────────────

lint:
	go vet ./...

# ── Tidy ─────────────────────────────────────────────────────────────

tidy:
	go mod tidy

# ── Clean ────────────────────────────────────────────────────────────

clean:
	rm -rf gen/
