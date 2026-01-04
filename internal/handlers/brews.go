package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/api2spec/api2spec-fixture-gin/internal/models"
	"github.com/api2spec/api2spec-fixture-gin/internal/store"
)

// BrewHandler handles brew-related endpoints
type BrewHandler struct {
	store *store.MemoryStore
}

// NewBrewHandler creates a new brew handler
func NewBrewHandler(store *store.MemoryStore) *BrewHandler {
	return &BrewHandler{store: store}
}

// List godoc
// @Summary List all brews
// @Description Get a paginated list of brews with optional filters
// @Tags brews
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(20) minimum(1) maximum(100)
// @Param status query string false "Filter by status" Enums(preparing, steeping, ready, served, cold)
// @Param teapotId query string false "Filter by teapot ID" format(uuid)
// @Param teaId query string false "Filter by tea ID" format(uuid)
// @Success 200 {object} models.BrewListResponse
// @Router /brews [get]
func (h *BrewHandler) List(c *gin.Context) {
	var query models.BrewQuery
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

	brews, total := h.store.ListBrews(query)
	totalPages := (total + query.Limit - 1) / query.Limit
	if totalPages < 0 {
		totalPages = 0
	}

	c.JSON(http.StatusOK, models.BrewListResponse{
		Data: brews,
		Pagination: models.Pagination{
			Page:       query.Page,
			Limit:      query.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// Create godoc
// @Summary Create a brew
// @Description Create a new brewing session
// @Tags brews
// @Accept json
// @Produce json
// @Param body body models.CreateBrewRequest true "Brew data"
// @Success 201 {object} models.Brew
// @Failure 400 {object} models.Error
// @Router /brews [post]
func (h *BrewHandler) Create(c *gin.Context) {
	var req models.CreateBrewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Verify teapot exists
	if _, found := h.store.GetTeapot(req.TeapotID); !found {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: "Teapot not found",
		})
		return
	}

	// Verify tea exists and get default temp
	tea, found := h.store.GetTea(req.TeaID)
	if !found {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: "Tea not found",
		})
		return
	}

	// Use tea's recommended temp if not provided
	waterTemp := tea.SteepTempCelsius
	if req.WaterTempCelsius != nil {
		waterTemp = *req.WaterTempCelsius
	}

	now := time.Now().UTC()
	brew := models.Brew{
		ID:               uuid.New().String(),
		TeapotID:         req.TeapotID,
		TeaID:            req.TeaID,
		Status:           models.BrewPreparing,
		WaterTempCelsius: waterTemp,
		Notes:            req.Notes,
		StartedAt:        now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	h.store.CreateBrew(brew)
	c.JSON(http.StatusCreated, brew)
}

// Get godoc
// @Summary Get a brew by ID
// @Description Get a single brew by its UUID
// @Tags brews
// @Accept json
// @Produce json
// @Param id path string true "Brew ID" format(uuid)
// @Success 200 {object} models.Brew
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /brews/{id} [get]
func (h *BrewHandler) Get(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid brew ID format",
		})
		return
	}

	brew, found := h.store.GetBrew(id)
	if !found {
		c.JSON(http.StatusNotFound, models.Error{
			Code:    "NOT_FOUND",
			Message: "Brew not found",
		})
		return
	}

	c.JSON(http.StatusOK, brew)
}

// Patch godoc
// @Summary Partially update a brew
// @Description Update specific fields of a brew
// @Tags brews
// @Accept json
// @Produce json
// @Param id path string true "Brew ID" format(uuid)
// @Param body body models.PatchBrewRequest true "Fields to update"
// @Success 200 {object} models.Brew
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /brews/{id} [patch]
func (h *BrewHandler) Patch(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid brew ID format",
		})
		return
	}

	existing, found := h.store.GetBrew(id)
	if !found {
		c.JSON(http.StatusNotFound, models.Error{
			Code:    "NOT_FOUND",
			Message: "Brew not found",
		})
		return
	}

	var req models.PatchBrewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Apply patches
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}
	if req.CompletedAt != nil {
		existing.CompletedAt = req.CompletedAt
	}
	existing.UpdatedAt = time.Now().UTC()

	h.store.UpdateBrew(existing)
	c.JSON(http.StatusOK, existing)
}

// Delete godoc
// @Summary Delete a brew
// @Description Delete a brew by ID
// @Tags brews
// @Accept json
// @Produce json
// @Param id path string true "Brew ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /brews/{id} [delete]
func (h *BrewHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid brew ID format",
		})
		return
	}

	if !h.store.DeleteBrew(id) {
		c.JSON(http.StatusNotFound, models.Error{
			Code:    "NOT_FOUND",
			Message: "Brew not found",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListByTeapot godoc
// @Summary List brews by teapot
// @Description Get a paginated list of brews for a specific teapot
// @Tags teapots
// @Accept json
// @Produce json
// @Param teapotId path string true "Teapot ID" format(uuid)
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(20) minimum(1) maximum(100)
// @Success 200 {object} models.BrewListResponse
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /teapots/{teapotId}/brews [get]
func (h *BrewHandler) ListByTeapot(c *gin.Context) {
	teapotID := c.Param("id")

	if _, err := uuid.Parse(teapotID); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid teapot ID format",
		})
		return
	}

	// Verify teapot exists
	if _, found := h.store.GetTeapot(teapotID); !found {
		c.JSON(http.StatusNotFound, models.Error{
			Code:    "NOT_FOUND",
			Message: "Teapot not found",
		})
		return
	}

	var query models.PaginationQuery
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

	brews, total := h.store.ListBrewsByTeapot(teapotID, query.Page, query.Limit)
	totalPages := (total + query.Limit - 1) / query.Limit
	if totalPages < 0 {
		totalPages = 0
	}

	c.JSON(http.StatusOK, models.BrewListResponse{
		Data: brews,
		Pagination: models.Pagination{
			Page:       query.Page,
			Limit:      query.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// ListSteeps godoc
// @Summary List steeps for a brew
// @Description Get a paginated list of steeps for a specific brew
// @Tags brews
// @Accept json
// @Produce json
// @Param brewId path string true "Brew ID" format(uuid)
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(20) minimum(1) maximum(100)
// @Success 200 {object} models.SteepListResponse
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /brews/{brewId}/steeps [get]
func (h *BrewHandler) ListSteeps(c *gin.Context) {
	brewID := c.Param("id")

	if _, err := uuid.Parse(brewID); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid brew ID format",
		})
		return
	}

	// Verify brew exists
	if _, found := h.store.GetBrew(brewID); !found {
		c.JSON(http.StatusNotFound, models.Error{
			Code:    "NOT_FOUND",
			Message: "Brew not found",
		})
		return
	}

	var query models.PaginationQuery
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

	steeps, total := h.store.ListSteepsByBrew(brewID, query.Page, query.Limit)
	totalPages := (total + query.Limit - 1) / query.Limit
	if totalPages < 0 {
		totalPages = 0
	}

	c.JSON(http.StatusOK, models.SteepListResponse{
		Data: steeps,
		Pagination: models.Pagination{
			Page:       query.Page,
			Limit:      query.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// CreateSteep godoc
// @Summary Create a steep for a brew
// @Description Add a new steeping cycle to a brew
// @Tags brews
// @Accept json
// @Produce json
// @Param brewId path string true "Brew ID" format(uuid)
// @Param body body models.CreateSteepRequest true "Steep data"
// @Success 201 {object} models.Steep
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /brews/{brewId}/steeps [post]
func (h *BrewHandler) CreateSteep(c *gin.Context) {
	brewID := c.Param("id")

	if _, err := uuid.Parse(brewID); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid brew ID format",
		})
		return
	}

	// Verify brew exists
	if _, found := h.store.GetBrew(brewID); !found {
		c.JSON(http.StatusNotFound, models.Error{
			Code:    "NOT_FOUND",
			Message: "Brew not found",
		})
		return
	}

	var req models.CreateSteepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Get next steep number
	steepNumber := h.store.CountSteepsByBrew(brewID) + 1

	steep := models.Steep{
		ID:              uuid.New().String(),
		BrewID:          brewID,
		SteepNumber:     steepNumber,
		DurationSeconds: req.DurationSeconds,
		Rating:          req.Rating,
		Notes:           req.Notes,
		CreatedAt:       time.Now().UTC(),
	}

	h.store.CreateSteep(steep)
	c.JSON(http.StatusCreated, steep)
}
