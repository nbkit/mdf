package actions

import (
	"github.com/nbkit/mdf/framework/widget"
	"github.com/nbkit/mdf/utils"
)

type commonEnable struct {
	register *widget.MDAction
}

func newCommonEnable() *commonEnable {
	return &commonEnable{
		register: &widget.MDAction{Code: "enable", Widget: "common", Action: "enable"},
	}
}

func (s commonEnable) Register() *widget.MDAction {
	return s.register
}

func (s commonEnable) Exec(flow *utils.FlowContext) {

}
