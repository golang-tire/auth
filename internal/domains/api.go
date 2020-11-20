package domains

import (
	"context"
	"net/http"

	"github.com/golang-tire/auth/internal/pkg/helpers"
	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang-tire/pkg/grpcgw"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API interface {
	grpcgw.Controller
}

type api struct {
	service Service
	auth.DomainServiceServer
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux, httpMux *http.ServeMux) {
	cl := auth.NewDomainServiceClient(conn)
	_ = auth.RegisterDomainServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	auth.RegisterDomainServiceServer(server, a)
}

func (a api) ListDomains(ctx context.Context, request *auth.ListDomainsRequest) (*auth.ListDomainsResponse, error) {
	offset, limit := helpers.GetOffsetAndLimit(request.Offset, request.Limit)
	res, err := a.service.Query(ctx, offset, limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) GetDomain(ctx context.Context, request *auth.GetDomainRequest) (*auth.Domain, error) {
	res, err := a.service.Get(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) CreateDomain(ctx context.Context, request *auth.CreateDomainRequest) (*auth.Domain, error) {
	res, err := a.service.Create(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) UpdateDomain(ctx context.Context, request *auth.UpdateDomainRequest) (*auth.Domain, error) {
	res, err := a.service.Update(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) DeleteDomain(ctx context.Context, request *auth.DeleteDomainRequest) (*empty.Empty, error) {
	_, err := a.service.Delete(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &empty.Empty{}, err
}

func New(srv Service) API {
	s := api{service: srv}
	grpcgw.RegisterController(s)
	return s
}
