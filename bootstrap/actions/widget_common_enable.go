package actions

import (
	"github.com/ggoop/mdf/utils"

	"github.com/ggoop/mdf/framework/md"
)

type commonEnable struct {
}

func newCommonEnable() *commonEnable {
	return &commonEnable{}
}

func (s commonEnable) Register() md.RuleRegister {
	return md.RuleRegister{Code: "enable", OwnerType: md.RuleType_Widget, OwnerCode: "common"}
}

func (s commonEnable) Exec(token *utils.TokenContext, req *utils.ReqContext, res *utils.ResContext) {

}
