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
	Name      string  `json:"name" example:"database"`
	Status    string  `json:"status" example:"ok" enums:"ok,degraded,down"`
	LatencyMs *int64  `json:"latencyMs,omitempty" example:"5"`
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
