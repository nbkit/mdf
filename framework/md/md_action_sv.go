package md

import (
	"fmt"
	"github.com/ggoop/mdf/db"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/utils"
)

const (
	RuleType_Widget string = "widget"
	RuleType_Entity string = "entity"
)

// 注册器
type RuleRegister struct {
	Domain    string //领域：
	Code      string //规则编码：save,delete,query,find
	OwnerCode string //规则拥有者：common,widgetID,entityID
	OwnerType string //规则拥有者类型：widget,entity
}

func (s RuleRegister) GetKey() string {
	return fmt.Sprintf("%s:%s:%s", s.OwnerType, s.OwnerCode, s.Code)
}

/**
规则通用接口
*/
type IActionRule interface {
	Exec(token *utils.TokenContext, req *utils.ReqContext, res *utils.ResContext)
	Register() RuleRegister
}
type IActionSv interface {
	DoAction(token *utils.TokenContext, req *utils.ReqContext) *utils.ResContext
	RegisterRule(rules ...IActionRule)
	RegisterAction(rules ...IActionRule)
}

func ActionSv() IActionSv {
	return actionSv
}

type actionSvImpl struct {
	rules   map[string]IActionRule
	actions map[string]IActionRule
}

var actionSv IActionSv = newActionSvImpl()

func newActionSvImpl() *actionSvImpl {
	return &actionSvImpl{rules: make(map[string]IActionRule), actions: make(map[string]IActionRule)}
}
func (s actionSvImpl) GetRule(reg RuleRegister) (IActionRule, bool) {
	if r, ok := s.rules[reg.GetKey()]; ok {
		return r, ok
	}
	return nil, false
}

func (s actionSvImpl) GetAction(reg RuleRegister) (IActionRule, bool) {
	if r, ok := s.actions[reg.GetKey()]; ok {
		return r, ok
	}
	return nil, false
}

//执行命令
func (s actionSvImpl) DoAction(token *utils.TokenContext, req *utils.ReqContext) *utils.ResContext {
	res := &utils.ResContext{}
	commonId := "common"
	// 查找动作执行
	var action IActionRule
	if a, ok := s.GetAction(RuleRegister{OwnerType: req.OwnerType, OwnerCode: req.OwnerCode, Code: req.Action}); ok {
		action = a
	}
	if action == nil {
		if a, ok := s.GetAction(RuleRegister{OwnerType: req.OwnerType, OwnerCode: commonId, Code: req.Action}); ok {
			action = a
		}
	}
	if action != nil {
		if action.Exec(token, req, res); res.Error() != nil {
			return res
		}
	}
	//执行规则集合
	rules := make([]IActionRule, 0)
	ownerIds := []string{commonId, req.OwnerCode}
	ruleDatas := make([]MDActionRule, 0)
	//查询拥有者规则
	db.Default().Order("sequence,code").
		Where("owner_type=? and owner_code in (?) and action=? and enabled=1", req.OwnerType, ownerIds, req.Action).
		Find(&ruleDatas)

	if len(ruleDatas) > 0 {
		replacedList := make(map[string]MDActionRule)
		for _, r := range ruleDatas {
			if r.Replaced != "" {
				replacedList[r.Replaced] = r
			}
		}
		for _, r := range ruleDatas {
			if replaced, ok := replacedList[fmt.Sprintf("%s:%s", r.OwnerCode, r.Code)]; ok {
				glog.Error("规则被替换", glog.Any("replaced", replaced.Code))
				continue
			}
			if rule, ok := s.GetRule(RuleRegister{Domain: r.Domain, OwnerType: r.OwnerType, OwnerCode: r.OwnerCode, Code: r.Code}); ok {
				rules = append(rules, rule)
			} else {
				glog.Error("找不到规则", glog.Any("rule", r))
			}
		}
	}
	if len(rules) == 0 {
		glog.Error("没有找到任何规则可执行！")
	} else {
		for _, rule := range rules {
			if rule.Exec(token, req, res); res.Error() != nil {
				return res
			}
		}
	}
	return res
}
func (s actionSvImpl) RegisterRule(rules ...IActionRule) {
	if len(rules) > 0 {
		for i, _ := range rules {
			rule := rules[i]
			s.rules[rule.Register().GetKey()] = rule
		}
	}
}
func (s actionSvImpl) RegisterAction(rules ...IActionRule) {
	if len(rules) > 0 {
		for i, _ := range rules {
			rule := rules[i]
			s.actions[rule.Register().GetKey()] = rule
		}
	}
}
