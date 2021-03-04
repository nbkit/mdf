package md

import (
	"fmt"
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/utils"
	"strings"
)

func (s *mdSvImpl) GetEnums() ([]MDEnum, error) {
	items := make([]MDEnum, 0)
	if err := db.Default().Model(&MDEnum{}).Where("entity_id in (?)", db.Default().Model(MDEntity{}).Select("id").Where("type=?", "enum").SubQuery()).Order("entity_id").Order("sequence").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
func (s *mdSvImpl) GetEntities() ([]MDEntity, error) {
	items := make([]MDEntity, 0)
	if err := db.Default().Model(&MDEntity{}).Preload("Fields").Order("id").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *mdSvImpl) GetEnum(typeId string, values ...string) *MDEnum {
	if s.enumCache == nil || typeId == "" || values == nil || len(values) == 0 {
		return nil
	}
	for _, v := range values {
		if v, ok := s.enumCache[strings.ToLower(typeId+":"+v)]; ok {
			return v
		}
	}
	return nil
}
func (s *mdSvImpl) GetEnumBy(typeId string) ([]MDEnum, error) {
	items := make([]MDEnum, 0)
	if err := db.Default().Model(&MDEnum{}).Where("entity_id=?", typeId).Order("sequence,id").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
func (s *mdSvImpl) GetEntity(id string) *MDEntity {
	if v, ok := s.entityCache[strings.ToLower(id)]; ok {
		return v
	}
	item := MDEntity{}
	db.Default().Preload("Fields").Order("id").Take(&item, "id=?", id)
	if item.ID != "" {
		s.cacheEntity(item)
		return &item
	}
	return nil
}

func (s *mdSvImpl) TakeDataByQ(flow *utils.FlowContext) {
	entity := s.GetEntity(flow.Request.Entity)
	if entity == nil {
		return
	}
	oql := NewOQL()
	oql.From(entity.TableName)
	codeField := &MDField{}
	nameField := &MDField{}
	entField := &MDField{}
	for i, f := range entity.Fields {
		if f.TypeType == utils.TYPE_SIMPLE {
			oql.Select("$$" + f.Code + " as \"" + f.DbName + "\"")
			if strings.Contains(f.Tags, "code") {
				codeField = &entity.Fields[i]
			}
			if strings.Contains(f.Tags, "name") {
				nameField = &entity.Fields[i]
			}
			if strings.Contains(f.Tags, "ent") {
				entField = &entity.Fields[i]
			}
		}
	}
	if entField == nil || entField.Code == "" {
		entField = entity.GetField("EntID")
	}
	if entField != nil && entField.Code != "" {
		oql.Where(fmt.Sprintf("%s=?", entField.DbName), flow.EntID())
	}

	if codeField == nil || codeField.Code == "" {
		codeField = entity.GetField("Code")
	}
	qwhere := oql.Or()
	if codeField != nil && codeField.Code != "" {
		qwhere.Where(codeField.Code+" = ?", flow.Request.Q)
	}
	if nameField == nil || nameField.Code == "" {
		nameField = entity.GetField("Name")
	}
	if nameField != nil && nameField.Code != "" {
		qwhere.Where(nameField.Code+" = ?", flow.Request.Q)
	}
	var data map[string]interface{}
	if err := oql.Take(&data).Error(); err != nil {
		flow.Error(err)
		return
	} else if len(data) > 0 {
		flow.SetData(data)
		return
	}
}

func (s *mdSvImpl) QuotedBy(m MD, ids []string, excludes ...MD) ([]MDEntity, []string) {
	if m == nil || ids == nil || len(ids) == 0 {
		return nil, nil
	}
	excludeIds := make([]string, 0)
	if excludes != nil && len(excludes) > 0 {
		for _, e := range excludes {
			excludeIds = append(excludeIds, e.MD().ID)
		}
	}

	items := make([]MDField, 0)
	query := db.Default().Table(fmt.Sprintf("%v as f", db.Default().NewScope(MDField{}).TableName()))
	query = query.Joins(fmt.Sprintf("inner join %v as e on e.id=f.entity_id", db.Default().NewScope(MDEntity{}).TableName()))
	query = query.Select("f.*")
	if len(excludeIds) > 0 {
		query = query.Where("f.entity_id not in (?)", excludeIds)
	}
	query.Where("f.type_id=? and f.type_type=? and f.kind=?", m.MD().ID, "entity", "belongs_to").Find(&items)
	if len(items) > 0 {
		rtns := make([]MDEntity, 0)
		count := 0
		for _, d := range items {
			entity := MDSv().GetEntity(d.EntityID)
			if entity == nil || entity.TableName == "" {
				continue
			}
			if d.Kind == "belongs_to" {
				field := entity.GetField(d.ForeignKey)
				if field == nil {
					continue
				}
				db.Default().Table(fmt.Sprintf("%v as t", entity.TableName)).Where(fmt.Sprintf("%v in (?)", field.DbName), ids).Count(&count)
				if count > 0 {
					item := MDEntity{ID: entity.ID, Type: entity.Type, Name: entity.Name, TableName: entity.TableName}
					rtns = append(rtns, item)
				}
			}
		}
		if len(rtns) > 0 {
			s := make([]string, 0)
			for _, item := range rtns {
				s = append(s, item.Name)
			}
			return rtns, s
		}
	}
	return nil, nil
}
