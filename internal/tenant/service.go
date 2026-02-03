package tenant

import (
	"errors"
	"strings"
)

var ErrInvalidName = errors.New("tenant name is required")

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(name string) (Tenant, error) {
	if strings.TrimSpace(name) == "" {
		return Tenant{}, ErrInvalidName
	}
	return s.store.Create(strings.TrimSpace(name)), nil
}

func (s *Service) Get(id string) (Tenant, bool) {
	return s.store.Get(id)
}

func (s *Service) List() []Tenant {
	return s.store.List()
}

func (s *Service) Delete(id string) bool {
	return s.store.Delete(id)
}
