package actions

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type commonDisable struct {
}

func newCommonDisable() *commonDisable {
	return &commonDisable{}
}

func (s commonDisable) Register() md.RuleRegister {
	return md.RuleRegister{Code: "disable", Widget: "common"}
}
func (s commonDisable) Exec(flow *utils.FlowContext) {

}
