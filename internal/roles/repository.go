package roles

import (
	"context"
	"time"

	"github.com/golang-tire/auth/internal/db"

	"github.com/google/uuid"

	"github.com/golang-tire/auth/internal/entity"
)

// Repository encapsulates the logic to access roles from the data source.
type Repository interface {
	// Get returns the role with the specified role UUID.
	Get(ctx context.Context, uuid string) (entity.Role, error)
	// GetByTitle returns the role with the specified role title.
	GetByTitle(ctx context.Context, title string) (entity.Role, error)
	// Count returns the number of roles.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of roles with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]entity.Role, int, error)
	// Create saves a new role in the storage.
	Create(ctx context.Context, role entity.Role) (string, error)
	// Update updates the role with given UUID in the storage.
	Update(ctx context.Context, role entity.Role) error
	// Delete removes the role with given UUID from the storage.
	Delete(ctx context.Context, role entity.Role) error
}

// repository persists roles in database
type repository struct {
	db *db.DB
}

// NewRepository creates a new role repository
func NewRepository(db *db.DB) Repository {
	return repository{db}
}

// Get reads the role with the specified ID from the database.
func (r repository) Get(ctx context.Context, uuid string) (entity.Role, error) {
	var role entity.Role
	err := r.db.With(ctx).Model(&role).Where("uuid = ?", uuid).First()
	return role, err
}

// GetByTitle returns the role with the specified role title.
func (r repository) GetByTitle(ctx context.Context, title string) (entity.Role, error) {
	var role entity.Role
	err := r.db.With(ctx).Model(&role).Where("title = ?", title).First()
	return role, err
}

// Create saves a new role record in the database.
// It returns the UUID of the newly inserted role record.
func (r repository) Create(ctx context.Context, role entity.Role) (string, error) {
	now := time.Now()
	role.UUID = uuid.New().String()
	role.CreatedAt = now
	role.UpdatedAt = now
	_, err := r.db.With(ctx).Model(&role).Insert()
	return role.UUID, err
}

// Update saves the changes to an role in the database.
func (r repository) Update(ctx context.Context, role entity.Role) error {
	_, err := r.db.With(ctx).Model(&role).WherePK().UpdateNotZero()
	return err
}

// Delete deletes an role with the specified ID from the database.
func (r repository) Delete(ctx context.Context, role entity.Role) error {
	_, err := r.db.With(ctx).Model(&role).WherePK().Delete()
	return err
}

// Count returns the number of the role records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int
	count, err := r.db.With(ctx).Model((*entity.Role)(nil)).Count()
	return int64(count), err
}

// Query retrieves the role records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]entity.Role, int, error) {
	var _roles []entity.Role
	count, err := r.db.With(ctx).Model(&_roles).
		Order("id ASC").Limit(int(limit)).
		Offset(int(offset)).SelectAndCount()
	return _roles, count, err
}
