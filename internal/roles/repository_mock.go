package roles

import (
	"context"
	"errors"

	"github.com/go-pg/pg"
	"github.com/golang-tire/auth/internal/entity"
	"github.com/google/uuid"
)

var errCRUD = errors.New("error crud")

func NewMockRepository() *mockRepository {
	return &mockRepository{}
}

type mockRepository struct {
	items []entity.Role
}

func (m mockRepository) GetByTitle(ctx context.Context, title string) (entity.Role, error) {
	for _, item := range m.items {
		if item.Title == title {
			return item, nil
		}
	}
	return entity.Role{}, pg.ErrNoRows
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Role, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return entity.Role{}, pg.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int64) ([]entity.Role, int, error) {
	return m.items, len(m.items), nil
}

func (m *mockRepository) Create(ctx context.Context, role entity.Role) (string, error) {
	Uuid := uuid.New().String()
	if role.Title == "error" {
		return Uuid, errCRUD
	}
	m.items = append(m.items, role)
	return Uuid, nil
}

func (m *mockRepository) Update(ctx context.Context, role entity.Role) error {
	if role.Title == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.UUID == role.UUID {
			m.items[i] = role
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, role entity.Role) error {
	for i, item := range m.items {
		if item.UUID == role.UUID {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}
