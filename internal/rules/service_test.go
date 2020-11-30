package rules

import (
	"context"
	"testing"

	"github.com/golang-tire/auth/internal/pkg/testutils"

	"github.com/golang-tire/auth/internal/entity"

	"github.com/golang-tire/auth/internal/domains"
	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang-tire/auth/internal/roles"
	"github.com/stretchr/testify/assert"
)

func TestCreateRuleRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.CreateRuleRequest
		wantError bool
	}{
		{"success", auth.CreateRuleRequest{
			Role:   "test",
			Domain: "test",
			Object: "test",
			Action: "test",
			Effect: auth.Effect_ALLOW,
		}, false},
		{"required", auth.CreateRuleRequest{
			Role:   "",
			Domain: "test",
			Object: "test",
			Action: "test",
			Effect: auth.Effect_ALLOW,
		}, true},
		{"too long", auth.CreateRuleRequest{
			Role:   "test",
			Domain: "test",
			Object: "test",
			Action: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			Effect: auth.Effect_ALLOW,
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateRuleRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.UpdateRuleRequest
		wantError bool
	}{
		{"success", auth.UpdateRuleRequest{
			Role:   "test",
			Domain: "test",
			Object: "test",
			Action: "test",
			Effect: auth.Effect_ALLOW,
		}, false},
		{"required", auth.UpdateRuleRequest{
			Role:   "",
			Domain: "test",
			Object: "test",
			Action: "test",
			Effect: auth.Effect_ALLOW,
		}, true},
		{"too long", auth.UpdateRuleRequest{
			Role:   "test",
			Domain: "test",
			Object: "test",
			Action: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			Effect: auth.Effect_ALLOW,
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {

	testutils.TestUp()

	domainRepo := domains.NewMockRepository()
	roleRepo := roles.NewMockRepository()

	s := NewService(&MockRepository{}, domainRepo, roleRepo)
	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, int64(0), count)

	dId, err := domainRepo.Create(ctx, entity.Domain{
		Name:   "foo.bar",
		Enable: true,
	})
	assert.Nil(t, err)

	rId, err := roleRepo.Create(ctx, entity.Role{
		Title:  "admin",
		Enable: true,
	})
	assert.Nil(t, err)

	domain, err := domainRepo.Get(ctx, dId)
	assert.Nil(t, err)
	role, err := roleRepo.Get(ctx, rId)
	assert.Nil(t, err)

	// successful creation
	rule, err := s.Create(ctx, &auth.CreateRuleRequest{
		Role:     role.Title,
		Resource: "products",
		Domain:   domain.Name,
		Object:   "test",
		Action:   "test",
		Effect:   auth.Effect_ALLOW,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, rule.Uuid)
	id := rule.Uuid
	assert.Equal(t, "admin", rule.Role)
	assert.NotEmpty(t, rule.CreatedAt)
	assert.NotEmpty(t, rule.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.Create(ctx, &auth.CreateRuleRequest{
		Role:     "",
		Resource: "products",
		Domain:   "test",
		Object:   "test",
		Action:   "test",
		Effect:   auth.Effect_ALLOW,
	})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.Create(ctx, &auth.CreateRuleRequest{
		Role:     role.Title,
		Resource: "error",
		Domain:   domain.Name,
		Object:   "test",
		Action:   "test",
		Effect:   auth.Effect_ALLOW,
	})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.Create(ctx, &auth.CreateRuleRequest{
		Role:     role.Title,
		Resource: "cars",
		Domain:   domain.Name,
		Object:   "test",
		Action:   "test",
		Effect:   auth.Effect_ALLOW,
	})

	// update
	rule, err = s.Update(ctx, &auth.UpdateRuleRequest{
		Role:     role.Title,
		Resource: "products",
		Domain:   domain.Name,
		Object:   "test-updated",
		Action:   "test",
		Effect:   auth.Effect_ALLOW,
		Uuid:     id,
	})
	assert.Nil(t, err)
	assert.Equal(t, "test-updated", rule.Object)

	_, err = s.Update(ctx, &auth.UpdateRuleRequest{
		Role:     role.Title,
		Resource: "products",
		Domain:   domain.Name,
		Object:   "test",
		Action:   "test",
		Uuid:     "none",
		Effect:   auth.Effect_ALLOW,
	})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, &auth.UpdateRuleRequest{
		Role:     "",
		Resource: "products",
		Domain:   domain.Name,
		Object:   "test",
		Action:   "test",
		Uuid:     id,
		Effect:   auth.Effect_ALLOW,
	})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// unexpected error in update
	_, err = s.Update(ctx, &auth.UpdateRuleRequest{
		Role:     role.Title,
		Resource: "error",
		Domain:   domain.Name,
		Object:   "test",
		Action:   "test",
		Effect:   auth.Effect_ALLOW,
		Uuid:     id,
	})
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	rule, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test-updated", rule.Object)
	assert.Equal(t, id, rule.Uuid)

	// query
	_rules, _ := s.Query(ctx, "", 0, 0)
	assert.Equal(t, 2, int(_rules.TotalCount))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	rule, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, rule.Uuid)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	testutils.TestDown()
}
