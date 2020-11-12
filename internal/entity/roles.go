package entity

import (
	"time"

	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang/protobuf/ptypes"
)

type Role struct {
	tableName struct{} `pg:"roles,alias:role"` //nolint
	ID        uint64   `pg:",pk"`
	UUID      string
	Title     string `pg:",unique"`
	Enable    bool   `pg:"default:FALSE,notnull,use_zero"`
	CreatedAt time.Time
	UpdatedAt time.Time
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

func RoleFromProto(role *auth.Role) Role {
	c, _ := ptypes.Timestamp(role.CreatedAt)
	u, _ := ptypes.Timestamp(role.UpdatedAt)

	return Role{
		UUID:      role.Uuid,
		Title:     role.Title,
		CreatedAt: c,
		UpdatedAt: u,
	}
}
