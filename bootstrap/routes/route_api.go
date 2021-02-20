package routes

import (
	"github.com/nbkit/mdf/gin"
)

func routeApi(engine *gin.Engine) {
	group := engine.Group("api")
	apiAuth(group.Group("auth"))
}
