package rules

import (
	"fmt"
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/utils"

	"github.com/nbkit/mdf/framework/md"
)

type commonEnable struct {
}

func newCommonEnable() *commonEnable {
	return &commonEnable{}
}

func (s commonEnable) Register() md.RuleRegister {
	return md.RuleRegister{Code: "enable", OwnerType: md.RuleType_Widget, OwnerCode: "common"}
}

func (s commonEnable) Exec(token *utils.TokenContext, req *utils.ReqContext, res *utils.ResContext) {
	if req.ID == "" {
		res.SetError("缺少 ID 参数！")
		return
	}
	//查找实体信息
	entity := md.MDSv().GetEntity(req.Entity)
	if entity == nil {
		res.SetError("找不到实体！")
		return
	}
	if err := db.Default().Exec(fmt.Sprintf("update %s set enabled=1 where id =?", entity.TableName), req.ID).Error; err != nil {
		res.SetError(err)
		return
	}
}
