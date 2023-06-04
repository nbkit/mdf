package rule

import (
	"fmt"
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
	"reflect"
	"sort"
)

/**
规则通用接口
*/

type IAction interface {
	Exec(flow *utils.FlowContext)
	Register() MDAction
}
type IRule interface {
	Exec(flow *utils.FlowContext)
	Register() MDRule
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
func (s *actionSvImpl) getRuleSequence(rule MDRule) int {
	if rule.Sequence > 0 {
		return rule.Sequence
	}
	return 50
}
func (s *actionSvImpl) getRules(flow *utils.FlowContext) []IRule {
	ruleList := make([]IRule, 0)
	if len(s.rules) == 0 {
		return ruleList
	}
	replaceList := make(map[string]string)
	key := ""
	for _, rule := range s.rules {
		g := rule.Register()
		if (g.Action == flow.Request.Action || g.Action == "*") && g.Widget == flow.Request.Widget {
			key = fmt.Sprintf("%v:%v", flow.Request.Action, s.getRuleSequence(g))
			replaceList[key] = key
			ruleList = append(ruleList, rule)
		}
	}
	for _, rule := range s.rules {
		g := rule.Register()
		if g.Widget == mdCommonTag && g.Action == flow.Request.Action {
			key = fmt.Sprintf("%v:%v", flow.Request.Action, s.getRuleSequence(g))
			if _, ok := replaceList[key]; !ok {
				replaceList[key] = key
				ruleList = append(ruleList, rule)
			}
		}
	}
	sort.Slice(ruleList, func(i, j int) bool {
		return s.getRuleSequence(ruleList[i].Register()) < s.getRuleSequence(ruleList[j].Register())
	})
	return ruleList
}

// 执行命令
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
			if rule.Register().Action == "*" {
				s.execCommonRule(rule, flow)
			} else {
				rule.Exec(flow)
			}
			if flow.Error() != nil {
				return flow
			}
			if flow.Canceled() {
				log.ErrorF("请求已被规则%s终止！", rule.Register().Widget)
				return flow
			}
		}
	}
	return flow
}
func (s actionSvImpl) execCommonRule(rule IRule, flow *utils.FlowContext) {
	reflectValue := reflect.ValueOf(rule)
	methodName := utils.CamelString(flow.Request.Action)
	if reflectValue.CanAddr() && reflectValue.Kind() != reflect.Ptr {
		reflectValue = reflectValue.Addr()
	}
	if methodValue := reflectValue.MethodByName(methodName); methodValue.IsValid() {
		switch method := methodValue.Interface().(type) {
		case func():
			method()
		case func(*utils.FlowContext):
			method(flow)
		default:
			flow.Error(fmt.Errorf("unsupported function %v", methodName))
		}
	}
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
