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
	return md.RuleRegister{Code: "query", Widget: "common"}
}

func (s commonQuery) Exec(flow *utils.FlowContext) {

}
