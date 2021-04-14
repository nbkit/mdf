package db

import (
	"fmt"
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/utils"
	"sync"

	"github.com/nbkit/mdf/db/gorm"
)

type Repo struct {
	*gorm.DB
}

var dbIns *Repo
var _cacheDefaultMu sync.Once

func Default() *Repo {
	if dbIns == nil {
		_cacheDefaultMu.Do(func() {
			dbIns = NewMysqlRepo()
		})
	}
	return dbIns
}
func SetDefault(d *Repo) {
	dbIns = d
}
func Open(args ...interface{}) *Repo {
	var (
		dialect string
		source  string
	)
	if len(args) > 0 {
		dialect = utils.Config.Db.Driver
		source = utils.Config.Db.GetDsnString(true)
	} else {
		switch value := args[0].(type) {
		case string:
			if len(args) == 1 {
				dialect = utils.Config.Db.Driver
				source = value
			} else if len(args) >= 2 {
				dialect = value
				source = args[1].(string)
			}
		case utils.DbConfig:
		case *utils.DbConfig:
			dialect = value.Driver
			source = value.GetDsnString(true)
		}
	}
	db, err := gorm.Open(dialect, source)
	if err != nil {
		log.ErrorF("orm failed to initialized: %v", err)
	}
	db.LogMode(utils.Config.Db.Mode == "debug" || utils.Config.Db.Mode == "")
	repo := &Repo{db}
	return repo
}
func (s *Repo) Close() error {
	return s.DB.Close()
}
func (s *Repo) Begin() *Repo {
	return &Repo{s.DB.Begin()}
}
func (s *Repo) New() *Repo {
	return &Repo{s.DB.New()}
}
func NewMysqlRepo() *Repo {
	db, err := gorm.Open(utils.Config.Db.Driver, utils.Config.Db.GetDsnString(true))
	if err != nil {
		log.ErrorF("orm failed to initialized: %v", err)
		panic(err)
	}
	db.LogMode(utils.Config.Db.Mode == "debug" || utils.Config.Db.Mode == "")
	repo := &Repo{db}
	return repo
}
func DestroyDB(name string) error {
	db, err := gorm.Open(utils.Config.Db.Driver, utils.Config.Db.GetDsnString(false))
	if err != nil {
		log.ErrorF("orm failed to initialized: %v", err)
	}
	defer db.Close()
	return db.Exec(fmt.Sprintf("Drop Database if exists %s;", name)).Error
}
func CreateDB(name string) error {
	db, err := gorm.Open(utils.Config.Db.Driver, utils.Config.Db.GetDsnString(false))
	if err != nil {
		return log.ErrorF("orm failed to initialized: %v", err)
	}
	defer db.Close()
	script := ""
	if utils.Config.Db.Driver == utils.ORM_DRIVER_MSSQL {
		script = fmt.Sprintf("if not exists (select * from sysdatabases where name='%s') begin create database %s end;", name, name)
	} else {
		script = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET %s COLLATE %s;", name, utils.Config.Db.Charset, utils.Config.Db.Collation)
	}
	err = db.Exec(script).Error
	if err != nil {
		return log.ErrorF("create DATABASE err: %v", err)
	}
	return nil
}
