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
	return md.RuleRegister{Code: "delete", OwnerType: md.RuleType_Widget, OwnerCode: "common"}
}

func (s CommonDelete) Exec(token *utils.TokenContext, req *utils.ReqContext, res *utils.ResContext) {

}
