package model

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type Profile struct {
	md.Model
	EntID        string      `gorm:"size:50" json:"ent_id"`
	Code         string      `gorm:"size:50"`
	Name         string      `gorm:"size:50"`
	Memo         string      `gorm:"size:50"`
	Tags         string      `gorm:"size:100;name:标签" json:"tags"`
	System       utils.SBool `gorm:"not null;default:0;name:系统的"`
	DataTypeID   string      `gorm:"size:20" json:"data_type_id"`
	DataType     *md.MDEnum  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:sys.dti.param.type"`
	DefaultValue string
	Value        string
}

func (t Profile) TableName() string {
	return "sys_profiles"
}
func (s *Profile) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".profile", Domain: MD_DOMAIN, Name: "参数"}
}
