package actions

import (
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/utils"
)

type commonEnable struct {
}

func newCommonEnable() commonEnable {
	return commonEnable{}
}

func (s commonEnable) Register() rule.MDAction {
	return rule.MDAction{Code: "enable", Widget: "common", Action: "enable"}
}

func (s commonEnable) Exec(flow *utils.FlowContext) {

}
