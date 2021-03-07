package model

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

/**
用户
*/
type User struct {
	md.Model
	Openid    string      `gorm:"size:50" json:"openid"`
	Mobile    string      `gorm:"size:50" json:"mobile"`
	Email     string      `gorm:"size:50" json:"email"`
	Account   string      `gorm:"size:50" json:"account"`
	Password  string      `json:"password"`
	Name      string      `gorm:"size:50" json:"name"`
	AvatarUrl string      `json:"avatar_url"`
	Memo      string      `json:"memo"`
	Token     string      `gorm:"size:50" json:"token"`
	IsSystem  utils.SBool `gorm:"not null;default:0;name:系统的" json:"is_system"`
	Enabled   utils.SBool `gorm:"not null;default:1;name:启用" json:"enabled"`
}

func (t User) TableName() string {
	return "sys_users"
}
func (s *User) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".user", Domain: MD_DOMAIN, Name: "用户"}
}
