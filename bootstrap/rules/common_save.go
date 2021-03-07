package rules

import (
	"fmt"
	"github.com/nbkit/mdf/db"
	"strings"

	"github.com/nbkit/mdf/framework/glog"
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type commonSave struct {
}

func newCommonSave() *commonSave {
	return &commonSave{}
}

func (s *commonSave) Register() md.RuleRegister {
	return md.RuleRegister{Code: "save", Widget: "common"}
}

func (s *commonSave) Exec(flow *utils.FlowContext) {
	reqData := make(map[string]interface{})
	if data, ok := flow.Request.Data.(map[string]interface{}); !ok {
		flow.Error("data数据格式不正确")
		return
	} else {
		reqData = data
	}
	if reqData == nil {
		flow.Error("没有要保存的数据！")
		return
	}
	//查找实体信息
	entity := md.MDSv().GetEntity(flow.Request.Entity)
	if entity == nil {
		flow.Error("找不到实体！")
		return
	}
	state := ""
	if s, ok := reqData[utils.STATE_FIELD]; ok && s != nil {
		state = s.(string)
	}
	if state == utils.STATE_UPDATED && flow.Request.ID == "" {

	}
	//如果有ID，则为修改保存
	if flow.Request.ID != "" {
		oldData := make(map[string]interface{})
		exector := md.NewOQL().From(entity.TableName)
		for _, f := range entity.Fields {
			if f.TypeType == utils.TYPE_SIMPLE {
				exector.Select(f.Code)
			}
		}
		exector.Where("id=?", flow.Request.ID)
		datas := make([]map[string]interface{}, 0)
		if err := exector.Find(&datas).Error(); err != nil {
			flow.Error(err)
			return
		} else if len(datas) > 0 {
			oldData = datas[0]
		}
		if len(oldData) == 0 {
			flow.Error("找不到要修改的数据！")
			return
		}
		s.doActionUpdate(flow, entity, reqData, oldData)
	} else {
		s.doActionCreate(flow, entity, reqData)
	}
}

func (s *commonSave) fillEntityDefaultValue(entity *md.MDEntity, data map[string]interface{}) map[string]interface{} {
	for _, field := range entity.Fields {
		//如果字段设置了默认值，且没有传入字段值时，取默认值
		if field.DbName != "" && field.DefaultValue != "" {
			if dv, ok := data[field.DbName]; !ok || dv == nil {
				data[field.DbName] = field.CompileValue(field.DefaultValue)
			}
		}
	}
	return data
}
func (s *commonSave) dataToEntityData(entity *md.MDEntity, data map[string]interface{}) map[string]interface{} {
	changeData := make(map[string]interface{})
	for di, dv := range data {
		field := entity.GetField(di)
		if field == nil || field.TypeType != utils.TYPE_SIMPLE {
			continue
		}
		dbFieldName := field.DbName
		if obj, is := dv.(map[string]interface{}); is && obj != nil && obj["code"] != nil {
			changeData[dbFieldName] = obj["code"]
		} else {
			changeData[dbFieldName] = field.CompileValue(dv)
		}
	}
	// 处理枚举和实体
	for di, dv := range data {
		field := entity.GetField(di)
		if field == nil {
			continue
		}
		if field.TypeType == utils.TYPE_ENTITY || field.TypeType == utils.TYPE_ENUM {
			dbField := entity.GetField(field.Code + "ID")
			if dbField == nil || dbField.TypeType != utils.TYPE_SIMPLE || dbField.DbName == "" {
				continue
			}
			if obj, is := dv.(map[string]interface{}); is && obj != nil && obj["id"] != nil {
				changeData[dbField.DbName] = obj["id"]
			} else {
				changeData[dbField.DbName] = ""
			}
			continue
		}
	}
	return changeData
}
func (s *commonSave) doActionCreate(flow *utils.FlowContext, entity *md.MDEntity, reqData map[string]interface{}) {
	reqData["id"] = utils.GUID()
	if sysField := entity.GetField("EntID"); sysField != nil && flow.EntID() != "" {
		reqData[sysField.DbName] = flow.EntID()
	}
	if sysField := entity.GetField("CreatedBy"); sysField != nil && flow.UserID() != "" {
		reqData[sysField.DbName] = flow.UserID()
	}
	fieldCreated := entity.GetField("CreatedAt")
	if fieldCreated != nil && fieldCreated.DbName != "" {
		reqData[fieldCreated.DbName] = utils.TimeNow()
	}
	fieldUpdatedAt := entity.GetField("UpdatedAt")
	if fieldUpdatedAt != nil && fieldUpdatedAt.DbName != "" {
		reqData[fieldUpdatedAt.DbName] = utils.TimeNow()
	}
	//取传入的值
	changeData := s.dataToEntityData(entity, reqData)
	//配置默认值
	changeData = s.fillEntityDefaultValue(entity, changeData)
	if len(changeData) == 0 {
		flow.Error("没有要保存的数据")
		return
	}
	//数据校验
	if err := s.dataCheck(flow, entity, changeData); err != nil {
		flow.Error(err)
		return
	}
	//树规则
	isTree := false
	fieldParent := entity.GetField("ParentID")
	if fieldParent != nil && fieldParent.DbName != "" {
		isTree = true
	}
	isLeafField := entity.GetField("IsLeaf")
	if isTree {
		if changeData["id"] != "" && changeData["id"] == changeData[fieldParent.DbName] {
			flow.Error("树结构，父节点不能等于当前节点!")
			return
		}
		if isLeafField != nil && isLeafField.DbName != "" {
			changeData[isLeafField.DbName] = 1
		}
	}

	//开始保存数据
	fields := make([]string, 0)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	for f, v := range changeData {
		if vv, is := v.(utils.SBool); is && !vv.Valid() {
			continue
		}
		if vv, is := v.(utils.SJson); is && !vv.Valid() {
			continue
		}
		fields = append(fields, db.Default().Dialect().Quote(f))
		placeholders = append(placeholders, "?")
		values = append(values, v)
	}
	sql := fmt.Sprintf("insert into %s (%s) values (%s)", db.Default().Dialect().Quote(entity.TableName), strings.Join(fields, ","), strings.Join(placeholders, ","))
	if err := db.Default().Table(entity.TableName).Exec(sql, values...).Error; err != nil {
		flow.Error(err)
		return
	}
	//处理树节点标识
	if isTree {
		//更新父节点标识
		if parentID, ok := changeData[fieldParent.DbName].(string); ok && parentID != "" && isLeafField != nil {
			updates := make(map[string]interface{})
			updates[isLeafField.DbName] = 0
			if fieldUpdatedAt != nil && fieldUpdatedAt.DbName != "" {
				updates[fieldUpdatedAt.DbName] = utils.TimeNow()
			}
			if err := db.Default().Table(entity.TableName).Where("id=?", parentID).Updates(updates).Error; err != nil {
				flow.Error(err)
				return
			}
		}
	}
	//保存关联实体
	if s.saveRelationData(flow, entity, reqData); flow.Error() != nil {
		return
	}
	flow.Set("data", changeData)
}
func (s *commonSave) doActionUpdate(flow *utils.FlowContext, entity *md.MDEntity, reqData map[string]interface{}, oldData map[string]interface{}) {
	fieldUpdatedAt := entity.GetField("UpdatedAt")
	if fieldUpdatedAt != nil && fieldUpdatedAt.DbName != "" {
		reqData[fieldUpdatedAt.DbName] = utils.TimeNow()
	}
	if sysField := entity.GetField("ID"); sysField != nil && flow.Request.ID != "" {
		reqData[sysField.DbName] = flow.Request.ID
	}
	data := s.dataToEntityData(entity, reqData)

	changeData := make(map[string]interface{})
	for nk, nv := range data {
		if nk == "id" {
			continue
		}
		isChanged := true
		oldValue := oldData[nk]
		field := entity.GetField(nk)
		if field == nil {
			continue
		}
		fieldType := strings.ToLower(field.TypeID)
		//布尔类型判断
		if fieldType == "bool" || fieldType == "boolean" {
			newV := utils.ToSBool(nv)
			oldV := utils.ToSBool(oldValue)
			if !newV.Valid() || newV.Equal(oldV) {
				isChanged = false
			}
		} else {
			if nv == oldValue {
				isChanged = false
			}
		}
		if isChanged {
			changeData[nk] = nv
		}
	}
	//树规则
	isTree := false
	fieldParent := entity.GetField("ParentID")
	if fieldParent != nil && fieldParent.DbName != "" {
		isTree = true
	}
	isLeafField := entity.GetField("IsLeaf")
	if isTree {
		if flow.Request.ID == changeData[fieldParent.DbName] {
			flow.Error("树结构，父节点不能等于当前节点!")
			return
		}
	}
	//数据校验
	if err := s.dataCheck(flow, entity, changeData); err != nil {
		flow.Error(err)
		return
	}
	if len(changeData) > 0 {
		//开始保存数据
		if err := db.Default().Table(entity.TableName).Where("id=?", flow.Request.ID).Updates(changeData).Error; err != nil {
			flow.Error(err)
			return
		}
	}
	//保存关联实体
	if s.saveRelationData(flow, entity, reqData); flow.Error() != nil {
		return
	}

	if len(changeData) > 0 {
		if isTree && isLeafField != nil {
			oldParentID := ""
			if mv, ok := oldData[fieldParent.DbName]; ok {
				oldParentID = mv.(string)
			}
			//如果修改了父节点
			if newParentID, ok := changeData[fieldParent.DbName]; ok {
				if newParentID != "" { //如果设置父节点不为空，则更新父节点为非叶子节点
					updates := make(map[string]interface{})
					updates[isLeafField.DbName] = utils.SBool_False
					if fieldUpdatedAt != nil && fieldUpdatedAt.DbName != "" {
						updates[fieldUpdatedAt.DbName] = utils.TimeNow()
					}
					if err := db.Default().Table(entity.TableName).Where("id=?", newParentID).Updates(updates).Error; err != nil {
						flow.Error(err)
						return
					}
				}
				if oldParentID != "" { //如果设置父节点为空，则更新父节点叶子节点状态
					count := 0
					updates := make(map[string]interface{})
					if fieldUpdatedAt != nil && fieldUpdatedAt.DbName != "" {
						updates[fieldUpdatedAt.DbName] = utils.TimeNow()
					}
					if db.Default().Table(entity.TableName).Where(fmt.Sprintf("%s=?", fieldParent.DbName), oldParentID).Count(&count); count == 0 {
						updates[isLeafField.DbName] = 1
					} else {
						updates[isLeafField.DbName] = 0
					}
					if err := db.Default().Table(entity.TableName).Where("id=?", oldParentID).Updates(updates).Error; err != nil {
						flow.Error(err)
						return
					}
				}
			}
		}
		flow.Set("data", changeData)
	}
	return
}

func (s *commonSave) saveRelationData(flow *utils.FlowContext, entity *md.MDEntity, reqData map[string]interface{}) {
	for _, nv := range entity.Fields {
		if nv.Kind == md.KIND_TYPE_HAS_MANT {
			if do, ok := reqData[nv.DbName].([]interface{}); ok && len(do) > 0 {
				for _, dr := range do {
					if ds, ok := dr.(map[string]interface{}); ok {
						state := ""

						if s, ok := ds[utils.STATE_FIELD]; ok && s != nil {
							state = s.(string)
						}
						if state == "" {
							glog.Error("实体对应状态为空，跳过更新！", glog.String("state", state))
							continue
						}
						newFlow := flow.Copy()
						refEntity := md.MDSv().GetEntity(nv.TypeID)
						if f := refEntity.GetField(nv.ForeignKey); f != nil {
							ds[f.DbName] = reqData["id"]
						}
						ruleID := ""
						if state == utils.STATE_CREATED || state == utils.STATE_UPDATED {
							ruleID = "save"
						}
						if state == utils.STATE_DELETED {
							ruleID = "delete"
						}
						if ruleID == "" {
							glog.Error("该状态找不到对应规则", glog.String("state", state))
							continue
						}
						if state == utils.STATE_UPDATED || state == utils.STATE_DELETED {
							if id, ok := ds["id"].(string); ok && id != "" {
								newFlow.Request.ID = id
							}
						}
						newFlow.Request.Data = ds
						newFlow.Request.Entity = refEntity.ID
						newFlow.Request.Rule = ruleID

						if md.ActionSv().DoAction(newFlow); newFlow.Error() != nil {
							flow.Error(newFlow.Error())
							return
						}
					}
				}

			}
		}
	}
}
func (s *commonSave) dataCheck(flow *utils.FlowContext, entity *md.MDEntity, data map[string]interface{}) error {
	return nil
}
