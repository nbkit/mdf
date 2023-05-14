package rule

import (
	"fmt"
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type MDAction struct {
	ID         string      `gorm:"primary_key;size:36" json:"id"`
	CreatedAt  utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt  utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	Widget     string      `gorm:"size:36;unique_index:uix_code;name:部件;not null" json:"widget"` //common为公共动作
	Code       string      `gorm:"size:36;name:编码;not null" json:"code"`
	Name       string      `gorm:"size:50;name:名称;not null" json:"name"`
	Type       string      `gorm:"size:20;unique_index:uix_code;name:类型;not null" json:"type"`
	Action     string      `gorm:"size:50;unique_index:uix_code;name:动作;not null" json:"action"`
	Url        string      `gorm:"size:100;name:服务路径" json:"url"`
	Parameter  utils.SJson `gorm:"type:text;name:参数" json:"parameter"`
	Method     string      `gorm:"size:20;name:请求方式" json:"method"`
	Target     string      `gorm:"size:36;name:目标" json:"target"`
	PrevScript string      `gorm:"type:text;name:前置脚本" json:"prev_script"`
	Script     string      `gorm:"type:text;name:脚本" json:"upgrade"`
	PostScript string      `gorm:"type:text;name:后置脚本" json:"post_script"`
	Enabled    utils.SBool `gorm:"default:true;not null" json:"enabled"`
}

func (s MDAction) GetKey() string {
	return fmt.Sprintf("%s:%s", s.Widget, s.Code)
}
func (s *MDAction) MD() *md.Mder {
	return &md.Mder{ID: "md.action", Domain: md.MD_domain, Name: "组件命令"}
}
