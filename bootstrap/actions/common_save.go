package actions

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type commonSave struct {
}

func newCommonSave() *commonSave {
	return &commonSave{}
}

func (s *commonSave) Register() md.RuleRegister {
	return md.RuleRegister{Code: "save", Widget: "common"}
}

func (s *commonSave) Exec(flow *utils.FlowContext) {

}
