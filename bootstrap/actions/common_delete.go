package actions

import (
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/utils"
)

type commonDelete struct {
}

func newCommonDelete() commonDelete {
	return commonDelete{}
}
func (s commonDelete) Register() rule.MDAction {
	return rule.MDAction{Code: "delete", Widget: "common", Action: "delete"}
}

func (s commonDelete) Exec(flow *utils.FlowContext) {

}
