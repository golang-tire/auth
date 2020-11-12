package rules

import (
	"context"

	"github.com/golang-tire/auth/internal/helpers"
	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang-tire/pkg/grpcgw"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API interface {
	grpcgw.Controller
	auth.RuleServiceServer
}

type api struct {
	service Service
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
	cl := auth.NewRuleServiceClient(conn)
	_ = auth.RegisterRuleServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	auth.RegisterRuleServiceServer(server, a)
}

func (a api) ListRules(ctx context.Context, request *auth.ListRulesRequest) (*auth.ListRulesResponse, error) {
	offset, limit := helpers.GetOffsetAndLimit(request.Offset, request.Limit)
	res, err := a.service.Query(ctx, offset, limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) GetRule(ctx context.Context, request *auth.GetRuleRequest) (*auth.Rule, error) {
	res, err := a.service.Get(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) CreateRule(ctx context.Context, request *auth.CreateRuleRequest) (*auth.Rule, error) {
	res, err := a.service.Create(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) UpdateRule(ctx context.Context, request *auth.UpdateRuleRequest) (*auth.Rule, error) {
	res, err := a.service.Update(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) DeleteRule(ctx context.Context, request *auth.DeleteRuleRequest) (*empty.Empty, error) {
	_, err := a.service.Delete(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &empty.Empty{}, err
}

func New(srv Service) API {
	s := api{service: srv}
	grpcgw.RegisterController(s)
	return s
}
