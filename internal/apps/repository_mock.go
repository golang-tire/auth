package apps

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/golang-tire/auth/internal/entity"
	"github.com/google/uuid"
)

var errCRUD = errors.New("error crud")

func NewMockRepository() *mockRepository {
	return &mockRepository{}
}

type mockRepository struct {
	apps      []entity.App
	resources []entity.Resource
	objects   []entity.Object
}

func (m mockRepository) GetAppByName(ctx context.Context, name string) (entity.App, error) {
	for _, item := range m.apps {
		if item.Name == name {
			return item, nil
		}
	}
	return entity.App{}, gorm.ErrRecordNotFound
}

func (m mockRepository) GetApp(ctx context.Context, uuid string) (entity.App, error) {
	for _, item := range m.apps {
		if item.UUID == uuid {
			return item, nil
		}
	}
	return entity.App{}, gorm.ErrRecordNotFound
}

func (m mockRepository) CountApps(ctx context.Context) (int64, error) {
	return int64(len(m.apps)), nil
}

func (m mockRepository) QueryApps(ctx context.Context, query string, offset, limit int64) ([]entity.App, int, error) {
	return m.apps, len(m.apps), nil
}

func (m *mockRepository) CreateApp(ctx context.Context, app entity.App) (string, error) {
	Uuid := uuid.New().String()
	app.UUID = Uuid
	if app.Name == "error" {
		return Uuid, errCRUD
	}
	m.apps = append(m.apps, app)
	return Uuid, nil
}

func (m *mockRepository) UpdateApp(ctx context.Context, app entity.App) error {
	if app.Name == "error" {
		return errCRUD
	}
	for i, item := range m.apps {
		if item.UUID == app.UUID {
			m.apps[i] = app
			break
		}
	}
	return nil
}

func (m *mockRepository) DeleteApp(ctx context.Context, app entity.App) error {
	for i, item := range m.apps {
		if item.UUID == app.UUID {
			m.apps[i] = m.apps[len(m.apps)-1]
			m.apps = m.apps[:len(m.apps)-1]
			break
		}
	}
	return nil
}

func (m mockRepository) GetResourceByName(ctx context.Context, name string) (entity.Resource, error) {
	for _, item := range m.resources {
		if item.Name == name {
			return item, nil
		}
	}
	return entity.Resource{}, gorm.ErrRecordNotFound
}

func (m mockRepository) GetResource(ctx context.Context, uuid string) (entity.Resource, error) {
	for _, item := range m.resources {
		if item.UUID == uuid {
			return item, nil
		}
	}
	return entity.Resource{}, gorm.ErrRecordNotFound
}

func (m mockRepository) CountResources(ctx context.Context) (int64, error) {
	return int64(len(m.resources)), nil
}

func (m mockRepository) QueryResources(ctx context.Context, query string, offset, limit int64) ([]entity.Resource, int, error) {
	return m.resources, len(m.resources), nil
}

func (m *mockRepository) CreateResource(ctx context.Context, app entity.App, resource entity.Resource) (string, error) {
	Uuid := uuid.New().String()
	resource.App = app
	resource.AppID = app.ID
	resource.UUID = Uuid
	if resource.Name == "error" {
		return Uuid, errCRUD
	}
	m.resources = append(m.resources, resource)
	return Uuid, nil
}

func (m *mockRepository) UpdateResource(ctx context.Context, app entity.App, resource entity.Resource) error {
	if resource.Name == "error" {
		return errCRUD
	}
	resource.App = app
	resource.AppID = app.ID
	for i, item := range m.resources {
		if item.UUID == resource.UUID {
			m.resources[i] = resource

			break
		}
	}
	return nil
}

func (m *mockRepository) DeleteResource(ctx context.Context, resource entity.Resource) error {
	for i, item := range m.resources {
		if item.UUID == resource.UUID {
			m.resources[i] = m.resources[len(m.resources)-1]
			m.resources = m.resources[:len(m.resources)-1]
			break
		}
	}
	return nil
}

func (m mockRepository) GetObjectByIdentifier(ctx context.Context, identifier string) (entity.Object, error) {
	for _, item := range m.objects {
		if item.Identifier == identifier {
			return item, nil
		}
	}
	return entity.Object{}, gorm.ErrRecordNotFound
}

func (m mockRepository) GetObject(ctx context.Context, uuid string) (entity.Object, error) {
	for _, item := range m.objects {
		if item.UUID == uuid {
			return item, nil
		}
	}
	return entity.Object{}, gorm.ErrRecordNotFound
}

func (m mockRepository) CountObjects(ctx context.Context) (int64, error) {
	return int64(len(m.objects)), nil
}

func (m mockRepository) QueryObjects(ctx context.Context, query string, offset, limit int64) ([]entity.Object, int, error) {
	return m.objects, len(m.objects), nil
}

func (m *mockRepository) CreateObject(ctx context.Context, app entity.App, object entity.Object) (string, error) {
	Uuid := uuid.New().String()
	object.App = app
	object.AppID = app.ID
	object.UUID = Uuid
	if object.Identifier == "error" {
		return Uuid, errCRUD
	}
	m.objects = append(m.objects, object)
	return Uuid, nil
}

func (m *mockRepository) UpdateObject(ctx context.Context, app entity.App, object entity.Object) error {
	if object.Identifier == "error" {
		return errCRUD
	}
	object.App = app
	object.AppID = app.ID
	for i, item := range m.objects {
		if item.UUID == object.UUID {
			m.objects[i] = object

			break
		}
	}
	return nil
}

func (m *mockRepository) DeleteObject(ctx context.Context, object entity.Object) error {
	for i, item := range m.objects {
		if item.UUID == object.UUID {
			m.objects[i] = m.objects[len(m.objects)-1]
			m.objects = m.objects[:len(m.objects)-1]
			break
		}
	}
	return nil
}
