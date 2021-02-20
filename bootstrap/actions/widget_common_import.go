package actions

import (
	"github.com/nbkit/mdf/framework/files"
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type commonImport struct {
}

func newCommonImport() *commonImport {
	return &commonImport{}
}
func (s *commonImport) Register() md.RuleRegister {
	return md.RuleRegister{Code: "import", OwnerType: md.RuleType_Widget, OwnerCode: "common"}
}
func (s *commonImport) Exec(token *utils.TokenContext, req *utils.ReqContext, res *utils.ResContext) {
	if len(req.Files) > 0 {
		datas := make([]files.ImportData, 0)
		for _, file := range req.Files {
			if f, err := file.Open(); err != nil {
				res.SetError(err)
				return
			} else {
				if ds, err := files.NewExcelSv().GetExcelDatasByReader(f); err != nil {
					res.SetError(err)
					return
				} else if len(ds) > 0 {
					datas = append(datas, ds...)
				}
			}
		}
		req.Data = datas
	}
}
