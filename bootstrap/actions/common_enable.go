package actions

import (
	"github.com/nbkit/mdf/utils"

	"github.com/nbkit/mdf/framework/md"
)

type commonEnable struct {
	register *md.MDAction
}

func newCommonEnable() *commonEnable {
	return &commonEnable{
		register: &md.MDAction{Code: "enable", Widget: "common", Action: "enable"},
	}
}

func (s commonEnable) Register() *md.MDAction {
	return s.register
}

func (s commonEnable) Exec(flow *utils.FlowContext) {

}
