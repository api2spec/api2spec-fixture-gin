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
