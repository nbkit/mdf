package md

import (
	"fmt"
	"github.com/nbkit/mdf/decimal"
	"strings"

	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
)

type MDEntity struct {
	ID        string     `gorm:"primary_key;size:50" json:"id"`
	CreatedAt utils.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time `gorm:"name:更新时间" json:"updated_at"`
	Type      string     `gorm:"size:50;not null;name:类型"` // simple，entity，enum，interface，dto,view
	Domain    string     `gorm:"size:50;name:领域" json:"domain"`
	Code      string     `gorm:"size:100;index:code_idx;not null;name:编码"`
	Name      string     `gorm:"size:100;name:实体名称"`
	TableName string     `gorm:"size:50;name:表名"`
	Memo      string     `gorm:"size:500;name:备注"`
	Tags      string     `gorm:"size:500;name:标签"`
	System    utils.SBool
	Fields    []MDField `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:EntityID"`
	cache     map[string]MDField
}

func (s MDEntity) TableComment() string {
	return "实体"
}
func (s MDEntity) String() string {
	return fmt.Sprintf("%s-%s-%s", s.Domain, s.Code, s.ID)
}
func (s *MDEntity) MD() *Mder {
	return &Mder{ID: "md.entity", Domain: MD_domain, Name: s.TableComment()}
}
func (s *MDEntity) GetField(code string) *MDField {
	if s.cache == nil {
		s.cache = make(map[string]MDField)
	}

	if s.Fields != nil && len(s.Fields) > 0 && len(s.cache) == 0 {
		for i, v := range s.Fields {
			s.cache[strings.ToLower(v.Code)] = s.Fields[i]
			if v.DbName != "" {
				s.cache[strings.ToLower(v.DbName)] = s.Fields[i]
			}
		}
	}
	if v, ok := s.cache[strings.ToLower(code)]; ok {
		return &v
	}
	return nil
}

type MDEntityRelation struct {
	ID        string     `gorm:"primary_key;size:50" json:"id"`
	CreatedAt utils.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time `gorm:"name:更新时间" json:"updated_at"`
	Code      string     `gorm:"size:100;index:code_idx;not null"`
	Name      string     `gorm:"size:100"`
	ParentID  string     `gorm:"size:50;name:父实体;not null"`
	ChildID   string     `gorm:"size:50;name:子实体;not null"`
	Kind      string     `gorm:"name:类型;not null"` //inherit:继承，interface:接口，
	ParentKey string     `gorm:"size:36;name:主键" json:"parent_key"`
	ChildKey  string     `gorm:"size:36;name:外键" json:"child_key"`
	Limit     string     `gorm:"size:500;name:限制"`
}

func (s MDEntityRelation) TableComment() string {
	return "实体关系"
}
func (s *MDEntityRelation) MD() *Mder {
	return &Mder{ID: "md.entity.relation", Domain: MD_domain, Name: s.TableComment()}
}

type MDField struct {
	ID             string      `gorm:"primary_key;size:50" json:"id"`
	CreatedAt      utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt      utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	EntityID       string      `gorm:"size:50;unique_index:uix_code;not null;name:实体ID"`
	Entity         *MDEntity   `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Code           string      `gorm:"size:50;unique_index:uix_code;not null;name:编码"`
	Name           string      `gorm:"size:50;name:名称"`
	DbName         string      `gorm:"size:50;name:数据列名"`
	IsNormal       utils.SBool `gorm:"name:是否数据库字段"`
	IsPrimaryKey   utils.SBool `gorm:"name:是否主键"`        //外键
	ForeignKey     string      `gorm:"size:50;name:外键"`  //外键
	AssociationKey string      `gorm:"size:50;name:关联键"` //Association
	Kind           string      `gorm:"size:50;name:关系"`
	TypeID         string      `gorm:"size:50;name:数据类型"`
	TypeType       string      `gorm:"size:50;name:数据类别"`
	Type           *MDEntity   `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Limit          string      `gorm:"size:500;name:限制"`
	Memo           string      `gorm:"size:500"`
	Tags           string      `gorm:"size:500"` // code,name,ent,import
	Sequence       int
	Nullable       utils.SBool
	Length         int
	Precision      int
	DefaultValue   string
	MinValue       string
	MaxValue       string
	SrcID          string `gorm:"size:50" json:"src_id"`
}

func (s MDField) String() string {
	return fmt.Sprintf("%s-%s-%s", s.Code, s.Name, s.TypeType)
}
func (s MDField) TableComment() string {
	return "属性"
}
func (s MDField) CompileValue(value interface{}) interface{} {
	if value == nil || value == "" || s.TypeID == "" {
		return nil
	}
	if s.TypeID == utils.FIELD_TYPE_STRING || s.TypeID == utils.FIELD_TYPE_TEXT || s.TypeID == utils.FIELD_TYPE_XML {
		return value
	}
	if s.TypeID == utils.FIELD_TYPE_INT {
		return utils.ToInt(value)
	}
	if s.TypeID == utils.FIELD_TYPE_BOOL {
		return utils.ToSBool(value)
	}
	if s.TypeID == utils.FIELD_TYPE_JSON {
		return utils.ToSBool(value)
	}
	if s.TypeID == utils.FIELD_TYPE_DATE || s.TypeID == utils.FIELD_TYPE_DATETIME {
		return utils.ToTime(value)
	}
	if s.TypeID == utils.FIELD_TYPE_DECIMAL {
		if v, err := decimal.NewFromString(utils.ToString(value)); err != nil {
			return log.ErrorD(err)
		} else {
			return v
		}
	}
	return nil
}
func (s *MDField) MD() *Mder {
	return &Mder{ID: "md.field", Domain: MD_domain, Name: s.TableComment()}
}
