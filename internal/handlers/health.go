package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/api2spec/api2spec-fixture-gin/internal/models"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
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
