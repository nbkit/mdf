package mdf

import (
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
)

func initSeedAction() {
	items := make([]rule.MDRule, 0)
	//widget common
	items = append(items, rule.MDRule{Widget: "common", Action: "delete", Name: "删除", Sequence: 50})
	items = append(items, rule.MDRule{Widget: "common", Action: "disable", Name: "停用", Sequence: 50})
	items = append(items, rule.MDRule{Widget: "common", Action: "enable", Name: "启用", Sequence: 50})
	items = append(items, rule.MDRule{Widget: "common", Action: "import", Name: "导入", Sequence: 50})
	items = append(items, rule.MDRule{Widget: "common", Action: "query", Name: "查询", Sequence: 50})
	items = append(items, rule.MDRule{Widget: "common", Action: "save", Name: "保存", Sequence: 50})
	items = append(items, rule.MDRule{Widget: "common", Action: "fetchMeta", Name: "获取元数据", Sequence: 50})

	//ui
	items = append(items, rule.MDRule{Widget: "ui", Action: "import", Name: "保存前规则", Sequence: 30})

	for i, _ := range items {
		item := items[i]
		if item.Domain == "" {
			item.Domain = "mdf"
		}
		item.Enabled = utils.SBool_True
		count := 0
		if err := db.Default().Model(rule.MDRule{}).Where("widget=? and action=?", item.Widget, item.Action).Count(&count).Error; err != nil {
			log.ErrorD(err)
		} else if count == 0 {
			db.Default().Create(&item)
		}
	}
}
