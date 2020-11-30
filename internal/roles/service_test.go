package roles

import (
	"context"
	"testing"

	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoleRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.CreateRoleRequest
		wantError bool
	}{
		{"success", auth.CreateRoleRequest{Title: "test"}, false},
		{"required", auth.CreateRoleRequest{Title: ""}, true},
		{"too long", auth.CreateRoleRequest{Title: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateRoleRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.UpdateRoleRequest
		wantError bool
	}{
		{"success", auth.UpdateRoleRequest{Title: "test"}, false},
		{"required", auth.UpdateRoleRequest{Title: ""}, true},
		{"too long", auth.UpdateRoleRequest{Title: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	s := NewService(&mockRepository{})
	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, int64(0), count)

	// successful creation
	role, err := s.Create(ctx, &auth.CreateRoleRequest{Title: "test"})
	assert.Nil(t, err)
	assert.NotEmpty(t, role.Uuid)
	id := role.Uuid
	assert.Equal(t, "test", role.Title)
	assert.NotEmpty(t, role.CreatedAt)
	assert.NotEmpty(t, role.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.Create(ctx, &auth.CreateRoleRequest{Title: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.Create(ctx, &auth.CreateRoleRequest{Title: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.Create(ctx, &auth.CreateRoleRequest{Title: "test2"})

	// update
	role, err = s.Update(ctx, &auth.UpdateRoleRequest{Title: "test updated", Uuid: id})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", role.Title)
	_, err = s.Update(ctx, &auth.UpdateRoleRequest{Title: "test updated", Uuid: "none"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, &auth.UpdateRoleRequest{Title: "", Uuid: id})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// unexpected error in update
	_, err = s.Update(ctx, &auth.UpdateRoleRequest{Title: "error", Uuid: id})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	role, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", role.Title)
	assert.Equal(t, id, role.Uuid)

	// query
	_roles, _ := s.Query(ctx, "", 0, 0)
	assert.Equal(t, 2, int(_roles.TotalCount))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	role, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, role.Uuid)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)
}
