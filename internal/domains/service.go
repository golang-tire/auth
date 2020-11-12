package domains

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-tire/auth/internal/entity"
	auth "github.com/golang-tire/auth/internal/proto/v1"
)

// Service encapsulates use case logic for domains.
type Service interface {
	Get(ctx context.Context, uuid string) (*auth.Domain, error)
	Query(ctx context.Context, offset, limit int64) (*auth.ListDomainsResponse, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, input *auth.CreateDomainRequest) (*auth.Domain, error)
	Update(ctx context.Context, input *auth.UpdateDomainRequest) (*auth.Domain, error)
	Delete(ctx context.Context, uuid string) (*auth.Domain, error)
}

// ValidateCreateRequest validates the CreateDomainRequest fields.
func ValidateCreateRequest(c *auth.CreateDomainRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.Length(0, 128)),
	)
}

// Validate validates the UpdateDomainRequest fields.
func ValidateUpdateRequest(u *auth.UpdateDomainRequest) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Name, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo Repository
}

// NewService creates a new domain service.
func NewService(repo Repository) Service {
	return service{repo}
}

// Get returns the domain with the specified the domain UUID.
func (s service) Get(ctx context.Context, UUID string) (*auth.Domain, error) {
	domain, err := s.repo.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return domain.ToProto(), nil
}

// Create creates a new domain.
func (s service) Create(ctx context.Context, req *auth.CreateDomainRequest) (*auth.Domain, error) {
	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}
	id, err := s.repo.Create(ctx, entity.Domain{
		Name:   req.Name,
		Enable: req.Enable,
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// Update updates the domain with the specified UUID.
func (s service) Update(ctx context.Context, req *auth.UpdateDomainRequest) (*auth.Domain, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	domain, err := s.repo.Get(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	domain.Name = req.Name
	domain.Enable = req.Enable
	domain.UpdatedAt = now

	if err := s.repo.Update(ctx, domain); err != nil {
		return nil, err
	}
	return domain.ToProto(), nil
}

// Delete deletes the domain with the specified UUID.
func (s service) Delete(ctx context.Context, UUID string) (*auth.Domain, error) {
	domain, err := s.repo.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Delete(ctx, domain); err != nil {
		return nil, err
	}
	return domain.ToProto(), nil
}

// Count returns the number of domains.
func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// Query returns the domains with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int64) (*auth.ListDomainsResponse, error) {
	items, count, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return &auth.ListDomainsResponse{
		Domains:    entity.DomainToProtoList(items),
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}
