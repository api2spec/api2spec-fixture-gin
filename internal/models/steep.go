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
