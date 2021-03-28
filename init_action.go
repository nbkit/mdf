package mdf

import (
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
)

func initSeedAction() {
	items := make([]md.MDRule, 0)
	//widget common
	items = append(items, md.MDRule{Widget: "common", Action: "delete", Code: "delete", Name: "删除", Sequence: 50})
	items = append(items, md.MDRule{Widget: "common", Action: "disable", Code: "disable", Name: "停用", Sequence: 50})
	items = append(items, md.MDRule{Widget: "common", Action: "enable", Code: "enable", Name: "启用", Sequence: 50})
	items = append(items, md.MDRule{Widget: "common", Action: "import", Code: "import", Name: "导入", Sequence: 50})
	items = append(items, md.MDRule{Widget: "common", Action: "query", Code: "query", Name: "查询", Sequence: 50})
	items = append(items, md.MDRule{Widget: "common", Action: "save", Code: "save", Name: "保存", Sequence: 50})
	items = append(items, md.MDRule{Widget: "common", Action: "fetchMeta", Code: "fetchMeta", Name: "获取元数据", Sequence: 50})

	//ui
	items = append(items, md.MDRule{Widget: "ui", Action: "import", Code: "import.before", Name: "保存前规则", Sequence: 30})

	for i, _ := range items {
		item := items[i]
		if item.Domain == "" {
			item.Domain = "mdf"
		}
		item.Enabled = utils.SBool_True
		item.Async = utils.SBool_False
		count := 0
		if err := db.Default().Model(md.MDRule{}).Where("widget=? and code=?", item.Widget, item.Code).Count(&count).Error; err != nil {
			log.ErrorD(err)
		} else if count == 0 {
			db.Default().Create(&item)
		}
	}
}
