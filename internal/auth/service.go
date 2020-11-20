package auth

import (
	"context"

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

// Service encapsulates use case logic for auth.
type Service interface {
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
	Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error)
	VerifyToken(ctx context.Context, req *auth.VerifyTokenRequest) (*auth.VerifyTokenResponse, error)
	RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error)
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

// NewService creates a new auth service.
func NewService(userService users.Service) Service {
	return service{userService}
}
