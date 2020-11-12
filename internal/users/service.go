package users

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/golang-tire/auth/internal/rules"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/golang-tire/auth/internal/entity"
	auth "github.com/golang-tire/auth/internal/proto/v1"
)

// Service encapsulates use case logic for users.
type Service interface {
	Get(ctx context.Context, uuid string) (*auth.User, error)
	Query(ctx context.Context, offset, limit int64) (*auth.ListUsersResponse, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, req *auth.CreateUserRequest) (*auth.User, error)
	Update(ctx context.Context, req *auth.UpdateUserRequest) (*auth.User, error)
	Delete(ctx context.Context, uuid string) (*auth.User, error)

	AddRule(ctx context.Context, req *auth.AddUserRuleRequest) (*auth.UserRule, error)
	UpdateRule(ctx context.Context, req *auth.UpdateUserRuleRequest) (*auth.UserRule, error)
	DeleteRule(ctx context.Context, req *auth.DeleteUserRuleRequest) (*auth.User, error)
	AddDomainRole(ctx context.Context, req *auth.AddDomainRoleRequest) (*auth.AddDomainRoleResponse, error)
	UpdateDomainRole(ctx context.Context, req *auth.UpdateDomainRoleRequest) (*auth.User, error)
	DeleteDomainRole(ctx context.Context, req *auth.DeleteDomainRoleRequest) (*auth.User, error)
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
		validation.Field(&u.Username, validation.Required, validation.Length(0, 128)),
	)
}

// ValidateAddUserRuleRequest validates the AddUserRuleRequest fields.
func ValidateAddUserRuleRequest(c *auth.AddUserRuleRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.UserUuid, validation.Required, validation.Length(36, 36), is.UUIDv4),
		validation.Field(&c.RuleUuid, validation.Required, validation.Length(36, 36), is.UUIDv4),
	)
}

// ValidateUpdateUserRuleRequest validates the UpdateUserRequest fields.
func ValidateUpdateUserRuleRequest(u *auth.UpdateUserRuleRequest) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.UserUuid, validation.Required, validation.Length(36, 36), is.UUIDv4),
		validation.Field(&u.RuleUuid, validation.Required, validation.Length(36, 36), is.UUIDv4),
		validation.Field(&u.UserRuleUuid, validation.Required, validation.Length(36, 36), is.UUIDv4),
	)
}

type service struct {
	repo     Repository
	ruleRepo rules.Repository
}

// NewService creates a new user service.
func NewService(repo Repository, ruleRepo rules.Repository) Service {
	return service{repo, ruleRepo}
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

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user.ToProto(), nil
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
	return user.ToProto(), nil
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

// AddRule add a rule item to user
func (s service) AddRule(ctx context.Context, req *auth.AddUserRuleRequest) (*auth.UserRule, error) {
	if err := ValidateAddUserRuleRequest(req); err != nil {
		return nil, err
	}
	user, err := s.repo.Get(ctx, req.UserUuid)
	if err != nil {
		return nil, err
	}

	rule, err := s.ruleRepo.Get(ctx, req.RuleUuid)
	if err != nil {
		return nil, err
	}

	userRule, err := s.repo.AddRule(ctx, user, rule)
	if err != nil {
		return nil, err
	}

	return &auth.UserRule{
		Uuid:     userRule.UUID,
		UserUuid: user.UUID,
		Rule:     rule.ToProto(),
	}, nil
}

// UpdateRule update rule item to user
func (s service) UpdateRule(ctx context.Context, req *auth.UpdateUserRuleRequest) (*auth.UserRule, error) {
	if err := ValidateUpdateUserRuleRequest(req); err != nil {
		return nil, err
	}
	user, err := s.repo.Get(ctx, req.UserUuid)
	if err != nil {
		return nil, err
	}

	rule, err := s.ruleRepo.Get(ctx, req.RuleUuid)
	if err != nil {
		return nil, err
	}

	oldUserRule, err := s.repo.GetRule(ctx, req.UserRuleUuid)
	if err != nil {
		return nil, err
	}

	userRule, err := s.repo.UpdateRule(ctx, oldUserRule.UUID, user, rule)
	if err != nil {
		return nil, err
	}

	c, _ := ptypes.TimestampProto(userRule.CreatedAt)
	u, _ := ptypes.TimestampProto(userRule.UpdatedAt)
	return &auth.UserRule{
		Uuid:      userRule.UUID,
		UserUuid:  user.UUID,
		Rule:      rule.ToProto(),
		CreatedAt: c,
		UpdatedAt: u,
	}, nil
}

func (s service) DeleteRule(ctx context.Context, req *auth.DeleteUserRuleRequest) (*auth.User, error) {
	user, err := s.repo.Get(ctx, req.UserRuleUuid)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Delete(ctx, user); err != nil {
		return nil, err
	}
	return user.ToProto(), nil
}

func (s service) AddDomainRole(ctx context.Context, req *auth.AddDomainRoleRequest) (*auth.AddDomainRoleResponse, error) {
	panic("implement me")
}

func (s service) UpdateDomainRole(ctx context.Context, req *auth.UpdateDomainRoleRequest) (*auth.User, error) {
	panic("implement me")
}

func (s service) DeleteDomainRole(ctx context.Context, req *auth.DeleteDomainRoleRequest) (*auth.User, error) {
	panic("implement me")
}
