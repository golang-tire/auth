package main

import (
	"context"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/golang-tire/auth/internal/domains"

	"github.com/golang-tire/auth/internal/users"

	"github.com/golang-tire/auth/internal/rules"

	"github.com/golang-tire/pkg/config"
	"github.com/golang-tire/pkg/grpcgw"

	"github.com/golang-tire/auth/internal/entity"

	"github.com/golang-tire/auth/internal/roles"

	"github.com/golang-tire/auth/internal/db"
	_ "github.com/golang-tire/auth/internal/helpers"
)

const (
	defaultHttpPort       = 8089
	defaultGrpcPort       = 9090
	defaultSwaggerBaseURL = "/v1/swagger"
)

var (
	httpPort       = config.RegisterInt("server.httpPort", defaultHttpPort)
	grpcPort       = config.RegisterInt("server.grpcPort", defaultGrpcPort)
	swaggerBaseURL = config.RegisterString("server.swaggerBaseURL", defaultSwaggerBaseURL)
)

func setupModules(ctx context.Context) error {

	// reload configs
	err := config.Load()
	if err != nil {
		return err
	}

	dbInstance, err := db.Init(ctx)
	if err != nil {
		return err
	}

	models := []interface{}{
		&entity.Domain{},
		&entity.Role{},
		&entity.Rule{},
		&entity.User{},
		&entity.DomainRole{},
		&entity.UserRule{},
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
	rulesSrv := rules.NewService(rulesRepo, domainsRepo)
	rules.New(rulesSrv)

	usersRepo := users.NewRepository(dbInstance)
	usersSrv := users.NewService(usersRepo, rulesRepo)
	users.New(usersSrv)

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
		),
	)

	return err
}
