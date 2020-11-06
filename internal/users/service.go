package users

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-tire/auth/internal/entity"
	auth "github.com/golang-tire/auth/internal/proto/v1"
)

// Service encapsulates use case logic for users.
type Service interface {
	Get(ctx context.Context, uuid string) (*auth.User, error)
	Query(ctx context.Context, offset, limit int64) (*auth.ListUsersResponse, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, input *auth.CreateUserRequest) (*auth.User, error)
	Update(ctx context.Context, input *auth.UpdateUserRequest) (*auth.User, error)
	Delete(ctx context.Context, uuid string) (*auth.User, error)
}

// ValidateCreateRequest validates the CreateUserRequest fields.
func ValidateCreateRequest(c *auth.CreateUserRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Username, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Password, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Email, validation.Required, validation.Length(0, 128)),
	)
}

// Validate validates the UpdateUserRequest fields.
func ValidateUpdateRequest(u *auth.UpdateUserRequest) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Username, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo Repository
}

// NewService creates a new user service.
func NewService(repo Repository) Service {
	return service{repo}
}

// Get returns the user with the specified the user UUID.
func (s service) Get(ctx context.Context, UUID string) (*auth.User, error) {
	user, err := s.repo.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return user.ToProto(), nil
}

// Create creates a new user.
func (s service) Create(ctx context.Context, req *auth.CreateUserRequest) (*auth.User, error) {
	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}
	id, err := s.repo.Create(ctx, entity.User{
		Firstname:   req.Firstname,
		Lastname:    req.Lastname,
		Username:    req.Username,
		Password:    req.Password,
		Gender:      req.Gender,
		AvatarURL:   req.AvatarUrl,
		Email:       req.Email,
		Enable:      req.Enable,
		RawData:     req.RawData,
		Rules:       nil,
		DomainRoles: nil,
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

	userModel := entity.User{
		ID:        user.ID,
		UUID:      user.UUID,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Username:  req.Username,
		Password:  req.Password,
		Gender:    req.Gender,
		AvatarURL: req.AvatarUrl,
		Email:     req.Email,
		Enable:    req.Enable,
		RawData:   req.RawData,
		CreatedAt: user.CreatedAt,
		UpdatedAt: now,
	}

	if err := s.repo.Update(ctx, userModel); err != nil {
		return nil, err
	}
	return user.ToProto(), nil
}

// Delete deletes the user with the specified UUID.
func (s service) Delete(ctx context.Context, UUID string) (*auth.User, error) {
	user, err := s.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Delete(ctx, UUID); err != nil {
		return nil, err
	}
	return user, nil
}

// Count returns the number of users.
func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// Query returns the users with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int64) (*auth.ListUsersResponse, error) {
	items, count, err := s.repo.Query(ctx, offset, limit)
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
