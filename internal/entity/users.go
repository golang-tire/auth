package entity

import (
	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang/protobuf/ptypes"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID      string `gorm:"index"`
	Firstname string
	Lastname  string
	Username  string `gorm:"index"`
	Password  string
	Gender    string
	AvatarURL string
	Email     string `gorm:"unique"`
	Enable    bool
	RawData   string
	UserRoles []UserRole
}

type UserRole struct {
	gorm.Model
	UUID     string `gorm:"index"`
	RoleID   uint
	Role     Role
	User     User
	UserID   uint
	DomainID uint
	Domain   Domain
	Enable   bool
}

func (r User) ToProto(secure bool) *auth.User {
	c, _ := ptypes.TimestampProto(r.CreatedAt)
	u, _ := ptypes.TimestampProto(r.UpdatedAt)

	user := &auth.User{
		Uuid:      r.UUID,
		Firstname: r.Firstname,
		Lastname:  r.Lastname,
		Gender:    r.Gender,
		AvatarUrl: r.AvatarURL,
		Username:  r.Username,
		Password:  "",
		Email:     r.Email,
		Enable:    r.Enable,
		RawData:   r.RawData,
		Roles:     UserRoleToProtoList(r.UserRoles),
		CreatedAt: c,
		UpdatedAt: u,
	}

	if !secure {
		user.Password = r.Password
	}

	return user
}

func UserToProtoList(rml []User) []*auth.User {
	var r []*auth.User
	for _, i := range rml {
		r = append(r, i.ToProto(true))
	}
	return r
}

func (dr UserRole) ToProto() *auth.UserRole {
	c, _ := ptypes.TimestampProto(dr.CreatedAt)
	u, _ := ptypes.TimestampProto(dr.UpdatedAt)
	domainRole := &auth.UserRole{
		Uuid:      dr.UUID,
		Role:      dr.Role.Title,
		Domain:    dr.Domain.Name,
		Enable:    dr.Enable,
		CreatedAt: c,
		UpdatedAt: u,
	}
	return domainRole
}

func UserRoleToProtoList(rml []UserRole) []*auth.UserRole {
	var r []*auth.UserRole
	for _, i := range rml {
		r = append(r, i.ToProto())
	}
	return r
}
