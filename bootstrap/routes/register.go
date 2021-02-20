package routes

import "github.com/nbkit/mdf/gin"

func Register(engine *gin.Engine) {
	routeView(engine)
	routeApi(engine)
	routeProxy(engine)
	routeDti(engine)
	routeMd(engine)
}
