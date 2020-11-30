package users

import (
	"context"
	"testing"

	"github.com/golang-tire/auth/internal/pkg/testutils"

	"github.com/golang-tire/auth/internal/roles"

	"github.com/golang-tire/auth/internal/domains"

	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.CreateUserRequest
		wantError bool
	}{
		{"success", auth.CreateUserRequest{
			Firstname: "foo",
			Lastname:  "bar",
			Username:  "user1user1",
			Password:  "pass1234pass",
			Gender:    "x",
			AvatarUrl: "https://foo.bar/foo.jpg",
			Email:     "foo@bar.com",
			Enable:    true,
			RawData:   "",
		}, false},
		{"required", auth.CreateUserRequest{
			Firstname: "",
			Lastname:  "bar",
			Username:  "",
			Password:  "pass",
			Gender:    "x",
			AvatarUrl: "https://foo.bar/foo.jpg",
			Email:     "foo@bar.com",
			Enable:    true,
			RawData:   "",
		}, true},
		{"too long", auth.CreateUserRequest{
			Firstname: "foo",
			Lastname:  "bar",
			Password:  "pass",
			Gender:    "x",
			AvatarUrl: "https://foo.bar/foo.jpg",
			Email:     "foo@bar.com",
			Enable:    true,
			RawData:   "",
			Username:  "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     auth.UpdateUserRequest
		wantError bool
	}{
		{"success", auth.UpdateUserRequest{
			Firstname: "foo",
			Lastname:  "bar",
			Username:  "foobarfoo",
			Password:  "pass1234qwer",
			Gender:    "x",
			AvatarUrl: "https://foo.bar/foo.jpg",
			Email:     "foo@bar.com",
			Enable:    true,
			RawData:   "",
		}, false},
		{"required", auth.UpdateUserRequest{
			Firstname: "",
			Lastname:  "bar",
			Username:  "",
			Password:  "pass",
			Gender:    "x",
			AvatarUrl: "https://foo.bar/foo.jpg",
			Email:     "foo@bar.com",
			Enable:    true,
			RawData:   "",
		}, true},
		{"too long", auth.UpdateUserRequest{
			Firstname: "foo",
			Lastname:  "bar",
			Password:  "pass",
			Gender:    "x",
			AvatarUrl: "https://foo.bar/foo.jpg",
			Email:     "foo@bar.com",
			Enable:    true,
			RawData:   "",
			Username:  "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
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
	s := NewService(&mockRepository{}, domains.NewMockRepository(), roles.NewMockRepository())
	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, int64(0), count)

	// successful creation
	user, err := s.Create(ctx, &auth.CreateUserRequest{
		Firstname: "foo",
		Lastname:  "bar",
		Username:  "foobarfoo",
		Password:  "pass1234qwer",
		Gender:    "x",
		AvatarUrl: "https://foo.bar/foo.jpg",
		Email:     "foo@bar.com",
		Enable:    true,
		RawData:   "",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, user.Uuid)
	id := user.Uuid
	assert.Equal(t, "foobarfoo", user.Username)
	assert.NotEmpty(t, user.CreatedAt)
	assert.NotEmpty(t, user.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.Create(ctx, &auth.CreateUserRequest{
		Firstname: "foo",
		Lastname:  "bar",
		Username:  "",
		Password:  "pass",
		Gender:    "x",
		AvatarUrl: "https://foo.bar/foo.jpg",
		Email:     "foo@bar.com",
		Enable:    true,
		RawData:   "",
	})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.Create(ctx, &auth.CreateUserRequest{
		Firstname: "foo",
		Lastname:  "barbar",
		Username:  "errorerror",
		Password:  "pass1234qwer",
		Gender:    "x",
		AvatarUrl: "https://foo.bar/foo.jpg",
		Email:     "foo@bar.com",
		Enable:    true,
		RawData:   "",
	})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.Create(ctx, &auth.CreateUserRequest{
		Firstname: "foo",
		Lastname:  "bar",
		Username:  "foobarfoo2",
		Password:  "pass1234qwer",
		Gender:    "x",
		AvatarUrl: "https://foo.bar/foo.jpg",
		Email:     "foo@bar.com",
		Enable:    true,
		RawData:   "",
	})

	// update
	user, err = s.Update(ctx, &auth.UpdateUserRequest{
		Firstname: "foo",
		Lastname:  "bar",
		Username:  "user_updated",
		Password:  "pass1234qwer",
		Gender:    "x",
		AvatarUrl: "https://foo.bar/foo.jpg",
		Email:     "foo@bar.com",
		Enable:    true,
		RawData:   "",
		Uuid:      id,
	})
	assert.Nil(t, err)
	assert.Equal(t, "user_updated", user.Username)
	_, err = s.Update(ctx, &auth.UpdateUserRequest{
		Firstname: "foo",
		Lastname:  "bar",
		Username:  "user_updated",
		Password:  "pass1234qwer",
		Gender:    "x",
		AvatarUrl: "https://foo.bar/foo.jpg",
		Email:     "foo@bar.com",
		Enable:    true,
		RawData:   "",
		Uuid:      "none",
	})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, &auth.UpdateUserRequest{
		Firstname: "foo",
		Lastname:  "bar",
		Username:  "",
		Password:  "pass1234qwer",
		Gender:    "x",
		AvatarUrl: "https://foo.bar/foo.jpg",
		Email:     "foo@bar.com",
		Enable:    true,
		RawData:   "",
		Uuid:      id,
	})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// unexpected error in update
	_, err = s.Update(ctx, &auth.UpdateUserRequest{
		Firstname: "foo",
		Lastname:  "bar",
		Username:  "errorerror",
		Password:  "pass1234qwer",
		Gender:    "x",
		AvatarUrl: "https://foo.bar/foo.jpg",
		Email:     "bar@foo.com",
		Enable:    true,
		RawData:   "",
		Uuid:      id,
	})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	user, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "user_updated", user.Username)
	assert.Equal(t, id, user.Uuid)

	// query
	_users, _ := s.Query(ctx, "", 0, 0)
	assert.Equal(t, 2, int(_users.TotalCount))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	user, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, user.Uuid)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	testutils.TestDown()
}
