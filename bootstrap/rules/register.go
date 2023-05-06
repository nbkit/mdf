package rules

import (
	"github.com/nbkit/mdf/framework/widget"
)

func Register() {
	//注册到mof框架
	widget.ActionSv().RegisterRule(
		newCommonQuery(),
		newCommonSave(),
		newCommonDelete(),
		newCommonEnable(),
		newCommonDisable(),
		newCommonImport(),
		newCommonFetchMeta(),

		//md
		newEntityImportBefore(),
		newUiImportBefore(),
	)
}
