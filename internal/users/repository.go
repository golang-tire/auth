package users

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/golang-tire/auth/internal/pkg/db"

	"github.com/google/uuid"

	"github.com/golang-tire/auth/internal/entity"
)

// Repository encapsulates the logic to access users from the data source.
type Repository interface {
	// Get returns the user with the specified user UUID.
	Get(ctx context.Context, Uuid string) (entity.User, error)
	// Count returns the number of users.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of users with the given offset and limit.
	Query(ctx context.Context, query string, offset, limit int64) ([]entity.User, int, error)
	// Create saves a new user in the storage.
	Create(ctx context.Context, user entity.User) (string, error)
	// Update updates the user with given UUID in the storage.
	Update(ctx context.Context, user entity.User) error
	// Delete removes the user from the storage.
	Delete(ctx context.Context, user entity.User) error
	// AddUserRole add a role to user
	AddUserRole(ctx context.Context, userRole entity.UserRole) (string, error)
	// GetUserRole reads the user role with the specified ID from the database.
	GetUserRole(ctx context.Context, Uuid string) (entity.UserRole, error)
	// UpdateUserRole updates the user role
	UpdateUserRole(ctx context.Context, userRole entity.UserRole) error
	// DeleteUserRole delete the user role
	DeleteUserRole(ctx context.Context, userRole entity.UserRole) error
	// GetUserRole reads the user role with the specified ID from the database.
	AllUserRole(ctx context.Context) ([]entity.UserRole, error)
	// FindOne returns the one of users with the given condition
	FindOne(ctx context.Context, condition string, params ...interface{}) (entity.User, error)
}

// repository persists users in database
type repository struct {
	db *db.DB
}

// NewRepository creates a new user repository
func NewRepository(db *db.DB) Repository {
	return repository{db}
}

// Get reads the user with the specified Uuid from the database.
func (r repository) Get(ctx context.Context, Uuid string) (entity.User, error) {
	var user entity.User
	res := r.db.With(ctx).
		Preload("UserRoles.Domain").
		Preload("UserRoles.Role").
		Where("users.uuid = ?", Uuid).First(&user)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		return entity.User{}, fmt.Errorf("user with uuid `%s` not found", Uuid)
	}
	return user, res.Error
}

// Create saves a new user record in the database.
// It returns the UUID of the newly inserted user record.
func (r repository) Create(ctx context.Context, user entity.User) (string, error) {
	now := time.Now()
	user.UUID = uuid.New().String()
	user.CreatedAt = now
	user.UpdatedAt = now
	res := r.db.With(ctx).Create(&user)
	return user.UUID, res.Error
}

// Update saves the changes to an user in the database.
func (r repository) Update(ctx context.Context, user entity.User) error {
	res := r.db.With(ctx).Save(&user)
	return res.Error
}

// Delete deletes an user from the database.
func (r repository) Delete(ctx context.Context, user entity.User) error {
	res := r.db.With(ctx).Delete(&user)
	return res.Error
}

// Count returns the number of the user records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int64
	res := r.db.With(ctx).Model(&entity.User{}).Count(&count)
	return count, res.Error
}

// Query retrieves the user records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, query string, offset, limit int64) ([]entity.User, int, error) {
	var _users []entity.User
	res := r.db.With(ctx).
		Limit(int(limit)).
		Offset(int(offset)).
		Order("users.id asc").
		Preload("UserRoles.Domain").
		Preload("UserRoles.Role")

	if len(query) >= 1 {
		res = res.Where("username LIKE ? OR email LIKE ? OR firstname LIKE ? OR lastname LIKE ?", "%"+query+"%").Find(&_users)
	} else {
		res = res.Find(&_users)
	}

	count, err := r.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return _users, int(count), res.Error
}

// FindOne returns the one of users with the given condition
func (r repository) FindOne(ctx context.Context, condition string, params ...interface{}) (entity.User, error) {
	var user entity.User
	res := r.db.With(ctx).
		Preload("UserRoles.Domain").
		Preload("UserRoles.Role").
		Where(condition, params...).First(&user)
	return user, res.Error
}

func (r repository) AddUserRole(ctx context.Context, userRole entity.UserRole) (string, error) {
	now := time.Now()
	userRole.UUID = uuid.New().String()
	userRole.CreatedAt = now
	userRole.UpdatedAt = now
	res := r.db.With(ctx).Create(&userRole)
	return userRole.UUID, res.Error
}

// GetUserRole reads the user role with the specified ID from the database.
func (r repository) GetUserRole(ctx context.Context, uuid string) (entity.UserRole, error) {
	var userRole entity.UserRole
	res := r.db.With(ctx).Where("uuid = ?", uuid).First(&userRole)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		return entity.UserRole{}, fmt.Errorf("userRole with uuid `%s` not found", uuid)
	}
	return userRole, res.Error
}

func (r repository) UpdateUserRole(ctx context.Context, userRole entity.UserRole) error {
	res := r.db.With(ctx).Save(&userRole)
	return res.Error
}

func (r repository) DeleteUserRole(ctx context.Context, userRole entity.UserRole) error {
	res := r.db.With(ctx).Delete(&userRole)
	return res.Error
}

func (r repository) AllUserRole(ctx context.Context) ([]entity.UserRole, error) {
	var _userRoles []entity.UserRole
	res := r.db.With(ctx).
		Order("user_roles.id asc").
		Preload("Domain").
		Preload("Role").
		Preload("User").
		Where("user_roles.enable = ?", true).
		Find(&_userRoles)
	return _userRoles, res.Error
}
