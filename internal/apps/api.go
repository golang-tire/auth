package apps

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
	auth.AppServiceServer
	service Service
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux, httpMux *http.ServeMux) {
	cl := auth.NewAppServiceClient(conn)
	_ = auth.RegisterAppServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	auth.RegisterAppServiceServer(server, a)
}

func (a api) ListApps(ctx context.Context, request *auth.ListAppsRequest) (*auth.ListAppsResponse, error) {
	offset, limit := helpers.GetOffsetAndLimit(request.Offset, request.Limit)
	res, err := a.service.QueryApps(ctx, request.Query, offset, limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) GetApp(ctx context.Context, request *auth.GetAppRequest) (*auth.App, error) {
	res, err := a.service.GetApp(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) CreateApp(ctx context.Context, request *auth.CreateAppRequest) (*auth.App, error) {
	res, err := a.service.CreateApp(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) UpdateApp(ctx context.Context, request *auth.UpdateAppRequest) (*auth.App, error) {
	res, err := a.service.UpdateApp(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) DeleteApp(ctx context.Context, request *auth.DeleteAppRequest) (*empty.Empty, error) {
	_, err := a.service.DeleteApp(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &empty.Empty{}, err
}

func (a api) ListResources(ctx context.Context, request *auth.ListResourcesRequest) (*auth.ListResourcesResponse, error) {
	offset, limit := helpers.GetOffsetAndLimit(request.Offset, request.Limit)
	res, err := a.service.QueryResources(ctx, request.Query, offset, limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) GetResource(ctx context.Context, request *auth.GetResourceRequest) (*auth.Resource, error) {
	res, err := a.service.GetResource(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) CreateResource(ctx context.Context, request *auth.CreateResourceRequest) (*auth.Resource, error) {
	res, err := a.service.CreateResource(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) UpdateResource(ctx context.Context, request *auth.UpdateResourceRequest) (*auth.Resource, error) {
	res, err := a.service.UpdateResource(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) DeleteResource(ctx context.Context, request *auth.DeleteResourceRequest) (*empty.Empty, error) {
	_, err := a.service.DeleteResource(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &empty.Empty{}, err
}

func (a api) ListObjects(ctx context.Context, request *auth.ListObjectsRequest) (*auth.ListObjectsResponse, error) {
	offset, limit := helpers.GetOffsetAndLimit(request.Offset, request.Limit)
	res, err := a.service.QueryObjects(ctx, request.Query, offset, limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) GetObject(ctx context.Context, request *auth.GetObjectRequest) (*auth.Object, error) {
	res, err := a.service.GetObject(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) CreateObject(ctx context.Context, request *auth.CreateObjectRequest) (*auth.Object, error) {
	res, err := a.service.CreateObject(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) UpdateObject(ctx context.Context, request *auth.UpdateObjectRequest) (*auth.Object, error) {
	res, err := a.service.UpdateObject(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) DeleteObject(ctx context.Context, request *auth.DeleteObjectRequest) (*empty.Empty, error) {
	_, err := a.service.DeleteObject(ctx, request.Uuid)
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
