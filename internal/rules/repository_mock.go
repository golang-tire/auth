package rules

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
	items []entity.Rule
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Rule, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return entity.Rule{}, gorm.ErrRecordNotFound
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int64) ([]entity.Rule, int, error) {
	return m.items, len(m.items), nil
}

func (m *mockRepository) Create(ctx context.Context, rule entity.Rule) (string, error) {
	Uuid := uuid.New().String()
	if rule.Object == "error" {
		return Uuid, errCRUD
	}
	m.items = append(m.items, rule)
	return Uuid, nil
}

func (m *mockRepository) Update(ctx context.Context, rule entity.Rule) error {
	if rule.Object == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.UUID == rule.UUID {
			m.items[i] = rule
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, rule entity.Rule) error {
	for i, item := range m.items {
		if item.UUID == rule.UUID {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}
