package entity

import (
	"gorm.io/gorm"

	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang/protobuf/ptypes"
)

type Role struct {
	gorm.Model
	UUID   string `gorm:"index"`
	Title  string `gorm:"index"`
	Enable bool
}

func (rm Role) ToProto() *auth.Role {
	c, _ := ptypes.TimestampProto(rm.CreatedAt)
	u, _ := ptypes.TimestampProto(rm.UpdatedAt)

	role := &auth.Role{
		Uuid:      rm.UUID,
		Title:     rm.Title,
		Enable:    rm.Enable,
		CreatedAt: c,
		UpdatedAt: u,
	}
	return role
}

func RoleToProtoList(rml []Role) []*auth.Role {
	var r []*auth.Role
	for _, i := range rml {
		r = append(r, i.ToProto())
	}
	return r
}
