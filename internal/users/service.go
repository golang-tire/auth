package users

import (
	"context"
	"time"

	"github.com/golang-tire/auth/internal/domains"

	"github.com/golang-tire/auth/internal/roles"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/golang-tire/auth/internal/entity"
	auth "github.com/golang-tire/auth/internal/proto/v1"
)

// Service encapsulates use case logic for users.
type Service interface {
	Get(ctx context.Context, Uuid string) (*auth.User, error)
	Query(ctx context.Context, query string, offset, limit int64) (*auth.ListUsersResponse, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, req *auth.CreateUserRequest) (*auth.User, error)
	Update(ctx context.Context, req *auth.UpdateUserRequest) (*auth.User, error)
	Delete(ctx context.Context, Uuid string) (*auth.User, error)

	// GetByUsername returns the users if username found
	GetByUsername(ctx context.Context, username string) (*auth.User, error)

	AddUserRole(ctx context.Context, req *auth.AddUserRoleRequest) (*auth.User, error)
	UpdateUserRole(ctx context.Context, req *auth.UpdateUserRoleRequest) (*auth.User, error)
	DeleteUserRole(ctx context.Context, req *auth.DeleteUserRoleRequest) (*auth.User, error)
	ListUserRoles(ctx context.Context) ([]entity.UserRole, error)
}

// ValidateCreateRequest validates the CreateUserRequest fields.
func ValidateCreateRequest(c *auth.CreateUserRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Username, validation.Required, validation.Length(6, 128)),
		validation.Field(&c.Password, validation.Required, validation.Length(8, 128)),
		validation.Field(&c.Email, validation.Required, validation.Length(4, 128), is.Email),
	)
}

// Validate validates the UpdateUserRequest fields.
func ValidateUpdateRequest(u *auth.UpdateUserRequest) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Username, validation.Required, validation.Length(6, 128)),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 128)),
		validation.Field(&u.Email, validation.Required, validation.Length(4, 128), is.Email),
	)
}

// ValidateAddUserRoleRequest validates the AddUserRoleRequest fields.
func ValidateAddUserRoleRequest(c *auth.AddUserRoleRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Uuid, validation.Required, is.UUID),
		validation.Field(&c.RoleUuid, validation.Required, is.UUID),
		validation.Field(&c.DomainUuid, validation.Required, is.UUID),
	)
}

// ValidateUpdateUserRoleRequest validates the AddUserRoleRequest fields.
func ValidateUpdateUserRoleRequest(c *auth.UpdateUserRoleRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Uuid, validation.Required, is.UUID),
		validation.Field(&c.UserRoleUuid, validation.Required, is.UUID),
		validation.Field(&c.RoleUuid, validation.Required, is.UUID),
		validation.Field(&c.DomainUuid, validation.Required, is.UUID),
	)
}

type service struct {
	repo        Repository
	domainsRepo domains.Repository
	rolesRepo   roles.Repository
}

// NewService creates a new user service.
func NewService(repo Repository, domainsRepo domains.Repository, rolesRepo roles.Repository) Service {
	return service{repo, domainsRepo, rolesRepo}
}

// Get returns the user with the specified the user Uuid.
func (s service) Get(ctx context.Context, Uuid string) (*auth.User, error) {
	user, err := s.repo.Get(ctx, Uuid)
	if err != nil {
		return nil, err
	}
	return user.ToProto(true), nil
}

// Create creates a new user.
func (s service) Create(ctx context.Context, req *auth.CreateUserRequest) (*auth.User, error) {
	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}
	id, err := s.repo.Create(ctx, entity.User{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Username:  req.Username,
		Password:  req.Password,
		Gender:    req.Gender,
		AvatarURL: req.AvatarUrl,
		Email:     req.Email,
		Enable:    req.Enable,
		RawData:   req.RawData,
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// Update updates the user with the specified UUID.
func (s service) Update(ctx context.Context, req *auth.UpdateUserRequest) (*auth.User, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	user, err := s.repo.Get(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	user.Firstname = req.Firstname
	user.Lastname = req.Lastname
	user.Username = req.Username
	user.Gender = req.Gender
	user.AvatarURL = req.AvatarUrl
	user.Enable = req.Enable
	user.RawData = req.RawData
	user.Email = req.Email
	user.UpdatedAt = now

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user.ToProto(true), nil
}

// Delete deletes the user with the specified UUID.
func (s service) Delete(ctx context.Context, UUID string) (*auth.User, error) {
	user, err := s.repo.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Delete(ctx, user); err != nil {
		return nil, err
	}
	return user.ToProto(true), nil
}

// Count returns the number of users.
func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// Query returns the users with the specified offset and limit.
func (s service) Query(ctx context.Context, query string, offset, limit int64) (*auth.ListUsersResponse, error) {
	items, count, err := s.repo.Query(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	return &auth.ListUsersResponse{
		Users:      entity.UserToProtoList(items),
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}

// GetByUsername returns the users if username found
func (s service) GetByUsername(ctx context.Context, username string) (*auth.User, error) {
	user, err := s.repo.FindOne(ctx, "users.username = ?", username)
	if err != nil {
		return nil, err
	}
	return user.ToProto(false), nil
}

func (s service) AddUserRole(ctx context.Context, req *auth.AddUserRoleRequest) (*auth.User, error) {
	if err := ValidateAddUserRoleRequest(req); err != nil {
		return nil, err
	}
	user, err := s.repo.Get(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	domain, err := s.domainsRepo.Get(ctx, req.DomainUuid)
	if err != nil {
		return nil, err
	}

	role, err := s.rolesRepo.Get(ctx, req.RoleUuid)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.AddUserRole(ctx, entity.UserRole{
		Role:   role,
		User:   user,
		Domain: domain,
		Enable: req.Enable,
	})
	if err != nil {
		return nil, err
	}

	// get updated user with its latest roles
	return s.Get(ctx, user.UUID)
}

func (s service) UpdateUserRole(ctx context.Context, req *auth.UpdateUserRoleRequest) (*auth.User, error) {

	if err := ValidateUpdateUserRoleRequest(req); err != nil {
		return nil, err
	}

	user, err := s.repo.Get(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	userRole, err := s.repo.GetUserRole(ctx, req.UserRoleUuid)
	if err != nil {
		return nil, err
	}

	domain, err := s.domainsRepo.Get(ctx, req.DomainUuid)
	if err != nil {
		return nil, err
	}

	role, err := s.rolesRepo.Get(ctx, req.RoleUuid)
	if err != nil {
		return nil, err
	}

	userRole.Domain = domain
	userRole.Role = role
	userRole.Enable = req.Enable

	err = s.repo.UpdateUserRole(ctx, userRole)
	if err != nil {
		return nil, err
	}

	// get updated user with its latest roles
	return s.Get(ctx, user.UUID)
}

func (s service) DeleteUserRole(ctx context.Context, req *auth.DeleteUserRoleRequest) (*auth.User, error) {
	userRole, err := s.repo.GetUserRole(ctx, req.UserRoleUuid)
	if err != nil {
		return nil, err
	}
	if err = s.repo.DeleteUserRole(ctx, userRole); err != nil {
		return nil, err
	}
	// get updated user with its latest roles
	return s.Get(ctx, req.Uuid)
}

func (s service) ListUserRoles(ctx context.Context) ([]entity.UserRole, error) {
	items, err := s.repo.AllUserRole(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}
