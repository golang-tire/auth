package rules

import (
	"context"
	"testing"

	auth "github.com/golang-tire/auth/internal/proto/v1"
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
		}, false},
		{"required", auth.CreateRuleRequest{
			Role:   "",
			Domain: "test",
			Object: "test",
			Action: "test",
		}, true},
		{"too long", auth.CreateRuleRequest{
			Role:   "test",
			Domain: "test",
			Object: "test",
			Action: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
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
		}, false},
		{"required", auth.UpdateRuleRequest{
			Role:   "",
			Domain: "test",
			Object: "test",
			Action: "test",
		}, true},
		{"too long", auth.UpdateRuleRequest{
			Role:   "test",
			Domain: "test",
			Object: "test",
			Action: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
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
	rule, err := s.Create(ctx, &auth.CreateRuleRequest{
		Role:   "test",
		Domain: "test",
		Object: "test",
		Action: "test",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, rule.Uuid)
	id := rule.Uuid
	assert.Equal(t, "test", rule.Role)
	assert.NotEmpty(t, rule.CreatedAt)
	assert.NotEmpty(t, rule.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.Create(ctx, &auth.CreateRuleRequest{
		Role:   "",
		Domain: "test",
		Object: "test",
		Action: "test",
	})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.Create(ctx, &auth.CreateRuleRequest{
		Role:   "error",
		Domain: "test",
		Object: "test",
		Action: "test",
	})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.Create(ctx, &auth.CreateRuleRequest{
		Role:   "test2",
		Domain: "test",
		Object: "test",
		Action: "test",
	})

	// update
	rule, err = s.Update(ctx, &auth.UpdateRuleRequest{
		Role:   "test updated",
		Domain: "test",
		Object: "test",
		Action: "test",
		Uuid:   id,
	})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", rule.Role)
	_, err = s.Update(ctx, &auth.UpdateRuleRequest{
		Role:   "test updated",
		Domain: "test",
		Object: "test",
		Action: "test",
		Uuid:   "none",
	})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, &auth.UpdateRuleRequest{
		Role:   "",
		Domain: "test",
		Object: "test",
		Action: "test",
		Uuid:   id,
	})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// unexpected error in update
	_, err = s.Update(ctx, &auth.UpdateRuleRequest{
		Role:   "error",
		Domain: "test",
		Object: "test",
		Action: "test",
		Uuid:   id,
	})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	rule, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", rule.Role)
	assert.Equal(t, id, rule.Uuid)

	// query
	_rules, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, int(_rules.TotalCount))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	rule, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, rule.Uuid)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)
}
