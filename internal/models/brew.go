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
