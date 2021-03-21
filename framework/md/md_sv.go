package md

import (
	"fmt"
	"github.com/nbkit/mdf/db"
	"strings"
	"sync"

	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
)

type IMDSv interface {
	Migrate(values ...interface{})
	AddEntities(items []MDEntity) error

	GetEntity(id string) *MDEntity
	GetEnum(typeId string, values ...string) *MDEnum

	TakeDataByQ(flow *utils.FlowContext)
	UpdateEntity(item MDEntity) error
	UpdateEnumType(enumType MDEnumType) error

	QuotedBy(m MD, ids []string, excludes ...MD) ([]MDEntity, []string)

	Cache()
}

type mdSvImpl struct {
	*sync.Mutex
	entityCache     map[string]*MDEntity
	enumCache       map[string]*MDEnum
	initMDCompleted bool
}

func MDSv() IMDSv {
	return mdSv
}

var mdSv IMDSv = newMDSv()

func newMDSv() *mdSvImpl {
	return &mdSvImpl{
		Mutex:       &sync.Mutex{},
		entityCache: make(map[string]*MDEntity),
		enumCache:   make(map[string]*MDEnum),
	}
}

func (s *mdSvImpl) Cache() {
	s.cacheEnums()
	s.cacheEntities()
}

func (s *mdSvImpl) cacheEnums() {
	s.enumCache = make(map[string]*MDEnum)
	items, _ := s.GetEnums()
	for i, _ := range items {
		s.cacheEnum(items[i])
	}
}
func (s *mdSvImpl) cacheEnum(item MDEnum) {
	s.enumCache[strings.ToLower(item.EntityID+":"+item.ID)] = &item
	s.enumCache[strings.ToLower(item.EntityID+":"+item.Name)] = &item
}
func (s *mdSvImpl) cacheEntities() {
	s.entityCache = make(map[string]*MDEntity)
	items, _ := s.GetEntities()
	for i, _ := range items {
		s.cacheEntity(items[i])
	}
}
func (s *mdSvImpl) cacheEntity(item MDEntity) {
	s.entityCache[strings.ToLower(item.ID)] = &item
}

func (s *mdSvImpl) Migrate(values ...interface{}) {
	//先增加模型表
	if !s.initMDCompleted {
		s.initMDCompleted = true
		mds := []interface{}{
			&MDEntity{}, &MDEntityRelation{}, &MDField{}, &MDEnum{},
			&MDAction{}, &MDRule{},
			&MDWidget{}, &MDWidgetDs{}, &MDWidgetLayout{}, &MDWidgetItem{},
			&MDToolbars{}, &MDToolbarItem{},
			&MDAction{}, &MDRule{},
			&MDFilters{}, &MDFilterSolution{}, &MDFilterItem{},
		}
		needDb := make([]interface{}, 0)
		for _, v := range mds {
			m := newMd(v)
			if dd := m.GetMder(); dd == nil || dd.Type == utils.TYPE_ENTITY || dd.Type == utils.TYPE_ENUM || dd.Type == "" {
				needDb = append(needDb, v)
			}
		}
		db.Default().AutoMigrate(needDb...)
		log.Print("AutoMigrate MD")
		for _, v := range mds {
			m := newMd(v)
			m.Migrate()
		}

		initData()
	}
	if len(values) > 0 {
		needDb := make([]interface{}, 0)
		for _, v := range values {
			m := newMd(v)
			if dd := m.GetMder(); dd == nil || dd.Type == utils.TYPE_ENTITY || dd.Type == utils.TYPE_ENUM || dd.Type == "" {
				needDb = append(needDb, v)
			}
			m.Migrate()
		}
		if err := db.Default().AutoMigrate(needDb...).Error; err != nil {
			log.ErrorD(err)
		}
	}

}
func (s *mdSvImpl) entityToTables(items []MDEntity, oldItems []MDEntity) error {
	if items == nil || len(items) == 0 {
		return nil
	}
	for i, _ := range items {
		entity := items[i]
		oldItem := MDEntity{}
		for oi, ov := range oldItems {
			if entity.ID == ov.ID {
				oldItem = oldItems[oi]
				break
			}
		}
		if entity.Type != utils.TYPE_ENTITY {
			continue
		}
		if entity.TableName == "" {
			entity.TableName = strings.ReplaceAll(entity.ID, ".", "_")
		}
		if db.Default().Dialect().HasTable(entity.TableName) {
			s.updateTable(entity, oldItem)
		} else {
			s.createTable(entity)
		}
	}
	return nil
}
func (s *mdSvImpl) createTable(item MDEntity) {
	if len(item.Fields) == 0 {
		return
	}
	var primaryKeys []string
	var tags []string
	for i, _ := range item.Fields {
		field := item.Fields[i]
		if field.IsNormal.IsTrue() {
			tags = append(tags, s.buildColumnNameString(field))
			if field.IsPrimaryKey.IsTrue() {
				primaryKeys = append(primaryKeys, s.quote(field.DbName))
			}
		}
	}
	var primaryKeyStr string
	if len(primaryKeys) > 0 {
		primaryKeyStr = fmt.Sprintf(", PRIMARY KEY (%v)", strings.Join(primaryKeys, ","))
	}
	var tableOptions string
	if err := db.Default().Exec(fmt.Sprintf("CREATE TABLE %v (%v %v)%s", s.quote(item.TableName), strings.Join(tags, ","), primaryKeyStr, tableOptions)).Error; err != nil {
		log.ErrorD(err)
	}
}
func (s *mdSvImpl) quote(str string) string {
	return db.Default().Dialect().Quote(str)
}
func (s *mdSvImpl) updateTable(item MDEntity, old MDEntity) {
	//更新栏目
	for i := range item.Fields {
		field := item.Fields[i]
		if field.DbName == "-" || field.DbName == "" {
			continue
		}
		oldField := MDField{}
		for oi, ov := range old.Fields {
			if strings.ToLower(field.Code) == strings.ToLower(ov.Code) {
				oldField = old.Fields[oi]
				break
			}
		}
		if !field.IsNormal.IsTrue() {
			continue
		}
		newString := s.buildColumnNameString(field)
		if db.Default().Dialect().HasColumn(item.TableName, field.DbName) { //字段已存在
			oldString := s.buildColumnNameString(oldField)
			//修改字段类型、类型长度、默认值、注释
			if oldString != newString && strings.Contains(item.Tags, "update") {
				if err := db.Default().Exec(fmt.Sprintf("ALTER TABLE %v MODIFY %v", s.quote(item.TableName), newString)).Error; err != nil {
					log.ErrorD(err)
				}
			}
		} else { //新增字段
			if err := db.Default().Exec(fmt.Sprintf("ALTER TABLE %v ADD %v", s.quote(item.TableName), newString)).Error; err != nil {
				log.ErrorD(err)
			}
		}
	}
}
func (s *mdSvImpl) buildColumnNameString(item MDField) string {
	/*
		column_definition:
		data_type [NOT NULL | NULL] [DEFAULT {literal | (expr)} ]
		[AUTO_INCREMENT] [UNIQUE [KEY]] [[PRIMARY] KEY]
		[COMMENT 'string']

	*/
	dialectName := db.Default().Dialect().GetName()
	if dialectName == "godror" || dialectName == "oracle" {
		return s.buildColumnNameString4Oracle(item)
	} else {
		return s.buildColumnNameString4Mysql(item)
	}
}

func (s *mdSvImpl) AddEntities(items []MDEntity) error {
	entityIds := make([]string, 0)
	oldEntities := make([]MDEntity, 0)
	for i, _ := range items {
		entity := items[i]
		if entity.ID == "" {
			continue
		}
		oldEntity := MDEntity{}
		if db.Default().Model(oldEntity).Preload("Fields").Order("id").Where("id=?", entity.ID).Take(&oldEntity); oldEntity.ID != "" {
			oldEntities = append(oldEntities, oldEntity)
			datas := make(map[string]interface{})
			if oldEntity.TableName != entity.TableName {
				datas["TableName"] = entity.TableName
			}
			if oldEntity.Type != entity.Type {
				datas["Type"] = entity.Type
			}
			if oldEntity.Code != entity.Code {
				datas["Code"] = entity.Code
			}
			if oldEntity.Tags != entity.Tags {
				datas["Tags"] = entity.Tags
			}
			if entity.System.Valid() && oldEntity.System.NotEqual(entity.System) {
				datas["System"] = entity.System
			}
			if oldEntity.Domain != entity.Domain {
				datas["Domain"] = entity.Domain
			}
			if oldEntity.Name != entity.Name {
				datas["Name"] = entity.Name
			}
			if oldEntity.Memo != entity.Memo {
				datas["Memo"] = entity.Memo
			}
			if len(datas) > 0 {
				db.Default().Model(MDEntity{}).Where("id=?", oldEntity.ID).Updates(datas)
			}
		} else {
			if entity.Type == utils.TYPE_ENTITY && entity.TableName == "" {
				entity.TableName = strings.ReplaceAll(entity.ID, ".", "_")
			}
			db.Default().Create(&entity)
		}
		entityIds = append(entityIds, entity.ID)
	}
	//属性字段
	for i, _ := range items {
		entity := items[i]
		if entity.ID != "" && entity.Type == utils.TYPE_ENTITY && len(entity.Fields) > 0 {
			itemCodes := make([]string, 0)
			for f, _ := range entity.Fields {
				field := entity.Fields[f]
				itemCodes = append(itemCodes, field.Code)
				field.IsNormal = utils.SBool_True
				if field.DbName == "-" {
					field.IsNormal = utils.SBool_False
				}
				if field.DbName == "" && field.IsNormal.IsTrue() {
					field.DbName = utils.SnakeString(field.Code)
				}
				oldField := MDField{}
				if fieldType := s.GetEntity(field.TypeID); fieldType != nil {
					field.TypeType = fieldType.Type
					field.TypeID = fieldType.ID
				}
				if field.TypeType == utils.TYPE_ENTITY { //实体
					if field.Kind == "" {
						field.Kind = "belongs_to"
					}
					if field.ForeignKey == "" {
						field.ForeignKey = fmt.Sprintf("%sID", field.Code)
					}
					if field.AssociationKey == "" {
						field.AssociationKey = "ID"
					}
					field.IsNormal = utils.SBool_False
				} else if field.TypeType == utils.TYPE_ENUM { //枚举
					if field.Kind == "" {
						field.Kind = "belongs_to"
					}
					if field.ForeignKey == "" {
						field.ForeignKey = fmt.Sprintf("%sID", field.Code)
					}
					if field.AssociationKey == "" {
						field.AssociationKey = "ID"
					}
					field.IsNormal = utils.SBool_False
				}
				if db.Default().Model(MDField{}).Order("id").Where("entity_id=? and code=?", entity.ID, field.Code).Take(&oldField); oldField.ID != "" {
					datas := make(map[string]interface{})
					if oldField.Name != field.Name {
						datas["Name"] = field.Name
					}
					if oldField.Tags != field.Tags {
						datas["Tags"] = field.Tags
					}
					if oldField.DbName != field.DbName {
						datas["DbName"] = field.DbName
					}
					if oldField.IsNormal != field.IsNormal {
						datas["IsNormal"] = field.IsNormal
					}
					if oldField.IsPrimaryKey != field.IsPrimaryKey {
						datas["IsPrimaryKey"] = field.IsPrimaryKey
					}
					if oldField.Length != field.Length {
						datas["Length"] = field.Length
					}
					if oldField.Nullable != field.Nullable {
						datas["Nullable"] = field.Nullable
					}
					if oldField.DefaultValue != field.DefaultValue {
						datas["DefaultValue"] = field.DefaultValue
					}
					if oldField.TypeID != field.TypeID {
						datas["TypeID"] = field.TypeID
					}
					if oldField.TypeType != field.TypeType {
						datas["TypeType"] = field.TypeType
					}
					if oldField.Limit != field.Limit {
						datas["Limit"] = field.Limit
					}
					if oldField.MinValue != field.MinValue {
						datas["MinValue"] = field.MinValue
					}
					if oldField.MaxValue != field.MaxValue {
						datas["MaxValue"] = field.MaxValue
					}
					if oldField.Precision != field.Precision {
						datas["Precision"] = field.Precision
					}
					if oldField.AssociationKey != field.AssociationKey {
						datas["AssociationKey"] = field.AssociationKey
					}
					if oldField.ForeignKey != field.ForeignKey {
						datas["ForeignKey"] = field.ForeignKey
					}
					if oldField.Kind != field.Kind {
						datas["Kind"] = field.Kind
					}
					if oldField.Sequence != field.Sequence {
						datas["Sequence"] = field.Sequence
					}
					if oldField.SrcID != field.SrcID && field.SrcID != "" {
						datas["SrcID"] = field.SrcID
					}
					if len(datas) > 0 {
						db.Default().Model(MDField{}).Where("entity_id=? and code=?", entity.ID, field.Code).Updates(datas)
					}
				} else {
					db.Default().Create(&field)
				}
			}
			db.Default().Delete(MDField{}, "entity_id=? and code not in (?)", entity.ID, itemCodes)
		}
	}
	//枚举
	for _, entity := range items {
		if entity.ID != "" && entity.Type == utils.TYPE_ENUM && len(entity.Fields) > 0 {
			itemCodes := make([]string, 0)
			for f, field := range entity.Fields {
				newEnum := MDEnum{ID: field.Code, EntityID: entity.ID, Sequence: f, Name: field.Name}
				oldEnum := MDEnum{}
				itemCodes = append(itemCodes, newEnum.ID)
				if db.Default().Model(oldEnum).Order("id").Where("entity_id=? and id=?", newEnum.EntityID, newEnum.ID).Take(&oldEnum); oldEnum.ID != "" {
					datas := make(map[string]interface{})
					if oldEnum.Name != field.Name {
						datas["Name"] = field.Name
					}
					if oldEnum.Sequence != newEnum.Sequence {
						datas["Sequence"] = newEnum.Sequence
					}
					if len(datas) > 0 {
						db.Default().Model(MDEnum{}).Where("entity_id=? and id=?", oldEnum.EntityID, oldEnum.ID).Updates(datas)
					}
				} else {
					db.Default().Create(&newEnum)
				}
			}
			db.Default().Delete(MDEnum{}, "entity_id=? and id not in (?)", entity.ID, itemCodes)
		}
	}
	if len(entityIds) > 0 {
		toTables := make([]MDEntity, 0)
		db.Default().Model(MDEntity{}).Preload("Fields").Where("id in (?) and type=?", entityIds, utils.TYPE_ENTITY).Find(&toTables)
		return s.entityToTables(toTables, oldEntities)
	}
	//缓存
	s.Cache()
	return nil
}
