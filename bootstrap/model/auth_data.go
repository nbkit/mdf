package model

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type AuthRolePermit struct {
	md.Model
	EntID    string `gorm:"size:50;not null;unique_index:uix_code;index:role"`
	PermitID string `gorm:"size:50;not null;unique_index:uix_code" json:"user_id"`
	RoleID   string `gorm:"size:50;not null;unique_index:uix_code;index:role" json:"role_id"`
}

func (s *AuthRolePermit) MD() *md.Mder {
	return &md.Mder{ID: "auth.role.permit", Domain: "auth", Name: "角色对应的权限"}
}

type AuthRoleUser struct {
	md.Model
	EntID  string `gorm:"size:50;not null;unique_index:uix_code;index:role"`
	UserID string `gorm:"size:50;not null;unique_index:uix_code" json:"user_id"`
	RoleID string `gorm:"size:50;not null;unique_index:uix_code;index:role" json:"role_id"`
}

func (s *AuthRoleUser) MD() *md.Mder {
	return &md.Mder{ID: "auth.role.user", Domain: "auth", Name: "角色对应的用户"}
}

type AuthRoleEntity struct {
	md.Model
	EntID    string      `gorm:"size:50;not null;index:idx_role;index:idx_entity"`
	RoleID   string      `gorm:"size:50;not null;index:idx_role" json:"role_id"`
	EntityID string      `gorm:"size:50;not null;index:idx_entity" json:"entity_id"` //类型ID
	Exp      string      `json:"exp"`
	Memo     string      `json:"memo"`
	IsRead   utils.SBool `gorm:"default:0;not null;name:读取权限" json:"is_read"`
	IsWrite  utils.SBool `gorm:"default:0;not null;name:写入权限" json:"is_write"`
	IsDelete utils.SBool `gorm:"default:0;not null;name:删除权限" json:"is_delete"`
}

func (s *AuthRoleEntity) MD() *md.Mder {
	return &md.Mder{ID: "auth.role.entity", Domain: "auth", Name: "角色对应的数据权限"}
}
