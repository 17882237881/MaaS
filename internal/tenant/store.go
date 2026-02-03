package tenant

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Store interface {
	Create(name string) Tenant
	Get(id string) (Tenant, bool)
	List() []Tenant
	Delete(id string) bool
}

type InMemoryStore struct {
	mu    sync.RWMutex
	items map[string]Tenant
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{items: make(map[string]Tenant)}
}

func (s *InMemoryStore) Create(name string) Tenant {
	tenant := Tenant{
		ID:        uuid.NewString(),
		Name:      name,
		Status:    "active",
		CreatedAt: time.Now().UTC(),
	}

	s.mu.Lock()
	s.items[tenant.ID] = tenant
	s.mu.Unlock()

	return tenant
}

func (s *InMemoryStore) Get(id string) (Tenant, bool) {
	s.mu.RLock()
	tenant, ok := s.items[id]
	s.mu.RUnlock()

	return tenant, ok
}

func (s *InMemoryStore) List() []Tenant {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Tenant, 0, len(s.items))
	for _, tenant := range s.items {
		out = append(out, tenant)
	}
	return out
}

func (s *InMemoryStore) Delete(id string) bool {
	s.mu.Lock()
	_, ok := s.items[id]
	if ok {
		delete(s.items, id)
	}
	s.mu.Unlock()

	return ok
}
