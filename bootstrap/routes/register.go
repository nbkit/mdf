package routes

import "github.com/nbkit/mdf/gin"

func Register(engine *gin.Engine) {
	routeView(engine)
	routeApi(engine)
	routeMd(engine)
}
