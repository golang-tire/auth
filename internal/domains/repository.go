package domains

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/golang-tire/auth/internal/pkg/db"

	"github.com/google/uuid"

	"github.com/golang-tire/auth/internal/entity"
)

// Repository encapsulates the logic to access domains from the data source.
type Repository interface {
	// Get returns the domain with the specified domain UUID.
	Get(ctx context.Context, uuid string) (entity.Domain, error)
	// GetByName returns the domain with the specified domain name.
	GetByName(ctx context.Context, name string) (entity.Domain, error)
	// Count returns the number of domains.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of domains with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]entity.Domain, int, error)
	// Create saves a new domain in the storage.
	Create(ctx context.Context, domain entity.Domain) (string, error)
	// Update updates the domain with given UUID in the storage.
	Update(ctx context.Context, domain entity.Domain) error
	// Delete removes the domain with given UUID from the storage.
	Delete(ctx context.Context, domain entity.Domain) error
}

// repository persists domains in database
type repository struct {
	db *db.DB
}

// NewRepository creates a new domain repository
func NewRepository(db *db.DB) Repository {
	return repository{db}
}

// Get reads the domain with the specified ID from the database.
func (r repository) Get(ctx context.Context, uuid string) (entity.Domain, error) {
	var domain entity.Domain
	res := r.db.With(ctx).Where("uuid = ?", uuid).First(&domain)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		return entity.Domain{}, fmt.Errorf("domain with uuid `%s` not found", uuid)
	}
	return domain, res.Error
}

// GetByName returns the domain with the specified domain name.
func (r repository) GetByName(ctx context.Context, name string) (entity.Domain, error) {
	var domain entity.Domain
	res := r.db.With(ctx).Where("name = ?", name).First(&domain)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		return entity.Domain{}, fmt.Errorf("domain with name `%s` not found", name)
	}
	return domain, res.Error
}

// Create saves a new domain record in the database.
// It returns the UUID of the newly inserted domain record.
func (r repository) Create(ctx context.Context, domain entity.Domain) (string, error) {
	now := time.Now()
	domain.UUID = uuid.New().String()
	domain.CreatedAt = now
	domain.UpdatedAt = now
	res := r.db.With(ctx).Create(&domain)
	return domain.UUID, res.Error
}

// Update saves the changes to an domain in the database.
func (r repository) Update(ctx context.Context, domain entity.Domain) error {
	res := r.db.With(ctx).Save(&domain)
	return res.Error
}

// Delete deletes an domain with the specified ID from the database.
func (r repository) Delete(ctx context.Context, domain entity.Domain) error {
	res := r.db.With(ctx).Delete(&domain)
	return res.Error
}

// Count returns the number of the domain records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int64
	res := r.db.With(ctx).Model(&entity.Domain{}).Count(&count)
	return count, res.Error
}

// Query retrieves the domain records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]entity.Domain, int, error) {
	var _domains []entity.Domain
	res := r.db.With(ctx).
		Limit(int(limit)).
		Offset(int(offset)).
		Order("id asc").
		Find(&_domains)

	count, err := r.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return _domains, int(count), res.Error
}
