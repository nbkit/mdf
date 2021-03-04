package model

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type DtiHook struct {
	md.Model
	EntID    string      `gorm:"size:36"`
	EntityID string      `gorm:"size:36"`
	Code     string      `json:"code"`
	Name     string      `json:"name"`
	Memo     string      `json:"memo"`
	Method   string      `gorm:"size:10" json:"method"` //请求类型
	Path     string      `gorm:"size:50" json:"path"`
	Header   utils.SJson `gorm:"size:200" json:"header"`
	Body     utils.SJson `gorm:"size:500" json:"body"`
	Query    utils.SJson `gorm:"size:200" json:"query"`
	Enabled  utils.SBool `gorm:"default:true" json:"enabled"`
	Sequence int         `json:"sequence"`
}

func (s DtiHook) TableName() string {
	return "sys_dti_hooks"
}
func (s *DtiHook) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".dti.hook", Domain: MD_DOMAIN, Name: "接口钩子"}
}
