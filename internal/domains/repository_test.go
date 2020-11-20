package domains

import (
	"context"
	"testing"

	"github.com/golang-tire/auth/internal/pkg/db"

	"github.com/golang-tire/auth/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	database := db.NewForTest(t, []interface{}{&entity.Domain{}})
	err := db.ResetTables(t, database, "domains")
	assert.Nil(t, err)
	repo := NewRepository(database)

	ctx := context.Background()
	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	testUuid, err := repo.Create(ctx, entity.Domain{
		Name:   "foo.bar",
		Enable: true,
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, int64(1), count2-count)

	// get
	domain, err := repo.Get(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "foo.bar", domain.Name)
	_, err = repo.Get(ctx, "test0")
	assert.NotNil(t, err)

	// update
	domain.Name = "bar.foo"
	err = repo.Update(ctx, domain)
	assert.Nil(t, err)
	domain, _ = repo.Get(ctx, testUuid)
	assert.Equal(t, "bar.foo", domain.Name)

	// query
	_, count3, err := repo.Query(ctx, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(count3))

	// delete
	err = repo.Delete(ctx, domain)
	assert.Nil(t, err)
	_, err = repo.Get(ctx, testUuid)
	assert.NotNil(t, err)
}
