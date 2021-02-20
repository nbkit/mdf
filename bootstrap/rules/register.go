package rules

import (
	"github.com/nbkit/mdf/framework/md"
)

func Register() {
	//注册到mof框架
	md.ActionSv().RegisterRule(
		newCommonQuery(),
		newCommonSave(),
		newCommonDelete(),
		newCommonEnable(),
		newCommonDisable(),
		newCommonImport(),
		newCommonFetchMeta(),

		//md
		newMdImportPre(),
	)
}
