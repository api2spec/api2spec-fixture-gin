package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/api2spec/api2spec-fixture-gin/internal/models"
	"github.com/api2spec/api2spec-fixture-gin/internal/store"
)

// TeaHandler handles tea-related endpoints
type TeaHandler struct {
	store *store.MemoryStore
}

// NewTeaHandler creates a new tea handler
func NewTeaHandler(store *store.MemoryStore) *TeaHandler {
	return &TeaHandler{store: store}
}

// List godoc
// @Summary List all teas
// @Description Get a paginated list of teas with optional filters
// @Tags teas
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(20) minimum(1) maximum(100)
// @Param type query string false "Filter by tea type" Enums(green, black, oolong, white, puerh, herbal, rooibos)
// @Param caffeineLevel query string false "Filter by caffeine level" Enums(none, low, medium, high)
// @Success 200 {object} models.TeaListResponse
// @Router /teas [get]
func (h *TeaHandler) List(c *gin.Context) {
	var query models.TeaQuery
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

	teas, total := h.store.ListTeas(query)
	totalPages := (total + query.Limit - 1) / query.Limit
	if totalPages < 0 {
		totalPages = 0
	}

	c.JSON(http.StatusOK, models.TeaListResponse{
		Data: teas,
		Pagination: models.Pagination{
			Page:       query.Page,
			Limit:      query.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// Create godoc
// @Summary Create a tea
// @Description Create a new tea
// @Tags teas
// @Accept json
// @Produce json
// @Param body body models.CreateTeaRequest true "Tea data"
// @Success 201 {object} models.Tea
// @Failure 400 {object} models.Error
// @Router /teas [post]
func (h *TeaHandler) Create(c *gin.Context) {
	var req models.CreateTeaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Set default caffeine level if not provided
	if req.CaffeineLevel == "" {
		req.CaffeineLevel = models.CaffeineMedium
	}

	now := time.Now().UTC()
	tea := models.Tea{
		ID:               uuid.New().String(),
		Name:             req.Name,
		Type:             req.Type,
		Origin:           req.Origin,
		CaffeineLevel:    req.CaffeineLevel,
		SteepTempCelsius: req.SteepTempCelsius,
		SteepTimeSeconds: req.SteepTimeSeconds,
		Description:      req.Description,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	h.store.CreateTea(tea)
	c.JSON(http.StatusCreated, tea)
}

// Get godoc
// @Summary Get a tea by ID
// @Description Get a single tea by its UUID
// @Tags teas
// @Accept json
// @Produce json
// @Param id path string true "Tea ID" format(uuid)
// @Success 200 {object} models.Tea
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /teas/{id} [get]
func (h *TeaHandler) Get(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid tea ID format",
		})
		return
	}

	tea, found := h.store.GetTea(id)
	if !found {
		c.JSON(http.StatusNotFound, models.Error{
			Code:    "NOT_FOUND",
			Message: "Tea not found",
		})
		return
	}

	c.JSON(http.StatusOK, tea)
}

// Update godoc
// @Summary Update a tea (full replacement)
// @Description Replace all fields of a tea
// @Tags teas
// @Accept json
// @Produce json
// @Param id path string true "Tea ID" format(uuid)
// @Param body body models.UpdateTeaRequest true "Tea data"
// @Success 200 {object} models.Tea
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /teas/{id} [put]
func (h *TeaHandler) Update(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid tea ID format",
		})
		return
	}

	existing, found := h.store.GetTea(id)
	if !found {
		c.JSON(http.StatusNotFound, models.Error{
			Code:    "NOT_FOUND",
			Message: "Tea not found",
		})
		return
	}

	var req models.UpdateTeaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	tea := models.Tea{
		ID:               id,
		Name:             req.Name,
		Type:             req.Type,
		Origin:           req.Origin,
		CaffeineLevel:    req.CaffeineLevel,
		SteepTempCelsius: req.SteepTempCelsius,
		SteepTimeSeconds: req.SteepTimeSeconds,
		Description:      req.Description,
		CreatedAt:        existing.CreatedAt,
		UpdatedAt:        time.Now().UTC(),
	}

	h.store.UpdateTea(tea)
	c.JSON(http.StatusOK, tea)
}

// Patch godoc
// @Summary Partially update a tea
// @Description Update specific fields of a tea
// @Tags teas
// @Accept json
// @Produce json
// @Param id path string true "Tea ID" format(uuid)
// @Param body body models.PatchTeaRequest true "Fields to update"
// @Success 200 {object} models.Tea
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /teas/{id} [patch]
func (h *TeaHandler) Patch(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid tea ID format",
		})
		return
	}

	existing, found := h.store.GetTea(id)
	if !found {
		c.JSON(http.StatusNotFound, models.Error{
			Code:    "NOT_FOUND",
			Message: "Tea not found",
		})
		return
	}

	var req models.PatchTeaRequest
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
	if req.Type != nil {
		existing.Type = *req.Type
	}
	if req.Origin != nil {
		existing.Origin = req.Origin
	}
	if req.CaffeineLevel != nil {
		existing.CaffeineLevel = *req.CaffeineLevel
	}
	if req.SteepTempCelsius != nil {
		existing.SteepTempCelsius = *req.SteepTempCelsius
	}
	if req.SteepTimeSeconds != nil {
		existing.SteepTimeSeconds = *req.SteepTimeSeconds
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	existing.UpdatedAt = time.Now().UTC()

	h.store.UpdateTea(existing)
	c.JSON(http.StatusOK, existing)
}

// Delete godoc
// @Summary Delete a tea
// @Description Delete a tea by ID
// @Tags teas
// @Accept json
// @Produce json
// @Param id path string true "Tea ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /teas/{id} [delete]
func (h *TeaHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid tea ID format",
		})
		return
	}

	if !h.store.DeleteTea(id) {
		c.JSON(http.StatusNotFound, models.Error{
			Code:    "NOT_FOUND",
			Message: "Tea not found",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
