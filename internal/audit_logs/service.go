package audit_logs

import (
	"context"

	"github.com/golang-tire/auth/internal/entity"

	auth "github.com/golang-tire/auth/internal/proto/v1"
)

// Service encapsulates use case logic for audit logs.
type Service interface {
	Get(ctx context.Context, uuid string) (*auth.AuditLog, error)
	Query(ctx context.Context, query string, offset, limit int64) (*auth.ListAuditLogsResponse, error)
	Count(ctx context.Context) (int64, error)
}

type service struct {
	repo Repository
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
func NewService(repo Repository) Service {
	return service{repo}
}
