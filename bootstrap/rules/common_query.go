package rules

import (
	"fmt"
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
)

type commonQuery struct {
	register *md.MDRule
}

func newCommonQuery() *commonQuery {
	return &commonQuery{
		register: &md.MDRule{Code: "query", Widget: "common"},
	}
}
func (s *commonQuery) Register() *md.MDRule {
	return s.register
}

func (s commonQuery) Exec(flow *utils.FlowContext) {
	if flow.Request.Entity == "" {
		flow.Error("缺少 MainEntity 参数！")
		return
	}
	//查找实体信息
	entity := md.MDSv().GetEntity(flow.Request.Entity)
	if entity == nil {
		flow.Error("找不到实体！")
		return
	}
	exector := md.NewOQL().From(entity.TableName)
	for _, f := range entity.Fields {
		if f.TypeType == utils.TYPE_SIMPLE {
			exector.Select(fmt.Sprintf("$$%s as \"%s\"", f.Code, f.DbName))
		}
	}
	if flow.Request.ID != "" {
		exector.Where("id=?", flow.Request.ID)
	}
	if sysField := entity.GetField("EntID"); sysField != nil && flow.EntID() != "" {
		exector.Where(sysField.Code+" = ?", flow.EntID())
	}
	count := 0
	if err := exector.Count(&count).Error(); err != nil {
		flow.Error(err)
		return
	}
	datas := make([]map[string]interface{}, 0)
	if err := exector.Find(datas).Error(); err != nil {
		flow.Error(err)
		return
	} else if len(datas) > 0 {
		s.loadEnums(datas, entity)
		s.loadEntities(datas, entity)
		flow.Set("data", datas)
		flow.Set("pager", utils.Pager{Total: count, PageSize: flow.Request.PageSize, Page: flow.Request.Page})
	}
}
func (s commonQuery) loadEnums(datas []map[string]interface{}, entity *md.MDEntity) error {
	for _, f := range entity.Fields {
		if f.TypeType == utils.TYPE_ENUM {
			for ri, data := range datas {
				if fv, ok := data[f.DbName+"_id"]; ok && fv != nil && fv.(string) != "" {
					datas[ri][f.DbName] = md.MDSv().GetEnum(f.Limit, fv.(string))
				}
			}
		}
	}
	return nil
}
func (s commonQuery) loadEntities(datas []map[string]interface{}, entity *md.MDEntity) error {
	for _, f := range entity.Fields {
		if f.TypeType == utils.TYPE_ENTITY && f.TypeID != "" && (f.Kind == md.KIND_TYPE_BELONGS_TO || f.Kind == md.KIND_TYPE_HAS_ONE) {
			ids := make([]interface{}, 0)
			for _, data := range datas {
				if fv, ok := data[f.DbName+"_id"]; ok && fv != nil && fv.(string) != "" {
					ids = append(ids, fv)
				}
			}
			if len(ids) > 0 {
				refEntity := md.MDSv().GetEntity(f.TypeID)
				if refEntity != nil {
					exector := md.NewOQL().From(refEntity.TableName)
					for _, f := range refEntity.Fields {
						if f.TypeType == utils.TYPE_SIMPLE {
							exector.Select(fmt.Sprintf("$$%s as \"%s\"", f.Code, f.DbName))
						}
					}
					exector.Where(fmt.Sprintf("%s in ( ? )", f.AssociationKey), ids)
					refDatas := make([]map[string]interface{}, 0)
					if err := exector.Find(&refDatas).Error(); err != nil {
						log.ErrorD(err)
					} else if len(refDatas) > 0 {
						dataMap := make(map[string]interface{})
						for i, _ := range refDatas {
							d := refDatas[i]
							dataMap[d["id"].(string)] = d
						}
						for i, data := range datas {
							if fv, ok := data[f.DbName+"_id"]; ok && fv != nil && fv.(string) != "" {
								datas[i][f.DbName] = dataMap[fv.(string)]
							}
						}
					}
				}
			}
		}
	}
	return nil
}
