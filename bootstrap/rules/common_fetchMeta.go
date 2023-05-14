package rules

import (
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/utils"
)

type commonFetchMeta struct {
}

func newCommonFetchMeta() commonFetchMeta {
	return commonFetchMeta{}
}
func (s commonFetchMeta) Register() rule.MDRule {
	return rule.MDRule{Action: "fetchMeta", Widget: "common", Sequence: 50}
}
func (s commonFetchMeta) Exec(flow *utils.FlowContext) {

	flow.Set("aaa", flow.Request.Action)
}
