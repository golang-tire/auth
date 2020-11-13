package entity

import (
	"time"

	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang/protobuf/ptypes"
)

type User struct {
	tableName struct{} `pg:"users,alias:user"` //nolint
	ID        uint64   `pg:",pk"`
	UUID      string   `pg:",unique"`
	Firstname string
	Lastname  string
	Username  string `pg:",unique"`
	Password  string
	Gender    string
	AvatarURL string
	Email     string `pg:",unique"`
	Enable    bool   `pg:"default:FALSE,notnull,use_zero"`
	RawData   string
	UserRoles []*UserRole `pg:"rel:has-many"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRole struct {
	tableName struct{} `pg:"user_roles,alias:user_roles"` //nolint
	ID        uint64   `pg:",pk"`
	UUID      string   `pg:",unique"`
	RoleId    uint64
	Role      *Role `pg:"rel:has-one"`
	UserID    uint64
	User      *User `pg:"rel:has-one"`
	DomainId  uint64
	Domain    *Domain `pg:"rel:has-one"`
	Enable    bool    `pg:"default:TRUE,notnull,use_zero"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r User) ToProto() *auth.User {
	c, _ := ptypes.TimestampProto(r.CreatedAt)
	u, _ := ptypes.TimestampProto(r.UpdatedAt)

	user := &auth.User{
		Uuid:      r.UUID,
		Firstname: r.Firstname,
		Lastname:  r.Lastname,
		Gender:    r.Gender,
		AvatarUrl: r.AvatarURL,
		Username:  r.Username,
		Email:     r.Email,
		Enable:    r.Enable,
		RawData:   r.RawData,
		Roles:     UserRoleToProtoList(r.UserRoles),
		CreatedAt: c,
		UpdatedAt: u,
	}
	return user
}

func UserToProtoList(rml []User) []*auth.User {
	var r []*auth.User
	for _, i := range rml {
		r = append(r, i.ToProto())
	}
	return r
}

func UserFromProto(user *auth.User) User {
	c, _ := ptypes.Timestamp(user.CreatedAt)
	u, _ := ptypes.Timestamp(user.UpdatedAt)

	return User{
		UUID:      user.Uuid,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Gender:    user.Gender,
		AvatarURL: user.AvatarUrl,
		Username:  user.Username,
		Email:     user.Email,
		Enable:    user.Enable,
		RawData:   user.RawData,
		UserRoles: UserRoleListFromProto(user.Roles),
		CreatedAt: c,
		UpdatedAt: u,
	}
}

func (dr UserRole) ToProto() *auth.UserRole {
	c, _ := ptypes.TimestampProto(dr.CreatedAt)
	u, _ := ptypes.TimestampProto(dr.UpdatedAt)
	domainRole := &auth.UserRole{
		Uuid:      dr.UUID,
		Role:      dr.Role.ToProto(),
		Enable:    dr.Enable,
		CreatedAt: c,
		UpdatedAt: u,
	}
	return domainRole
}

func UserRoleToProtoList(rml []*UserRole) []*auth.UserRole {
	var r []*auth.UserRole
	for _, i := range rml {
		r = append(r, i.ToProto())
	}
	return r
}

func UserRoleFromProto(domainRole *auth.UserRole) UserRole {
	c, _ := ptypes.Timestamp(domainRole.CreatedAt)
	u, _ := ptypes.Timestamp(domainRole.UpdatedAt)

	return UserRole{
		UUID:      domainRole.Uuid,
		Enable:    domainRole.Enable,
		CreatedAt: c,
		UpdatedAt: u,
	}
}

func UserRoleListFromProto(domainRoles []*auth.UserRole) []*UserRole {
	var d []*UserRole
	for _, i := range domainRoles {
		dr := UserRoleFromProto(i)
		d = append(d, &dr)
	}
	return d
}
