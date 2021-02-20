package model

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type AuthRole struct {
	md.Model
	EntID   string      `gorm:"size:36;not null;unique_index:uix_code"`
	Code    string      `gorm:"size:36;not null;unique_index:uix_code" json:"code"`
	Name    string      `gorm:"not null" json:"name"`
	Memo    string      `json:"memo"`
	Enabled utils.SBool `gorm:"not null;default:1;name:启用"`
	System  utils.SBool `gorm:"not null;default:0;name:系统的"`
}

func (s *AuthRole) MD() *md.Mder {
	return &md.Mder{ID: "auth.role", Domain: "auth", Name: "角色"}
}
