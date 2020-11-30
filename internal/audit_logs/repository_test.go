package audit_logs

import (
	"context"
	"testing"

	"github.com/golang-tire/auth/internal/users"

	"github.com/golang-tire/auth/internal/entity"
	"github.com/golang-tire/auth/internal/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	database := db.NewForTest(t, []interface{}{&entity.User{}, &entity.UserRole{}, &entity.AuditLog{}})
	err := db.ResetTables(t, database, "audit_logs")
	assert.Nil(t, err)

	userRepo := users.NewMockRepository()
	repo := NewRepository(database)

	ctx := context.Background()

	userUuid, err := userRepo.Create(ctx, entity.User{
		Firstname: "test",
		Lastname:  "test",
		Username:  "test-user-name",
		Password:  "testpass",
		Gender:    "other",
		AvatarURL: "",
		Email:     "email@example.com",
		Enable:    true,
		RawData:   "",
	})
	assert.Nil(t, err)
	user, err := userRepo.Get(ctx, userUuid)
	assert.Nil(t, err)

	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	testUuid, err := repo.Create(ctx, entity.AuditLog{
		User:     user,
		UserID:   user.ID,
		Object:   "foo.bar",
		Action:   "GET",
		OldValue: "ABC",
		NewValue: "CBA",
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, int64(1), count2-count)

	// get
	domain, err := repo.Get(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "foo.bar", domain.Object)
	_, err = repo.Get(ctx, "test0")
	assert.NotNil(t, err)

	// query
	_, count3, err := repo.Query(ctx, "", 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(count3))
}
