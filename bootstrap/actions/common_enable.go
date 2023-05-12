package actions

import (
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/utils"
)

type commonEnable struct {
	register *rule.MDAction
}

func newCommonEnable() *commonEnable {
	return &commonEnable{
		register: &rule.MDAction{Code: "enable", Widget: "common", Action: "enable"},
	}
}

func (s commonEnable) Register() *rule.MDAction {
	return s.register
}

func (s commonEnable) Exec(flow *utils.FlowContext) {

}
