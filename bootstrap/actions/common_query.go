package actions

import (
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/utils"
)

type commonQuery struct {
	register *rule.MDAction
}

func newCommonQuery() *commonQuery {
	return &commonQuery{
		register: &rule.MDAction{Code: "query", Widget: "common", Action: "query"},
	}
}
func (s *commonQuery) Register() *rule.MDAction {
	return s.register
}

func (s commonQuery) Exec(flow *utils.FlowContext) {

}
