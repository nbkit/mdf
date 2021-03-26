package actions

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type commonQuery struct {
	register *md.MDAction
}

func newCommonQuery() *commonQuery {
	return &commonQuery{
		register: &md.MDAction{Code: "query", Widget: "common", Action: "query"},
	}
}
func (s *commonQuery) Register() *md.MDAction {
	return s.register
}

func (s commonQuery) Exec(flow *utils.FlowContext) {

}
