package audit_logs

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/golang-tire/auth/internal/pkg/db"

	"github.com/golang-tire/auth/internal/entity"
)

// Repository encapsulates the logic to access audit logs from the data source.
type Repository interface {
	// Get returns the auditLog with the specified auditLog UUID.
	Get(ctx context.Context, uuid string) (entity.AuditLog, error)
	// Count returns the number of auditLogs.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of auditLogs with the given offset and limit.
	Query(ctx context.Context, query string, offset, limit int64) ([]entity.AuditLog, int, error)
}

// repository persists audit logs in database
type repository struct {
	db *db.DB
}

func (r repository) Get(ctx context.Context, uuid string) (entity.AuditLog, error) {
	var auditLog entity.AuditLog
	res := r.db.With(ctx).Where("uuid = ?", uuid).First(&auditLog)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		return entity.AuditLog{}, fmt.Errorf("auditLog with uuid `%s` not found", uuid)
	}
	return auditLog, res.Error
}

func (r repository) Count(ctx context.Context) (int64, error) {
	var count int64
	res := r.db.With(ctx).Model(&entity.AuditLog{}).Count(&count)
	return count, res.Error
}

func (r repository) Query(ctx context.Context, query string, offset, limit int64) ([]entity.AuditLog, int, error) {
	var _auditLogs []entity.AuditLog
	res := r.db.With(ctx).
		Limit(int(limit)).
		Offset(int(offset)).
		Order("id asc")

	if len(query) >= 1 {
		subQuery := r.db.With(ctx).Select("id").Where("username LIKE ?", "%"+query+"%").Table("users")
		res = res.Where("user_id IN (?)", subQuery).Find(&_auditLogs)
	} else {
		res = res.Find(&_auditLogs)
	}

	count, err := r.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return _auditLogs, int(count), res.Error
}

// NewRepository creates a new audit log repository
func NewRepository(db *db.DB) Repository {
	return repository{db}
}
