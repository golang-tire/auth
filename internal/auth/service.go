package auth

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"

	"google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/golang-tire/pkg/session"

	"github.com/go-ozzo/ozzo-validation/v4/is"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-tire/auth/internal/pkg/helpers"
	"github.com/golang-tire/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang-tire/auth/internal/users"

	auth "github.com/golang-tire/auth/internal/proto/v1"
)

const (
	xForwardedURI    = "x-forwarded-uri"
	xForwardedHost   = "x-forwarded-host"
	xForwardedMethod = "x-forwarded-method"

	//authorizationHeader = "authorization"
	xAuthUsername  = "x-auth-user-name"
	xAuthUserEmail = "x-auth-user-email"
	xAuthUserUuid  = "x-auth-user-uuid"
)

type Headers struct {
	ForwardedHost   string
	ForwardedURI    string
	ForwardedMethod string
}

// Service encapsulates use case logic for auth.
type Service interface {
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
	Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error)
	VerifyToken(ctx context.Context, req *auth.VerifyTokenRequest) (*auth.VerifyTokenResponse, error)
	RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error)
	Validate(ctx context.Context, req *auth.ValidateRequest) (*empty.Empty, error)
}

// ValidateLoginRequest validates the LoginRequest fields.
func ValidateLoginRequest(c *auth.LoginRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Username, validation.Required, validation.Length(6, 128)),
		validation.Field(&c.Password, validation.Required, validation.Length(6, 128)),
	)
}

// ValidateRegisterRequest validates the RegisterRequest fields.
func ValidateRegisterRequest(c *auth.RegisterRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Username, validation.Required, validation.Length(6, 128)),
		validation.Field(&c.Password, validation.Required, validation.Length(6, 128)),
		validation.Field(&c.Email, validation.Required, is.Email),
	)
}

type service struct {
	rbac        *rbacService
	userService users.Service
}

func (s service) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {

	if err := ValidateLoginRequest(req); err != nil {
		return nil, err
	}

	hostname, err := ExtractHostName(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "hostname not found")
	}

	log.Info("hostname", log.String("hostname", hostname))

	user, err := s.userService.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "username %s not found", req.Username)
	}

	if !helpers.CheckPasswordHash(req.Password, user.Password) {
		return nil, status.Error(codes.Unauthenticated, "username or password is not valid")
	}

	tokens, err := createToken(user)
	if err != nil {
		log.Error("error on create user token", log.String("user", user.Username), log.Err(err))
		return nil, status.Errorf(codes.Internal, "internal server error, create token")
	}

	err = saveTokens(tokens)
	if err != nil {
		log.Error("error on set user session", log.String("user", user.Username), log.Err(err))
		return nil, status.Errorf(codes.Internal, "internal server error, session")
	}

	return &auth.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s service) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {

	if err := ValidateRegisterRequest(req); err != nil {
		return nil, err
	}

	req.Password, _ = helpers.HashPassword(req.Password)
	user, err := s.userService.Create(ctx, &auth.CreateUserRequest{
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

func (s service) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	token, err := ExtractToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token not found")
	}

	vToken, err := extractTokenData(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	err = session.Delete(*vToken.AccessUuid)
	if err != nil {
		log.Error("failed to remove token", log.Err(err))
		return nil, status.Errorf(codes.InvalidArgument, "session already expired")
	}
	return &auth.LogoutResponse{RedirectTo: "/v1/auth/login"}, nil
}

func (s service) VerifyToken(ctx context.Context, req *auth.VerifyTokenRequest) (*auth.VerifyTokenResponse, error) {
	_, err := verifyToken(req.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	return &auth.VerifyTokenResponse{AccessToken: req.AccessToken}, nil
}

func (s service) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	vToken, err := extractTokenData(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token")
	}

	td, err := loadTokenDetails(*vToken.RefreshUuid)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "refresh token expired")
	}

	dbUser, err := s.userService.GetByUsername(ctx, td.Username)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid user")
	}

	if !dbUser.Enable {
		return nil, status.Errorf(codes.Unauthenticated, "user is not active")
	}

	err = deleteToken(vToken)
	if err != nil {
		log.Error("failed to remove token", log.Err(err))
		return nil, status.Errorf(codes.Internal, "logout failed")
	}

	tokens, err := createToken(dbUser)
	if err != nil {
		log.Error("error on create user token", log.String("user", dbUser.Username), log.Err(err))
		return nil, status.Errorf(codes.Internal, "internal server error, create token")
	}

	err = saveTokens(tokens)
	if err != nil {
		log.Error("error on set user session", log.String("user", dbUser.Username), log.Err(err))
		return nil, status.Errorf(codes.Internal, "internal server error, session")
	}

	return &auth.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s service) Validate(ctx context.Context, req *auth.ValidateRequest) (*empty.Empty, error) {

	headers, err := getHeaders(ctx)
	if err != nil {
		return &empty.Empty{}, err
	}
	user := ctx.Value(userKey).(*auth.User)

	// check for rbac
	ok, err := s.checkRbac(headers.ForwardedURI, headers.ForwardedHost, headers.ForwardedMethod, user)
	if err != nil {
		log.Error("check rbac permission failed", log.Err(err))
		return &empty.Empty{}, status.Errorf(codes.Internal, "check permission failed")
	}

	if !ok {
		return &empty.Empty{}, status.Errorf(codes.PermissionDenied, "forbidden")
	}

	err = grpc.SendHeader(ctx, metadata.New(map[string]string{
		xAuthUsername:  user.Username,
		xAuthUserEmail: user.Email,
		xAuthUserUuid:  user.Uuid,
	}))

	if err != nil {
		return &empty.Empty{}, status.Errorf(codes.Internal, "complete auth process failed")
	}
	return &empty.Empty{}, nil
}

func (s service) checkRbac(uri, domain, method string, user *auth.User) (bool, error) {
	resource, object, err := s.parseURI(uri)
	if err != nil {
		log.Error("parse uri failed", log.Err(err))
		return false, errors.New("parse forwarded uri failed")
	}

	return s.rbac.enforcer.Enforce(user.Username, domain, resource, method, object)
}

func (s service) parseURI(uri string) (string, string, error) {

	var resource, object string
	for _, p := range s.rbac.regexPatterns {
		var n = 0
		s := p.String()
		if strings.Contains(s, "resource") {
			n++
		}
		if strings.Contains(s, "object") {
			n++
		}

		res := p.FindStringSubmatch(uri)
		names := p.SubexpNames()

		if len(res)-1 != n {
			continue
		}

		for i := range res {
			if i == 0 {
				continue
			}

			if names[i] == "resource" {
				resource = res[i]
			}

			if names[i] == "object" {
				object = res[i]
			}
		}
	}
	if resource == "" && object == "" {
		return "", "", errors.New("resource or object not found")
	}

	if resource == "" {
		resource = "*"
	}

	if object == "" {
		object = "*"
	}

	return resource, object, nil
}

func getHeaders(ctx context.Context) (*Headers, error) {
	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "metadata is not readable!")
	}

	headers := &Headers{
		ForwardedHost:   m.Get(xForwardedHost)[0],
		ForwardedURI:    m.Get(xForwardedURI)[0],
		ForwardedMethod: m.Get(xForwardedMethod)[0],
	}

	if headers.ForwardedHost == "" {
		return nil, status.Errorf(codes.InvalidArgument, xForwardedHost+" is required ")
	}

	if headers.ForwardedURI == "" {
		return nil, status.Errorf(codes.InvalidArgument, xForwardedURI+" is required ")
	}

	if headers.ForwardedMethod == "" {
		return nil, status.Errorf(codes.InvalidArgument, xForwardedMethod+" is required ")
	}

	return headers, nil
}

// NewService creates a new auth service.
func NewService(userService users.Service, rbac *rbacService) Service {
	return service{rbac, userService}
}
