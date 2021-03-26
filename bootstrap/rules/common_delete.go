package rules

import (
	"fmt"
	"github.com/nbkit/mdf/db"

	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type CommonDelete struct {
	register *md.MDRule
}

func newCommonDelete() *CommonDelete {
	return &CommonDelete{
		register: &md.MDRule{Action: "delete", Code: "delete", Widget: "common", Sequence: 50},
	}
}

func (s CommonDelete) Register() *md.MDRule {
	return s.register
}

func (s CommonDelete) Exec(flow *utils.FlowContext) {
	if flow.Request.ID == "" {
		flow.Error("缺少 ID 参数！")
		return
	}
	//查找实体信息
	entity := md.MDSv().GetEntity(flow.Request.Entity)
	if entity == nil {
		flow.Error("找不到实体")
		return
	}
	if field := entity.GetField("System"); field != nil && field.DbName != "" {
		count := 0
		db.Default().Table(entity.TableName).Where(fmt.Sprintf("id=? and %s = 1", db.Default().Dialect().Quote("system")), flow.Request.ID).Count(&count)
		if count > 0 {
			flow.Error("系统预制数据不可删除")
			return
		}
	}
	if df := entity.GetField("DeletedAt"); df != nil {
		if err := db.Default().Exec(fmt.Sprintf("update %s set %s=? where id=?", entity.TableName, df.DbName), utils.TimeNow(), flow.Request.ID).Error; err != nil {
			flow.Error(err)
			return
		}
	} else {
		if err := db.Default().Exec(fmt.Sprintf("delete from %s where id=?", entity.TableName), flow.Request.ID).Error; err != nil {
			flow.Error(err)
			return
		}
	}
}
