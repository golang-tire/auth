package rules

import (
	"context"
	"testing"

	"github.com/golang-tire/auth/internal/pkg/testutils"

	"github.com/golang-tire/auth/internal/domains"
	"github.com/golang-tire/auth/internal/roles"

	"github.com/golang-tire/auth/internal/pkg/db"

	"github.com/golang-tire/auth/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	testutils.TestUp()
	ctx := context.Background()

	database := db.NewForTest(t, []interface{}{&entity.Domain{}, &entity.Role{}, &entity.Rule{}})
	err := db.ResetTables(t, database, "domain", "roles", "rules")
	assert.Nil(t, err)
	repo := NewRepository(database)

	domainRepo := domains.NewMockRepository()
	roleRepo := roles.NewMockRepository()

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

	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	testUuid, err := repo.Create(ctx, entity.Rule{
		Domain:   domain,
		Role:     role,
		Resource: "rules",
		Object:   "rules",
		Action:   "get",
		Effect:   "allow",
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

	// update
	rule.Resource = "products"
	err = repo.Update(ctx, rule)
	assert.Nil(t, err)
	rule, _ = repo.Get(ctx, testUuid)
	assert.Equal(t, "products", rule.Resource)

	// query
	_, count3, err := repo.Query(ctx, "", 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(count3))

	// delete
	err = repo.Delete(ctx, rule)
	assert.Nil(t, err)
	_, err = repo.Get(ctx, testUuid)
	assert.NotNil(t, err)

	testutils.TestDown()
}
