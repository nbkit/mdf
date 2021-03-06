package bootstrap

import (
	"fmt"
	"github.com/nbkit/mdf/bootstrap/actions"
	"github.com/nbkit/mdf/bootstrap/model"
	"github.com/nbkit/mdf/bootstrap/routes"
	"github.com/nbkit/mdf/bootstrap/rules"
	"github.com/nbkit/mdf/bootstrap/seeds"
	"github.com/nbkit/mdf/bootstrap/services"
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/framework/reg"
	"github.com/nbkit/mdf/gin"
	"github.com/nbkit/mdf/middleware/token"
	"github.com/nbkit/mdf/utils"
	"os"
	"time"
)

func Start() {

	engine := gin.New()
	runArg := ""
	if os.Args != nil && len(os.Args) > 0 {
		if len(os.Args) > 1 {
			runArg = os.Args[1]
		}
	}

	initContext(engine)
	if runArg == "upgrade" || runArg == "init" {
		model.Register()
		seeds.Register()
	}
	//动作注册
	actions.Register()
	//规则注册
	rules.Register()

	//初始化缓存
	md.ActionSv().Cache()
	md.MDSv().Cache()

	//使用token中间件
	engine.Use(token.Default())
	//设置模板
	utils.CreatePath("dist")
	engine.LoadHTMLGlob("dist/*.html")
	//注册路由
	routes.Register(engine)

	//注册中心
	go reg.StartServer()

	//启动 JOB
	go startCron()

	//启动引擎
	engine.Run(fmt.Sprintf(":%s", utils.DefaultConfig.App.Port))
}
func initContext(engine *gin.Engine) {
	db.CreateDB(utils.DefaultConfig.Db.Database)
}

//启动 JOB
func startCron() {
	time.Sleep(10 * time.Second)
	services.CronSv().Start()
}
