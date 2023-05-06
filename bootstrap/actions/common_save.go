package actions

import (
	"github.com/nbkit/mdf/framework/widget"
	"github.com/nbkit/mdf/utils"
)

type commonSave struct {
	register *widget.MDAction
}

func newCommonSave() *commonSave {
	return &commonSave{
		register: &widget.MDAction{Code: "save", Widget: "common", Action: "save"},
	}
}

func (s *commonSave) Register() *widget.MDAction {
	return s.register
}

func (s *commonSave) Exec(flow *utils.FlowContext) {

}
