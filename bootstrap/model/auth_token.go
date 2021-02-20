package model

import (
	"encoding/json"

	"github.com/nbkit/mdf/framework/glog"
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type AuthToken struct {
	md.Model
	ClientID  string `gorm:"size:50"`
	EntID     string `gorm:"size:50;index:idx_type"`
	EntName   string `gorm:"size:50"`
	UserID    string `gorm:"size:50"`
	UserName  string `gorm:"size:50"`
	ProductID string `gorm:"size:50"`
	ServiceID string `gorm:"size:50"`
	Token     string `gorm:"size:100;index:idx_token;not null"`
	Name      string
	Type      string `gorm:"size:50;index:idx_type;not null"`
	Content   string
	Scope     string
	ExpireAt  utils.Time
}

func (s *AuthToken) MD() *md.Mder {
	return &md.Mder{ID: "auth.token", Domain: "auth", Name: "令牌"}
}

func (s *AuthToken) SetContent(value interface{}) error {
	str, err := json.Marshal(value)
	if err != nil {
		glog.Errorf("error:%v", err)
		return err
	}
	s.Content = string(str)
	return nil
}
func (s *AuthToken) GetContent(value interface{}) error {
	if err := json.Unmarshal([]byte(s.Content), value); err != nil {
		glog.Errorf("parse content error:%v", err)
		return err
	}
	return nil
}
