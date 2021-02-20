package actions

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
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
