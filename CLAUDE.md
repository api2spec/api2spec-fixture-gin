# Project: api2spec-fixture-gin

Gin REST API fixture for testing api2spec OpenAPI generation. Tea brewing domain with TIF compliance.

## Tech Stack

- Go 1.22+ / Gin 1.10+
- In-memory storage (no database)
- go-playground/validator for validation
- testify for testing

## Commands

```bash
go run ./cmd/server/main.go    # Run server (port 3000)
go test -v ./...               # Run all tests
go build ./...                 # Build
go mod tidy                    # Tidy dependencies
```

## Project Structure

```
cmd/server/main.go             # Entry point
internal/handlers/*.go         # HTTP handlers (swag comments)
internal/models/*.go           # Request/response structs
internal/store/memory.go       # Thread-safe in-memory store
internal/router/router.go      # Route configuration
docs/SPEC.md                   # Full specification
```

## Key Patterns

- **Swag comments** on all handlers for api2spec to parse
- **Gin route params:** Use `:id` consistently (Gin doesn't allow `:id` and `:teapotId` on same path segment)
- **Validation:** `binding` tags on structs
- **Error responses:** Use `models.Error` struct
- **List endpoints:** Return `{data: [], pagination: {}}`

## Domain

```
Teapot → Brew → Steep
Tea → Brew
```

## TIF Compliance

`GET /brew` returns 418 "I'm a teapot"
