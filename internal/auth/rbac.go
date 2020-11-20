package auth

import (
	"context"
	"strings"

	"github.com/golang-tire/pkg/pubsub"
	"google.golang.org/protobuf/proto"

	zaplogger "github.com/casbin/zap-logger"

	"github.com/golang-tire/pkg/log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/golang-tire/auth/internal/rules"
	"github.com/golang-tire/auth/internal/users"
	"github.com/golang-tire/pkg/config"
)

var (
	rbacConfig = config.RegisterString("rbac.conf", "")
)

type rbacService struct {
	enforcer *casbin.Enforcer
}

type adapter struct {
	lines    []line
	ctx      context.Context
	usersSrv users.Service
	rulesSrv rules.Service
}

type line struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

type Policy struct {
	Subject  string
	Domain   string
	Resource string
	Action   string
	Object   string
	Effect   string
}

type Group struct {
	Subject string
	Role    string
	Domain  string
}

func newAdapter(ctx context.Context, ruleSrv rules.Service, userSrv users.Service) persist.Adapter {
	return &adapter{
		lines:    []line{},
		ctx:      ctx,
		usersSrv: userSrv,
		rulesSrv: ruleSrv,
	}
}

func (a *adapter) loadFromDb() error {
	userRoles, err := a.usersSrv.ListUserRoles(a.ctx)
	if err != nil {
		return err
	}

	for _, ur := range userRoles {
		if ur.Domain.Name == "" {
			ur.Domain.Name = "*"
		}
		a.lines = append(a.lines, line{
			PType: "g",
			V0:    ur.User.Username,
			V1:    ur.Role.Title,
			V2:    ur.Domain.Name,
		})
	}

	ruleItems, err := a.rulesSrv.All(a.ctx)
	for _, rule := range ruleItems {

		if rule.Domain.Name == "" {
			rule.Domain.Name = "*"
		}

		a.lines = append(a.lines, line{
			PType: "p",
			V0:    rule.Role.Title,
			V1:    rule.Domain.Name,
			V2:    rule.Resource,
			V3:    rule.Action,
			V4:    rule.Object,
			V5:    strings.ToLower(rule.Effect),
		})
	}
	return nil
}

func loadPolicyLine(line line, model model.Model) {
	lineText := line.PType
	if line.V0 != "" {
		lineText += ", " + line.V0
	}
	if line.V1 != "" {
		lineText += ", " + line.V1
	}
	if line.V2 != "" {
		lineText += ", " + line.V2
	}
	if line.V3 != "" {
		lineText += ", " + line.V3
	}
	if line.V4 != "" {
		lineText += ", " + line.V4
	}
	if line.V5 != "" {
		lineText += ", " + line.V5
	}

	persist.LoadPolicyLine(lineText, model)
}

func (a *adapter) LoadPolicy(model model.Model) error {

	err := a.loadFromDb()
	if err != nil {
		return err
	}
	for _, line := range a.lines {
		loadPolicyLine(line, model)
	}
	return nil
}

func (a *adapter) SavePolicy(model model.Model) error {
	panic("implement me")
}

func (a *adapter) AddPolicy(sec string, ptype string, rule []string) error {
	panic("implement me")
}

func (a *adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	panic("implement me")
}

func (a *adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	panic("implement me")
}

func (a *rbacService) OnPolicyChange(ctx context.Context, msg proto.Message) {
	err := a.enforcer.LoadPolicy()
	if err != nil {
		log.Error("reload polices failed", log.Err(err))
	}
}

func InitRbac(ctx context.Context, rulesSrv rules.Service, usersSrv users.Service) (*casbin.Enforcer, error) {

	log.Info("init rbac module")
	err := config.Load()
	if err != nil {
		return nil, err
	}

	a := newAdapter(ctx, rulesSrv, usersSrv)
	m, err := model.NewModelFromString(rbacConfig.String())
	if err != nil {
		return nil, err
	}

	logger := zaplogger.NewLoggerByZap(log.Logger(), true)
	enf, err := casbin.NewEnforcer(m, a, logger, true)

	rbacSrv := &rbacService{enforcer: enf}

	pubsub.Get().Subscribe(ctx, "rule-change", rbacSrv.OnPolicyChange)
	pubsub.Get().Subscribe(ctx, "user-change", rbacSrv.OnPolicyChange)
	return enf, err
}
