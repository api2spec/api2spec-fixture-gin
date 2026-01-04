package handlers_test

import (
	"bytes"
	"encoding/json"
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

func setupTeaRouter(s *store.MemoryStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := handlers.NewTeaHandler(s)
	router.GET("/teas", handler.List)
	router.POST("/teas", handler.Create)
	router.GET("/teas/:id", handler.Get)
	router.PUT("/teas/:id", handler.Update)
	router.PATCH("/teas/:id", handler.Patch)
	router.DELETE("/teas/:id", handler.Delete)
	return router
}

func TestTeaHandler_List(t *testing.T) {
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
				s.CreateTea(models.Tea{
					ID:               uuid.New().String(),
					Name:             "Earl Grey",
					Type:             models.TeaBlack,
					CaffeineLevel:    models.CaffeineHigh,
					SteepTempCelsius: 95,
					SteepTimeSeconds: 240,
				})
			},
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedTotal:  1,
		},
		{
			name: "filter by type",
			setupStore: func(s *store.MemoryStore) {
				s.CreateTea(models.Tea{
					ID:               uuid.New().String(),
					Name:             "Earl Grey",
					Type:             models.TeaBlack,
					CaffeineLevel:    models.CaffeineHigh,
					SteepTempCelsius: 95,
					SteepTimeSeconds: 240,
				})
				s.CreateTea(models.Tea{
					ID:               uuid.New().String(),
					Name:             "Sencha",
					Type:             models.TeaGreen,
					CaffeineLevel:    models.CaffeineMedium,
					SteepTempCelsius: 75,
					SteepTimeSeconds: 120,
				})
			},
			queryParams:    "?type=green",
			expectedStatus: http.StatusOK,
			expectedTotal:  1,
		},
		{
			name: "filter by caffeine level",
			setupStore: func(s *store.MemoryStore) {
				s.CreateTea(models.Tea{
					ID:               uuid.New().String(),
					Name:             "Earl Grey",
					Type:             models.TeaBlack,
					CaffeineLevel:    models.CaffeineHigh,
					SteepTempCelsius: 95,
					SteepTimeSeconds: 240,
				})
				s.CreateTea(models.Tea{
					ID:               uuid.New().String(),
					Name:             "Chamomile",
					Type:             models.TeaHerbal,
					CaffeineLevel:    models.CaffeineNone,
					SteepTempCelsius: 100,
					SteepTimeSeconds: 300,
				})
			},
			queryParams:    "?caffeineLevel=none",
			expectedStatus: http.StatusOK,
			expectedTotal:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			tt.setupStore(s)
			router := setupTeaRouter(s)

			req := httptest.NewRequest(http.MethodGet, "/teas"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response models.TeaListResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedTotal, response.Pagination.Total)
			assert.NotNil(t, response.Data)
		})
	}
}

func TestTeaHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		expectedStatus int
	}{
		{
			name: "valid tea",
			body: models.CreateTeaRequest{
				Name:             "Earl Grey",
				Type:             models.TeaBlack,
				CaffeineLevel:    models.CaffeineHigh,
				SteepTempCelsius: 95,
				SteepTimeSeconds: 240,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "valid tea without caffeine level (defaults to medium)",
			body: map[string]interface{}{
				"name":             "Green Tea",
				"type":             "green",
				"steepTempCelsius": 80,
				"steepTimeSeconds": 180,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing name",
			body: map[string]interface{}{
				"type":             "black",
				"steepTempCelsius": 95,
				"steepTimeSeconds": 240,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid type",
			body: map[string]interface{}{
				"name":             "Test",
				"type":             "coffee",
				"steepTempCelsius": 95,
				"steepTimeSeconds": 240,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "temp too low",
			body: map[string]interface{}{
				"name":             "Test",
				"type":             "black",
				"steepTempCelsius": 50,
				"steepTimeSeconds": 240,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "temp too high",
			body: map[string]interface{}{
				"name":             "Test",
				"type":             "black",
				"steepTempCelsius": 110,
				"steepTimeSeconds": 240,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			router := setupTeaRouter(s)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/teas", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response models.Tea
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.NotEmpty(t, response.ID)
				assert.False(t, response.CreatedAt.IsZero())
			}
		})
	}
}

func TestTeaHandler_Get(t *testing.T) {
	tests := []struct {
		name           string
		setupStore     func(*store.MemoryStore) string
		getID          func(string) string
		expectedStatus int
	}{
		{
			name: "existing tea",
			setupStore: func(s *store.MemoryStore) string {
				id := uuid.New().String()
				s.CreateTea(models.Tea{
					ID:               id,
					Name:             "Earl Grey",
					Type:             models.TeaBlack,
					CaffeineLevel:    models.CaffeineHigh,
					SteepTempCelsius: 95,
					SteepTimeSeconds: 240,
				})
				return id
			},
			getID:          func(id string) string { return id },
			expectedStatus: http.StatusOK,
		},
		{
			name: "non-existent tea",
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
			router := setupTeaRouter(s)

			req := httptest.NewRequest(http.MethodGet, "/teas/"+tt.getID(id), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestTeaHandler_Update(t *testing.T) {
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
				s.CreateTea(models.Tea{
					ID:               id,
					Name:             "Old Name",
					Type:             models.TeaBlack,
					CaffeineLevel:    models.CaffeineHigh,
					SteepTempCelsius: 95,
					SteepTimeSeconds: 240,
				})
				return id
			},
			getID: func(id string) string { return id },
			body: models.UpdateTeaRequest{
				Name:             "New Name",
				Type:             models.TeaGreen,
				CaffeineLevel:    models.CaffeineMedium,
				SteepTempCelsius: 80,
				SteepTimeSeconds: 180,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "non-existent tea",
			setupStore: func(s *store.MemoryStore) string {
				return uuid.New().String()
			},
			getID: func(id string) string { return id },
			body: models.UpdateTeaRequest{
				Name:             "New Name",
				Type:             models.TeaGreen,
				CaffeineLevel:    models.CaffeineMedium,
				SteepTempCelsius: 80,
				SteepTimeSeconds: 180,
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			id := tt.setupStore(s)
			router := setupTeaRouter(s)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/teas/"+tt.getID(id), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestTeaHandler_Patch(t *testing.T) {
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
				s.CreateTea(models.Tea{
					ID:               id,
					Name:             "Old Name",
					Type:             models.TeaBlack,
					CaffeineLevel:    models.CaffeineHigh,
					SteepTempCelsius: 95,
					SteepTimeSeconds: 240,
				})
				return id
			},
			getID: func(id string) string { return id },
			body: map[string]interface{}{
				"name": "New Name",
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.Tea
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "New Name", response.Name)
				assert.Equal(t, models.TeaBlack, response.Type)
			},
		},
		{
			name: "non-existent tea",
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
			router := setupTeaRouter(s)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPatch, "/teas/"+tt.getID(id), bytes.NewReader(body))
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

func TestTeaHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		setupStore     func(*store.MemoryStore) string
		getID          func(string) string
		expectedStatus int
	}{
		{
			name: "delete existing tea",
			setupStore: func(s *store.MemoryStore) string {
				id := uuid.New().String()
				s.CreateTea(models.Tea{
					ID:               id,
					Name:             "Earl Grey",
					Type:             models.TeaBlack,
					CaffeineLevel:    models.CaffeineHigh,
					SteepTempCelsius: 95,
					SteepTimeSeconds: 240,
				})
				return id
			},
			getID:          func(id string) string { return id },
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "delete non-existent tea",
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
			router := setupTeaRouter(s)

			req := httptest.NewRequest(http.MethodDelete, "/teas/"+tt.getID(id), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
