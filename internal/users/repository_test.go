package users

import (
	"context"
	"testing"

	"github.com/golang-tire/auth/internal/pkg/testutils"

	"github.com/golang-tire/auth/internal/pkg/db"

	"github.com/golang-tire/auth/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {

	testutils.TestUp()
	database := db.NewForTest(t, []interface{}{&entity.Role{}, &entity.Domain{}, &entity.User{}, &entity.UserRole{}})
	err := db.ResetTables(t, database, "roles", "domains", "users", "user_roles")
	assert.Nil(t, err)
	repo := NewRepository(database)

	ctx := context.Background()
	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	testUuid, err := repo.Create(ctx, entity.User{
		Firstname: "foo",
		Lastname:  "bar",
		Username:  "user1",
		Password:  "pass",
		Gender:    "x",
		AvatarURL: "https://foo.bar/foo.jpg",
		Email:     "foo1@bar.com",
		Enable:    true,
		RawData:   "",
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, int64(1), count2-count)

	// get
	user, err := repo.Get(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "foo", user.Firstname)
	_, err = repo.Get(ctx, "test0")
	assert.NotNil(t, err)

	// update
	user.Firstname = "bar"
	err = repo.Update(ctx, user)
	assert.Nil(t, err)
	user, _ = repo.Get(ctx, testUuid)
	assert.Equal(t, "bar", user.Firstname)

	// query
	_, count3, err := repo.Query(ctx, "", 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(count3))

	// delete
	err = repo.Delete(ctx, user)
	assert.Nil(t, err)
	_, err = repo.Get(ctx, testUuid)
	assert.NotNil(t, err)
	testutils.TestDown()
}
