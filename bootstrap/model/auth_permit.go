package model

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type AuthPermit struct {
	md.Model
	Code      string      `gorm:"unique_index:uix_code;not null" json:"code"`
	Name      string      `gorm:"not null" json:"name"`
	Memo      string      `json:"memo"`
	OwnerType string      `gorm:"size:36;unique_index:uix_code;name:拥有者类型;not null" json:"owner_type"`
	OwnerCode string      `gorm:"size:36;unique_index:uix_code;name:拥有者Code;not null" json:"owner_code"` //common为公共动作
	Enabled   utils.SBool `gorm:"not null;default:1;name:启用"`
}

func (s *AuthPermit) MD() *md.Mder {
	return &md.Mder{ID: "auth.permit", Domain: "auth", Name: "权限"}
}
