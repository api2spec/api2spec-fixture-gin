package models

import "time"

// TeapotMaterial represents valid teapot materials
// @Description Teapot material type
// @Enum ceramic,cast-iron,glass,porcelain,clay,stainless-steel
type TeapotMaterial string

const (
	MaterialCeramic        TeapotMaterial = "ceramic"
	MaterialCastIron       TeapotMaterial = "cast-iron"
	MaterialGlass          TeapotMaterial = "glass"
	MaterialPorcelain      TeapotMaterial = "porcelain"
	MaterialClay           TeapotMaterial = "clay"
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
