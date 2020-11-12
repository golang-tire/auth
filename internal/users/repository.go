package users

import (
	"context"
	"time"

	"github.com/golang-tire/auth/internal/db"

	"github.com/google/uuid"

	"github.com/golang-tire/auth/internal/entity"
)

// Repository encapsulates the logic to access users from the data source.
type Repository interface {
	// Get returns the user with the specified user UUID.
	Get(ctx context.Context, uuid string) (entity.User, error)
	// Count returns the number of users.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of users with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]entity.User, int, error)
	// Create saves a new user in the storage.
	Create(ctx context.Context, user entity.User) (string, error)
	// Update updates the user with given UUID in the storage.
	Update(ctx context.Context, user entity.User) error
	// Delete removes the user from the storage.
	Delete(ctx context.Context, user entity.User) error

	AddRule(ctx context.Context, user entity.User, rule entity.Rule) (*entity.UserRule, error)
	// GetRule returns the user rule with the specified user UUID.
	GetRule(ctx context.Context, uuid string) (entity.UserRule, error)
	// UpdateRule updates the user rule with given UUID in the storage.
	UpdateRule(ctx context.Context, uuid string, user entity.User, rule entity.Rule) (*entity.UserRule, error)
}

// repository persists users in database
type repository struct {
	db *db.DB
}

// NewRepository creates a new user repository
func NewRepository(db *db.DB) Repository {
	return repository{db}
}

// Get reads the user with the specified ID from the database.
func (r repository) Get(ctx context.Context, uuid string) (entity.User, error) {
	var user entity.User
	err := r.db.With(ctx).Model(&user).Where("uuid = ?", uuid).First()
	return user, err
}

// Create saves a new user record in the database.
// It returns the UUID of the newly inserted user record.
func (r repository) Create(ctx context.Context, user entity.User) (string, error) {
	now := time.Now()
	user.UUID = uuid.New().String()
	user.CreatedAt = now
	user.UpdatedAt = now
	_, err := r.db.With(ctx).Model(&user).Insert()
	return user.UUID, err
}

// Update saves the changes to an user in the database.
func (r repository) Update(ctx context.Context, user entity.User) error {
	_, err := r.db.With(ctx).Model(&user).WherePK().Update()
	return err
}

// Delete deletes an user from the database.
func (r repository) Delete(ctx context.Context, user entity.User) error {
	_, err := r.db.With(ctx).Model(&user).WherePK().Delete()
	return err
}

// Count returns the number of the user records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int
	count, err := r.db.With(ctx).Model((*entity.User)(nil)).Count()
	return int64(count), err
}

// Query retrieves the user records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]entity.User, int, error) {
	var _users []entity.User
	count, err := r.db.With(ctx).Model(&_users).
		Order("id ASC").Limit(int(limit)).
		Offset(int(offset)).SelectAndCount()
	return _users, count, err
}

// AddRule saves a new user rule record in the database.
func (r repository) AddRule(ctx context.Context, user entity.User, rule entity.Rule) (*entity.UserRule, error) {
	userRule := entity.UserRule{
		UUID:   uuid.New().String(),
		RuleID: rule.ID,
		UserID: user.ID,
	}
	_, err := r.db.With(ctx).Model(&userRule).Insert()
	if err != nil {
		return nil, err
	}
	return &userRule, nil
}

// GetRule returns the user rule with the specified user UUID.
func (r repository) GetRule(ctx context.Context, uuid string) (entity.UserRule, error) {
	var userRule entity.UserRule
	err := r.db.With(ctx).Model(&userRule).Where("uuid = ?", uuid).First()
	return userRule, err
}

// UpdateRule updates the user rule with given UUID in the storage.
func (r repository) UpdateRule(ctx context.Context, uuid string, user entity.User, rule entity.Rule) (*entity.UserRule, error) {
	userRule, err := r.GetRule(ctx, uuid)
	if err != nil {
		return nil, err
	}

	userRule.UserID = user.ID
	userRule.RuleID = rule.ID
	_, err = r.db.With(ctx).Model(&userRule).WherePK().Delete()
	return &userRule, err
}

// DeleteRule deletes an user rule with the specified UUID from the database.
func (r repository) DeleteRule(ctx context.Context, uuid string) error {
	userRule, err := r.GetRule(ctx, uuid)
	if err != nil {
		return err
	}
	_, err = r.db.With(ctx).Model(&userRule).WherePK().Delete()
	return err
}
