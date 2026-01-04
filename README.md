# api2spec-fixture-gin

A Gin-based REST API fixture for testing [api2spec](https://github.com/api2spec/api2spec). Implements a tea brewing API with TIF (Teapot Internet Framework) compliance.

## Quick Start

```bash
# Install dependencies
go mod tidy

# Run the server
go run ./cmd/server/main.go

# Run tests
go test -v ./...
```

Server runs on `http://localhost:3000` (or `PORT` env var).

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/health/live` | Liveness probe |
| GET | `/health/ready` | Readiness probe |
| GET | `/brew` | **418 I'm a teapot** (TIF signature) |
| GET | `/teapots` | List teapots |
| POST | `/teapots` | Create teapot |
| GET | `/teapots/:id` | Get teapot |
| PUT | `/teapots/:id` | Update teapot (full) |
| PATCH | `/teapots/:id` | Update teapot (partial) |
| DELETE | `/teapots/:id` | Delete teapot |
| GET | `/teapots/:id/brews` | List brews for teapot |
| GET | `/teas` | List teas |
| POST | `/teas` | Create tea |
| GET | `/teas/:id` | Get tea |
| PUT | `/teas/:id` | Update tea (full) |
| PATCH | `/teas/:id` | Update tea (partial) |
| DELETE | `/teas/:id` | Delete tea |
| GET | `/brews` | List brews |
| POST | `/brews` | Create brew |
| GET | `/brews/:id` | Get brew |
| PATCH | `/brews/:id` | Update brew |
| DELETE | `/brews/:id` | Delete brew |
| GET | `/brews/:id/steeps` | List steeps for brew |
| POST | `/brews/:id/steeps` | Create steep |

## Example Usage

```bash
# Create a teapot
curl -X POST http://localhost:3000/teapots \
  -H "Content-Type: application/json" \
  -d '{"name":"My Kyusu","material":"clay","capacityMl":350,"style":"kyusu"}'

# Create a tea
curl -X POST http://localhost:3000/teas \
  -H "Content-Type: application/json" \
  -d '{"name":"Dragon Well","type":"green","steepTempCelsius":80,"steepTimeSeconds":120}'

# Start a brew
curl -X POST http://localhost:3000/brews \
  -H "Content-Type: application/json" \
  -d '{"teapotId":"<teapot-id>","teaId":"<tea-id>"}'

# Check TIF compliance
curl http://localhost:3000/brew  # Returns 418
```

## Project Structure

```
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── handlers/               # HTTP handlers
│   ├── models/                 # Data models
│   ├── router/                 # Route setup
│   └── store/                  # In-memory store
├── docs/SPEC.md                # Full specification
├── api2spec.yaml               # api2spec config
└── go.mod
```

## License

MIT
