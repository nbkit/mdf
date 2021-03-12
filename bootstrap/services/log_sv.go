package services

import (
	"github.com/nbkit/mdf/bootstrap/model"
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
)

type ILogSv interface {
	CreateOp(item model.Log)
	CreateLog(item model.Log)
	Create(item model.Log)
	Log(item model.Log)
}
type logSvImpl struct {
}

func LogSv() ILogSv {
	return logSvInstance
}

var logSvInstance ILogSv = newLogSvImpl()

func newLogSvImpl() *logSvImpl {
	return &logSvImpl{}
}
func (s *logSvImpl) CreateOp(item model.Log) {
	item.Type = "op"
	s.Create(item)
}
func (s *logSvImpl) CreateLog(item model.Log) {
	item.Type = "log"
	s.Create(item)
}
func (s *logSvImpl) Create(item model.Log) {
	item.ID = utils.GUID()
	if item.Type == "" {
		item.Type = "log"
	}
	if err := db.Default().Create(&item).Error; err != nil {
		log.Error(err)
	}
	log.Errorf("%s-%s: %v", item.NodeType, item.NodeID, item.Msg)
}

func (s *logSvImpl) Log(item model.Log) {
	log.Errorf("%s-%s: %v", item.NodeType, item.NodeID, item.Msg)
}
