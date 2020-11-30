package users

import (
	"context"
	"errors"

	"github.com/golang-tire/auth/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var errCRUD = errors.New("error crud")

func NewMockRepository() *mockRepository {
	return &mockRepository{}
}

type mockRepository struct {
	items []entity.User
}

func (m mockRepository) AllUserRole(ctx context.Context) ([]entity.UserRole, error) {
	panic("implement me")
}

func (m mockRepository) FindOne(ctx context.Context, condition string, params ...interface{}) (entity.User, error) {
	panic("implement me")
}

func (m mockRepository) AddUserRole(ctx context.Context, userRole entity.UserRole) (string, error) {
	panic("implement me")
}

func (m mockRepository) GetUserRole(ctx context.Context, uuid string) (entity.UserRole, error) {
	panic("implement me")
}

func (m mockRepository) UpdateUserRole(ctx context.Context, userRole entity.UserRole) error {
	panic("implement me")
}

func (m mockRepository) DeleteUserRole(ctx context.Context, userRole entity.UserRole) error {
	panic("implement me")
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.User, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return entity.User{}, gorm.ErrRecordNotFound
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, query string, offset, limit int64) ([]entity.User, int, error) {
	return m.items, len(m.items), nil
}

func (m *mockRepository) Create(ctx context.Context, user entity.User) (string, error) {
	Uuid := uuid.New().String()
	user.UUID = Uuid
	if user.Username == "errorerror" {
		return Uuid, errCRUD
	}
	m.items = append(m.items, user)
	return Uuid, nil
}

func (m *mockRepository) Update(ctx context.Context, user entity.User) error {
	if user.Username == "errorerror" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.UUID == user.UUID {
			m.items[i] = user
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, user entity.User) error {
	for i, item := range m.items {
		if item.UUID == user.UUID {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}
