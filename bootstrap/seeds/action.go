package seeds

import (
	"github.com/ggoop/mdf/db"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

func seedAction() {
	items := make([]md.MDActionRule, 0)
	//widget common
	items = append(items, md.MDActionRule{OwnerType: "widget", OwnerCode: "common", Action: "delete", Code: "delete", Name: "删除", Sequence: 50})
	items = append(items, md.MDActionRule{OwnerType: "widget", OwnerCode: "common", Action: "disable", Code: "disable", Name: "停用", Sequence: 50})
	items = append(items, md.MDActionRule{OwnerType: "widget", OwnerCode: "common", Action: "enable", Code: "enable", Name: "启用", Sequence: 50})
	items = append(items, md.MDActionRule{OwnerType: "widget", OwnerCode: "common", Action: "import", Code: "import", Name: "导入", Sequence: 50})
	items = append(items, md.MDActionRule{OwnerType: "widget", OwnerCode: "common", Action: "query", Code: "query", Name: "查询", Sequence: 50})
	items = append(items, md.MDActionRule{OwnerType: "widget", OwnerCode: "common", Action: "save", Code: "save", Name: "保存", Sequence: 50})
	items = append(items, md.MDActionRule{OwnerType: "widget", OwnerCode: "common", Action: "fetchMeta", Code: "fetchMeta", Name: "获取元数据", Sequence: 50})

	//md
	items = append(items, md.MDActionRule{OwnerType: "widget", OwnerCode: "md", Action: "import", Code: "importPre", Name: "保存前规则", Sequence: 30})

	for i, _ := range items {
		item := items[i]
		if item.Domain == "" {
			item.Domain = "mdf"
		}
		item.Enabled = utils.SBool_True
		item.Async = utils.SBool_False
		count := 0
		if err := db.Default().Model(md.MDActionRule{}).Where("owner_type=? and owner_code=? and code=?", item.OwnerType, item.OwnerCode, item.Code).Count(&count).Error; err != nil {
			glog.Error(err)
		} else if count == 0 {
			db.Default().Create(&item)
		}
	}
}
