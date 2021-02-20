package routes

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/gin"
	"github.com/ggoop/mdf/middleware/token"
	"github.com/ggoop/mdf/utils"
	"net/http"
)

func routeMd(engine *gin.Engine) {
	engine.POST("md", func(c *gin.Context) {
		md.ActionSv().DoAction(token.Get(c), utils.NewReqContext().Bind(c)).Bind(c)
	})
	engine.GET("md/:widget", func(c *gin.Context) {
		widget := c.Param("widget")
		c.HTML(http.StatusOK, "index.html", utils.Map{"title": utils.DefaultConfig.App.Name, "widget": widget})
	})
}
