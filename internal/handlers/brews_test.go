package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/api2spec/api2spec-fixture-gin/internal/handlers"
	"github.com/api2spec/api2spec-fixture-gin/internal/models"
	"github.com/api2spec/api2spec-fixture-gin/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupBrewRouter(s *store.MemoryStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := handlers.NewBrewHandler(s)
	router.GET("/brews", handler.List)
	router.POST("/brews", handler.Create)
	router.GET("/brews/:id", handler.Get)
	router.PATCH("/brews/:id", handler.Patch)
	router.DELETE("/brews/:id", handler.Delete)
	return router
}

func setupBrewSteepRouter(s *store.MemoryStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := handlers.NewBrewHandler(s)
	router.GET("/brews/:id/steeps", handler.ListSteeps)
	router.POST("/brews/:id/steeps", handler.CreateSteep)
	return router
}

func setupTeapotBrewRouter(s *store.MemoryStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := handlers.NewBrewHandler(s)
	router.GET("/teapots/:id/brews", handler.ListByTeapot)
	return router
}

func createTestTeapot(s *store.MemoryStore) string {
	id := uuid.New().String()
	s.CreateTeapot(models.Teapot{
		ID:         id,
		Name:       "Test Teapot",
		Material:   models.MaterialCeramic,
		CapacityMl: 1000,
		Style:      models.StyleEnglish,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	return id
}

func createTestTea(s *store.MemoryStore) string {
	id := uuid.New().String()
	s.CreateTea(models.Tea{
		ID:               id,
		Name:             "Test Tea",
		Type:             models.TeaBlack,
		CaffeineLevel:    models.CaffeineHigh,
		SteepTempCelsius: 95,
		SteepTimeSeconds: 240,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	})
	return id
}

func TestBrewHandler_List(t *testing.T) {
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
				teapotID := createTestTeapot(s)
				teaID := createTestTea(s)
				s.CreateBrew(models.Brew{
					ID:               uuid.New().String(),
					TeapotID:         teapotID,
					TeaID:            teaID,
					Status:           models.BrewPreparing,
					WaterTempCelsius: 95,
					StartedAt:        time.Now(),
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				})
			},
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedTotal:  1,
		},
		{
			name: "filter by status",
			setupStore: func(s *store.MemoryStore) {
				teapotID := createTestTeapot(s)
				teaID := createTestTea(s)
				s.CreateBrew(models.Brew{
					ID:               uuid.New().String(),
					TeapotID:         teapotID,
					TeaID:            teaID,
					Status:           models.BrewPreparing,
					WaterTempCelsius: 95,
					StartedAt:        time.Now(),
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				})
				s.CreateBrew(models.Brew{
					ID:               uuid.New().String(),
					TeapotID:         teapotID,
					TeaID:            teaID,
					Status:           models.BrewReady,
					WaterTempCelsius: 95,
					StartedAt:        time.Now(),
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				})
			},
			queryParams:    "?status=ready",
			expectedStatus: http.StatusOK,
			expectedTotal:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			tt.setupStore(s)
			router := setupBrewRouter(s)

			req := httptest.NewRequest(http.MethodGet, "/brews"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response models.BrewListResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedTotal, response.Pagination.Total)
			assert.NotNil(t, response.Data)
		})
	}
}

func TestBrewHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		setupStore     func(*store.MemoryStore) (string, string)
		body           func(string, string) interface{}
		expectedStatus int
	}{
		{
			name: "valid brew",
			setupStore: func(s *store.MemoryStore) (string, string) {
				teapotID := createTestTeapot(s)
				teaID := createTestTea(s)
				return teapotID, teaID
			},
			body: func(teapotID, teaID string) interface{} {
				return models.CreateBrewRequest{
					TeapotID:         teapotID,
					TeaID:            teaID,
					WaterTempCelsius: intPtr(90),
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "valid brew without temp (uses tea default)",
			setupStore: func(s *store.MemoryStore) (string, string) {
				teapotID := createTestTeapot(s)
				teaID := createTestTea(s)
				return teapotID, teaID
			},
			body: func(teapotID, teaID string) interface{} {
				return models.CreateBrewRequest{
					TeapotID: teapotID,
					TeaID:    teaID,
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "non-existent teapot",
			setupStore: func(s *store.MemoryStore) (string, string) {
				teaID := createTestTea(s)
				return uuid.New().String(), teaID
			},
			body: func(teapotID, teaID string) interface{} {
				return models.CreateBrewRequest{
					TeapotID: teapotID,
					TeaID:    teaID,
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "non-existent tea",
			setupStore: func(s *store.MemoryStore) (string, string) {
				teapotID := createTestTeapot(s)
				return teapotID, uuid.New().String()
			},
			body: func(teapotID, teaID string) interface{} {
				return models.CreateBrewRequest{
					TeapotID: teapotID,
					TeaID:    teaID,
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid teapot UUID",
			setupStore: func(s *store.MemoryStore) (string, string) {
				teaID := createTestTea(s)
				return "not-a-uuid", teaID
			},
			body: func(teapotID, teaID string) interface{} {
				return map[string]interface{}{
					"teapotId": teapotID,
					"teaId":    teaID,
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			teapotID, teaID := tt.setupStore(s)
			router := setupBrewRouter(s)

			body, _ := json.Marshal(tt.body(teapotID, teaID))
			req := httptest.NewRequest(http.MethodPost, "/brews", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response models.Brew
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.NotEmpty(t, response.ID)
				assert.Equal(t, models.BrewPreparing, response.Status)
				assert.False(t, response.CreatedAt.IsZero())
			}
		})
	}
}

func TestBrewHandler_Get(t *testing.T) {
	tests := []struct {
		name           string
		setupStore     func(*store.MemoryStore) string
		getID          func(string) string
		expectedStatus int
	}{
		{
			name: "existing brew",
			setupStore: func(s *store.MemoryStore) string {
				teapotID := createTestTeapot(s)
				teaID := createTestTea(s)
				id := uuid.New().String()
				s.CreateBrew(models.Brew{
					ID:               id,
					TeapotID:         teapotID,
					TeaID:            teaID,
					Status:           models.BrewPreparing,
					WaterTempCelsius: 95,
					StartedAt:        time.Now(),
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				})
				return id
			},
			getID:          func(id string) string { return id },
			expectedStatus: http.StatusOK,
		},
		{
			name: "non-existent brew",
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
			router := setupBrewRouter(s)

			req := httptest.NewRequest(http.MethodGet, "/brews/"+tt.getID(id), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestBrewHandler_Patch(t *testing.T) {
	tests := []struct {
		name           string
		setupStore     func(*store.MemoryStore) string
		getID          func(string) string
		body           interface{}
		expectedStatus int
		validate       func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "patch status",
			setupStore: func(s *store.MemoryStore) string {
				teapotID := createTestTeapot(s)
				teaID := createTestTea(s)
				id := uuid.New().String()
				s.CreateBrew(models.Brew{
					ID:               id,
					TeapotID:         teapotID,
					TeaID:            teaID,
					Status:           models.BrewPreparing,
					WaterTempCelsius: 95,
					StartedAt:        time.Now(),
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				})
				return id
			},
			getID: func(id string) string { return id },
			body: map[string]interface{}{
				"status": "steeping",
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.Brew
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, models.BrewSteeping, response.Status)
			},
		},
		{
			name: "non-existent brew",
			setupStore: func(s *store.MemoryStore) string {
				return uuid.New().String()
			},
			getID: func(id string) string { return id },
			body: map[string]interface{}{
				"status": "steeping",
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			id := tt.setupStore(s)
			router := setupBrewRouter(s)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPatch, "/brews/"+tt.getID(id), bytes.NewReader(body))
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

func TestBrewHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		setupStore     func(*store.MemoryStore) string
		getID          func(string) string
		expectedStatus int
	}{
		{
			name: "delete existing brew",
			setupStore: func(s *store.MemoryStore) string {
				teapotID := createTestTeapot(s)
				teaID := createTestTea(s)
				id := uuid.New().String()
				s.CreateBrew(models.Brew{
					ID:               id,
					TeapotID:         teapotID,
					TeaID:            teaID,
					Status:           models.BrewPreparing,
					WaterTempCelsius: 95,
					StartedAt:        time.Now(),
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				})
				return id
			},
			getID:          func(id string) string { return id },
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "delete non-existent brew",
			setupStore: func(s *store.MemoryStore) string {
				return uuid.New().String()
			},
			getID:          func(id string) string { return id },
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			id := tt.setupStore(s)
			router := setupBrewRouter(s)

			req := httptest.NewRequest(http.MethodDelete, "/brews/"+tt.getID(id), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestBrewHandler_ListByTeapot(t *testing.T) {
	tests := []struct {
		name           string
		setupStore     func(*store.MemoryStore) string
		getID          func(string) string
		expectedStatus int
		expectedTotal  int
	}{
		{
			name: "list brews for teapot",
			setupStore: func(s *store.MemoryStore) string {
				teapotID := createTestTeapot(s)
				teaID := createTestTea(s)
				s.CreateBrew(models.Brew{
					ID:               uuid.New().String(),
					TeapotID:         teapotID,
					TeaID:            teaID,
					Status:           models.BrewPreparing,
					WaterTempCelsius: 95,
					StartedAt:        time.Now(),
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				})
				return teapotID
			},
			getID:          func(id string) string { return id },
			expectedStatus: http.StatusOK,
			expectedTotal:  1,
		},
		{
			name: "non-existent teapot",
			setupStore: func(s *store.MemoryStore) string {
				return uuid.New().String()
			},
			getID:          func(id string) string { return id },
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			id := tt.setupStore(s)
			router := setupTeapotBrewRouter(s)

			req := httptest.NewRequest(http.MethodGet, "/teapots/"+tt.getID(id)+"/brews", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response models.BrewListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedTotal, response.Pagination.Total)
			}
		})
	}
}

func TestBrewHandler_ListSteeps(t *testing.T) {
	tests := []struct {
		name           string
		setupStore     func(*store.MemoryStore) string
		getID          func(string) string
		expectedStatus int
		expectedTotal  int
	}{
		{
			name: "list steeps for brew",
			setupStore: func(s *store.MemoryStore) string {
				teapotID := createTestTeapot(s)
				teaID := createTestTea(s)
				brewID := uuid.New().String()
				s.CreateBrew(models.Brew{
					ID:               brewID,
					TeapotID:         teapotID,
					TeaID:            teaID,
					Status:           models.BrewPreparing,
					WaterTempCelsius: 95,
					StartedAt:        time.Now(),
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				})
				s.CreateSteep(models.Steep{
					ID:              uuid.New().String(),
					BrewID:          brewID,
					SteepNumber:     1,
					DurationSeconds: 30,
					CreatedAt:       time.Now(),
				})
				return brewID
			},
			getID:          func(id string) string { return id },
			expectedStatus: http.StatusOK,
			expectedTotal:  1,
		},
		{
			name: "non-existent brew",
			setupStore: func(s *store.MemoryStore) string {
				return uuid.New().String()
			},
			getID:          func(id string) string { return id },
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			id := tt.setupStore(s)
			router := setupBrewSteepRouter(s)

			req := httptest.NewRequest(http.MethodGet, "/brews/"+tt.getID(id)+"/steeps", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response models.SteepListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedTotal, response.Pagination.Total)
			}
		})
	}
}

func TestBrewHandler_CreateSteep(t *testing.T) {
	tests := []struct {
		name           string
		setupStore     func(*store.MemoryStore) string
		getID          func(string) string
		body           interface{}
		expectedStatus int
	}{
		{
			name: "valid steep",
			setupStore: func(s *store.MemoryStore) string {
				teapotID := createTestTeapot(s)
				teaID := createTestTea(s)
				brewID := uuid.New().String()
				s.CreateBrew(models.Brew{
					ID:               brewID,
					TeapotID:         teapotID,
					TeaID:            teaID,
					Status:           models.BrewPreparing,
					WaterTempCelsius: 95,
					StartedAt:        time.Now(),
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				})
				return brewID
			},
			getID: func(id string) string { return id },
			body: models.CreateSteepRequest{
				DurationSeconds: 30,
				Rating:          intPtr(4),
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "non-existent brew",
			setupStore: func(s *store.MemoryStore) string {
				return uuid.New().String()
			},
			getID: func(id string) string { return id },
			body: models.CreateSteepRequest{
				DurationSeconds: 30,
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "missing duration",
			setupStore: func(s *store.MemoryStore) string {
				teapotID := createTestTeapot(s)
				teaID := createTestTea(s)
				brewID := uuid.New().String()
				s.CreateBrew(models.Brew{
					ID:               brewID,
					TeapotID:         teapotID,
					TeaID:            teaID,
					Status:           models.BrewPreparing,
					WaterTempCelsius: 95,
					StartedAt:        time.Now(),
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				})
				return brewID
			},
			getID:          func(id string) string { return id },
			body:           map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			id := tt.setupStore(s)
			router := setupBrewSteepRouter(s)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/brews/"+tt.getID(id)+"/steeps", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response models.Steep
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.NotEmpty(t, response.ID)
				assert.Equal(t, 1, response.SteepNumber)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
