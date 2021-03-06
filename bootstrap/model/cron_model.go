package model

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
)

const (
	CRON_UNIT_ENUM_once  = "once"
	CRON_UNIT_ENUM_month = "month"
	CRON_UNIT_ENUM_week  = "week"
	CRON_UNIT_ENUM_day   = "day"
	CRON_UNIT_ENUM_hour  = "hour"
)

type CronEndpoint struct {
	md.Model
	Code      string      `gorm:"size:36" json:"code"`
	Name      string      `gorm:"size:36" json:"name"`
	Memo      string      `gorm:"size:255" json:"memo"`
	Type      string      `gorm:"size:36;not null" json:"type"` // 类型：action,http
	Domain    string      `gorm:"size:36;name:模块" json:"domain"`
	Action    string      `gorm:"size:36" json:"action"`
	OwnerType string      `gorm:"size:36;name:拥有者类型" json:"owner_type"`
	OwnerCode string      `gorm:"size:36;name:拥有者Code" json:"owner_code"`
	Tag       string      `gorm:"size:255" json:"tag"`
	Method    string      `gorm:"size:10" json:"method"` //请求类型
	Path      string      `gorm:"size:50" json:"path"`
	Header    utils.SJson `gorm:"size:200" json:"header"`
	Body      string      `gorm:"size:500" json:"body"`
	Query     utils.SJson `gorm:"size:200" json:"query"`
	Enabled   utils.SBool `gorm:"default:true" json:"enabled"`
	System    utils.SBool `gorm:"not null;default:0;name:系统的" json:"system"`
}

func (s CronEndpoint) TableName() string {
	return "sys_cron_endpoints"
}
func (s *CronEndpoint) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".cron.endpoint", Domain: MD_DOMAIN, Name: "计划"}
}

type CronParam struct {
	md.Model
	EndpointID string        `gorm:"size:36" json:"endpoint_id"`
	Endpoint   *CronEndpoint `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Code       string        `json:"code"`
	Name       string        `json:"name"`
	Memo       string        `json:"memo"`
	Type       string        `gorm:"size:20;name:类型;not null" json:"type"`
	Required   utils.SBool   `gorm:"default:0" json:"required"`
	Hidden     utils.SBool   `gorm:"default:0" json:"hidden"`
	Sequence   int           `gorm:"default:0" json:"sequence"`
	Operator   string        `gorm:"size:36;name:操作符号" json:"operator"`
	RefType    string        `gorm:"size:36;name:参照类型" json:"ref_type"`
	RefCode    string        `gorm:"size:36;name:参照编码" json:"ref_code"`
	RefReturn  string        `gorm:"size:36;name:参照返回" json:"ref_return"`
	RefFilter  string        `gorm:"type:text;name:参照查询条件" json:"ref_filter"`
	Value1     utils.SJson   `gorm:"type:text;name:值1" json:"value1"`
	Value2     utils.SJson   `gorm:"type:text;name:值2" json:"value2"`
	Enabled    utils.SBool   `gorm:"default:true" json:"enabled"`
}

func (s CronParam) TableName() string {
	return "sys_cron_params"
}
func (s *CronParam) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".cron.param", Domain: MD_DOMAIN, Name: "参数"}
}

type CronTask struct {
	md.Model
	ClientID     string             `gorm:"size:50"`
	EntID        string             `gorm:"size:50"`
	UserID       string             `gorm:"size:50"`
	Tag          string             `json:"tag"`
	Code         string             `json:"code"`
	Name         string             `json:"name"`
	Memo         string             `json:"memo"`
	StatusID     string             `gorm:"size:50" json:"status_id"` //waiting,running,completed
	Status       *md.MDEnum         `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:sys.corn.status"`
	Enabled      utils.SBool        `gorm:"default:true" json:"enabled"`
	Context      utils.TokenContext `gorm:"type:text;name:上下文" json:"context"` //上下文参数
	NeTime       int64              `json:"ne_time"`                           //下次执行时间 ,unix
	FmTime       utils.Time         `json:"fm_time"`                           //开始
	ToTime       utils.Time         `json:"to_time"`                           //结束
	UnitID       string             `json:"unit_id"`                           //sys.cron.unit,频率 执行一次,每月,每周,每天,每小时
	Unit         *md.MDEnum         `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:sys.cron.unit"`
	Cycle        int                `json:"cycle"` //周期,1天，2月，4小时
	Spec         string             `json:"spec"`
	Retry        int                `json:"retry"`
	EndpointID   string             `gorm:"size:36" json:"endpoint_id"`
	Endpoint     *CronEndpoint      `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Header       utils.SJson        `gorm:"size:200" json:"header"`
	Body         string             `gorm:"size:500" json:"body"`
	Query        utils.SJson        `gorm:"size:200" json:"query"`
	NumRun       int                //执行次数
	NumSuccess   int                //成功次数
	NumFailed    int                //失败次数
	NumPeriod    int                //当前频率执行次数
	LastMsg      string             `json:"last_msg"`
	LastStatusID string             `json:"last_status_id"` //running,succeed,failed
	LastStatus   *md.MDEnum         `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:sys.corn.status"`
	LastTime     utils.Time         `json:"last_time"`
	System       utils.SBool        `gorm:"not null;default:0;name:系统的" json:"system"`
}

func (s CronTask) TableName() string {
	return "sys_cron_tasks"
}
func (s *CronTask) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".cron.task", Domain: MD_DOMAIN, Name: "计划任务"}
}

/**
日志
*/
type CronLog struct {
	md.Model
	EntID      string     `gorm:"size:50;index:idx_ent_type"`
	EndpointID string     `gorm:"size:36" json:"endpoint_id"`
	TaskID     string     `gorm:"size:50;index:idx_ent_type" json:"task_id"`
	Task       *CronTask  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Title      string     `json:"title"`
	Msg        string     `gorm:"type:text" json:"msg"`
	StatusID   string     `json:"status_id"` //running,succeed,failed
	Status     *md.MDEnum `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:sys.corn.status"`
}

func (s CronLog) TableName() string {
	return "sys_cron_logs"
}
func (s *CronLog) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".cron.log", Domain: MD_DOMAIN, Name: "任务日志"}
}
