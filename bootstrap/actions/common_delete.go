package actions

import (
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/utils"
)

type CommonDelete struct {
	register *rule.MDAction
}

func newCommonDelete() *CommonDelete {
	return &CommonDelete{
		register: &rule.MDAction{Code: "delete", Widget: "common", Action: "delete"},
	}
}
func (s CommonDelete) Register() *rule.MDAction {
	return s.register
}

func (s CommonDelete) Exec(flow *utils.FlowContext) {

}
