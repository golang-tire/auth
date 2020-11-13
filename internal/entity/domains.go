package entity

import (
	"gorm.io/gorm"

	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang/protobuf/ptypes"
)

type Domain struct {
	gorm.Model
	UUID   string `gorm:"index"`
	Name   string `gorm:"index"`
	Enable bool
}

func (dm Domain) ToProto() *auth.Domain {
	c, _ := ptypes.TimestampProto(dm.CreatedAt)
	u, _ := ptypes.TimestampProto(dm.UpdatedAt)

	role := &auth.Domain{
		Uuid:      dm.UUID,
		Name:      dm.Name,
		Enable:    dm.Enable,
		CreatedAt: c,
		UpdatedAt: u,
	}
	return role
}

func DomainToProtoList(dml []Domain) []*auth.Domain {
	var r []*auth.Domain
	for _, i := range dml {
		r = append(r, i.ToProto())
	}
	return r
}
