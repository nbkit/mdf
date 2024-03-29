package rules

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type commonFetchMeta struct {
	register *md.MDRule
}

func newCommonFetchMeta() *commonFetchMeta {
	return &commonFetchMeta{
		register: &md.MDRule{Action: "fetchMeta", Code: "fetchMeta", Widget: "common", Sequence: 50},
	}
}
func (s *commonFetchMeta) Register() *md.MDRule {
	return s.register
}
func (s *commonFetchMeta) Exec(flow *utils.FlowContext) {

	flow.Set("aaa", flow.Request.Action)
}
