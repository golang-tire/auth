package auth

import (
	"context"
	"errors"

	"github.com/casbin/casbin/v2"

	"github.com/golang-tire/pkg/kv"

	auth "github.com/golang-tire/auth/internal/proto/v1"

	"github.com/golang-tire/auth/internal/users"
	"github.com/golang-tire/pkg/grpcgw"

	"google.golang.org/grpc/metadata"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

type contextKey int

const (
	resourceKey contextKey = iota
	userKey
	tokenKey
	fullMethodKey
	hostNameKey
)

type Middleware struct {
	enforcer    *casbin.Enforcer
	userService users.Service
}

func streamExtractor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return handler(srv, ss)
}

func unaryExtractor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.InvalidArgument, "metadata is not readable!")
	}

	// TODO find a better way to read x-forward-host and authority field
	forwardedHost := m.Get("x-forwarded-host")
	if len(forwardedHost) == 1 {
		// its came from grpc-gateway
		ctx = context.WithValue(ctx, hostNameKey, forwardedHost[0])
	} else {
		// its a grpc call
		ctx = context.WithValue(ctx, hostNameKey, m.Get(":authority")[0])
	}

	res, ok := kv.Memory().Get(info.FullMethod)
	if !ok {
		ctx = context.WithValue(ctx, resourceKey, res)
	}
	ctx = context.WithValue(ctx, fullMethodKey, info.FullMethod)
	return handler(ctx, req)
}

func (m Middleware) authHandler(ctx context.Context) (context.Context, error) {
	r := ctx.Value(resourceKey)
	if r == nil { // No user requested here
		return ctx, nil
	}
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return ctx, status.Errorf(codes.InvalidArgument, "token required")
	}

	vToken, err := extractTokenData(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	td, err := loadTokenDetails(*vToken.AccessUuid)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "session expired")
	}

	user, err := m.userService.Get(ctx, td.UserUuid)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "invalid credential")
	}

	if !user.Enable {
		return ctx, status.Errorf(codes.Unauthenticated, "access denied")
	}

	return context.WithValue(context.WithValue(ctx, userKey, user), tokenKey, token), nil
}

// ExtractUser try to extract the current user from the context
func ExtractUser(ctx context.Context) (*auth.User, error) {
	u, ok := ctx.Value(userKey).(*auth.User)
	if !ok {
		return nil, errors.New("no user in context")
	}
	return u, nil
}

// ExtractToken try to extract token from context
func ExtractToken(ctx context.Context) (string, error) {
	tok, ok := ctx.Value(tokenKey).(string)
	if !ok {
		return "", errors.New("no token in context")
	}
	return tok, nil
}

// ExtractHostName try to extract hostname from context
func ExtractHostName(ctx context.Context) (string, error) {
	tok, ok := ctx.Value(hostNameKey).(string)
	if !ok {
		return "", errors.New("no hostname in context")
	}
	return tok, nil
}

func InitMiddleware(userService users.Service) {

	middleware := Middleware{userService: userService}
	grpcgw.RegisterInterceptors(grpcgw.Interceptor{
		Unary:  unaryExtractor,
		Stream: streamExtractor,
	})

	grpcgw.RegisterInterceptors(grpcgw.Interceptor{
		Unary:  grpc_auth.UnaryServerInterceptor(middleware.authHandler),
		Stream: grpc_auth.StreamServerInterceptor(middleware.authHandler),
	})
}
