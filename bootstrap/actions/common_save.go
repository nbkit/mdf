package actions

import (
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/utils"
)

type commonSave struct {
}

func newCommonSave() commonSave {
	return commonSave{}
}

func (s commonSave) Register() rule.MDAction {
	return rule.MDAction{Code: "save", Widget: "common", Action: "save"}
}

func (s commonSave) Exec(flow *utils.FlowContext) {

}
