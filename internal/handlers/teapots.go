package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/api2spec/api2spec-fixture-gin/internal/models"
	"github.com/api2spec/api2spec-fixture-gin/internal/store"
)

// TeapotHandler handles teapot-related endpoints
type TeapotHandler struct {
	store *store.MemoryStore
}

// NewTeapotHandler creates a new teapot handler
func NewTeapotHandler(store *store.MemoryStore) *TeapotHandler {
	return &TeapotHandler{store: store}
}

// List godoc
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
	if totalPages < 0 {
		totalPages = 0
	}

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

// Create godoc
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

// Get godoc
// @Summary Get a teapot by ID
// @Description Get a single teapot by its UUID
// @Tags teapots
// @Accept json
// @Produce json
// @Param id path string true "Teapot ID" format(uuid)
// @Success 200 {object} models.Teapot
// @Failure 400 {object} models.Error
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

// Update godoc
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

// Patch godoc
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

// Delete godoc
// @Summary Delete a teapot
// @Description Delete a teapot by ID
// @Tags teapots
// @Accept json
// @Produce json
// @Param id path string true "Teapot ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} models.Error
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
