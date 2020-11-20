package domains

import (
	"context"
	"testing"

	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/stretchr/testify/assert"
)

func TestCreateDomainRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.CreateDomainRequest
		wantError bool
	}{
		{"success", auth.CreateDomainRequest{Name: "test"}, false},
		{"required", auth.CreateDomainRequest{Name: ""}, true},
		{"too long", auth.CreateDomainRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateDomainRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.UpdateDomainRequest
		wantError bool
	}{
		{"success", auth.UpdateDomainRequest{Name: "test"}, false},
		{"required", auth.UpdateDomainRequest{Name: ""}, true},
		{"too long", auth.UpdateDomainRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
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
	domain, err := s.Create(ctx, &auth.CreateDomainRequest{Name: "test", Enable: true})
	assert.Nil(t, err)
	assert.NotEmpty(t, domain.Uuid)
	id := domain.Uuid
	assert.Equal(t, "test", domain.Name)
	assert.NotEmpty(t, domain.CreatedAt)
	assert.NotEmpty(t, domain.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.Create(ctx, &auth.CreateDomainRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.Create(ctx, &auth.CreateDomainRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.Create(ctx, &auth.CreateDomainRequest{Name: "test2"})

	// update
	domain, err = s.Update(ctx, &auth.UpdateDomainRequest{Name: "test updated", Uuid: id})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", domain.Name)
	_, err = s.Update(ctx, &auth.UpdateDomainRequest{Name: "test updated", Uuid: "none"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, &auth.UpdateDomainRequest{Name: "", Uuid: id})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// unexpected error in update
	_, err = s.Update(ctx, &auth.UpdateDomainRequest{Name: "error", Uuid: id})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	domain, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", domain.Name)
	assert.Equal(t, id, domain.Uuid)

	// query
	_domains, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, int(_domains.TotalCount))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	domain, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, domain.Uuid)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)
}
