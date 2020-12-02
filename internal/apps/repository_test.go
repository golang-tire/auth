package apps

import (
	"context"
	"testing"

	"github.com/golang-tire/auth/internal/pkg/db"

	"github.com/golang-tire/auth/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	database := db.NewForTest(t, []interface{}{&entity.App{}, &entity.Resource{}, &entity.Object{}})
	err := db.ResetTables(t, database, "apps", "resources", "objects")
	assert.Nil(t, err)
	repo := NewRepository(database)

	ctx := context.Background()

	// Apps
	// initial count
	appCount, err := repo.CountApps(ctx)
	assert.Nil(t, err)

	// create
	testUuid, err := repo.CreateApp(ctx, entity.App{
		Name:   "foo.bar",
		Enable: true,
	})
	assert.Nil(t, err)
	count2, _ := repo.CountApps(ctx)
	assert.Equal(t, int64(1), count2-appCount)

	// get
	app, err := repo.GetApp(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "foo.bar", app.Name)
	_, err = repo.GetApp(ctx, "test0")
	assert.NotNil(t, err)

	// update
	app.Name = "bar.foo"
	err = repo.UpdateApp(ctx, app)
	assert.Nil(t, err)
	app, _ = repo.GetApp(ctx, testUuid)
	assert.Equal(t, "bar.foo", app.Name)

	// query
	_, count3, err := repo.QueryApps(ctx, "", 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(count3))

	// delete
	err = repo.DeleteApp(ctx, app)
	assert.Nil(t, err)
	_, err = repo.GetApp(ctx, testUuid)
	assert.NotNil(t, err)

	// Resources
	appUuid, err := repo.CreateApp(ctx, entity.App{
		Name:   "foo.bar",
		Enable: true,
	})
	assert.Nil(t, err)
	testApp, err := repo.GetApp(ctx, appUuid)
	assert.Nil(t, err)

	// initial count
	resourceCount, err := repo.CountResources(ctx)
	assert.Nil(t, err)

	// create
	testUuid, err = repo.CreateResource(ctx, testApp, entity.Resource{
		Name: "foo.bar",
	})
	assert.Nil(t, err)
	count2, _ = repo.CountResources(ctx)
	assert.Equal(t, int64(1), count2-resourceCount)

	// get
	resource, err := repo.GetResource(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "foo.bar", resource.Name)
	_, err = repo.GetResource(ctx, "test0")
	assert.NotNil(t, err)

	// update
	resource.Name = "bar.foo"
	err = repo.UpdateResource(ctx, testApp, resource)
	assert.Nil(t, err)
	resource, _ = repo.GetResource(ctx, testUuid)
	assert.Equal(t, "bar.foo", resource.Name)

	// query
	_, count3, err = repo.QueryResources(ctx, "", 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(count3))

	// delete
	err = repo.DeleteResource(ctx, resource)
	assert.Nil(t, err)
	_, err = repo.GetResource(ctx, testUuid)
	assert.NotNil(t, err)

	// initial count
	objectCount, err := repo.CountObjects(ctx)
	assert.Nil(t, err)

	// create
	testUuid, err = repo.CreateObject(ctx, testApp, entity.Object{
		Identifier: "foo.bar",
	})
	assert.Nil(t, err)
	count2, _ = repo.CountObjects(ctx)
	assert.Equal(t, int64(1), count2-objectCount)

	// get
	object, err := repo.GetObject(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "foo.bar", object.Identifier)
	_, err = repo.GetObject(ctx, "test0")
	assert.NotNil(t, err)

	// update
	object.Identifier = "bar.foo"
	err = repo.UpdateObject(ctx, testApp, object)
	assert.Nil(t, err)
	object, _ = repo.GetObject(ctx, testUuid)
	assert.Equal(t, "bar.foo", object.Identifier)

	// query
	_, count3, err = repo.QueryObjects(ctx, "", 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(count3))

	// delete
	err = repo.DeleteObject(ctx, object)
	assert.Nil(t, err)
	_, err = repo.GetObject(ctx, testUuid)
	assert.NotNil(t, err)
}
