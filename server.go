package mdf

import (
	"fmt"
	"github.com/nbkit/mdf/bootstrap/actions"
	"github.com/nbkit/mdf/bootstrap/model"
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
	"path"
	"path/filepath"
)

type Config struct {
}

type Server struct {
	runArg  string
	option  *Option
	engine  *gin.Engine
	useDone []func(server *Server)
}

func NewServer(cfg Config) *Server {
	s := newServer(cfg)
	return s
}
func newServer(cfg Config) *Server {
	option := newOption()
	gin.SetMode(utils.Config.App.Mode)
	gin.ForceConsoleColor()
	ser := &Server{
		engine:  gin.New(),
		option:  option,
		useDone: make([]func(server *Server), 0),
	}
	if os.Args != nil && len(os.Args) > 0 {
		if len(os.Args) > 1 {
			ser.runArg = os.Args[len(os.Args)-1]
		}
	}
	if ser.runArg == "migrate" || ser.runArg == "upgrade" || ser.runArg == "init" || ser.runArg == "debug" {
		ser.option.isMigrate = true
	}
	ser.initHtmlTemplate()

	return ser
}

var initArgs = []string{"install", "uninstall"}

func (s *Server) Start(o ...func(*Option)) {
	//初始启动参数
	s.initOption(o...)
	// 日志输出
	s.engine.Use(gin.Logger())
	//数据库迁移
	if s.option.isMigrate {
		s.initMigrate()
	}
	//动作注册,规则 注册
	if s.option.enableMDF {
		actions.Register()
		rules.Register()
	}
	// 处理升级
	if s.option.isUpgrade {
		s.upgrade()
	}
	//使用token中间件
	if s.option.enableAuthToken {
		s.engine.Use(token.Default())
	}
	// 初始化路由
	s.initRoute()
	//执行中间插件
	s.initDone()
	//如果是安装或者卸载，则不需要启动和执行后边逻辑
	if utils.StringsContains(initArgs, s.runArg) > -1 {
		return
	}
	// 初始化缓存
	s.initCache()
	// 启动注册中心
	if s.option.isRegistry {
		s.startReg()
	}
	// 启动JOB
	if s.option.enableCron {
		s.startCron()
	}
	//启动引擎
	s.engine.Run(fmt.Sprintf(":%s", utils.Config.App.Port))
}

func (s *Server) Use(done func(server *Server)) *Server {
	s.useDone = append(s.useDone, done)
	return s
}
func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}
func (s *Server) GetRunArg() []string {
	return []string{s.runArg}
}
func (s *Server) IsMigrate() bool {
	return s.option.isMigrate
}

func (s *Server) initOption(o ...func(*Option)) {
	for _, f := range o {
		f(s.option)
	}
}
func (s *Server) initDone() {
	for _, f := range s.useDone {
		if f != nil {
			f(s)
		}
	}
}
func (s *Server) initHtmlTemplate() {
	viewPath := utils.Config.GetValue("view.path")
	if viewPath == "" {
		viewPath = "./storage/template"
	}
	isBinary := utils.Config.GetBool("view.binary")
	//设置模板
	if utils.PathExists(viewPath) {
		if isBinary {

		}
		pattern := utils.JoinCurrentPath(path.Join(viewPath, "*.html"))
		if filenames, err := filepath.Glob(pattern); err != nil {
			log.Error().Error(err)
		} else if len(filenames) > 0 {
			s.engine.LoadHTMLGlob(pattern)
		}
	}
}
func (s *Server) initMigrate() {
	if utils.Config.Db.Database != "" && s.option.isMigrate {
		db.CreateDB(utils.Config.Db.Database)
		db.Default().DB.DB().SetConnMaxLifetime(0)
	}
	md.MDSv().Migrate()
	initSeedAction()

	if s.option.isBaseDataCenter {
		model.Register()
	}
}
func (s *Server) initCache() *Server {
	if utils.Config.Db.Database != "" && s.option.enableMDF {
		md.MDSv().Cache()
	}
	return s
}
func (s *Server) upgrade() *Server {
	if utils.Config.Db.Database != "" && s.option.isMigrate {
		upgrade.Script(upgrade.ScriptOption{Path: "./storage/script/pre"}).Exec()

		upgrade.Script(upgrade.ScriptOption{Path: "./storage/script/seeds"}).Exec()

		upgrade.Script(upgrade.ScriptOption{Path: "./storage/script/post"}).Exec()
	}
	return s
}
func (s *Server) startReg() {
	go reg.StartServer()

}

//启动 JOB
func (s *Server) startCron() {
	go services.CronSv().Start()

}
