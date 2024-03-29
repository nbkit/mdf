package services

import (
	"bytes"
	"encoding/json"
	"github.com/nbkit/mdf/bootstrap/errors"
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/log"
	"github.com/robfig/cron"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/nbkit/mdf/utils"

	"github.com/nbkit/mdf/bootstrap/model"
)

var cronCache *cron.Cron

type ICronSv interface {
	Start()
	CreateTask(entID string, item *model.CronTask) (*model.CronTask, error)
	DestroyTask(entID string, ids []string) error
}
type cronSvImpl struct {
	*sync.Mutex
}

func CronSv() ICronSv {
	return cronSv
}

var cronSv ICronSv = newCronSv()

func newCronSv() *cronSvImpl {
	return &cronSvImpl{Mutex: &sync.Mutex{}}
}
func (s *cronSvImpl) Start() {
	if cronCache == nil {
		cronCache = cron.New()
		cronCache.AddFunc("@every 20s", s.runJob)
		cronCache.Start()
	}
}

func (s *cronSvImpl) runJob() {
	s.Lock()
	defer func() {
		s.Unlock()
	}()
	// 获取需要执行的任务,可用的，等待执行，且已到执行时间
	items := make([]model.CronTask, 0)
	db.Default().Where("enabled=1 and status_id in (?) and ne_time<=?", []string{"waiting"}, utils.TimeNow().Unix()).Preload("Endpoint").Find(&items)
	if len(items) == 0 {
		return
	}
	for i, _ := range items {
		s.jobHandle(&items[i])
	}
	time.Sleep(time.Second)
}

func (s *cronSvImpl) jobHandle(item *model.CronTask) {
	defer func() {
		if r := recover(); r != nil {
			s.runJobItemReset(item)
		}
	}()
	item.NumRun = item.NumRun + 1
	db.Default().Model(item).Where("id=?", item.ID).Updates(map[string]interface{}{"StatusID": "running", "NumRun": item.NumRun, "LastStatusID": "running", "LastTime": utils.TimeNow()})

	s.createLog(&model.CronLog{EntID: item.EntID, TaskID: item.ID, Title: "开始执行", StatusID: "running"})

	go s.runJobItem(item)
}

func (s *cronSvImpl) runJobItem(item *model.CronTask) {
	defer func() {
		if r := recover(); r != nil {
			s.runJobItemReset(item)
		}
	}()
	tokenContext := item.Context
	if item.ClientID != "" {
		tokenContext.Set("ClientID", item.ClientID)
	}
	if item.UserID != "" {
		tokenContext.Set("UserID", item.UserID)
	}
	if item.EntID != "" {
		tokenContext.Set("EntID", item.EntID)
	}

	client := &http.Client{}
	if item.Endpoint == nil || item.Endpoint.Path == "" || item.Endpoint.Code == "" {
		s.runJobItemFailed(item, log.ErrorD("没有配置可执行接口"))
		return
	}

	postBody := item.Body
	remoteUrl, _ := url.Parse(item.Endpoint.Path)
	req, err := http.NewRequest("POST", remoteUrl.String(), bytes.NewBuffer([]byte(postBody)))
	if err != nil {
		s.runJobItemFailed(item, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenContext.ToTokenString())
	req.Header.Set("JOB_ID", item.ID)
	req.Header.Set("JOB_CODE", item.Code)
	req.Header.Set("USER_ID", tokenContext.UserID())
	req.Header.Set("ENT_ID", tokenContext.EntID())
	resp, err := client.Do(req)
	if err != nil {
		s.runJobItemFailed(item, err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.runJobItemFailed(item, err)
		return
	}
	resBodyObj := resBodyDTO{}
	if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
		log.ErrorD(utils.ToString(resBody))
		s.runJobItemFailed(item, err)
		return
	}
	if resp.StatusCode != 200 || resBodyObj.Msg != "" {
		s.runJobItemFailed(item, errors.New(resBodyObj.Msg, 3020))
		return
	}
	s.runJobItemSucceed(item)
}
func (s *cronSvImpl) runJobItemFailed(item *model.CronTask, err error) {

	s.createLog(&model.CronLog{EntID: item.EntID, TaskID: item.ID, Title: "执行计划出错", StatusID: "failed", Msg: err.Error()})

	updates := make(map[string]interface{})
	item.NumPeriod = item.NumPeriod + 1
	item.NumFailed = item.NumFailed + 1
	updates["NumFailed"] = item.NumFailed
	updates["NumPeriod"] = item.NumPeriod
	updates["LastStatusID"] = "failed"
	updates["LastMsg"] = err.Error()
	if item.UnitID == "once" {
		if item.NumPeriod > item.Retry { //大于重试次数，则停止
			updates["StatusID"] = "completed"
		} else {
			updates["StatusID"] = "waiting"
		}
	} else {
		updates["StatusID"] = "waiting"
		if item.NumPeriod > item.Retry {
			currTime := time.Unix(item.NeTime, 0)
			if item.UnitID == "month" {
				currTime = currTime.AddDate(0, 1, 0)
			} else if item.UnitID == "week" {
				currTime = currTime.AddDate(0, 0, 7)
			} else if item.UnitID == "day" {
				currTime = currTime.AddDate(0, 0, 1)
			} else if item.UnitID == "hour" {
				currTime = currTime.Add(1 * time.Hour)
			} else {
				currTime = currTime.AddDate(1, 0, 0)
			}
			updates["NumPeriod"] = 0
			updates["NeTime"] = currTime.Unix()
		}
	}
	db.Default().Model(item).Where("id=?", item.ID).Updates(updates)
}
func (s *cronSvImpl) runJobItemSucceed(item *model.CronTask) {
	s.createLog(&model.CronLog{EntID: item.EntID, TaskID: item.ID, Title: "执行成功", StatusID: "succeed"})

	updates := make(map[string]interface{})
	item.NumPeriod = item.NumPeriod + 1
	item.NumSuccess = item.NumSuccess + 1

	updates["NumSuccess"] = item.NumSuccess
	updates["NumPeriod"] = item.NumPeriod
	updates["LastStatusID"] = "succeed"
	updates["LastMsg"] = ""
	if item.UnitID == "once" {
		updates["StatusID"] = "completed"
	} else {
		updates["StatusID"] = "waiting"
		currTime := time.Unix(item.NeTime, 0)
		if item.UnitID == "month" {
			currTime = currTime.AddDate(0, 1, 0)
		} else if item.UnitID == "week" {
			currTime = currTime.AddDate(0, 0, 7)
		} else if item.UnitID == "day" {
			currTime = currTime.AddDate(0, 0, 1)
		} else if item.UnitID == "hour" {
			currTime = currTime.Add(1 * time.Hour)
		} else {
			currTime = currTime.AddDate(1, 0, 0)
		}
		updates["NumPeriod"] = 0
		updates["NeTime"] = currTime.Unix()
	}
	db.Default().Model(item).Where("id=?", item.ID).Updates(updates)
}
func (s *cronSvImpl) runJobItemReset(item *model.CronTask) {
	db.Default().Model(item).Where("id=? and status_id=?", item.ID, "running").Updates(map[string]interface{}{"StatusID": "waiting", "LastMsg": "被异常中断"})
	s.createLog(&model.CronLog{EntID: item.EntID, TaskID: item.ID, Title: "被异常中断", StatusID: "completed"})
}

func (s *cronSvImpl) createLog(item *model.CronLog) {
	if item.ID != "" {
		item.ID = utils.GUID()
	}
	db.Default().Create(item)
}

/**
创建一个调度任务
*/
func (s *cronSvImpl) CreateTask(entID string, item *model.CronTask) (*model.CronTask, error) {
	if item.UnitID == "" {
		return nil, errors.ParamsRequired("频率")
	}
	if item.Code == "" {
		item.Code = utils.GUID()
	}
	if item.Cycle < 1 {
		item.Cycle = 1
	}
	if item.FmTime.IsZero() {
		item.FmTime = utils.TimeNow()
	}
	if item.FmTime.Unix() < utils.TimeNow().Unix() {
		item.FmTime = utils.TimeNow()
	}
	//相同的 tag+url+UnitID 只能出现一次
	if item.Tag != "" {
		count := 0
		if db.Default().Model(model.CronTask{}).Where("enabled=1 and status_id in (?) tag=? and unit_id=? and ent_id=?", []string{"waiting", "running"}, item.Tag, item.UnitID, item.EntID).Count(&count); count > 0 {
			return nil, errors.ExistError("频率相同的任务")
		}
	}
	if item.EndpointID == "" && item.Endpoint != nil {
		if item.Endpoint.ID != "" {
			item.EndpointID = item.Endpoint.ID
		} else if item.Endpoint.Code != "" {
			if e, err := s.GetEndpoint(item.Endpoint.Code); err != nil {
				return nil, err
			} else {
				item.EndpointID = e.ID
			}
		}
	}
	if item.EndpointID == "" {
		return nil, log.ErrorD("endpoint为空")
	}
	//下次执行时间
	item.NeTime = item.FmTime.Unix()
	item.EntID = entID
	item.Enabled = utils.SBool_True
	item.StatusID = "waiting"
	item.EntID = entID
	item.ID = utils.GUID()
	if err := db.Default().Create(item).Error; err != nil {
		return nil, err
	}
	return s.GetTaskBy(item.ID)
}

func (s *cronSvImpl) GetEndpoint(id string) (*model.CronEndpoint, error) {
	item := model.CronEndpoint{}
	if err := db.Default().Model(item).Where("id=?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *cronSvImpl) GetTaskBy(id string) (*model.CronTask, error) {
	item := model.CronTask{}
	if err := db.Default().Model(item).Where("id=?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *cronSvImpl) DestroyTask(entID string, ids []string) error {
	systemDatas := 0
	db.Default().Model(model.CronTask{}).Where("id in (?) and `system`=1", ids).Count(&systemDatas)
	if systemDatas > 0 {
		return utils.ToError("系统预制不能删除!")
	}
	if err := db.Default().Delete(model.CronTask{}, "ent_id=? and id in (?)", entID, ids).Error; err != nil {
		return err
	}
	return nil
}
