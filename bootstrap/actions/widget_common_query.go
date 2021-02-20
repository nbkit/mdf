package actions

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type commonQuery struct {
}

func newCommonQuery() *commonQuery {
	return &commonQuery{}
}
func (s *commonQuery) Register() md.RuleRegister {
	return md.RuleRegister{Code: "query", OwnerType: md.RuleType_Widget, OwnerCode: "common"}
}

func (s commonQuery) Exec(token *utils.TokenContext, req *utils.ReqContext, res *utils.ResContext) {

}
