# Gin Fixture Specification

**Repository:** `api2spec-fixture-gin`  
**Purpose:** Target fixture (no native OpenAPI generation)

---

## Quick Reference

| Property | Value |
|----------|-------|
| Language | Go 1.22+ |
| Framework | Gin 1.10+ |
| Schema Library | Go structs + validator tags + swag comments |
| ORM | None (in-memory for fixture) |
| Test Runner | Go testing + testify |

---

## Project Setup

### Initialize

```bash
mkdir api2spec-fixture-gin
cd api2spec-fixture-gin
go mod init github.com/api2spec/api2spec-fixture-gin
go get github.com/gin-gonic/gin
go get github.com/go-playground/validator/v10
go get github.com/google/uuid
go get github.com/stretchr/testify
```

### go.mod

```go
module github.com/api2spec/api2spec-fixture-gin

go 1.22

require (
    github.com/gin-gonic/gin v1.10.0
    github.com/go-playground/validator/v10 v10.22.0
    github.com/google/uuid v1.6.0
    github.com/stretchr/testify v1.9.0
)
```

---

## Directory Structure

```
api2spec-fixture-gin/
├── docs/
│   └── SPEC.md                  # This file
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── handlers/
│   │   ├── teapots.go           # Teapot handlers
│   │   ├── teas.go              # Tea handlers
│   │   ├── brews.go             # Brew handlers
│   │   └── health.go            # Health + TIF 418
│   ├── models/
│   │   ├── teapot.go            # Teapot structs
│   │   ├── tea.go               # Tea structs
│   │   ├── brew.go              # Brew structs
│   │   ├── steep.go             # Steep structs
│   │   └── common.go            # Pagination, Error, etc.
│   ├── store/
│   │   └── memory.go            # In-memory data store
│   └── router/
│       └── router.go            # Route setup
├── expected/
│   └── openapi.yaml             # Expected api2spec output
├── api2spec.config.yaml         # api2spec configuration
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Models to Implement

### internal/models/common.go

```go
package models

import "time"

// PaginationQuery represents pagination query parameters
// @Description Pagination query parameters
type PaginationQuery struct {
    Page  int `form:"page" binding:"omitempty,min=1" default:"1"`
    Limit int `form:"limit" binding:"omitempty,min=1,max=100" default:"20"`
}

// Pagination represents pagination metadata in responses
// @Description Pagination metadata
type Pagination struct {
    Page       int `json:"page" example:"1"`
    Limit      int `json:"limit" example:"20"`
    Total      int `json:"total" example:"100"`
    TotalPages int `json:"totalPages" example:"5"`
}

// PaginatedResponse is a generic paginated response wrapper
type PaginatedResponse[T any] struct {
    Data       []T        `json:"data"`
    Pagination Pagination `json:"pagination"`
}

// Error represents an API error response
// @Description API error response
type Error struct {
    Code    string            `json:"code" example:"VALIDATION_ERROR"`
    Message string            `json:"message" example:"Invalid request body"`
    Details map[string]string `json:"details,omitempty"`
}

// HealthCheck represents a single health check result
// @Description Health check result
type HealthCheck struct {
    Name      string `json:"name" example:"database"`
    Status    string `json:"status" example:"ok" enums:"ok,degraded,down"`
    LatencyMs *int64 `json:"latencyMs,omitempty" example:"5"`
    Message   *string `json:"message,omitempty"`
}

// HealthResponse represents the health endpoint response
// @Description Health check response
type HealthResponse struct {
    Status    string        `json:"status" example:"ok" enums:"ok,degraded,down"`
    Timestamp time.Time     `json:"timestamp" example:"2025-01-04T12:00:00Z"`
    Version   *string       `json:"version,omitempty" example:"1.0.0"`
    Checks    []HealthCheck `json:"checks,omitempty"`
}

// TeapotResponse represents the TIF 418 response
// @Description TIF 418 I'm a teapot response
type TeapotResponse struct {
    Error   string `json:"error" example:"I'm a teapot"`
    Message string `json:"message" example:"This server is TIF-compliant"`
    Spec    string `json:"spec" example:"https://teapotframework.dev"`
}
```

### internal/models/teapot.go

```go
package models

import "time"

// TeapotMaterial represents valid teapot materials
// @Description Teapot material type
// @Enum ceramic,cast-iron,glass,porcelain,clay,stainless-steel
type TeapotMaterial string

const (
    MaterialCeramic       TeapotMaterial = "ceramic"
    MaterialCastIron      TeapotMaterial = "cast-iron"
    MaterialGlass         TeapotMaterial = "glass"
    MaterialPorcelain     TeapotMaterial = "porcelain"
    MaterialCite          TeapotMaterial = "clay"
    MaterialStainlessSteel TeapotMaterial = "stainless-steel"
)

// TeapotStyle represents valid teapot styles
// @Description Teapot style type
// @Enum kyusu,gaiwan,english,moroccan,turkish,yixing
type TeapotStyle string

const (
    StyleKyusu    TeapotStyle = "kyusu"
    StyleGaiwan   TeapotStyle = "gaiwan"
    StyleEnglish  TeapotStyle = "english"
    StyleMoroccan TeapotStyle = "moroccan"
    StyleTurkish  TeapotStyle = "turkish"
    StyleYixing   TeapotStyle = "yixing"
)

// Teapot represents a teapot entity
// @Description Teapot entity
type Teapot struct {
    ID          string         `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
    Name        string         `json:"name" example:"Classic English Teapot"`
    Material    TeapotMaterial `json:"material" example:"ceramic"`
    CapacityMl  int            `json:"capacityMl" example:"1200"`
    Style       TeapotStyle    `json:"style" example:"english"`
    Description *string        `json:"description" example:"A traditional English teapot"`
    CreatedAt   time.Time      `json:"createdAt" example:"2025-01-04T12:00:00Z"`
    UpdatedAt   time.Time      `json:"updatedAt" example:"2025-01-04T12:00:00Z"`
}

// CreateTeapotRequest represents the request body for creating a teapot
// @Description Create teapot request
type CreateTeapotRequest struct {
    Name        string         `json:"name" binding:"required,min=1,max=100" example:"My Kyusu"`
    Material    TeapotMaterial `json:"material" binding:"required,oneof=ceramic cast-iron glass porcelain clay stainless-steel" example:"clay"`
    CapacityMl  int            `json:"capacityMl" binding:"required,min=1,max=5000" example:"350"`
    Style       TeapotStyle    `json:"style" binding:"omitempty,oneof=kyusu gaiwan english moroccan turkish yixing" example:"kyusu"`
    Description *string        `json:"description" binding:"omitempty,max=500"`
}

// UpdateTeapotRequest represents the request body for PUT (full replacement)
// @Description Update teapot request (full replacement)
type UpdateTeapotRequest struct {
    Name        string         `json:"name" binding:"required,min=1,max=100" example:"Updated Teapot"`
    Material    TeapotMaterial `json:"material" binding:"required,oneof=ceramic cast-iron glass porcelain clay stainless-steel" example:"ceramic"`
    CapacityMl  int            `json:"capacityMl" binding:"required,min=1,max=5000" example:"1000"`
    Style       TeapotStyle    `json:"style" binding:"required,oneof=kyusu gaiwan english moroccan turkish yixing" example:"english"`
    Description *string        `json:"description" binding:"omitempty,max=500"`
}

// PatchTeapotRequest represents the request body for PATCH (partial update)
// @Description Patch teapot request (partial update)
type PatchTeapotRequest struct {
    Name        *string         `json:"name" binding:"omitempty,min=1,max=100"`
    Material    *TeapotMaterial `json:"material" binding:"omitempty,oneof=ceramic cast-iron glass porcelain clay stainless-steel"`
    CapacityMl  *int            `json:"capacityMl" binding:"omitempty,min=1,max=5000"`
    Style       *TeapotStyle    `json:"style" binding:"omitempty,oneof=kyusu gaiwan english moroccan turkish yixing"`
    Description *string         `json:"description" binding:"omitempty,max=500"`
}

// TeapotQuery represents query parameters for listing teapots
// @Description Teapot list query parameters
type TeapotQuery struct {
    PaginationQuery
    Material *TeapotMaterial `form:"material" binding:"omitempty,oneof=ceramic cast-iron glass porcelain clay stainless-steel"`
    Style    *TeapotStyle    `form:"style" binding:"omitempty,oneof=kyusu gaiwan english moroccan turkish yixing"`
}

// TeapotListResponse represents a paginated list of teapots
// @Description Paginated teapot list response
type TeapotListResponse struct {
    Data       []Teapot   `json:"data"`
    Pagination Pagination `json:"pagination"`
}
```

### internal/models/tea.go

```go
package models

import "time"

// TeaType represents valid tea types
// @Description Tea type
// @Enum green,black,oolong,white,puerh,herbal,rooibos
type TeaType string

const (
    TeaGreen   TeaType = "green"
    TeaBlack   TeaType = "black"
    TeaOolong  TeaType = "oolong"
    TeaWhite   TeaType = "white"
    TeaPuerh   TeaType = "puerh"
    TeaHerbal  TeaType = "herbal"
    TeaRooibos TeaType = "rooibos"
)

// CaffeineLevel represents caffeine content levels
// @Description Caffeine level
// @Enum none,low,medium,high
type CaffeineLevel string

const (
    CaffeineNone   CaffeineLevel = "none"
    CaffeineLow    CaffeineLevel = "low"
    CaffeineMedium CaffeineLevel = "medium"
    CaffeineHigh   CaffeineLevel = "high"
)

// Tea represents a tea entity
// @Description Tea entity
type Tea struct {
    ID               string        `json:"id" example:"550e8400-e29b-41d4-a716-446655440001"`
    Name             string        `json:"name" example:"Dragon Well Green Tea"`
    Type             TeaType       `json:"type" example:"green"`
    Origin           *string       `json:"origin,omitempty" example:"Hangzhou, China"`
    CaffeineLevel    CaffeineLevel `json:"caffeineLevel" example:"medium"`
    SteepTempCelsius int           `json:"steepTempCelsius" example:"80"`
    SteepTimeSeconds int           `json:"steepTimeSeconds" example:"180"`
    Description      *string       `json:"description,omitempty" example:"A famous Chinese green tea"`
    CreatedAt        time.Time     `json:"createdAt" example:"2025-01-04T12:00:00Z"`
    UpdatedAt        time.Time     `json:"updatedAt" example:"2025-01-04T12:00:00Z"`
}

// CreateTeaRequest represents the request body for creating a tea
// @Description Create tea request
type CreateTeaRequest struct {
    Name             string        `json:"name" binding:"required,min=1,max=100" example:"Earl Grey"`
    Type             TeaType       `json:"type" binding:"required,oneof=green black oolong white puerh herbal rooibos" example:"black"`
    Origin           *string       `json:"origin" binding:"omitempty,max=100" example:"England"`
    CaffeineLevel    CaffeineLevel `json:"caffeineLevel" binding:"omitempty,oneof=none low medium high" example:"high"`
    SteepTempCelsius int           `json:"steepTempCelsius" binding:"required,min=60,max=100" example:"95"`
    SteepTimeSeconds int           `json:"steepTimeSeconds" binding:"required,min=1,max=600" example:"240"`
    Description      *string       `json:"description" binding:"omitempty,max=1000"`
}

// UpdateTeaRequest represents the request body for PUT (full replacement)
// @Description Update tea request (full replacement)
type UpdateTeaRequest struct {
    Name             string        `json:"name" binding:"required,min=1,max=100"`
    Type             TeaType       `json:"type" binding:"required,oneof=green black oolong white puerh herbal rooibos"`
    Origin           *string       `json:"origin" binding:"omitempty,max=100"`
    CaffeineLevel    CaffeineLevel `json:"caffeineLevel" binding:"required,oneof=none low medium high"`
    SteepTempCelsius int           `json:"steepTempCelsius" binding:"required,min=60,max=100"`
    SteepTimeSeconds int           `json:"steepTimeSeconds" binding:"required,min=1,max=600"`
    Description      *string       `json:"description" binding:"omitempty,max=1000"`
}

// PatchTeaRequest represents the request body for PATCH (partial update)
// @Description Patch tea request (partial update)
type PatchTeaRequest struct {
    Name             *string        `json:"name" binding:"omitempty,min=1,max=100"`
    Type             *TeaType       `json:"type" binding:"omitempty,oneof=green black oolong white puerh herbal rooibos"`
    Origin           *string        `json:"origin" binding:"omitempty,max=100"`
    CaffeineLevel    *CaffeineLevel `json:"caffeineLevel" binding:"omitempty,oneof=none low medium high"`
    SteepTempCelsius *int           `json:"steepTempCelsius" binding:"omitempty,min=60,max=100"`
    SteepTimeSeconds *int           `json:"steepTimeSeconds" binding:"omitempty,min=1,max=600"`
    Description      *string        `json:"description" binding:"omitempty,max=1000"`
}

// TeaQuery represents query parameters for listing teas
// @Description Tea list query parameters
type TeaQuery struct {
    PaginationQuery
    Type          *TeaType       `form:"type" binding:"omitempty,oneof=green black oolong white puerh herbal rooibos"`
    CaffeineLevel *CaffeineLevel `form:"caffeineLevel" binding:"omitempty,oneof=none low medium high"`
}

// TeaListResponse represents a paginated list of teas
// @Description Paginated tea list response
type TeaListResponse struct {
    Data       []Tea      `json:"data"`
    Pagination Pagination `json:"pagination"`
}
```

### internal/models/brew.go

```go
package models

import "time"

// BrewStatus represents the status of a brew
// @Description Brew status
// @Enum preparing,steeping,ready,served,cold
type BrewStatus string

const (
    BrewPreparing BrewStatus = "preparing"
    BrewSteeping  BrewStatus = "steeping"
    BrewReady     BrewStatus = "ready"
    BrewServed    BrewStatus = "served"
    BrewCold      BrewStatus = "cold"
)

// Brew represents a brewing session
// @Description Brew session entity
type Brew struct {
    ID               string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440002"`
    TeapotID         string     `json:"teapotId" example:"550e8400-e29b-41d4-a716-446655440000"`
    TeaID            string     `json:"teaId" example:"550e8400-e29b-41d4-a716-446655440001"`
    Status           BrewStatus `json:"status" example:"steeping"`
    WaterTempCelsius int        `json:"waterTempCelsius" example:"85"`
    Notes            *string    `json:"notes,omitempty" example:"Using filtered water"`
    StartedAt        time.Time  `json:"startedAt" example:"2025-01-04T12:00:00Z"`
    CompletedAt      *time.Time `json:"completedAt,omitempty" example:"2025-01-04T12:05:00Z"`
    CreatedAt        time.Time  `json:"createdAt" example:"2025-01-04T12:00:00Z"`
    UpdatedAt        time.Time  `json:"updatedAt" example:"2025-01-04T12:00:00Z"`
}

// BrewWithDetails includes the related teapot and tea
// @Description Brew session with related entities
type BrewWithDetails struct {
    Brew
    Teapot Teapot `json:"teapot"`
    Tea    Tea    `json:"tea"`
}

// CreateBrewRequest represents the request body for creating a brew
// @Description Create brew request
type CreateBrewRequest struct {
    TeapotID         string  `json:"teapotId" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
    TeaID            string  `json:"teaId" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440001"`
    WaterTempCelsius *int    `json:"waterTempCelsius" binding:"omitempty,min=60,max=100" example:"85"`
    Notes            *string `json:"notes" binding:"omitempty,max=500"`
}

// PatchBrewRequest represents the request body for PATCH
// @Description Patch brew request
type PatchBrewRequest struct {
    Status      *BrewStatus `json:"status" binding:"omitempty,oneof=preparing steeping ready served cold"`
    Notes       *string     `json:"notes" binding:"omitempty,max=500"`
    CompletedAt *time.Time  `json:"completedAt" binding:"omitempty"`
}

// BrewQuery represents query parameters for listing brews
// @Description Brew list query parameters
type BrewQuery struct {
    PaginationQuery
    Status   *BrewStatus `form:"status" binding:"omitempty,oneof=preparing steeping ready served cold"`
    TeapotID *string     `form:"teapotId" binding:"omitempty,uuid"`
    TeaID    *string     `form:"teaId" binding:"omitempty,uuid"`
}

// BrewListResponse represents a paginated list of brews
// @Description Paginated brew list response
type BrewListResponse struct {
    Data       []Brew     `json:"data"`
    Pagination Pagination `json:"pagination"`
}
```

### internal/models/steep.go

```go
package models

import "time"

// Steep represents a single steeping cycle within a brew
// @Description Steep cycle entity
type Steep struct {
    ID              string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440003"`
    BrewID          string    `json:"brewId" example:"550e8400-e29b-41d4-a716-446655440002"`
    SteepNumber     int       `json:"steepNumber" example:"1"`
    DurationSeconds int       `json:"durationSeconds" example:"30"`
    Rating          *int      `json:"rating,omitempty" example:"4"`
    Notes           *string   `json:"notes,omitempty" example:"Light and floral"`
    CreatedAt       time.Time `json:"createdAt" example:"2025-01-04T12:01:00Z"`
}

// CreateSteepRequest represents the request body for creating a steep
// @Description Create steep request
type CreateSteepRequest struct {
    DurationSeconds int     `json:"durationSeconds" binding:"required,min=1" example:"30"`
    Rating          *int    `json:"rating" binding:"omitempty,min=1,max=5" example:"4"`
    Notes           *string `json:"notes" binding:"omitempty,max=200"`
}

// SteepListResponse represents a paginated list of steeps
// @Description Paginated steep list response
type SteepListResponse struct {
    Data       []Steep    `json:"data"`
    Pagination Pagination `json:"pagination"`
}
```

---

## Handlers to Implement

### Route Summary Table

| Method | Path | Request Body | Query Params | Success | Errors |
|--------|------|--------------|--------------|---------|--------|
| GET | `/teapots` | — | page, limit, material, style | 200 | — |
| POST | `/teapots` | CreateTeapotRequest | — | 201 | 400 |
| GET | `/teapots/:id` | — | — | 200 | 404 |
| PUT | `/teapots/:id` | UpdateTeapotRequest | — | 200 | 400, 404 |
| PATCH | `/teapots/:id` | PatchTeapotRequest | — | 200 | 400, 404 |
| DELETE | `/teapots/:id` | — | — | 204 | 404 |
| GET | `/teapots/:teapotId/brews` | — | page, limit | 200 | 404 |
| GET | `/teas` | — | page, limit, type, caffeineLevel | 200 | — |
| POST | `/teas` | CreateTeaRequest | — | 201 | 400 |
| GET | `/teas/:id` | — | — | 200 | 404 |
| PUT | `/teas/:id` | UpdateTeaRequest | — | 200 | 400, 404 |
| PATCH | `/teas/:id` | PatchTeaRequest | — | 200 | 400, 404 |
| DELETE | `/teas/:id` | — | — | 204 | 404 |
| GET | `/brews` | — | page, limit, status, teapotId, teaId | 200 | — |
| POST | `/brews` | CreateBrewRequest | — | 201 | 400 |
| GET | `/brews/:id` | — | — | 200 | 404 |
| PATCH | `/brews/:id` | PatchBrewRequest | — | 200 | 400, 404 |
| DELETE | `/brews/:id` | — | — | 204 | 404 |
| GET | `/brews/:brewId/steeps` | — | page, limit | 200 | 404 |
| POST | `/brews/:brewId/steeps` | CreateSteepRequest | — | 201 | 400, 404 |
| GET | `/health` | — | — | 200 | — |
| GET | `/health/live` | — | — | 200 | — |
| GET | `/health/ready` | — | — | 200/503 | — |
| GET | `/brew` | — | — | **418** | — |

### internal/handlers/teapots.go

```go
package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/api2spec/api2spec-fixture-gin/internal/models"
    "github.com/api2spec/api2spec-fixture-gin/internal/store"
)

type TeapotHandler struct {
    store *store.MemoryStore
}

func NewTeapotHandler(store *store.MemoryStore) *TeapotHandler {
    return &TeapotHandler{store: store}
}

// ListTeapots godoc
// @Summary List all teapots
// @Description Get a paginated list of teapots with optional filters
// @Tags teapots
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(20) minimum(1) maximum(100)
// @Param material query string false "Filter by material" Enums(ceramic, cast-iron, glass, porcelain, clay, stainless-steel)
// @Param style query string false "Filter by style" Enums(kyusu, gaiwan, english, moroccan, turkish, yixing)
// @Success 200 {object} models.TeapotListResponse
// @Router /teapots [get]
func (h *TeapotHandler) List(c *gin.Context) {
    var query models.TeapotQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        c.JSON(http.StatusBadRequest, models.Error{
            Code:    "VALIDATION_ERROR",
            Message: err.Error(),
        })
        return
    }

    // Set defaults
    if query.Page == 0 {
        query.Page = 1
    }
    if query.Limit == 0 {
        query.Limit = 20
    }

    teapots, total := h.store.ListTeapots(query)
    totalPages := (total + query.Limit - 1) / query.Limit

    c.JSON(http.StatusOK, models.TeapotListResponse{
        Data: teapots,
        Pagination: models.Pagination{
            Page:       query.Page,
            Limit:      query.Limit,
            Total:      total,
            TotalPages: totalPages,
        },
    })
}

// CreateTeapot godoc
// @Summary Create a teapot
// @Description Create a new teapot
// @Tags teapots
// @Accept json
// @Produce json
// @Param body body models.CreateTeapotRequest true "Teapot data"
// @Success 201 {object} models.Teapot
// @Failure 400 {object} models.Error
// @Router /teapots [post]
func (h *TeapotHandler) Create(c *gin.Context) {
    var req models.CreateTeapotRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.Error{
            Code:    "VALIDATION_ERROR",
            Message: err.Error(),
        })
        return
    }

    // Set default style if not provided
    if req.Style == "" {
        req.Style = models.StyleEnglish
    }

    now := time.Now().UTC()
    teapot := models.Teapot{
        ID:          uuid.New().String(),
        Name:        req.Name,
        Material:    req.Material,
        CapacityMl:  req.CapacityMl,
        Style:       req.Style,
        Description: req.Description,
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    h.store.CreateTeapot(teapot)
    c.JSON(http.StatusCreated, teapot)
}

// GetTeapot godoc
// @Summary Get a teapot by ID
// @Description Get a single teapot by its UUID
// @Tags teapots
// @Accept json
// @Produce json
// @Param id path string true "Teapot ID" format(uuid)
// @Success 200 {object} models.Teapot
// @Failure 404 {object} models.Error
// @Router /teapots/{id} [get]
func (h *TeapotHandler) Get(c *gin.Context) {
    id := c.Param("id")
    
    if _, err := uuid.Parse(id); err != nil {
        c.JSON(http.StatusBadRequest, models.Error{
            Code:    "VALIDATION_ERROR",
            Message: "Invalid teapot ID format",
        })
        return
    }

    teapot, found := h.store.GetTeapot(id)
    if !found {
        c.JSON(http.StatusNotFound, models.Error{
            Code:    "NOT_FOUND",
            Message: "Teapot not found",
        })
        return
    }

    c.JSON(http.StatusOK, teapot)
}

// UpdateTeapot godoc
// @Summary Update a teapot (full replacement)
// @Description Replace all fields of a teapot
// @Tags teapots
// @Accept json
// @Produce json
// @Param id path string true "Teapot ID" format(uuid)
// @Param body body models.UpdateTeapotRequest true "Teapot data"
// @Success 200 {object} models.Teapot
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /teapots/{id} [put]
func (h *TeapotHandler) Update(c *gin.Context) {
    id := c.Param("id")

    if _, err := uuid.Parse(id); err != nil {
        c.JSON(http.StatusBadRequest, models.Error{
            Code:    "VALIDATION_ERROR",
            Message: "Invalid teapot ID format",
        })
        return
    }

    existing, found := h.store.GetTeapot(id)
    if !found {
        c.JSON(http.StatusNotFound, models.Error{
            Code:    "NOT_FOUND",
            Message: "Teapot not found",
        })
        return
    }

    var req models.UpdateTeapotRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.Error{
            Code:    "VALIDATION_ERROR",
            Message: err.Error(),
        })
        return
    }

    teapot := models.Teapot{
        ID:          id,
        Name:        req.Name,
        Material:    req.Material,
        CapacityMl:  req.CapacityMl,
        Style:       req.Style,
        Description: req.Description,
        CreatedAt:   existing.CreatedAt,
        UpdatedAt:   time.Now().UTC(),
    }

    h.store.UpdateTeapot(teapot)
    c.JSON(http.StatusOK, teapot)
}

// PatchTeapot godoc
// @Summary Partially update a teapot
// @Description Update specific fields of a teapot
// @Tags teapots
// @Accept json
// @Produce json
// @Param id path string true "Teapot ID" format(uuid)
// @Param body body models.PatchTeapotRequest true "Fields to update"
// @Success 200 {object} models.Teapot
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /teapots/{id} [patch]
func (h *TeapotHandler) Patch(c *gin.Context) {
    id := c.Param("id")

    if _, err := uuid.Parse(id); err != nil {
        c.JSON(http.StatusBadRequest, models.Error{
            Code:    "VALIDATION_ERROR",
            Message: "Invalid teapot ID format",
        })
        return
    }

    existing, found := h.store.GetTeapot(id)
    if !found {
        c.JSON(http.StatusNotFound, models.Error{
            Code:    "NOT_FOUND",
            Message: "Teapot not found",
        })
        return
    }

    var req models.PatchTeapotRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.Error{
            Code:    "VALIDATION_ERROR",
            Message: err.Error(),
        })
        return
    }

    // Apply patches
    if req.Name != nil {
        existing.Name = *req.Name
    }
    if req.Material != nil {
        existing.Material = *req.Material
    }
    if req.CapacityMl != nil {
        existing.CapacityMl = *req.CapacityMl
    }
    if req.Style != nil {
        existing.Style = *req.Style
    }
    if req.Description != nil {
        existing.Description = req.Description
    }
    existing.UpdatedAt = time.Now().UTC()

    h.store.UpdateTeapot(existing)
    c.JSON(http.StatusOK, existing)
}

// DeleteTeapot godoc
// @Summary Delete a teapot
// @Description Delete a teapot by ID
// @Tags teapots
// @Accept json
// @Produce json
// @Param id path string true "Teapot ID" format(uuid)
// @Success 204 "No Content"
// @Failure 404 {object} models.Error
// @Router /teapots/{id} [delete]
func (h *TeapotHandler) Delete(c *gin.Context) {
    id := c.Param("id")

    if _, err := uuid.Parse(id); err != nil {
        c.JSON(http.StatusBadRequest, models.Error{
            Code:    "VALIDATION_ERROR",
            Message: "Invalid teapot ID format",
        })
        return
    }

    if !h.store.DeleteTeapot(id) {
        c.JSON(http.StatusNotFound, models.Error{
            Code:    "NOT_FOUND",
            Message: "Teapot not found",
        })
        return
    }

    c.Status(http.StatusNoContent)
}
```

### internal/handlers/health.go

```go
package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/api2spec/api2spec-fixture-gin/internal/models"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
    return &HealthHandler{}
}

// Health godoc
// @Summary Health check
// @Description Get service health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
    version := "1.0.0"
    c.JSON(http.StatusOK, models.HealthResponse{
        Status:    "ok",
        Timestamp: time.Now().UTC(),
        Version:   &version,
    })
}

// Live godoc
// @Summary Liveness probe
// @Description Kubernetes liveness probe endpoint
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health/live [get]
func (h *HealthHandler) Live(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Ready godoc
// @Summary Readiness probe
// @Description Kubernetes readiness probe endpoint
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Failure 503 {object} models.HealthResponse
// @Router /health/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
    checks := []models.HealthCheck{
        {Name: "memory", Status: "ok"},
        {Name: "database", Status: "ok"},
    }

    allOk := true
    for _, check := range checks {
        if check.Status != "ok" {
            allOk = false
            break
        }
    }

    status := "ok"
    statusCode := http.StatusOK
    if !allOk {
        status = "degraded"
        statusCode = http.StatusServiceUnavailable
    }

    c.JSON(statusCode, models.HealthResponse{
        Status:    status,
        Timestamp: time.Now().UTC(),
        Checks:    checks,
    })
}

// Brew godoc
// @Summary TIF 418 signature endpoint
// @Description Returns 418 I'm a teapot - TIF compliance signature
// @Tags health
// @Accept json
// @Produce json
// @Success 418 {object} models.TeapotResponse
// @Router /brew [get]
func (h *HealthHandler) Brew(c *gin.Context) {
    c.JSON(http.StatusTeapot, models.TeapotResponse{
        Error:   "I'm a teapot",
        Message: "This server is TIF-compliant and cannot brew coffee",
        Spec:    "https://teapotframework.dev",
    })
}
```

### internal/router/router.go

```go
package router

import (
    "github.com/gin-gonic/gin"
    "github.com/api2spec/api2spec-fixture-gin/internal/handlers"
    "github.com/api2spec/api2spec-fixture-gin/internal/store"
)

func Setup() *gin.Engine {
    r := gin.Default()

    // Initialize store
    memStore := store.NewMemoryStore()

    // Initialize handlers
    teapotHandler := handlers.NewTeapotHandler(memStore)
    teaHandler := handlers.NewTeaHandler(memStore)
    brewHandler := handlers.NewBrewHandler(memStore)
    healthHandler := handlers.NewHealthHandler()

    // Health routes
    r.GET("/health", healthHandler.Health)
    r.GET("/health/live", healthHandler.Live)
    r.GET("/health/ready", healthHandler.Ready)
    r.GET("/brew", healthHandler.Brew)

    // Teapot routes
    teapots := r.Group("/teapots")
    {
        teapots.GET("", teapotHandler.List)
        teapots.POST("", teapotHandler.Create)
        teapots.GET("/:id", teapotHandler.Get)
        teapots.PUT("/:id", teapotHandler.Update)
        teapots.PATCH("/:id", teapotHandler.Patch)
        teapots.DELETE("/:id", teapotHandler.Delete)
        teapots.GET("/:teapotId/brews", brewHandler.ListByTeapot)
    }

    // Tea routes
    teas := r.Group("/teas")
    {
        teas.GET("", teaHandler.List)
        teas.POST("", teaHandler.Create)
        teas.GET("/:id", teaHandler.Get)
        teas.PUT("/:id", teaHandler.Update)
        teas.PATCH("/:id", teaHandler.Patch)
        teas.DELETE("/:id", teaHandler.Delete)
    }

    // Brew routes
    brews := r.Group("/brews")
    {
        brews.GET("", brewHandler.List)
        brews.POST("", brewHandler.Create)
        brews.GET("/:id", brewHandler.Get)
        brews.PATCH("/:id", brewHandler.Patch)
        brews.DELETE("/:id", brewHandler.Delete)
        brews.GET("/:brewId/steeps", brewHandler.ListSteeps)
        brews.POST("/:brewId/steeps", brewHandler.CreateSteep)
    }

    return r
}
```

### cmd/server/main.go

```go
package main

import (
    "log"
    "os"

    "github.com/api2spec/api2spec-fixture-gin/internal/router"
)

func main() {
    r := router.Setup()

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    log.Printf("🫖 Tea API running at http://localhost:%s", port)
    log.Printf("   TIF signature: http://localhost:%s/brew", port)

    if err := r.Run(":" + port); err != nil {
        log.Fatal(err)
    }
}
```

---

## In-Memory Store

### internal/store/memory.go

```go
package store

import (
    "sync"

    "github.com/api2spec/api2spec-fixture-gin/internal/models"
)

type MemoryStore struct {
    mu       sync.RWMutex
    teapots  map[string]models.Teapot
    teas     map[string]models.Tea
    brews    map[string]models.Brew
    steeps   map[string]models.Steep
}

func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        teapots: make(map[string]models.Teapot),
        teas:    make(map[string]models.Tea),
        brews:   make(map[string]models.Brew),
        steeps:  make(map[string]models.Steep),
    }
}

// Teapot methods
func (s *MemoryStore) ListTeapots(query models.TeapotQuery) ([]models.Teapot, int) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var filtered []models.Teapot
    for _, t := range s.teapots {
        if query.Material != nil && t.Material != *query.Material {
            continue
        }
        if query.Style != nil && t.Style != *query.Style {
            continue
        }
        filtered = append(filtered, t)
    }

    total := len(filtered)
    start := (query.Page - 1) * query.Limit
    end := start + query.Limit

    if start >= total {
        return []models.Teapot{}, total
    }
    if end > total {
        end = total
    }

    return filtered[start:end], total
}

func (s *MemoryStore) CreateTeapot(t models.Teapot) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.teapots[t.ID] = t
}

func (s *MemoryStore) GetTeapot(id string) (models.Teapot, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    t, ok := s.teapots[id]
    return t, ok
}

func (s *MemoryStore) UpdateTeapot(t models.Teapot) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.teapots[t.ID] = t
}

func (s *MemoryStore) DeleteTeapot(id string) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    if _, ok := s.teapots[id]; !ok {
        return false
    }
    delete(s.teapots, id)
    return true
}

// Similar methods for Tea, Brew, Steep...
// (implement the same pattern for each entity)
```

---

## api2spec Configuration

### api2spec.config.yaml

```yaml
framework: gin
entry:
  - "internal/handlers/**/*.go"
  - "internal/router/**/*.go"
exclude:
  - "**/*_test.go"
output:
  path: generated/openapi.yaml
  format: yaml
openapi:
  info:
    title: Tea Brewing API
    version: 1.0.0
    description: Gin fixture API for api2spec. TIF-compliant.
  servers:
    - url: http://localhost:3000
      description: Development
  tags:
    - name: teapots
      description: Teapot management
    - name: teas
      description: Tea catalog
    - name: brews
      description: Brewing sessions
    - name: health
      description: Health checks
schemas:
  include:
    - "internal/models/**/*.go"
frameworkOptions:
  gin:
    swagStyle: true
    routerSetup: "internal/router/router.go"
```

---

## Makefile

```makefile
.PHONY: run build test clean

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

tidy:
	go mod tidy

lint:
	golangci-lint run

validate:
	api2spec check --ci
```

---

## Implementation Checklist

### Phase 1: Setup
- [ ] Initialize go module
- [ ] Create directory structure
- [ ] Create go.mod with dependencies
- [ ] Create Makefile

### Phase 2: Models
- [ ] internal/models/common.go (Pagination, Error, Health, TeapotResponse)
- [ ] internal/models/teapot.go (all Teapot types)
- [ ] internal/models/tea.go (all Tea types)
- [ ] internal/models/brew.go (all Brew types)
- [ ] internal/models/steep.go (all Steep types)

### Phase 3: Store
- [ ] internal/store/memory.go (in-memory data store with all CRUD methods)

### Phase 4: Handlers
- [ ] internal/handlers/teapots.go (List, Create, Get, Update, Patch, Delete)
- [ ] internal/handlers/teas.go (List, Create, Get, Update, Patch, Delete)
- [ ] internal/handlers/brews.go (List, Create, Get, Patch, Delete, ListByTeapot, ListSteeps, CreateSteep)
- [ ] internal/handlers/health.go (Health, Live, Ready, Brew/418)

### Phase 5: Router & Entry
- [ ] internal/router/router.go (route setup)
- [ ] cmd/server/main.go (entry point)

### Phase 6: Config & Expected Output
- [ ] api2spec.config.yaml
- [ ] expected/openapi.yaml (manually created gold standard)
- [ ] README.md

### Phase 7: Validation
- [ ] Run `make run` and test all endpoints
- [ ] Verify 418 response at GET /brew
- [ ] Run api2spec and compare output

---

## Notes for Claude Code

1. **Swag comments are the source of truth** — Include `// @Summary`, `// @Tags`, `// @Param`, `// @Success`, `// @Failure`, `// @Router` comments on every handler
2. **Go structs define schemas** — Use proper json tags and binding tags for validation
3. **Keep handlers simple** — This is a fixture, not production. In-memory storage is fine.
4. **Status codes matter** — Use correct codes (201 for create, 204 for delete, 400/404 for errors)
5. **The 418 endpoint is required** — This is the TIF signature
6. **PUT vs PATCH** — PUT requires all fields (use required binding), PATCH uses pointers for optional fields
7. **Use gin.H sparingly** — Prefer typed response structs for api2spec to parse

---

## Testing the Fixture

```bash
# Start the server
make run

# Test endpoints
curl http://localhost:3000/health
curl http://localhost:3000/brew  # Should return 418

# Create a teapot
curl -X POST http://localhost:3000/teapots \
  -H "Content-Type: application/json" \
  -d '{"name":"My Kyusu","material":"clay","capacityMl":350,"style":"kyusu"}'

# List teapots
curl http://localhost:3000/teapots

# Get teapot by ID
curl http://localhost:3000/teapots/{id}
```
