package entity

import (
	"time"

	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang/protobuf/ptypes"
)

type Rule struct {
	tableName struct{} `pg:"rules,alias:ru"` //nolint
	ID        uint64   `pg:",pk"`
	UUID      string   `pg:"default:gen_random_uuid()"`
	Subject   string
	DomainId  uint64
	Domain    *Domain `pg:"rel:has-one"`
	Object    string
	Action    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r Rule) ToProto() *auth.Rule {
	c, _ := ptypes.TimestampProto(r.CreatedAt)
	u, _ := ptypes.TimestampProto(r.UpdatedAt)

	rule := &auth.Rule{
		Uuid:      r.UUID,
		Subject:   r.Subject,
		Object:    r.Object,
		Action:    r.Action,
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

func RuleFromProto(rule *auth.Rule) Rule {
	c, _ := ptypes.Timestamp(rule.CreatedAt)
	u, _ := ptypes.Timestamp(rule.UpdatedAt)

	return Rule{
		UUID:      rule.Uuid,
		Subject:   rule.Subject,
		Object:    rule.Object,
		Action:    rule.Action,
		CreatedAt: c,
		UpdatedAt: u,
	}
}

func RuleListFromProto(rules []*auth.Rule) []*Rule {
	var d []*Rule
	for _, i := range rules {
		dr := RuleFromProto(i)
		d = append(d, &dr)
	}
	return d
}
