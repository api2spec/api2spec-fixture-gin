package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/api2spec/api2spec-fixture-gin/internal/handlers"
	"github.com/api2spec/api2spec-fixture-gin/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHealthHandler_Health(t *testing.T) {
	handler := handlers.NewHealthHandler()
	router := gin.New()
	router.GET("/health", handler.Health)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "ok", response.Status)
	assert.NotNil(t, response.Version)
	assert.Equal(t, "1.0.0", *response.Version)
	assert.False(t, response.Timestamp.IsZero())
}

func TestHealthHandler_Live(t *testing.T) {
	handler := handlers.NewHealthHandler()
	router := gin.New()
	router.GET("/health/live", handler.Live)

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "ok", response["status"])
}

func TestHealthHandler_Ready(t *testing.T) {
	handler := handlers.NewHealthHandler()
	router := gin.New()
	router.GET("/health/ready", handler.Ready)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "ok", response.Status)
	assert.NotEmpty(t, response.Checks)
	assert.False(t, response.Timestamp.IsZero())
}

func TestHealthHandler_Brew(t *testing.T) {
	handler := handlers.NewHealthHandler()
	router := gin.New()
	router.GET("/brew", handler.Brew)

	req := httptest.NewRequest(http.MethodGet, "/brew", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTeapot, w.Code)

	var response models.TeapotResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "I'm a teapot", response.Error)
	assert.Contains(t, response.Message, "TIF-compliant")
	assert.Equal(t, "https://teapotframework.dev", response.Spec)
}
