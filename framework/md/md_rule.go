package md

import (
	"fmt"
	"github.com/nbkit/mdf/utils"
)

type MDRule struct {
	ID        string      `gorm:"primary_key;size:50" json:"id"` //领域.规则：md.save，ui.save
	CreatedAt utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	Domain    string      `gorm:"size:36;name:模块" json:"domain"`
	Widget    string      `gorm:"size:36;unique_index:uix_code;name:部件;not null" json:"widget"` //common为公共动作
	Code      string      `gorm:"size:50;unique_index:uix_code;name:编码;not null" json:"code"`
	Name      string      `gorm:"size:50;name:名称" json:"name"`
	Action    string      `gorm:"size:50;index:idx_action;name:动作;not null" json:"action"`
	Url       string      `gorm:"size:100;name:服务路径" json:"url"`
	Sequence  int         `gorm:"size:3;name:顺序;default:50;not null" json:"sequence"`
	Replaced  string      `gorm:"size:50;name:被替换的" json:"replaced"`
	Async     utils.SBool `gorm:"default:false;not null;name:异步的" json:"async"`
	Enabled   utils.SBool `gorm:"default:true;not null" json:"enabled"`
}

func (s MDRule) GetKey() string {
	return fmt.Sprintf("%s:%s", s.Widget, s.Code)
}
func (s *MDRule) MD() *Mder {
	return &Mder{ID: "md.rule", Domain: md_domain, Name: "动作规则"}
}
