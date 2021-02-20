package routes

import (
	"github.com/ggoop/mdf/gin"
)

func routeApi(engine *gin.Engine) {
	group := engine.Group("api")
	apiAuth(group.Group("auth"))
}
