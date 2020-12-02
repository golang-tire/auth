package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-tire/auth/internal/apps"

	"github.com/golang-tire/auth/internal/audit_logs"

	"github.com/golang-tire/pkg/pubsub"

	"github.com/go-redis/redis/v8"

	"github.com/golang-tire/pkg/kv"

	"github.com/golang-tire/auth/internal/auth"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/golang-tire/auth/internal/domains"

	"github.com/golang-tire/auth/internal/users"

	"github.com/golang-tire/auth/internal/rules"

	"github.com/golang-tire/pkg/config"
	"github.com/golang-tire/pkg/grpcgw"

	"github.com/golang-tire/auth/internal/entity"

	"github.com/golang-tire/auth/internal/roles"

	"github.com/golang-tire/auth/internal/pkg/db"
	_ "github.com/golang-tire/auth/internal/pkg/helpers"
)

const (
	defaultHTTPPort       = 8089
	defaultGrpcPort       = 9090
	defaultSwaggerBaseURL = "/v1/swagger"
)

var (
	httpPort       = config.RegisterInt("server.httpPort", defaultHTTPPort)
	grpcPort       = config.RegisterInt("server.grpcPort", defaultGrpcPort)
	swaggerBaseURL = config.RegisterString("server.swaggerBaseURL", defaultSwaggerBaseURL)

	redisHost     = config.RegisterString("redis.host", "localhost")
	redisPort     = config.RegisterInt("redis.port", 6379)
	redisDb       = config.RegisterInt64("redis.db", 1)
	redisPassword = config.RegisterString("redis.password", "")
)

func setupModules(ctx context.Context) error {

	// reload configs
	err := config.Load()
	if err != nil {
		return err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisHost.String(), redisPort.Int()),
		Password: redisPassword.String(),
		DB:       redisDb.Int(),
	})

	_, err = kv.InitWithConn(ctx, rdb)
	if err != nil {
		return err
	}

	pubsub.New(rdb)

	dbInstance, err := db.Init(ctx)
	if err != nil {
		return err
	}

	models := []interface{}{
		&entity.Domain{},
		&entity.Role{},
		&entity.Rule{},
		&entity.User{},
		&entity.UserRole{},
		&entity.AuditLog{},
		&entity.App{},
		&entity.Resource{},
		&entity.Object{},
	}

	err = db.CreateSchema(dbInstance.DB(), models)
	if err != nil {
		return err
	}

	domainsRepo := domains.NewRepository(dbInstance)
	domainsSrv := domains.NewService(domainsRepo)
	domains.New(domainsSrv)

	rolesRepo := roles.NewRepository(dbInstance)
	rolesSrv := roles.NewService(rolesRepo)
	roles.New(rolesSrv)

	rulesRepo := rules.NewRepository(dbInstance)
	rulesSrv := rules.NewService(rulesRepo, domainsRepo, rolesRepo)
	rules.New(rulesSrv)

	usersRepo := users.NewRepository(dbInstance)
	usersSrv := users.NewService(usersRepo, domainsRepo, rolesRepo)
	users.New(usersSrv)

	auditLogRepo := audit_logs.NewRepository(dbInstance)
	auditLogSrv := audit_logs.NewService(auditLogRepo, usersRepo)
	audit_logs.New(auditLogSrv)

	appsRepo := apps.NewRepository(dbInstance)
	appsSrv := apps.NewService(appsRepo)
	apps.New(appsSrv)

	rbacSrv, err := auth.InitRbac(ctx, rulesSrv, usersSrv)
	if err != nil {
		return err
	}

	authService := auth.NewService(usersSrv, rbacSrv)
	_, err = auth.New(ctx, authService, rulesSrv, usersSrv)
	if err != nil {
		return err
	}

	jsonpb := &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}

	err = grpcgw.Serve(ctx,
		grpcgw.GrpcPort(grpcPort.Int()),
		grpcgw.HttpPort(httpPort.Int()),
		grpcgw.SwaggerBaseURL(swaggerBaseURL.String()),
		grpcgw.ServeMuxOptions(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
			runtime.WithIncomingHeaderMatcher(grpcHeaderMatcher),
			runtime.WithOutgoingHeaderMatcher(grpcHeaderMatcher),
		),
	)

	return err
}

func grpcHeaderMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "x-forwarded-uri", "x-forwarded-method", "x-auth-user-email", "x-auth-user-uuid", "x-auth-user-name":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
