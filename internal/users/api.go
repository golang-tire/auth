package users

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
	auth.UserServiceServer
}

type api struct {
	service Service
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
	cl := auth.NewUserServiceClient(conn)
	_ = auth.RegisterUserServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	auth.RegisterUserServiceServer(server, a)
}

func (a api) ListUsers(ctx context.Context, request *auth.ListUsersRequest) (*auth.ListUsersResponse, error) {
	offset, limit := helpers.GetOffsetAndLimit(request.Offset, request.Limit)
	res, err := a.service.Query(ctx, offset, limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) GetUser(ctx context.Context, request *auth.GetUserRequest) (*auth.User, error) {
	res, err := a.service.Get(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) CreateUser(ctx context.Context, request *auth.CreateUserRequest) (*auth.User, error) {
	res, err := a.service.Create(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) UpdateUser(ctx context.Context, request *auth.UpdateUserRequest) (*auth.User, error) {
	res, err := a.service.Update(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) DeleteUser(ctx context.Context, request *auth.DeleteUserRequest) (*empty.Empty, error) {
	_, err := a.service.Delete(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &empty.Empty{}, err
}

func (a api) AddRule(ctx context.Context, request *auth.AddUserRuleRequest) (*auth.UserRule, error) {
	panic("implement me")
}

func (a api) UpdateRule(ctx context.Context, request *auth.UpdateUserRuleRequest) (*auth.UserRule, error) {
	panic("implement me")
}

func (a api) DeleteRule(ctx context.Context, request *auth.DeleteUserRuleRequest) (*empty.Empty, error) {
	panic("implement me")
}

func (a api) AddDomainRole(ctx context.Context, request *auth.AddDomainRoleRequest) (*auth.AddDomainRoleResponse, error) {
	panic("implement me")
}

func (a api) UpdateDomainRole(ctx context.Context, request *auth.UpdateDomainRoleRequest) (*auth.User, error) {
	panic("implement me")
}

func (a api) DeleteDomainRole(ctx context.Context, request *auth.DeleteDomainRoleRequest) (*empty.Empty, error) {
	panic("implement me")
}

func New(srv Service) API {
	s := api{service: srv}
	grpcgw.RegisterController(s)
	return s
}
