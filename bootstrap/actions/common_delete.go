package actions

import (
	"github.com/nbkit/mdf/framework/widget"
	"github.com/nbkit/mdf/utils"
)

type CommonDelete struct {
	register *widget.MDAction
}

func newCommonDelete() *CommonDelete {
	return &CommonDelete{
		register: &widget.MDAction{Code: "delete", Widget: "common", Action: "delete"},
	}
}
func (s CommonDelete) Register() *widget.MDAction {
	return s.register
}

func (s CommonDelete) Exec(flow *utils.FlowContext) {

}
