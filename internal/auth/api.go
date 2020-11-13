package auth

import (
	"context"

	"github.com/golang-tire/auth/internal/helpers"

	"github.com/golang-tire/pkg/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	userService users.Service
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
	hostname, err := ExtractHostName(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "hostname not found")
	}

	log.Info("hostname", log.String("hostname", hostname))

	user, err := a.userService.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "username %s not found", req.Username)
	}

	if !helpers.CheckPasswordHash(req.Password, user.Password) {
		return nil, status.Error(codes.Unauthenticated, "username or password is not valid")
	}

	tokens, err := CreateToken(user)
	if err != nil {
		log.Error("error on create user token", log.String("user", user.Username), log.Err(err))
		return nil, status.Errorf(codes.Internal, "internal server error, create token")
	}

	err = SaveTokens(user, tokens)
	if err != nil {
		log.Error("error on set user session", log.String("user", user.Username), log.Err(err))
		return nil, status.Errorf(codes.Internal, "internal server error, session")
	}

	log.Info("save token pass")
	return tokens, nil
}

func (a api) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {

	req.Password, _ = helpers.HashPassword(req.Password)
	user, err := a.userService.Create(ctx, &auth.CreateUserRequest{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Gender:    req.Gender,
		AvatarUrl: req.AvatarUrl,
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		RawData:   req.RawData,
		Enable:    true,
	})

	if err != nil {
		return nil, status.Errorf(codes.OutOfRange, "registration failed with %s", err.Error())
	}

	return &auth.RegisterResponse{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Gender:    user.Gender,
		AvatarUrl: user.AvatarUrl,
		Username:  user.Username,
		Email:     user.Email,
		RawData:   user.RawData,
	}, nil
}

func (a api) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	token, err := ExtractToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	err = DeleteToken(token, true /* isAccessToken */)
	if err != nil {
		log.Error("failed to remove token", log.Err(err))
		return nil, status.Errorf(codes.Internal, "logout failed")
	}
	return nil, nil
}

func (a api) VerifyToken(ctx context.Context, req *auth.VerifyTokenRequest) (*auth.VerifyTokenResponse, error) {
	vTokens, err := VerifyToken(req.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	return &auth.VerifyTokenResponse{AccessToken: vTokens}, nil
}

func (a api) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	vToken, err := VerifyToken(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token")
	}

	user, err := ExtractSessionUser(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "session expired")
	}

	dbUser, err := a.userService.GetByUsername(ctx, user.Username)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid user")
	}

	if !dbUser.Enable {
		return nil, status.Errorf(codes.Unauthenticated, "user is not active")
	}
	err = DeleteToken(vToken, false /* isAccessToken */)
	if err != nil {
		log.Error("failed to remove token", log.Err(err))
		return nil, status.Errorf(codes.Internal, "logout failed")
	}

	tokens, err := CreateToken(dbUser)
	if err != nil {
		log.Error("error on create user token", log.String("user", dbUser.Username), log.Err(err))
		return nil, status.Errorf(codes.Internal, "internal server error, create token")
	}

	return &auth.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func New(userService users.Service) API {
	s := api{userService: userService}
	grpcgw.RegisterController(s)
	InitMiddleware(userService)

	kv.Memory().SetString("/authV1.AuthService/Login", "open")
	kv.Memory().SetString("/authV1.AuthService/Register", "open")
	kv.Memory().SetString("/authV1.AuthService/VerifyToken", "open")
	kv.Memory().SetString("/authV1.AuthService/RefreshToken", "open")

	return s
}
