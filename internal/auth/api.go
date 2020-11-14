package auth

import (
	"context"

	"github.com/golang-tire/pkg/log"

	"github.com/golang-tire/auth/internal/rules"

	"github.com/golang-tire/pkg/kv"

	"github.com/golang-tire/auth/internal/users"

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
	auth.AuthServiceServer
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
	cl := auth.NewAuthServiceClient(conn)
	_ = auth.RegisterAuthServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	auth.RegisterAuthServiceServer(server, a)
}

func (a api) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	res, err := a.service.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a api) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	res, err := a.service.Register(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a api) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	res, err := a.service.Logout(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a api) VerifyToken(ctx context.Context, req *auth.VerifyTokenRequest) (*auth.VerifyTokenResponse, error) {
	res, err := a.service.VerifyToken(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a api) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	res, err := a.service.RefreshToken(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func New(ctx context.Context, srv Service, rulesService rules.Service, userService users.Service) (API, error) {
	s := api{service: srv}
	grpcgw.RegisterController(s)

	enforcer, err := InitRbac(ctx, rulesService, userService)
	if err != nil {
		return nil, err
	}
	log.Info("load rbac polices...")
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}
	InitMiddleware(userService, enforcer)
	kv.Memory().SetString("/authV1.AuthService/Login", "open")
	kv.Memory().SetString("/authV1.AuthService/Register", "open")
	kv.Memory().SetString("/authV1.AuthService/VerifyToken", "open")
	kv.Memory().SetString("/authV1.AuthService/RefreshToken", "open")
	return s, nil
}
