package rules

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type commonFetchMeta struct {
}

func newCommonFetchMeta() *commonFetchMeta {
	return &commonFetchMeta{}
}
func (s *commonFetchMeta) Register() md.RuleRegister {
	return md.RuleRegister{Code: "fetchMeta", OwnerType: md.RuleType_Widget, OwnerCode: "common"}
}
func (s *commonFetchMeta) Exec(token *utils.TokenContext, req *utils.ReqContext, res *utils.ResContext) {

	res.Set("aaa", req.Action)
}
