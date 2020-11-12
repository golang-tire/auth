package entity

import (
	"time"

	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang/protobuf/ptypes"
)

type Domain struct {
	tableName struct{} `pg:"domains,alias:domain"` //nolint
	ID        uint64   `pg:",pk"`
	UUID      string
	Name      string `pg:",unique"`
	Enable    bool   `pg:"default:FALSE,notnull,use_zero"`
	CreatedAt time.Time
	UpdatedAt time.Time
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

func DomainFromProto(domain *auth.Domain) Domain {
	c, _ := ptypes.Timestamp(domain.CreatedAt)
	u, _ := ptypes.Timestamp(domain.UpdatedAt)

	return Domain{
		UUID:      domain.Uuid,
		Name:      domain.Name,
		Enable:    domain.Enable,
		CreatedAt: c,
		UpdatedAt: u,
	}
}
