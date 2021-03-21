package routes

import (
	"github.com/nbkit/mdf/gin"
	"github.com/nbkit/mdf/utils"
	"net/http"
)

func routeView(engine *gin.Engine) {
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", utils.Map{"title": utils.Config.App.Name})
	})
}
