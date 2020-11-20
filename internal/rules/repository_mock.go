package rules

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/golang-tire/auth/internal/entity"
	"github.com/google/uuid"
)

var errCRUD = errors.New("error crud")

// NewMockRepository create a mock repository for rules
func NewMockRepository() *MockRepository {
	return &MockRepository{items: []entity.Rule{}}
}

// MockRepository rules mock repository
type MockRepository struct {
	items []entity.Rule
}

func (m MockRepository) Get(ctx context.Context, id string) (entity.Rule, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return entity.Rule{}, gorm.ErrRecordNotFound
}

func (m MockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m MockRepository) Query(ctx context.Context, offset, limit int64) ([]entity.Rule, int, error) {
	return m.items, len(m.items), nil
}

func (m *MockRepository) Create(ctx context.Context, rule entity.Rule) (string, error) {
	UUID := uuid.New().String()
	rule.UUID = UUID
	if rule.Resource == "error" {
		return UUID, errCRUD
	}
	m.items = append(m.items, rule)
	return UUID, nil
}

func (m *MockRepository) Update(ctx context.Context, rule entity.Rule) error {
	if rule.Resource == "error" {
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

func (m *MockRepository) Delete(ctx context.Context, rule entity.Rule) error {
	for i, item := range m.items {
		if item.UUID == rule.UUID {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}

func (m *MockRepository) All(ctx context.Context) ([]entity.Rule, error) {
	return m.items, nil
}
