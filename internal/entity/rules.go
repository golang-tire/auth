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
	Resource string
	Object   string
	Action   string
	Effect   string
}

func (r Rule) ToProto() *auth.Rule {
	c, _ := ptypes.TimestampProto(r.CreatedAt)
	u, _ := ptypes.TimestampProto(r.UpdatedAt)

	var effect auth.Effect = auth.Effect_DENY
	if r.Effect == "ALLOW" {
		effect = auth.Effect_ALLOW
	}

	rule := &auth.Rule{
		Uuid:      r.UUID,
		Role:      r.Role.Title,
		Object:    r.Object,
		Action:    r.Action,
		Domain:    r.Domain.Name,
		Resource:  r.Resource,
		Effect:    effect,
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
