package md

import (
	"github.com/nbkit/mdf/utils"
)

type Model struct {
	ID        string     `gorm:"primary_key;size:50" json:"id"`
	CreatedAt utils.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time `gorm:"name:更新时间" json:"updated_at"`
	SysState  string     `gorm:"-" json:"_state"`
}

func (s Model) String() string {
	return s.ID
}

func (s Model) GetKey() string {
	return s.ID
}
