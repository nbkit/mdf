package rules

import (
	"github.com/nbkit/mdf/framework/rule"
)

func Register() {
	//注册到mof框架
	rule.ActionSv().RegisterRule(
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
