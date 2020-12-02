package apps

import (
	"context"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4/is"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-tire/auth/internal/entity"
	auth "github.com/golang-tire/auth/internal/proto/v1"
)

// Service encapsulates use case logic for apps.
type Service interface {
	GetApp(ctx context.Context, uuid string) (*auth.App, error)
	QueryApps(ctx context.Context, query string, offset, limit int64) (*auth.ListAppsResponse, error)
	CountApps(ctx context.Context) (int64, error)
	CreateApp(ctx context.Context, input *auth.CreateAppRequest) (*auth.App, error)
	UpdateApp(ctx context.Context, input *auth.UpdateAppRequest) (*auth.App, error)
	DeleteApp(ctx context.Context, uuid string) (*auth.App, error)

	GetResource(ctx context.Context, uuid string) (*auth.Resource, error)
	QueryResources(ctx context.Context, query string, offset, limit int64) (*auth.ListResourcesResponse, error)
	CountResources(ctx context.Context) (int64, error)
	CreateResource(ctx context.Context, input *auth.CreateResourceRequest) (*auth.Resource, error)
	UpdateResource(ctx context.Context, input *auth.UpdateResourceRequest) (*auth.Resource, error)
	DeleteResource(ctx context.Context, uuid string) (*auth.Resource, error)

	GetObject(ctx context.Context, uuid string) (*auth.Object, error)
	QueryObjects(ctx context.Context, query string, offset, limit int64) (*auth.ListObjectsResponse, error)
	CountObjects(ctx context.Context) (int64, error)
	CreateObject(ctx context.Context, input *auth.CreateObjectRequest) (*auth.Object, error)
	UpdateObject(ctx context.Context, input *auth.UpdateObjectRequest) (*auth.Object, error)
	DeleteObject(ctx context.Context, uuid string) (*auth.Object, error)
}

// ValidateAppCreateRequest validates the CreateAppRequest fields.
func ValidateAppCreateRequest(c *auth.CreateAppRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Enable, validation.Required),
	)
}

// ValidateAppUpdateRequest validates the UpdateAppRequest fields.
func ValidateAppUpdateRequest(u *auth.UpdateAppRequest) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Name, validation.Required, validation.Length(0, 128)),
		validation.Field(&u.Enable, validation.Required),
	)
}

// ValidateResourceCreateRequest validates the CreateResourceRequest fields.
func ValidateResourceCreateRequest(c *auth.CreateResourceRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Uuid, validation.Required, is.UUIDv4),
	)
}

// ValidateResourceUpdateRequest validates the UpdateResourceRequest fields.
func ValidateResourceUpdateRequest(u *auth.UpdateResourceRequest) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Name, validation.Required, validation.Length(0, 128)),
		validation.Field(&u.Uuid, validation.Required, is.UUIDv4),
	)
}

// ValidateObjectCreateRequest validates the CreateObjectRequest fields.
func ValidateObjectCreateRequest(c *auth.CreateObjectRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Identifier, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Uuid, validation.Required, is.UUIDv4),
	)
}

// ValidateObjectUpdateRequest validates the UpdateObjectRequest fields.
func ValidateObjectUpdateRequest(u *auth.UpdateObjectRequest) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Identifier, validation.Required, validation.Length(0, 128)),
		validation.Field(&u.Uuid, validation.Required, is.UUIDv4),
	)
}

type service struct {
	repo Repository
}

// NewService creates a new app service.
func NewService(repo Repository) Service {
	return service{repo}
}

// GetApp returns the app with the specified the app UUID.
func (s service) GetApp(ctx context.Context, UUID string) (*auth.App, error) {
	app, err := s.repo.GetApp(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return app.ToProto(), nil
}

// CreateApp creates a new app.
func (s service) CreateApp(ctx context.Context, req *auth.CreateAppRequest) (*auth.App, error) {
	if err := ValidateAppCreateRequest(req); err != nil {
		return nil, err
	}
	id, err := s.repo.CreateApp(ctx, entity.App{
		Name:   req.Name,
		Enable: req.Enable,
	})
	if err != nil {
		return nil, err
	}
	return s.GetApp(ctx, id)
}

// UpdateApp updates the app with the specified UUID.
func (s service) UpdateApp(ctx context.Context, req *auth.UpdateAppRequest) (*auth.App, error) {
	if err := ValidateAppUpdateRequest(req); err != nil {
		return nil, err
	}

	app, err := s.repo.GetApp(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	app.Name = req.Name
	app.Enable = req.Enable
	app.UpdatedAt = now

	if err := s.repo.UpdateApp(ctx, app); err != nil {
		return nil, err
	}
	return app.ToProto(), nil
}

// DeleteApp deletes the app with the specified UUID.
func (s service) DeleteApp(ctx context.Context, UUID string) (*auth.App, error) {
	app, err := s.repo.GetApp(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if err = s.repo.DeleteApp(ctx, app); err != nil {
		return nil, err
	}
	return app.ToProto(), nil
}

// CountApp returns the number of apps.
func (s service) CountApps(ctx context.Context) (int64, error) {
	return s.repo.CountApps(ctx)
}

// QueryApps returns the apps with the specified offset and limit.
func (s service) QueryApps(ctx context.Context, query string, offset, limit int64) (*auth.ListAppsResponse, error) {
	items, count, err := s.repo.QueryApps(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	return &auth.ListAppsResponse{
		Apps:       entity.AppToProtoList(items),
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}

// GetResource returns the app with the specified the resource UUID.
func (s service) GetResource(ctx context.Context, UUID string) (*auth.Resource, error) {
	app, err := s.repo.GetResource(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return app.ToProto(), nil
}

// CreateResource creates a new resource.
func (s service) CreateResource(ctx context.Context, req *auth.CreateResourceRequest) (*auth.Resource, error) {
	if err := ValidateResourceCreateRequest(req); err != nil {
		return nil, err
	}

	app, err := s.repo.GetApp(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	id, err := s.repo.CreateResource(ctx, app, entity.Resource{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	return s.GetResource(ctx, id)
}

// UpdateResource updates the resource with the specified UUID.
func (s service) UpdateResource(ctx context.Context, req *auth.UpdateResourceRequest) (*auth.Resource, error) {
	if err := ValidateResourceUpdateRequest(req); err != nil {
		return nil, err
	}

	resource, err := s.repo.GetResource(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	app, err := s.repo.GetApp(ctx, resource.App.UUID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	resource.Name = req.Name
	resource.UpdatedAt = now

	if err := s.repo.UpdateResource(ctx, app, resource); err != nil {
		return nil, err
	}
	return resource.ToProto(), nil
}

// DeleteResource deletes the resource with the specified UUID.
func (s service) DeleteResource(ctx context.Context, UUID string) (*auth.Resource, error) {
	app, err := s.repo.GetResource(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if err = s.repo.DeleteResource(ctx, app); err != nil {
		return nil, err
	}
	return app.ToProto(), nil
}

// CountResource returns the number of resources.
func (s service) CountResources(ctx context.Context) (int64, error) {
	return s.repo.CountResources(ctx)
}

// QueryResources returns the resources with the specified offset and limit.
func (s service) QueryResources(ctx context.Context, query string, offset, limit int64) (*auth.ListResourcesResponse, error) {
	items, count, err := s.repo.QueryResources(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	return &auth.ListResourcesResponse{
		Resources:  entity.ResourceToProtoList(items),
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}

// GetObject returns the object with the specified the object UUID.
func (s service) GetObject(ctx context.Context, UUID string) (*auth.Object, error) {
	object, err := s.repo.GetObject(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return object.ToProto(), nil
}

// CreateObject creates a new object.
func (s service) CreateObject(ctx context.Context, req *auth.CreateObjectRequest) (*auth.Object, error) {
	if err := ValidateObjectCreateRequest(req); err != nil {
		return nil, err
	}

	object, err := s.repo.GetApp(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	id, err := s.repo.CreateObject(ctx, object, entity.Object{
		Identifier: req.Identifier,
	})
	if err != nil {
		return nil, err
	}
	return s.GetObject(ctx, id)
}

// UpdateObject updates the object with the specified UUID.
func (s service) UpdateObject(ctx context.Context, req *auth.UpdateObjectRequest) (*auth.Object, error) {
	if err := ValidateObjectUpdateRequest(req); err != nil {
		return nil, err
	}

	object, err := s.repo.GetObject(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	app, err := s.repo.GetApp(ctx, object.App.UUID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	object.Identifier = req.Identifier
	object.UpdatedAt = now

	if err := s.repo.UpdateObject(ctx, app, object); err != nil {
		return nil, err
	}
	return object.ToProto(), nil
}

// DeleteObject deletes the object with the specified UUID.
func (s service) DeleteObject(ctx context.Context, UUID string) (*auth.Object, error) {
	object, err := s.repo.GetObject(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if err = s.repo.DeleteObject(ctx, object); err != nil {
		return nil, err
	}
	return object.ToProto(), nil
}

// CountObjects returns the number of objects.
func (s service) CountObjects(ctx context.Context) (int64, error) {
	return s.repo.CountObjects(ctx)
}

// QueryObjects returns the objects with the specified offset and limit.
func (s service) QueryObjects(ctx context.Context, query string, offset, limit int64) (*auth.ListObjectsResponse, error) {
	items, count, err := s.repo.QueryObjects(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	return &auth.ListObjectsResponse{
		Objects:    entity.ObjectToProtoList(items),
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}
