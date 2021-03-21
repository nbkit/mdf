package actions

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type commonDisable struct {
	register *md.MDAction
}

func newCommonDisable() *commonDisable {
	return &commonDisable{
		register: &md.MDAction{Code: "disable", Widget: "common"},
	}
}

func (s commonDisable) Register() *md.MDAction {
	return s.register
}

func (s commonDisable) Exec(flow *utils.FlowContext) {

}
