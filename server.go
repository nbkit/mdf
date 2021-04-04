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
	"github.com/nbkit/mdf/middleware/token"
	"github.com/nbkit/mdf/utils"
	"os"
)

type Option struct {
	EnabledFeature bool
}
type Server interface {
	Upgrade() Server
	Cache() Server
	Start()
	Use(func(engine *gin.Engine)) Server
}
type serverImpl struct {
	engine   *gin.Engine
	runArg   string
	option   Option
	entities []interface{}
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
			ser.runArg = os.Args[1]
		}
	}
	ser.initContext()
	return ser
}
func (s *serverImpl) Start() {
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
	if utils.Config.Db.Database != "" && (s.runArg == "upgrade" || s.runArg == "init" || s.runArg == "debug") {
		upgrade.Script(upgrade.ScriptOption{Path: "./storage/script/pre"}).Exec()

		upgrade.Script(upgrade.ScriptOption{Path: "./storage/script/seeds"}).Exec()

		upgrade.Script(upgrade.ScriptOption{Path: "./storage/script/post"}).Exec()
	}
	return s
}
func (s *serverImpl) Use(done func(engine *gin.Engine)) Server {
	if done != nil {
		done(s.engine)
	}
	return s
}
func (s *serverImpl) initContext() {
	if utils.Config.Db.Database != "" {
		db.CreateDB(utils.Config.Db.Database)
		db.Default().DB.DB().SetConnMaxLifetime(0)
	}
	//设置模板
	if utils.PathExists("dist") {
		s.engine.LoadHTMLGlob(utils.JoinCurrentPath("dist/*.html"))
	}
	if s.runArg == "upgrade" || s.runArg == "init" || s.runArg == "debug" {
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

func (s *serverImpl) commonRoute() {
	s.engine.GET("ping", func(c *gin.Context) {
		utils.NewResContext().Set("data", true).Bind(c)
	})
	s.engine.GET("id", func(c *gin.Context) {
		action, _ := c.GetQuery("action")
		if action == "encrypt" {
			str := ""
			if s, ok := c.GetQuery("q"); ok && s != "" {
				str, _ = utils.AesCFBEncrypt(s, utils.Config.App.Token)
			}
			utils.NewResContext().Set("data", str).Bind(c)
		} else {
			ids := make([]string, 0)
			for i := 0; i < 10; i++ {
				ids = append(ids, utils.GUID())
			}
			utils.NewResContext().Set("data", ids).Bind(c)
		}

	})
}

//启动 JOB
func (s *serverImpl) startCron() {
	if utils.Config.GetBool("CRON") {
		go services.CronSv().Start()
	}
}
