package actions

import (
	"github.com/nbkit/mdf/framework/rule"
)

func Register() {
	//注册到mof框架
	rule.ActionSv().RegisterAction(
		newCommonDelete(),
		newCommonDisable(),
		newCommonEnable(),
		newCommonImport(),
		newCommonQuery(),
		newCommonSave(),
	)
}
