package rules

import (
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/utils"
)

type commonFetchMeta struct {
	register *rule.MDRule
}

func newCommonFetchMeta() *commonFetchMeta {
	return &commonFetchMeta{
		register: &rule.MDRule{Action: "fetchMeta", Code: "fetchMeta", Widget: "common", Sequence: 50},
	}
}
func (s *commonFetchMeta) Register() *rule.MDRule {
	return s.register
}
func (s *commonFetchMeta) Exec(flow *utils.FlowContext) {

	flow.Set("aaa", flow.Request.Action)
}
