package actions

import (
	"github.com/nbkit/mdf/framework/widget"
	"github.com/nbkit/mdf/utils"
)

type commonQuery struct {
	register *widget.MDAction
}

func newCommonQuery() *commonQuery {
	return &commonQuery{
		register: &widget.MDAction{Code: "query", Widget: "common", Action: "query"},
	}
}
func (s *commonQuery) Register() *widget.MDAction {
	return s.register
}

func (s commonQuery) Exec(flow *utils.FlowContext) {

}
