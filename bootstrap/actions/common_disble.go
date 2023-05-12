package actions

import (
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/utils"
)

type commonDisable struct {
	register *rule.MDAction
}

func newCommonDisable() *commonDisable {
	return &commonDisable{
		register: &rule.MDAction{Code: "disable", Widget: "common", Action: "disable"},
	}
}

func (s commonDisable) Register() *rule.MDAction {
	return s.register
}

func (s commonDisable) Exec(flow *utils.FlowContext) {

}
