package services

import (
	"github.com/nbkit/mdf/bootstrap/errors"
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/framework/glog"
	"github.com/robfig/cron"
	"sync"

	"github.com/nbkit/mdf/utils"

	"github.com/nbkit/mdf/bootstrap/model"
)

var m_cronCache *cron.Cron

type ICronSv interface {
}
type cronSvImpl struct {
	*sync.Mutex
}

var cronSv ICronSv = newCronSv()

func newCronSv() *cronSvImpl {
	return &cronSvImpl{Mutex: &sync.Mutex{}}
}
func (s *cronSvImpl) Start() {
	if m_cronCache == nil {
		m_cronCache := cron.New()
		m_cronCache.AddFunc("@every 20s", s.runJob)
		m_cronCache.Start()
	}
}

func (s *cronSvImpl) runJob() {

}

/**
创建一个调度任务
*/
func (s *cronSvImpl) Create(entID string, item *model.CronTask) (*model.CronTask, error) {
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
		return nil, glog.Error("endpoint为空")
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
func (s *cronSvImpl) Destroy(entID string, ids []string) error {
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
