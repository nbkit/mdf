package actions

import (
	"github.com/nbkit/mdf/framework/widget"
)

func Register() {
	//注册到mof框架
	widget.ActionSv().RegisterAction(
		newCommonDelete(),
		newCommonDisable(),
		newCommonEnable(),
		newCommonImport(),
		newCommonQuery(),
		newCommonSave(),
	)
}
