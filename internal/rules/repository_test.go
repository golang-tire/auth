package rules

import (
	"context"
	"testing"

	"gorm.io/gorm"

	"github.com/golang-tire/auth/internal/db"

	"github.com/golang-tire/auth/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	database := db.NewForTest(t, []interface{}{(*entity.Rule)(nil)})
	err := db.ResetTables(t, database, "rules")
	assert.Nil(t, err)
	repo := NewRepository(database)

	ctx := context.Background()
	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	testUuid, err := repo.Create(ctx, entity.Rule{
		Object: "rules",
		Action: "get",
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, int64(1), count2-count)

	// get
	rule, err := repo.Get(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "rules", rule.Object)
	_, err = repo.Get(ctx, "test0")
	assert.NotNil(t, err)
	assert.EqualError(t, gorm.ErrRecordNotFound, err.Error())

	// update
	err = repo.Update(ctx, entity.Rule{
		UUID:   testUuid,
		Object: "products",
		Action: "get",
	})
	assert.Nil(t, err)
	rule, _ = repo.Get(ctx, testUuid)
	assert.Equal(t, "products", rule.Object)

	// query
	_, count3, err := repo.Query(ctx, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(count3))

	// delete
	err = repo.Delete(ctx, rule)
	assert.Nil(t, err)
	_, err = repo.Get(ctx, testUuid)
	assert.NotNil(t, err)
	assert.EqualError(t, gorm.ErrRecordNotFound, err.Error())
	err = repo.Delete(ctx, rule)
	assert.NotNil(t, err)
	assert.EqualError(t, gorm.ErrRecordNotFound, err.Error())
}
