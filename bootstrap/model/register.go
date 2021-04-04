package model

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

// 注册数据模型,提供数据层，按模块注册数据模型
func Register() {
	//log
	md.MDSv().Migrate(&Log{})
	if utils.Config.GetBool("EnabledFeature") {
		//sys
		md.MDSv().Migrate(&Client{}, &AuthToken{}, &CodeRule{}, &CodeValue{})
		//profile
		md.MDSv().Migrate(&Profile{})
		//product
		md.MDSv().Migrate(&Product{}, &ProductModule{}, &ProductService{})
		//user
		md.MDSv().Migrate(&User{})
		//ent
		md.MDSv().Migrate(&Ent{}, &EntUser{})
		//role
		md.MDSv().Migrate(&AuthRole{}, &AuthPermit{}, &AuthToken{}, &AuthRoleUser{}, &AuthRolePermit{}, &AuthRoleEntity{})
		//cron
		md.MDSv().Migrate(&CronEndpoint{}, &CronParam{}, &CronTask{}, &CronLog{})
		//oss
		md.MDSv().Migrate(&Oss{}, &OssObject{})
		//dti
		md.MDSv().Migrate(&DtiHook{})
	}
}
