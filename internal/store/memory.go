package store

import (
	"sort"
	"sync"

	"github.com/api2spec/api2spec-fixture-gin/internal/models"
)

// MemoryStore provides thread-safe in-memory storage for all entities
type MemoryStore struct {
	mu      sync.RWMutex
	teapots map[string]models.Teapot
	teas    map[string]models.Tea
	brews   map[string]models.Brew
	steeps  map[string]models.Steep
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		teapots: make(map[string]models.Teapot),
		teas:    make(map[string]models.Tea),
		brews:   make(map[string]models.Brew),
		steeps:  make(map[string]models.Steep),
	}
}

// ===== Teapot Methods =====

// ListTeapots returns a paginated and filtered list of teapots
func (s *MemoryStore) ListTeapots(query models.TeapotQuery) ([]models.Teapot, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []models.Teapot
	for _, t := range s.teapots {
		if query.Material != nil && t.Material != *query.Material {
			continue
		}
		if query.Style != nil && t.Style != *query.Style {
			continue
		}
		filtered = append(filtered, t)
	}

	// Sort by CreatedAt descending for consistent ordering
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	total := len(filtered)
	start := (query.Page - 1) * query.Limit
	end := start + query.Limit

	if start >= total {
		return []models.Teapot{}, total
	}
	if end > total {
		end = total
	}

	return filtered[start:end], total
}

// CreateTeapot adds a new teapot to the store
func (s *MemoryStore) CreateTeapot(t models.Teapot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.teapots[t.ID] = t
}

// GetTeapot retrieves a teapot by ID
func (s *MemoryStore) GetTeapot(id string) (models.Teapot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.teapots[id]
	return t, ok
}

// UpdateTeapot updates an existing teapot
func (s *MemoryStore) UpdateTeapot(t models.Teapot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.teapots[t.ID] = t
}

// DeleteTeapot removes a teapot by ID
func (s *MemoryStore) DeleteTeapot(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.teapots[id]; !ok {
		return false
	}
	delete(s.teapots, id)
	return true
}

// ===== Tea Methods =====

// ListTeas returns a paginated and filtered list of teas
func (s *MemoryStore) ListTeas(query models.TeaQuery) ([]models.Tea, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []models.Tea
	for _, t := range s.teas {
		if query.Type != nil && t.Type != *query.Type {
			continue
		}
		if query.CaffeineLevel != nil && t.CaffeineLevel != *query.CaffeineLevel {
			continue
		}
		filtered = append(filtered, t)
	}

	// Sort by CreatedAt descending for consistent ordering
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	total := len(filtered)
	start := (query.Page - 1) * query.Limit
	end := start + query.Limit

	if start >= total {
		return []models.Tea{}, total
	}
	if end > total {
		end = total
	}

	return filtered[start:end], total
}

// CreateTea adds a new tea to the store
func (s *MemoryStore) CreateTea(t models.Tea) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.teas[t.ID] = t
}

// GetTea retrieves a tea by ID
func (s *MemoryStore) GetTea(id string) (models.Tea, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.teas[id]
	return t, ok
}

// UpdateTea updates an existing tea
func (s *MemoryStore) UpdateTea(t models.Tea) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.teas[t.ID] = t
}

// DeleteTea removes a tea by ID
func (s *MemoryStore) DeleteTea(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.teas[id]; !ok {
		return false
	}
	delete(s.teas, id)
	return true
}

// ===== Brew Methods =====

// ListBrews returns a paginated and filtered list of brews
func (s *MemoryStore) ListBrews(query models.BrewQuery) ([]models.Brew, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []models.Brew
	for _, b := range s.brews {
		if query.Status != nil && b.Status != *query.Status {
			continue
		}
		if query.TeapotID != nil && b.TeapotID != *query.TeapotID {
			continue
		}
		if query.TeaID != nil && b.TeaID != *query.TeaID {
			continue
		}
		filtered = append(filtered, b)
	}

	// Sort by CreatedAt descending for consistent ordering
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	total := len(filtered)
	start := (query.Page - 1) * query.Limit
	end := start + query.Limit

	if start >= total {
		return []models.Brew{}, total
	}
	if end > total {
		end = total
	}

	return filtered[start:end], total
}

// ListBrewsByTeapot returns brews filtered by teapot ID with pagination
func (s *MemoryStore) ListBrewsByTeapot(teapotID string, page, limit int) ([]models.Brew, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []models.Brew
	for _, b := range s.brews {
		if b.TeapotID == teapotID {
			filtered = append(filtered, b)
		}
	}

	// Sort by CreatedAt descending for consistent ordering
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	total := len(filtered)
	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		return []models.Brew{}, total
	}
	if end > total {
		end = total
	}

	return filtered[start:end], total
}

// CreateBrew adds a new brew to the store
func (s *MemoryStore) CreateBrew(b models.Brew) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.brews[b.ID] = b
}

// GetBrew retrieves a brew by ID
func (s *MemoryStore) GetBrew(id string) (models.Brew, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.brews[id]
	return b, ok
}

// UpdateBrew updates an existing brew
func (s *MemoryStore) UpdateBrew(b models.Brew) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.brews[b.ID] = b
}

// DeleteBrew removes a brew by ID
func (s *MemoryStore) DeleteBrew(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.brews[id]; !ok {
		return false
	}
	delete(s.brews, id)
	return true
}

// ===== Steep Methods =====

// ListSteepsByBrew returns steeps filtered by brew ID with pagination
func (s *MemoryStore) ListSteepsByBrew(brewID string, page, limit int) ([]models.Steep, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []models.Steep
	for _, steep := range s.steeps {
		if steep.BrewID == brewID {
			filtered = append(filtered, steep)
		}
	}

	// Sort by SteepNumber ascending
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].SteepNumber < filtered[j].SteepNumber
	})

	total := len(filtered)
	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		return []models.Steep{}, total
	}
	if end > total {
		end = total
	}

	return filtered[start:end], total
}

// CountSteepsByBrew returns the number of steeps for a brew
func (s *MemoryStore) CountSteepsByBrew(brewID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, steep := range s.steeps {
		if steep.BrewID == brewID {
			count++
		}
	}
	return count
}

// CreateSteep adds a new steep to the store
func (s *MemoryStore) CreateSteep(steep models.Steep) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.steeps[steep.ID] = steep
}

// GetSteep retrieves a steep by ID
func (s *MemoryStore) GetSteep(id string) (models.Steep, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	steep, ok := s.steeps[id]
	return steep, ok
}
