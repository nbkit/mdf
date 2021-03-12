package model

import (
	"encoding/json"
	"github.com/nbkit/mdf/log"

	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type AuthToken struct {
	md.Model
	OwnerType string `gorm:"size:36;index:owner_code;name:拥有者类型;not null" json:"owner_type"`
	OwnerCode string `gorm:"size:36;index:owner_code;name:拥有者Code;not null" json:"owner_code"`
	OwnerID   string `gorm:"size:36;index:owner_code;name:拥有者Code;not null" json:"owner_id"`
	ClientID  string `gorm:"size:50"`
	EntID     string `gorm:"size:50;index:idx_type"`
	EntName   string `gorm:"size:50"`
	UserID    string `gorm:"size:50"`
	UserName  string `gorm:"size:50"`
	Token     string `gorm:"size:100;unique_index:idx_token;not null"`
	Name      string `gorm:"size:50"`
	Content   string
	Scope     string `gorm:"size:200"`
	ExpireAt  utils.Time
}

func (s *AuthToken) MD() *md.Mder {
	return &md.Mder{ID: "auth.token", Domain: "auth", Name: "令牌"}
}

func (s *AuthToken) SetContent(value interface{}) error {
	str, err := json.Marshal(value)
	if err != nil {
		log.Errorf("error:%v", err)
		return err
	}
	s.Content = string(str)
	return nil
}
func (s *AuthToken) GetContent(value interface{}) error {
	if err := json.Unmarshal([]byte(s.Content), value); err != nil {
		log.Errorf("parse content error:%v", err)
		return err
	}
	return nil
}
