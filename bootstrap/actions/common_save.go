package actions

import (
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/utils"
)

type commonSave struct {
	register *rule.MDAction
}

func newCommonSave() *commonSave {
	return &commonSave{
		register: &rule.MDAction{Code: "save", Widget: "common", Action: "save"},
	}
}

func (s *commonSave) Register() *rule.MDAction {
	return s.register
}

func (s *commonSave) Exec(flow *utils.FlowContext) {

}
