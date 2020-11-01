package rules

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-tire/auth/internal/entity"
	auth "github.com/golang-tire/auth/protobuf"
)

// Service encapsulates use case logic for rules.
type Service interface {
	Get(ctx context.Context, uuid string) (*auth.Rule, error)
	Query(ctx context.Context, offset, limit int64) (*auth.ListRulesResponse, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, input *auth.CreateRuleRequest) (*auth.Rule, error)
	Update(ctx context.Context, input *auth.UpdateRuleRequest) (*auth.Rule, error)
	Delete(ctx context.Context, uuid string) (*auth.Rule, error)
}

// ValidateCreateRequest validates the CreateRuleRequest fields.
func ValidateCreateRequest(c *auth.CreateRuleRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Subject, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Object, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Domain, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Action, validation.Required, validation.Length(0, 128)),
	)
}

// Validate validates the UpdateRuleRequest fields.
func ValidateUpdateRequest(u *auth.UpdateRuleRequest) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Subject, validation.Required, validation.Length(0, 128)),
		validation.Field(&u.Object, validation.Required, validation.Length(0, 128)),
		validation.Field(&u.Domain, validation.Required, validation.Length(0, 128)),
		validation.Field(&u.Action, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo Repository
}

// NewService creates a new rule service.
func NewService(repo Repository) Service {
	return service{repo}
}

// Get returns the rule with the specified the rule UUID.
func (s service) Get(ctx context.Context, UUID string) (*auth.Rule, error) {
	rule, err := s.repo.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return rule.ToProto(), nil
}

// Create creates a new rule.
func (s service) Create(ctx context.Context, req *auth.CreateRuleRequest) (*auth.Rule, error) {
	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}
	id, err := s.repo.Create(ctx, entity.Rule{
		Subject: req.Subject,
		Domain:  req.Domain,
		Object:  req.Object,
		Action:  req.Action,
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// Update updates the rule with the specified UUID.
func (s service) Update(ctx context.Context, req *auth.UpdateRuleRequest) (*auth.Rule, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	rule, err := s.repo.Get(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	rule.Subject = req.Subject
	rule.Domain = req.Domain
	rule.Object = req.Object
	rule.Action = req.Action
	rule.UpdatedAt = now

	ruleModel := entity.Rule{
		ID:        rule.ID,
		UUID:      rule.UUID,
		Subject:   req.Subject,
		Domain:    req.Domain,
		Object:    req.Object,
		Action:    req.Action,
		CreatedAt: rule.CreatedAt,
		UpdatedAt: now,
	}

	if err := s.repo.Update(ctx, ruleModel); err != nil {
		return nil, err
	}
	return rule.ToProto(), nil
}

// Delete deletes the rule with the specified UUID.
func (s service) Delete(ctx context.Context, UUID string) (*auth.Rule, error) {
	rule, err := s.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Delete(ctx, UUID); err != nil {
		return nil, err
	}
	return rule, nil
}

// Count returns the number of rules.
func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// Query returns the rules with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int64) (*auth.ListRulesResponse, error) {
	items, count, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return &auth.ListRulesResponse{
		Rules:      entity.RuleToProtoList(items),
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}
