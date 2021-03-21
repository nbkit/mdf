package actions

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type CommonDelete struct {
	register *md.MDAction
}

func newCommonDelete() *CommonDelete {
	return &CommonDelete{
		register: &md.MDAction{Code: "delete", Widget: "common"},
	}
}
func (s CommonDelete) Register() *md.MDAction {
	return s.register
}

func (s CommonDelete) Exec(flow *utils.FlowContext) {

}
