package entity

import (
	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang/protobuf/ptypes"
	"gorm.io/gorm"
)

type Rule struct {
	gorm.Model
	UUID     string `gorm:"index"`
	RoleID   uint
	Role     Role `gorm:"foreignKey:RoleID"`
	DomainID uint
	Domain   Domain `gorm:"foreignKey:DomainID"`
	Object   string
	Action   string
}

func (r Rule) ToProto() *auth.Rule {
	c, _ := ptypes.TimestampProto(r.CreatedAt)
	u, _ := ptypes.TimestampProto(r.UpdatedAt)

	rule := &auth.Rule{
		Uuid:      r.UUID,
		Role:      r.Role.Title,
		Object:    r.Object,
		Action:    r.Action,
		Domain:    r.Domain.Name,
		CreatedAt: c,
		UpdatedAt: u,
	}
	return rule
}

func RuleToProtoList(rml []Rule) []*auth.Rule {
	var r []*auth.Rule
	for _, i := range rml {
		r = append(r, i.ToProto())
	}
	return r
}
