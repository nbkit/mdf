package actions

import (
	"github.com/nbkit/mdf/framework/widget"
	"github.com/nbkit/mdf/utils"
)

type commonDisable struct {
	register *widget.MDAction
}

func newCommonDisable() *commonDisable {
	return &commonDisable{
		register: &widget.MDAction{Code: "disable", Widget: "common", Action: "disable"},
	}
}

func (s commonDisable) Register() *widget.MDAction {
	return s.register
}

func (s commonDisable) Exec(flow *utils.FlowContext) {

}
