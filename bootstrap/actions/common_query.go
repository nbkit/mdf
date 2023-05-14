package actions

import (
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/utils"
)

type commonQuery struct {
}

func newCommonQuery() commonQuery {
	return commonQuery{}
}
func (s commonQuery) Register() rule.MDAction {
	return rule.MDAction{Code: "query", Widget: "common", Action: "query"}
}

func (s commonQuery) Exec(flow *utils.FlowContext) {

}
