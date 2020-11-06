package rules

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/go-pg/pg/v10"
	"github.com/golang-tire/auth/internal/entity"
	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/stretchr/testify/assert"
)

var errCRUD = errors.New("error crud")

func TestCreateRuleRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.CreateRuleRequest
		wantError bool
	}{
		{"success", auth.CreateRuleRequest{
			Subject: "test",
			Domain:  "test",
			Object:  "test",
			Action:  "test",
		}, false},
		{"required", auth.CreateRuleRequest{
			Subject: "",
			Domain:  "test",
			Object:  "test",
			Action:  "test",
		}, true},
		{"too long", auth.CreateRuleRequest{
			Subject: "test",
			Domain:  "test",
			Object:  "test",
			Action:  "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
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
			Subject: "test",
			Domain:  "test",
			Object:  "test",
			Action:  "test",
		}, false},
		{"required", auth.UpdateRuleRequest{
			Subject: "",
			Domain:  "test",
			Object:  "test",
			Action:  "test",
		}, true},
		{"too long", auth.UpdateRuleRequest{
			Subject: "test",
			Domain:  "test",
			Object:  "test",
			Action:  "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
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
		Subject: "test",
		Domain:  "test",
		Object:  "test",
		Action:  "test",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, rule.Uuid)
	id := rule.Uuid
	assert.Equal(t, "test", rule.Subject)
	assert.NotEmpty(t, rule.CreatedAt)
	assert.NotEmpty(t, rule.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.Create(ctx, &auth.CreateRuleRequest{
		Subject: "",
		Domain:  "test",
		Object:  "test",
		Action:  "test",
	})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.Create(ctx, &auth.CreateRuleRequest{
		Subject: "error",
		Domain:  "test",
		Object:  "test",
		Action:  "test",
	})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.Create(ctx, &auth.CreateRuleRequest{
		Subject: "test2",
		Domain:  "test",
		Object:  "test",
		Action:  "test",
	})

	// update
	rule, err = s.Update(ctx, &auth.UpdateRuleRequest{
		Subject: "test updated",
		Domain:  "test",
		Object:  "test",
		Action:  "test",
		Uuid:    id,
	})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", rule.Subject)
	_, err = s.Update(ctx, &auth.UpdateRuleRequest{
		Subject: "test updated",
		Domain:  "test",
		Object:  "test",
		Action:  "test",
		Uuid:    "none",
	})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, &auth.UpdateRuleRequest{
		Subject: "",
		Domain:  "test",
		Object:  "test",
		Action:  "test",
		Uuid:    id,
	})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// unexpected error in update
	_, err = s.Update(ctx, &auth.UpdateRuleRequest{
		Subject: "error",
		Domain:  "test",
		Object:  "test",
		Action:  "test",
		Uuid:    id,
	})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	rule, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", rule.Subject)
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

type mockRepository struct {
	items []entity.Rule
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Rule, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return entity.Rule{}, pg.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int64) ([]entity.Rule, int, error) {
	return m.items, len(m.items), nil
}

func (m *mockRepository) Create(ctx context.Context, rule entity.Rule) (string, error) {
	Uuid := uuid.New().String()
	if rule.Subject == "error" {
		return Uuid, errCRUD
	}
	m.items = append(m.items, rule)
	return Uuid, nil
}

func (m *mockRepository) Update(ctx context.Context, rule entity.Rule) error {
	if rule.Subject == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.UUID == rule.UUID {
			m.items[i] = rule
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	for i, item := range m.items {
		if item.UUID == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}
