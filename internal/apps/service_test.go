package apps

import (
	"context"
	"testing"

	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/stretchr/testify/assert"
)

func TestCreateAppRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.CreateAppRequest
		wantError bool
	}{
		{"success", auth.CreateAppRequest{Name: "test", Enable: true}, false},
		{"required", auth.CreateAppRequest{Name: "", Enable: true}, true},
		{"too long", auth.CreateAppRequest{Enable: true, Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAppCreateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateAppRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.UpdateAppRequest
		wantError bool
	}{
		{"success", auth.UpdateAppRequest{Name: "test", Enable: true}, false},
		{"required", auth.UpdateAppRequest{Name: "", Enable: true}, true},
		{"too long", auth.UpdateAppRequest{Enable: true, Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAppUpdateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestCreateResourceRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.CreateResourceRequest
		wantError bool
	}{
		{"success", auth.CreateResourceRequest{Name: "test"}, false},
		{"required", auth.CreateResourceRequest{Name: ""}, true},
		{"too long", auth.CreateResourceRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateResourceCreateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateResourceRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.UpdateResourceRequest
		wantError bool
	}{
		{"success", auth.UpdateResourceRequest{Name: "test"}, false},
		{"required", auth.UpdateResourceRequest{Name: ""}, true},
		{"too long", auth.UpdateResourceRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateResourceUpdateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestCreateObjectRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.CreateObjectRequest
		wantError bool
	}{
		{"success", auth.CreateObjectRequest{Identifier: "test"}, false},
		{"required", auth.CreateObjectRequest{Identifier: ""}, true},
		{"too long", auth.CreateObjectRequest{Identifier: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateObjectCreateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateObjectRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.UpdateObjectRequest
		wantError bool
	}{
		{"success", auth.UpdateObjectRequest{Identifier: "test"}, false},
		{"required", auth.UpdateObjectRequest{Identifier: ""}, true},
		{"too long", auth.UpdateObjectRequest{Identifier: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateObjectUpdateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	s := NewService(&mockRepository{})
	ctx := context.Background()

	// initial count
	count, _ := s.CountApps(ctx)
	assert.Equal(t, int64(0), count)

	// successful creation
	app, err := s.CreateApp(ctx, &auth.CreateAppRequest{Name: "test", Enable: true})
	assert.Nil(t, err)
	assert.NotEmpty(t, app.Uuid)
	id := app.Uuid
	assert.Equal(t, "test", app.Name)
	assert.NotEmpty(t, app.CreatedAt)
	assert.NotEmpty(t, app.UpdatedAt)
	count, _ = s.CountApps(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.CreateApp(ctx, &auth.CreateAppRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.CountApps(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.CreateApp(ctx, &auth.CreateAppRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.CountApps(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.CreateApp(ctx, &auth.CreateAppRequest{Name: "test2"})

	// update
	app, err = s.UpdateApp(ctx, &auth.UpdateAppRequest{Name: "test updated", Uuid: id})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", app.Name)
	_, err = s.UpdateApp(ctx, &auth.UpdateAppRequest{Name: "test updated", Uuid: "none"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.UpdateApp(ctx, &auth.UpdateAppRequest{Name: "", Uuid: id})
	assert.NotNil(t, err)
	count, _ = s.CountApps(ctx)
	assert.Equal(t, int64(2), count)

	// unexpected error in update
	_, err = s.UpdateApp(ctx, &auth.UpdateAppRequest{Name: "error", Uuid: id})
	assert.Equal(t, errCRUD, err)
	count, _ = s.CountApps(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.GetApp(ctx, "none")
	assert.NotNil(t, err)
	app, err = s.GetApp(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", app.Name)
	assert.Equal(t, id, app.Uuid)

	// query
	_apps, _ := s.QueryApps(ctx, "", 0, 0)
	assert.Equal(t, 2, int(_apps.TotalCount))

	// delete
	_, err = s.DeleteApp(ctx, "none")
	assert.NotNil(t, err)
	app, err = s.DeleteApp(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, app.Uuid)
	count, _ = s.CountApps(ctx)
	assert.Equal(t, int64(1), count)

	// resources
	// initial count
	resourceCount, _ := s.CountResources(ctx)
	assert.Equal(t, int64(0), resourceCount)

	// successful creation
	resource, err := s.CreateResource(ctx, &auth.CreateResourceRequest{Name: "test", Uuid: app.Uuid})
	assert.Nil(t, err)
	assert.NotEmpty(t, resource.Uuid)
	resourceId := resource.Uuid
	assert.Equal(t, "test", resource.Name)
	assert.NotEmpty(t, resource.CreatedAt)
	assert.NotEmpty(t, resource.UpdatedAt)
	resourceCount, _ = s.CountResources(ctx)
	assert.Equal(t, int64(1), resourceCount)

	// validation error in creation
	_, err = s.CreateResource(ctx, &auth.CreateResourceRequest{Name: "", Uuid: app.Uuid})
	assert.NotNil(t, err)
	resourceCount, _ = s.CountResources(ctx)
	assert.Equal(t, int64(1), resourceCount)

	// unexpected error in creation
	_, err = s.CreateResource(ctx, &auth.CreateResourceRequest{Name: "error", Uuid: app.Uuid})
	assert.Equal(t, errCRUD, err)
	resourceCount, _ = s.CountResources(ctx)
	assert.Equal(t, int64(1), resourceCount)

	_, _ = s.CreateResource(ctx, &auth.CreateResourceRequest{Name: "test2", Uuid: app.Uuid})

	// update
	resource, err = s.UpdateResource(ctx, &auth.UpdateResourceRequest{Name: "test updated", Uuid: resourceId})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", resource.Name)
	_, err = s.UpdateResource(ctx, &auth.UpdateResourceRequest{Name: "test updated", Uuid: "none"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.UpdateResource(ctx, &auth.UpdateResourceRequest{Name: "", Uuid: resourceId})
	assert.NotNil(t, err)
	resourceCount, _ = s.CountResources(ctx)
	assert.Equal(t, int64(2), resourceCount)

	// unexpected error in update
	_, err = s.UpdateResource(ctx, &auth.UpdateResourceRequest{Name: "error", Uuid: resourceId})
	assert.Equal(t, errCRUD, err)
	count, _ = s.CountResources(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.GetResource(ctx, "none")
	assert.NotNil(t, err)
	resource, err = s.GetResource(ctx, resourceId)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", resource.Name)
	assert.Equal(t, id, resource.Uuid)

	// query
	_resources, _ := s.QueryResources(ctx, "", 0, 0)
	assert.Equal(t, 2, int(_resources.TotalCount))

	// delete
	_, err = s.DeleteResource(ctx, "none")
	assert.NotNil(t, err)
	resource, err = s.DeleteResource(ctx, resourceId)
	assert.Nil(t, err)
	assert.Equal(t, id, resource.Uuid)
	resourceCount, _ = s.CountResources(ctx)
	assert.Equal(t, int64(1), resourceCount)

	// objects
	// initial count
	objectCount, _ := s.CountObjects(ctx)
	assert.Equal(t, int64(0), objectCount)

	// successful creation
	object, err := s.CreateObject(ctx, &auth.CreateObjectRequest{Identifier: "test", Uuid: app.Uuid})
	assert.Nil(t, err)
	assert.NotEmpty(t, object.Uuid)
	objectId := object.Uuid
	assert.Equal(t, "test", object.Identifier)
	assert.NotEmpty(t, object.CreatedAt)
	assert.NotEmpty(t, object.UpdatedAt)
	objectCount, _ = s.CountObjects(ctx)
	assert.Equal(t, int64(1), objectCount)

	// validation error in creation
	_, err = s.CreateObject(ctx, &auth.CreateObjectRequest{Identifier: "", Uuid: app.Uuid})
	assert.NotNil(t, err)
	objectCount, _ = s.CountObjects(ctx)
	assert.Equal(t, int64(1), objectCount)

	// unexpected error in creation
	_, err = s.CreateObject(ctx, &auth.CreateObjectRequest{Identifier: "error", Uuid: app.Uuid})
	assert.Equal(t, errCRUD, err)
	objectCount, _ = s.CountObjects(ctx)
	assert.Equal(t, int64(1), objectCount)

	_, _ = s.CreateObject(ctx, &auth.CreateObjectRequest{Identifier: "test2", Uuid: app.Uuid})

	// update
	object, err = s.UpdateObject(ctx, &auth.UpdateObjectRequest{Identifier: "test updated", Uuid: objectId})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", object.Identifier)
	_, err = s.UpdateObject(ctx, &auth.UpdateObjectRequest{Identifier: "test updated", Uuid: "none"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.UpdateObject(ctx, &auth.UpdateObjectRequest{Identifier: "", Uuid: objectId})
	assert.NotNil(t, err)
	objectCount, _ = s.CountObjects(ctx)
	assert.Equal(t, int64(2), objectCount)

	// unexpected error in update
	_, err = s.UpdateObject(ctx, &auth.UpdateObjectRequest{Identifier: "error", Uuid: objectId})
	assert.Equal(t, errCRUD, err)
	count, _ = s.CountObjects(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.GetObject(ctx, "none")
	assert.NotNil(t, err)
	object, err = s.GetObject(ctx, objectId)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", object.Identifier)
	assert.Equal(t, id, object.Uuid)

	// query
	_objects, _ := s.QueryObjects(ctx, "", 0, 0)
	assert.Equal(t, 2, int(_objects.TotalCount))

	// delete
	_, err = s.DeleteObject(ctx, "none")
	assert.NotNil(t, err)
	object, err = s.DeleteObject(ctx, objectId)
	assert.Nil(t, err)
	assert.Equal(t, id, object.Uuid)
	objectCount, _ = s.CountObjects(ctx)
	assert.Equal(t, int64(1), objectCount)
}
