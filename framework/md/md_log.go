package md

import (
	"github.com/nbkit/mdf/utils"
)

type MDLog struct {
	ID         string     `gorm:"primary_key;size:50" json:"id"`
	CreatedAt  utils.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt  utils.Time `gorm:"name:更新时间" json:"updated_at"`
	EntID      string     `gorm:"size:50;index:idx_ent_type"`
	Type       string     `gorm:"not null;size:10;index:idx_ent_type;name:日志类型"` //login,op
	UserID     string     `gorm:"size:50"`
	UserName   string     `gorm:"size:50"`
	NodeType   string     `gorm:"size:50;index:idx_ent_type" json:"node_type"`
	NodeID     string     `gorm:"size:120;index:idx_ent_type" json:"node_id"`
	NodeAction string     `gorm:"size:50;name:动作" json:"node_action"`
	DataID     string     `gorm:"size:50" json:"data_id"`
	DataType   string     `gorm:"size:50" json:"data_type"`
	Title      string     `gorm:"size:50"`
	Msg        string     `gorm:"type:text" json:"msg"`
	Level      string     `gorm:"size:50" json:"level"` //error,warn,info,debug
	Status     string     `gorm:"size:36;name:状态" json:"status"`
	ReqIP      string     `gorm:"size:50"`
	ReqClient  string     `gorm:"size:50;name:设备" json:"req_client"`
	ReqAgent   string     `gorm:"size:500"`
}

func (s MDLog) Clone() *MDLog {
	return &MDLog{
		Type:       s.Type,
		EntID:      s.EntID,
		UserID:     s.UserID,
		UserName:   s.UserName,
		Title:      s.Title,
		NodeID:     s.NodeID,
		NodeAction: s.NodeAction,
		NodeType:   s.NodeType,
		DataID:     s.DataID,
		DataType:   s.DataType,
		Msg:        s.Msg,
		Level:      s.Level,
		Status:     s.Status,
		ReqIP:      s.ReqIP,
		ReqAgent:   s.ReqAgent,
	}
}
func (s *MDLog) SetMsg(msg interface{}) *MDLog {
	if msg == nil {
		return s
	} else if m, ok := msg.(string); ok {
		s.Msg = m
	} else if m, ok := msg.(error); ok {
		s.Msg = m.Error()
	} else if m, ok := msg.(utils.GError); ok {
		s.Msg = m.Error()
	}
	return s
}
func (s *MDLog) MD() *Mder {
	return &Mder{ID: md_domain + ".log", Domain: md_domain, Name: "日志"}
}
