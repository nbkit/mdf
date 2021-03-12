package md

import (
	"fmt"
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
	"sort"
)

// 注册器
type RuleRegister struct {
	Domain string //领域：
	Code   string //规则编码：save,delete,query,find
	Widget string //规则拥有者：common,widgetID,entityID
}

func (s RuleRegister) GetKey() string {
	return fmt.Sprintf("%s:%s", s.Widget, s.Code)
}

/**
规则通用接口
*/
type IActionRule interface {
	Exec(flow *utils.FlowContext)
	Register() RuleRegister
}
type IActionSv interface {
	DoAction(flow *utils.FlowContext) *utils.FlowContext
	RegisterRule(rules ...IActionRule)
	RegisterAction(rules ...IActionRule)
	Cache()
}

func ActionSv() IActionSv {
	return actionSv
}

var mdCommonTag string = "common"

type actionSvImpl struct {
	rules   map[string]IActionRule
	actions map[string]IActionRule

	mdRules map[string]*MDActionRule
}

var actionSv IActionSv = newActionSvImpl()

func newActionSvImpl() *actionSvImpl {
	return &actionSvImpl{
		rules:   make(map[string]IActionRule),
		actions: make(map[string]IActionRule),
		mdRules: make(map[string]*MDActionRule),
	}
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

func (s *actionSvImpl) Cache() {
	s.mdRules = make(map[string]*MDActionRule)
	ruleDatas := make([]MDActionRule, 0)
	db.Default().Where("enabled=1").Find(&ruleDatas)
	for i, _ := range ruleDatas {
		rule := ruleDatas[i]
		s.mdRules[fmt.Sprintf("%s:%s:%s", rule.Widget, rule.Action, rule.Code)] = &rule
	}
}
func (s *actionSvImpl) getActionRule(flow *utils.FlowContext) []MDActionRule {
	ruleList := make([]MDActionRule, 0)
	for _, r := range s.mdRules {
		if r.Action == flow.Request.Action && (r.Widget == mdCommonTag || r.Widget == flow.Request.Widget) {
			ruleList = append(ruleList, *r)
		}
	}
	sort.Slice(ruleList, func(i, j int) bool {
		return ruleList[i].Sequence < ruleList[i].Sequence
	})
	return ruleList
}

//执行命令
func (s actionSvImpl) DoAction(flow *utils.FlowContext) *utils.FlowContext {
	// 查找动作执行
	var action IActionRule
	if a, ok := s.GetAction(RuleRegister{Widget: flow.Request.Widget, Code: flow.Request.Action}); ok {
		action = a
	}
	if action == nil {
		if a, ok := s.GetAction(RuleRegister{Widget: mdCommonTag, Code: flow.Request.Action}); ok {
			action = a
		}
	}
	if action != nil {
		if action.Exec(flow); flow.Error() != nil {
			return flow
		}
	}
	//执行规则集合
	rules := make([]IActionRule, 0)
	ruleDatas := s.getActionRule(flow)

	if len(ruleDatas) > 0 {
		replacedList := make(map[string]MDActionRule)
		for _, r := range ruleDatas {
			if r.Replaced != "" {
				replacedList[r.Replaced] = r
			}
		}
		for _, r := range ruleDatas {
			if replaced, ok := replacedList[fmt.Sprintf("%s:%s", r.Widget, r.Code)]; ok {
				log.Error("规则被替换", log.Any("replaced", replaced.Code))
				continue
			}
			if rule, ok := s.GetRule(RuleRegister{Domain: r.Domain, Widget: r.Widget, Code: r.Code}); ok {
				rules = append(rules, rule)
			} else {
				log.Error("找不到规则", log.Any("rule", r))
			}
		}
	}
	if len(rules) == 0 {
		log.Error("没有找到任何规则可执行！")
	} else {
		for _, rule := range rules {
			if rule.Exec(flow); flow.Error() != nil {
				return flow
			}
			if flow.Canceled() {
				log.Errorf("请求已被规则%s终止！", rule.Register().Code)
				return flow
			}
		}
	}
	return flow
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
