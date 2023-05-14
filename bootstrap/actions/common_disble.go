package actions

import (
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/utils"
)

type commonDisable struct {
}

func newCommonDisable() commonDisable {
	return commonDisable{}
}

func (s commonDisable) Register() rule.MDAction {
	return rule.MDAction{Code: "disable", Widget: "common", Action: "disable"}
}

func (s commonDisable) Exec(flow *utils.FlowContext) {

}
