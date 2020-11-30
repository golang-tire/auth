package entity

import (
	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang/protobuf/ptypes"
	"gorm.io/gorm"
)

type AuditLog struct {
	gorm.Model
	UUID     string `gorm:"index"`
	User     User   `gorm:"foreignKey:UserID"`
	UserID   uint
	Action   string
	Object   string
	OldValue string
	NewValue string
}

func (al AuditLog) ToProto() *auth.AuditLog {
	c, _ := ptypes.TimestampProto(al.CreatedAt)
	u, _ := ptypes.TimestampProto(al.UpdatedAt)

	role := &auth.AuditLog{
		Uuid:      al.UUID,
		User:      al.User.ToProto(true),
		Action:    al.Action,
		Object:    al.Object,
		OldValue:  al.OldValue,
		NewValue:  al.NewValue,
		CreatedAt: c,
		UpdatedAt: u,
	}
	return role
}

func AuditLogToProtoList(dml []AuditLog) []*auth.AuditLog {
	var r []*auth.AuditLog
	for _, i := range dml {
		r = append(r, i.ToProto())
	}
	return r
}
