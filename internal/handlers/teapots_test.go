package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/api2spec/api2spec-fixture-gin/internal/handlers"
	"github.com/api2spec/api2spec-fixture-gin/internal/models"
	"github.com/api2spec/api2spec-fixture-gin/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTeapotRouter(s *store.MemoryStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := handlers.NewTeapotHandler(s)
	router.GET("/teapots", handler.List)
	router.POST("/teapots", handler.Create)
	router.GET("/teapots/:id", handler.Get)
	router.PUT("/teapots/:id", handler.Update)
	router.PATCH("/teapots/:id", handler.Patch)
	router.DELETE("/teapots/:id", handler.Delete)
	return router
}

func TestTeapotHandler_List(t *testing.T) {
	tests := []struct {
		name           string
		setupStore     func(*store.MemoryStore)
		queryParams    string
		expectedStatus int
		expectedTotal  int
	}{
		{
			name:           "empty list",
			setupStore:     func(s *store.MemoryStore) {},
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedTotal:  0,
		},
		{
			name: "list with items",
			setupStore: func(s *store.MemoryStore) {
				s.CreateTeapot(models.Teapot{
					ID:         uuid.New().String(),
					Name:       "Test Teapot",
					Material:   models.MaterialCeramic,
					CapacityMl: 1000,
					Style:      models.StyleEnglish,
				})
			},
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedTotal:  1,
		},
		{
			name: "filter by material",
			setupStore: func(s *store.MemoryStore) {
				s.CreateTeapot(models.Teapot{
					ID:         uuid.New().String(),
					Name:       "Ceramic Teapot",
					Material:   models.MaterialCeramic,
					CapacityMl: 1000,
					Style:      models.StyleEnglish,
				})
				s.CreateTeapot(models.Teapot{
					ID:         uuid.New().String(),
					Name:       "Glass Teapot",
					Material:   models.MaterialGlass,
					CapacityMl: 800,
					Style:      models.StyleEnglish,
				})
			},
			queryParams:    "?material=ceramic",
			expectedStatus: http.StatusOK,
			expectedTotal:  1,
		},
		{
			name: "filter by style",
			setupStore: func(s *store.MemoryStore) {
				s.CreateTeapot(models.Teapot{
					ID:         uuid.New().String(),
					Name:       "English Teapot",
					Material:   models.MaterialCeramic,
					CapacityMl: 1000,
					Style:      models.StyleEnglish,
				})
				s.CreateTeapot(models.Teapot{
					ID:         uuid.New().String(),
					Name:       "Kyusu",
					Material:   models.MaterialClay,
					CapacityMl: 350,
					Style:      models.StyleKyusu,
				})
			},
			queryParams:    "?style=kyusu",
			expectedStatus: http.StatusOK,
			expectedTotal:  1,
		},
		{
			name: "pagination",
			setupStore: func(s *store.MemoryStore) {
				for i := 0; i < 25; i++ {
					s.CreateTeapot(models.Teapot{
						ID:         uuid.New().String(),
						Name:       fmt.Sprintf("Teapot %d", i),
						Material:   models.MaterialCeramic,
						CapacityMl: 1000,
						Style:      models.StyleEnglish,
					})
				}
			},
			queryParams:    "?page=1&limit=10",
			expectedStatus: http.StatusOK,
			expectedTotal:  25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			tt.setupStore(s)
			router := setupTeapotRouter(s)

			req := httptest.NewRequest(http.MethodGet, "/teapots"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response models.TeapotListResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedTotal, response.Pagination.Total)
			assert.NotNil(t, response.Data)
		})
	}
}

func TestTeapotHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		expectedStatus int
	}{
		{
			name: "valid teapot",
			body: models.CreateTeapotRequest{
				Name:       "My Teapot",
				Material:   models.MaterialCeramic,
				CapacityMl: 1000,
				Style:      models.StyleEnglish,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "valid teapot without style (defaults to english)",
			body: models.CreateTeapotRequest{
				Name:       "My Teapot",
				Material:   models.MaterialCeramic,
				CapacityMl: 1000,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing name",
			body: map[string]interface{}{
				"material":   "ceramic",
				"capacityMl": 1000,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid material",
			body: map[string]interface{}{
				"name":       "Test",
				"material":   "plastic",
				"capacityMl": 1000,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "capacity too large",
			body: map[string]interface{}{
				"name":       "Test",
				"material":   "ceramic",
				"capacityMl": 10000,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			router := setupTeapotRouter(s)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/teapots", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response models.Teapot
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.NotEmpty(t, response.ID)
				assert.False(t, response.CreatedAt.IsZero())
			}
		})
	}
}

func TestTeapotHandler_Get(t *testing.T) {
	tests := []struct {
		name           string
		setupStore     func(*store.MemoryStore) string
		getID          func(string) string
		expectedStatus int
	}{
		{
			name: "existing teapot",
			setupStore: func(s *store.MemoryStore) string {
				id := uuid.New().String()
				s.CreateTeapot(models.Teapot{
					ID:         id,
					Name:       "Test Teapot",
					Material:   models.MaterialCeramic,
					CapacityMl: 1000,
					Style:      models.StyleEnglish,
				})
				return id
			},
			getID:          func(id string) string { return id },
			expectedStatus: http.StatusOK,
		},
		{
			name: "non-existent teapot",
			setupStore: func(s *store.MemoryStore) string {
				return uuid.New().String()
			},
			getID:          func(id string) string { return id },
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "invalid UUID",
			setupStore: func(s *store.MemoryStore) string {
				return ""
			},
			getID:          func(id string) string { return "not-a-uuid" },
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			id := tt.setupStore(s)
			router := setupTeapotRouter(s)

			req := httptest.NewRequest(http.MethodGet, "/teapots/"+tt.getID(id), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestTeapotHandler_Update(t *testing.T) {
	tests := []struct {
		name           string
		setupStore     func(*store.MemoryStore) string
		getID          func(string) string
		body           interface{}
		expectedStatus int
	}{
		{
			name: "valid update",
			setupStore: func(s *store.MemoryStore) string {
				id := uuid.New().String()
				s.CreateTeapot(models.Teapot{
					ID:         id,
					Name:       "Old Name",
					Material:   models.MaterialCeramic,
					CapacityMl: 1000,
					Style:      models.StyleEnglish,
				})
				return id
			},
			getID: func(id string) string { return id },
			body: models.UpdateTeapotRequest{
				Name:       "New Name",
				Material:   models.MaterialGlass,
				CapacityMl: 800,
				Style:      models.StyleKyusu,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "non-existent teapot",
			setupStore: func(s *store.MemoryStore) string {
				return uuid.New().String()
			},
			getID: func(id string) string { return id },
			body: models.UpdateTeapotRequest{
				Name:       "New Name",
				Material:   models.MaterialGlass,
				CapacityMl: 800,
				Style:      models.StyleKyusu,
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "missing required field",
			setupStore: func(s *store.MemoryStore) string {
				id := uuid.New().String()
				s.CreateTeapot(models.Teapot{
					ID:         id,
					Name:       "Old Name",
					Material:   models.MaterialCeramic,
					CapacityMl: 1000,
					Style:      models.StyleEnglish,
				})
				return id
			},
			getID: func(id string) string { return id },
			body: map[string]interface{}{
				"name":     "New Name",
				"material": "glass",
				// missing capacityMl and style
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			id := tt.setupStore(s)
			router := setupTeapotRouter(s)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/teapots/"+tt.getID(id), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestTeapotHandler_Patch(t *testing.T) {
	tests := []struct {
		name           string
		setupStore     func(*store.MemoryStore) string
		getID          func(string) string
		body           interface{}
		expectedStatus int
		validate       func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "patch name only",
			setupStore: func(s *store.MemoryStore) string {
				id := uuid.New().String()
				s.CreateTeapot(models.Teapot{
					ID:         id,
					Name:       "Old Name",
					Material:   models.MaterialCeramic,
					CapacityMl: 1000,
					Style:      models.StyleEnglish,
				})
				return id
			},
			getID: func(id string) string { return id },
			body: map[string]interface{}{
				"name": "New Name",
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.Teapot
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "New Name", response.Name)
				assert.Equal(t, models.MaterialCeramic, response.Material)
			},
		},
		{
			name: "patch material only",
			setupStore: func(s *store.MemoryStore) string {
				id := uuid.New().String()
				s.CreateTeapot(models.Teapot{
					ID:         id,
					Name:       "Test",
					Material:   models.MaterialCeramic,
					CapacityMl: 1000,
					Style:      models.StyleEnglish,
				})
				return id
			},
			getID: func(id string) string { return id },
			body: map[string]interface{}{
				"material": "glass",
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.Teapot
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "Test", response.Name)
				assert.Equal(t, models.MaterialGlass, response.Material)
			},
		},
		{
			name: "non-existent teapot",
			setupStore: func(s *store.MemoryStore) string {
				return uuid.New().String()
			},
			getID: func(id string) string { return id },
			body: map[string]interface{}{
				"name": "New Name",
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			id := tt.setupStore(s)
			router := setupTeapotRouter(s)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPatch, "/teapots/"+tt.getID(id), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestTeapotHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		setupStore     func(*store.MemoryStore) string
		getID          func(string) string
		expectedStatus int
	}{
		{
			name: "delete existing teapot",
			setupStore: func(s *store.MemoryStore) string {
				id := uuid.New().String()
				s.CreateTeapot(models.Teapot{
					ID:         id,
					Name:       "Test Teapot",
					Material:   models.MaterialCeramic,
					CapacityMl: 1000,
					Style:      models.StyleEnglish,
				})
				return id
			},
			getID:          func(id string) string { return id },
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "delete non-existent teapot",
			setupStore: func(s *store.MemoryStore) string {
				return uuid.New().String()
			},
			getID:          func(id string) string { return id },
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "invalid UUID",
			setupStore: func(s *store.MemoryStore) string {
				return ""
			},
			getID:          func(id string) string { return "not-a-uuid" },
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			id := tt.setupStore(s)
			router := setupTeapotRouter(s)

			req := httptest.NewRequest(http.MethodDelete, "/teapots/"+tt.getID(id), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
