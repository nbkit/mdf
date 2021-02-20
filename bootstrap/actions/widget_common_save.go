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
	return md.RuleRegister{Code: "save", OwnerType: md.RuleType_Widget, OwnerCode: "common"}
}

func (s *commonSave) Exec(token *utils.TokenContext, req *utils.ReqContext, res *utils.ResContext) {

}
