package actions

import (
	"github.com/nbkit/mdf/framework/files"
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/utils"
)

type commonImport struct {
}

func newCommonImport() commonImport {
	return commonImport{}
}
func (s commonImport) Register() rule.MDAction {
	return rule.MDAction{Code: "import", Widget: "common", Action: "import"}
}

func (s commonImport) Exec(flow *utils.FlowContext) {
	if len(flow.Request.Files) > 0 {
		datas := make([]files.ImportData, 0)
		for _, file := range flow.Request.Files {
			if f, err := file.Open(); err != nil {
				flow.Error(err)
				return
			} else {
				if ds, err := files.NewExcelSv().GetExcelDatasByReader(f); err != nil {
					flow.Error(err)
					return
				} else if len(ds) > 0 {
					datas = append(datas, ds...)
				}
			}
		}
		flow.Request.Data = datas
	}
}
