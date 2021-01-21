package entity

import (
	"google.golang.org/protobuf/proto"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/golang-tire/auth/internal/pkg/pubsub"

	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang-tire/pkg/log"
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

func (r *Rule) AfterCreate(tx *gorm.DB) (err error) {
	b, err := r.Bytes()
	if err != nil {
		log.Error("encode create rule message to bytes failed", log.Err(err))
	}
	pubErr := pubsub.Publish("create-rule", message.NewMessage(r.UUID, b))
	if pubErr != nil {
		log.Error("send new-rule event failed", log.Err(pubErr))
	}
	return nil
}

func (r *Rule) AfterUpdate(tx *gorm.DB) (err error) {
	b, err := r.Bytes()
	if err != nil {
		log.Error("encode update rule message to bytes failed", log.Err(err))
	}

	pubErr := pubsub.Publish("update-rule", message.NewMessage(r.UUID, b))
	if pubErr != nil {
		log.Error("send new-rule event failed", log.Err(pubErr))
	}
	return
}

func (r *Rule) AfterDelete(tx *gorm.DB) (err error) {
	return
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

func (r Rule) Bytes() ([]byte, error) {
	bytes, err := proto.Marshal(r.ToProto())
	if err != nil {
		return nil, err
	}
	return bytes, err
}

func RuleToProtoList(rml []Rule) []*auth.Rule {
	var r []*auth.Rule
	for _, i := range rml {
		r = append(r, i.ToProto())
	}
	return r
}
