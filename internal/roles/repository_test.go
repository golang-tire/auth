package roles

import (
	"context"
	"testing"

	"gorm.io/gorm"

	"github.com/golang-tire/auth/internal/db"

	"github.com/golang-tire/auth/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	database := db.NewForTest(t, []interface{}{(*entity.Role)(nil)})
	err := db.ResetTables(t, database, "roles")
	assert.Nil(t, err)
	repo := NewRepository(database)

	ctx := context.Background()
	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	testUuid, err := repo.Create(ctx, entity.Role{
		Title:  "admin",
		Enable: true,
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, int64(1), count2-count)

	// get
	role, err := repo.Get(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "admin", role.Title)
	_, err = repo.Get(ctx, "test0")
	assert.NotNil(t, err)
	assert.EqualError(t, gorm.ErrRecordNotFound, err.Error())

	// get by title
	role, err = repo.GetByTitle(ctx, "admin")
	assert.Nil(t, err)
	assert.Equal(t, "admin", role.Title)
	_, err = repo.GetByTitle(ctx, "test0")
	assert.NotNil(t, err)
	assert.EqualError(t, gorm.ErrRecordNotFound, err.Error())

	// update
	err = repo.Update(ctx, entity.Role{
		UUID:   testUuid,
		Title:  "manager",
		Enable: true,
	})
	assert.Nil(t, err)
	role, _ = repo.Get(ctx, testUuid)
	assert.Equal(t, "manager", role.Title)

	// query
	_, count3, err := repo.Query(ctx, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(count3))

	// delete
	err = repo.Delete(ctx, role)
	assert.Nil(t, err)
	_, err = repo.Get(ctx, testUuid)
	assert.NotNil(t, err)
	assert.EqualError(t, gorm.ErrRecordNotFound, err.Error())
	err = repo.Delete(ctx, role)
	assert.NotNil(t, err)
	assert.EqualError(t, gorm.ErrRecordNotFound, err.Error())
}
