package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-tire/pkg/log"

	auth "github.com/golang-tire/auth/internal/proto/v1"
)

const (
	xForwardedURI    = "X-Forwarded-URI"
	xForwardedHost   = "X-Forwarded-Host"
	xForwardedMethod = "X-Forwarded-Method"

	authorizationHeader = "authorization"
	xAuthUsername       = "x-auth-user-name"
	xAuthUserEmail      = "x-auth-user-email"
	xAuthUserUuid       = "x-auth-user-uuid"
)

func (a api) CheckCredential(w http.ResponseWriter, r *http.Request) {

	// check for method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("method not allowed"))
		return
	}

	// check for xForwardedURI
	forwardedURI, err := getXForwardedURI(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	// check for xForwardedMethod
	forwardedMethod, err := getXForwardedMethod(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	// check for XForwardedHost
	forwardedHost, err := getXForwardedHost(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	// check for token
	token, statusCode, err := getToken(r)
	if err != nil {
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	vToken, err := extractTokenData(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("invalid token"))
		return
	}

	// get token data
	td, err := loadTokenDetails(*vToken.AccessUuid)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("session expired"))
		return
	}

	// get user
	user, err := a.usersSrv.Get(a.ctx, td.UserUuid)
	if err != nil || !user.Enable {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("invalid credential"))
		return
	}

	// check for rbac
	ok, err := a.checkRbac(forwardedURI, forwardedHost, forwardedMethod, user)
	if err != nil {
		log.Error("check rbac permission failed", log.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("check permission failed"))
		return
	}
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("forbidden"))
		return
	}

	w.Header().Set(xAuthUsername, user.Username)
	w.Header().Set(xAuthUserEmail, user.Email)
	w.Header().Set(xAuthUserUuid, user.Uuid)
	w.WriteHeader(http.StatusOK)
}

func (a api) checkRbac(uri, domain, method string, user *auth.User) (bool, error) {
	resource, object, err := a.parseURI(uri)
	if err != nil {
		log.Error("parse uri failed", log.Err(err))
		return false, errors.New("parse forwarded uri failed")
	}

	return a.rbac.enforcer.Enforce(user.Username, domain, resource, method, object)
}

func (a api) parseURI(uri string) (string, string, error) {

	var resource, object string
	for _, p := range a.rbac.regexPatterns {
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

func getToken(r *http.Request) (string, int, error) {
	val := r.Header.Get(authorizationHeader)
	if val == "" {
		return "", http.StatusNonAuthoritativeInfo, errors.New("authoritative info not provided")
	}

	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", http.StatusUnauthorized, errors.New("bad authorization string")
	}

	if !strings.EqualFold(splits[0], "bearer") {
		return "", http.StatusUnauthorized, errors.New("request unauthenticated with bearer")
	}
	return splits[1], 0, nil
}

func getXForwardedMethod(r *http.Request) (string, error) {
	val := r.Header.Get(xForwardedMethod)
	if val == "" {
		return "", errors.New(xForwardedMethod + " header is required")
	}
	return val, nil
}

func getXForwardedHost(r *http.Request) (string, error) {
	val := r.Header.Get(xForwardedHost)
	if val == "" {
		return "", errors.New(xForwardedHost + " header is required")
	}
	return val, nil
}

func getXForwardedURI(r *http.Request) (string, error) {
	val := r.Header.Get(xForwardedURI)
	if val == "" {
		return "", errors.New(xForwardedURI + " header is required")
	}
	return val, nil
}
