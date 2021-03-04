package model

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

type Oss struct {
	md.Model
	EntID           string      `gorm:"size:50"`
	Region          string      `gorm:"name:地域或者数据中心"`
	Endpoint        string      `gorm:"name:OSS访问域名"`
	Bucket          string      `gorm:"name:存储空间"`
	AccessKeyID     string      `gorm:"name:访问密钥" json:"access_key_id"`
	AccessKeySecret string      `gorm:"name:访问密钥" json:"access_key_secret"`
	Memo            string      `gorm:"name:备注"`
	IsDefault       utils.SBool `gorm:"name:默认"`
	Type            string      `gorm:"name:OSS类型"` //ent:自建OSS，lease：租用，local:本地OSS
}

func (t Oss) TableName() string {
	return "sys_osses"
}
func (s *Oss) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".oss", Domain: MD_DOMAIN, Name: "对象存储"}
}

type OssObject struct {
	md.Model
	EntID        string      `gorm:"size:50"`
	Type         string      `gorm:"size:50;name:对象类型"  json:"type"` //obj,dir
	Folder       string      `gorm:"name:文件夹" json:"folder" form:"folder"`
	Path         string      `gorm:"name:资源Path"`
	Code         string      `gorm:"size:50;name:代码" json:"code"`
	Name         string      `gorm:"size:100;name:原始文件" json:"name"`
	OriginalName string      `gorm:"size:100;name:原始文件名"`
	Tag          string      `gorm:"size:50;name:标识"  json:"tag"  form:"tag"`
	Size         int64       `gorm:"name:文件大小"  json:"size"`
	MimeType     string      `gorm:"size:100;name:Mime类型"   json:"mime_type"`
	Ext          string      `gorm:"size:50;name:扩展名" json:"ext"`
	OwnerID      string      `gorm:"size:50;name:文件拥有者" json:"owner_id"  form:"owner_id"`       //单据ID
	OwnerType    string      `gorm:"size:50;name:文件拥有者类型" json:"owner_type"  form:"owner_type"` //单据类型
	OssID        string      `gorm:"size:50" json:"oss_id"  form:"oss_id"`                      //存储空间ID
	OssType      string      `gorm:"size:50" json:"oss_type"`                                   //存储空间类型
	OssBucket    string      `gorm:"name:存储空间" json:"oss_bucket"`                               //存储空间
	Oss          *Oss        `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Permit       utils.SBool `gorm:"name:权限控制" json:"permit"   form:"permit"`
}

func (t OssObject) TableName() string {
	return "sys_oss_objects"
}
func (s *OssObject) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".oss.object", Domain: MD_DOMAIN, Name: "对象存储对象"}
}
