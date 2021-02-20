package routes

import (
	"github.com/ggoop/mdf/gin"
	"github.com/ggoop/mdf/utils"
	"net/http"
)

func routeView(engine *gin.Engine) {
	engine.GET("id", func(c *gin.Context) {
		ids := make([]string, 0)
		for i := 0; i < 10; i++ {
			ids = append(ids, utils.GUID())
		}
		utils.NewResContext().Set("data", ids).Bind(c)
	})
	engine.GET("ping", func(c *gin.Context) {
		utils.NewResContext().Set("data", true).Bind(c)
	})
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", utils.Map{"title": utils.DefaultConfig.App.Name})
	})
}
