package auth

import (
	"fmt"
	"time"

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

type sessionData struct {
	User   *auth.User
	Tokens *auth.LoginResponse
}

// SaveTokens save user tokens after login
func SaveTokens(user *auth.User, response *auth.LoginResponse) error {
	var d = sessionData{
		User:   user,
		Tokens: response,
	}
	err := session.Set(response.AccessToken, &d, time.Minute*time.Duration(accessTokenLife.Int()))
	if err != nil {
		return err
	}
	return session.Set(response.RefreshToken, &d, time.Hour*time.Duration(refreshTokenLife.Int()))
}

// LoadTokens get tokens by access token
func LoadTokens(token string) (*sessionData, error) {
	var data sessionData
	err := session.Get(token, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func ExtractSessionUser(token string) (*auth.User, error) {
	data, err := LoadTokens(token)
	if err != nil {
		return nil, err
	}
	return data.User, nil
}

// DeleteToken delete user token
func DeleteToken(token string, isAccessToken bool) error {
	data, err := LoadTokens(token)
	if err != nil {
		return err
	}
	err = session.Delete(token)
	if err != nil {
		return err
	}
	if isAccessToken {
		return session.Delete(data.Tokens.RefreshToken)
	}
	return session.Delete(data.Tokens.AccessToken)
}

// CreateToken will create access and refresh token
func CreateToken(user *auth.User) (*auth.LoginResponse, error) {

	accessExpiresAt := time.Now().Add(time.Minute * time.Duration(accessTokenLife.Int())).Unix()
	refreshExpiresAt := time.Now().Add(time.Hour * time.Duration(refreshTokenLife.Int())).Unix()
	aToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_uuid":   user.Uuid,
		"access_uuid": user.Uuid,
		"exp":         accessExpiresAt,
	})

	accessToken, err := aToken.SignedString([]byte(jwtSecret.String()))
	if err != nil {
		return nil, err
	}

	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_uuid":    user.Uuid,
		"refresh_uuid": user.Uuid,
		"exp":          refreshExpiresAt,
	})

	refreshToken, err := rToken.SignedString([]byte(jwtSecret.String()))
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func VerifyToken(req string) (string, error) {
	_, err := jwt.Parse(req, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret.String()), nil
	})
	if err != nil {
		return "", err
	}
	return req, nil
}
