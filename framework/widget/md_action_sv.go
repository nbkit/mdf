package widget

import (
	"fmt"
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
	"sort"
)

/**
规则通用接口
*/

type IAction interface {
	Exec(flow *utils.FlowContext)
	Register() *MDAction
}
type IRule interface {
	Exec(flow *utils.FlowContext)
	Register() *rule.MDRule
}
type IActionSv interface {
	DoAction(flow *utils.FlowContext) *utils.FlowContext
	RegisterRule(rules ...IRule)
	RegisterAction(rules ...IAction)
}

func ActionSv() IActionSv {
	return actionSv
}

var mdCommonTag string = "common"

type actionSvImpl struct {
	rules   map[string]IRule
	actions map[string]IAction
}

var actionSv IActionSv = newActionSvImpl()

func newActionSvImpl() *actionSvImpl {
	return &actionSvImpl{
		rules:   make(map[string]IRule),
		actions: make(map[string]IAction),
	}
}
func (s actionSvImpl) getAction(flow *utils.FlowContext) IAction {
	key := fmt.Sprintf("%s:%s", flow.Request.Widget, flow.Request.Action)
	if r, ok := s.actions[key]; ok {
		return r
	}
	key = fmt.Sprintf("%s:%s", mdCommonTag, flow.Request.Action)
	if r, ok := s.actions[key]; ok {
		return r
	}
	return nil
}
func (s *actionSvImpl) getRules(flow *utils.FlowContext) []IRule {
	ruleList := make([]IRule, 0)
	if len(s.rules) == 0 {
		return ruleList
	}
	replacedList := make(map[string]IRule)
	for i, _ := range s.rules {
		rule := s.rules[i]
		g := rule.Register()
		if g.Action != flow.Request.Action {
			continue
		}
		if g.Replaced != "" {
			replacedList[g.Replaced] = rule
		}
	}
	for i, _ := range s.rules {
		rule := s.rules[i]
		g := rule.Register()
		if g.Action != flow.Request.Action {
			continue
		}
		if replaced, ok := replacedList[fmt.Sprintf("%s:%s", g.Widget, g.Code)]; ok {
			log.Error().Any("replaced", replaced.Register().Code).Msg("规则被替换")
			continue
		}
		if g.Widget == mdCommonTag || g.Widget == flow.Request.Widget {
			ruleList = append(ruleList, rule)
		}
	}
	sort.Slice(ruleList, func(i, j int) bool {
		return ruleList[i].Register().Sequence < ruleList[j].Register().Sequence
	})
	return ruleList
}

//执行命令
func (s actionSvImpl) DoAction(flow *utils.FlowContext) *utils.FlowContext {
	// 查找动作执行
	if action := s.getAction(flow); action != nil {
		if action.Exec(flow); flow.Error() != nil {
			return flow
		}
	}
	//执行规则集合
	if rules := s.getRules(flow); len(rules) == 0 {
		log.Error().Msg("没有找到任何规则可执行！")
	} else {
		for _, rule := range rules {
			if rule.Exec(flow); flow.Error() != nil {
				return flow
			}
			if flow.Canceled() {
				log.ErrorF("请求已被规则%s终止！", rule.Register().Code)
				return flow
			}
		}
	}
	return flow
}
func (s actionSvImpl) RegisterRule(rules ...IRule) {
	if len(rules) > 0 {
		for i, _ := range rules {
			rule := rules[i]
			s.rules[rule.Register().GetKey()] = rule
		}
	}
}
func (s actionSvImpl) RegisterAction(rules ...IAction) {
	if len(rules) > 0 {
		for i, _ := range rules {
			rule := rules[i]
			s.actions[rule.Register().GetKey()] = rule
		}
	}
}
