package entity

import (
	"time"

	auth "github.com/golang-tire/auth/protobuf"
	"github.com/golang/protobuf/ptypes"
)

type Rule struct {
	tableName struct{} `pg:"rules,alias:ru"` //nolint
	ID        uint64   `pg:",pk"`
	UUID      string   `pg:"default:gen_random_uuid()"`
	Subject   string
	Domain    string
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
		Domain:    r.Domain,
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
		Domain:    rule.Domain,
		Object:    rule.Object,
		Action:    rule.Action,
		CreatedAt: c,
		UpdatedAt: u,
	}
}
