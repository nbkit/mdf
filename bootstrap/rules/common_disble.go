package rules

import (
	"fmt"
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/framework/rule"
	"github.com/nbkit/mdf/utils"

	"github.com/nbkit/mdf/framework/md"
)

type commonDisable struct {
}

func newCommonDisable() commonDisable {
	return commonDisable{}
}

func (s commonDisable) Register() rule.MDRule {
	//return rule.MDRule{Action: "disable", Widget: "common", Sequence: 50}
	return rule.MDRule{Action: "*", Widget: "aa", Sequence: 50}
}
func (s commonDisable) query(flow *utils.FlowContext) {
	if flow.Request.ID == "" {
		flow.Error("缺少 ID 参数！")
		return
	}
}
func (s commonDisable) Exec(flow *utils.FlowContext) {
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
	if err := db.Default().Exec(fmt.Sprintf("update %s set enabled=0 where id =?", entity.TableName), flow.Request.ID).Error; err != nil {
		flow.Error(err)
		return
	}
}
