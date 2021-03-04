package rules

import (
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/framework/files"
	"github.com/nbkit/mdf/framework/glog"
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type uiImportPre struct {
}

func newUiImportPre() *uiImportPre {
	return &uiImportPre{}
}
func (s *uiImportPre) Register() md.RuleRegister {
	return md.RuleRegister{Code: "importPre", OwnerType: md.RuleType_Widget, OwnerCode: "md"}
}
func (s *uiImportPre) Exec(flow *utils.FlowContext) {
	if flow.Request.Data == nil {
		flow.Error("没有要导入的数据")
		return
	}
	if items, ok := flow.Request.Data.([]files.ImportData); !ok {
		flow.Error("导入的数据非法！")
		return
	} else {
		s.deleteData(flow, items)
		s.doProcess(flow, items)
	}
}
func (s *uiImportPre) deleteData(flow *utils.FlowContext, data []files.ImportData) {
	widgetCodes := make([]string, 0)
	filterCodes := make([]string, 0)
	for i, _ := range data {
		d := data[i]
		if d.EntityCode == "md.widget" {
			for _, r := range d.Data {
				if cv, co := r["Code"]; co && cv != "" {
					widgetCodes = append(widgetCodes, cv)
				}
			}
		}
		if d.EntityCode == "md.filters" {
			for _, r := range d.Data {
				if cv, co := r["Code"]; co && cv != "" {
					filterCodes = append(filterCodes, cv)
				}
			}
		}
	}
	var sql string
	//先按删除数据
	if len(widgetCodes) > 0 {
		sql = "delete from md_widget_ds where widget_id in (select id from md_widgets where code in (?))"
		if err := db.Default().Exec(sql, widgetCodes).Error; err != nil {
			glog.Error(flow.Error(err))
			return
		}

		sql = "delete from md_widget_layouts where widget_id in (select id from md_widgets where code in (?))"
		if err := db.Default().Exec(sql, widgetCodes).Error; err != nil {
			glog.Error(flow.Error(err))
			return
		}

		sql = "delete from md_widget_items where widget_id in (select id from md_widgets where code in (?))"
		if err := db.Default().Exec(sql, widgetCodes).Error; err != nil {
			glog.Error(flow.Error(err))
			return
		}
		sql = "delete from md_toolbars where widget_id in (select id from md_widgets where code in (?))"
		if err := db.Default().Exec(sql, widgetCodes).Error; err != nil {
			glog.Error(flow.Error(err))
			return
		}
		sql = "delete from md_toolbar_items where widget_id in (select id from md_widgets where code in (?))"
		if err := db.Default().Exec(sql, widgetCodes).Error; err != nil {
			glog.Error(flow.Error(err))
			return
		}

		sql = "delete from md_action_commands where owner_type ='widget' and owner_code in (?)"
		if err := db.Default().Exec(sql, widgetCodes).Error; err != nil {
			glog.Error(flow.Error(err))
			return
		}
		sql = "delete from md_action_rules where owner_type ='widget' and owner_code in (?)"
		if err := db.Default().Exec(sql, widgetCodes).Error; err != nil {
			glog.Error(flow.Error(err))
			return
		}

		sql = "delete from auth_permits where owner_type ='widget' and owner_code in (?)"
		if err := db.Default().Exec(sql, widgetCodes).Error; err != nil {
			glog.Error(flow.Error(err))
			return
		}

		sql = "delete from md_widgets where code in (?)"
		if err := db.Default().Exec(sql, widgetCodes).Error; err != nil {
			glog.Error(flow.Error(err))
			return
		}
	}
	if len(filterCodes) > 0 {
		sql = "delete from md_filter_items where filter_id in (select id from md_filters where code in (?))"
		if err := db.Default().Exec(sql, filterCodes).Error; err != nil {
			glog.Error(flow.Error(err))
			return
		}
		sql = "delete from md_filter_solutions where filter_id in (select id from md_filters where code in (?))"
		if err := db.Default().Exec(sql, filterCodes).Error; err != nil {
			glog.Error(flow.Error(err))
			return
		}
		sql = "delete from md_filters where code in (?)"
		if err := db.Default().Exec(sql, filterCodes).Error; err != nil {
			glog.Error(flow.Error(err))
			return
		}
	}

}

func (s *uiImportPre) doProcess(flow *utils.FlowContext, data []files.ImportData) {

}
