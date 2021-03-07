package db

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/nbkit/mdf/framework/glog"
	"github.com/nbkit/mdf/utils"

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
func Open() *Repo {
	db, err := gorm.Open(utils.Config.Db.Driver, getDsnString(true))
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
	}

	db.LogMode(utils.Config.App.Debug)
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
	db, err := gorm.Open(utils.Config.Db.Driver, getDsnString(true))
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
		panic(err)
	}
	db.LogMode(utils.Config.App.Debug)
	repo := &Repo{db}
	return repo
}
func getDsnString(inDb bool) string {
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	//mssql:   =>  sqlserver://username:password@localhost:1433?database=dbname
	//mysql => user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
	str := ""
	// 创建连接
	if utils.Config.Db.Driver == utils.ORM_DRIVER_MSSQL {
		var buf bytes.Buffer
		buf.WriteString("sqlserver://")
		buf.WriteString(utils.Config.Db.Username)
		if utils.Config.Db.Password != "" {
			buf.WriteByte(':')
			buf.WriteString(utils.Config.Db.Password)
		}
		buf.WriteByte('@')
		if utils.Config.Db.Host != "" {
			buf.WriteString(utils.Config.Db.Host)
			if utils.Config.Db.Port != "" {
				buf.WriteByte(':')
				buf.WriteString(utils.Config.Db.Port)
			} else {
				buf.WriteString(":1433")
			}
		}
		if utils.Config.Db.Database != "" && inDb {
			buf.WriteString("?database=")
			buf.WriteString(utils.Config.Db.Database)
		} else {
			buf.WriteString("?database=master")
		}
		str = buf.String()
		return str
	}
	{
		config := mysql.Config{
			User:   utils.Config.Db.Username,
			Passwd: utils.Config.Db.Password, Net: "tcp", Addr: utils.Config.Db.Host,
			AllowNativePasswords: true,
			ParseTime:            true,
			Loc:                  time.Local,
		}
		if inDb {
			config.DBName = utils.Config.Db.Database
		}
		if utils.Config.Db.Port != "" {
			config.Addr = fmt.Sprintf("%s:%s", utils.Config.Db.Host, utils.Config.Db.Port)
		}
		str = config.FormatDSN()
	}
	return str
}
func DestroyDB(name string) error {
	db, err := gorm.Open(utils.Config.Db.Driver, getDsnString(false))
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
	}
	defer db.Close()
	return db.Exec(fmt.Sprintf("Drop Database if exists %s;", name)).Error
}
func CreateDB(name string) error {
	db, err := gorm.Open(utils.Config.Db.Driver, getDsnString(false))
	if err != nil {
		return glog.Errorf("orm failed to initialized: %v", err)
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
		return glog.Errorf("create DATABASE err: %v", err)
	}
	return nil
}
