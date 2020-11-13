package rules

import (
	"context"
	"time"

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
	err := r.db.With(ctx).Model(&rule).
		Relation("Domain").
		Where("rule.uuid = ?", uuid).First()
	return rule, err
}

// Create saves a new rule record in the database.
// It returns the UUID of the newly inserted rule record.
func (r repository) Create(ctx context.Context, rule entity.Rule) (string, error) {
	now := time.Now()
	rule.UUID = uuid.New().String()
	rule.CreatedAt = now
	rule.UpdatedAt = now
	_, err := r.db.With(ctx).Model(&rule).Insert()
	return rule.UUID, err
}

// Update saves the changes to an rule in the database.
func (r repository) Update(ctx context.Context, rule entity.Rule) error {
	_, err := r.db.With(ctx).Model(&rule).WherePK().UpdateNotZero()
	return err
}

// Delete deletes an rule with the specified ID from the database.
func (r repository) Delete(ctx context.Context, rule entity.Rule) error {
	_, err := r.db.With(ctx).Model(&rule).WherePK().Delete()
	return err
}

// Count returns the number of the rule records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int
	count, err := r.db.With(ctx).Model((*entity.Rule)(nil)).Count()
	return int64(count), err
}

// Query retrieves the rule records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]entity.Rule, int, error) {
	var _rules []entity.Rule
	count, err := r.db.With(ctx).Model(&_rules).
		Relation("Domain").
		Order("id ASC").Limit(int(limit)).
		Offset(int(offset)).SelectAndCount()
	return _rules, count, err
}
