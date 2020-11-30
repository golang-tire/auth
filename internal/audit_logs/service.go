package audit_logs

import (
	"context"

	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/golang-tire/auth/internal/users"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/golang-tire/auth/internal/entity"

	auth "github.com/golang-tire/auth/internal/proto/v1"
)

// Service encapsulates use case logic for audit logs.
type Service interface {
	Get(ctx context.Context, uuid string) (*auth.AuditLog, error)
	Query(ctx context.Context, query string, offset, limit int64) (*auth.ListAuditLogsResponse, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, input *auth.CreateAuditLogRequest) (*auth.AuditLog, error)
}

type service struct {
	repo     Repository
	userRepo users.Repository
}

// ValidateCreateRequest validates the CreateAuditLogRequest fields.
func ValidateCreateRequest(c *auth.CreateAuditLogRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.UserUuid, validation.Required, validation.Length(1, 128), is.UUIDv4),
		validation.Field(&c.Object, validation.Required, validation.Length(1, 128)),
		validation.Field(&c.Action, validation.Required, validation.Length(1, 128)),
		validation.Field(&c.OldValue, validation.Length(0, 128)),
		validation.Field(&c.NewValue, validation.Length(0, 128)),
	)
}

func (s service) Create(ctx context.Context, req *auth.CreateAuditLogRequest) (*auth.AuditLog, error) {
	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}

	user, err := s.userRepo.Get(ctx, req.UserUuid)
	if err != nil {
		return nil, err
	}

	id, err := s.repo.Create(ctx, entity.AuditLog{
		User:     user,
		UserID:   user.ID,
		Action:   req.Action,
		Object:   req.Object,
		OldValue: req.OldValue,
		NewValue: req.NewValue,
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

func (s service) Get(ctx context.Context, Uuid string) (*auth.AuditLog, error) {
	auditLog, err := s.repo.Get(ctx, Uuid)
	if err != nil {
		return nil, err
	}
	return auditLog.ToProto(), nil
}

func (s service) Query(ctx context.Context, query string, offset, limit int64) (*auth.ListAuditLogsResponse, error) {
	items, count, err := s.repo.Query(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	return &auth.ListAuditLogsResponse{
		AuditLogs:  entity.AuditLogToProtoList(items),
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}

func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// NewService creates a new audit logs service.
func NewService(repo Repository, userRepo users.Repository) Service {
	return service{repo, userRepo}
}
