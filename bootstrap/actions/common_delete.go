package actions

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type CommonDelete struct {
}

func newCommonDelete() *CommonDelete {
	return &CommonDelete{}
}
func (s CommonDelete) Register() md.RuleRegister {
	return md.RuleRegister{Code: "delete", Widget: "common"}
}

func (s CommonDelete) Exec(flow *utils.FlowContext) {

}
