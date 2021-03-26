package rules

import (
	"fmt"
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/utils"

	"github.com/nbkit/mdf/framework/md"
)

type commonEnable struct {
	register *md.MDRule
}

func newCommonEnable() *commonEnable {
	return &commonEnable{
		register: &md.MDRule{Action: "enable", Code: "enable", Widget: "common", Sequence: 50},
	}
}

func (s commonEnable) Register() *md.MDRule {
	return s.register
}

func (s commonEnable) Exec(flow *utils.FlowContext) {
	if flow.Request.ID == "" {
		flow.Error("缺少 ID 参数！")
		return
	}
	//查找实体信息
	entity := md.MDSv().GetEntity(flow.Request.Entity)
	if entity == nil {
		flow.Error("找不到实体！")
		return
	}
	if err := db.Default().Exec(fmt.Sprintf("update %s set enabled=1 where id =?", entity.TableName), flow.Request.ID).Error; err != nil {
		flow.Error(err)
		return
	}
}
