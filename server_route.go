package mdf

import (
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/gin"
	"github.com/nbkit/mdf/utils"
	"net/http"
)

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
	s.engine.POST("md", func(c *gin.Context) {
		md.ActionSv().DoAction(utils.NewFlowContext().Bind(c)).Output()
	})
	s.engine.GET("md/:widget", func(c *gin.Context) {
		widget := c.Param("widget")
		c.HTML(http.StatusOK, "index.html", utils.Map{"app": utils.Config.App.Name, "time": utils.TimeNow().Format(utils.Layout_YYYYMMDDHHIISS), "widget": widget})
	})
}
