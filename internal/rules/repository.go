package rules

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/golang-tire/auth/internal/db"

	"github.com/google/uuid"

	"github.com/golang-tire/auth/internal/entity"
)

// Repository encapsulates the logic to access rules from the data source.
type Repository interface {
	// Get returns the rule with the specified rule UUID.
	Get(ctx context.Context, uuid string) (entity.Rule, error)
	// Count returns the number of rules.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of rules with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]entity.Rule, int, error)
	// Create saves a new rule in the storage.
	Create(ctx context.Context, rule entity.Rule) (string, error)
	// Update updates the rule with given UUID in the storage.
	Update(ctx context.Context, rule entity.Rule) error
	// Delete removes the rule with given UUID from the storage.
	Delete(ctx context.Context, rule entity.Rule) error
	// All retrieves all rules records from the database.
	All(ctx context.Context) ([]entity.Rule, error)
}

// repository persists rules in database
type repository struct {
	db *db.DB
}

// NewRepository creates a new rule repository
func NewRepository(db *db.DB) Repository {
	return repository{db}
}

// Get reads the rule with the specified ID from the database.
func (r repository) Get(ctx context.Context, uuid string) (entity.Rule, error) {
	var rule entity.Rule
	res := r.db.With(ctx).
		Preload("Domain").
		Preload("Role").
		Where("rules.uuid = ?", uuid).First(&rule)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		return entity.Rule{}, fmt.Errorf("rule with uuid `%s` not found", uuid)
	}
	return rule, res.Error
}

// Create saves a new rule record in the database.
// It returns the UUID of the newly inserted rule record.
func (r repository) Create(ctx context.Context, rule entity.Rule) (string, error) {
	now := time.Now()
	rule.UUID = uuid.New().String()
	rule.CreatedAt = now
	rule.UpdatedAt = now
	res := r.db.With(ctx).Create(&rule)
	return rule.UUID, res.Error
}

// Update saves the changes to an rule in the database.
func (r repository) Update(ctx context.Context, rule entity.Rule) error {
	res := r.db.With(ctx).Save(&rule)
	return res.Error
}

// Delete deletes an rule with the specified ID from the database.
func (r repository) Delete(ctx context.Context, rule entity.Rule) error {
	res := r.db.With(ctx).Delete(&rule)
	return res.Error
}

// Count returns the number of the rule records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int64
	res := r.db.With(ctx).Model(&entity.Rule{}).Count(&count)
	return count, res.Error
}

// Query retrieves the rule records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]entity.Rule, int, error) {
	var _rule []entity.Rule
	res := r.db.With(ctx).
		Limit(int(limit)).
		Offset(int(offset)).
		Order("rules.id asc").
		Preload("Domain").
		Preload("Role").
		Find(&_rule)

	count, err := r.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return _rule, int(count), res.Error
}

// All retrieves all rules records from the database.
func (r repository) All(ctx context.Context) ([]entity.Rule, error) {
	var _rule []entity.Rule
	res := r.db.With(ctx).
		Order("rules.id asc").
		Preload("Domain").
		Preload("Role").
		Find(&_rule)

	return _rule, res.Error
}
