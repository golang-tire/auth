package audit_logs

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/golang-tire/auth/internal/entity"
	"gorm.io/gorm"
)

var errCRUD = errors.New("error crud")

func NewMockRepository() *mockRepository {
	return &mockRepository{}
}

type mockRepository struct {
	items []entity.AuditLog
}

func (m mockRepository) GetByName(ctx context.Context, name string) (entity.AuditLog, error) {
	for _, item := range m.items {
		if item.Object == name {
			return item, nil
		}
	}
	return entity.AuditLog{}, gorm.ErrRecordNotFound
}

func (m *mockRepository) Create(ctx context.Context, auditLog entity.AuditLog) (string, error) {
	Uuid := uuid.New().String()
	auditLog.UUID = Uuid
	if auditLog.Object == "error" {
		return Uuid, errCRUD
	}
	m.items = append(m.items, auditLog)
	return Uuid, nil
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.AuditLog, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return entity.AuditLog{}, gorm.ErrRecordNotFound
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, query string, offset, limit int64) ([]entity.AuditLog, int, error) {
	return m.items, len(m.items), nil
}
