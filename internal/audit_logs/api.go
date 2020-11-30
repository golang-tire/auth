package audit_logs

import (
	"context"
	"net/http"

	"github.com/golang-tire/auth/internal/pkg/helpers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang-tire/pkg/grpcgw"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type API interface {
	grpcgw.Controller
}

type api struct {
	service Service
	auth.AuditLogServiceServer
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux, httpMux *http.ServeMux) {
	cl := auth.NewAuditLogServiceClient(conn)
	_ = auth.RegisterAuditLogServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	auth.RegisterAuditLogServiceServer(server, a)
}

func (a api) ListAuditLogs(ctx context.Context, request *auth.ListAuditLogsRequest) (*auth.ListAuditLogsResponse, error) {
	offset, limit := helpers.GetOffsetAndLimit(request.Offset, request.Limit)
	res, err := a.service.Query(ctx, request.Query, offset, limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) GetAuditLog(ctx context.Context, request *auth.GetAuditLogRequest) (*auth.AuditLog, error) {
	res, err := a.service.Get(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func New(srv Service) API {
	s := api{service: srv}
	grpcgw.RegisterController(s)
	return s
}
