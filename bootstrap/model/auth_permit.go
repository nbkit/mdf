package model

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type AuthPermit struct {
	md.Model
	Code    string      `gorm:"unique_index:uix_code;not null" json:"code"`
	Name    string      `gorm:"not null" json:"name"`
	Memo    string      `json:"memo"`
	Widget  string      `gorm:"size:36;unique_index:uix_code;name:组件;not null" json:"widget"`
	Enabled utils.SBool `gorm:"not null;default:1;name:启用"`
}

func (s *AuthPermit) MD() *md.Mder {
	return &md.Mder{ID: "auth.permit", Domain: "auth", Name: "权限"}
}
