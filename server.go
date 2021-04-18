package mdf

import (
	"fmt"
	"github.com/nbkit/mdf/bootstrap/actions"
	"github.com/nbkit/mdf/bootstrap/model"
	"github.com/nbkit/mdf/bootstrap/routes"
	"github.com/nbkit/mdf/bootstrap/rules"
	"github.com/nbkit/mdf/bootstrap/services"
	"github.com/nbkit/mdf/bootstrap/upgrade"
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/framework/reg"
	"github.com/nbkit/mdf/gin"
	"github.com/nbkit/mdf/log"
	"github.com/nbkit/mdf/middleware/token"
	"github.com/nbkit/mdf/utils"
	"os"
	"path/filepath"
)

type Option struct {
	EnabledFeature bool
}
type Server interface {
	Upgrade() Server
	Cache() Server
	Start()
	Use(func(server Server)) Server
	GetEngine() *gin.Engine
	IsMigrate() bool
	GetRunArg() []string
}
type serverImpl struct {
	engine    *gin.Engine
	runArg    string
	option    Option
	isMigrate bool
	entities  []interface{}
}

func NewServer(options ...Option) Server {
	return newServer(options...)
}
func newServer(options ...Option) *serverImpl {
	option := Option{EnabledFeature: false}
	if len(options) > 0 {
		option = options[0]
	}
	utils.Config.SetValue("EnabledFeature", option.EnabledFeature)
	gin.SetMode(utils.Config.App.Mode)
	gin.ForceConsoleColor()
	ser := &serverImpl{
		engine:   gin.New(),
		option:   option,
		entities: make([]interface{}, 0),
	}
	if os.Args != nil && len(os.Args) > 0 {
		if len(os.Args) > 1 {
			ser.runArg = os.Args[len(os.Args)-1]
		}
	}
	if ser.runArg == "migrate" || ser.runArg == "upgrade" || ser.runArg == "init" || ser.runArg == "debug" {
		ser.isMigrate = true
	}
	ser.initContext()
	return ser
}

var initArgs = []string{"install", "uninstall"}

func (s *serverImpl) Start() {
	if utils.StringsContains(initArgs, s.runArg) > -1 {
		return
	}
	//注册中心
	s.startReg()
	//启动 JOB
	s.startCron()
	//启动引擎
	s.engine.Run(fmt.Sprintf(":%s", utils.Config.App.Port))
}
func (s *serverImpl) Cache() Server {
	md.MDSv().Cache()
	return s
}
func (s *serverImpl) Upgrade() Server {
	if utils.Config.Db.Database != "" && s.isMigrate {
		upgrade.Script(upgrade.ScriptOption{Path: "./storage/script/pre"}).Exec()

		upgrade.Script(upgrade.ScriptOption{Path: "./storage/script/seeds"}).Exec()

		upgrade.Script(upgrade.ScriptOption{Path: "./storage/script/post"}).Exec()
	}
	return s
}
func (s *serverImpl) Use(done func(server Server)) Server {
	if done != nil {
		done(s)
	}
	return s
}

func (s *serverImpl) GetEngine() *gin.Engine {
	return s.engine
}
func (s *serverImpl) GetRunArg() []string {
	return []string{s.runArg}
}
func (s *serverImpl) IsMigrate() bool {
	return s.isMigrate
}
func (s *serverImpl) initContext() {
	if utils.Config.Db.Database != "" {
		db.CreateDB(utils.Config.Db.Database)
		db.Default().DB.DB().SetConnMaxLifetime(0)
	}
	//设置模板
	if utils.PathExists("storage/template") {
		pattern := utils.JoinCurrentPath("storage/template/*.html")
		if filenames, err := filepath.Glob(pattern); err != nil {
			log.Error().Error(err)
		} else if len(filenames) > 0 {
			s.engine.LoadHTMLGlob(pattern)
		}
	}
	if s.isMigrate {
		md.MDSv().Migrate()
		model.Register()
		initSeedAction()
	}
	//动作注册
	actions.Register()
	//规则 注册
	rules.Register()
	//使用token中间件
	s.engine.Use(token.Default())

	// 日志输出
	s.engine.Use(gin.Logger())

	if s.option.EnabledFeature {
		//注册路由
		routes.Register(s.engine)
	}
	// 通用路由
	s.commonRoute()
	//缓存
	s.Cache()
}
func (s *serverImpl) startReg() {
	if s.option.EnabledFeature {
		go reg.StartServer()
	}
}

//启动 JOB
func (s *serverImpl) startCron() {
	if utils.Config.GetBool("CRON") {
		go services.CronSv().Start()
	}
}
