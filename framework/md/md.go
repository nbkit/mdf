package md

import (
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/decimal"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/nbkit/mdf/db/gorm"
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
)

type MD interface {
	MD() *Mder
}
type Mder struct {
	ID     string
	Type   string
	Name   string
	Domain string
}

type md struct {
	Value interface{}
}

func newMd(value interface{}) *md {
	item := md{Value: value}
	return &item
}
func (m *md) GetMder() *Mder {
	if mder, ok := m.Value.(MD); ok {
		return mder.MD()
	}
	return nil
}
func (m *md) GetEntity() *MDEntity {
	mdInfo := m.GetMder()
	if mdInfo == nil {
		return nil
	}
	item := MDEntity{}
	query := db.Default().Model(item).Preload("Fields").Order("id").Where("id=?", mdInfo.ID)
	if err := query.Take(&item).Error; err != nil {
		log.ErrorD(err)
	} else {
		return &item
	}
	return nil
}

// Get Data Type for MySQL Dialect
func (s *md) dataTypeOf(field *gorm.StructField) string {
	size := 0
	if num, ok := field.TagSettingsGet("SIZE"); ok {
		size, _ = strconv.Atoi(num)
	} else {
		size = 255
	}
	var (
		reflectType = field.Struct.Type
	)

	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	fieldValue := reflect.Indirect(reflect.New(reflectType))
	sqlType := ""

	if sqlType == "" {
		switch fieldValue.Kind() {
		case reflect.Bool:
			sqlType = "boolean"
		case reflect.Int8:
		case reflect.Uint8:
		case reflect.Int, reflect.Int16, reflect.Int32:
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		case reflect.Int64:
		case reflect.Uint64:
			sqlType = "int"
		case reflect.Float32, reflect.Float64:
			sqlType = "decimal"
		case reflect.String:
			if size > 100 {

			}
			sqlType = "string"
		case reflect.Struct:
			if _, ok := fieldValue.Interface().(time.Time); ok {
				sqlType = "datetime"
			}
			if _, ok := fieldValue.Interface().(utils.Time); ok {
				sqlType = "datetime"
			}
			if _, ok := fieldValue.Interface().(decimal.Decimal); ok {
				sqlType = "decimal"
			}
		default:
			sqlType = "string"
		}
	}
	if sqlType == "" {
		sqlType = "string"
	}
	return sqlType
}

func (m *md) Migrate() {
	mdInfo := m.GetMder()
	if mdInfo == nil {
		return
	}
	if mdInfo.ID == "" {
		log.Error().String("Name", mdInfo.Name).Msg("元数据ID为空")
		return
	}
	scope := db.Default().NewScope(m.Value)

	entity := m.GetEntity()
	vt := reflect.ValueOf(m.Value).Elem().Type()
	newEntity := &MDEntity{TableName: scope.TableName(), Name: mdInfo.Name, Domain: mdInfo.Domain, Code: vt.Name(), Type: mdInfo.Type}
	if newEntity.Type == "" {
		newEntity.Type = utils.TYPE_ENTITY
	}
	newEntity.System = utils.SBool_True

	if entity == nil {
		entity = newEntity
		entity.ID = mdInfo.ID
		db.Default().Create(entity)
		entity = m.GetEntity()
	} else {
		updates := make(map[string]interface{})
		if entity.Name != newEntity.Name {
			updates["Name"] = newEntity.Name
		}
		if entity.Code != newEntity.Code {
			updates["Code"] = newEntity.Code
		}
		if entity.Type != newEntity.Type {
			updates["Type"] = newEntity.Type
		}
		if entity.TableName != newEntity.TableName {
			updates["TableName"] = newEntity.TableName
		}
		if entity.Domain != newEntity.Domain {
			updates["Domain"] = newEntity.Domain
		}
		if entity.System.NotEqual(newEntity.System) {
			updates["System"] = newEntity.System
		}
		if len(updates) > 0 {
			db.Default().Model(MDEntity{}).Where("id=?", entity.ID).Updates(updates)
			entity = m.GetEntity()
		}
	}
	if entity == nil {
		log.Error().String("Name", mdInfo.Name).Msg("元数据ID为空")
		return
	}
	codes := make([]string, 0)
	for _, field := range scope.GetModelStruct().StructFields {
		if field.IsIgnored {
			continue
		}

		newField := MDField{
			Code:         field.Name,
			DbName:       field.DBName,
			IsPrimaryKey: utils.ToSBool(field.IsPrimaryKey),
			IsNormal:     utils.ToSBool(field.IsNormal),
			Name:         field.TagSettings["NAME"],
			EntityID:     entity.ID,
			Tags:         field.TagSettings["TAG"],
			DefaultValue: field.TagSettings["DEFAULT"],
		}
		if size, ok := field.TagSettings["SIZE"]; ok {
			newField.Length = utils.ToInt(size)
		}
		if _, ok := field.TagSettings["NOT NULL"]; ok {
			newField.Nullable = utils.SBool_False
		}

		if newField.Name == "" {
			newField.Name = newField.Code
		}
		//普通数据库字段
		if field.IsNormal {
		}
		reflectType := field.Struct.Type
		if reflectType.Kind() == reflect.Slice {
			reflectType = field.Struct.Type.Elem()
		}
		if reflectType.Kind() == reflect.Ptr {
			reflectType = reflectType.Elem()
		}
		newField.Limit = field.TagSettings["LIMIT"]
		if relationship := field.Relationship; relationship != nil {
			newField.Kind = relationship.Kind
			newField.ForeignKey = strings.Join(relationship.ForeignFieldNames, ".")
			newField.AssociationKey = strings.Join(relationship.AssociationForeignFieldNames, ".")

			fieldValue := reflect.New(reflectType)
			if e, ok := fieldValue.Interface().(MD); ok {
				if eMd := e.MD(); eMd != nil {
					newField.TypeID = eMd.ID
					newField.TypeType = eMd.Type
				}
			}
		} else {
			fieldValue := reflect.New(reflectType)
			if e, ok := fieldValue.Interface().(MD); ok {
				if eMd := e.MD(); eMd != nil {
					newField.TypeID = eMd.ID
					newField.TypeType = eMd.Type
				}
			} else if e := m.dataTypeOf(field); e != "" {
				newField.TypeID = e
				newField.TypeType = utils.TYPE_SIMPLE
			}
		}
		if newField.TypeID != "" && newField.TypeType == "" {
			if typeEntity := MDSv().GetEntity(newField.TypeID); typeEntity != nil {
				newField.TypeType = typeEntity.Type
			}
		}
		codes = append(codes, newField.Code)
		oldField := entity.GetField(newField.Code)

		if oldField == nil {
			db.Default().Create(&newField)
		} else {
			updates := make(map[string]interface{})
			if oldField.Name != newField.Name {
				updates["Name"] = newField.Name
			}
			if oldField.Length != newField.Length {
				updates["Length"] = newField.Length
			}
			if oldField.DbName != newField.DbName {
				updates["DbName"] = newField.DbName
			}
			if oldField.AssociationKey != newField.AssociationKey {
				updates["AssociationKey"] = newField.AssociationKey
			}
			if oldField.ForeignKey != newField.ForeignKey {
				updates["ForeignKey"] = newField.ForeignKey
			}
			if oldField.IsNormal != newField.IsNormal {
				updates["IsNormal"] = newField.IsNormal
			}
			if oldField.IsPrimaryKey != newField.IsPrimaryKey {
				updates["IsPrimaryKey"] = newField.IsPrimaryKey
			}
			if oldField.Kind != newField.Kind {
				updates["Kind"] = newField.Kind
			}
			if oldField.TypeID != newField.TypeID {
				updates["TypeID"] = newField.TypeID
			}
			if oldField.TypeType != newField.TypeType {
				updates["TypeType"] = newField.TypeType
			}
			if oldField.DefaultValue != newField.DefaultValue {
				updates["DefaultValue"] = newField.DefaultValue
			}
			if oldField.Nullable.NotEqual(newField.Nullable) {
				updates["Nullable"] = newField.Nullable
			}
			if oldField.Limit != newField.Limit {
				updates["Limit"] = newField.Limit
			}
			if oldField.SrcID != newField.SrcID && newField.SrcID != "" {
				updates["SrcID"] = newField.SrcID
			}
			if oldField.Tags != newField.Tags {
				updates["Tags"] = newField.Tags
			}
			if len(updates) > 0 {
				db.Default().Model(MDField{}).Where("id=?", oldField.ID).Updates(updates)
			}
		}
	}
	//删除不存在的
	db.Default().Delete(MDField{}, "entity_id=? and code not in (?)", entity.ID, codes)
}
