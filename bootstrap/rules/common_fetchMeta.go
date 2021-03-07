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
	return md.RuleRegister{Code: "fetchMeta", Widget: "common"}
}
func (s *commonFetchMeta) Exec(flow *utils.FlowContext) {

	flow.Set("aaa", flow.Request.Action)
}
