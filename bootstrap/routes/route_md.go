package routes

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/gin"
	"github.com/nbkit/mdf/utils"
	"net/http"
)

func routeMd(engine *gin.Engine) {
	engine.POST("md", func(c *gin.Context) {
		md.ActionSv().DoAction(utils.NewFlowContext().Bind(c)).Output()
	})
	engine.GET("md/:widget", func(c *gin.Context) {
		widget := c.Param("widget")
		c.HTML(http.StatusOK, "index.html", utils.Map{"title": utils.Config.App.Name, "widget": widget})
	})
}
