package service

import (
	"context"
	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/gogf/gf/frame/g"
	"gowebserver/app/common/dao"
	"gowebserver/app/common/model"
	"os"
	"sync"
)

//ICasbinRuleManager casbin规则管理接口
type ICasbinRuleManager interface {
}

//casbinRuleManagerImpl casbin rule
type casbinRuleManagerImpl struct{}

var (
	casbinService = casbinRuleManagerImpl{}
	syncOnce      sync.Once
	casbinAdapter *adapterCasbin
)

type adapterCasbin struct {
	Enforcer    *casbin.SyncedEnforcer
	EnforcerErr error
	ctx         context.Context
}

//CasbinEnforcer 获取casbinAdapter单例对象
func CasbinEnforcer(ctx context.Context) (enforcer *casbin.SyncedEnforcer, err error) {
	syncOnce.Do(func() {
		casbinAdapter = casbinService.newAdapter(ctx)
	})
	enforcer = casbinAdapter.Enforcer
	err = casbinAdapter.EnforcerErr
	return
}

//ReloadDBCasbinPolicy 重新加载数据权限规则
func ReloadDBCasbinPolicy() error {
	enforcer, err := CasbinEnforcer(context.TODO())
	if err != nil {
		g.Log().Error(err)
		return err
	}
	if err = enforcer.LoadPolicy(); err != nil {
		g.Log().Error(err)
		return err
	}
	return nil
}

//初始化casbinAdapter单例对象
func (s *casbinRuleManagerImpl) newAdapter(ctx context.Context) (a *adapterCasbin) {
	a = new(adapterCasbin)
	a.initPolicy(ctx)
	a.ctx = context.Background()
	return
}

//初始化casbin_rule
func (a *adapterCasbin) initPolicy(ctx context.Context) {
	currentDir, _ := os.Getwd()
	casbinModelFile := g.Cfg().GetString("casbin.model")
	g.Log().Noticef("currentPath: %s casbin.model file: %s", currentDir, casbinModelFile)

	//这里通过数据库适配器的方式初始化
	//内部会自动调用适配器实现的 LoadPolicy 接口加载casbin_rule表中的策略
	e, err := casbin.NewSyncedEnforcer(casbinModelFile, a)
	if err != nil {
		a.EnforcerErr = err
		g.Log().Error(err)
		return
	}

	// Debug
	if g.Cfg().GetBool("debug.enable") {
		//e.EnableLog(true)
		e.GetModel().PrintModel()
	}

	// set adapter
	e.SetAdapter(a)

	// Clear the current policy.
	e.ClearPolicy()
	a.Enforcer = e

	// Load the policy from DB.
	err = a.LoadPolicy(e.GetModel())
	if err != nil {
		a.EnforcerErr = err
		return
	}
}

// SavePolicy saves policy to database.
func (a *adapterCasbin) SavePolicy(model casbinModel.Model) (err error) {
	err = a.dropTable()
	if err != nil {
		return
	}

	err = a.createTable()
	if err != nil {
		return
	}

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			_, err = dao.CasbinRule.Ctx(a.ctx).Data(line).Insert()
			if err != nil {
				return err
			}
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			_, err = dao.CasbinRule.Ctx(a.ctx).Data(line).Insert()
			if err != nil {
				return err
			}
		}
	}
	return
}

func (a *adapterCasbin) dropTable() (err error) {
	return
}

func (a *adapterCasbin) createTable() (err error) {
	return
}

// LoadPolicy loads policy from database.
func (a *adapterCasbin) LoadPolicy(casbinModel casbinModel.Model) error {
	var lines []*model.CasbinRule
	if err := dao.CasbinRule.Ctx(a.ctx).Scan(&lines); err != nil {
		return err
	}
	for _, line := range lines {
		loadPolicyLine(line, casbinModel)
	}
	return nil
}

// AddPolicy adds a policy rule to the storage.
func (a *adapterCasbin) AddPolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	_, err := dao.CasbinRule.Ctx(a.ctx).Data(line).Insert()
	return err
}

// RemovePolicy removes a policy rule from the storage.
func (a *adapterCasbin) RemovePolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	err := rawDelete(a, line)
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *adapterCasbin) RemoveFilteredPolicy(sec string, ptype string,
	fieldIndex int, fieldValues ...string) error {
	line := &model.CasbinRule{}
	line.Ptype = ptype
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
	}
	err := rawDelete(a, line)
	return err
}

func loadPolicyLine(line *model.CasbinRule, casbinModel casbinModel.Model) {
	lineText := line.Ptype
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
	persist.LoadPolicyLine(lineText, casbinModel)
}

func savePolicyLine(ptype string, rule []string) *model.CasbinRule {
	line := &model.CasbinRule{}
	line.Ptype = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}
	return line
}

func rawDelete(a *adapterCasbin, line *model.CasbinRule) error {
	db := dao.CasbinRule.Ctx(a.ctx).Where("ptype = ?", line.Ptype)
	if line.V0 != "" {
		db = db.Where("v0 = ?", line.V0)
	}
	if line.V1 != "" {
		db = db.Where("v1 = ?", line.V1)
	}
	if line.V2 != "" {
		db = db.Where("v2 = ?", line.V2)
	}
	if line.V3 != "" {
		db = db.Where("v3 = ?", line.V3)
	}
	if line.V4 != "" {
		db = db.Where("v4 = ?", line.V4)
	}
	if line.V5 != "" {
		db = db.Where("v5 = ?", line.V5)
	}
	_, err := db.Delete()
	return err
}
