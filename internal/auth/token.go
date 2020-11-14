package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/golang-tire/pkg/config"
	"github.com/golang-tire/pkg/session"

	auth "github.com/golang-tire/auth/internal/proto/v1"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtSecret        = config.RegisterString("auth.jwtSecret", "this-is-for-test-dont-use-in-production")
	accessTokenLife  = config.RegisterInt("auth.accessTokenLife", 15)
	refreshTokenLife = config.RegisterInt("auth.refreshTokenLife", 170)
)

type tokenData struct {
	UserUuid      string
	IsAccessToken bool
	AccessUuid    *string
	RefreshUuid   *string
}

type tokenDetails struct {
	AccessExpireAt  int64
	RefreshExpireAt int64
	AccessToken     string
	RefreshToken    string
	AccessUuid      string
	RefreshUuid     string
	UserUuid        string
	Username        string
}

// saveTokens save user tokens after login
func saveTokens(tokens *tokenDetails) error {
	err := session.Set(tokens.AccessUuid, tokens, time.Minute*time.Duration(accessTokenLife.Int()))
	if err != nil {
		return err
	}
	return session.Set(tokens.RefreshUuid, tokens, time.Hour*time.Duration(refreshTokenLife.Int()))
}

// loadTokenDetails get tokens by access token
func loadTokenDetails(token string) (*tokenDetails, error) {
	var data tokenDetails
	err := session.Get(token, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// deleteToken delete user token
func deleteToken(data *tokenData) error {
	if data.IsAccessToken {
		return session.Delete(*data.AccessUuid)
	}
	return session.Delete(*data.RefreshUuid)
}

func extractTokenData(token string) (*tokenData, error) {

	td := &tokenData{
		UserUuid:      "",
		IsAccessToken: false,
		AccessUuid:    nil,
		RefreshUuid:   nil,
	}
	vToken, err := verifyToken(token)
	if err != nil {
		return nil, err
	}

	var invalidErr = fmt.Errorf("invalid token structure")
	claims, ok := vToken.Claims.(jwt.MapClaims)
	if !ok || !vToken.Valid {
		return nil, invalidErr
	}

	if v, found := claims["access_uuid"]; found {
		td.IsAccessToken = true

		ac, ok := v.(string)
		if !ok {
			return nil, invalidErr
		}
		td.AccessUuid = &ac
	}

	if v, found := claims["refresh_uuid"]; found {
		td.IsAccessToken = false

		rf, ok := v.(string)
		if !ok {
			return nil, invalidErr
		}
		td.RefreshUuid = &rf
	}

	if td.AccessUuid == nil && td.RefreshUuid == nil {
		return nil, invalidErr
	}

	if v, found := claims["user_uuid"]; found {
		us, ok := v.(string)
		if !ok {
			return nil, invalidErr
		}
		td.UserUuid = us
	}

	return td, nil
}

// createToken will create access and refresh token
func createToken(user *auth.User) (*tokenDetails, error) {

	td := &tokenDetails{}
	td.AccessExpireAt = time.Now().Add(time.Minute * time.Duration(accessTokenLife.Int())).Unix()
	td.RefreshExpireAt = time.Now().Add(time.Hour * time.Duration(refreshTokenLife.Int())).Unix()

	td.AccessUuid = uuid.New().String()
	td.RefreshUuid = uuid.New().String()
	td.UserUuid = user.Uuid
	td.Username = user.Username

	aToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_uuid":   user.Uuid,
		"access_uuid": td.AccessUuid,
		"exp":         td.AccessExpireAt,
	})

	var err error
	td.AccessToken, err = aToken.SignedString([]byte(jwtSecret.String()))
	if err != nil {
		return nil, err
	}

	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_uuid":    user.Uuid,
		"refresh_uuid": td.RefreshUuid,
		"exp":          td.RefreshExpireAt,
	})

	td.RefreshToken, err = rToken.SignedString([]byte(jwtSecret.String()))
	if err != nil {
		return nil, err
	}

	return td, nil
}

// verifyToken verify if a token is valid, will return token object
func verifyToken(req string) (*jwt.Token, error) {
	token, err := jwt.Parse(req, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret.String()), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
