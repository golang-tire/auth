package users

import (
	"context"
	"net/http"

	"github.com/golang-tire/auth/internal/pkg/helpers"
	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang-tire/pkg/grpcgw"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API interface {
	grpcgw.Controller
}

type api struct {
	service Service
	auth.UserServiceServer
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux, httpMux *http.ServeMux) {
	cl := auth.NewUserServiceClient(conn)
	_ = auth.RegisterUserServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	auth.RegisterUserServiceServer(server, a)
}

func (a api) ListUsers(ctx context.Context, request *auth.ListUsersRequest) (*auth.ListUsersResponse, error) {
	offset, limit := helpers.GetOffsetAndLimit(request.Offset, request.Limit)
	res, err := a.service.Query(ctx, request.Query, offset, limit)
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
	request.Password, _ = helpers.HashPassword(request.Password)
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

func (a api) AddUserRole(ctx context.Context, request *auth.AddUserRoleRequest) (*auth.User, error) {
	domainRole, err := a.service.AddUserRole(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return domainRole, nil
}

func (a api) UpdateUserRole(ctx context.Context, request *auth.UpdateUserRoleRequest) (*auth.User, error) {
	user, err := a.service.UpdateUserRole(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return user, nil
}

func (a api) DeleteUserRole(ctx context.Context, request *auth.DeleteUserRoleRequest) (*empty.Empty, error) {
	_, err := a.service.DeleteUserRole(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &empty.Empty{}, nil
}

func New(srv Service) API {
	s := api{service: srv}
	grpcgw.RegisterController(s)
	return s
}
