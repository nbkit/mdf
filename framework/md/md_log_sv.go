package md

import (
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
)

type ILogSv interface {
	CreateOp(item *MDLog)
	CreateLog(item *MDLog)
	Create(item *MDLog)
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
func (s *logSvImpl) CreateOp(item *MDLog) {
	item.Type = "op"
	s.Create(item)
}
func (s *logSvImpl) CreateLog(item *MDLog) {
	item.Type = "log"
	s.Create(item)
}
func (s *logSvImpl) Create(item *MDLog) {
	item.ID = utils.GUID()
	if item.Type == "" {
		item.Type = "log"
	}
	if err := db.Default().Create(item).Error; err != nil {
		log.ErrorD(err)
	}
	log.InfoF("%s-%s: %v", item.NodeType, item.NodeID, item.Msg)
}
