package rules

import (
	"context"
	"testing"
	"time"

	"github.com/golang-tire/auth/internal/db"

	"github.com/go-pg/pg"
	"github.com/golang-tire/auth/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	database := db.NewForTest(t, []interface{}{(*entity.Rule)(nil)})
	db.ResetTables(t, database, "rules")
	repo := NewRepository(database)

	ctx := context.Background()
	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	now := time.Now()
	testUuid, err := repo.Create(ctx, entity.Rule{
		Subject: "admin",
		Domain:  "domain1",
		Object:  "rules",
		Action:  "get",
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, int64(1), count2-count)

	// get
	rule, err := repo.Get(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "admin", rule.Subject)
	_, err = repo.Get(ctx, "test0")
	assert.NotNil(t, err)
	assert.EqualError(t, pg.ErrNoRows, err.Error())

	// update
	err = repo.Update(ctx, entity.Rule{
		ID:        rule.ID,
		UUID:      testUuid,
		Subject:   "manager",
		Domain:    "domain1",
		Object:    "rules",
		Action:    "get",
		CreatedAt: now,
		UpdatedAt: now,
	})
	assert.Nil(t, err)
	rule, _ = repo.Get(ctx, testUuid)
	assert.Equal(t, "manager", rule.Subject)

	// query
	_, count3, err := repo.Query(ctx, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(count3))

	// delete
	err = repo.Delete(ctx, testUuid)
	assert.Nil(t, err)
	_, err = repo.Get(ctx, testUuid)
	assert.NotNil(t, err)
	assert.EqualError(t, pg.ErrNoRows, err.Error())
	err = repo.Delete(ctx, testUuid)
	assert.NotNil(t, err)
	assert.EqualError(t, pg.ErrNoRows, err.Error())
}
