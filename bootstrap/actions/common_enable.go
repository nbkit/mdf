package actions

import (
	"github.com/nbkit/mdf/utils"

	"github.com/nbkit/mdf/framework/md"
)

type commonEnable struct {
}

func newCommonEnable() *commonEnable {
	return &commonEnable{}
}

func (s commonEnable) Register() md.RuleRegister {
	return md.RuleRegister{Code: "enable", Widget: "common"}
}

func (s commonEnable) Exec(flow *utils.FlowContext) {

}
