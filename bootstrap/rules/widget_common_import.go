package rules

import (
	"fmt"
	"github.com/nbkit/mdf/bootstrap/services"
	"github.com/nbkit/mdf/db"
	"strings"

	"github.com/nbkit/mdf/bootstrap/model"
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
	logData := model.Log{EntID: token.EntID(), UserID: token.UserID(), NodeType: req.OwnerType, NodeID: req.OwnerCode, DataID: req.Entity}
	services.LogSv().CreateLog(logData.Clone().SetMsg("导入开始======begin======"))

	defer func() {
		services.LogSv().CreateLog(logData.Clone().SetMsg("导入结束"))
	}()
	if req.Data == nil {
		err := res.Error("没有要导入的数据")
		services.LogSv().CreateLog(logData.Clone().SetMsg(err))
		return
	}
	data := req.Data
	if items, ok := data.([]files.ImportData); ok {
		for _, data := range items {
			s.importMapData(token, req, res, data)
		}
	} else if items, ok := data.(files.ImportData); ok {
		s.importMapData(token, req, res, items)
	}
}
func (s *commonImport) importMapData(token *utils.TokenContext, req *utils.ReqContext, res *utils.ResContext, data files.ImportData) {
	log := services.LogSv()
	logData := model.Log{EntID: token.EntID(), UserID: token.UserID(), NodeType: req.OwnerType, NodeID: req.OwnerCode, DataID: req.Entity}
	log.CreateLog(logData.Clone().SetMsg(fmt.Sprintf("接收到需要导入的数据-%s：%v条", req.Entity, len(data.Data))))

	entity := md.MDSv().GetEntity(data.EntityCode)
	if entity == nil {
		entity = md.MDSv().GetEntity(req.Entity)
	}
	if entity == nil {
		res.Error("没有配置导入实体")
		return
	}
	dbDatas := make([]map[string]interface{}, 0)
	quotedMap := make(map[string]string)

	for _, item := range data.Data {
		dbItem := make(map[string]interface{})
		if v, ok := item[utils.STATE_FIELD]; ok && (v == utils.STATE_TEMP || v == utils.STATE_NORMAL || v == utils.STATE_IGNORED) {
			continue
		}
		for kk, kv := range item {
			field := entity.GetField(kk)
			if field == nil || kv == "" {
				continue
			}
			fieldName := ""
			if field.TypeType == utils.TYPE_ENTITY {
				fieldName = field.DbName + "_id"
				qreq := utils.ReqContext{Entity: field.TypeID, Q: kv, Data: item}
				if obj, err := md.MDSv().TakeDataByQ(token, &qreq); err != nil {
					log.CreateLog(logData.Clone().SetMsg(fmt.Sprintf("数据[%s]=[%s],查询失败：%v", qreq.Entity, qreq.Q, err.Error())))
				} else if len(obj) > 0 && obj["id"] != nil {
					dbItem[fieldName] = obj["id"]
					quotedMap[fieldName] = fieldName
				} else if len(obj) > 0 && obj["ID"] != nil {
					dbItem[fieldName] = obj["ID"]
					quotedMap[fieldName] = fieldName
				} else {
					log.CreateLog(logData.Clone().SetMsg(fmt.Sprintf("关联对象[%s],找不到[%s]对应数据!", qreq.Entity, qreq.Q)))
				}
			} else if field.TypeType == utils.TYPE_ENUM {
				fieldName = field.DbName + "_id"
				if vv := md.MDSv().GetEnum(field.Limit, kv); vv != nil {
					dbItem[fieldName] = vv.ID
					quotedMap[fieldName] = fieldName
				} else {
					log.CreateLog(logData.Clone().SetMsg(fmt.Sprintf("关联枚举[%s],找不到[%s]对应数据!", field.Limit, kv)))
				}
			} else if field.TypeType == utils.TYPE_SIMPLE {
				fieldName = field.DbName
				if field.TypeID == utils.FIELD_TYPE_BOOL {
					dbItem[fieldName] = utils.ToSBool(files.GetCellValue(kk, item))
					quotedMap[fieldName] = fieldName
				} else if field.TypeID == utils.FIELD_TYPE_DATETIME || field.TypeID == utils.FIELD_TYPE_DATE {
					dbItem[fieldName] = files.GetCellValue(kk, item)
					quotedMap[fieldName] = fieldName
				} else if field.TypeID == utils.FIELD_TYPE_DECIMAL || field.TypeID == utils.FIELD_TYPE_INT {
					dbItem[fieldName] = files.GetCellValue(kk, item)
					quotedMap[fieldName] = fieldName
				} else {
					dbItem[fieldName] = kv
					quotedMap[fieldName] = fieldName
				}
			}
		}
		if field := entity.GetField("ID"); field != nil {
			fieldName := field.DbName
			if _, ok := dbItem[fieldName]; !ok {
				dbItem[fieldName] = utils.GUID()
			}
			quotedMap[fieldName] = fieldName
		}
		if field := entity.GetField("EntID"); field != nil && field.DbName != "" {
			fieldName := field.DbName
			if _, ok := dbItem[fieldName]; !ok {
				dbItem[fieldName] = token.EntID()
			}
			quotedMap[fieldName] = fieldName
		}
		if field := entity.GetField("CreatedBy"); field != nil && field.DbName != "" {
			fieldName := field.DbName
			if _, ok := dbItem[fieldName]; !ok {
				dbItem[fieldName] = token.UserID()
			}
			quotedMap[fieldName] = fieldName
		}
		if field := entity.GetField("CreatedAt"); field != nil && field.DbName != "" {
			fieldName := field.DbName
			dbItem[fieldName] = utils.TimeNow()
			quotedMap[fieldName] = fieldName
		}
		if field := entity.GetField("UpdatedAt"); field != nil && field.DbName != "" {
			fieldName := field.DbName
			dbItem[fieldName] = utils.TimeNow()
			quotedMap[fieldName] = fieldName
		}
		if len(dbItem) > 0 {
			dbDatas = append(dbDatas, dbItem)
		}
	}
	quoted := make([]string, 0, len(quotedMap))

	for fk, _ := range quotedMap {
		quoted = append(quoted, fk)
	}

	placeholdersArr := make([]string, 0, len(quotedMap))
	valueVars := make([]interface{}, 0)
	var itemCount uint = 0
	var MaxBatchs uint = 100

	for _, data := range dbDatas {
		itemCount = itemCount + 1
		placeholders := make([]string, 0, len(quoted))
		for _, f := range quoted {
			placeholders = append(placeholders, "?")
			valueVars = append(valueVars, data[f])
		}
		placeholdersArr = append(placeholdersArr, "("+strings.Join(placeholders, ", ")+")")

		if itemCount >= MaxBatchs {
			if err := s.batchInsertSave(entity, quoted, placeholdersArr, valueVars...); err != nil {
				log.CreateLog(logData.Clone().SetMsg(fmt.Sprintf("数据库存储[%v]条记录出错了:%s!", itemCount, err.Error())))
				res.Error(err)
				return
			}
			itemCount = 0
			placeholdersArr = make([]string, 0, len(quotedMap))
			valueVars = make([]interface{}, 0)
		}
	}
	if itemCount > 0 {
		if err := s.batchInsertSave(entity, quoted, placeholdersArr, valueVars...); err != nil {
			log.CreateLog(logData.Clone().SetMsg(fmt.Sprintf("数据库存储[%v]条记录出错了:%s!", itemCount, err.Error())))
			res.Error(err)
			return
		}
	}
}

func (s *commonImport) batchInsertSave(entity *md.MDEntity, quoted []string, placeholders []string, valueVars ...interface{}) error {
	var sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", db.Default().Dialect().Quote(entity.TableName), strings.Join(quoted, ", "), strings.Join(placeholders, ", "))

	if err := db.Default().Exec(sql, valueVars...).Error; err != nil {
		return err
	}
	return nil
}
