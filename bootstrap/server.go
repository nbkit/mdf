package bootstrap

import (
	"fmt"
	"github.com/nbkit/mdf/bootstrap/actions"
	"github.com/nbkit/mdf/bootstrap/model"
	"github.com/nbkit/mdf/bootstrap/routes"
	"github.com/nbkit/mdf/bootstrap/rules"
	"github.com/nbkit/mdf/bootstrap/services"
	"github.com/nbkit/mdf/bootstrap/upgrade"
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/framework/glog"
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/framework/reg"
	"github.com/nbkit/mdf/gin"
	"github.com/nbkit/mdf/middleware/token"
	"github.com/nbkit/mdf/utils"
	"os"
)

type Server interface {
	Start()
}
type serverImpl struct {
	engine *gin.Engine
	runArg string
}

func NewServer() Server {
	return newServer()
}
func newServer() *serverImpl {
	ser := &serverImpl{
		engine: gin.New(),
	}
	if os.Args != nil && len(os.Args) > 0 {
		if len(os.Args) > 1 {
			ser.runArg = os.Args[1]
		}
	}

	ser.initContext()

	return ser
}

func (s *serverImpl) Start() {
	if s.runArg == "upgrade" || s.runArg == "init" {
		upgrade.Script(upgrade.ScriptOption{Path: "./storage/script/pre"})

		upgrade.Script(upgrade.ScriptOption{Path: "./storage/script/seeds"})

		upgrade.Script(upgrade.ScriptOption{Path: "./storage/script/post"})
	}
	//初始化缓存
	md.ActionSv().Cache()
	md.MDSv().Cache()

	//设置模板
	utils.CreatePath("dist")
	s.engine.LoadHTMLGlob("dist/*.html")

	//注册中心
	s.startReg()

	//启动 JOB
	s.startCron()

	//启动引擎
	s.engine.Run(fmt.Sprintf(":%s", utils.Config.App.Port))
}

func (s *serverImpl) initContext() {
	db.CreateDB(utils.Config.Db.Database)
	db.Default().SetLogger(glog.GetLogger(""))
	db.Default().DB.DB().SetConnMaxLifetime(0)

	if s.runArg == "upgrade" || s.runArg == "init" {
		model.Register()
		initSeedAction()
	}
	//动作注册
	actions.Register()
	//规则注册
	rules.Register()
	//使用token中间件
	s.engine.Use(token.Default())
	//注册路由
	routes.Register(s.engine)

}
func (s *serverImpl) startReg() {
	go reg.StartServer()
}

//启动 JOB
func (s *serverImpl) startCron() {
	if utils.Config.GetBool("CRON") {
		go services.CronSv().Start()
	}
}
