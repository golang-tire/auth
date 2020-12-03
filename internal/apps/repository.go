package apps

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/golang-tire/auth/internal/pkg/db"

	"github.com/google/uuid"

	"github.com/golang-tire/auth/internal/entity"
)

// Repository encapsulates the logic to access apps from the data source.
type Repository interface {
	// GetApp returns the app with the specified app UUID.
	GetApp(ctx context.Context, uuid string) (entity.App, error)
	// GetAppByName returns the app with the specified app name.
	GetAppByName(ctx context.Context, name string) (entity.App, error)
	// CountApps returns the number of apps.
	CountApps(ctx context.Context) (int64, error)
	// QueryApps returns the list of apps with the given offset and limit.
	QueryApps(ctx context.Context, query string, offset, limit int64) ([]entity.App, int, error)
	// CreateApp saves a new app in the storage.
	CreateApp(ctx context.Context, app entity.App) (string, error)
	// UpdateApp updates the app with given UUID in the storage.
	UpdateApp(ctx context.Context, app entity.App) error
	// DeleteApp removes the app with given UUID from the storage.
	DeleteApp(ctx context.Context, app entity.App) error

	// GetResource returns the Resource with the specified Resource UUID.
	GetResource(ctx context.Context, uuid string) (entity.Resource, error)
	// GetResourceByName returns the app with the specified Resource name.
	GetResourceByName(ctx context.Context, name string) (entity.Resource, error)
	// CountResources returns the number of Resources.
	CountResources(ctx context.Context) (int64, error)
	// QueryResources returns the list of Resources with the given offset and limit.
	QueryResources(ctx context.Context, query string, offset, limit int64) ([]entity.Resource, int, error)
	// CreateResource saves a new Resource in the storage.
	CreateResource(ctx context.Context, app entity.App, resource entity.Resource) (string, error)
	// UpdateResource updates the Resource with given UUID in the storage.
	UpdateResource(ctx context.Context, app entity.App, resource entity.Resource) error
	// DeleteResource removes the Resource with given UUID from the storage.
	DeleteResource(ctx context.Context, resource entity.Resource) error

	// GetObject returns the Object with the specified Object UUID.
	GetObject(ctx context.Context, uuid string) (entity.Object, error)
	// GetObjectByIdentifier returns the app with the specified Object Identifier.
	GetObjectByIdentifier(ctx context.Context, identifier string) (entity.Object, error)
	// CountObjects returns the number of Objects.
	CountObjects(ctx context.Context) (int64, error)
	// QueryObjects returns the list of Objects with the given offset and limit.
	QueryObjects(ctx context.Context, query string, offset, limit int64) ([]entity.Object, int, error)
	// CreateObject saves a new Object in the storage.
	CreateObject(ctx context.Context, app entity.App, object entity.Object) (string, error)
	// UpdateObject updates the Object with given UUID in the storage.
	UpdateObject(ctx context.Context, app entity.App, object entity.Object) error
	// DeleteObject removes the Object with given UUID from the storage.
	DeleteObject(ctx context.Context, object entity.Object) error
}

// repository persists apps in database
type repository struct {
	db *db.DB
}

// NewRepository creates a new app repository
func NewRepository(db *db.DB) Repository {
	return repository{db}
}

// Get reads the app with the specified ID from the database.
func (r repository) GetApp(ctx context.Context, uuid string) (entity.App, error) {
	var app entity.App
	res := r.db.With(ctx).Where("uuid = ?", uuid).First(&app)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		return entity.App{}, fmt.Errorf("app with uuid `%s` not found", uuid)
	}
	return app, res.Error
}

// GetAppByName returns the app with the specified app name.
func (r repository) GetAppByName(ctx context.Context, name string) (entity.App, error) {
	var app entity.App
	res := r.db.With(ctx).Where("name = ?", name).First(&app)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		return entity.App{}, fmt.Errorf("app with name `%s` not found", name)
	}
	return app, res.Error
}

// CreateApp saves a new app record in the database.
// It returns the UUID of the newly inserted app record.
func (r repository) CreateApp(ctx context.Context, app entity.App) (string, error) {
	now := time.Now()
	app.UUID = uuid.New().String()
	app.CreatedAt = now
	app.UpdatedAt = now
	res := r.db.With(ctx).Create(&app)
	return app.UUID, res.Error
}

// UpdateApp saves the changes to an app in the database.
func (r repository) UpdateApp(ctx context.Context, app entity.App) error {
	res := r.db.With(ctx).Save(&app)
	return res.Error
}

// DeleteApp deletes an app with the specified ID from the database.
func (r repository) DeleteApp(ctx context.Context, app entity.App) error {
	res := r.db.With(ctx).Delete(&app)
	return res.Error
}

// CountApps returns the number of the app records in the database.
func (r repository) CountApps(ctx context.Context) (int64, error) {
	var count int64
	res := r.db.With(ctx).Model(&entity.App{}).Count(&count)
	return count, res.Error
}

// QueryApps retrieves the app records with the specified offset and limit from the database.
func (r repository) QueryApps(ctx context.Context, query string, offset, limit int64) ([]entity.App, int, error) {
	var _apps []entity.App
	res := r.db.With(ctx).
		Limit(int(limit)).
		Offset(int(offset)).
		Order("id asc")

	if len(query) >= 1 {
		res = res.Where("name LIKE ?", "%"+query+"%").Find(&_apps)
	} else {
		res = res.Find(&_apps)
	}

	count, err := r.CountApps(ctx)
	if err != nil {
		return nil, 0, err
	}
	return _apps, int(count), res.Error
}

// GetResource reads the resource with the specified ID from the database.
func (r repository) GetResource(ctx context.Context, uuid string) (entity.Resource, error) {
	var resource entity.Resource
	res := r.db.With(ctx).Where("uuid = ?", uuid).Preload("App").First(&resource)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		return entity.Resource{}, fmt.Errorf("resource with uuid `%s` not found", uuid)
	}
	return resource, res.Error
}

// GetResourceByName returns the resource with the specified resource name.
func (r repository) GetResourceByName(ctx context.Context, name string) (entity.Resource, error) {
	var resource entity.Resource
	res := r.db.With(ctx).Where("name = ?", name).Preload("App").First(&resource)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		return entity.Resource{}, fmt.Errorf("resource with name `%s` not found", name)
	}
	return resource, res.Error
}

// CreateResource saves a new resource record in the database.
// It returns the UUID of the newly inserted resource record.
func (r repository) CreateResource(ctx context.Context, app entity.App, resource entity.Resource) (string, error) {
	now := time.Now()
	resource.UUID = uuid.New().String()
	resource.CreatedAt = now
	resource.UpdatedAt = now
	resource.App = app
	resource.AppID = app.ID
	res := r.db.With(ctx).Create(&resource)
	return resource.UUID, res.Error
}

// UpdateResource saves the changes to an resource in the database.
func (r repository) UpdateResource(ctx context.Context, app entity.App, resource entity.Resource) error {
	resource.App = app
	resource.AppID = app.ID
	res := r.db.With(ctx).Save(&resource)
	return res.Error
}

// DeleteResource deletes an resource with the specified ID from the database.
func (r repository) DeleteResource(ctx context.Context, resource entity.Resource) error {
	res := r.db.With(ctx).Delete(&resource)
	return res.Error
}

// CountResources returns the number of the resource records in the database.
func (r repository) CountResources(ctx context.Context) (int64, error) {
	var count int64
	res := r.db.With(ctx).Model(&entity.Resource{}).Count(&count)
	return count, res.Error
}

// QueryResources retrieves the resource records with the specified offset and limit from the database.
func (r repository) QueryResources(ctx context.Context, query string, offset, limit int64) ([]entity.Resource, int, error) {
	var _resources []entity.Resource
	res := r.db.With(ctx).
		Limit(int(limit)).
		Offset(int(offset)).
		Order("id asc").
		Preload("App")

	if len(query) >= 1 {
		res = res.Where("name LIKE ?", "%"+query+"%").Find(&_resources)
	} else {
		res = res.Find(&_resources)
	}

	count, err := r.CountResources(ctx)
	if err != nil {
		return nil, 0, err
	}
	return _resources, int(count), res.Error
}

// GetObject reads the object with the specified ID from the database.
func (r repository) GetObject(ctx context.Context, uuid string) (entity.Object, error) {
	var object entity.Object
	res := r.db.With(ctx).Where("uuid = ?", uuid).Preload("App").First(&object)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		return entity.Object{}, fmt.Errorf("object with uuid `%s` not found", uuid)
	}
	return object, res.Error
}

// GetObjectByName returns the object with the specified object identifier.
func (r repository) GetObjectByIdentifier(ctx context.Context, identifier string) (entity.Object, error) {
	var object entity.Object
	res := r.db.With(ctx).Where("identifier = ?", identifier).Preload("App").First(&object)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		return entity.Object{}, fmt.Errorf("object with identifier `%s` not found", identifier)
	}
	return object, res.Error
}

// CreateObject saves a new object record in the database.
// It returns the UUID of the newly inserted object record.
func (r repository) CreateObject(ctx context.Context, app entity.App, object entity.Object) (string, error) {
	now := time.Now()
	object.UUID = uuid.New().String()
	object.CreatedAt = now
	object.UpdatedAt = now
	object.AppID = app.ID
	object.App = app
	res := r.db.With(ctx).Create(&object)
	return object.UUID, res.Error
}

// UpdateObject saves the changes to an object in the database.
func (r repository) UpdateObject(ctx context.Context, app entity.App, object entity.Object) error {
	object.AppID = app.ID
	object.App = app
	res := r.db.With(ctx).Save(&object)
	return res.Error
}

// DeleteObject deletes an object with the specified ID from the database.
func (r repository) DeleteObject(ctx context.Context, object entity.Object) error {
	res := r.db.With(ctx).Delete(&object)
	return res.Error
}

// CountObjects returns the number of the object records in the database.
func (r repository) CountObjects(ctx context.Context) (int64, error) {
	var count int64
	res := r.db.With(ctx).Model(&entity.Object{}).Count(&count)
	return count, res.Error
}

// QueryObjects retrieves the object records with the specified offset and limit from the database.
func (r repository) QueryObjects(ctx context.Context, query string, offset, limit int64) ([]entity.Object, int, error) {
	var _objects []entity.Object
	res := r.db.With(ctx).
		Limit(int(limit)).
		Offset(int(offset)).
		Order("id asc").
		Preload("App")

	if len(query) >= 1 {
		res = res.Where("identifier LIKE ?", "%"+query+"%").Find(&_objects)
	} else {
		res = res.Find(&_objects)
	}

	count, err := r.CountObjects(ctx)
	if err != nil {
		return nil, 0, err
	}
	return _objects, int(count), res.Error
}
