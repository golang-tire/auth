package audit_logs

import (
	"context"
	"testing"

	"github.com/golang-tire/auth/internal/entity"

	"github.com/golang-tire/auth/internal/users"

	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/stretchr/testify/assert"
)

func Test_service_CR(t *testing.T) {

	userRepo := users.NewMockRepository()
	s := NewService(&mockRepository{}, userRepo)
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

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, int64(0), count)

	// successful create
	auditLog, err := s.Create(ctx, &auth.CreateAuditLogRequest{
		UserUuid: userUuid,
		Action:   "POST",
		Object:   "product",
		OldValue: "ABC",
		NewValue: "ACB",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, auditLog.Uuid)
	id := auditLog.Uuid
	assert.Equal(t, "product", auditLog.Object)
	assert.NotEmpty(t, auditLog.CreatedAt)
	assert.NotEmpty(t, auditLog.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.Create(ctx, &auth.CreateAuditLogRequest{
		UserUuid: userUuid,
		Action:   "",
		Object:   "",
		OldValue: "ABC",
		NewValue: "ACB",
	})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.Create(ctx, &auth.CreateAuditLogRequest{
		UserUuid: userUuid,
		Action:   "POST",
		Object:   "error",
		OldValue: "ABC",
		NewValue: "ACB",
	})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.Create(ctx, &auth.CreateAuditLogRequest{
		UserUuid: userUuid,
		Action:   "POST",
		Object:   "cars",
		OldValue: "qwer",
		NewValue: "rewq",
	})

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	auditLog, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "product", auditLog.Object)
	assert.Equal(t, id, auditLog.Uuid)

	// query
	_domains, _ := s.Query(ctx, "", 0, 0)
	assert.Equal(t, 2, int(_domains.TotalCount))
}
