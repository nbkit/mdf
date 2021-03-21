package actions

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type commonSave struct {
	register *md.MDAction
}

func newCommonSave() *commonSave {
	return &commonSave{
		register: &md.MDAction{Code: "save", Widget: "common"},
	}
}

func (s *commonSave) Register() *md.MDAction {
	return s.register
}

func (s *commonSave) Exec(flow *utils.FlowContext) {

}
