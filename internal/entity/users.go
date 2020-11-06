package entity

import (
	"time"

	"github.com/go-pg/pg/v10/orm"

	auth "github.com/golang-tire/auth/internal/proto/v1"
	"github.com/golang/protobuf/ptypes"
)

type User struct {
	tableName   struct{} `pg:"users,alias:u"` //nolint
	ID          uint64   `pg:",pk"`
	UUID        string   `pg:"default:gen_random_uuid()"`
	Firstname   string
	Lastname    string
	Username    string
	Password    string
	Gender      string
	AvatarURL   string
	Email       string
	Enable      bool
	RawData     string
	Rules       []*Rule       `pg:"many2many:user_rules"`
	DomainRoles []*DomainRole `pg:"rel:has-many"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserRule struct {
	tableName struct{} `pg:"user_rules,alias:ur"` //nolint
	RuleId    uint64
	UserId    uint64
}

type DomainRole struct {
	tableName struct{} `pg:"user_domain_roles,alias:dr"` //nolint
	ID        uint64   `pg:",pk"`
	UUID      string   `pg:"default:gen_random_uuid()"`
	RoleId    uint64
	Role      *Role `pg:"rel:has-one"`
	UserID    uint64
	DomainId  uint64
	Domain    *Domain `pg:"rel:has-one"`
	Enable    bool
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
		UUID:        user.Uuid,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Gender:      user.Gender,
		AvatarURL:   user.AvatarUrl,
		Username:    user.Username,
		Email:       user.Email,
		Enable:      user.Enable,
		RawData:     user.RawData,
		DomainRoles: DomainRoleListFromProto(user.DomainRoles),
		Rules:       RuleListFromProto(user.Rules),
		CreatedAt:   c,
		UpdatedAt:   u,
	}
}

func (dr DomainRole) ToProto() *auth.DomainRole {
	c, _ := ptypes.TimestampProto(dr.CreatedAt)
	u, _ := ptypes.TimestampProto(dr.UpdatedAt)
	domainRole := &auth.DomainRole{
		Uuid:      dr.UUID,
		Role:      dr.Role.ToProto(),
		Enable:    dr.Enable,
		CreatedAt: c,
		UpdatedAt: u,
	}
	return domainRole
}

func DomainRoleToProtoList(rml []*DomainRole) []*auth.DomainRole {
	var r []*auth.DomainRole
	for _, i := range rml {
		r = append(r, i.ToProto())
	}
	return r
}

func DomainRoleFromProto(domainRole *auth.DomainRole) DomainRole {
	c, _ := ptypes.Timestamp(domainRole.CreatedAt)
	u, _ := ptypes.Timestamp(domainRole.UpdatedAt)

	return DomainRole{
		UUID:      domainRole.Uuid,
		Enable:    domainRole.Enable,
		CreatedAt: c,
		UpdatedAt: u,
	}
}

func DomainRoleListFromProto(domainRoles []*auth.DomainRole) []*DomainRole {
	var d []*DomainRole
	for _, i := range domainRoles {
		dr := DomainRoleFromProto(i)
		d = append(d, &dr)
	}
	return d
}

func init() {
	// Register many to many model so ORM can better recognize m2m relation.
	// This should be done before dependant models are used.
	orm.RegisterTable((*UserRule)(nil))
}
