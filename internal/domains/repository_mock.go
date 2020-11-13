package domains

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/golang-tire/auth/internal/entity"
	"github.com/google/uuid"
)

var errCRUD = errors.New("error crud")

func NewMockRepository() *mockRepository {
	return &mockRepository{}
}

type mockRepository struct {
	items []entity.Domain
}

func (m mockRepository) GetByName(ctx context.Context, name string) (entity.Domain, error) {
	for _, item := range m.items {
		if item.Name == name {
			return item, nil
		}
	}
	return entity.Domain{}, gorm.ErrRecordNotFound
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Domain, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return entity.Domain{}, gorm.ErrRecordNotFound
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int64) ([]entity.Domain, int, error) {
	return m.items, len(m.items), nil
}

func (m *mockRepository) Create(ctx context.Context, domain entity.Domain) (string, error) {
	Uuid := uuid.New().String()
	if domain.Name == "error" {
		return Uuid, errCRUD
	}
	m.items = append(m.items, domain)
	return Uuid, nil
}

func (m *mockRepository) Update(ctx context.Context, domain entity.Domain) error {
	if domain.Name == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.UUID == domain.UUID {
			m.items[i] = domain
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, domain entity.Domain) error {
	for i, item := range m.items {
		if item.UUID == domain.UUID {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}
