package entity

import (
	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang/protobuf/ptypes"
	"gorm.io/gorm"
)

type App struct {
	gorm.Model
	UUID   string `gorm:"index"`
	Name   string
	Enable bool
}

type Resource struct {
	gorm.Model
	UUID  string `gorm:"index"`
	Name  string `gorm:"index:idx_app_resource"`
	AppID uint   `gorm:"index:idx_app_resource"`
	App   App
}

type Object struct {
	gorm.Model
	UUID       string `gorm:"index"`
	Identifier string `gorm:"index:idx_app_object"`
	AppID      uint   `gorm:"index:idx_app_object"`
	App        App
}

func (ap App) ToProto() *auth.App {
	c, _ := ptypes.TimestampProto(ap.CreatedAt)
	u, _ := ptypes.TimestampProto(ap.UpdatedAt)

	app := &auth.App{
		Uuid:      ap.UUID,
		Name:      ap.Name,
		Enable:    ap.Enable,
		CreatedAt: c,
		UpdatedAt: u,
	}
	return app
}

func AppToProtoList(apl []App) []*auth.App {
	var a []*auth.App
	for _, i := range apl {
		a = append(a, i.ToProto())
	}
	return a
}

func (r Resource) ToProto() *auth.Resource {
	c, _ := ptypes.TimestampProto(r.CreatedAt)
	u, _ := ptypes.TimestampProto(r.UpdatedAt)

	resource := &auth.Resource{
		Uuid:      r.UUID,
		Name:      r.Name,
		App:       r.App.ToProto(),
		CreatedAt: c,
		UpdatedAt: u,
	}
	return resource
}

func ResourceToProtoList(rl []Resource) []*auth.Resource {
	var r []*auth.Resource
	for _, i := range rl {
		r = append(r, i.ToProto())
	}
	return r
}

func (o Object) ToProto() *auth.Object {
	c, _ := ptypes.TimestampProto(o.CreatedAt)
	u, _ := ptypes.TimestampProto(o.UpdatedAt)

	object := &auth.Object{
		Uuid:       o.UUID,
		Identifier: o.Identifier,
		App:        o.App.ToProto(),
		CreatedAt:  c,
		UpdatedAt:  u,
	}
	return object
}

func ObjectToProtoList(ol []Object) []*auth.Object {
	var o []*auth.Object
	for _, i := range ol {
		o = append(o, i.ToProto())
	}
	return o
}
